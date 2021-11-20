package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.uber.org/zap"

	"github.com/allisson/postmand"
)

type requestFilters struct {
	Limit        int         `json:"limit,string"`
	Offset       int         `json:"offset,string"`
	Active       bool        `json:"active,string"`
	Success      bool        `json:"success,string"`
	WebhookID    postmand.ID `json:"webhook_id"`
	DeliveryID   postmand.ID `json:"delivery_id"`
	Status       string      `json:"status"`
	CreatedAtGt  time.Time   `json:"created_at.gt"`
	CreatedAtGte time.Time   `json:"created_at.gte"`
	CreatedAtLt  time.Time   `json:"created_at.lt"`
	CreatedAtLte time.Time   `json:"created_at.lte"`
}

func makeResponse(w http.ResponseWriter, body []byte, statusCode int, contentType string, logger *zap.Logger) {
	w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", contentType))
	w.WriteHeader(statusCode)
	if _, err := w.Write(body); err != nil {
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

func requestBind(r *http.Request, into interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	m := make(map[string]string)
	for key := range r.Form {
		m[key] = r.Form.Get(key)
	}
	rawJSON, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(rawJSON, into)
}

func makeListOptions(r *http.Request, filters []string) postmand.RepositoryListOptions {
	listOptions := postmand.RepositoryListOptions{
		Limit:  50,
		Offset: 0,
	}
	rf := requestFilters{}

	if err := requestBind(r, &rf); err != nil {
		return listOptions
	}

	if rf.Limit != 0 {
		listOptions.Limit = rf.Limit
	}
	listOptions.Offset = rf.Offset

	// Parse filters
	f := make(map[string]interface{})
	for _, filter := range filters {
		filterValue := r.Form.Get(filter)
		if filterValue != "" {
			f[filter] = filterValue
		}
	}
	listOptions.Filters = f

	return listOptions
}
