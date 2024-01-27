package mocks

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/mock"
)

type MockSNSPublisher struct {
	mock.Mock
}

func (m *MockSNSPublisher) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sns.PublishOutput), args.Error(1)
}
