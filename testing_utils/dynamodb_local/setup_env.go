package testing_utils

import (
	"fmt"
	"log"
	"os"

	pkg "github.com/rudvlad473/invoice-app-backend/pkg/constants"
	dynamodb_local "github.com/rudvlad473/invoice-app-backend/testing_utils/dynamodb_local/constants"
)

var envKeyValueMap = map[string]string{
	string(pkg.EnvKeyDynamodbUrl):        fmt.Sprintf("http://localhost:%d", dynamodb_local.DynamodbLocalPort),
	string(pkg.EnvKeyAWSRegion):          "us-west-2",
	string(pkg.EnvKeyAWSSecretAccessKey): "dummy",
	string(pkg.EnvKeyAWSAccessKeyId):     "dummy",
}

func SetupDynamodbLocalEnv() {
	for k, v := range envKeyValueMap {
		err := os.Setenv(k, v)

		if err != nil {
			log.Fatalf("couldn't set env variable '%s' to '%s', %v", k, v, err)
		}
	}
}
