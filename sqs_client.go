package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type AWSSQSReceiver struct {
	Client SQSReceiver
}

func (r *AWSSQSReceiver) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	return r.Client.ReceiveMessage(ctx, params)
}

func (r *AWSSQSReceiver) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	return r.Client.DeleteMessage(ctx, params)
}
