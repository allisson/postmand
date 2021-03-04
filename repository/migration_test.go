package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigration(t *testing.T) {
	ctx := context.Background()
	th := newTestHelper()
	defer th.db.Close()
	migrationDir := "file://../db/migrations"
	migration := NewMigration(th.db, migrationDir)
	err := migration.Run(ctx)
	assert.Nil(t, err)
}
