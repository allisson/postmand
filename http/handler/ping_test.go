package handler

import (
	"errors"
	nethttp "net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/allisson/postmand/http"
	"github.com/allisson/postmand/mocks"
)

func TestPing(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("With success", func(t *testing.T) {
		pingService := &mocks.PingService{}
		pingHandler := NewPing(pingService, logger)
		router := http.NewRouter(logger)
		router.Get("/healthz", pingHandler.Healthz)

		pingService.On("Run", mock.Anything).Return(nil)
		apitest.New().
			Handler(router).
			Get("/healthz").
			Expect(t).
			Body(`{"success":true}`).
			Status(nethttp.StatusOK).
			End()

		pingService.AssertExpectations(t)
	})

	t.Run("With error", func(t *testing.T) {
		pingService := &mocks.PingService{}
		pingHandler := NewPing(pingService, logger)
		router := http.NewRouter(logger)
		router.Get("/healthz", pingHandler.Healthz)

		pingService.On("Run", mock.Anything).Return(errors.New("BOOM"))
		apitest.New().
			Handler(router).
			Get("/healthz").
			Expect(t).
			Body(`{"success":false}`).
			Status(nethttp.StatusInternalServerError).
			End()

		pingService.AssertExpectations(t)
	})
}
