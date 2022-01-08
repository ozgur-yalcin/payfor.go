package qnbpay

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"log"
	"net/http"
	"strings"
)

var EndPoints map[string]string = map[string]string{
	"TEST": "https://vpostest.qnbfinansbank.com/Gateway/XmlGate.aspx",
	"PROD": "https://vpos.qnbfinansbank.com/Gateway/XmlGate.aspx",
}

var Currencies map[string]string = map[string]string{
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
	XMLName    xml.Name    `xml:"PayforRequest,omitempty"`
	MbrId      interface{} `xml:"MbrId,omitempty"`
	MerchantId interface{} `xml:"MerchantId,omitempty"`
	UserCode   interface{} `xml:"UserCode,omitempty"`
	UserPass   interface{} `xml:"UserPass,omitempty"`
	SecureType interface{} `xml:"SecureType,omitempty"`
	TxnType    interface{} `xml:"TxnType,omitempty"`
	Amount     interface{} `xml:"PurchAmount,omitempty"`
	Currency   interface{} `xml:"Currency,omitempty"`
	Instalment interface{} `xml:"InstallmentCount,omitempty"`
	CardHolder interface{} `xml:"CardHolderName,omitempty"`
	CardNumber interface{} `xml:"Pan,omitempty"`
	CardExpiry interface{} `xml:"Expiry,omitempty"`
	CardCode   interface{} `xml:"Cvv2,omitempty"`
	OrderId    interface{} `xml:"OrderId,omitempty"`
	OrgOrderId interface{} `xml:"OrgOrderId,omitempty"`
	OkUrl      interface{} `xml:"OkUrl,omitempty"`
	FailUrl    interface{} `xml:"FailUrl,omitempty"`
	Rnd        interface{} `xml:"Rnd,omitempty"`
	Hash       interface{} `xml:"Hash,omitempty"`
	MOTO       interface{} `xml:"MOTO,omitempty"`
	Lang       interface{} `xml:"Lang,omitempty"`
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

func SHA1(data string) (hash string) {
	h := sha1.New()
	h.Write([]byte(data))
	hash = hex.EncodeToString(h.Sum(nil))
	return hash
}

func Api(mbrid, merchantid, usercode, userpass string) (*API, *Request) {
	api := new(API)
	request := new(Request)
	request.MbrId = mbrid
	request.MerchantId = merchantid
	request.UserCode = usercode
	request.UserPass = userpass
	return api, request
}

func (api *API) SetMode(mode string) {
	api.Mode = mode
}

func (api *API) SetKey(key string) {
	api.Key = key
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

func (request *Request) SetAmount(total string) {
	request.Amount = total
}

func (request *Request) SetInstalment(ins string) {
	request.Instalment = ins
}

func (request *Request) SetCurrency(currency string) {
	request.Currency = Currencies[currency]
}

func (request *Request) SetOrderId(oid string) {
	request.OrderId = oid
}

func (request *Request) SetOrgOrderId(orgid string) {
	request.OrgOrderId = orgid
}

func (request *Request) SetMoto(moto string) {
	request.MOTO = moto
}

func (request *Request) SetLang(lang string) {
	request.Lang = lang
}

func (api *API) Pay(ctx context.Context, req *Request) Response {
	req.TxnType = "Auth"
	req.SecureType = "NonSecure"
	if req.MOTO == nil {
		req.MOTO = "0"
	}
	return api.Transaction(ctx, req)
}

func (api *API) Refund(ctx context.Context, req *Request) Response {
	req.TxnType = "Refund"
	req.SecureType = "NonSecure"
	return api.Transaction(ctx, req)
}

func (api *API) Cancel(ctx context.Context, req *Request) Response {
	req.TxnType = "Void"
	req.SecureType = "NonSecure"
	return api.Transaction(ctx, req)
}

func (api *API) Transaction(ctx context.Context, req *Request) (res Response) {
	postdata, err := xml.Marshal(req)
	if err != nil {
		log.Println(err)
		return res
	}
	request, err := http.NewRequestWithContext(ctx, "POST", EndPoints[api.Mode], strings.NewReader(xml.Header+string(postdata)))
	if err != nil {
		log.Println(err)
		return res
	}
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return res
	}
	defer response.Body.Close()
	decoder := xml.NewDecoder(response.Body)
	decoder.Decode(&res)
	return res
}
