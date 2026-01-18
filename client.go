package expresspay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	DefaultBaseURL = "https://api.express-pay.by/v1/"
	SandboxBaseURL = "https://sandbox-api.express-pay.by/v1/"
)

type Client struct {
	BaseURL       string
	Token         string
	Secret        string
	HTTPClient    *http.Client
	UseSignature  bool
	SignatureFunc SignatureFunc
}

type Option func(*Client)

func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) {
		if h != nil {
			c.HTTPClient = h
		}
	}
}

func WithSignature() Option {
	return func(c *Client) {
		c.UseSignature = true
	}
}

func WithSignatureFunc(fn SignatureFunc) Option {
	return func(c *Client) {
		if fn != nil {
			c.SignatureFunc = fn
		}
	}
}

func NewClient(baseURL, token, secret string, opts ...Option) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	c := &Client{
		BaseURL:       baseURL,
		Token:         token,
		Secret:        secret,
		HTTPClient:    http.DefaultClient,
		UseSignature:  secret != "",
		SignatureFunc: DefaultSignature,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) ListInvoices(ctx context.Context, p ListInvoicesParams) ([]Invoice, error) {
	query := url.Values{}
	query.Set("token", c.Token)
	addIfNotEmpty(query, "From", p.From)
	addIfNotEmpty(query, "To", p.To)
	addIfNotEmpty(query, "AccountNo", p.AccountNo)
	addIfNotEmpty(query, "Status", p.Status)

	sigParams := map[string]string{
		"Token":     c.Token,
		"From":      p.From,
		"To":        p.To,
		"AccountNo": p.AccountNo,
		"Status":    p.Status,
	}
	if err := c.applySignature("get-list-invoices", sigParams, query, nil, true); err != nil {
		return nil, err
	}

	data, err := c.do(ctx, http.MethodGet, "invoices", query, nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Items []Invoice `json:"Items"`
	}
	if err := decodeJSON(data, &resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (c *Client) CreateInvoice(ctx context.Context, r AddInvoiceRequest) (*AddInvoiceResponse, error) {
	query := url.Values{}
	query.Set("token", c.Token)
	form := url.Values{}
	form.Set("AccountNo", r.AccountNo)
	form.Set("Amount", r.Amount)
	form.Set("Currency", r.Currency)
	addIfNotEmpty(form, "Expiration", r.Expiration)
	addIfNotEmpty(form, "Info", r.Info)
	addIfNotEmpty(form, "Surname", r.Surname)
	addIfNotEmpty(form, "FirstName", r.FirstName)
	addIfNotEmpty(form, "Patronymic", r.Patronymic)
	addIfNotEmpty(form, "City", r.City)
	addIfNotEmpty(form, "Street", r.Street)
	addIfNotEmpty(form, "House", r.House)
	addIfNotEmpty(form, "Building", r.Building)
	addIfNotEmpty(form, "Apartment", r.Apartment)
	addIfNotEmpty(form, "IsNameEditable", r.IsNameEditable)
	addIfNotEmpty(form, "IsAddressEditable", r.IsAddressEditable)
	addIfNotEmpty(form, "IsAmountEditable", r.IsAmountEditable)
	addIfNotEmpty(form, "EmailNotification", r.EmailNotification)
	addIfNotEmpty(form, "SmsPhone", r.SmsPhone)
	addIfNotEmpty(form, "ReturnInvoiceUrl", r.ReturnInvoiceURL)

	sigParams := map[string]string{
		"Token":             c.Token,
		"AccountNo":         r.AccountNo,
		"Amount":            r.Amount,
		"Currency":          r.Currency,
		"Expiration":        r.Expiration,
		"Info":              r.Info,
		"Surname":           r.Surname,
		"FirstName":         r.FirstName,
		"Patronymic":        r.Patronymic,
		"City":              r.City,
		"Street":            r.Street,
		"House":             r.House,
		"Building":          r.Building,
		"Apartment":         r.Apartment,
		"IsNameEditable":    r.IsNameEditable,
		"IsAddressEditable": r.IsAddressEditable,
		"IsAmountEditable":  r.IsAmountEditable,
		"EmailNotification": r.EmailNotification,
		"SmsPhone":          r.SmsPhone,
		"ReturnInvoiceUrl":  r.ReturnInvoiceURL,
	}
	if err := c.applySignature("add-invoice", sigParams, query, nil, true); err != nil {
		return nil, err
	}

	data, err := c.do(ctx, http.MethodPost, "invoices", query, form)
	if err != nil {
		return nil, err
	}
	var resp AddInvoiceResponse
	if err := decodeJSON(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetInvoice(ctx context.Context, invoiceNo int) (*InvoiceDetails, error) {
	query := url.Values{}
	query.Set("token", c.Token)

	sigParams := map[string]string{
		"Token": c.Token,
		"Id":    fmt.Sprintf("%d", invoiceNo),
	}
	if err := c.applySignature("get-details-invoice", sigParams, query, nil, true); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("invoices/%d", invoiceNo)
	data, err := c.do(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, err
	}
	var resp InvoiceDetails
	if err := decodeJSON(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetInvoiceStatus(ctx context.Context, invoiceNo int) (*InvoiceStatusResponse, error) {
	query := url.Values{}
	query.Set("token", c.Token)

	sigParams := map[string]string{
		"Token":     c.Token,
		"InvoiceId": fmt.Sprintf("%d", invoiceNo),
	}
	if err := c.applySignature("status-invoice", sigParams, query, nil, true); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("invoices/%d/status", invoiceNo)
	data, err := c.do(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, err
	}
	var resp InvoiceStatusResponse
	if err := decodeJSON(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CancelInvoice(ctx context.Context, invoiceNo int) error {
	query := url.Values{}
	query.Set("token", c.Token)

	sigParams := map[string]string{
		"Token": c.Token,
		"Id":    fmt.Sprintf("%d", invoiceNo),
	}
	if err := c.applySignature("cancel-invoice", sigParams, query, nil, true); err != nil {
		return err
	}

	path := fmt.Sprintf("invoices/%d", invoiceNo)
	_, err := c.do(ctx, http.MethodDelete, path, query, nil)
	return err
}

func (c *Client) ListPayments(ctx context.Context, p ListPaymentsParams) ([]Payment, error) {
	query := url.Values{}
	query.Set("token", c.Token)
	addIfNotEmpty(query, "From", p.From)
	addIfNotEmpty(query, "To", p.To)
	addIfNotEmpty(query, "AccountNo", p.AccountNo)

	sigParams := map[string]string{
		"Token":     c.Token,
		"From":      p.From,
		"To":        p.To,
		"AccountNo": p.AccountNo,
	}
	if err := c.applySignature("get-list-payments", sigParams, query, nil, true); err != nil {
		return nil, err
	}

	data, err := c.do(ctx, http.MethodGet, "payments", query, nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Items []Payment `json:"Items"`
	}
	if err := decodeJSON(data, &resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (c *Client) GetPayment(ctx context.Context, paymentNo int) (*PaymentDetails, error) {
	query := url.Values{}
	query.Set("token", c.Token)

	sigParams := map[string]string{
		"Token": c.Token,
		"Id":    fmt.Sprintf("%d", paymentNo),
	}
	if err := c.applySignature("get-details-payment", sigParams, query, nil, true); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("payments/%d", paymentNo)
	data, err := c.do(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, err
	}
	var resp PaymentDetails
	if err := decodeJSON(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetQRCode(ctx context.Context, invoiceID int, p QRCodeParams) (*QRCodeResponse, error) {
	query := url.Values{}
	query.Set("token", c.Token)
	query.Set("id", fmt.Sprintf("%d", invoiceID))
	addIfNotEmpty(query, "ViewType", p.ViewType)
	addIfNotEmpty(query, "ImageWidth", p.ImageWidth)
	addIfNotEmpty(query, "ImageHeight", p.ImageHeight)

	sigParams := map[string]string{
		"Token":       c.Token,
		"InvoiceId":   fmt.Sprintf("%d", invoiceID),
		"ViewType":    p.ViewType,
		"ImageWidth":  p.ImageWidth,
		"ImageHeight": p.ImageHeight,
	}
	if err := c.applySignature("get-qr-code", sigParams, query, nil, true); err != nil {
		return nil, err
	}

	data, err := c.do(ctx, http.MethodGet, "qrcode/getqrcode", query, nil)
	if err != nil {
		return nil, err
	}
	var resp QRCodeResponse
	if err := decodeJSON(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateCardInvoice(ctx context.Context, r AddCardInvoiceRequest) (*AddCardInvoiceResponse, error) {
	query := url.Values{}
	query.Set("token", c.Token)
	form := url.Values{}
	form.Set("AccountNo", r.AccountNo)
	form.Set("Amount", r.Amount)
	form.Set("Currency", r.Currency)
	form.Set("Info", r.Info)
	form.Set("ReturnUrl", r.ReturnURL)
	form.Set("FailUrl", r.FailURL)
	addIfNotEmpty(form, "Expiration", r.Expiration)
	addIfNotEmpty(form, "Language", r.Language)
	addIfNotEmpty(form, "SessionTimeoutSecs", r.SessionTimeoutSecs)
	addIfNotEmpty(form, "ExpirationDate", r.ExpirationDate)
	addIfNotEmpty(form, "ReturnInvoiceUrl", r.ReturnInvoiceURL)

	sigParams := map[string]string{
		"Token":              c.Token,
		"AccountNo":          r.AccountNo,
		"Expiration":         r.Expiration,
		"Amount":             r.Amount,
		"Currency":           r.Currency,
		"Info":               r.Info,
		"ReturnUrl":          r.ReturnURL,
		"FailUrl":            r.FailURL,
		"Language":           r.Language,
		"SessionTimeoutSecs": r.SessionTimeoutSecs,
		"ExpirationDate":     r.ExpirationDate,
		"ReturnInvoiceUrl":   r.ReturnInvoiceURL,
	}
	if err := c.applySignature("add-card-invoice", sigParams, query, nil, true); err != nil {
		return nil, err
	}

	data, err := c.do(ctx, http.MethodPost, "cardinvoices", query, form)
	if err != nil {
		return nil, err
	}
	var resp AddCardInvoiceResponse
	if err := decodeJSON(data, &resp); err != nil {
		return nil, err
	}
	if err := resp.check(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetCardInvoicePaymentForm(ctx context.Context, cardInvoiceNo int) (*CardInvoiceFormResponse, error) {
	query := url.Values{}
	query.Set("token", c.Token)

	sigParams := map[string]string{
		"Token":         c.Token,
		"CardInvoiceNo": fmt.Sprintf("%d", cardInvoiceNo),
	}
	if err := c.applySignature("card-invoice-form", sigParams, query, nil, true); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("cardinvoices/%d/payment", cardInvoiceNo)
	data, err := c.do(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, err
	}
	var resp CardInvoiceFormResponse
	if err := decodeJSON(data, &resp); err != nil {
		return nil, err
	}
	if err := resp.check(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetCardInvoiceStatus(ctx context.Context, cardInvoiceNo int, language string) (*CardInvoiceStatusResponse, error) {
	query := url.Values{}
	query.Set("token", c.Token)
	addIfNotEmpty(query, "Language", language)

	sigParams := map[string]string{
		"Token":         c.Token,
		"CardInvoiceNo": fmt.Sprintf("%d", cardInvoiceNo),
		"Language":      language,
	}
	if err := c.applySignature("status-card-invoice", sigParams, query, nil, true); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("cardinvoices/%d/status", cardInvoiceNo)
	data, err := c.do(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, err
	}
	var resp CardInvoiceStatusResponse
	if err := decodeJSON(data, &resp); err != nil {
		return nil, err
	}
	if err := resp.check(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) ReverseCardInvoice(ctx context.Context, cardInvoiceNo int) (*CardInvoiceReverseResponse, error) {
	query := url.Values{}
	query.Set("token", c.Token)
	form := url.Values{}
	form.Set("Token", c.Token)
	form.Set("CardInvoiceNo", fmt.Sprintf("%d", cardInvoiceNo))

	sigParams := map[string]string{
		"Token":         c.Token,
		"CardInvoiceNo": fmt.Sprintf("%d", cardInvoiceNo),
	}
	if err := c.applySignature("reverse-card-invoice", sigParams, query, nil, true); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("cardinvoices/%d/reverse", cardInvoiceNo)
	data, err := c.do(ctx, http.MethodPost, path, query, form)
	if err != nil {
		return nil, err
	}
	var resp CardInvoiceReverseResponse
	if err := decodeJSON(data, &resp); err != nil {
		return nil, err
	}
	if err := resp.check(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateWebInvoice(ctx context.Context, r AddWebInvoiceRequest) (*AddWebInvoiceResponse, error) {
	form := url.Values{}
	form.Set("ServiceId", r.ServiceID)
	form.Set("AccountNo", r.AccountNo)
	form.Set("Amount", r.Amount)
	form.Set("Currency", r.Currency)
	form.Set("ReturnType", r.ReturnType)
	form.Set("ReturnUrl", r.ReturnURL)
	form.Set("FailUrl", r.FailURL)
	addIfNotEmpty(form, "Expiration", r.Expiration)
	addIfNotEmpty(form, "Info", r.Info)
	addIfNotEmpty(form, "Surname", r.Surname)
	addIfNotEmpty(form, "FirstName", r.FirstName)
	addIfNotEmpty(form, "Patronymic", r.Patronymic)
	addIfNotEmpty(form, "City", r.City)
	addIfNotEmpty(form, "Street", r.Street)
	addIfNotEmpty(form, "House", r.House)
	addIfNotEmpty(form, "Building", r.Building)
	addIfNotEmpty(form, "Apartment", r.Apartment)
	addIfNotEmpty(form, "IsNameEditable", r.IsNameEditable)
	addIfNotEmpty(form, "IsAddressEditable", r.IsAddressEditable)
	addIfNotEmpty(form, "IsAmountEditable", r.IsAmountEditable)
	addIfNotEmpty(form, "EmailNotification", r.EmailNotification)
	addIfNotEmpty(form, "SmsPhone", r.SmsPhone)
	addIfNotEmpty(form, "ReturnInvoiceUrl", r.ReturnInvoiceURL)

	sigParams := map[string]string{
		"Token":             c.Token,
		"ServiceId":         r.ServiceID,
		"AccountNo":         r.AccountNo,
		"Amount":            r.Amount,
		"Currency":          r.Currency,
		"Expiration":        r.Expiration,
		"Info":              r.Info,
		"Surname":           r.Surname,
		"FirstName":         r.FirstName,
		"Patronymic":        r.Patronymic,
		"City":              r.City,
		"Street":            r.Street,
		"House":             r.House,
		"Building":          r.Building,
		"Apartment":         r.Apartment,
		"IsNameEditable":    r.IsNameEditable,
		"IsAddressEditable": r.IsAddressEditable,
		"IsAmountEditable":  r.IsAmountEditable,
		"EmailNotification": r.EmailNotification,
		"SmsPhone":          r.SmsPhone,
		"ReturnType":        r.ReturnType,
		"ReturnUrl":         r.ReturnURL,
		"FailUrl":           r.FailURL,
		"ReturnInvoiceUrl":  r.ReturnInvoiceURL,
	}
	if err := c.applySignature("add-web-invoice", sigParams, nil, form, false); err != nil {
		return nil, err
	}

	data, err := c.do(ctx, http.MethodPost, "web_invoices", nil, form)
	if err != nil {
		return nil, err
	}
	var resp AddWebInvoiceResponse
	if err := decodeJSON(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateWebCardInvoice(ctx context.Context, r AddWebCardInvoiceRequest) (*AddWebCardInvoiceResponse, error) {
	form := url.Values{}
	form.Set("ServiceId", r.ServiceID)
	form.Set("AccountNo", r.AccountNo)
	form.Set("Amount", r.Amount)
	form.Set("Currency", r.Currency)
	form.Set("Info", r.Info)
	form.Set("ReturnType", r.ReturnType)
	form.Set("ReturnUrl", r.ReturnURL)
	form.Set("FailUrl", r.FailURL)
	addIfNotEmpty(form, "Expiration", r.Expiration)
	addIfNotEmpty(form, "Language", r.Language)
	addIfNotEmpty(form, "SessionTimeoutSecs", r.SessionTimeoutSecs)
	addIfNotEmpty(form, "ExpirationDate", r.ExpirationDate)
	addIfNotEmpty(form, "ReturnInvoiceUrl", r.ReturnInvoiceURL)

	sigParams := map[string]string{
		"Token":              c.Token,
		"ServiceId":          r.ServiceID,
		"AccountNo":          r.AccountNo,
		"Expiration":         r.Expiration,
		"Amount":             r.Amount,
		"Currency":           r.Currency,
		"Info":               r.Info,
		"ReturnUrl":          r.ReturnURL,
		"FailUrl":            r.FailURL,
		"Language":           r.Language,
		"SessionTimeoutSecs": r.SessionTimeoutSecs,
		"ExpirationDate":     r.ExpirationDate,
		"ReturnType":         r.ReturnType,
		"ReturnInvoiceUrl":   r.ReturnInvoiceURL,
	}
	if err := c.applySignature("add-webcard-invoice", sigParams, nil, form, false); err != nil {
		return nil, err
	}

	data, err := c.do(ctx, http.MethodPost, "web_cardinvoices", nil, form)
	if err != nil {
		return nil, err
	}
	var resp AddWebCardInvoiceResponse
	if err := decodeJSON(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) applySignature(action string, params map[string]string, query url.Values, form url.Values, inQuery bool) error {
	if !c.UseSignature && c.Secret == "" {
		return nil
	}
	if c.SignatureFunc == nil {
		return fmt.Errorf("signature function is not configured")
	}
	sig, err := c.SignatureFunc(action, params, c.Secret)
	if err != nil {
		return err
	}
	if inQuery {
		if query != nil {
			query.Set("signature", sig)
		}
		return nil
	}
	if form != nil {
		form.Set("Signature", sig)
	}
	return nil
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, form url.Values) ([]byte, error) {
	fullURL := c.BaseURL + path
	if query != nil && len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	var body io.Reader
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete {
		if form != nil {
			body = strings.NewReader(form.Encode())
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, err
	}
	if form != nil && (method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if apiErr := parseAPIError(data); apiErr != nil {
		apiErr.HTTPStatus = resp.StatusCode
		return nil, apiErr
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return data, nil
}

func decodeJSON(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	return dec.Decode(v)
}

func addIfNotEmpty(v url.Values, key, value string) {
	if value == "" {
		return
	}
	v.Set(key, value)
}
