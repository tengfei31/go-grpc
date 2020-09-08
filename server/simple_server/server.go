package main

import (
	"context"
	"time"

	//"crypto/tls"
	//"crypto/x509"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	//"io/ioutil"
	"log"
	"net"
	"runtime/debug"

	pb "go-grpc/proto"
)

//rpc接收结构
type SearchService struct{}

func (s *SearchService) Search(c context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	for i := 0; i <= 5; i++ {
		if c.Err() == context.Canceled {
			return nil, status.Errorf(codes.Canceled, "SearchService.Search canceled")
		}
		time.Sleep(time.Second * 1)
	}

	log.Printf("recv from client: %s", r.GetRequest())
	return &pb.SearchResponse{Response: r.GetRequest() + " Server"}, nil
}

//端口号
const PORT = "9001"

func main() {
	//cert, err := tls.LoadX509KeyPair("/Users/admin/Documents/www/go/src/github.com/tengfei31/go-grpc/conf/server/server.pem", "/Users/admin/Documents/www/go/src/github.com/tengfei31/go-grpc/conf/server/server.key")
	//if err != nil {
	//	log.Fatalf("tls.LoadX509KeyPair error: %v", err)
	//}
	//certPool := x509.NewCertPool()
	//ca, err := ioutil.ReadFile("/Users/admin/Documents/www/go/src/github.com/tengfei31/go-grpc/conf/ca.pem")
	//if err != nil {
	//	log.Fatalf("ioutil.ReadFile err: %v", err)
	//}
	//ok := certPool.AppendCertsFromPEM(ca)
	//if ok == false {
	//	log.Fatalf("certPool.AppendCertsFromPEM err")
	//}
	//c := credentials.NewTLS(&tls.Config{
	//	Certificates: []tls.Certificate{cert},
	//	ClientAuth: tls.RequireAndVerifyClientCert,
	//	ClientCAs: certPool,
	//})

	opts := []grpc.ServerOption{
		//grpc.Creds(c),
		grpc_middleware.WithUnaryServerChain(RecoveryInterceptor, LoggingInterceptor),
	}
	server := grpc.NewServer(opts...)
	pb.RegisterSearchServiceServer(server, &SearchService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	server.Serve(lis)
}

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	return resp, err
}

func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()
	return handler(ctx, req)
}
