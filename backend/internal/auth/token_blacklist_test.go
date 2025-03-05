// backend/internal/auth/token_blacklist_test.go
package auth

import (
	"testing"
	"time"
)

func TestTokenBlacklist(t *testing.T) {
	// Criar blacklist com intervalo de limpeza curto para testes
	bl := NewTokenBlacklist(100 * time.Millisecond)

	// Token que expira em 1 segundo
	token1 := "token1"
	expiry1 := time.Now().Add(1 * time.Second)
	bl.Add(token1, expiry1)

	// Token que já expirou
	token2 := "token2"
	expiry2 := time.Now().Add(-1 * time.Second)
	bl.Add(token2, expiry2)

	// Verificar se token1 está na blacklist
	if !bl.IsBlacklisted(token1) {
		t.Error("Token1 deveria estar na blacklist")
	}

	// Verificar se token2 não está na blacklist (já expirou)
	if bl.IsBlacklisted(token2) {
		t.Error("Token2 não deveria estar na blacklist")
	}

	// Aguardar para o token1 expirar
	time.Sleep(1100 * time.Millisecond)

	// Token1 deve ter sido removido após expirar
	if bl.IsBlacklisted(token1) {
		t.Error("Token1 deveria ter sido removido da blacklist após expirar")
	}

	// Verificar se a limpeza automática está funcionando
	bl.mutex.RLock()
	tokenCount := len(bl.tokens)
	bl.mutex.RUnlock()

	if tokenCount > 0 {
		t.Errorf("Blacklist deveria estar vazia após limpeza, contém %d tokens", tokenCount)
	}
}
