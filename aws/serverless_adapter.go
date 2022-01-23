package aws

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// ServerlessAdapter implements functions to allow TetraControl to run serverlessly
type ServerlessAdapter struct{}

// IsServerless returns true, if running in Lambda
func (serverless ServerlessAdapter) IsServerless() bool {
	_, isAwsLambda := os.LookupEnv("AWS_EXECUTION_ENV")
	return isAwsLambda
}

// Start intializes and starts up the mux adapter
func (serverless ServerlessAdapter) Start() {
	intializeAws()
	adapter := baseServerlessAdapter{}
	adapter.Start()
}

// baseServerlessAdapter implements functions to allow TetraControl to run serverlessly
type baseServerlessAdapter struct{}

// Start passes core.Router to AWS Lambda
func (adapter *baseServerlessAdapter) Start() {
	lambda.StartHandler(adapter)
}

// Invoke works as the entrypoint to the lambda function and marshalls the event
func (adapter *baseServerlessAdapter) Invoke(ctx context.Context, event []byte) ([]byte, error) {

	request := &events.APIGatewayProxyRequest{}
	if err := json.Unmarshal(event, &request); err == nil {
		if len(request.Path) > 0 {
			result, err := ServerlessAPIHandler(ctx, *request)
			jsonResult, _ := json.Marshal(result)
			return jsonResult, err
		}
	}

	stateEvent := &ECSInstanceStateEvent{}
	if err := json.Unmarshal(event, &stateEvent); err == nil {
		if len(stateEvent.Detail.DesiredStatus) > 0 {
			result, err := ServerlessECSHandler(ctx, *stateEvent)
			jsonResult, _ := json.Marshal(result)
			return jsonResult, err
		}
	}

	err := errors.New("could not find proper event handler")
	return []byte(err.Error()), err
}
