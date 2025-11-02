package invoice

import (
	"context"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	testingutils "github.com/rudvlad473/invoice-app-backend/testing_utils/dynamodb_local"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// client setup
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	dynamodbClient := dynamodb.NewFromConfig(cfg)

	initializeDynamoDBLocalCommand := testingutils.NewInitializeDynamoDBLocalCommand(ctx, dynamodbClient)

	initializeDynamoDBLocalCommand.Execute()

	m.Run()

	initializeDynamoDBLocalCommand.Undo()
}
