package adyen

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// DefaultCurrency default currency for transactions
const DefaultCurrency = "EUR"

// New - creates Adyen instance
func New(username, password, clientID, merchantAccount string) *Adyen {
	return &Adyen{
		Username:        username,
		Password:        password,
		ClientID:        clientID,
		MerchantAccount: merchantAccount,
		Currency:        DefaultCurrency,
	}
}

// Adyen - base structure with configuration options
type Adyen struct {
	Username        string
	Password        string
	ClientID        string
	MerchantAccount string
	Currency        string
	Logger          *log.Logger
}

// Version of a current Adyen API
const (
	APIVersion         = "v25"
	AdyenTestURL       = "https://pal-test.adyen.com/pal/servlet/Payment"
	AdyenClientTestURL = "https://test.adyen.com/hpp/cse/js/"
)

// ClientURL - returns URl, that need to loaded in UI, to encrypt Credit Card information
func (a *Adyen) ClientURL() string {
	return AdyenClientTestURL + a.ClientID + ".shtml"
}

// AdyenURL returns Adyen backend URL
func (a *Adyen) AdyenURL(requestType string) string {
	return AdyenTestURL + "/" + APIVersion + "/" + requestType
}

// AttachLogger attach logger to API instance
func (a *Adyen) AttachLogger(Logger *log.Logger) {
	a.Logger = Logger
}

// SetCurrency set default currency for transactions
func (a *Adyen) SetCurrency(currency string) {
	a.Currency = currency
}

// execute request on Adyen side, transforms "requestEntity" into JSON representation
func (a *Adyen) execute(method string, requestEntity interface{}) (*Response, error) {
	body, err := json.Marshal(requestEntity)
	if err != nil {
		return nil, err
	}

	url := a.AdyenURL(method)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if a.Logger != nil {
		a.Logger.Printf("[Request]: %s %s\n%s", method, url, body)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(a.Username, a.Password)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	if a.Logger != nil {
		a.Logger.Printf("[Response]: %s %s\n%s", method, url, buf.String())
	}

	providerResponse := &Response{
		Response: resp,
		Body:     buf.Bytes(),
	}

	err = providerResponse.handleHTTPError()

	if err != nil {
		return nil, err
	}

	return providerResponse, nil
}

// Payment - returns PaymentGateway
func (a *Adyen) Payment() *PaymentGateway {
	return &PaymentGateway{a}
}

// Modification - returns ModificationGateway
func (a *Adyen) Modification() *ModificationGateway {
	return &ModificationGateway{a}
}
