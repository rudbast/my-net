package core

import (
	"database/sql"
	"log"
)

type (
	Service struct {
		logger   *log.Logger
		database *sql.DB
	}
)

func New(lg *log.Logger, db *sql.DB) *Service {
	return &Service{
		logger:   lg,
		database: db,
	}
}
