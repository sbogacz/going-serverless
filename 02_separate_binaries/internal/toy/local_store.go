package toy

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/sbogacz/gophercon18-kickoff-talk/third/internal/httperrs"
)

// LocalStore is an in-memory implementation of our
// store interface
type LocalStore struct {
	store map[string]string
	lock  *sync.RWMutex
}

// NewLocalStore returns a local implementation of the
// Store interface
func NewLocalStore() *LocalStore {
	return &LocalStore{
		store: make(map[string]string),
		lock:  &sync.RWMutex{},
	}
}

// Set stores the given data at the given key
func (l *LocalStore) Set(ctx context.Context, key, data string) error {
	l.lock.Lock()
	l.store[key] = data
	l.lock.Unlock()

	return nil
}

// Get retrieves the data stored at the given key
func (l *LocalStore) Get(ctx context.Context, key string) (string, error) {
	l.lock.RLock()
	data, ok := l.store[key]
	l.lock.RUnlock()

	if !ok {
		return "", httperrs.NotFound(errors.Errorf("no item found at %s", key), "")
	}
	return data, nil
}

// Del removes the data stored at the given key
func (l *LocalStore) Del(ctx context.Context, key string) error {
	l.lock.RLock()
	_, ok := l.store[key]
	l.lock.RUnlock()

	if !ok {
		return httperrs.NotFound(errors.Errorf("no item found at %s", key), "")
	}

	// if we have the item, delete
	l.lock.Lock()
	delete(l.store, key)
	l.lock.Unlock()

	return nil
}
