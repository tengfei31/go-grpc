package main

import (
	"context"
	"go-grpc/pkg/gtls"
	pb "go-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

const PORT string = "9004"

type AUTH struct {
	AppKey    string
	AppSecret string
}

func (a *AUTH) SetAUTH(appKey, appSecret string) {
	a.AppKey = appKey
	a.AppSecret = appSecret
}
func (a *AUTH) GetAppKey() string {
	return a.AppKey
}
func (a *AUTH) GetAppSecret() string {
	return a.AppSecret
}

func (a *AUTH) Check(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok == false {
		return status.Errorf(codes.Unauthenticated, "自定义认证token失败")
	}
	var (
		appKey    string
		appSecret string
	)
	if value, ok := md["app_key"]; ok {
		appKey = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}
	if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
		return status.Errorf(codes.Unauthenticated, "自定义认证token无效")
	}
	return nil
}

func (a *AUTH) RequireTransportSecurity() bool {
	return true
}

type SearchService struct {
	auth *AUTH
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	if err := s.auth.Check(ctx); err != nil {
		return nil, err
	}
	return &pb.SearchResponse{Response: r.GetRequest() + "server"}, nil
}

func main() {
	tlsServer := gtls.Server{
		//CaFile: "/Users/admin/Documents/www/go/src/github.com/tengfei31/go-grpc/conf/ca.pem",
		CertFile: "/Users/admin/Documents/www/go/src/github.com/tengfei31/go-grpc/conf/server/server.pem",
		KeyFile:  "/Users/admin/Documents/www/go/src/github.com/tengfei31/go-grpc/conf/server/server.key",
	}
	c, err := tlsServer.GetTLSCredentials()
	if err != nil {
		log.Fatalf("tlsServer.GetTLSCredentials error: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.Creds(c),
		//grpc_middleware.WithUnaryServerChain(RecoveryInterceptor, LoggingInterceptor),
	}
	server := grpc.NewServer(opts...)
	pb.RegisterSearchServiceServer(server, &SearchService{})
	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	server.Serve(lis)

}
