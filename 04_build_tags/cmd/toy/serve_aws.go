// +build aws

package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/sbogacz/going-serverless/04_build_tags/internal/toy"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"
)

var (
	chiLambda *chiadapter.ChiLambda
)

// Handler satisfies the AWS Lambda Go interface
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return chiLambda.Proxy(req)
}

func serve(ctx *cli.Context) error {
	store, cleanup, err := getStore()
	if err != nil {
		log.WithError(err).Fatal("failed to initialize a store")
	}
	defer cleanup()

	s = toy.New(config, store)
	go s.Start()

	chiLambda = chiadapter.New(s.Router)

	lambda.Start(Handler)

	s.Stop()

	return nil
}

func getStore() (*blob.Bucket, func(), error) {
	cfg := &aws.Config{
		Credentials: credentials.NewEnvCredentials(),
	}
	sess := session.Must(session.NewSession(cfg))
	store, err := s3blob.OpenBucket(context.TODO(), sess, config.BucketName, nil)
	if err != nil {
		return nil, nil, err
	}
	return store, nil, nil
}
