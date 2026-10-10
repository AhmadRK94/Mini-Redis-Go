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
func (s *Store) StartCleanup(interval time.Duration) {
	s.stopClean = make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.cleanupExpired()
			case <-s.stopClean:
				return
			}
		}
	}()
}

func (s *Store) StopCleanup() {
	s.cleanOnce.Do(func() {
		close(s.stopClean)
	})
}
