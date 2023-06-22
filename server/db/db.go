package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type Database struct {
	Client *sqlx.DB
}

func NewDatabase(conStr string) (*Database, error) {
	dbConn, err := sqlx.Connect("postgres", conStr)
	if err != nil {
		log.Debug().Msg(conStr)
		return &Database{}, fmt.Errorf("failed to connect to database: %s", err)
	}
	
	return &Database{
		Client: dbConn,
	}, nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.Client.PingContext(ctx)
}