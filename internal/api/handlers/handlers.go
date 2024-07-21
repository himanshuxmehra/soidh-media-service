package handlers

import (
	"database/sql"

	"go.uber.org/zap"
)

type Handlers struct {
	db     *sql.DB
	logger *zap.Logger
}

func New(db *sql.DB, logger *zap.Logger) *Handlers {
	return &Handlers{
		db:     db,
		logger: logger,
	}
}