package main

import (
	"context"
	"fmt"
	"io/ioutil"
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
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/blob/s3blob"
)

var (
	config     = &toy.Config{}
	s          *toy.Server
	localStore bool
)

func flags() []cli.Flag {
	return append(config.Flags(),
		cli.BoolFlag{
			Name:        "local-store",
			Usage:       "use an in-memory backing store",
			Destination: &localStore,
		})
}

func main() {
	app := cli.NewApp()
	app.Usage = "this is the CLI version of our toy app"
	app.Flags = flags()
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
	if localStore {
		store, cleanup, err = getLocalStore()
	} else {
		store, cleanup, err = getS3Store()
	}
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

func getLocalStore() (*blob.Bucket, func(), error) {
	dir, err := ioutil.TempDir("", "toy-test-files")
	if err != nil {
		return nil, nil, err
	}
	store, err := fileblob.OpenBucket(dir, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize local store: %w", err)
	}
	return store, func() { os.RemoveAll(dir) }, nil
}

func getS3Store() (*blob.Bucket, func(), error) {
	noop := func() {
		return
	}
	cfg := &aws.Config{
		Credentials: credentials.NewEnvCredentials(),
	}
	sess := session.Must(session.NewSession(cfg))
	store, err := s3blob.OpenBucket(context.TODO(), sess, config.BucketName, nil)
	if err != nil {
		return nil, nil, err
	}
	return store, noop, nil
}
