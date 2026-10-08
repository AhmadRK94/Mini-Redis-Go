package store

import (
	"errors"
	"sync"
	"time"
)

type Entry struct {
	Value     string
	ExpiresAt time.Time
}

type Store struct {
	mu   sync.RWMutex
	data map[string]Entry
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]Entry),
	}
}

func (s *Store) set(key, value string, forceUpdate bool) error {
	if key == "" {
		return errors.New("Invalid key: Empty key.")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.data[key]
	if exists && !forceUpdate {
		return errors.New("Key already exist.\n")
	}
	e := Entry{
		Value:     value,
		ExpiresAt: time.Time{},
	}

	s.data[key] = e
	return nil
}

// Set create or overwrite the existing data.
func (s *Store) Set(key, value string) error {
	return s.set(key, value, true)
}

// SetNX create if it doesn't exist.
func (s *Store) SetNX(key, value string) error {
	return s.set(key, value, false)
}

// SetTTL create or overwrite the existing data with expire time.
func (s *Store) SetWithTTL(key, value string, ttl time.Duration) error {
	if key == "" {
		return errors.New("Invalid key: Empty key.")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.data[key]
	if exists {
		return errors.New("Key already exist.\n")
	}
	e := Entry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}

	s.data[key] = e
	return nil
}

// set expiration for existing key
func (s *Store) SetExpire(key string, ttl time.Duration) error {
	if key == "" {
		return errors.New("Invalid key: Empty key.")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, exists := s.data[key]
	if !exists {
		return errors.New("Key doesn't exist.\n")
	}
	v.ExpiresAt = time.Now().Add(ttl)
	s.data[key] = v
	return nil

}

// return value for specify key if exist.
func (s *Store) Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("Invalid key: Empty key.")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, exists := s.data[key]

	if !exists {
		return "", errors.New("key doesn't exist.\n")
	}

	expired := !v.ExpiresAt.IsZero() &&
		!time.Now().Before(v.ExpiresAt)

	if expired {
		return "", errors.New("key doesn't exist.\n")
	}

	return v.Value, nil
}

// get expiration if exist. -1 for data with no expiration, -2 if not exist, positive values is in seconds
func (s *Store) GetTTL(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, exists := s.data[key]
	if !exists {
		return -2
	}
	if v.ExpiresAt.IsZero() {
		return -1
	}
	ttl := int(time.Until(v.ExpiresAt).Seconds())
	if ttl < 0 {
		return -2
	}
	return ttl
}

// delete data for specified key
func (s *Store) Delete(key string) error {
	if key == "" {
		return errors.New("Invalid key: Empty key.")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.data[key]
	if !exists {
		return errors.New("Key doesn't exist.\n")
	}
	delete(s.data, key)
	return nil
}
