// package goingserverless is here instead of package main, because GCP actually requires
// function code to be at the root of the directory. To achieve this, we rely on our
// Makefule
package goingserverless

import (
	"context"
	"flag"
	"net/http"

	"github.com/sbogacz/going-serverless/03_separate_binaries_and_gocloud/internal/toy"
	log "github.com/sirupsen/logrus"
	"gocloud.dev/blob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/gcp"
)

var (
	config = &toy.Config{}
	s      *toy.Server //nolint:unused
)

func init() {
	for _, f := range config.Flags() {
		f.Apply(flag.CommandLine)
	}
	store, err := getStore()
	if err != nil {
		log.WithError(err).Fatal("failed to initialize a store")
	}
	s = toy.New(config, store)
}

// Handle is the function that will be called by GCP function trigger
func Handle(w http.ResponseWriter, r *http.Request) { //nolint:deadcode,unused
	s.Router.ServeHTTP(w, r)
}

func getStore() (*blob.Bucket, error) {
	ctx := context.Background()

	creds, err := gcp.DefaultCredentials(ctx)
	if err != nil {
		return nil, err
	}

	client, err := gcp.NewHTTPClient(
		gcp.DefaultTransport(),
		gcp.CredentialsTokenSource(creds))
	if err != nil {
		return nil, err
	}

	bucket, err := gcsblob.OpenBucket(ctx, client, config.BucketName, nil)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}
