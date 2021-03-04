package service

import (
	"context"
	"testing"

	"github.com/allisson/postmand/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMigration(t *testing.T) {
	migrationRepository := &mocks.MigrationRepository{}
	migrationService := NewMigration(migrationRepository)
	migrationRepository.On("Run", mock.Anything).Return(nil)
	ctx := context.Background()
	err := migrationService.Run(ctx)
	assert.Nil(t, err)
	migrationRepository.AssertExpectations(t)
}
