package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func NewDatabase(connectionUrl string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), connectionUrl)

	if err != nil {
		return nil, err
	}

	return conn, nil
}