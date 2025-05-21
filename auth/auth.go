package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/WadeCappa/authmaster/authmaster"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type testResult struct {
	err    error
	userId *int64
}

const (
	CURRENT_USER_KEY              = "current-user"
	AUTH_HEADER_KEY               = "Authorization"
	AUTH_SERVICE_URL_ENV_VARIABLE = "AUTH_URL"
)

func Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		authenticateUser(c)
	}
}

func authenticateUser(c *gin.Context) {
	token := c.Request.Header[AUTH_HEADER_KEY]
	if token == nil {
		fmt.Println("No auth header provided")
		c.Abort()
		return
	}

	if len(token) != 1 {
		fmt.Printf("too many auth tokens provided, this is unexpected (provided %d)\n", len(token))
		c.Abort()
		return
	}

	authUrl := os.Getenv(AUTH_SERVICE_URL_ENV_VARIABLE)
	testResult, err := runWithConnection(authUrl, func(conn *grpc.ClientConn) testResult {
		userId, err := test(conn, token[0])
		return testResult{userId: userId, err: err}
	})

	if err != nil {
		fmt.Printf("Encountered some error reaching out to auth service: %v\n", err)
		c.Abort()
		return
	}

	if testResult.err != nil {
		fmt.Printf("Encountered some error: %v\n", testResult.err)
		c.Abort()
		return
	}

	c.Set(CURRENT_USER_KEY, *testResult.userId)
	fmt.Printf("found user of id %d\n", *testResult.userId)
	c.Next()
}

func test(conn *grpc.ClientConn, token string) (*int64, error) {
	newMetadata := metadata.Pairs("Authorization", token)
	newContext := metadata.NewOutgoingContext(context.Background(), newMetadata)

	c := authmaster.NewAuthmasterClient(conn)
	response, err := c.TestAuth(newContext, &authmaster.TestAuthRequest{})

	if err != nil {
		return nil, err
	}

	return &response.UserId, nil
}

func runWithConnection[T any](addr string, f func(conn *grpc.ClientConn) T) (*T, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("did not connect: %v\n", err)
		return nil, err
	}
	defer conn.Close()
	res := f(conn)
	return &res, nil
}
