package types

import "time"

// GitHubAppConfig holds the configuration for a GitHub App
type GitHubAppConfig struct {
	AppID          string `json:"app_id"`
	PrivateKey     string `json:"private_key"`
	InstallationID string `json:"installation_id,omitempty"`
	BaseURL        string `json:"base_url,omitempty"`
}

// GitHubAppToken represents an installation access token
type GitHubAppToken struct {
	Token               string            `json:"token"`
	ExpiresAt           time.Time         `json:"expires_at"`
	Permissions         map[string]string `json:"permissions"`
	RepositorySelection string            `json:"repository_selection"`
	Repositories        []Repository      `json:"repositories,omitempty"`
}

// Repository represents a GitHub repository
type Repository struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
}

// GitHubAppInstallation represents a GitHub App installation
type GitHubAppInstallation struct {
	ID                     int                    `json:"id"`
	Account                Account                `json:"account"`
	RepositorySelection    string                 `json:"repository_selection"`
	Permissions            map[string]string      `json:"permissions"`
	SuspendedAt            *time.Time             `json:"suspended_at"`
	SuspendedBy            interface{}            `json:"suspended_by"`
	CreatedAt              time.Time              `json:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at"`
	SingleFileName         *string                `json:"single_file_name"`
	HasMultipleSingleFiles bool                   `json:"has_multiple_single_files"`
	SingleFilePaths        []string               `json:"single_file_paths"`
	AppID                  int                    `json:"app_id"`
	AppSlug                string                 `json:"app_slug"`
	TargetID               int                    `json:"target_id"`
	TargetType             string                 `json:"target_type"`
	Events                 []string               `json:"events"`
}

// Account represents a GitHub account (user or organization)
type Account struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	Type  string `json:"type"`
}

// GitHubApp represents a GitHub App
type GitHubApp struct {
	ID                int               `json:"id"`
	Slug              string            `json:"slug"`
	NodeID            string            `json:"node_id"`
	Owner             Account           `json:"owner"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	ExternalURL       string            `json:"external_url"`
	HTMLURL           string            `json:"html_url"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	Permissions       map[string]string `json:"permissions"`
	Events            []string          `json:"events"`
	InstallationsCount int              `json:"installations_count"`
	ClientID          string            `json:"client_id"`
	ClientSecret      string            `json:"client_secret"`
	PEM               string            `json:"pem"`
}

// GitHubAPIError represents an error response from the GitHub API
type GitHubAPIError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url,omitempty"`
	Errors           []struct {
		Resource string `json:"resource"`
		Field    string `json:"field"`
		Code     string `json:"code"`
	} `json:"errors,omitempty"`
}

// InstallationTokenRequest represents the request body for creating an installation token
type InstallationTokenRequest struct {
	RepositoryIDs []int  `json:"repository_ids,omitempty"`
	Permissions   map[string]string `json:"permissions,omitempty"`
}

// InstallationTokenResponse represents the response from creating an installation token
type InstallationTokenResponse struct {
	Token               string            `json:"token"`
	ExpiresAt           time.Time         `json:"expires_at"`
	Permissions         map[string]string `json:"permissions"`
	RepositorySelection string            `json:"repository_selection"`
	Repositories        []Repository      `json:"repositories,omitempty"`
} 