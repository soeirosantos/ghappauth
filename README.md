# GHAppAuth

## ⚠️ Disclaimer

This is an **educational project** and is not officially supported. I don't recommend adding this as a direct dependency to your production applications. However, the code itself is production-ready and follows Go best practices - you're welcome to copy and adapt it for your own use.

## What This Is

This is a Go library for handling [GitHub App](https://docs.github.com/en/apps) authentication. GitHub Apps are a way to integrate with GitHub's API that's more secure and flexible than personal access tokens.

The main use case is when you want your application to act on behalf of itself (not on behalf of users) and need to interact with GitHub repositories, issues, pull requests, etc. Think CI/CD systems, automation tools, or any service that needs to access GitHub data. Check the GitHub documentation [About authentication with a GitHub App](https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/about-authentication-with-a-github-app) for more details.

## Features

- **JWT Generation**: Generate JWT tokens for GitHub App authentication
- **Installation Token Management**: Retrieve and manage installation access tokens
- **Automatic Token Renewal**: Built-in caching with automatic token renewal before expiration
- **Thread-Safe**: Concurrent access support with proper locking mechanisms
- **Error Handling**: Comprehensive error handling with detailed error messages

## Quick Start

### Setting Up a GitHub App

Before you can use this library, you need to create a GitHub App:

1. Go to [GitHub Developer Settings](https://github.com/settings/apps)
2. Click "New GitHub App"
3. Fill in the basic info (name, description, homepage URL)
4. Set the **Repository permissions** you need (e.g., `Contents: Read` for repository access)
5. Choose **Repository access**: "All repositories" or "Only select repositories"
6. Click "Create GitHub App"
7. Generate a private key (this downloads a `.pem` file - keep it secure!)
8. Note your **App ID** (you'll see it in the app settings)

Once your app is created, you need to install it on the repositories/organizations you want to access:

1. Go to your app's "Install App" page
2. Click "Install" on the organization/user you want to access
3. Note the **Installation ID** (you can find this in the URL or by calling the [List installations API](https://docs.github.com/en/rest/apps/apps#list-installations-for-the-authenticated-app))

### Running the Example

1. **Set up your environment variables:**

```bash
export GITHUB_APP_ID="your_app_id"
export GITHUB_PRIVATE_KEY="$(cat /path/to/your/private-key.pem)"
export GITHUB_INSTALLATION_ID="your_installation_id"
```

2. **Run the example:**

```bash
go run cmd/example/main.go
```

The example will:
- Get your app information
- Get installation details
- Generate an installation token
- Demonstrate token caching
- Show cache statistics
- Generate a JWT token
- Check token expiration

## Basic Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "time"
    
    "ghappauth/internal/auth"
    "ghappauth/internal/types"
)

func main() {
    config := &types.GitHubAppConfig{
        AppID:          "your_app_id",
        PrivateKey:     "your_private_key",
        InstallationID: "your_installation_id",
    }

    githubAuth, err := auth.NewGitHubAppAuth(config)
    if err != nil {
        log.Fatal(err)
    }

    tokenManager := auth.NewTokenManager(githubAuth, 5*time.Minute)
    token, err := tokenManager.GetToken()
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Token: %s", token.Token)
    log.Printf("Expires at: %s", token.ExpiresAt)

    // Use the token to list repositories, assuming the app has repo:read permissions
    client := &http.Client{}
    req, err := http.NewRequest("GET", "https://api.github.com/installation/repositories", nil)
    if err != nil {
        log.Fatalf("Failed to create request: %v", err)
    }
    
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))
    req.Header.Set("Accept", "application/vnd.github.v3+json")

    resp, err := client.Do(req)
    if err != nil {
        log.Fatalf("Failed to make request: %v", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Failed to read response: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        log.Fatalf("API request failed with status %d: %s", resp.StatusCode, string(body))
    }

    var reposResponse struct {
        Repositories []struct {
            Name     string `json:"name"`
            FullName string `json:"full_name"`
            Private  bool   `json:"private"`
        } `json:"repositories"`
    }

    if err := json.Unmarshal(body, &reposResponse); err != nil {
        log.Fatalf("Failed to parse response: %v", err)
    }

    log.Printf("Found %d repositories:", len(reposResponse.Repositories))
    for _, repo := range reposResponse.Repositories {
        visibility := "public"
        if repo.Private {
            visibility = "private"
        }
        log.Printf("  - %s (%s)", repo.FullName, visibility)
    }
}
```
