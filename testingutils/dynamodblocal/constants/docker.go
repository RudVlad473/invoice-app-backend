package testing_utils

import "fmt"

const (
	DynamodbLocalPort   = 8000
	DockerContainerName = "appdynamodb-local-test"
)

var LocalDynamodbUrl = fmt.Sprintf("http://localhost:%d", DynamodbLocalPort)
