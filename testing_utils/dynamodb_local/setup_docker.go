package testing_utils

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	dynamodb_local "github.com/rudvlad473/invoice-app-backend/testing_utils/dynamodb_local/constants"
)

type SetupDynamodbLocalDockerCommand struct{}

func (c *SetupDynamodbLocalDockerCommand) Execute() {
	cmd := exec.Command(
		"docker",
		"run",
		"--rm",
		"-d",
		"--name",
		"dynamodb-local-test",
		"-p",
		fmt.Sprintf("%d:%d", dynamodb_local.DynamodbLocalPort, dynamodb_local.DynamodbLocalPort),
		"amazon/dynamodb-local",
	)

	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to start DynamoDB Local: %v", err)
	}

	/*
		Wait here before container is actually initialize
		Not the best approach, possibly rework it to await the
		actual container itself somehow
	*/
	time.Sleep(2 * time.Second)
}

func (c *SetupDynamodbLocalDockerCommand) Undo() {
	_ = exec.Command("docker", "stop", "dynamodb-local-test").Run()
}
