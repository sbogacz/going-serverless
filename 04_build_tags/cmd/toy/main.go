package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/sbogacz/going-serverless/04_build_tags/internal/toy"
	"github.com/urfave/cli"
)

var (
	config = &toy.Config{}
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
