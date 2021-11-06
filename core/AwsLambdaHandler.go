package core

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	muxadapter "github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
)

var muxLambda *muxadapter.GorillaMuxAdapter

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return muxLambda.ProxyWithContext(ctx, req)
}

// AwsLambdaHandler passes core.Router to AWS Lambda
func AwsLambdaHandler() {
	muxLambda = muxadapter.New(Router)
	lambda.Start(handleRequest)
}
