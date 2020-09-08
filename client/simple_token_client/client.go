package main

import (
	"context"
	"go-grpc/pkg/gtls"

	//"go-grpc/pkg/gtls"
	pb "go-grpc/proto"
	"google.golang.org/grpc"
	"log"
)

const PORT string = "9004"

type AUTH struct {
	AppKey    string
	AppSecret string
}

func (a *AUTH) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"app_key":    a.AppKey,
		"app_secret": a.AppSecret,
	}, nil
}

func (a *AUTH) RequireTransportSecurity() bool {
	return true
}

func main() {
	tlsClient := gtls.Client{
		ServerName: "go-grpc",
		CertFile:   "/Users/admin/Documents/www/go/src/github.com/tengfei31/go-grpc/conf/client/client.pem",
		KeyFile:    "/Users/admin/Documents/www/go/src/github.com/tengfei31/go-grpc/conf/client/client.key",
	}
	c, err := tlsClient.GetTLSCredentials()
	if err != nil {
		log.Fatalf("tlsClient.GetTLSCredentials error: %v", err)
	}
	auth := AUTH{
		AppKey:    "go-grpc",
		AppSecret: "wtf",
	}
	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c), grpc.WithPerRPCCredentials(&auth))
	if err != nil {
		log.Fatalf("grpc.Dial error: %v", err)
	}
	defer conn.Close()
	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(context.Background(), &pb.SearchRequest{Request: "grpc"})
	if err != nil {
		log.Fatalf("client.Search error: %v", err)
	}
	log.Printf("resp: %s", resp.GetResponse())

}
