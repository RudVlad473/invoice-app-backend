package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rudvlad473/invoice-app-backend/invoice"
	pkg "github.com/rudvlad473/invoice-app-backend/pkg/constants"
	dynamodb_local "github.com/rudvlad473/invoice-app-backend/testing_utils/dynamodb_local"
)

func main() {
	/* 'context' variables */
	ctx := context.Background()
	mode := pkg.Mode(os.Getenv(string(pkg.EnvKeyMode)))
	/**/

	/* Database */

	// 1.1 test setup
	/*
		For test setup we need:
		1. Spin up docker container with dynamodb local
		2. Populate env variables needed for dynamodb client to connect to that local instance
		  - AWS_ENDPOINT_URL_DYNAMODB
		  - AWS_REGION
		  - AWS_ACCESS_KEY_ID
		  - AWS_SECRET_ACCESS_KEY

		3. Initialize dynamodb client
		4. Populate that db with tables needed for the app
	*/
	if mode == pkg.ModeDev {
		initializeDynamoDBLocalCommand := &dynamodb_local.InitializeDynamodbLocalCommand{}

		initializeDynamoDBLocalCommand.Execute()

		defer initializeDynamoDBLocalCommand.Undo()
	}

	// client setup
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	dynamodbClient := dynamodb.NewFromConfig(cfg)

	// 1.2 test setup
	if mode == pkg.ModeDev {
		err := dynamodb_local.SetupDynamoDBLocalTables(dynamodbClient, ctx)

		if err != nil {
			log.Fatalf("unable to setup DynamoDB local tables, %v", err)
		}
	}

	/**/

	/* Repositories */
	invoice.NewRepository(dynamodbClient)
	/**/

	/* Services */

	/**/

	/* Handles (controllers) */
	http.Handle(pkg.URL_INVOICES, invoice.NewHandler())
	/**/

	http.ListenAndServe(":8080", nil)
}
