package testingutils

import (
	"context"
)

// InitializeDynamodbLocalCommand https://refactoring.guru/design-patterns/command/go/example

type InitializeDynamodbLocalCommand struct {
	setupDocker *SetupDynamodbLocalDockerCommand
	ctx         context.Context
	appDynamodb *AppDynamodb
}

func NewInitializeDynamoDBLocalCommand(ctx context.Context, appDynamodb *AppDynamodb) *InitializeDynamodbLocalCommand {
	return &InitializeDynamodbLocalCommand{
		setupDocker: &SetupDynamodbLocalDockerCommand{},
		ctx:         ctx,
		appDynamodb: appDynamodb,
	}
}

// Execute /*
/*
For test setup we need:
1. Spin up docker container with appdynamodb local
2. Initialize appdynamodb client
3. Populate that db with tables needed for the app
*/
func (c *InitializeDynamodbLocalCommand) Execute() {
	// 1
	c.setupDocker.Execute()

	// 2,3 - should be done per test
}

func (c *InitializeDynamodbLocalCommand) Undo() {
	c.setupDocker.Undo()
}
