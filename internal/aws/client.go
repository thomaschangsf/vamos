package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Client represents an AWS client
type Client struct {
	s3Client *s3.Client
}

// NewClient creates a new AWS client
func NewClient(ctx context.Context, region string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	return &Client{
		s3Client: s3.NewFromConfig(cfg),
	}, nil
}

// ListBuckets lists all S3 buckets
func (c *Client) ListBuckets(ctx context.Context) ([]string, error) {
	result, err := c.s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %v", err)
	}

	buckets := make([]string, len(result.Buckets))
	for i, bucket := range result.Buckets {
		buckets[i] = *bucket.Name
	}

	return buckets, nil
}
