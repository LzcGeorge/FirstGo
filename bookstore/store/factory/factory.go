package factory

import (
	"bookstore/store"
	"fmt"
	"sync"
)

var (
	providerMu sync.RWMutex
	providers  = make(map[string]store.Store)
)

func Regitster(name string, provider store.Store) {
	providerMu.Lock()
	defer providerMu.Unlock()
	if provider == nil {
		panic("store: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("store: Register called twice for provider " + name)
	}
	providers[name] = provider
}

func New(name string) (store.Store, error) {
	providerMu.RLock()
	p, ok := providers[name]
	providerMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("store: unknown provider %q (forgotten import?)", name)
	}
	return p, nil
}
