package main

import (
	"context"
	"flag"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"snssqsintegration/mocks"
	"testing"
)

// TestMain_IsNumber verifies the IsNumber function's ability to distinguish numeric strings from non-numeric ones.
func TestMain_IsNumber(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"42", true},
		{"-42", true},
		{"0", true},
		{"notANumber", false},
		{"123abc", false},
		{"", false},
	}

	for _, testCase := range testCases {
		if IsNumber(testCase.input) != testCase.expected {
			t.Errorf("IsNumber(%s) expected %t, got %t", testCase.input, testCase.expected, !testCase.expected)
		}
	}
}

// TestMain_ReadAndValidateArguments tests the ReadAndValidateArguments function with various sets of command-line arguments.
func TestMain_ReadAndValidateArguments(t *testing.T) {
	// Define a helper function to simulate command-line argument parsing.
	setCommandLineArgs := func(args []string) {
		os.Args = args
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}

	// Test cases
	tests := []struct {
		name          string
		args          []string
		expectErr     bool
		expectedQueue string
		expectedArn   string
	}{
		{
			"Valid Arguments",
			[]string{"cmd", "-q", "testQueue", "-t", "testTopic", "-a", "testArn"},
			false,
			"testQueue",
			"testArn",
		},
		{
			"Missing Queue",
			[]string{"cmd", "-t", "testTopic", "-a", "testArn"},
			true,
			"",
			"",
		},
		{
			"Missing Topic Arn",
			[]string{"cmd", "-q", "testQueue", "-t", "testTopic"},
			true,
			"",
			"",
		},
		// Add more test cases as needed.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setCommandLineArgs(tt.args)
			queue, topicArn, err := ReadAndValidateArguments()

			if (err != nil) != tt.expectErr {
				t.Errorf("ReadAndValidateArguments() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if queue != nil && *queue != tt.expectedQueue {
				t.Errorf("Expected queue = %s, got %s", tt.expectedQueue, *queue)
			}
			if topicArn != nil && *topicArn != tt.expectedArn {
				t.Errorf("Expected topicArn = %s, got %s", tt.expectedArn, *topicArn)
			}
		})
	}
}

func TestMain_Publish(t *testing.T) {
	input := "test message"
	topicArn := "arn:aws:sns:us-east-1:123456789012:myTopic"

	// Create an instance of our test object
	mockSNSPublisher := new(mocks.MockSNSPublisher)

	// Setup expectations
	mockSNSPublisher.On("Publish", mock.Anything, mock.MatchedBy(func(params *sns.PublishInput) bool {
		return *params.Message == input && *params.TopicArn == topicArn
	})).Return(&sns.PublishOutput{MessageId: aws.String("mockMessageId")}, nil)

	// Call the function under test
	err := Publish(mockSNSPublisher, input, &topicArn)

	// Assert expectations
	mockSNSPublisher.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestMain_ReceiveAndDeleteMessage(t *testing.T) {
	ctx := context.TODO()
	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789012/myQueue"

	// Mock setup
	mockSQSReceiver := new(mocks.MockSQSReceiver)

	// Mock ReceiveMessage expectation
	mockSQSReceiver.On("ReceiveMessage", mock.Anything, mock.MatchedBy(func(params *sqs.ReceiveMessageInput) bool {
		return *params.QueueUrl == queueURL
	})).Return(&sqs.ReceiveMessageOutput{
		Messages: []types.Message{{Body: aws.String("Test Message")}},
	}, nil)

	// Mock DeleteMessage expectation
	receiptHandle := "mockReceiptHandle"
	mockSQSReceiver.On("DeleteMessage", mock.Anything, mock.MatchedBy(func(params *sqs.DeleteMessageInput) bool {
		return *params.QueueUrl == queueURL && *params.ReceiptHandle == receiptHandle
	})).Return(&sqs.DeleteMessageOutput{}, nil)

	// Test ReceiveMessage
	msgOutput, err := mockSQSReceiver.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueURL),
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, msgOutput.Messages)

	// Test DeleteMessage
	err = DeleteMsg(mockSQSReceiver, aws.String(queueURL), types.Message{ReceiptHandle: aws.String(receiptHandle)})
	assert.NoError(t, err)

	// Verify expectations
	mockSQSReceiver.AssertExpectations(t)
}
