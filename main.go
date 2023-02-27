package main

import (
	"database/sql"
	"github.com/dbssensei/articlewebservice/api"
	"github.com/dbssensei/articlewebservice/db/seeder"
	"github.com/go-redis/redis/v8"
	"os"

	db "github.com/dbssensei/articlewebservice/db/sqlc"
	"github.com/dbssensei/articlewebservice/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	redis := redis.NewClient(&redis.Options{
		Addr:     config.RedisSource,
		Password: config.RedisPassword,
		DB:       0,
	})

	store := db.NewStore(conn)
	runDBMigration(config.MigrationURL, config.DBSource)
	runDBSeeder(store)
	runGinServer(config, store, redis)
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}

func runDBSeeder(store db.Store) {
	seeder.Execute(store)
}

func runGinServer(config util.Config, store db.Store, redis *redis.Client) {
	server, err := api.NewServer(config, store, redis)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}
