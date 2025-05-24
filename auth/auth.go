package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/WadeCappa/authmaster/authmaster"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type testResult struct {
	err    error
	userId *int64
}

const (
	CURRENT_USER_KEY = "current-user"
	AUTH_HEADER_KEY  = "Authorization"
)

func getUser(token, authUrl string) (*int64, error) {
	testResult, err := runWithConnection(authUrl, func(conn *grpc.ClientConn) testResult {
		userId, err := test(conn, token)
		return testResult{userId: userId, err: err}
	})

	if err != nil {
		return nil, err
	}

	if testResult.err != nil {
		return nil, testResult.err
	}

	return testResult.userId, nil
}

func test(conn *grpc.ClientConn, token string) (*int64, error) {
	newMetadata := metadata.Pairs(AUTH_HEADER_KEY, token)
	newContext := metadata.NewOutgoingContext(context.Background(), newMetadata)

	c := authmaster.NewAuthmasterClient(conn)
	response, err := c.TestAuth(newContext, &authmaster.TestAuthRequest{})

	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, fmt.Errorf("did not receive either an error or a response from the auth server")
	}

	return &response.UserId, nil
}

func runWithConnection[T any](addr string, f func(conn *grpc.ClientConn) T) (*T, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	log.Printf("attempting to connect to %s which resolved to %s", addr, conn.CanonicalTarget())
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
		return nil, err
	}
	defer conn.Close()
	res := f(conn)
	return &res, nil
}

func AuthenticateUser(ctx context.Context, authServiceUrl string) (*int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	fmt.Println(md)
	if !ok {
		return nil, fmt.Errorf("could not find auth header")
	}

	t, ok := md["authorization"]
	if !ok {
		return nil, fmt.Errorf("no auth header provided")
	}

	if len(t) != 1 {
		return nil, fmt.Errorf("too many auth headers")
	}

	userId, err := getUser(t[0], authServiceUrl)
	if err != nil {
		return nil, fmt.Errorf("Encountered some error reaching out to auth service: %v\n", err)
	}

	return userId, nil
}
