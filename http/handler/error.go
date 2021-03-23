package handler

import "net/http"

type errorResponseCode int

const (
	internalServerErrorCode errorResponseCode = iota + 1
	invalidIDCode
	malformedRequestBodyCode
	requestValidationFailedCode
	webhookNotFoundCode
	deliveryNotFoundCode
	deliveryAttemptNotFoundCode
)

var errorResponses = map[string]errorResponse{
	"internal_server_error": {
		Code:       internalServerErrorCode,
		Message:    "internal server error",
		StatusCode: http.StatusInternalServerError,
	},
	"invalid_id": {
		Code:       invalidIDCode,
		Message:    "invalid id",
		StatusCode: http.StatusNotFound,
	},
	"malformed_request_body": {
		Code:       malformedRequestBodyCode,
		Message:    "malformed request body",
		StatusCode: http.StatusBadRequest,
	},
	"request_validation_failed": {
		Code:       requestValidationFailedCode,
		Message:    "request validation failed",
		StatusCode: http.StatusBadRequest,
	},
	"webhook_not_found": {
		Code:       webhookNotFoundCode,
		Message:    "webhook not found",
		StatusCode: http.StatusNotFound,
	},
	"delivery_not_found": {
		Code:       deliveryNotFoundCode,
		Message:    "delivery not found",
		StatusCode: http.StatusNotFound,
	},
	"delivery_attempt_not_found": {
		Code:       deliveryAttemptNotFoundCode,
		Message:    "delivery attempt not found",
		StatusCode: http.StatusNotFound,
	},
}

type errorResponse struct {
	Code       errorResponseCode `json:"code"`
	Message    string            `json:"message"`
	Details    string            `json:"details,omitempty"`
	StatusCode int               `json:"-"`
} //@name Error
