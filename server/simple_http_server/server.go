package main

import (
	"context"
	"go-grpc/pkg/gtls"
	pb "go-grpc/proto"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strings"
	//grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	//"crypto/tls"
	//"crypto/x509"
	//"io/ioutil"
	//"google.golang.org/grpc/credentials"
)

const PORT string = "9001"

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + " HTTP Server"}, nil
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

	mux := GetHttpServeMux()
	server := grpc.NewServer(opts...)
	pb.RegisterSearchServiceServer(server, &SearchService{})
	http.ListenAndServe(
		":"+PORT,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				server.ServeHTTP(w, r)
			} else {
				mux.ServeHTTP(w, r)
			}
			return
		}),
	)
}

func GetHttpServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("eddycjy: go-grpc"))
	})
	return mux
}
