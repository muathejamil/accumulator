package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type AWSSNSPublisher struct {
	Client SNSPublisher
}

func (p *AWSSNSPublisher) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	return p.Client.Publish(ctx, params)
}
