package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"os"
)

func main() {
	cert := ""
	key := ""
	ca := ""
	endpoint := ""
	c, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		panic(err)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(ca))
	cfg := &tls.Config{}
	cfg.Certificates = append(cfg.Certificates, c)
	cfg.RootCAs = pool
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: cfg,
		},
	}
	resp, err := client.Get(endpoint)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}
