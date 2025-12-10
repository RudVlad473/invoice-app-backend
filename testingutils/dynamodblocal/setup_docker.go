package testingutils

import (
	"fmt"
	"log"
	"os/exec"

	dynamodb_local "github.com/rudvlad473/invoice-app-backend/testingutils/dynamodblocal/constants"
)

// SetupDynamodbLocalDockerCommand
// *Should've probably used official docker api for go for this*/
type SetupDynamodbLocalDockerCommand struct{}

func (c *SetupDynamodbLocalDockerCommand) Execute() {
	// Check if container already running
	checkCmd := exec.Command(
		"docker", "ps", "-q", "-f",
		fmt.Sprintf("name=%s", dynamodb_local.DockerContainerName),
	)
	output, _ := checkCmd.Output()

	if len(output) > 0 {
		log.Printf("Container %s already running\n", dynamodb_local.DockerContainerName)
		return
	}

	cmd := exec.Command(
		"docker", "run", "--rm", "-d",
		"--name", dynamodb_local.DockerContainerName,
		"-p", fmt.Sprintf("%d:%d", dynamodb_local.DynamodbLocalPort, dynamodb_local.DynamodbLocalPort),
		"amazon/dynamodb-local",
	)

	log.Printf("Starting Dynamodb local docker container...")

	if _, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("failed to start DynamoDB Local: %v", err)
	}
}

func (c *SetupDynamodbLocalDockerCommand) Undo() {
	cmd := exec.Command("docker", "stop", dynamodb_local.DockerContainerName)

	if _, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("failed to pause local DynamoDB container: %v", err)
	}
}
