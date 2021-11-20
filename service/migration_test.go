package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/allisson/postmand/mocks"
)

func TestMigration(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	migrationRepository := &mocks.MigrationRepository{}
	migrationService := NewMigration(migrationRepository, logger)
	migrationRepository.On("Run", mock.Anything).Return(nil)
	ctx := context.Background()
	err := migrationService.Run(ctx)
	assert.Nil(t, err)
	migrationRepository.AssertExpectations(t)
}
