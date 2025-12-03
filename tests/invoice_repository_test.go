package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/fatih/structs"
	"github.com/google/go-cmp/cmp"
	"github.com/rudvlad473/invoice-app-backend/invoice"
	invoiceModels "github.com/rudvlad473/invoice-app-backend/invoice/models"
	"github.com/rudvlad473/invoice-app-backend/testingutils/dynamodblocal"
	testing_utils "github.com/rudvlad473/invoice-app-backend/testingutils/fakes"
)

var appDynamodb = testingutils.NewAppDynamodb()
var invoiceRepository = invoice.NewRepository(appDynamodb.DynamodbClient)
var ctx = context.Background()

func Setup(t *testing.T, shouldPopulateInvoices bool) []invoiceModels.Invoice {
	t.Helper()

	err := appDynamodb.SetupTables()

	if err != nil {
		t.Fatalf("unable to setup DynamoDB local tables, %v", err)
	}

	invoices, err := appDynamodb.PopulateTables(shouldPopulateInvoices)

	if err != nil {
		t.Fatalf("unable to populate Tables, %v", err)
	}

	t.Cleanup(func() {
		err := appDynamodb.CleanupTables()

		if err != nil {
			t.Fatalf("unable to cleanup DynamoDB local tables, %v", err)
		}
	})

	return invoices
}

func TestFindById(t *testing.T) {
	t.Run("should return when invoices exist", func(t *testing.T) {
		// arrange
		invoices := Setup(t, true)
		invoiceId := invoices[gofakeit.Number(0, len(invoices)-1)].Id

		// act
		foundInvoice, err := invoiceRepository.FindById(ctx, invoiceId)

		// assert
		if err != nil {
			t.Fatalf("invoice with id of '%s' wasn't found in repository, %s", invoiceId, err)
		}
		if foundInvoice.Id != invoiceId {
			t.Fatalf("invoice with id of '%s' wasn't found in repository, instead found: '%s'", invoiceId, foundInvoice.Id)
		}
	})

	t.Run("should NOT return when NO invoices exist", func(t *testing.T) {
		// arrange
		Setup(t, false)
		randomId := gofakeit.UUID()

		// act
		emptyInvoice, err := invoiceRepository.FindById(ctx, randomId)

		// assert
		if err == nil {
			t.Fatalf("error was not returned when it should've (no invoice exists)")
		}
		if !reflect.DeepEqual(emptyInvoice, invoiceModels.Invoice{}) {
			t.Fatalf("returned invoice was not empty, %v", emptyInvoice)
		}
	})
}

func TestSave(t *testing.T) {
	t.Run("should save invoice", func(t *testing.T) {
		// arrange
		Setup(t, false)
		invoiceToSave := testing_utils.CreateSaveInvoiceDTO()

		// act
		savedInvoice, err := invoiceRepository.Save(ctx, invoiceToSave)

		// assert
		if err != nil {
			t.Fatalf("couldn't save invoice \n %s", err)
		}
		// TODO: add tests for DTOs
		if _, err = invoiceRepository.FindById(ctx, savedInvoice.Id); err != nil {
			t.Fatalf("couldn't find the saved invoice \n %s", err)
		}
	})
}

func TestCountAll(t *testing.T) {
	t.Run("should return TOTAL number of invoices in the table", func(t *testing.T) {
		// arrange
		invoices := Setup(t, true)

		// act
		amountOfInvoices, err := invoiceRepository.CountAll(ctx)

		// assert
		if err != nil {
			t.Fatalf("couldn't count invoices \n %s", err)
		}
		if amountOfInvoices != len(invoices) {
			t.Fatalf("invalid amount of invoices returned, actual amount: '%d', returned amount: '%d' \n %s", len(invoices), amountOfInvoices, err)
		}
	})
}

func TestDeleteById(t *testing.T) {
	t.Run("should delete invoice by id", func(t *testing.T) {
		// arrange
		invoices := Setup(t, true)
		invoiceIdToDelete := invoices[gofakeit.Number(0, len(invoices)-1)].Id

		// act
		err := invoiceRepository.DeleteById(ctx, invoiceIdToDelete)

		// assert
		if err != nil {
			t.Fatalf("couldn't delete invoice \n %s", err)
		}
		if _, err = invoiceRepository.FindById(ctx, invoiceIdToDelete); err == nil {
			t.Fatalf("invoice was still found after deletion, id: `%s` \n %s", invoiceIdToDelete, err)
		}
	})

	t.Run("should NOT delete invoice by id when id doesn't exist", func(t *testing.T) {
		// arrange
		invoices := Setup(t, true)
		randomInvoiceIdToDelete := gofakeit.UUID()

		// act
		err := invoiceRepository.DeleteById(ctx, randomInvoiceIdToDelete)

		// assert
		if err == nil {
			t.Fatalf("deleted invoice when shouldn't have, id: `%s` \n %s", randomInvoiceIdToDelete, err)
		}
		if countAfterDeletion, err := invoiceRepository.CountAll(ctx); err != nil || countAfterDeletion != len(invoices) {
			t.Fatalf("amount of invoices was changed after deletion, amount after deletion: `%d` \n %s", countAfterDeletion, err)
		}
	})
}

/*
Here we assume that all passed dto's are correct
Tests that validate DTOs should be written separately
*/
func TestUpdateById(t *testing.T) {
	/*
		We need this to iterate over all keys of DTO
		to make sure update logic works as expected for all cases
	*/
	updateDtoKeyValueMap := structs.Map(testing_utils.CreateUpdateInvoiceDTO())
	for key := range updateDtoKeyValueMap {
		t.Run(fmt.Sprintf("should partially update invoice (without updating items), key: '%s', value: '%+v'", key, updateDtoKeyValueMap[key]), func(t *testing.T) {
			// arrange
			invoices := Setup(t, true)
			invoiceToUpdate := invoices[gofakeit.Number(0, len(invoices)-1)]

			/*
				Here we set 1 field that we want to update
			*/
			finalUpdateDto := invoiceModels.UpdateInvoiceDTO{}
			jsonBytes, _ := json.Marshal(map[string]interface{}{
				key: updateDtoKeyValueMap[key],
			})
			err := json.Unmarshal(jsonBytes, &finalUpdateDto)

			if err != nil {
				t.Fatalf("couldn't unmarshal updated invoice \n %s", err)
			}

			// act
			updatedInvoice, err := invoiceRepository.UpdateById(ctx, invoiceToUpdate.Id, finalUpdateDto)

			// assert
			if err != nil {
				t.Fatalf("couldn't update invoice \n %s", err)
			}
			if _, err = invoiceRepository.FindById(ctx, invoiceToUpdate.Id); err != nil {
				t.Fatalf("couldn't find the updated invoice \n %s", err)
			}
			// cmp is a better alternative to reflect.DeepEqual
			if cmp.Equal(invoiceToUpdate, updatedInvoice) {
				t.Fatalf("invoice stayed the same after update")
			}
			// field should be updated as expected
			if !cmp.Equal(structs.Map(updatedInvoice)[key], updateDtoKeyValueMap[key]) {
				t.Fatalf("expected '%s' field to be equal to '%+v' value, instead was '%+v'", key, updateDtoKeyValueMap[key], structs.Map(updatedInvoice)[key])
			}
			for invoiceDtoKey := range updateDtoKeyValueMap {
				if invoiceDtoKey == key {
					continue
				}

				if !cmp.Equal(structs.Map(updatedInvoice)[key], updateDtoKeyValueMap[key]) {
					t.Fatalf("other fields of the entity SHOULD NOT be updated, instead saw '%s' field being updated from '%+v' to '%+v'", invoiceDtoKey, structs.Map(invoiceToUpdate)[invoiceDtoKey], structs.Map(updatedInvoice)[invoiceDtoKey])
				}
			}
		})
	}
}
