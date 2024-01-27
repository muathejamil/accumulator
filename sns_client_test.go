package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// MockSNSPublisher is a mock type for the SNSPublisher interface
type MockSNSPublisher struct {
	mock.Mock
}

func (m *MockSNSPublisher) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sns.PublishOutput), args.Error(1)
}

func TestSNSPublisher_Publish(t *testing.T) {
	// Create an instance of our test object
	mockSNS := new(MockSNSPublisher)

	// Setup expectation
	mockSNS.On("Publish", mock.Anything, mock.AnythingOfType("*sns.PublishInput")).Return(&sns.PublishOutput{
		MessageId: aws.String("testMessageId"),
	}, nil)

	publisher := AWSSNSPublisher{
		Client: mockSNS,
	}

	// Test the Publish method
	_, err := publisher.Publish(context.TODO(), &sns.PublishInput{
		Message:  aws.String("Hello, world!"),
		TopicArn: aws.String("arn:aws:sns:us-east-1:123456789012:myTopic"),
	})

	// Assertions
	assert.NoError(t, err)
	mockSNS.AssertExpectations(t) // Assert that Publish was called as expected
}
