package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/WadeCappa/real_time_chat/channel-manager/external_channel_manager"
	"google.golang.org/grpc"
)

const (
	DEFAULT_PORT = 50055
)

var (
	port = flag.Int("port", DEFAULT_PORT, "port for this service")
)

func getPostgresUrl() string {
	postgresHostname := os.Getenv("CHANNEL_MANAGER_POSTGRES_HOSTNAME")
	return fmt.Sprintf("postgres://postgres:pass@%s/channel_manager_db", postgresHostname)
}

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
