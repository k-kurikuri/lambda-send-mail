package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"strconv"
)

type Request struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	Message string `json:"message"`
}

func Handler(request Request) (Response, error) {
	id := request.ID
	name := request.Name
	msgPrefix := strconv.Itoa(id) + " : " + name

	return Response{
		Message: msgPrefix + " Go Serverless v1.0! Your function executed successfully!",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
