package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/allisson/postmand"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.uber.org/zap"
)

func makeResponse(w http.ResponseWriter, body []byte, statusCode int, contentType string, logger *zap.Logger) {
	w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", contentType))
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		logger.Error("http-failed-to-write-response-body", zap.Error(err))
	}
}

func makeJSONResponse(w http.ResponseWriter, statusCode int, body interface{}, logger *zap.Logger) {
	d, err := json.Marshal(body)
	if err != nil {
		logger.Error("http-failed-to-marshal-body", zap.Error(err))
	}
	c := new(bytes.Buffer)
	err = json.Compact(c, d)
	if err != nil {
		logger.Error("http-failed-to-compact-json", zap.Error(err))
	}
	makeResponse(w, c.Bytes(), statusCode, "application/json", logger)
}

func makeErrorResponse(w http.ResponseWriter, er *errorResponse, logger *zap.Logger) {
	makeJSONResponse(w, er.StatusCode, er, logger)
}

func readBodyJSON(r *http.Request, into interface{}, logger *zap.Logger) *errorResponse {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read-request-body-error", zap.Error(err))
		er := errorResponses["internal_server_error"]
		return &er

	}

	if err := json.Unmarshal(requestBody, into); err != nil {
		logger.Error("request-json-unmarshal-error", zap.Error(err))
		er := errorResponses["malformed_request_body"]
		return &er
	}

	if val, ok := into.(validation.Validatable); ok {
		if err := val.Validate(); err != nil {
			if e, ok := err.(validation.InternalError); ok {
				logger.Error("read-request-validate-error", zap.Error(e))
				er := errorResponses["internal_server_error"]
				return &er
			}
			er := errorResponses["request_validation_failed"]
			er.Details = err.Error()
			return &er
		}
	}

	return nil
}

func makeListOptions(r *http.Request, filters []string) (postmand.RepositoryListOptions, error) {
	listOptions := postmand.RepositoryListOptions{}

	if err := r.ParseForm(); err != nil {
		return listOptions, err
	}

	// Parse limit and offset
	limit := 50
	offset := 0
	if r.Form.Get("limit") != "" {
		v, err := strconv.Atoi(r.Form.Get("limit"))
		if err == nil && v <= limit {
			limit = v
		}
	}
	if r.Form.Get("offset") != "" {
		v, err := strconv.Atoi(r.Form.Get("offset"))
		if err == nil {
			offset = v
		}
	}
	listOptions.Limit = limit
	listOptions.Offset = offset

	// Parse filters
	f := make(map[string]interface{})
	for _, filter := range filters {
		if r.Form.Get(filter) != "" {
			f[filter] = r.Form.Get(filter)
		}
	}
	listOptions.Filters = f

	return listOptions, nil
}
