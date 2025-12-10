package tests

import (
	"context"
	"os"
	"testing"

	testingutils "github.com/rudvlad473/invoice-app-backend/testingutils/dynamodblocal"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	appDynamodb := testingutils.NewTestDynamodbClient()
	initializeDynamoDBLocalCommand := testingutils.NewInitializeDynamoDBLocalCommand(ctx, appDynamodb)

	initializeDynamoDBLocalCommand.Execute()

	code := m.Run()

	initializeDynamoDBLocalCommand.Undo()

	os.Exit(code)
}
