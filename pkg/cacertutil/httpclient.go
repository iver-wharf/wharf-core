package cacertutil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/iver-wharf/wharf-core/pkg/logger"
)

var log = logger.NewScoped("CA-CERT-UTIL")

// NewHTTPClientWithCerts creates a fresh net/http.Client populated with some
// root CA certificates from file.
// Argument must point to an existing file with PEM formatted certificates.
//
// Based on https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/
func NewHTTPClientWithCerts(localCertFile string) (*http.Client, error) {
	rootCAs := getCertPoolFromEnvironment()
	if err := appendCertsFromFile(rootCAs, localCertFile); err != nil {
		return nil, err
	}

	// Trust the augmented cert pool in our client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: rootCAs,
			},
		},
	}

	return client, nil
}

func getCertPoolFromEnvironment() (rootCAs *x509.CertPool) {
	rootCAs, _ = x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
		log.Debug().Message("Using empty cert pool.")
	} else {
		log.Debug().Message("Using system's cert pool.")
	}
	return
}

func appendCertsFromFile(certPool *x509.CertPool, certFile string) error {
	certs, err := ioutil.ReadFile(certFile)
	if err != nil {
		return fmt.Errorf("failed to append %q to certPool: %v", certFile, err)
	}

	log.Debug().WithString("file", certFile).Message("Loaded certs.")

	if ok := certPool.AppendCertsFromPEM(certs); !ok {
		log.Debug().
			WithString("file", certFile).
			Message("No certs appended from file.")
	}
	return nil
}
