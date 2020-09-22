package toy

import (
	"github.com/urfave/cli"
)

// Config holds the static config our server needs
type Config struct {
	BucketName string
	Port       int
}

// Flags returns a slice of urfave.Flag to allow
// our configuration to be populated by flags or
// env vars
func (c *Config) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "bucket-name, b",
			EnvVar:      "BUCKET_NAME",
			Usage:       "the bucket in which to store our data",
			Destination: &c.BucketName,
		},
		cli.IntFlag{
			Name:        "port, p",
			EnvVar:      "PORT",
			Usage:       "the port that we want to server requests on",
			Value:       8080,
			Destination: &c.Port,
		},
	}
}
