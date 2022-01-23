package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// ECSInstanceStateEvent holds top level event details
// full event: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs_cwe_events.html
type ECSInstanceStateEvent struct {
	ID     string                 `json:"id"`
	Detail ECSInstanceStateDetail `json:"detail"`
}

// ECSInstanceStateDetail holds detail level info
type ECSInstanceStateDetail struct {
	Status        string `json:"status"`
	DesiredStatus string `json:"desiredStatus"`
	TaskArn       string `json:"taskArn"`
}

// ServerlessECSHandler runs when a game server task changes state
func ServerlessECSHandler(ctx context.Context, event ECSInstanceStateEvent) (interface{}, error) {

	if event.Detail.DesiredStatus == "STOPPED" {

		// Check if item exists
		result, err := dynamodbClient.GetItem(&dynamodb.GetItemInput{
			TableName: &config.DynamoTable,
			Key: map[string]*dynamodb.AttributeValue{
				"task": {
					S: aws.String(event.Detail.TaskArn),
				},
			},
		})
		if err != nil {
			return err.Error(), err
		}

		// Delete item if found
		if result.Item != nil {
			_, err := dynamodbClient.DeleteItem(&dynamodb.DeleteItemInput{
				TableName: &config.DynamoTable,
				Key: map[string]*dynamodb.AttributeValue{
					"task": {
						S: aws.String(event.Detail.TaskArn),
					},
				},
			})
			return fmt.Sprintf("called delete on %v", event.Detail.TaskArn), err
		}
		return fmt.Sprintf("%v does not exist, nothing to do", event.Detail.TaskArn), nil
	}

	return fmt.Sprintf("desired state for %v is %v, ignoring", event.Detail.TaskArn, event.Detail.Status), nil
}
