package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	// TEST OMIT
	"github.com/stretchr/testify/require" // HL
)

func TestHelpers(t *testing.T) {
	t.Run("extract key should parse the path correctly", func(t *testing.T) {
		req := &events.APIGatewayProxyRequest{
			Path: "blobs/1234",
		}
		key, err := extractKey(req)
		require.NoError(t, err)
		require.Equal(t, "1234", key)
	})
}
