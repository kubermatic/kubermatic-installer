package installer

import (
	"fmt"
)

type DexOrganization struct {
	Name string `yaml:"name"`
}

type DexConnectorConfig struct {
	Issuer        string            `yaml:"issuer,omitempty"`
	ClientID      string            `yaml:"clientID"`
	ClientSecret  string            `yaml:"clientSecret"`
	RedirectURI   string            `yaml:"redirectURI"`
	Organizations []DexOrganization `yaml:"orgs,omitempty"`
}

type DexConnector struct {
	Type   string             `yaml:"type"`
	ID     string             `yaml:"id"`
	Name   string             `yaml:"name"`
	Config DexConnectorConfig `yaml:"config"`
}

func NewGoogleDexConnector(clientID string, secretKey string, baseURL string) DexConnector {
	return DexConnector{
		Type: "oidc",
		ID:   "google",
		Name: "google",
		Config: DexConnectorConfig{
			Issuer:       "https://accounts.google.com",
			ClientID:     clientID,
			ClientSecret: secretKey,
			RedirectURI:  fmt.Sprintf("%s/dex/callback", baseURL),
		},
	}
}

func NewGitHubDexConnector(clientID string, secretKey string, baseURL string, org string) DexConnector {
	connector := DexConnector{
		Type: "github",
		ID:   "github",
		Name: "GitHub",
		Config: DexConnectorConfig{
			ClientID:     clientID,
			ClientSecret: secretKey,
			RedirectURI:  fmt.Sprintf("%s/dex/callback", baseURL),
		},
	}

	if len(org) > 0 {
		connector.Config.Organizations = []DexOrganization{
			{
				Name: org,
			},
		}
	}

	return connector
}

type DexClient struct {
	ID           string   `yaml:"id"`
	Name         string   `yaml:"name"`
	Secret       string   `yaml:"secret"`
	RedirectURIs []string `yaml:"RedirectURIs"`
}
