package s3store

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
	"github.com/sbogacz/gophercon18-kickoff-talk/third/internal/httperrs"
)

// S3Store is an S3 backed implementation of our Store
// interaface
type S3Store struct {
	client *s3.S3
	bucket string
}

// New takes an S3 client and a bucket name, and
// returns the S3 implementation of the Store interface
func New(client *s3.S3, bucket string) *S3Store {
	return &S3Store{
		client: client,
		bucket: bucket,
	}
}

// Set takes key and data as a string, saves it to S3, and returns the uuid
// under which the file is stored
func (ss *S3Store) Set(ctx context.Context, key, data string) error {
	input := putInput(ss.bucket, key, data)

	putReq := ss.client.PutObjectRequest(input)
	if _, err := putReq.Send(); err != nil {
		return errors.Wrap(err, "failed to put file in S3")
	}
	return nil
}

// Get takes a key, and attempts to fetch the data from S3
func (ss *S3Store) Get(ctx context.Context, key string) (string, error) {
	input := getInput(ss.bucket, key)
	getReq := ss.client.GetObjectRequest(input)
	resp, err := getReq.Send()
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchKey {
				return "", httperrs.NotFound(err, "no such file")
			}
		}
		return "", httperrs.InternalServer(err, "failed to get file from S3")
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", httperrs.InternalServer(err, "failed to read object body")
	}
	return string(b), nil
}

// Del takes a key, and attempts to delete the data from S3
func (ss *S3Store) Del(ctx context.Context, key string) error {
	input := deleteInput(ss.bucket, key)
	deleteReq := ss.client.DeleteObjectRequest(input)
	_, err := deleteReq.Send()
	if err != nil {
		return httperrs.InternalServer(err, "failed to delete file from S3")
	}
	return nil
}

func putInput(bucket, key, data string) *s3.PutObjectInput {
	return &s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(key),
		ServerSideEncryption: s3.ServerSideEncryptionAes256,
		Body:                 aws.ReadSeekCloser(strings.NewReader(data)),
	}
}

func getInput(bucket, key string) *s3.GetObjectInput {
	return &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
}

func deleteInput(bucket, key string) *s3.DeleteObjectInput {
	return &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
}
