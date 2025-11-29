package invoice

import (
	"time"
)

type SaveAddressDTO struct {
	Street   string `json:"street" validate:"required,min=1,max=200"`
	City     string `json:"city" validate:"required,min=1,max=100"`
	PostCode string `json:"postCode" validate:"required,min=1,max=20"`
	Country  string `json:"country" validate:"required,min=1,max=100"`
}

type SaveItemDTO struct {
	Name     string  `json:"name" validate:"required,min=1,max=100"`
	Quantity int     `json:"quantity" validate:"required,min=1"`
	Price    float64 `json:"price" validate:"required,gt=0"`
}

type SaveInvoiceDTO struct {
	PaymentDue    time.Time      `json:"paymentDue" validate:"required"`
	Description   string         `json:"description" validate:"required,min=1,max=500"`
	ClientName    string         `json:"clientName" validate:"required,min=1,max=100"`
	ClientEmail   string         `json:"clientEmail" validate:"required,email"`
	SenderAddress SaveAddressDTO `json:"senderAddress" validate:"required"`
	ClientAddress SaveAddressDTO `json:"clientAddress" validate:"required"`
	Items         []SaveItemDTO  `json:"items" validate:"required,min=1,dive"`
}
