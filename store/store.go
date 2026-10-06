package store

import (
	"errors"
	"sync"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

func (s *Store) set(key, value string, forceUpdate bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.data[key]
	if exists && !forceUpdate {
		return errors.New("Key already exist.\n")
	}
	s.data[key] = value
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

// return value for specify key if exist.
func (s *Store) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, exist := s.data[key]
	if !exist {
		return "", errors.New("Key doesn't exist.\n")
	}
	return v, nil
}
