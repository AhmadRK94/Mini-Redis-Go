package store

import (
	"encoding/json"
	"os"
	"time"
)

func (s *Store) Save(filename string) error {
	s.mu.RLock()

	data, err := json.MarshalIndent(s.data, "", "    ")

	s.mu.RUnlock()

	if err != nil {
		return err
	}

	temp := filename + ".temp"

	if err := os.WriteFile(temp, data, 0644); err != nil {
		return err
	}
	return os.Rename(temp, filename)
}

func (s *Store) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var enteries map[string]Entry
	if err := json.Unmarshal(data, &enteries); err != nil {
		return err
	}

	now := time.Now()

	for key, value := range enteries {
		if !value.ExpiresAt.IsZero() && !now.Before(value.ExpiresAt) {
			continue
		}
		s.data[key] = value
	}
	return nil
}
