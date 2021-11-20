package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/allisson/postmand"
	"github.com/allisson/postmand/mocks"
)

func TestWorker(t *testing.T) {
	ctx := context.Background()
	pollingInterval := 10 * time.Millisecond

	t.Run("run with dispatch error", func(t *testing.T) {
		deliveryRepository := &mocks.DeliveryRepository{}
		logger, _ := zap.NewDevelopment()
		workerService := NewWorker(deliveryRepository, logger, pollingInterval)

		deliveryRepository.On("Dispatch", mock.Anything).Return(nil, errors.New("error"))
		// Wait 15 miliseconds before call shutdown.
		go func() {
			workerService.Shutdown(ctx)
		}()
		workerService.run(ctx)

		deliveryRepository.AssertExpectations(t)
	})

	t.Run("run with no dispatch", func(t *testing.T) {
		deliveryRepository := &mocks.DeliveryRepository{}
		logger, _ := zap.NewDevelopment()
		workerService := NewWorker(deliveryRepository, logger, pollingInterval)

		deliveryRepository.On("Dispatch", mock.Anything).Return(nil, nil)
		// Wait 15 miliseconds before call shutdown.
		go func() {
			workerService.Shutdown(ctx)
		}()
		workerService.run(ctx)

		deliveryRepository.AssertExpectations(t)
	})

	t.Run("run with dispatch", func(t *testing.T) {
		deliveryRepository := &mocks.DeliveryRepository{}
		logger, _ := zap.NewDevelopment()
		workerService := NewWorker(deliveryRepository, logger, pollingInterval)

		deliveryRepository.On("Dispatch", mock.Anything).Return(&postmand.DeliveryAttempt{}, nil)
		// Wait 15 miliseconds before call shutdown.
		go func() {
			workerService.Shutdown(ctx)
		}()
		workerService.run(ctx)

		deliveryRepository.AssertExpectations(t)
	})
}
