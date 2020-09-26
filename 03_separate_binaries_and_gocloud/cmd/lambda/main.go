package main

import (
	"context"
	"flag"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/sbogacz/going-serverless/03_separate_binaries_and_gocloud/internal/toy"
	log "github.com/sirupsen/logrus"
	"gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"
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

	chiLambda = chiadapter.New(s.Router)

	lambda.Start(Handler)
}

func getStore() (*blob.Bucket, error) {
	cfg := &aws.Config{
		Credentials: credentials.NewEnvCredentials(),
	}
	sess := session.Must(session.NewSession(cfg))
	store, err := s3blob.OpenBucket(context.TODO(), sess, config.BucketName, nil)
	if err != nil {
		return nil, err
	}
	return store, nil
}
