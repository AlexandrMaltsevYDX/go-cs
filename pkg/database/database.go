package database

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/AlexandrMaltsevYDX/go-cs/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDbPool(config *config.DatabaseConfig) *pgxpool.Pool {
	pgxpool, err := pgxpool.New(context.Background(), config.URL)

	if err != nil {
		log.Error().Msg("Failed to create database pool")
		panic(err)
	}

	log.Info().Msg("Database pool created successfully")
	return pgxpool
}
