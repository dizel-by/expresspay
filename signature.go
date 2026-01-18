package expresspay

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

type SignatureFunc func(action string, params map[string]string, secret string) (string, error)

var signatureOrder = map[string][]string{
	"add-invoice": {
		"token",
		"accountno",
		"amount",
		"currency",
		"expiration",
		"info",
		"surname",
		"firstname",
		"patronymic",
		"city",
		"street",
		"house",
		"building",
		"apartment",
		"isnameeditable",
		"isaddresseditable",
		"isamounteditable",
	},
	"get-details-invoice": {
		"token",
		"id",
	},
	"cancel-invoice": {
		"token",
		"id",
	},
	"status-invoice": {
		"token",
		"invoiceid",
	},
	"get-list-invoices": {
		"token",
		"from",
		"to",
		"accountno",
		"status",
	},
	"get-list-payments": {
		"token",
		"from",
		"to",
		"accountno",
	},
	"get-details-payment": {
		"token",
		"id",
	},
	"add-card-invoice": {
		"token",
		"accountno",
		"expiration",
		"amount",
		"currency",
		"info",
		"returnurl",
		"failurl",
		"language",
		"sessiontimeoutsecs",
		"expirationdate",
		"returninvoiceurl",
	},
	"card-invoice-form": {
		"token",
		"cardinvoiceno",
	},
	"status-card-invoice": {
		"token",
		"cardinvoiceno",
		"language",
	},
	"reverse-card-invoice": {
		"token",
		"cardinvoiceno",
	},
	"get-qr-code": {
		"token",
		"invoiceid",
		"viewtype",
		"imagewidth",
		"imageheight",
	},
	"add-web-invoice": {
		"token",
		"serviceid",
		"accountno",
		"amount",
		"currency",
		"expiration",
		"info",
		"surname",
		"firstname",
		"patronymic",
		"city",
		"street",
		"house",
		"building",
		"apartment",
		"isnameeditable",
		"isaddresseditable",
		"isamounteditable",
		"emailnotification",
		"smsphone",
		"returntype",
		"returnurl",
		"failurl",
	},
	"add-webcard-invoice": {
		"token",
		"serviceid",
		"accountno",
		"expiration",
		"amount",
		"currency",
		"info",
		"returnurl",
		"failurl",
		"language",
		"sessiontimeoutsecs",
		"expirationdate",
		"returntype",
	},
}

func DefaultSignature(action string, params map[string]string, secret string) (string, error) {
	order, ok := signatureOrder[action]
	if !ok {
		return "", fmt.Errorf("unknown signature action: %s", action)
	}

	normalized := make(map[string]string, len(params))
	for k, v := range params {
		normalized[strings.ToLower(k)] = v
	}

	var builder strings.Builder
	for _, key := range order {
		if value, ok := normalized[key]; ok {
			builder.WriteString(value)
		}
	}

	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(builder.String()))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil))), nil
}

func parseAPIError(data []byte) *APIError {
	var envelope struct {
		Error        *APIError   `json:"Error"`
		ErrorCode    json.Number `json:"ErrorCode"`
		ErrorMessage string      `json:"ErrorMessage"`
	}
	if err := decodeJSON(data, &envelope); err != nil {
		return nil
	}
	if envelope.Error != nil {
		return envelope.Error
	}
	if envelope.ErrorMessage != "" {
		return &APIError{ErrorCode: intFromNumber(envelope.ErrorCode), ErrorMessage: envelope.ErrorMessage}
	}
	return nil
}
