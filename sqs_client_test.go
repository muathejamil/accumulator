package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// MockSQSReceiver is a mock type for the SQSReceiver interface
type MockSQSReceiver struct {
	mock.Mock
}

func (m *MockSQSReceiver) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sqs.ReceiveMessageOutput), args.Error(1)
}

func (m *MockSQSReceiver) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sqs.DeleteMessageOutput), args.Error(1)
}

func TestAWSSQSReceiver_ReceiveMessage(t *testing.T) {
	mockSQS := new(MockSQSReceiver)
	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789012/myQueue"

	mockSQS.On("ReceiveMessage", mock.Anything, mock.MatchedBy(func(params *sqs.ReceiveMessageInput) bool {
		return *params.QueueUrl == queueURL
	})).Return(&sqs.ReceiveMessageOutput{
		Messages: []types.Message{{Body: aws.String("Test Message")}},
	}, nil)

	receiver := AWSSQSReceiver{
		Client: mockSQS,
	}

	_, err := receiver.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl: &queueURL,
	})

	assert.NoError(t, err)
	mockSQS.AssertExpectations(t)
}

func TestAWSSQSReceiver_DeleteMessage(t *testing.T) {
	mockSQS := new(MockSQSReceiver)
	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789012/myQueue"
	receiptHandle := "testReceiptHandle"

	mockSQS.On("DeleteMessage", mock.Anything, mock.MatchedBy(func(params *sqs.DeleteMessageInput) bool {
		return *params.QueueUrl == queueURL && *params.ReceiptHandle == receiptHandle
	})).Return(&sqs.DeleteMessageOutput{}, nil)

	receiver := AWSSQSReceiver{
		Client: mockSQS,
	}

	_, err := receiver.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: &receiptHandle,
	})

	assert.NoError(t, err)
	mockSQS.AssertExpectations(t)
}
