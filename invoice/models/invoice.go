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

// BaseItem
// /*Should act as baseline for dtos & entities (what we put in db)*/
type BaseItem struct {
	Name     string  `json:"name,omitempty" dynamodbav:"name,omitempty" validate:"required,min=1,max=100"`
	Quantity int     `json:"quantity,omitempty" dynamodbav:"quantity,omitempty" validate:"required,min=1"`
	Price    float64 `json:"price,omitempty" dynamodbav:"price,omitempty" validate:"required,gt=0"`
}

type Item struct {
	Id string `json:"id,omitempty" dynamodbav:"id,omitempty" validate:"required,uuid"`
	BaseItem
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
