package testing_utils

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dynamodbConstants "github.com/rudvlad473/invoice-app-backend/dynamodb/constants"
)

type TableDef struct {
	Name       string
	HashKey    string
	RangeKey   string
	Attributes []types.AttributeDefinition
}

func SetupDynamoDBLocalTables(dynamoClient *dynamodb.Client, ctx context.Context) error {
	tableDefs := []TableDef{
		TableDef{
			Name:    dynamodbConstants.TableNameInvoices,
			HashKey: "Id",
			Attributes: []types.AttributeDefinition{
				{
					AttributeName: aws.String("Id"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
		},
	}

	for _, tableDef := range tableDefs {
		err := setupTable(dynamoClient, ctx, tableDef)

		if err != nil {
			return fmt.Errorf("failed creating t %s: %w", tableDef.Name, err)
		}
	}

	return nil
}

func setupTable(dynamoClient *dynamodb.Client, ctx context.Context, t TableDef) error {
	// Check if t exists
	_, err := dynamoClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(t.Name),
	})

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
