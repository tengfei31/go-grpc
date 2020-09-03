package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "go-grpc/proto"
)

//rpc接收结构
type SearchService struct{}

func (s *SearchService) Search(c context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Printf("recv from client: %s", r.GetRequest())
	return &pb.SearchResponse{Response: r.GetRequest() + " Server"}, nil
}

//端口号
const PORT = "9001"

func main() {
	server := grpc.NewServer()
	pb.RegisterSearchServiceServer(server, &SearchService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	server.Serve(lis)
}
