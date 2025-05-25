package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/WadeCappa/real_time_chat/channel-manager/external_channel_manager"
	"google.golang.org/grpc"
)

const (
	DEFAULT_POSTGRES_URL  = "postgres://postgres:pass@localhost:5432/channel_manager_db"
	DEFAULT_AUTH_HOSTNAME = "localhost:50051"
	DEFAULT_PORT          = 50055
)

var (
	authHostname = flag.String("auth-hostname", DEFAULT_AUTH_HOSTNAME, "the hostname for the auth service")
	postgresUrl  = flag.String("postgres-url", DEFAULT_POSTGRES_URL, "the hostname for postgres")
	port         = flag.Int("port", DEFAULT_PORT, "port for this service")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	external_channel_manager.RegisterExternalchannelmanagerServer(s, &externalChatManangerServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
