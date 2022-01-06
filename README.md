[![license](https://img.shields.io/:license-mit-blue.svg)](https://github.com/ozgur-soft/qnbpay.go/blob/main/LICENSE.md)
[![documentation](https://pkg.go.dev/badge/github.com/ozgur-soft/qnbpay.go)](https://pkg.go.dev/github.com/ozgur-soft/qnbpay.go/src)

# Qnbpay.go
Qnb Finansbank Sanal POS API with golang

# Installation
```bash
go get github.com/ozgur-soft/qnbpay.go
```

# Sanalpos satış işlemi

# Sanalpos iade işlemi
```go
package main

import (
	"context"
	"encoding/xml"
	"fmt"

	qnbpay "github.com/ozgur-soft/qnbpay.go/src"
)

func main() {
	api, req := qnbpay.Api("5", "MerchantID ", "Usercode", "Userpass")
	// Test : "TEST" - Production "PROD" (zorunlu)
	api.SetMode("PROD")
	// Sipariş numarası (zorunlu)
	req.SetOrgOrderId("SYS_")
	// İade tutarı (zorunlu)
	req.SetAmount("1.00")
	// Para birimi (zorunlu)
	req.SetCurrency("TRY")
	// Dil (zorunlu)
	req.SetLang("TR")

	// İade
	ctx := context.Background()
	res := api.Refund(ctx, req)
	pretty, _ := xml.MarshalIndent(res, " ", " ")
	fmt.Println(string(pretty))
}
```

# Sanalpos iptal işlemi
```go
package main

import (
	"context"
	"encoding/xml"
	"fmt"

	qnbpay "github.com/ozgur-soft/qnbpay.go/src"
)

func main() {
	api, req := qnbpay.Api("5", "MerchantID ", "Usercode", "Userpass")
	// Test : "TEST" - Production "PROD" (zorunlu)
	api.SetMode("PROD")
	// Sipariş numarası (zorunlu)
	req.SetOrgOrderId("SYS_")
	// Para birimi (zorunlu)
	req.SetCurrency("TRY")
	// Dil (zorunlu)
	req.SetLang("TR")

	// İptal
	ctx := context.Background()
	res := api.Cancel(ctx, req)
	pretty, _ := xml.MarshalIndent(res, " ", " ")
	fmt.Println(string(pretty))
}
```
