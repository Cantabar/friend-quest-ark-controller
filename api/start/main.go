package main

import (
	"bytes"
	"encoding/json"

	"github.com/Cantabar/friend-quest-ark-controller/core"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse

type BodyRequest struct {
	AcsHost string `json:"acs-host"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request Request) (Response, error) {
	bodyRequest := BodyRequest{
		AcsHost: "",
	}

	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		return Response{Body: err.Error(), StatusCode: 404}, nil
	}

	instances := core.GetAcsInstances()
	instanceIndex := core.FindInstanceByAcsHost(instances, bodyRequest.AcsHost)
	var message string
	if instanceIndex != -1 {
		core.StartAcsInstance(instances[instanceIndex])
		message = "Instance is starting"
	} else {
		message = "Unable to find host"
	}

	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"message": message,
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Content-Type":                     "application/json",
			"X-ACS":                            "start-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
