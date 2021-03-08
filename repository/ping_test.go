package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	ctx := context.Background()
	th := newTestHelper()
	defer th.db.Close()

	err := th.pingRepository.Run(ctx)
	assert.Nil(t, err)
}
