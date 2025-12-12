package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	connMaxLifetime   = time.Hour
	connMaxIdleTime   = time.Minute * 30
	maxIdleConns      = 10
	maxOpenConns      = 10
	healthCheckPeriod = time.Minute
)

func Config(dbUrl string) *pgxpool.Config {
	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Panic("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = int32(maxOpenConns)
	dbConfig.MaxConnLifetime = connMaxLifetime
	dbConfig.MaxConnIdleTime = connMaxIdleTime
	dbConfig.HealthCheckPeriod = healthCheckPeriod

	return dbConfig
}

func NewPostgresDB(dbUrl string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connPool, err := pgxpool.NewWithConfig(ctx, Config(dbUrl))
	if err != nil {
		log.Panic("error while creating connection to the database!!", err)
	}

	connection, err := connPool.Acquire(ctx)
	if err != nil {
		log.Panic("error while acquiring connection from the database pool!!", err)
	}
	defer connection.Release()

	err = connection.Ping(ctx)
	if err != nil {
		log.Panic("could not ping database", err)
	}

	return connPool
}
