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
	"github.com/loudsmilestudios/TetraControl/core"
)

var isIntialized bool
var awsSession session.Session
var dynamodbClient *dynamodb.DynamoDB
var ecsClient *ecs.ECS
var ec2Client *ec2.EC2

const configFile = "aws.yaml"

type awsConfig struct {
	DynamoTable         string    `yaml:"table" env:"DYNAMODB_TABLE"`
	EcsCluster          string    `yaml:"ecs_cluster" env:"ECS_CLUSTER"`
	TaskDefinition      string    `yaml:"task_definition" env:"TASK_DEFINITION"`
	SecurityGroups      []*string `yaml:"security_groups" env:"SecurityGroups"`
	VpcID               string    `yaml:"vpc_id" env:"VPC_ID"`
	AWSProfile          string    `yaml:"profile" env:"AWS_PROFILE"`
	AWSRegion           string    `yaml:"region" env:"AWS_REGION"`
	GameServerContainer string    `yaml:"container_name" env:"CONTAINER_NAME" env-default:"game_server"`
	GameServerPort      uint16    `yaml:"server_port" env:"SERVER_PORT" env-default:"7777"`
}

var config awsConfig

func intializeAws() error {
	if !isIntialized {
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

		if config.DynamoTable == "" {
			log.Fatal("DynamoTable not set!")
		}
		if config.EcsCluster == "" {
			log.Fatal("ECS Cluster not set!")
		}
		if config.TaskDefinition == "" {
			log.Fatal("Task definition is not set!")
		}
		if config.VpcID == "" {
			log.Fatal("VPC ID is not set!")
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
		isIntialized = true
		return nil
	}
	return nil
}

// RegisterModules registers all AWS modules
func RegisterModules() {
	core.RegisterServerModule("aws", ServerManager{})
	core.RegisterServerlessModule("aws", ServerlessAdapter{})
}
