package main

import (
	"log"
	"os"

	"github.com/allisson/go-env"
	"github.com/allisson/postmand/repository"
	"github.com/allisson/postmand/service"
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
			Usage:   "execute database migration",
			Action: func(c *cli.Context) error {
				migrationRepository := repository.NewMigration(
					db,
					env.GetString("POSTMAND_DATABASE_MIGRATION_DIR", "file:///db/migrations"),
				)
				migrationService := service.NewMigration(migrationRepository)
				return migrationService.Run(c.Context)
			},
		},
	}
	err = app.Run(os.Args)
	if err != nil {
		logger.Fatal("cli-error", zap.Error(err))
	}
}
