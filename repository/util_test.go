package repository

import (
	"fmt"
	"os"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func init() {
	txdb.Register("pgx", "postgres", os.Getenv("POSTMAND_TEST_DATABASE_URL"))
	//NewRand(NewSource(seed))
	//rand.Seed(time.Now().UnixNano())
}

type testHelper struct {
	db                        *sqlx.DB
	webhookRepository         *Webhook
	deliveryRepository        *Delivery
	deliveryAttemptRepository *DeliveryAttempt
	pingRepository            *Ping
}

func newTestHelper() testHelper {
	cName := fmt.Sprintf("connection_%d", time.Now().UnixNano())
	db, _ := sqlx.Open("pgx", cName)
	return testHelper{
		db:                        db,
		webhookRepository:         NewWebhook(db),
		deliveryRepository:        NewDelivery(db),
		deliveryAttemptRepository: NewDeliveryAttempt(db),
		pingRepository:            NewPing(db),
	}
}
