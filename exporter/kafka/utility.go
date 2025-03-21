package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

func createTLSConfiguration(tlsSkipVerify bool, caFile, certFile, keyFile string) (t *tls.Config, err error) {
	t = &tls.Config{
		InsecureSkipVerify: tlsSkipVerify,
	}
	if certFile != "" && keyFile != "" && caFile != "" {
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return t, err
		}

		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return t, err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: tlsSkipVerify,
		}
	}
	return
}
