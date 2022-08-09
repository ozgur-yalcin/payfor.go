package payfor

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var EndPoints = map[string]string{
	"TEST": "https://vpostest.qnbfinansbank.com/Gateway/XmlGate.aspx",
	"PROD": "https://vpos.qnbfinansbank.com/Gateway/XmlGate.aspx",

	"TEST3D": "https://vpostest.qnbfinansbank.com/Gateway/Default.aspx",
	"PROD3D": "https://vpos.qnbfinansbank.com/Gateway/Default.aspx",
}

var CurrencyCode = map[string]string{
	"TRY": "949",
	"YTL": "949",
	"TRL": "949",
	"TL":  "949",
	"USD": "840",
	"EUR": "978",
	"GBP": "826",
	"JPY": "392",
}

type API struct {
	Mode string
	Key  string
}

type Request struct {
	XMLName         xml.Name `xml:"PayforRequest,omitempty"`
	RequestGuid     string   `xml:"RequestGuid,omitempty"`
	MbrId           string   `xml:"MbrId,omitempty" form:"MbrId,omitempty"`
	MerchantId      string   `xml:"MerchantId,omitempty" form:"MerchantId,omitempty"`
	UserCode        string   `xml:"UserCode,omitempty" form:"UserCode,omitempty"`
	UserPass        string   `xml:"UserPass,omitempty"`
	OrgOrderId      string   `xml:"OrgOrderId,omitempty"`
	OrderId         string   `xml:"OrderId,omitempty" form:"OrderId,omitempty"`
	TransactionType string   `xml:"TxnType,omitempty" form:"TxnType,omitempty"`
	Amount          string   `xml:"PurchAmount,omitempty" form:"PurchAmount,omitempty"`
	Currency        string   `xml:"Currency,omitempty" form:"Currency,omitempty"`
	Installment     string   `xml:"InstallmentCount,omitempty" form:"InstallmentCount,omitempty"`
	CardHolder      string   `xml:"CardHolderName,omitempty" form:"CardHolderName,omitempty"`
	CardNumber      string   `xml:"Pan,omitempty" form:"Pan,omitempty"`
	CardExpiry      string   `xml:"Expiry,omitempty" form:"Expiry,omitempty"`
	CardCode        string   `xml:"Cvv2,omitempty" form:"Cvv2,omitempty"`
	SecureType      string   `xml:"SecureType,omitempty" form:"SecureType,omitempty"`
	OkUrl           string   `xml:"OkUrl,omitempty" form:"OkUrl,omitempty"`
	FailUrl         string   `xml:"FailUrl,omitempty" form:"FailUrl,omitempty"`
	Random          string   `xml:"Rnd,omitempty" form:"Rnd,omitempty"`
	Hash            string   `xml:"Hash,omitempty" form:"Hash,omitempty"`
	MOTO            string   `xml:"MOTO,omitempty" form:"MOTO,omitempty"`
	Lang            string   `xml:"Lang,omitempty" form:"Lang,omitempty"`
}

type Response struct {
	XMLName        xml.Name `xml:"PayforResponse,omitempty"`
	OrderId        string   `xml:"OrderId,omitempty"`
	TransId        string   `xml:"TransId,omitempty"`
	Response       string   `xml:"Response,omitempty"`
	AuthCode       string   `xml:"AuthCode,omitempty"`
	HostRefNum     string   `xml:"HostRefNum,omitempty"`
	ProcReturnCode string   `xml:"ProcReturnCode,omitempty"`
	Status3D       string   `xml:"3DStatus,omitempty"`
	ResponseRnd    string   `xml:"ResponseRnd,omitempty"`
	ResponseHash   string   `xml:"ResponseHash,omitempty"`
	TxnResult      string   `xml:"TxnResult,omitempty"`
	ErrMsg         string   `xml:"ErrMsg,omitempty"`
}

func Random(n int) string {
	const alphanum = "123456789"
	var bytes = make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func HEX(data string) (hash string) {
	b, err := hex.DecodeString(data)
	if err != nil {
		log.Println(err)
		return hash
	}
	hash = string(b)
	return hash
}

func SHA1(data string) (hash string) {
	h := sha1.New()
	h.Write([]byte(data))
	hash = hex.EncodeToString(h.Sum(nil))
	return hash
}

func B64(data string) (hash string) {
	hash = base64.StdEncoding.EncodeToString([]byte(data))
	return hash
}

func D64(data string) []byte {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Println(err)
		return nil
	}
	return b
}

func Hash(data string) string {
	return B64(HEX(SHA1(data)))
}

func Api(mbr, merchant, username, password string) (*API, *Request) {
	api := new(API)
	request := new(Request)
	request.MbrId = mbr
	request.MerchantId = merchant
	request.UserCode = username
	request.UserPass = password
	return api, request
}

func (api *API) SetStoreKey(key string) {
	api.Key = key
}

func (api *API) SetMode(mode string) {
	api.Mode = mode
}

func (request *Request) SetCardHolder(holder string) {
	request.CardHolder = holder
}

func (request *Request) SetCardNumber(number string) {
	request.CardNumber = number
}

func (request *Request) SetCardExpiry(month, year string) {
	request.CardExpiry = month + year
}

func (request *Request) SetCardCode(code string) {
	request.CardCode = code
}

func (request *Request) SetAmount(total, currency string) {
	request.Amount = total
	request.Currency = CurrencyCode[currency]
}

func (request *Request) SetCurrency(currency string) {
	request.Currency = CurrencyCode[currency]
}

func (request *Request) SetInstallment(ins string) {
	request.Installment = ins
}

func (request *Request) SetOrderId(oid string) {
	request.OrderId = oid
}

func (request *Request) SetOrgOrderId(orgid string) {
	request.OrgOrderId = orgid
}

func (request *Request) SetLang(lang string) {
	request.Lang = lang
}

func (api *API) PreAuth(ctx context.Context, req *Request) (Response, error) {
	req.TransactionType = "PreAuth"
	req.SecureType = "NonSecure"
	if req.MOTO == "" {
		req.MOTO = "0"
	}
	return api.Transaction(ctx, req)
}

func (api *API) Auth(ctx context.Context, req *Request) (Response, error) {
	req.TransactionType = "Auth"
	req.SecureType = "NonSecure"
	if req.MOTO == "" {
		req.MOTO = "0"
	}
	return api.Transaction(ctx, req)
}

func (api *API) PreAuth3D(ctx context.Context, req *Request) (Response, error) {
	req.MbrId = ""
	req.MerchantId = ""
	req.SecureType = "3DModelPayment"
	return api.Transaction(ctx, req)
}

func (api *API) Auth3D(ctx context.Context, req *Request) (Response, error) {
	req.MbrId = ""
	req.MerchantId = ""
	req.SecureType = "3DModelPayment"
	return api.Transaction(ctx, req)
}

func (api *API) PreAuth3Dhtml(ctx context.Context, req *Request) (string, error) {
	req.TransactionType = "PreAuth"
	req.SecureType = "3DModel"
	req.Random = Random(6)
	req.Hash = Hash(req.MbrId + req.OrderId + req.Amount + req.OkUrl + req.FailUrl + req.TransactionType + req.Installment + req.Random + api.Key)
	return api.Transaction3D(ctx, req)
}

func (api *API) Auth3Dhtml(ctx context.Context, req *Request) (string, error) {
	req.TransactionType = "Auth"
	req.SecureType = "3DModel"
	req.Random = Random(6)
	req.Hash = Hash(req.MbrId + req.OrderId + req.Amount + req.OkUrl + req.FailUrl + req.TransactionType + req.Installment + req.Random + api.Key)
	return api.Transaction3D(ctx, req)
}

func (api *API) PostAuth(ctx context.Context, req *Request) (Response, error) {
	req.TransactionType = "PostAuth"
	req.SecureType = "NonSecure"
	if req.MOTO == "" {
		req.MOTO = "0"
	}
	return api.Transaction(ctx, req)
}

func (api *API) Refund(ctx context.Context, req *Request) (Response, error) {
	req.TransactionType = "Refund"
	req.SecureType = "NonSecure"
	return api.Transaction(ctx, req)
}

func (api *API) Cancel(ctx context.Context, req *Request) (Response, error) {
	req.TransactionType = "Void"
	req.SecureType = "NonSecure"
	return api.Transaction(ctx, req)
}

func (api *API) Transaction(ctx context.Context, req *Request) (res Response, err error) {
	postdata, err := xml.Marshal(req)
	if err != nil {
		return res, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", EndPoints[api.Mode], strings.NewReader(xml.Header+string(postdata)))
	if err != nil {
		return res, err
	}
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	decoder := xml.NewDecoder(response.Body)
	if err := decoder.Decode(&res); err != nil {
		return res, err
	}
	switch res.ProcReturnCode {
	case "00":
		return res, nil
	default:
		return res, errors.New(res.ErrMsg)
	}
}

func (api *API) Transaction3D(ctx context.Context, req *Request) (res string, err error) {
	postdata, err := QueryString(req)
	if err != nil {
		return res, err
	}
	html := []string{}
	html = append(html, `<!DOCTYPE html>`)
	html = append(html, `<html>`)
	html = append(html, `<head>`)
	html = append(html, `<script type="text/javascript">function submitonload() {document.payment.submit();document.getElementById('button').remove();document.getElementById('body').insertAdjacentHTML("beforeend", "Lütfen bekleyiniz...");}</script>`)
	html = append(html, `</head>`)
	html = append(html, `<body onload="javascript:submitonload();" id="body" style="text-align:center;margin:10px;font-family:Arial;font-weight:bold;">`)
	html = append(html, `<form action="`+EndPoints[api.Mode+"3D"]+`" method="post" name="payment">`)
	for k := range postdata {
		html = append(html, `<input type="hidden" name="`+k+`" value="`+postdata.Get(k)+`">`)
	}
	html = append(html, `<input type="submit" value="Gönder" id="button">`)
	html = append(html, `</form>`)
	html = append(html, `</body>`)
	html = append(html, `</html>`)
	res = B64(strings.Join(html, "\n"))
	return res, err
}
