package repository

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func init() {
	txdb.Register("pgx", "postgres", os.Getenv("POSTMAND_DATABASE_URL"))
	rand.Seed(time.Now().UnixNano())
}

type testHelper struct {
	db                *sqlx.DB
	webhookRepository *Webhook
}

func newTestHelper() testHelper {
	cName := fmt.Sprintf("connection_%d", time.Now().UnixNano())
	db, _ := sqlx.Open("pgx", cName)
	return testHelper{
		db:                db,
		webhookRepository: NewWebhook(db),
	}
}
