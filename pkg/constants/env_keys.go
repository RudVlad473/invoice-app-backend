package pkg

type EnvKey string

const (
	/*
		Set via CLI
	*/
	EnvKeyMode EnvKey = "MODE"

	/*
		Set programmatically
	*/
	EnvKeyDynamodbUrl         EnvKey = "AWS_ENDPOINT_URL_DYNAMODB"
	EnvKeyAWSRegion           EnvKey = "AWS_REGION"
	EnvKeyAWSAccessKeyId      EnvKey = "AWS_ACCESS_KEY_ID"
	EnvKeyAWSSecretAccessKey  EnvKey = "AWS_SECRET_ACCESS_KEY"
	EnvKeyAwsEndpointUrl      EnvKey = "AWS_ENDPOINT_URL"
	EnvKeyEC2MetadataDisabled EnvKey = "AWS_EC2_METADATA_DISABLED"
)
