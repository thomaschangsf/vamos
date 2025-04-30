package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockS3Client implements the S3ClientInterface
type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.ListBucketsOutput), args.Error(1)
}

func TestNewClient(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient(ctx, "us-west-2")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.s3Client)
}

func TestListBuckets(t *testing.T) {
	tests := []struct {
		name          string
		buckets       []types.Bucket
		expectedError error
	}{
		{
			name: "Success with buckets",
			buckets: []types.Bucket{
				{Name: aws.String("bucket1")},
				{Name: aws.String("bucket2")},
			},
			expectedError: nil,
		},
		{
			name:          "Success with no buckets",
			buckets:       []types.Bucket{},
			expectedError: nil,
		},
		{
			name:          "Error from API",
			buckets:       nil,
			expectedError: errors.New("API error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockS3 := new(MockS3Client)
			client := &Client{
				s3Client: mockS3,
			}

			output := &s3.ListBucketsOutput{
				Buckets: tt.buckets,
			}
			mockS3.On("ListBuckets", mock.Anything, mock.Anything, mock.Anything).Return(output, tt.expectedError)

			buckets, err := client.ListBuckets(context.Background())
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.buckets, buckets)
			}

			mockS3.AssertExpectations(t)
		})
	}
}
