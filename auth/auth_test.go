package auth

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/WadeCappa/authmaster/authmaster"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type createAccountResult struct {
	err   error
	token string
}

const (
	TESTING_ADDRESS = "localhost:50051"
)

var (
	addr = flag.String("addr", TESTING_ADDRESS, "the address to connect to")
)

func TestGettingUserId(t *testing.T) {
	flag.Parse()

	username := fmt.Sprintf("%d", rand.Int())
	password := fmt.Sprintf("%d", rand.Int())

	token, err := runWithConnection(*addr, func(conn *grpc.ClientConn) createAccountResult {
		newMetadata := metadata.Pairs()
		newContext := metadata.NewOutgoingContext(context.Background(), newMetadata)

		c := authmaster.NewAuthmasterClient(conn)
		_, err := c.CreateUser(newContext, &authmaster.CreateUserRequest{Username: username, Password: password})
		if err != nil {
			return createAccountResult{err: err}
		}

		login, err := c.Login(newContext, &authmaster.LoginRequest{Username: username, Password: password})
		if err != nil {
			return createAccountResult{err: err}
		}

		return createAccountResult{err: nil, token: login.Token}
	})

	if err != nil || token.err != nil {
		t.Errorf("Failed to login: %v, %v", err, token.err)
		return
	}

	userId, err := getUser(token.token, *addr)
	if err != nil {
		t.Errorf("failed to test user %v", err)
		return
	}

	if userId == nil {
		t.Errorf("Did not find user id")
		return
	}

	log.Printf("user logged in %d", *userId)
}
