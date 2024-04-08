package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/crypitor/postmand/mocks"
)

func TestPing(t *testing.T) {
	pingRepository := &mocks.PingRepository{}
	pingService := NewPing(pingRepository)
	pingRepository.On("Run", mock.Anything).Return(nil)
	ctx := context.Background()
	err := pingService.Run(ctx)
	assert.Nil(t, err)
	pingRepository.AssertExpectations(t)
}
