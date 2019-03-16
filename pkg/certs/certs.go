package certs

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"
)

var (
	certDir = path.Join(os.Getenv("GOPATH"), "src/github.com/ninnemana/vinyl/certs")
)

func ServerCertificate() (*tls.Certificate, *x509.CertPool, error) {
	certDir := path.Join(os.Getenv("GOPATH"), "src/github.com/ninnemana/vinyl/certs")
	certFile, err := os.Open(path.Join(certDir, "vinyltap.alexninneman.com.crt"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to open certificate")
	}

	cert, err := ioutil.ReadAll(certFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to read certificate")
	}

	keyFile, err := os.Open(path.Join(certDir, "vinyltap.alexninneman.com.key"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to open certificate key")
	}

	key, err := ioutil.ReadAll(keyFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to read certificate key")
	}

	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to load certificate")
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		return nil, nil, errors.Wrap(err, "failed to parse certificate")
	}

	return &pair, certPool, nil
}

func ServerCredentials() (credentials.TransportCredentials, error) {
	creds, err := credentials.NewServerTLSFromFile(
		path.Join(certDir, "vinyltap.alexninneman.com.crt"),
		path.Join(certDir, "vinyltap.alexninneman.com.key"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transport credentials")
	}

	return creds, nil
}

func ClientCredentials() (credentials.TransportCredentials, error) {
	certDir := path.Join(os.Getenv("GOPATH"), "src/github.com/ninnemana/vinyl/certs")

	creds, err := credentials.NewClientTLSFromFile(
		path.Join(certDir, "vinyltap.alexninneman.com.crt"),
		"",
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client certificate")
	}

	return creds, nil
}
