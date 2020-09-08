package gtls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
)

type Client struct {
	ServerName string
	CaFile     string
	CertFile   string
	KeyFile    string
}

func (t *Client) GetCredentialsByCA() (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(t.CertFile, t.KeyFile)
	if err != nil {
		return nil, err
	}
	ca, err := ioutil.ReadFile(t.CaFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(ca)
	if ok == false {
		return nil, errors.New("certPool.AppendCertsFromPEM error")
	}
	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   t.ServerName,
		RootCAs:      certPool,
	})
	return c, err
}

func (t *Client) GetTLSCredentials() (credentials.TransportCredentials, error) {
	c, err := credentials.NewClientTLSFromFile(t.CertFile, t.ServerName)
	if err != nil {
		return nil, err
	}
	return c, err
}
