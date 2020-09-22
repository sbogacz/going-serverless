package main

import (
	"flag"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/sbogacz/gophercon18-kickoff-talk/third/internal/s3store"
	"github.com/sbogacz/gophercon18-kickoff-talk/third/internal/toy"
)

var (
	config = &toy.Config{}
	s      *toy.Server

	chiLambda *chiadapter.ChiLambda
)

// Handler satisfies the AWS Lambda Go interface
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// If no name is provided in the HTTP request body, throw an error
	return chiLambda.Proxy(req)
}

func main() {
	for _, f := range config.Flags() {
		f.Apply(flag.CommandLine)
	}
	store, err := getStore()
	if err != nil {
		log.WithError(err).Fatal("failed to initialize a store")
	}
	s = toy.New(config, store)
	go s.Start()

	chiLambda = chiadapter.New(s.Router())

	lambda.Start(Handler)
}

func getStore() (toy.Store, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}
	return s3store.New(s3.New(cfg), config.BucketName), nil
}
