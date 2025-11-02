package invoice

import (
	"time"

	invoice "github.com/rudvlad473/invoice-app-backend/invoice/constants"
)

type Address struct {
	Street   string `json:"street"`
	City     string `json:"city"`
	PostCode string `json:"postCode"`
	Country  string `json:"country"`
}

type Item struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Total    float64 `json:"total"`
}

type Invoice struct {
	Id            string         `json:"id"`
	CreatedAt     time.Time      `json:"createdAt"`
	PaymentDue    time.Time      `json:"paymentDue"`
	Description   string         `json:"description"`
	PaymentTerms  int            `json:"paymentTerms"`
	ClientName    string         `json:"clientName"`
	ClientEmail   string         `json:"clientEmail"`
	Status        invoice.Status `json:"status"`
	SenderAddress Address        `json:"senderAddress"`
	ClientAddress Address        `json:"clientAddress"`
	Items         []Item         `json:"items"`
	Total         float64        `json:"total"`
}
