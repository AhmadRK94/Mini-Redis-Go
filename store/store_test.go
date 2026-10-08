package store

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	s := NewStore()

	err := s.Set("name", "ahmad")
	if err != nil {
		t.Fatal(err)
	}

	got, err := s.Get("name")
	if err != nil {
		t.Fatal(err)
	}
	if got != "ahmad" {
		t.Errorf("expected to get %q but get %q\n", "ahmad", got)
	}
}

func TestSetNX(t *testing.T) {
	s := NewStore()
	if err := s.SetNX("name", "ahmad"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetNX("name", "reza"); err == nil {
		t.Fatal("recreating an existing key.")
	}

}

func TestSetWithTTL(t *testing.T) {
	s := NewStore()
	if err := s.SetWithTTL("temporary", "value", 3*time.Second); err != nil {
		t.Fatal(err)
	}

}

func TestGetTTL(t *testing.T) {
	s := NewStore()

	if got := s.GetTTL("missing"); got != -2 {
		t.Errorf("GetTTL(missing) = %d, want -2.", got)
	}
	if err := s.Set("permanent", "value"); err != nil {
		t.Fatal(err)
	}
	if got := s.GetTTL("permanent"); got != -1 {
		t.Errorf("GetTTL(permanent) = %d, want -1", got)
	}
	if err := s.SetWithTTL("temporary", "value", 3*time.Second); err != nil {
		t.Fatal(err)
	}

	if got := s.GetTTL("temporary"); got <= 0 {
		t.Errorf("GetTTL(temporary) = %d, want a positive TTL", got)
	}
}

func TestDelete(t *testing.T) {
	s := NewStore()

	if err := s.Delete("missing"); err == nil {
		t.Error("Delete(missing) expected to get error but return nil")
	}
	if err := s.Set("name", "ahmad"); err != nil {
		t.Fatal(err)
	}

	if err := s.Delete("name"); err != nil {
		t.Errorf("Delete(name) expected to get no error but get: %q", err)
	}
}

func TestSetExpire(t *testing.T) {
	s := NewStore()
	if err := s.Set("name", "ahmad"); err != nil {
		t.Fatal(err)
	}

	if err := s.SetExpire("missing", 3*time.Second); err == nil {
		t.Error("SetExpire(missing) expected to get error but got nil")
	}

	if err := s.SetExpire("name", 3*time.Second); err != nil {
		t.Errorf("SetExpire(name) expected to get no error but got: %q\n", err)
	}
	if got := s.GetTTL("name"); got <= 0 {
		t.Errorf("GetTTL(name) = %d, want a positive TTL", got)
	}

}

func TestConcurrentSetAndGet(t *testing.T) {
	s := NewStore()

	var wg sync.WaitGroup

	for i := range 100 {
		wg.Go(func() {
			key := fmt.Sprintf("key-%d", i)
			_ = s.Set(key, "value")
			_, _ = s.Get(key)
		})
	}

	wg.Wait()
}
