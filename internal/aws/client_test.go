package aws

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	s3_types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	sts_types "github.com/aws/aws-sdk-go-v2/service/sts/types"
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

func (m *MockS3Client) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.ListObjectsV2Output), args.Error(1)
}

func (m *MockS3Client) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.GetBucketLocationOutput), args.Error(1)
}

// MockSTSClient implements the STSClientInterface
type MockSTSClient struct {
	mock.Mock
}

func (m *MockSTSClient) GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sts.GetCallerIdentityOutput), args.Error(1)
}

func (m *MockSTSClient) GetSessionToken(ctx context.Context, params *sts.GetSessionTokenInput, optFns ...func(*sts.Options)) (*sts.GetSessionTokenOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sts.GetSessionTokenOutput), args.Error(1)
}

func TestNewClient(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient(ctx, "us-west-2")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.s3Client)
	assert.NotNil(t, client.stsClient)
}

func TestListBuckets(t *testing.T) {
	tests := []struct {
		name          string
		buckets       []s3_types.Bucket
		expectedError error
	}{
		{
			name: "Success with buckets",
			buckets: []s3_types.Bucket{
				{Name: aws.String("bucket1")},
				{Name: aws.String("bucket2")},
			},
			expectedError: nil,
		},
		{
			name:          "Success with no buckets",
			buckets:       []s3_types.Bucket{},
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

func TestListObjects(t *testing.T) {
	ctx := context.Background()
	client, _ := NewClient(ctx, "us-west-2")

	tests := []struct {
		name          string
		location      string
		expectedError string
		objects       []types.Object
		setupMock     func(*MockS3Client)
	}{
		{
			name:     "Success with objects",
			location: "s3://my-bucket/path/",
			objects: []types.Object{
				{
					Key:          aws.String("path/file1.txt"),
					Size:         aws.Int64(1024),
					LastModified: aws.Time(time.Now()),
				},
			},
			setupMock: func(m *MockS3Client) {
				m.On("GetBucketLocation", mock.Anything, mock.Anything, mock.Anything).
					Return(&s3.GetBucketLocationOutput{
						LocationConstraint: types.BucketLocationConstraint("us-west-2"),
					}, nil)
				m.On("ListObjectsV2", mock.Anything, mock.Anything, mock.Anything).
					Return(&s3.ListObjectsV2Output{
						Contents: []types.Object{
							{
								Key:          aws.String("path/file1.txt"),
								Size:         aws.Int64(1024),
								LastModified: aws.Time(time.Now()),
							},
						},
					}, nil)
			},
		},
		{
			name:          "Invalid location",
			location:      "invalid-location",
			expectedError: "invalid S3 location: invalid-location",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockS3 := &MockS3Client{}
			if tt.setupMock != nil {
				tt.setupMock(mockS3)
			}
			client.s3Client = mockS3

			objects, err := client.ListObjects(ctx, tt.location)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.objects), len(objects))
				if len(objects) > 0 {
					assert.Equal(t, tt.objects[0].Key, objects[0].Key)
				}
			}

			mockS3.AssertExpectations(t)
		})
	}
}

func TestGetCallerIdentity(t *testing.T) {
	tests := []struct {
		name          string
		identity      *sts.GetCallerIdentityOutput
		expectedError error
	}{
		{
			name: "Success",
			identity: &sts.GetCallerIdentityOutput{
				Account: aws.String("123456789012"),
				Arn:     aws.String("arn:aws:iam::123456789012:user/test-user"),
				UserId:  aws.String("AIDATESTUSERID"),
			},
			expectedError: nil,
		},
		{
			name:          "Error from API",
			identity:      nil,
			expectedError: errors.New("API error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSTS := new(MockSTSClient)
			client := &Client{
				stsClient: mockSTS,
			}

			mockSTS.On("GetCallerIdentity", mock.Anything, mock.Anything, mock.Anything).Return(tt.identity, tt.expectedError)

			identity, err := client.GetCallerIdentity(context.Background())
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.identity, identity)
			}

			mockSTS.AssertExpectations(t)
		})
	}
}

func TestGetSessionToken(t *testing.T) {
	tests := []struct {
		name          string
		credentials   *sts.GetSessionTokenOutput
		expectedError error
	}{
		{
			name: "Success",
			credentials: &sts.GetSessionTokenOutput{
				Credentials: &sts_types.Credentials{
					AccessKeyId:     aws.String("AKIAEXAMPLE"),
					SecretAccessKey: aws.String("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
					SessionToken:    aws.String("AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKwRcOIfrRh3c/LTo6UDdyJwOOvEVPvLXCrrrUtdnniCEXAMPLE/IvU1dYUg2RVAJBanLiHb4IgRmpRV3zrkuWJOgQs8IZZaIv2BXIa2R4OlgkBN9bkUDNCJiBeb/AXlzBBko7b15fjrBs2+cTQtpZ3CYWFXG8C5zqx37wnOE49mRl/+OtkIKGO7fAE"),
					Expiration:      aws.Time(time.Now().Add(12 * time.Hour)),
				},
			},
			expectedError: nil,
		},
		{
			name:          "Error from API",
			credentials:   nil,
			expectedError: errors.New("API error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSTS := new(MockSTSClient)
			client := &Client{
				stsClient: mockSTS,
			}

			mockSTS.On("GetSessionToken", mock.Anything, mock.Anything, mock.Anything).Return(tt.credentials, tt.expectedError)

			credentials, err := client.GetSessionToken(context.Background(), 3600)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.credentials, credentials)
			}

			mockSTS.AssertExpectations(t)
		})
	}
}

func TestGetBucketRegion(t *testing.T) {
	tests := []struct {
		name           string
		bucket         string
		location       string
		expectedRegion string
		expectedError  error
	}{
		{
			name:           "Success with us-east-1",
			bucket:         "my-bucket",
			location:       "",
			expectedRegion: "us-east-1",
			expectedError:  nil,
		},
		{
			name:           "Success with us-west-2",
			bucket:         "my-bucket",
			location:       "us-west-2",
			expectedRegion: "us-west-2",
			expectedError:  nil,
		},
		{
			name:           "Error from API",
			bucket:         "my-bucket",
			location:       "",
			expectedRegion: "",
			expectedError:  errors.New("API error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockS3 := new(MockS3Client)
			client := &Client{
				s3Client: mockS3,
				region:   "us-west-2",
			}

			output := &s3.GetBucketLocationOutput{
				LocationConstraint: types.BucketLocationConstraint(tt.location),
			}
			mockS3.On("GetBucketLocation", mock.Anything, &s3.GetBucketLocationInput{
				Bucket: aws.String(tt.bucket),
			}, mock.Anything).Return(output, tt.expectedError)

			region, err := client.GetBucketRegion(context.Background(), tt.bucket)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRegion, region)
			}

			mockS3.AssertExpectations(t)
		})
	}
}
