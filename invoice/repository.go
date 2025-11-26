package invoice

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dynamodb_client "github.com/rudvlad473/invoice-app-backend/appdynamodb/constants"
	invoice "github.com/rudvlad473/invoice-app-backend/invoice/models"
)

type Repository struct {
	dynamodbClient *dynamodb.Client
}

func NewRepository(dynamoDBClient *dynamodb.Client) *Repository {
	return &Repository{dynamodbClient: dynamoDBClient}
}

func (r *Repository) FindById(ctx context.Context, id string) (invoice.Invoice, error) {
	v, err := r.dynamodbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(dynamodb_client.TableNameInvoices),
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return invoice.Invoice{}, err
	}

	if v == nil || v.Item == nil {
		return invoice.Invoice{}, fmt.Errorf("item with id '%s' not found", id)
	}

	var foundInvoice invoice.Invoice
	err = attributevalue.UnmarshalMap(v.Item, &foundInvoice)

	if err != nil {
		return invoice.Invoice{}, err
	}

	return foundInvoice, nil
}

func (r *Repository) Save(ctx context.Context, invoice invoice.Invoice) error {
	item, err := attributevalue.MarshalMap(invoice)

	if err != nil {
		return err
	}

	_, err = r.dynamodbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dynamodb_client.TableNameInvoices),
		Item:      item,
	})

	if err != nil {
		return err
	}

	return nil
}
