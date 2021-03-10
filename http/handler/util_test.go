package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestBind(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	q := req.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "5")
	q.Add("active", "true")
	q.Add("success", "true")
	q.Add("webhook_id", "7fe789b4-dec6-4eab-8144-c50e95b866ee")
	q.Add("delivery_id", "f5eedad9-76b8-4262-b2b6-bdc3771a80b2")
	q.Add("status", "pending")
	q.Add("created_at.gt", "2021-03-08T20:50:08.353038Z")
	q.Add("created_at.gte", "2021-03-08T20:50:08.353038Z")
	q.Add("created_at.lt", "2021-03-08T20:50:08.353038Z")
	q.Add("created_at.lte", "2021-03-08T20:50:08.353038Z")
	req.URL.RawQuery = q.Encode()
	rf := requestFilters{}
	err := requestBind(req, &rf)
	assert.Nil(t, err)
}
