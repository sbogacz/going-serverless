// TODO: The Azure deployment isn't complete yet
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sbogacz/going-serverless/03_separate_binaries_and_gocloud/internal/toy"
	"github.com/urfave/cli"
	"gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"
)

var (
	config = &toy.Config{}
	s      *toy.Server
)

func main() {
	app := cli.NewApp()
	app.Usage = "this is the CLI version of our toy app"
	app.Flags = config.Flags()
	app.Action = serve

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func serve(c *cli.Context) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var store *blob.Bucket
	var cleanup func()
	var err error
	store, cleanup, err = getS3Store()

	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}
	defer cleanup()
	s = toy.New(config, store)
	go s.Start()

	<-sigs
	s.Stop()
	return nil
}

func getS3Store() (*blob.Bucket, func(), error) {
	noop := func() {
		return
	}
	cfg := &aws.Config{
		Credentials: credentials.NewEnvCredentials(),
	}
	sess := session.Must(session.NewSession(cfg))
	store, err := s3blob.OpenBucket(context.TODO(), sess, config.BucketName)
	if err != nil {
		return nil, nil, err
	}
	return store, noop, nil
}
