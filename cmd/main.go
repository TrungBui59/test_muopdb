package main

import (
	"fmt"
	"github.com/TrungBui59/test_muopdb/interval/configs"
	"github.com/TrungBui59/test_muopdb/interval/muopdbclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	cfg, err := configs.NewConfig("")
	if err != nil {
		panic(err)
	}

	conn, err := createGRPCClientConn(fmt.Sprintf("%s:%d", cfg.MuopDBConfig.Host, cfg.MuopDBConfig.Port))
	if err != nil {
		log.Fatalf("failed to create grpc client conn: %v", err)
	}

	muopDBclient := muopdbclient.NewClient(conn)
	defer muopDBclient.Close()

}

func createGRPCClientConn(serverAddress string) (*grpc.ClientConn, error) {
	// Create a connection to the server
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(serverAddress, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
