package expresspay

import "encoding/json"

const (
	CurrencyBYN = "933"
	CurrencyEUR = "978"
	CurrencyUSD = "840"
	CurrencyRUB = "643"
)

const (
	InvoiceStatusPendingPayment   = "1"
	InvoiceStatusExpired          = "2"
	InvoiceStatusPaid             = "3"
	InvoiceStatusPartiallyPaid    = "4"
	InvoiceStatusCanceled         = "5"
	InvoiceStatusPaidByBankCard   = "6"
	InvoiceStatusPaymentReturned  = "7"
)

const (
	QRCodeViewTypeBase64 = "base64"
	QRCodeViewTypeText   = "text"
)

type APIError struct {
	Code         int    `json:"Code"`
	Msg          string `json:"Msg"`
	MsgCode      int    `json:"MsgCode"`
	ErrorCode    int    `json:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage"`
	HTTPStatus   int    `json:"-"`
	Raw          string `json:"-"`
}

func (e *APIError) Error() string {
	if e == nil {
		return "api error"
	}
	if e.ErrorCode != 0 || e.ErrorMessage != "" {
		return e.ErrorMessage
	}
	if e.Msg != "" {
		return e.Msg
	}
	if e.Raw != "" {
		return e.Raw
	}
	return "api error"
}

type ListInvoicesParams struct {
	From      string
	To        string
	AccountNo string
	Status    string
}

type AddInvoiceRequest struct {
	AccountNo         string
	Amount            string
	Currency          string
	Expiration        string
	Info              string
	Surname           string
	FirstName         string
	Patronymic        string
	City              string
	Street            string
	House             string
	Building          string
	Apartment         string
	IsNameEditable    string
	IsAddressEditable string
	IsAmountEditable  string
	EmailNotification string
	SmsPhone          string
	ReturnInvoiceURL  string
}

type AddInvoiceResponse struct {
	InvoiceNo  json.Number `json:"InvoiceNo"`
	InvoiceURL string      `json:"InvoiceUrl"`
}

type Invoice struct {
	InvoiceNo     json.Number `json:"InvoiceNo"`
	AccountNo     string      `json:"AccountNo"`
	Status        json.Number `json:"Status"`
	Created       string      `json:"Created"`
	Expiration    string      `json:"Expiration"`
	Amount        json.Number `json:"Amount"`
	Currency      json.Number `json:"Currency"`
	CardInvoiceNo json.Number `json:"CardInvoiceNo"`
}

type InvoiceDetails struct {
	Status            json.Number `json:"Status"`
	Created           string      `json:"Created"`
	Expiration        string      `json:"Expiration"`
	Amount            json.Number `json:"Amount"`
	Currency          json.Number `json:"Currency"`
	Info              string      `json:"Info"`
	Surname           string      `json:"Surname"`
	FirstName         string      `json:"FirstName"`
	Patronymic        string      `json:"Patronymic"`
	City              string      `json:"City"`
	Street            string      `json:"Street"`
	House             string      `json:"House"`
	Building          string      `json:"Building"`
	Apartment         string      `json:"Apartment"`
	IsNameEditable    json.Number `json:"IsNameEditable"`
	IsAddressEditable json.Number `json:"IsAddressEditable"`
	IsAmountEditable  json.Number `json:"IsAmountEditable"`
}

type InvoiceStatusResponse struct {
	Status json.Number `json:"Status"`
}

type ListPaymentsParams struct {
	From      string
	To        string
	AccountNo string
}

type Payment struct {
	PaymentNo  json.Number `json:"PaymentNo"`
	AccountNo  string      `json:"AccountNo"`
	Created    string      `json:"Created"`
	Amount     json.Number `json:"Amount"`
	Currency   json.Number `json:"Currency"`
	Info       string      `json:"Info"`
	Surname    string      `json:"Surname"`
	FirstName  string      `json:"FirstName"`
	Patronymic string      `json:"Patronymic"`
	City       string      `json:"City"`
	Street     string      `json:"Street"`
	House      string      `json:"House"`
	Building   string      `json:"Building"`
	Apartment  string      `json:"Apartment"`
}

type PaymentDetails struct {
	AccountNo  string      `json:"AccountNo"`
	Created    string      `json:"Created"`
	Amount     json.Number `json:"Amount"`
	Currency   json.Number `json:"Currency"`
	Info       string      `json:"Info"`
	Surname    string      `json:"Surname"`
	FirstName  string      `json:"FirstName"`
	Patronymic string      `json:"Patronymic"`
	City       string      `json:"City"`
	Street     string      `json:"Street"`
	House      string      `json:"House"`
	Building   string      `json:"Building"`
	Apartment  string      `json:"Apartment"`
}

type QRCodeParams struct {
	ViewType    string
	ImageWidth  string
	ImageHeight string
}

type QRCodeResponse struct {
	QrCodeBody string `json:"QrCodeBody"`
}

type AddCardInvoiceRequest struct {
	AccountNo          string
	Expiration         string
	Amount             string
	Currency           string
	Info               string
	ReturnURL          string
	FailURL            string
	Language           string
	SessionTimeoutSecs string
	ExpirationDate     string
	ReturnInvoiceURL   string
}

type AddCardInvoiceResponse struct {
	CardInvoiceNo json.Number `json:"CardInvoiceNo"`
	InvoiceURL    string      `json:"InvoiceUrl"`
	ErrorCode     json.Number `json:"ErrorCode"`
	ErrorMessage  string      `json:"ErrorMessage"`
}

func (r AddCardInvoiceResponse) check() error {
	if intFromNumber(r.ErrorCode) != 0 || r.ErrorMessage != "" {
		return &APIError{ErrorCode: intFromNumber(r.ErrorCode), ErrorMessage: r.ErrorMessage}
	}
	return nil
}

type CardInvoiceFormResponse struct {
	FormURL      string      `json:"FormUrl"`
	ErrorCode    json.Number `json:"ErrorCode"`
	ErrorMessage string      `json:"ErrorMessage"`
}

func (r CardInvoiceFormResponse) check() error {
	if intFromNumber(r.ErrorCode) != 0 || r.ErrorMessage != "" {
		return &APIError{ErrorCode: intFromNumber(r.ErrorCode), ErrorMessage: r.ErrorMessage}
	}
	return nil
}

type CardInvoiceStatusResponse struct {
	Amount            json.Number `json:"Amount"`
	CardInvoiceStatus json.Number `json:"CardInvoiceStatus"`
	ErrorCode         json.Number `json:"ErrorCode"`
	ErrorMessage      string      `json:"ErrorMessage"`
}

func (r CardInvoiceStatusResponse) check() error {
	if intFromNumber(r.ErrorCode) != 0 || r.ErrorMessage != "" {
		return &APIError{ErrorCode: intFromNumber(r.ErrorCode), ErrorMessage: r.ErrorMessage}
	}
	return nil
}

type CardInvoiceReverseResponse struct {
	ErrorCode    json.Number `json:"ErrorCode"`
	ErrorMessage string      `json:"ErrorMessage"`
}

func (r CardInvoiceReverseResponse) check() error {
	if intFromNumber(r.ErrorCode) != 0 || r.ErrorMessage != "" {
		return &APIError{ErrorCode: intFromNumber(r.ErrorCode), ErrorMessage: r.ErrorMessage}
	}
	return nil
}

type AddWebInvoiceRequest struct {
	ServiceID         string
	AccountNo         string
	Amount            string
	Currency          string
	ReturnType        string
	ReturnURL         string
	FailURL           string
	Expiration        string
	Info              string
	Surname           string
	FirstName         string
	Patronymic        string
	City              string
	Street            string
	House             string
	Building          string
	Apartment         string
	IsNameEditable    string
	IsAddressEditable string
	IsAmountEditable  string
	EmailNotification string
	SmsPhone          string
	ReturnInvoiceURL  string
}

type AddWebInvoiceResponse struct {
	InvoiceNo               json.Number `json:"InvoiceNo"`
	InvoiceURL              string      `json:"InvoiceUrl"`
	ExpressPayAccountNumber string      `json:"ExpressPayAccountNumber"`
	ExpressPayInvoiceNo     json.Number `json:"ExpressPayInvoiceNo"`
	Signature               string      `json:"Signature"`
}

type AddWebCardInvoiceRequest struct {
	ServiceID          string
	AccountNo          string
	Expiration         string
	Amount             string
	Currency           string
	Info               string
	ReturnType         string
	ReturnURL          string
	FailURL            string
	Language           string
	SessionTimeoutSecs string
	ExpirationDate     string
	ReturnInvoiceURL   string
}

type AddWebCardInvoiceResponse struct {
	FormURL                 string      `json:"FormUrl"`
	InvoiceURL              string      `json:"InvoiceUrl"`
	ExpressPayAccountNumber string      `json:"ExpressPayAccountNumber"`
	ExpressPayInvoiceNo     json.Number `json:"ExpressPayInvoiceNo"`
	Signature               string      `json:"Signature"`
}

func intFromNumber(n json.Number) int {
	if n == "" {
		return 0
	}
	v, _ := n.Int64()
	return int(v)
}
