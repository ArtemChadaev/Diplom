package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockMailer is a testify/mock implementation of mailer.Mailer.
type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) SendOTP(ctx context.Context, toEmail, code string) error {
	args := m.Called(ctx, toEmail, code)
	return args.Error(0)
}
