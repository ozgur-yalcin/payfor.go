[![license](https://img.shields.io/:license-mit-blue.svg)](https://github.com/ozgur-soft/payfor.go/blob/main/LICENSE.md)
[![documentation](https://pkg.go.dev/badge/github.com/ozgur-soft/payfor.go)](https://pkg.go.dev/github.com/ozgur-soft/payfor.go/src)

# Payfor.go
QNB Finansbank (PayFor) POS API with golang

# Installation
```bash
go get github.com/ozgur-yalcin/payfor.go
```

# Satış
```go
package main

import (
	"context"
	"encoding/xml"
	"fmt"

	payfor "github.com/ozgur-yalcin/payfor.go/src"
)

// Pos bilgileri
const (
	envmode  = "TEST"                    // Çalışma ortamı (Production : "PROD" - Test : "TEST")
	mbr      = "5"                       // Kurum kodu
	merchant = "085300000009704"         // İşyeri numarası
	username = "QNB_API_KULLANICI_3DPAY" // Kullanıcı adı
	password = "UcBN0"                   // Şifre
	lang     = "TR"                      // Dil
)

func main() {
	api, req := payfor.Api(mbr, merchant, username, password)
	api.SetMode(envmode)

	req.SetCardHolder("")                 // Kart sahibi (zorunlu)
	req.SetCardNumber("4155650100416111") // Kart numarası (zorunlu)
	req.SetCardExpiry("01", "25")         // Son kullanma tarihi - AA,YY (zorunlu)
	req.SetCardCode("123")                // Kart arkasındaki 3 haneli numara (zorunlu)
	req.SetAmount("1.00", "TRY")          // Satış tutarı ve para birimi (zorunlu)
	req.SetLang(lang)

	// Satış
	ctx := context.Background()
	if res, err := api.Auth(ctx, req); err == nil {
		pretty, _ := xml.MarshalIndent(res, " ", " ")
		fmt.Println(string(pretty))
	} else {
		fmt.Println(err)
	}
}
```

# İade
```go
package main

import (
	"context"
	"encoding/xml"
	"fmt"

	payfor "github.com/ozgur-yalcin/payfor.go/src"
)

// Pos bilgileri
const (
	envmode  = "TEST"                    // Çalışma ortamı (Production : "PROD" - Test : "TEST")
	mbr      = "5"                       // Kurum kodu
	merchant = "085300000009704"         // İşyeri numarası
	username = "QNB_API_KULLANICI_3DPAY" // Kullanıcı adı
	password = "UcBN0"                   // Şifre
	lang     = "TR"                      // Dil
)

func main() {
	api, req := payfor.Api(mbr, merchant, username, password)
	api.SetMode(envmode)

	req.SetOrgOrderId("SYS_")    // Sipariş numarası (zorunlu)
	req.SetAmount("1.00", "TRY") // İade tutarı ve para birimi (zorunlu)
	req.SetLang(lang)

	// İade
	ctx := context.Background()
	if res, err := api.Refund(ctx, req); err == nil {
		pretty, _ := xml.MarshalIndent(res, " ", " ")
		fmt.Println(string(pretty))
	} else {
		fmt.Println(err)
	}
}
```

# İptal
```go
package main

import (
	"context"
	"encoding/xml"
	"fmt"

	payfor "github.com/ozgur-yalcin/payfor.go/src"
)

// Pos bilgileri
const (
	envmode  = "TEST"                    // Çalışma ortamı (Production : "PROD" - Test : "TEST")
	mbr      = "5"                       // Kurum kodu
	merchant = "085300000009704"         // İşyeri numarası
	username = "QNB_API_KULLANICI_3DPAY" // Kullanıcı adı
	password = "UcBN0"                   // Şifre
	lang     = "TR"                      // Dil
)

func main() {
	api, req := payfor.Api(mbr, merchant, username, password)
	api.SetMode(envmode)

	req.SetOrgOrderId("SYS_") // Sipariş numarası (zorunlu)
	req.SetCurrency("TRY")    // Para birimi (zorunlu)
	req.SetLang(lang)

	// İptal
	ctx := context.Background()
	if res, err := api.Cancel(ctx, req); err == nil {
		pretty, _ := xml.MarshalIndent(res, " ", " ")
		fmt.Println(string(pretty))
	} else {
		fmt.Println(err)
	}
}
```
