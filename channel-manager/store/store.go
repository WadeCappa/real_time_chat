package store

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func Call[T any](postgresUrl string, query func(*pgx.Conn) T) (*T, error) {
	// this is probably insecure. Will want to change how we access this in the future
	conn, err := pgx.Connect(context.Background(), postgresUrl)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())
	res := query(conn)
	return &res, nil
}
