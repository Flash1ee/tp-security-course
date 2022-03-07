package database

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type PostgresConn struct {
	Conn *pgx.Conn
}

func NewPostgresConn(databaseUrl string) *PostgresConn {
	db, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		logrus.Fatalf("can not connect to database - url: %s", databaseUrl)
	}
	return &PostgresConn{
		Conn: db,
	}
}
