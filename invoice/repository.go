package invoice

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	dynamodb_client "github.com/rudvlad473/invoice-app-backend/appdynamodb/constants"
	invoiceModels "github.com/rudvlad473/invoice-app-backend/invoice/models"
)

type Repository struct {
	dynamodbClient *dynamodb.Client
}

func NewRepository(dynamoDBClient *dynamodb.Client) *Repository {
	return &Repository{dynamodbClient: dynamoDBClient}
}

func (r *Repository) FindById(ctx context.Context, id string) (invoiceModels.Invoice, error) {
	v, err := r.dynamodbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(dynamodb_client.TableNameInvoices),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return invoiceModels.Invoice{}, err
	}

	if v == nil || v.Item == nil {
		return invoiceModels.Invoice{}, fmt.Errorf("item with id '%s' not found", id)
	}

	var foundInvoice invoiceModels.Invoice
	err = attributevalue.UnmarshalMap(v.Item, &foundInvoice)

	if err != nil {
		return invoiceModels.Invoice{}, err
	}

	return foundInvoice, nil
}

func (r *Repository) Save(ctx context.Context, invoice invoiceModels.SaveInvoiceDTO) (invoiceModels.Invoice, error) {
	item, err := attributevalue.MarshalMap(invoice)
	item["id"] = &types.AttributeValueMemberS{Value: uuid.New().String()}

	if err != nil {
		return invoiceModels.Invoice{}, err
	}

	_, err = r.dynamodbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dynamodb_client.TableNameInvoices),
		Item:      item,
	})

	if err != nil {
		return invoiceModels.Invoice{}, err
	}

	var savedInvoice invoiceModels.Invoice
	err = attributevalue.UnmarshalMap(item, &savedInvoice)

	if err != nil {
		return invoiceModels.Invoice{}, err
	}

	return savedInvoice, nil
}
