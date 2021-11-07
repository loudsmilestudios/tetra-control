package aws

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/loudsmilestudios/TetraControl/core"
)

// ServerManager implements core.ServerManager to manage game servers on AWS
type ServerManager struct {
	config core.ServerConfig
}

// GetServer looks up a
func (manager *ServerManager) GetServer(identifier string) (core.Server, error) {

	result, err := dynamodbClient.GetItem(&dynamodb.GetItemInput{
		TableName: &config.dynamoTable,
		Key: map[string]*dynamodb.AttributeValue{
			"identifier": {
				S: aws.String(identifier),
			},
		},
	})

	// Pass along error
	if err != nil {
		return nil, err
	}

	// Dynamo Item could not be found
	if result.Item == nil {
		return nil, nil
	}

	server := Server{}
	err = dynamodbattribute.UnmarshalMap(result.Item, server)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal server: %v", err)
	}

	return &server, nil
}

// DeleteServerByIdentifier gets a server object and passes it to DeleteServer()
func (manager *ServerManager) DeleteServerByIdentifier(identifier string) error {

	// Grab server from identifier
	server, err := manager.GetServer(identifier)
	if err != nil {
		return err
	}

	// If server exists
	if server != nil {
		// Pass server object to delete server
		return manager.DeleteServer(server)
	}
	return nil
}

// DeleteServer stops the ECS task associated with a server
func (manager *ServerManager) DeleteServer(server core.Server) error {

	// Cast core.Server to aws.Server
	awsServer, isAwsServer := server.(*Server)
	if !isAwsServer {
		return errors.New("cannot delete server using AWS server manager, invalid type")
	}

	stopReason := fmt.Sprintf("Stopped by server relates to %s", config.dynamoTable)
	_, err := ecsClient.StopTask(&ecs.StopTaskInput{
		Task:    &awsServer.TaskArn,
		Cluster: &config.ecsCluster,
		Reason:  &stopReason,
	})

	return err
}

// GetServerCount returns the number of tasks running in the ECS cluster
func (manager *ServerManager) GetServerCount() (uint, error) {
	result, err := ecsClient.DescribeClusters(&ecs.DescribeClustersInput{
		Clusters: []*string{&config.ecsCluster},
	})
	if err != nil {
		return 0, err
	}

	return uint(*result.Clusters[0].RunningTasksCount), nil
}

// CreateServer runs a new ECS task with an associated identifier
func (manager *ServerManager) CreateServer(identifier string) (core.Server, error) {

	subnets, err := manager.GetVpcSubnets()
	if err != nil {
		return nil, err
	}

	result, err := ecsClient.RunTask(&ecs.RunTaskInput{
		Cluster:        &config.ecsCluster,
		TaskDefinition: &config.taskDefinition,
		StartedBy:      core.Strpointer("TetraControl"),
		LaunchType:     core.Strpointer("FARGATE"),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				Subnets:        subnets,
				SecurityGroups: config.SecurityGroups,
				AssignPublicIp: core.Strpointer("ENABLED"),
			},
		},
		Tags: []*ecs.Tag{
			{Key: core.Strpointer("identifier"),
				Value: &identifier},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(result.Tasks) <= 0 {
		return nil, errors.New("server created with 0 tasks")
	}

	server := Server{
		Identifier: identifier,
		TaskArn:    *result.Tasks[0].TaskArn,
	}
	if err = manager.AddServerToDatabase(server); err != nil {
		return &server, err
	}

	return &server, err
}

// CleanServerFromDatabase removes server info from Dynamo based on an identifier
func (manager *ServerManager) CleanServerFromDatabase(identifier string) error {
	_, err := dynamodbClient.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &config.dynamoTable,
		Key: map[string]*dynamodb.AttributeValue{
			"identifier": {
				S: aws.String(identifier),
			},
		},
	})
	return err
}

// AddServerToDatabase creates a new entry in Dynamo for a server object
func (manager *ServerManager) AddServerToDatabase(server Server) error {

	data, err := dynamodbattribute.MarshalMap(server)
	if err != nil {
		return err
	}

	_, err = dynamodbClient.PutItem(&dynamodb.PutItemInput{
		Item: data,
	})

	return err
}

// GetVpcSubnets returns a array of all subnets in the associated VPC
func (manager *ServerManager) GetVpcSubnets() ([]*string, error) {
	result, err := ec2Client.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   core.Strpointer("vpc-id"),
				Values: []*string{&config.VpcID},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	foundSubnets := []*string{}
	for _, subnet := range result.Subnets {
		foundSubnets = append(foundSubnets, subnet.SubnetId)
	}

	return foundSubnets, nil
}
