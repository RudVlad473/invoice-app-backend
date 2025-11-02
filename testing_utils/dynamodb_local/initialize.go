package testing_utils

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// InitializeDynamodbLocalCommand https://refactoring.guru/design-patterns/command/go/example

type InitializeDynamodbLocalCommand struct {
	setupDocker    *SetupDynamodbLocalDockerCommand
	ctx            context.Context
	dynamodbClient *dynamodb.Client
}

func NewInitializeDynamoDBLocalCommand(ctx context.Context, dynamodbClient *dynamodb.Client) *InitializeDynamodbLocalCommand {
	return &InitializeDynamodbLocalCommand{
		setupDocker:    &SetupDynamodbLocalDockerCommand{},
		ctx:            ctx,
		dynamodbClient: dynamodbClient,
	}
}

// Execute /*
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
func (c *InitializeDynamodbLocalCommand) Execute() {
	// 1
	c.setupDocker.Execute()

	// 2
	SetupDynamodbLocalEnv()

	// 4
	err := SetupDynamoDBLocalTables(c.dynamodbClient, c.ctx)

	if err != nil {
		log.Fatalf("unable to setup DynamoDB local tables, %v", err)
	}
}

func (c *InitializeDynamodbLocalCommand) Undo() {
	c.setupDocker.Undo()
}
