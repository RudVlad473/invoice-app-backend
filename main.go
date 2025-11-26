package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rudvlad473/invoice-app-backend/invoice"
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

	/* Repositories */
	invoice.NewRepository(dynamodbClient)
	/**/

	/* Services */

	/**/

	/* Handles (controllers) */
	http.Handle(pkg.UrlInvoices, invoice.NewHandler())
	/**/

	http.ListenAndServe(":8080", nil)
}
