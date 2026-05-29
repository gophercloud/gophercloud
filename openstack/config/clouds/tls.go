package clouds

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path"
	"strings"
)

func computeTLSConfig(cloud Cloud, options cloudOpts) (*tls.Config, error) {
	tlsConfig := new(tls.Config)
	if caCertPath := coalesce(options.caCertPath, os.Getenv("OS_CACERT"), cloud.CACertFile); caCertPath != "" {
		caCertPath, err := resolveTilde(caCertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve user home directory: %w", err)
		}

		caCert, err := os.ReadFile(caCertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open the CA cert file: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(bytes.TrimSpace(caCert)); !ok {
			return nil, fmt.Errorf("failed to parse the CA Cert from %q", caCertPath)
		}
		tlsConfig.RootCAs = caCertPool
	}

	tlsConfig.InsecureSkipVerify = func() bool {
		if options.insecure != nil {
			return *options.insecure
		}
		if cloud.Verify != nil {
			return !*cloud.Verify
		}
		return false
	}()

	if clientCertPath, clientKeyPath := coalesce(options.clientCertPath, os.Getenv("OS_CERT"), cloud.ClientCertFile), coalesce(options.clientKeyPath, os.Getenv("OS_KEY"), cloud.ClientKeyFile); clientCertPath != "" && clientKeyPath != "" {
		clientCertPath, err := resolveTilde(clientCertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve user home directory in client cert path: %w", err)
		}
		clientKeyPath, err := resolveTilde(clientKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve user home directory in client cert key path: %w", err)
		}

		clientCert, err := os.ReadFile(clientCertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read the client cert file: %w", err)
		}

		clientKey, err := os.ReadFile(clientKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read the client cert key file: %w", err)
		}

		cert, err := tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			return nil, err
		}

		tlsConfig.Certificates = []tls.Certificate{cert}
	} else if clientCertPath != "" && clientKeyPath == "" {
		return nil, fmt.Errorf("client cert is set, but client cert key is missing")
	} else if clientCertPath == "" && clientKeyPath != "" {
		return nil, fmt.Errorf("client cert key is set, but client cert is missing")
	}

	return tlsConfig, nil
}

func resolveTilde(p string) (string, error) {
	if after := strings.TrimPrefix(p, "~/"); after != p {
		h, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve user home directory: %w", err)
		}
		return path.Join(h, after), nil
	}
	return p, nil
}
