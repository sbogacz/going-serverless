package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
	"github.com/sbogacz/gophercon18-kickoff-talk/third/internal/s3store"
	"github.com/sbogacz/gophercon18-kickoff-talk/third/internal/toy"
	"github.com/urfave/cli"
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

	store, err := getStore()
	if err != nil {
		return errors.Wrap(err, "failed to initialize a store")
	}
	s = toy.New(config, store)
	go s.Start()

	<-sigs
	s.Stop()
	return nil
}

func getStore() (toy.Store, error) {
	if localStore {
		return toy.NewLocalStore(), nil
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}
	return s3store.New(s3.New(cfg), config.BucketName), nil
}
