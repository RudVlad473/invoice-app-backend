package testingutils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/brianvoe/gofakeit/v6"
	dynamodbConstants "github.com/rudvlad473/invoice-app-backend/appdynamodb/constants"
	invoiceConstants "github.com/rudvlad473/invoice-app-backend/invoice/constants"
	invoiceModels "github.com/rudvlad473/invoice-app-backend/invoice/models"
	testing_utils "github.com/rudvlad473/invoice-app-backend/testingutils/dynamodblocal/constants"
)

type AppDynamodb struct {
	DynamodbClient *dynamodb.Client
}

func NewAppDynamodb() *AppDynamodb {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("us-west-2"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
		config.WithHTTPClient(&http.Client{
			Timeout: 30 * time.Second, // Effectively disables IMDS
		}),
		config.WithBaseEndpoint(testing_utils.LocalDynamodbUrl),
	)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	dynamodbClient := dynamodb.NewFromConfig(cfg)

	return &AppDynamodb{DynamodbClient: dynamodbClient}
}

type TableDef struct {
	Name       string
	HashKey    string
	RangeKey   string
	Attributes []types.AttributeDefinition
}

var tableDefs = []TableDef{
	{
		Name:    dynamodbConstants.TableNameInvoices,
		HashKey: "id",
		Attributes: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
	},
}

// SetupTables
// /* FOR TESTS ONLY */
func (appDynamodb *AppDynamodb) SetupTables() error {
	ctx := context.Background()

	for _, tableDef := range tableDefs {
		err := setupTable(appDynamodb.DynamodbClient, ctx, tableDef)

		if err != nil {
			return fmt.Errorf("failed creating table '%s': %w", tableDef.Name, err)
		}
	}

	return nil
}

func setupTable(dynamoClient *dynamodb.Client, ctx context.Context, t TableDef) error {
	// Check if t exists
	_, err := dynamoClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(t.Name),
	})

	// means that table exists
	if err == nil {
		return nil
	}

	input := &dynamodb.CreateTableInput{
		TableName:            aws.String(t.Name),
		AttributeDefinitions: t.Attributes,
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String(t.HashKey), KeyType: types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	}

	if t.RangeKey != "" {
		input.KeySchema = append(input.KeySchema, types.KeySchemaElement{
			AttributeName: aws.String(t.RangeKey),
			KeyType:       types.KeyTypeRange,
		})
	}

	_, err = dynamoClient.CreateTable(ctx, input)
	if err != nil {
		return fmt.Errorf("failed creating t %s: %w", t.Name, err)
	}

	if err := waitUntilTableActive(dynamoClient, ctx, t.Name); err != nil {
		return fmt.Errorf("t %s not active: %w", t.Name, err)
	}

	return nil
}

func waitUntilTableActive(dynamoClient *dynamodb.Client, ctx context.Context, table string) error {
	for {
		out, err := dynamoClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(table),
		})

		if err != nil {
			return err
		}

		if out.Table.TableStatus == types.TableStatusActive {
			return nil
		}

		time.Sleep(200 * time.Millisecond)
	}
}

// CleanupTables
// /* FOR TESTS ONLY */
func (appDynamodb *AppDynamodb) CleanupTables() error {
	ctx := context.Background()

	out, err := appDynamodb.DynamodbClient.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		return err
	}

	for _, name := range out.TableNames {
		_, err := appDynamodb.DynamodbClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{
			TableName: aws.String(name),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// PopulateTables
// /* FOR TESTS ONLY */
func (appDynamodb *AppDynamodb) PopulateTables(shouldPopulateInvoices bool) ([]invoiceModels.Invoice, error) {
	ctx := context.Background()

	invoiceCount := 20
	if !shouldPopulateInvoices {
		invoiceCount = 0
	}

	var invoices []invoiceModels.Invoice

	for i := 0; i < invoiceCount; i++ {
		invoice := getFakeInvoice()
		invoices = append(invoices, invoice)

		av, err := attributevalue.MarshalMap(invoice)

		if err != nil {
			return nil, err
		}

		_, err = appDynamodb.DynamodbClient.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(dynamodbConstants.TableNameInvoices),
			Item:      av,
		})

		if err != nil {
			return nil, err
		}
	}

	return invoices, nil
}

func getFakeInvoice() invoiceModels.Invoice {
	now := time.Now().UTC()

	var items []invoiceModels.Item
	for i := 0; i < gofakeit.Number(1, 10); i++ {
		items = append(items, getFakeItem())
	}

	statuses := []invoiceConstants.Status{
		invoiceConstants.StatusDraft,
		invoiceConstants.StatusPending,
		invoiceConstants.StatusPaid,
	}

	senderAddress := gofakeit.Address()
	clientAddress := gofakeit.Address()

	return invoiceModels.Invoice{
		Id:          gofakeit.UUID(),
		CreatedAt:   now,
		PaymentDue:  now.AddDate(0, 0, gofakeit.Number(3, 15)),
		Description: gofakeit.ProductDescription(),
		ClientName:  gofakeit.Name(),
		ClientEmail: gofakeit.Email(),
		Status:      statuses[gofakeit.Number(0, len(statuses)-1)],
		SenderAddress: invoiceModels.Address{
			City:     senderAddress.City,
			Country:  senderAddress.Country,
			PostCode: senderAddress.Zip,
			Street:   senderAddress.Street,
		},
		ClientAddress: invoiceModels.Address{
			City:     clientAddress.City,
			Country:  clientAddress.Country,
			PostCode: clientAddress.Zip,
			Street:   clientAddress.Street,
		},
		Items: items,
	}
}

func getFakeItem() invoiceModels.Item {
	return invoiceModels.Item{
		Name:     gofakeit.ProductName(),
		Quantity: gofakeit.Number(1, 35),
		Price:    gofakeit.Price(0.25, 1000.0),
	}
}
