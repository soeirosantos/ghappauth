package auth

import (
	"testing"
	"time"

	"ghappauth/internal/types"
)

// Fake private key for testing
const testPrivateKey = `-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCz2Zi1a56fXTol
hfKUxw01Iv76IyztPWwZ66LL0VQoTCiWg0hQsTyu9Hbb1/pB7p5Cfgp0+J7Jjqlw
tPNM48V1qGCD6nY/WuvtLiNafLbxSrP23cDVcqBrkrEYNi1iGe8aDBi8sZe6tjjQ
nKPRokBOgroeywf4CIVk+m0xKG0wV2nzBlrYQPPu5nOBEXJZRgPz8VWQONHib79T
jz2u9u5RCK9a7mkGLpD3SDB4cR34NXMM3GFyymEMul8wrokT7oLy+tOWZ0akCzj4
6jjNDObTkg5Ds3V60K6CzHwuX97Zid6qTgTsXgIAYTzA7XtBOJXsLtf+QQzrTxfQ
e6WGZJ6pAgMBAAECggEAGkwYLCQkTsu5k4BWdrKe7RoaNkZWbLStIyJ1QgrFnMQm
BdFZuD1yJhLQyQHM3JR2HOWDxJQQg7WRyPfBwN29IvW6cboiuCN57nWZ7cKL7WmC
5EyLGaq9EihM2hbHOVpEFTGXnF+gqUhjs6fJEaoBGnm679hivehhDjbKfdmaEus2
ECL1ClotN0MQg7NfKN332nc60AJ+HTg/+1TntEX7lFM8f1EsfoIAzitmAsInx+sF
7FB7cM8pySDkOmDQTiEaCQUlv9T933Z8V7eSkIMi5AEB6u66JcjXiFnYzAh+IMKO
UCGiQkd1HBmJF01am4axpmPivnJqxyyQDGTUlxkiEQKBgQDAgm2+H8W0abIGI4mc
gU4v0HN9Ld3WaHFa09xjNIQNxTjPp3gAf9yppejzGz6M7JBQS00y6EpP0PHPbuhE
Y5+b25mWL1Xat8N83utvH2T8wRl7ZDUr81ne/gpOU1tIwEkklXv6/ivd1HNcIxa4
6dM6olx6srKv8rk7J/AW4l0vgwKBgQDvKlOVxoA/Fw+H72WilIuPhe8L8kvJgnGX
O5HcDLrniYdCMWfGCQykK3eQMaFGyz2cYUko4OpuH1Kovvfi55+4RAiLK5S03D3f
jkLCxA0OG1fSxtO+/WSS3JiLYPMlZkfEzWdvl4XB6IBZoWgk4MhZMYYJT/pCNrA9
zcexXzyVYwKBgGIJShf7mDRjazzDFk50bzvcXSQPmpyY/bkykVaYJPPaTy846tze
QKLIkhRT+IvN4URyxLK7JzT0hGCN640AawT1VYbtPjyvPse1wpIJm+U39WEoTAfA
2zC7kMYIn0EyY01VLxlIHVDP45u1ZtnughqnGo+Ft4fxBTHCCfutdaU/AoGALmiL
MZwEFLn31Ivar/KdJit6GFpa5G5AdnUjt4xs1DL2oRyPI3lsD4sztzI6Nk+H1Al4
tcr3EolXc9Eirs/9STdCZSb+wx2dj/y97ac3VU5u+0KDoiLvWiQeIaWdaNtw/7pP
4PKJDPh9t2a/m7BWkCAw/yuaxzBvgH6myj9NtTsCgYB7Kq71lb0M2AUDia3uWNbH
p2H7JUEy1syS/iUjyeM6KvQEvtF5xSpci8abRYMOJ7maJRHfE/vmxEVHydjYZn2N
T4SUfsuwcvVzX7klfnpJlHQP2+jYjVHv/BjKvaQkEHFCh7F3+6qB8PuBGYRlSewg
+yk+EW4NuLecU0irv/dreg==
-----END PRIVATE KEY-----`

func TestNewGitHubAppAuth(t *testing.T) {
	tests := []struct {
		name    string
		config  *types.GitHubAppConfig
		wantErr bool
	}{
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "missing app_id",
			config: &types.GitHubAppConfig{
				PrivateKey: testPrivateKey,
			},
			wantErr: true,
		},
		{
			name: "empty app_id",
			config: &types.GitHubAppConfig{
				AppID:          "",
				PrivateKey:     testPrivateKey,
				InstallationID: "67890",
			},
			wantErr: true,
		},
		{
			name: "missing private_key",
			config: &types.GitHubAppConfig{
				AppID: "12345",
			},
			wantErr: true,
		},
		{
			name: "empty private_key",
			config: &types.GitHubAppConfig{
				AppID:          "12345",
				PrivateKey:     "",
				InstallationID: "67890",
			},
			wantErr: true,
		},
		{
			name: "missing installation_id",
			config: &types.GitHubAppConfig{
				AppID:      "12345",
				PrivateKey: testPrivateKey,
			},
			wantErr: true,
		},
		{
			name: "empty installation_id",
			config: &types.GitHubAppConfig{
				AppID:          "12345",
				PrivateKey:     testPrivateKey,
				InstallationID: "",
			},
			wantErr: true,
		},
		{
			name: "invalid private key",
			config: &types.GitHubAppConfig{
				AppID:      "12345",
				PrivateKey: "invalid-key",
			},
			wantErr: true,
		},
		{
			name: "malformed PEM private key",
			config: &types.GitHubAppConfig{
				AppID:          "12345",
				PrivateKey:     "-----BEGIN RSA PRIVATE KEY-----\ninvalid-content\n-----END RSA PRIVATE KEY-----",
				InstallationID: "67890",
			},
			wantErr: true,
		},
		{
			name: "invalid app_id format",
			config: &types.GitHubAppConfig{
				AppID:          "invalid-app-id",
				PrivateKey:     testPrivateKey,
				InstallationID: "67890",
			},
			wantErr: true,
		},
		{
			name: "invalid installation_id format",
			config: &types.GitHubAppConfig{
				AppID:          "12345",
				PrivateKey:     testPrivateKey,
				InstallationID: "invalid-installation-id",
			},
			wantErr: true,
		},
		{
			name: "valid config",
			config: &types.GitHubAppConfig{
				AppID:          "12345",
				PrivateKey:     testPrivateKey,
				InstallationID: "67890",
			},
			wantErr: false,
		},
		{
			name: "valid config with base URL",
			config: &types.GitHubAppConfig{
				AppID:          "12345",
				PrivateKey:     testPrivateKey,
				InstallationID: "67890",
				BaseURL:        "https://api.github.com",
			},
			wantErr: false,
		},
		{
			name: "valid config with empty base URL (should use default)",
			config: &types.GitHubAppConfig{
				AppID:          "12345",
				PrivateKey:     testPrivateKey,
				InstallationID: "67890",
				BaseURL:        "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth, err := NewGitHubAppAuth(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGitHubAppAuth() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			// For valid configs, test that BaseURL is set correctly
			if !tt.wantErr && auth != nil {
				if tt.config.BaseURL == "" {
					// Should use default
					if auth.baseURL != "https://api.github.com" {
						t.Errorf("Expected default BaseURL 'https://api.github.com', got '%s'", auth.baseURL)
					}
				} else {
					// Should use provided value
					if auth.baseURL != tt.config.BaseURL {
						t.Errorf("Expected BaseURL '%s', got '%s'", tt.config.BaseURL, auth.baseURL)
					}
				}
			}
		})
	}
}

func TestGitHubAppAuth_GenerateJWT(t *testing.T) {
	config := &types.GitHubAppConfig{
		AppID:          "12345",
		PrivateKey:     testPrivateKey,
		InstallationID: "67890",
	}

	auth, err := NewGitHubAppAuth(config)
	if err != nil {
		t.Fatalf("Failed to create auth: %v", err)
	}

	jwt, err := auth.GenerateJWT()
	if err != nil {
		t.Fatalf("GenerateJWT() error = %v", err)
	}

	if jwt == "" {
		t.Error("GenerateJWT() returned empty token")
	}

	// Test with invalid app ID
	config.AppID = "invalid"
	_, err = NewGitHubAppAuth(config)
	if err == nil {
		t.Error("NewGitHubAppAuth() should fail with invalid app ID")
	}
}

func TestTokenManager_IsTokenExpired(t *testing.T) {
	config := &types.GitHubAppConfig{
		AppID:          "12345",
		PrivateKey:     testPrivateKey,
		InstallationID: "67890",
	}

	auth, err := NewGitHubAppAuth(config)
	if err != nil {
		t.Fatalf("Failed to create auth: %v", err)
	}

	tm := NewTokenManager(auth, 5*time.Minute)

	tests := []struct {
		name     string
		token    *types.GitHubAppToken
		buffer   time.Duration
		expected bool
	}{
		{
			name:     "nil token",
			token:    nil,
			buffer:   0,
			expected: true,
		},
		{
			name: "expired token",
			token: &types.GitHubAppToken{
				ExpiresAt: time.Now().Add(-1 * time.Hour),
			},
			buffer:   0,
			expected: true,
		},
		{
			name: "valid token",
			token: &types.GitHubAppToken{
				ExpiresAt: time.Now().Add(1 * time.Hour),
			},
			buffer:   0,
			expected: false,
		},
		{
			name: "token expiring within buffer",
			token: &types.GitHubAppToken{
				ExpiresAt: time.Now().Add(2 * time.Minute),
			},
			buffer:   5 * time.Minute,
			expected: true,
		},
		{
			name: "token not expiring within buffer",
			token: &types.GitHubAppToken{
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			buffer:   5 * time.Minute,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tm.IsTokenExpired(tt.token, tt.buffer)
			if result != tt.expected {
				t.Errorf("IsTokenExpired() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewTokenManager(t *testing.T) {
	config := &types.GitHubAppConfig{
		AppID:          "12345",
		PrivateKey:     testPrivateKey,
		InstallationID: "67890",
	}

	auth, err := NewGitHubAppAuth(config)
	if err != nil {
		t.Fatalf("Failed to create auth: %v", err)
	}

	// Test with default buffer
	tm := NewTokenManager(auth, 0)
	if tm == nil {
		t.Error("NewTokenManager() returned nil")
	}
	if tm.GetRenewBuffer() != 5*time.Minute {
		t.Errorf("Expected default buffer of 5 minutes, got %v", tm.GetRenewBuffer())
	}

	// Test with custom buffer
	customBuffer := 10 * time.Minute
	tm = NewTokenManager(auth, customBuffer)
	if tm.GetRenewBuffer() != customBuffer {
		t.Errorf("Expected buffer of %v, got %v", customBuffer, tm.GetRenewBuffer())
	}
}

func TestTokenManager_GetCacheStats(t *testing.T) {
	config := &types.GitHubAppConfig{
		AppID:          "12345",
		PrivateKey:     testPrivateKey,
		InstallationID: "67890",
	}

	auth, err := NewGitHubAppAuth(config)
	if err != nil {
		t.Fatalf("Failed to create auth: %v", err)
	}

	tm := NewTokenManager(auth, 5*time.Minute)
	stats := tm.GetCacheStats()

	if stats["total_cached"] != 0 {
		t.Errorf("Expected 0 cached tokens, got %v", stats["total_cached"])
	}

	if stats["renew_buffer"] != "5m0s" {
		t.Errorf("Expected renew buffer of 5m0s, got %v", stats["renew_buffer"])
	}
}

func TestTokenManager_SetRenewBuffer(t *testing.T) {
	config := &types.GitHubAppConfig{
		AppID:          "12345",
		PrivateKey:     testPrivateKey,
		InstallationID: "67890",
	}

	auth, err := NewGitHubAppAuth(config)
	if err != nil {
		t.Fatalf("Failed to create auth: %v", err)
	}

	tm := NewTokenManager(auth, 5*time.Minute)
	
	newBuffer := 15 * time.Minute
	tm.SetRenewBuffer(newBuffer)
	
	if tm.GetRenewBuffer() != newBuffer {
		t.Errorf("Expected buffer of %v, got %v", newBuffer, tm.GetRenewBuffer())
	}
}

func TestTokenManager_ClearCache(t *testing.T) {
	config := &types.GitHubAppConfig{
		AppID:          "12345",
		PrivateKey:     testPrivateKey,
		InstallationID: "67890",
	}

	auth, err := NewGitHubAppAuth(config)
	if err != nil {
		t.Fatalf("Failed to create auth: %v", err)
	}

	tm := NewTokenManager(auth, 5*time.Minute)
	
	// Clear empty cache
	tm.ClearCache()
	stats := tm.GetCacheStats()
	if stats["total_cached"] != 0 {
		t.Errorf("Expected 0 cached tokens after clear, got %v", stats["total_cached"])
	}
} 