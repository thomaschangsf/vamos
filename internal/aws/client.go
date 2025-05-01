package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// S3ClientInterface defines the interface for S3 operations
type S3ClientInterface interface {
	ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error)
}

// STSClientInterface defines the interface for STS operations
type STSClientInterface interface {
	GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
	GetSessionToken(ctx context.Context, params *sts.GetSessionTokenInput, optFns ...func(*sts.Options)) (*sts.GetSessionTokenOutput, error)
}

// Client represents an AWS client
type Client struct {
	region    string
	s3Client  S3ClientInterface
	stsClient STSClientInterface
}

// NewClient creates a new AWS client
func NewClient(ctx context.Context, region string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	return &Client{
		region:    region,
		s3Client:  s3.NewFromConfig(cfg),
		stsClient: sts.NewFromConfig(cfg),
	}, nil
}

// ListBuckets lists all S3 buckets
func (c *Client) ListBuckets(ctx context.Context) ([]types.Bucket, error) {
	result, err := c.s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	return result.Buckets, nil
}

// GetBucketRegion gets the region for a bucket
func (c *Client) GetBucketRegion(ctx context.Context, bucket string) (string, error) {
	// First try with the current client
	output, err := c.s3Client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return "", err
	}

	// The location constraint can be empty for us-east-1
	if output.LocationConstraint == "" {
		return "us-east-1", nil
	}
	return string(output.LocationConstraint), nil
}

// ListObjects lists objects in a bucket
func (c *Client) ListObjects(ctx context.Context, location string) ([]types.Object, error) {
	// Parse s3://bucket/path format
	bucket, prefix := parseS3Location(location)
	if bucket == "" {
		return nil, fmt.Errorf("invalid S3 location: %s", location)
	}

	// Try to get the bucket's region, but don't fail if we can't
	bucketRegion := c.region
	_, err := c.GetBucketRegion(ctx, bucket)
	if err != nil {
		// If we can't get the bucket region, try us-east-1 first (Common Crawl's region)
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
		if err != nil {
			return nil, fmt.Errorf("failed to create client for us-east-1: %v", err)
		}
		s3Client := s3.NewFromConfig(cfg)

		// Try listing objects with us-east-1 client
		result, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(prefix),
		})
		if err == nil {
			return result.Contents, nil
		}

		// If that fails, try with the user's specified region
		if c.region != "us-east-1" {
			cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(c.region))
			if err != nil {
				return nil, fmt.Errorf("failed to create client for region %s: %v", c.region, err)
			}
			s3Client = s3.NewFromConfig(cfg)

			result, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
				Bucket: aws.String(bucket),
				Prefix: aws.String(prefix),
			})
			if err == nil {
				return result.Contents, nil
			}
		}

		return nil, fmt.Errorf("failed to list objects. The bucket might be in a different region. Try specifying the region with -region flag")
	}

	// If we got the bucket region, use it
	var s3Client S3ClientInterface = c.s3Client
	if bucketRegion != c.region {
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(bucketRegion))
		if err != nil {
			return nil, fmt.Errorf("failed to create client for region %s: %v", bucketRegion, err)
		}
		s3Client = s3.NewFromConfig(cfg)
	}

	// List objects
	result, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, err
	}
	return result.Contents, nil
}

// GetCallerIdentity gets information about the current IAM identity
func (c *Client) GetCallerIdentity(ctx context.Context) (*sts.GetCallerIdentityOutput, error) {
	return c.stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
}

// GetSessionToken gets temporary credentials
func (c *Client) GetSessionToken(ctx context.Context, durationSeconds int32) (*sts.GetSessionTokenOutput, error) {
	return c.stsClient.GetSessionToken(ctx, &sts.GetSessionTokenInput{
		DurationSeconds: &durationSeconds,
	})
}

// parseS3Location parses an S3 location in the format s3://bucket/path
func parseS3Location(location string) (bucket, prefix string) {
	if !strings.HasPrefix(location, "s3://") {
		return "", ""
	}

	parts := strings.SplitN(strings.TrimPrefix(location, "s3://"), "/", 2)
	bucket = parts[0]
	if len(parts) > 1 {
		prefix = parts[1]
	}
	return bucket, prefix
}
