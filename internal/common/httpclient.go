package common

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// MakeHTTPClient creates a ready to use http client. It configures authentication based on app configuration
func MakeHTTPClient(request *http.Request) (*http.Client, error) {
	certFile := viper.GetString("AUTH.client_cert")
	keyFile := viper.GetString("AUTH.client_key")
	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Error("Failed to load client key and certificate")
			return nil, err
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		caFile := viper.GetString("AUTH.root_ca")
		if caFile != "" {
			caCerts := x509.NewCertPool()

			caData, err := ioutil.ReadFile(caFile)
			if err != nil {
				log.Error("Failed to load root CA")
				return nil, err
			}
			caCerts.AppendCertsFromPEM(caData)
			tlsConfig.RootCAs = caCerts
		}

		tlsConfig.BuildNameToCertificate()
		transport := &http.Transport{TLSClientConfig: tlsConfig, Proxy: http.ProxyFromEnvironment}
		return &http.Client{Transport: transport}, nil
	}

	//dev only
	if rhIdentity := viper.GetString("AUTH.x_rh_identity"); rhIdentity != "" {
		request.Header.Set("x-rh-identity", rhIdentity)
	}
	user := viper.GetString("AUTH.user")
	password := viper.GetString("AUTH.password")
	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}
	return &http.Client{}, nil
}
