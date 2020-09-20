package main

import (
	"log"
	"os"
)

type config struct {
	BucketName string
}

func (c *config) ParseEnv() {
	c.BucketName = os.Getenv("BUCKET_NAME")

	if c.BucketName == "" {
		log.Fatal("bucket name must be specified")
	}
}
