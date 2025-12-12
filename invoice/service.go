package invoice

import (
	"context"
	"fmt"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
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
	return r.updateById(ctx, id, invoice)
}

/*
You're supposed to pass a partial invoice here (or a full one) with all the fields you want updated
So don't pass fields that you haven't updated, because it'll cause a dry update
Sidenote:I spent HOURS trying to understand how to make it type-safe instead of putting `any` here
There isn't a way to do that in Go at the moment of writing of this function
It could potentially be moved to a reusable dynamodb client for reusability
*/
func (r *Repository) updateById(ctx context.Context, id string, invoice any) (updatedInvoice invoiceModels.Invoice, err error) {
	_, err = r.FindById(ctx, id)

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

	_, err = r.dynamodbClient.UpdateItem(ctx, input)

	if err != nil {
		return invoiceModels.Invoice{}, err
	}

	updatedInvoice, err = r.FindById(ctx, id)

	if err != nil {
		return invoiceModels.Invoice{}, err
	}

	return updatedInvoice, nil
}

func (r *Repository) AddItemByInvoiceId(ctx context.Context, invoiceId string, item invoiceModels.SaveItemDTO) (updatedInvoice invoiceModels.Invoice, err error) {
	invoice, err := r.FindById(ctx, invoiceId)

	if err != nil {
		return invoiceModels.Invoice{}, err
	}

	itemToAppend := invoiceModels.Item{}
	/*
		This logic allows us to 'merge' (js-style) two structure,
		so we don't have to map out all fields manually
		In js, it would be a simple `const itemToAppend = {...item}`
		But golang doesn't have a similar spread functionality
	*/
	err = mapstructure.Decode(item, &itemToAppend)

	if err != nil {
		return invoiceModels.Invoice{}, err
	}

	/*
		We only pass `Items` here because we don't want any other fields updated
		When updating the items here, we basically replace old items with the new ones, adding an item at the end
		This is not the most efficient solution, but should work smoothly for small amount of data
		(and the data is expected to be small)
	*/
	return r.updateById(ctx, invoiceId, invoiceModels.UpdateInvoiceItemDTO{
		Items: append(invoice.Items, itemToAppend),
	})
}

func (r *Repository) RemoveItemByInvoiceId(ctx context.Context, invoiceId string, itemId string) (updatedInvoice invoiceModels.Invoice, err error) {
	invoice, err := r.FindById(ctx, invoiceId)

	if err != nil {
		return invoiceModels.Invoice{}, err
	}

	indexOfItemToDelete := slices.IndexFunc(invoice.Items, func(item invoiceModels.Item) bool { return item.Id == itemId })

	if indexOfItemToDelete == -1 {
		return invoiceModels.Invoice{}, fmt.Errorf("item with id '%s' not found", itemId)
	}

	items := slices.Delete(invoice.Items, indexOfItemToDelete, indexOfItemToDelete+1)

	return r.updateById(ctx, invoiceId, invoiceModels.UpdateInvoiceItemDTO{
		/*
			2nd parameter is starting index to delete from, 3rd - ending index
			Example: slices.Delete([1,2,3], 1, 2) => [1,3]
		*/
		Items: items,
	})
}

func (r *Repository) UpdateItemByInvoiceId(ctx context.Context, invoiceId string,
	itemId string, item invoiceModels.UpdateItemDTO) (updatedInvoice invoiceModels.Invoice, err error) {
	invoice, err := r.FindById(ctx, invoiceId)

	if err != nil {
		return invoiceModels.Invoice{}, err
	}

	itemToUpdateIndex := slices.IndexFunc(invoice.Items, func(item invoiceModels.Item) bool { return item.Id == itemId })

	if itemToUpdateIndex == -1 {
		return invoiceModels.Invoice{}, fmt.Errorf("item with id '%s' not found", itemId)
	}

	/*Overriding existing item with passed dto*/
	err = mapstructure.Decode(item, &invoice.Items[itemToUpdateIndex])

	if err != nil {
		return invoiceModels.Invoice{}, err
	}

	return r.updateById(ctx, invoiceId, invoiceModels.UpdateInvoiceItemDTO{
		Items: invoice.Items,
	})
}
