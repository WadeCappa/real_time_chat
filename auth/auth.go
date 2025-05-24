package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

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
	CURRENT_USER_KEY = "current-user"
	AUTH_HEADER_KEY  = "Authorization"
)

func Build(authServiceUrl string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authenticateUser(c, authServiceUrl)
	}
}

func getUser(token, authUrl string) (*int64, error) {
	testResult, err := runWithConnection(authUrl, func(conn *grpc.ClientConn) testResult {
		userId, err := test(conn, token)
		return testResult{userId: userId, err: err}
	})

	if err != nil {
		return nil, err
	}

	if testResult.err != nil {
		return nil, err
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

func authenticateUser(c *gin.Context, authServiceUrl string) {
	token := c.Request.Header[AUTH_HEADER_KEY]
	if token == nil {
		fmt.Println("No auth header provided")
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	if len(token) != 1 {
		fmt.Printf("too many auth tokens provided, this is unexpected (provided %d)\n", len(token))
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	userId, err := getUser(token[0], authServiceUrl)
	if err != nil {
		fmt.Printf("Encountered some error: %v\n", err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
	}

	c.Set(CURRENT_USER_KEY, userId)
	fmt.Printf("found user of id %d\n", userId)
	c.Next()
}
