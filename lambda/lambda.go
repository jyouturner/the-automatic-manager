package Lambda

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

//LambaResponse specifies the data structure of Lambda output to API gateway.
type LambdaResponse struct {
	StatusCode int         `json:"statusCode"`
	Body       string      `json:"body"`
	Headers    HeadersType `json:"headers"`
}

type HeadersType map[string]interface{}

type Response struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//Success() create the message for API gateway with code 200
func Success(message string) LambdaResponse {
	return LambdaResponse{
		StatusCode: 200,
		Body:       message,
	}
}

//redirect
func Redirect(code int, url string) LambdaResponse {
	return LambdaResponse{
		StatusCode: code,
		Body:       "",
		Headers: HeadersType{
			"Location": url,
		},
	}

}

//FailureMessage create a message for API gateway with non-2xx code
func FailureMessage(code int, message string) LambdaResponse {

	return LambdaResponse{
		StatusCode: code,
		Body:       message,
	}

}

type ValidateBodyFunctionType func(string) error

type Handler struct {
	ValidateBodyFunction ValidateBodyFunctionType
}

//Handler is the Lambda handler, it does (very) basic validation, then pass the event to SQS queue to be processed.
func (p *Handler) HandleExample(ctx context.Context, event json.RawMessage) (LambdaResponse, error) {
	eventBody := gjson.Get(string(event), "body")
	err := p.ValidateBodyFunction(eventBody.String())
	// with AWS API Gateway and Lambda proxy, we should put the "error" in the response statusCode, than raising error from Lambda integration.
	// If we raise error here, all we get is 502 error, which is not desired.
	if err != nil {
		return FailureMessage(400, fmt.Sprintf("invalid event %v", err)), nil
	}
	//err = sendToQueue(eventBody.String(), os.Getenv("QUEUE_NAME"))
	//if err != nil {
	//		log.Errorf("failed to send message to queue %v ", err)
	//	}
	return Success("ok"), nil
}
