package main

import (
	"context"
	"log"
	"net/http"
	"slices"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-playground/validator/v10"
	"github.com/rudvlad473/invoice-app-backend/invoice"
	invoiceConstants "github.com/rudvlad473/invoice-app-backend/invoice/constants"
	pkg "github.com/rudvlad473/invoice-app-backend/pkg/constants"
)

func main() {
	/* 'context' variables */
	ctx := context.Background()
	/**/

	/* Database */

	// client setup
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	dynamodbClient := dynamodb.NewFromConfig(cfg)

	/**/

	/* Register validators for DTOs */

	var validate = validator.New()

	err = validate.RegisterValidation("status", func(fl validator.FieldLevel) bool {
		status := fl.Field().Interface().(invoiceConstants.Status)

		return slices.Contains(invoiceConstants.Statuses, status)
	})

	if err != nil {
		log.Fatalf("unable to register validator: %v", err)
	}

	/**/

	/* Repositories */
	invoice.NewRepository(dynamodbClient)
	/**/

	/* Services */

	/**/

	/* Handles (controllers) */
	http.Handle(pkg.UrlInvoices, invoice.NewHandler())
	/**/

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalf("unable to start server: %v", err)
	}
}
