package gtls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
)

type Server struct {
	CaFile   string
	CertFile string
	KeyFile  string
}

func (t *Server) GetCredentialsByCA() (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(t.CertFile, t.KeyFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(t.CertFile)
	if err != nil {
		return nil, err
	}
	ok := certPool.AppendCertsFromPEM(ca)
	if ok == false {
		return nil, errors.New("certPool.AppendCertsFromPEM err")
	}
	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{
			cert,
		},
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  certPool,
	})
	return c, err
}

func (t *Server) GetTLSCredentials() (credentials.TransportCredentials, error) {
	c, err := credentials.NewServerTLSFromFile(t.CertFile, t.KeyFile)
	if err != nil {
		return nil, err
	}
	return c, err
}
