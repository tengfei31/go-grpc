package main

import (
	"context"
	"go-grpc/pkg/gtls"
	pb "go-grpc/proto"
	"google.golang.org/grpc"
	"log"
)

const PORT string = "9001"

func main() {
	tlsClient := gtls.Client{
		CertFile:   "/Users/admin/Documents/www/go/src/github.com/tengfei31/go-grpc/conf/client/client.pem",
		KeyFile:    "/Users/admin/Documents/www/go/src/github.com/tengfei31/go-grpc/conf/client/client.key",
		ServerName: "go-grpc",
	}
	c, err := tlsClient.GetTLSCredentials()
	if err != nil {
		log.Fatalf("tlsClient.GetTLSCredentials error: %v", err)
	}
	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c))
	//conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(context.Background(), &pb.SearchRequest{
		Request: "gRPC",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}

	log.Printf("resp: %s", resp.GetResponse())
}
