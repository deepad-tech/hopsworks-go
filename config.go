package hopsworks

import "net/http"

type ClientConfig struct {
	apiKey string

	Host    string
	Port    int
	Project string

	HTTPClient *http.Client
}

const (
	hopsworksHost = "c.app.hopsworks.ai"
	hopsworksPort = 443
)

// DefaultConfig returns default configuration for Hopsworks API client.
func DefaultConfig(apiKey string) *ClientConfig {
	return &ClientConfig{
		apiKey:     apiKey,
		Host:       hopsworksHost,
		Port:       hopsworksPort,
		Project:    "",
		HTTPClient: &http.Client{},
	}
}
