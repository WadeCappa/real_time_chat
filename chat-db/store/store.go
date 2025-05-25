package store

import (
	"fmt"

	"github.com/gocql/gocql"
)

func Call[T any](dbUrl string, query func(*gocql.Session) (*T, error)) (*T, error) {
	cluster := gocql.NewCluster(dbUrl)

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: "cassandra",
	}

	cluster.Keyspace = "posts_db"
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4

	// this is probably insecure. Will want to change how we access this in the future
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}
	if session == nil {
		return nil, fmt.Errorf("failed to create session, but had no error")
	}
	defer session.Close()

	return query(session)
}
