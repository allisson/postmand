package handler

import "net/http"

var errorResponses = map[string]errorResponse{
	"internal_server_error": {
		Code:       1,
		Message:    "internal server error",
		StatusCode: http.StatusInternalServerError,
	},
	"invalid_id": {
		Code:       2,
		Message:    "invalid id",
		StatusCode: http.StatusNotFound,
	},
	"malformed_request_body": {
		Code:       3,
		Message:    "malformed request body",
		StatusCode: http.StatusBadRequest,
	},
	"request_validation_failed": {
		Code:       4,
		Message:    "request validation failed",
		StatusCode: http.StatusBadRequest,
	},
	"webhook_not_found": {
		Code:       5,
		Message:    "webhook not found",
		StatusCode: http.StatusNotFound,
	},
	"delivery_not_found": {
		Code:       6,
		Message:    "delivery not found",
		StatusCode: http.StatusNotFound,
	},
}

type errorResponse struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	StatusCode int    `json:"-"`
}
