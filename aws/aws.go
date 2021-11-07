package aws

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/ilyakaznacheev/cleanenv"
)

var awsSession session.Session
var dynamodbClient *dynamodb.DynamoDB
var ecsClient *ecs.ECS
var ec2Client *ec2.EC2

const configFile = "aws.yaml"

type awsConfig struct {
	dynamoTable    string    `yaml:"dynamodb_table" env:"DYNAMODB_TABLE"`
	ecsCluster     string    `yaml:"ecs_cluster" env:"ECS_CLUSTER"`
	taskDefinition string    `yaml:"task_definition" env:"TASK_DEFINITION"`
	SecurityGroups []*string `yaml:"security_groups" env:"SecurityGroups"`
	VpcID          string    `yaml:"vpc_id" env:"VPC_ID"`
	AWSProfile     string    `yaml:"profile" env:"AWS_PROFILE"`
	AWSRegion      string    `yaml:"region" env:"AWS_REGION"`
}

var config awsConfig

func init() {

	if _, err := os.Stat(configFile); err == nil {
		// Load config from file
		if err := cleanenv.ReadConfig(configFile, &config); err != nil {
			log.Fatalf("could not load aws config from file: %v", err)
		}

	} else if os.IsNotExist(err) {
		// Load config from environment
		if err := cleanenv.ReadEnv(&config); err != nil {
			log.Fatalf("could not load aws config from environment: %v", err)
		}
	}

	// Load TetraControl AWS Config -> AWS Config
	sessionConfig := &aws.Config{}
	if len(config.AWSProfile) > 0 {
		sessionConfig.Credentials = credentials.NewSharedCredentials("", config.AWSProfile)
	}
	if len(config.AWSRegion) > 0 {
		sessionConfig.Region = aws.String(config.AWSRegion)
	}

	awsSession := session.Must(session.NewSession(sessionConfig))
	dynamodbClient = dynamodb.New(awsSession)
	ecsClient = ecs.New(awsSession)
	ec2Client = ec2.New(awsSession)
}
