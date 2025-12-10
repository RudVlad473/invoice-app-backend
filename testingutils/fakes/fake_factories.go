package testing_utils

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
	invoiceModels "github.com/rudvlad473/invoice-app-backend/invoice/models"
)

func CreateSaveItemDTO() invoiceModels.SaveItemDTO {
	return invoiceModels.SaveItemDTO{
		Id:       gofakeit.UUID(),
		Name:     gofakeit.ProductName(),
		Price:    gofakeit.Price(1.0, 1000.0),
		Quantity: gofakeit.Number(10, 20),
	}
}

func CreateSaveAddressDTO() invoiceModels.SaveAddressDTO {
	return invoiceModels.SaveAddressDTO{
		City:     gofakeit.City(),
		Country:  gofakeit.Country(),
		Street:   gofakeit.Street(),
		PostCode: gofakeit.Zip(),
	}
}

func CreateSaveInvoiceDTO() invoiceModels.SaveInvoiceDTO {
	itemCount := gofakeit.Number(2, 10)

	var items []invoiceModels.SaveItemDTO

	for i := 0; i < itemCount; i++ {
		items = append(items, CreateSaveItemDTO())
	}

	clientAddress := CreateSaveAddressDTO()
	senderAddress := CreateSaveAddressDTO()

	return invoiceModels.SaveInvoiceDTO{
		ClientAddress: &clientAddress,
		SenderAddress: &senderAddress,
		ClientEmail:   gofakeit.Email(),
		ClientName:    gofakeit.Name(),
		Description:   gofakeit.ProductDescription(),
		PaymentDue:    time.Now().Add(24 * time.Hour * 7),
		Items:         items,
	}
}

// CreateUpdateInvoiceDTO /*
/*
	Doesn't populate items since it requires IDs of existing items,
	add it separately if needed
*/
func CreateUpdateInvoiceDTO() invoiceModels.UpdateInvoiceDTO {
	clientAddress := CreateSaveAddressDTO()
	senderAddress := CreateSaveAddressDTO()

	return invoiceModels.UpdateInvoiceDTO{
		ClientAddress: &clientAddress,
		SenderAddress: &senderAddress,
		ClientEmail:   gofakeit.Email(),
		ClientName:    gofakeit.Name(),
		Description:   gofakeit.ProductDescription(),
		PaymentDue:    time.Now().Add(24 * time.Hour * 7),
	}
}
