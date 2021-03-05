package main

import (
	"log"
	"os"
	"time"

	"github.com/allisson/go-env"
	"github.com/allisson/postmand/http"
	"github.com/allisson/postmand/http/handler"
	"github.com/allisson/postmand/repository"
	"github.com/allisson/postmand/service"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	// Setup logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("logger-setup-error: %v\n", err)
	}
	// nolint:errcheck
	defer logger.Sync()

	// Setup postgresql database
	db, err := sqlx.Open("postgres", env.GetString("POSTMAND_DATABASE_URL", ""))
	if err != nil {
		logger.Fatal("database-setup-error", zap.Error(err))
	}
	if err := db.Ping(); err != nil {
		logger.Fatal("database-ping-error", zap.Error(err))
	}
	db.SetMaxOpenConns(env.GetInt("POSTMAND_DATABASE_MAX_OPEN_CONNS", 2))

	// Setup cli
	app := cli.NewApp()
	app.Name = "postmand"
	app.Usage = "CLI"
	app.Authors = []*cli.Author{
		{
			Name:  "Allisson Azevedo",
			Email: "allisson@gmail.com",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:    "migrate",
			Aliases: []string{"m"},
			Usage:   "executes database migration",
			Action: func(c *cli.Context) error {
				migrationRepository := repository.NewMigration(
					db,
					env.GetString("POSTMAND_DATABASE_MIGRATION_DIR", "file:///db/migrations"),
				)
				migrationService := service.NewMigration(migrationRepository)
				return migrationService.Run(c.Context)
			},
		},
		{
			Name:    "worker",
			Aliases: []string{"w"},
			Usage:   "executes worker to dispatch webhooks",
			Action: func(c *cli.Context) error {
				deliveryRepository := repository.NewDelivery(db)
				pollingInterval := time.Duration(env.GetInt("POSTMAND_POLLING_INTERVAL", 1000)) * time.Millisecond
				workerService := service.NewWorker(deliveryRepository, logger, pollingInterval)
				workerService.Run(c.Context)
				return nil
			},
		},
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "executes http server",
			Action: func(c *cli.Context) error {
				webhookRepository := repository.NewWebhook(db)
				deliveryRepository := repository.NewDelivery(db)
				webhookService := service.NewWebhook(webhookRepository)
				deliveryService := service.NewDelivery(deliveryRepository)
				webhookHandler := handler.NewWebhook(webhookService, logger)
				deliveryHandler := handler.NewDelivery(deliveryService, logger)

				mux := http.NewRouter(logger)
				mux.Route("/v1/webhooks", func(r chi.Router) {
					r.Get("/", webhookHandler.List)
					r.Post("/", webhookHandler.Create)
					r.Get("/{webhook_id}", webhookHandler.Get)
					r.Put("/{webhook_id}", webhookHandler.Update)
					r.Delete("/{webhook_id}", webhookHandler.Delete)
				})
				mux.Route("/v1/deliveries", func(r chi.Router) {
					r.Get("/", deliveryHandler.List)
					r.Post("/", deliveryHandler.Create)
					r.Get("/{delivery_id}", deliveryHandler.Get)
					r.Delete("/{delivery_id}", deliveryHandler.Delete)
				})

				server := http.NewServer(mux, env.GetInt("POSTMAND_HTTP_PORT", 8000), logger)
				server.Run()

				return nil
			},
		},
	}
	err = app.Run(os.Args)
	if err != nil {
		logger.Fatal("cli-error", zap.Error(err))
	}
}
