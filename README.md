# Unofficial Express Pay Belarus Go client

Go client for the "Express Pay" API based on https://express-pay.by/docs/api/v1

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	"expresspay"
)

func main() {
	client := expresspay.NewClient(expresspay.SandboxBaseURL, "API_TOKEN", "SECRET", expresspay.WithSignature())

	invoice, err := client.CreateInvoice(context.Background(), expresspay.AddInvoiceRequest{
		AccountNo: "123456",
		Amount:    "10,00",
		Currency:  expresspay.CurrencyBYN,
		Info:      "test",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("InvoiceNo:", invoice.InvoiceNo.String())
}
```

## Notes

- `DefaultBaseURL` uses production, `SandboxBaseURL` points to the test stand.
- Signature generation follows the parameter order described in the documentation.
- Set `WithSignature()` if you need a signature even with an empty secret.
