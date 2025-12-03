package invoice

import (
	"time"

	invoice "github.com/rudvlad473/invoice-app-backend/invoice/constants"
)

type Address struct {
	Street   string `json:"street" dynamodbav:"street"`
	City     string `json:"city" dynamodbav:"city"`
	PostCode string `json:"postCode" dynamodbav:"postCode"`
	Country  string `json:"country" dynamodbav:"country"`
}

type Item struct {
	Name     string  `json:"name" dynamodbav:"name"`
	Quantity int     `json:"quantity" dynamodbav:"quantity"`
	Price    float64 `json:"price" dynamodbav:"price"`
}

type Invoice struct {
	Id            string         `json:"id" dynamodbav:"id"`
	CreatedAt     time.Time      `json:"createdAt" dynamodbav:"createdAt"`
	PaymentDue    time.Time      `json:"paymentDue" dynamodbav:"paymentDue"`
	Description   string         `json:"description" dynamodbav:"description"`
	ClientName    string         `json:"clientName" dynamodbav:"clientName"`
	ClientEmail   string         `json:"clientEmail" dynamodbav:"clientEmail"`
	Status        invoice.Status `json:"status" dynamodbav:"status"`
	SenderAddress Address        `json:"senderAddress" dynamodbav:"senderAddress"`
	ClientAddress Address        `json:"clientAddress" dynamodbav:"clientAddress"`
	Items         []Item         `json:"items" dynamodbav:"items"`
}
