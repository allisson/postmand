// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// WorkerService is an autogenerated mock type for the WorkerService type
type WorkerService struct {
	mock.Mock
}

// Run provides a mock function with given fields: ctx
func (_m *WorkerService) Run(ctx context.Context) {
	_m.Called(ctx)
}

// Shutdown provides a mock function with given fields: ctx
func (_m *WorkerService) Shutdown(ctx context.Context) {
	_m.Called(ctx)
}