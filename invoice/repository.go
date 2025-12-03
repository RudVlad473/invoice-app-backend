package invoice

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
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

func (r *Repository) DeleteById(ctx context.Context, id string) error {
	_, err := r.FindById(ctx, id)

	if err != nil {
		return err
	}

	_, err = r.dynamodbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(dynamodb_client.TableNameInvoices),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CountAll(ctx context.Context) (int, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(dynamodb_client.TableNameInvoices),
		Select:    types.SelectCount,
	}

	totalCount := 0

	for {
		output, err := r.dynamodbClient.Scan(ctx, input)
		if err != nil {
			return 0, err
		}

		totalCount += int(output.Count)

		if output.LastEvaluatedKey == nil {
			break
		}

		input.ExclusiveStartKey = output.LastEvaluatedKey
	}

	return totalCount, nil
}

func (r *Repository) UpdateById(ctx context.Context, id string, invoice invoiceModels.UpdateInvoiceDTO) (invoiceModels.Invoice, error) {
	_, err := r.FindById(ctx, id)

	if err != nil {
		return invoiceModels.Invoice{}, err
	}

	// Marshal DTO to DynamoDB attribute value map
	avMap, err := attributevalue.MarshalMap(invoice)

	if err != nil {
		return invoiceModels.Invoice{}, err
	}

	// Remove zero values if you want omitempty behavior
	// This step depends on your requirements

	// Build update expression from the map
	updateBuilder := expression.UpdateBuilder{}
	for key, val := range avMap {
		updateBuilder = updateBuilder.Set(expression.Name(key), expression.Value(val))
	}

	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()

	if err != nil {
		return invoiceModels.Invoice{}, err
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(dynamodb_client.TableNameInvoices),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	item, err := r.dynamodbClient.UpdateItem(ctx, input)

	if err != nil || item == nil {
		return invoiceModels.Invoice{}, err
	}

	var updatedInvoice invoiceModels.Invoice
	err = attributevalue.UnmarshalMap(item.Attributes, &updatedInvoice)

	if err != nil {
		return invoiceModels.Invoice{}, err
	}

	return updatedInvoice, nil
}
