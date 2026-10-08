package store

import "time"

// Cleanup functionality for deleting expired data
func (s *Store) cleanupExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	for key, entry := range s.data {
		if !entry.ExpiresAt.IsZero() && !now.Before(entry.ExpiresAt) {
			delete(s.data, key)
		}
	}
}

// cleanup goroutine with specified ticker and graceful shutdown.
func (s *Store) StartCleanup(
	interval time.Duration,
	stop <-chan struct{},
) {
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.cleanupExpired()
			case <-stop:
				return
			}
		}
	}()
}
