// +build !aws

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/sbogacz/going-serverless/04_build_tags/internal/toy"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
)

func serve(c *cli.Context) error {
	store, cleanup, err := getStore()
	if err != nil {
		log.WithError(err).Fatal("failed to initialize a store")
	}
	defer cleanup()

	s = toy.New(config, store)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go s.Start()

	<-sigs
	s.Stop()
	return nil
}

func getStore() (*blob.Bucket, func(), error) {
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
