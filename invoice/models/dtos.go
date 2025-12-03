package invoice

import (
	"time"
)

type SaveAddressDTO struct {
	Street   string `json:"street" dynamodbav:"street" validate:"required,min=1,max=200"`
	City     string `json:"city" dynamodbav:"city" validate:"required,min=1,max=100"`
	PostCode string `json:"postCode" dynamodbav:"postCode" validate:"required,min=1,max=20"`
	Country  string `json:"country" dynamodbav:"country" validate:"required,min=1,max=100"`
}

type SaveItemDTO struct {
	Name     string  `json:"name" dynamodbav:"name" validate:"required,min=1,max=100"`
	Quantity int     `json:"quantity" dynamodbav:"quantity" validate:"required,min=1"`
	Price    float64 `json:"price" dynamodbav:"price" validate:"required,gt=0"`
}

type SaveInvoiceDTO struct {
	PaymentDue    time.Time      `json:"paymentDue" dynamodbav:"paymentDue,omitempty" validate:"required"`
	Description   string         `json:"description" dynamodbav:"description,omitempty" validate:"required,min=1,max=500"`
	ClientName    string         `json:"clientName" dynamodbav:"clientName,omitempty" validate:"required,min=1,max=100"`
	ClientEmail   string         `json:"clientEmail" dynamodbav:"clientEmail,omitempty" validate:"required,email"`
	SenderAddress SaveAddressDTO `json:"senderAddress" dynamodbav:"senderAddress,omitempty" validate:"required"`
	ClientAddress SaveAddressDTO `json:"clientAddress" dynamodbav:"clientAddress,omitempty" validate:"required"`
	Items         []SaveItemDTO  `json:"items" dynamodbav:"items,omitempty" validate:"required,min=1,dive"`
}

type UpdateInvoiceDTO struct {
	PaymentDue    time.Time      `json:"paymentDue" dynamodbav:"paymentDue,omitempty" validate:"required"`
	Description   string         `json:"description" dynamodbav:"description,omitempty" validate:"required,min=1,max=500"`
	ClientName    string         `json:"clientName" dynamodbav:"clientName,omitempty" validate:"required,min=1,max=100"`
	ClientEmail   string         `json:"clientEmail" dynamodbav:"clientEmail,omitempty" validate:"required,email"`
	SenderAddress SaveAddressDTO `json:"senderAddress" dynamodbav:"senderAddress,omitempty" validate:"required"`
	ClientAddress SaveAddressDTO `json:"clientAddress" dynamodbav:"clientAddress,omitempty" validate:"required"`
	/*
		When passing this array, we should pass the desired state
		i.e. if we currently have items '1, 2, 3' and we pass '1, 3' in here,
		it would mean that we want to delete second item
	*/
	Items []string `json:"items" dynamodbav:"items" validate:"required,dive,uuid"`
}
