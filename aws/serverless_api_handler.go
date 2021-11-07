package aws

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	muxadapter "github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/loudsmilestudios/TetraControl/core"
)

var muxLambda *muxadapter.GorillaMuxAdapter

// ServerlessAPIHandler passes the API gateway request event to Mux
func ServerlessAPIHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	muxLambda = muxadapter.New(core.Router)
	return muxLambda.ProxyWithContext(ctx, req)
}
