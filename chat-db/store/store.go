package store

import (
	"context"
	"log"

	"github.com/WadeCappa/real_time_chat/chat-db/result"
	"github.com/jackc/pgx/v5"
)

func Call[T any](postgresUrl string, query func(*pgx.Conn) result.Result[T]) result.Result[T] {
	log.Printf("attempting to connect to %s", postgresUrl)
	// this is probably insecure. Will want to change how we access this in the future
	conn, err := pgx.Connect(context.Background(), postgresUrl)
	if err != nil {
		return result.Failed[T](err)
	}
	defer conn.Close(context.Background())
	return query(conn)
}
