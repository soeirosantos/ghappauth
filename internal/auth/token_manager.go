package auth

import (
	"fmt"
	"sync"
	"time"

	"ghappauth/internal/types"
)

// TokenManager handles caching and automatic renewal of installation tokens
type TokenManager struct {
	auth        *GitHubAppAuth
	cache       map[string]*cachedToken
	mutex       sync.RWMutex
	renewBuffer time.Duration // How much time before expiry to renew the token
}

// cachedToken represents a cached installation token
type cachedToken struct {
	token      *types.GitHubAppToken
	createdAt  time.Time
	lastUsed   time.Time
	renewing   bool
	renewMutex sync.Mutex
}

// NewTokenManager creates a new token manager
func NewTokenManager(auth *GitHubAppAuth, renewBuffer time.Duration) *TokenManager {
	if renewBuffer == 0 {
		renewBuffer = 5 * time.Minute // Default 5 minutes buffer
	}

	return &TokenManager{
		auth:        auth,
		cache:       make(map[string]*cachedToken),
		renewBuffer: renewBuffer,
	}
}

// GetToken retrieves a valid installation token, renewing if necessary
func (tm *TokenManager) GetToken() (*types.GitHubAppToken, error) {
	tm.mutex.RLock()
	cached, exists := tm.cache[tm.auth.config.InstallationID]
	tm.mutex.RUnlock()

	if exists && cached != nil {
		cached.lastUsed = time.Now()

		if !tm.IsTokenExpired(cached.token, tm.renewBuffer) {
			return cached.token, nil
		}

		return tm.renewToken(cached)
	}

	return tm.createNewToken()
}

// renewToken renews an existing cached token
func (tm *TokenManager) renewToken(cached *cachedToken) (*types.GitHubAppToken, error) {
	cached.renewMutex.Lock()
	defer cached.renewMutex.Unlock()

	if !tm.IsTokenExpired(cached.token, tm.renewBuffer) {
		return cached.token, nil
	}

	cached.renewing = true

	newToken, err := tm.auth.GetInstallationToken()
	if err != nil {
		cached.renewing = false
		return nil, fmt.Errorf("failed to renew token: %w", err)
	}

	cached.token = newToken
	cached.createdAt = time.Now()
	cached.renewing = false

	return newToken, nil
}

// createNewToken creates a new token and caches it
func (tm *TokenManager) createNewToken() (*types.GitHubAppToken, error) {
	token, err := tm.auth.GetInstallationToken()
	if err != nil {
		return nil, fmt.Errorf("failed to create new token: %w", err)
	}

	tm.mutex.Lock()
	tm.cache[tm.auth.config.InstallationID] = &cachedToken{
		token:     token,
		createdAt: time.Now(),
		lastUsed:  time.Now(),
		renewing:  false,
	}
	tm.mutex.Unlock()

	return token, nil
}

// InvalidateToken removes the token from cache, forcing renewal on next request
func (tm *TokenManager) InvalidateToken() {
	tm.mutex.Lock()
	delete(tm.cache, tm.auth.config.InstallationID)
	tm.mutex.Unlock()
}

// ClearCache removes all cached tokens
func (tm *TokenManager) ClearCache() {
	tm.mutex.Lock()
	tm.cache = make(map[string]*cachedToken)
	tm.mutex.Unlock()
}

// GetCacheStats returns statistics about the token cache
func (tm *TokenManager) GetCacheStats() map[string]interface{} {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	stats := map[string]interface{}{
		"total_cached": len(tm.cache),
		"renew_buffer": tm.renewBuffer.String(),
	}

	cacheDetails := make(map[string]interface{})
	for installationID, cached := range tm.cache {
		if cached != nil {
			cacheDetails[installationID] = map[string]interface{}{
				"created_at": cached.createdAt,
				"last_used":  cached.lastUsed,
				"renewing":   cached.renewing,
				"expires_at": cached.token.ExpiresAt,
				"is_expired": tm.IsTokenExpired(cached.token, 0),
			}
		}
	}
	stats["cache_details"] = cacheDetails

	return stats
}

// SetRenewBuffer sets the renewal buffer duration
func (tm *TokenManager) SetRenewBuffer(buffer time.Duration) {
	tm.renewBuffer = buffer
}

// GetRenewBuffer returns the current renewal buffer duration
func (tm *TokenManager) GetRenewBuffer() time.Duration {
	return tm.renewBuffer
}

// IsTokenExpired checks if a token is expired or will expire soon
func (tm *TokenManager) IsTokenExpired(token *types.GitHubAppToken, buffer time.Duration) bool {
	if token == nil {
		return true
	}
	
	return time.Now().Add(buffer).After(token.ExpiresAt)
} 