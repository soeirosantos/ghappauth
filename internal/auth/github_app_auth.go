package auth

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"ghappauth/internal/types"
)

// GitHubAppAuth handles GitHub App authentication
type GitHubAppAuth struct {
	config     *types.GitHubAppConfig
	privateKey *rsa.PrivateKey
	baseURL    string
	httpClient *HTTPClient
}

// NewGitHubAppAuth creates a new GitHub App authentication instance
func NewGitHubAppAuth(config *types.GitHubAppConfig) (*GitHubAppAuth, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if config.AppID == "" {
		return nil, fmt.Errorf("app_id is required")
	}

	if config.PrivateKey == "" {
		return nil, fmt.Errorf("private_key is required")
	}

	if config.InstallationID == "" {
		return nil, fmt.Errorf("installation_id is required")
	}

	_, err := strconv.Atoi(config.AppID)
	if err != nil {
		return nil, fmt.Errorf("invalid app_id: %w", err)
	}

	_, err = strconv.Atoi(config.InstallationID)
	if err != nil {
		return nil, fmt.Errorf("invalid installation_id: %w", err)
	}

	privateKey, err := parsePrivateKey(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}

	return &GitHubAppAuth{
		config:     config,
		privateKey: privateKey,
		baseURL:    baseURL,
		httpClient: NewHTTPClient(nil),
	}, nil
}

// parsePrivateKey parses a PEM-encoded RSA private key
func parsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	var privateKey *rsa.PrivateKey
	var err error

	switch block.Type {
	case "RSA PRIVATE KEY":
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKCS8 private key: %w", err)
		}
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("key is not an RSA private key")
		}
	default:
		return nil, fmt.Errorf("unsupported private key type: %s", block.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
	}

	return privateKey, nil
}

// GenerateJWT generates a JWT token for GitHub App authentication
func (g *GitHubAppAuth) GenerateJWT() (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    g.config.AppID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute)),
		Subject:   g.config.AppID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(g.privateKey)
}

// GetInstallationToken retrieves an installation access token from GitHub
func (g *GitHubAppAuth) GetInstallationToken() (*types.GitHubAppToken, error) {
	jwt, err := g.GenerateJWT()
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	url := fmt.Sprintf("%s/app/installations/%s/access_tokens", g.baseURL, g.config.InstallationID)
	
	var tokenResponse types.InstallationTokenResponse
	err = g.httpClient.DoRequest(context.Background(), &RequestConfig{
		Method:        "POST",
		URL:           url,
		AuthToken:     jwt,
		ExpectedStatus: http.StatusCreated,
	}, &tokenResponse)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get installation token: %w", err)
	}

	return &types.GitHubAppToken{
		Token:               tokenResponse.Token,
		ExpiresAt:           tokenResponse.ExpiresAt,
		Permissions:         tokenResponse.Permissions,
		RepositorySelection: tokenResponse.RepositorySelection,
		Repositories:        tokenResponse.Repositories,
	}, nil
}

// GetAppInfo retrieves information about the GitHub App
func (g *GitHubAppAuth) GetAppInfo() (*types.GitHubApp, error) {
	jwt, err := g.GenerateJWT()
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	url := fmt.Sprintf("%s/app", g.baseURL)
	
	var app types.GitHubApp
	err = g.httpClient.DoRequest(context.Background(), &RequestConfig{
		Method:        "GET",
		URL:           url,
		AuthToken:     jwt,
		ExpectedStatus: http.StatusOK,
	}, &app)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get app info: %w", err)
	}

	return &app, nil
}

// GetInstallation retrieves information about the configured installation
func (g *GitHubAppAuth) GetInstallation() (*types.GitHubAppInstallation, error) {
	jwt, err := g.GenerateJWT()
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	url := fmt.Sprintf("%s/app/installations/%s", g.baseURL, g.config.InstallationID)
	
	var installation types.GitHubAppInstallation
	err = g.httpClient.DoRequest(context.Background(), &RequestConfig{
		Method:        "GET",
		URL:           url,
		AuthToken:     jwt,
		ExpectedStatus: http.StatusOK,
	}, &installation)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get installation: %w", err)
	}

	return &installation, nil
}

 