package invoice

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dynamodb_client "github.com/rudvlad473/invoice-app-backend/dynamodb/constants"
	invoice "github.com/rudvlad473/invoice-app-backend/invoice/models"
)

type Repository struct {
	dynamodbClient *dynamodb.Client
}

func NewRepository(dynamoDBClient *dynamodb.Client) *Repository {
	return &Repository{dynamodbClient: dynamoDBClient}
}

func (r *Repository) GetInvoices(ctx context.Context, id string) ([]invoice.Invoice, error) {
	v, err := r.dynamodbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(dynamodb_client.TableNameInvoices),
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return nil, err
	}

	var invoices []invoice.Invoice
	err = attributevalue.UnmarshalMap(v.Item, &invoices)

	if err != nil {
		return nil, err
	}

	return invoices, nil
}
