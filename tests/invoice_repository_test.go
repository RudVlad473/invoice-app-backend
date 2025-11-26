package tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rudvlad473/invoice-app-backend/invoice"
	invoiceModels "github.com/rudvlad473/invoice-app-backend/invoice/models"
	"github.com/rudvlad473/invoice-app-backend/testingutils/dynamodblocal"
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
