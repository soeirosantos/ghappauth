package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"ghappauth/internal/auth"
	"ghappauth/internal/types"
)


func main() {
	// Load configuration from environment variables
	config := &types.GitHubAppConfig{
		AppID:          os.Getenv("GITHUB_APP_ID"),
		PrivateKey:     os.Getenv("GITHUB_PRIVATE_KEY"),
		InstallationID: os.Getenv("GITHUB_INSTALLATION_ID"),
	}

	// Validate required environment variables
	if config.AppID == "" {
		log.Fatal("GITHUB_APP_ID environment variable is required")
	}
	if config.PrivateKey == "" {
		log.Fatal("GITHUB_PRIVATE_KEY environment variable is required")
	}
	if config.InstallationID == "" {
		log.Fatal("GITHUB_INSTALLATION_ID environment variable is required")
	}

	// Create GitHub App authentication instance
	githubAuth, err := auth.NewGitHubAppAuth(config)
	if err != nil {
		log.Fatalf("Failed to create GitHub App auth: %v", err)
	}

	// Create token manager with 5-minute renewal buffer
	tokenManager := auth.NewTokenManager(githubAuth, 5*time.Minute)

	fmt.Println("=== GitHub App Authentication Example ===")

	// Example 1: Get app information
	fmt.Println("1. Getting app information...")
	appInfo, err := githubAuth.GetAppInfo()
	if err != nil {
		log.Printf("Failed to get app info: %v", err)
	} else {
		fmt.Printf("App Name: %s\n", appInfo.Name)
		fmt.Printf("App Description: %s\n", appInfo.Description)
		fmt.Printf("App URL: %s\n", appInfo.HTMLURL)
		fmt.Printf("Installations Count: %d\n", appInfo.InstallationsCount)
	}

	// Example 2: Get installation information
	fmt.Println("2. Getting installation information...")
	installation, err := githubAuth.GetInstallation()
	if err != nil {
		log.Printf("Failed to get installation info: %v", err)
	} else {
		fmt.Printf("Installation ID: %d\n", installation.ID)
		fmt.Printf("Account: %s (%s)\n", installation.Account.Login, installation.Account.Type)
		fmt.Printf("Repository Selection: %s\n", installation.RepositorySelection)
		fmt.Printf("Permissions: %+v\n", installation.Permissions)
		fmt.Println()
	}

	// Example 3: Get installation token (with caching)
	fmt.Println("3. Getting installation token...")
	token, err := tokenManager.GetToken()
	if err != nil {
		log.Printf("Failed to get installation token: %v", err)
	} else {
		fmt.Printf("Token: %s...\n", token.Token[:20])
		fmt.Printf("Expires At: %s\n", token.ExpiresAt.Format(time.RFC3339))
		fmt.Printf("Repository Selection: %s\n", token.RepositorySelection)
		fmt.Printf("Permissions: %+v\n", token.Permissions)
		if len(token.Repositories) > 0 {
			fmt.Printf("Repositories: %d accessible\n", len(token.Repositories))
		}
		fmt.Println()
	}

	// Example 4: Demonstrate token caching
	fmt.Println("4. Demonstrating token caching...")
	start := time.Now()
	token2, err := tokenManager.GetToken()
	duration := time.Since(start)
	if err != nil {
		log.Printf("Failed to get cached token: %v", err)
	} else {
		fmt.Printf("Cached token retrieved in: %v\n", duration)
		fmt.Printf("Token matches: %t\n", token.Token == token2.Token)
		fmt.Println()
	}

	// Example 5: Show cache statistics
	fmt.Println("5. Cache statistics...")
	stats := tokenManager.GetCacheStats()
	statsJSON, _ := json.MarshalIndent(stats, "", "  ")
	fmt.Printf("Cache Stats:\n%s\n", string(statsJSON))
	fmt.Println()

	// Example 6: Generate JWT token
	fmt.Println("6. Generating JWT token...")
	jwt, err := githubAuth.GenerateJWT()
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
	} else {
		fmt.Printf("JWT Token: %s...\n", jwt[:50])
		fmt.Println()
	}

	// Example 7: Check token expiration
	fmt.Println("7. Checking token expiration...")
	isExpired := tokenManager.IsTokenExpired(token, 0)
	fmt.Printf("Token is expired: %t\n", isExpired)

	willExpireSoon := tokenManager.IsTokenExpired(token, 10*time.Minute)
	fmt.Printf("Token will expire within 10 minutes: %t\n", willExpireSoon)
	fmt.Println()

	fmt.Println("=== Example completed successfully ===")
}
