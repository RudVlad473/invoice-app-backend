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

type UpdateInvoiceItemDTO struct {
	/*
		We don't validate it here because we should only validate an item that we're adding
	*/
	Items []Item `json:"items,omitempty" dynamodbav:"items,omitempty"`
}
