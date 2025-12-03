package invoice

import (
	"time"
)

type SaveAddressDTO struct {
	Street   string `json:"street,omitempty" dynamodbav:"street,omitempty" validate:"required,min=1,max=200"`
	City     string `json:"city,omitempty" dynamodbav:"city,omitempty" validate:"required,min=1,max=100"`
	PostCode string `json:"postCode,omitempty" dynamodbav:"postCode,omitempty" validate:"required,min=1,max=20"`
	Country  string `json:"country,omitempty" dynamodbav:"country,omitempty" validate:"required,min=1,max=100"`
}

type SaveItemDTO struct {
	Name     string  `json:"name,omitempty" dynamodbav:"name,omitempty" validate:"required,min=1,max=100"`
	Quantity int     `json:"quantity,omitempty" dynamodbav:"quantity,omitempty" validate:"required,min=1"`
	Price    float64 `json:"price,omitempty" dynamodbav:"price,omitempty" validate:"required,gt=0"`
}

type SaveInvoiceDTO struct {
	PaymentDue    time.Time       `json:"paymentDue,omitempty" dynamodbav:"paymentDue,omitempty" validate:"required"`
	Description   string          `json:"description,omitempty" dynamodbav:"description,omitempty" validate:"required,min=1,max=500"`
	ClientName    string          `json:"clientName,omitempty" dynamodbav:"clientName,omitempty" validate:"required,min=1,max=100"`
	ClientEmail   string          `json:"clientEmail,omitempty" dynamodbav:"clientEmail,omitempty" validate:"required,email"`
	SenderAddress *SaveAddressDTO `json:"senderAddress,omitempty" dynamodbav:"senderAddress,omitempty" validate:"required"`
	ClientAddress *SaveAddressDTO `json:"clientAddress,omitempty" dynamodbav:"clientAddress,omitempty" validate:"required"`
	Items         []SaveItemDTO   `json:"items,omitempty" dynamodbav:"items,omitempty" validate:"required,min=1,dive"`
}

type UpdateInvoiceDTO struct {
	PaymentDue  time.Time `json:"paymentDue,omitempty" dynamodbav:"paymentDue,omitempty" validate:"required"`
	Description string    `json:"description,omitempty" dynamodbav:"description,omitempty" validate:"required,min=1,max=500"`
	ClientName  string    `json:"clientName,omitempty" dynamodbav:"clientName,omitempty" validate:"required,min=1,max=100"`
	ClientEmail string    `json:"clientEmail,omitempty" dynamodbav:"clientEmail,omitempty" validate:"required,email"`
	/*
		We need to have these as pointers so we can automatically omit nil values
		Auto-omit logic doesn't work as great with non-pointer structs
	*/
	SenderAddress *SaveAddressDTO `json:"senderAddress,omitempty" dynamodbav:"senderAddress,omitempty" validate:"required"`
	ClientAddress *SaveAddressDTO `json:"clientAddress,omitempty" dynamodbav:"clientAddress,omitempty" validate:"required"`
}
