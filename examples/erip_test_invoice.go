package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dizel-by/expresspay"
)

func main() {
	client := expresspay.NewClient(
		expresspay.DefaultBaseURL,
		"YOUR_TOKEN",
		"",
	)

	resp, err := client.CreateInvoice(context.Background(), expresspay.AddInvoiceRequest{
		AccountNo:         "1",
		Amount:            "0,01",
		Currency:          "933",
		Expiration:        time.Now().AddDate(0, 0, 1).Format("20060102"),
		Info:              "Test ERIP invoice",
		IsNameEditable:    "0",
		IsAddressEditable: "0",
		IsAmountEditable:  "0",
		ReturnInvoiceURL:  "1",
		EmailNotification: "email@example.org",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("InvoiceNo: %s\n", resp.InvoiceNo.String())
	if resp.InvoiceURL != "" {
		fmt.Printf("InvoiceURL: %s\n", resp.InvoiceURL)
	}
}
