// backend/internal/auth/token_blacklist.go

package auth

import (
	"sync"
	"time"
)

type TokenBlacklist struct {
	tokens          map[string]time.Time
	mutex           sync.RWMutex
	cleanupInterval time.Duration
}

func NewTokenBlacklist(cleanupInterval time.Duration) *TokenBlacklist {
	bl := &TokenBlacklist{
		tokens:          make(map[string]time.Time),
		mutex:           sync.RWMutex{},
		cleanupInterval: cleanupInterval,
	}

	go bl.periodicCleanup()

	return bl
}

func (bl *TokenBlacklist) Add(token string, expiry time.Time) {
	bl.mutex.Lock()
	defer bl.mutex.Unlock()

	bl.tokens[token] = expiry
}

func (bl *TokenBlacklist) IsBlacklisted(token string) bool {
	bl.mutex.RLock()
	defer bl.mutex.RUnlock()

	expiry, exists := bl.tokens[token]
	if !exists {
		return false
	}

	// Se o token expirou, podemos removê-lo da blacklist
	if time.Now().After(expiry) {
		// Não remover durante leitura para evitar condição de corrida
		return false
	}

	return true
}

func (bl *TokenBlacklist) cleanup() {
	bl.mutex.Lock()
	defer bl.mutex.Unlock()

	now := time.Now()
	for token, expiry := range bl.tokens {
		if now.After(expiry) {
			delete(bl.tokens, token)
		}
	}
}

func (bl *TokenBlacklist) periodicCleanup() {
	ticker := time.NewTicker(bl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		bl.cleanup()
	}
}
