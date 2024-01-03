package lock

import "sync"

type TryLocker struct {
	mu    sync.Mutex
	lockC chan struct{}
}

func NewTryLocker() *TryLocker {
	return &TryLocker{
		lockC: make(chan struct{}, 1),
	}
}

// TryLock attempts to acquire the lock without blocking. Returns true if successful.
func (t *TryLocker) TryLock() bool {
	select {
	case t.lockC <- struct{}{}:
		t.mu.Lock() // Lock acquired
		return true
	default:
		return false // Lock not acquired
	}
}

// Unlock releases the lock.
func (t *TryLocker) Unlock() {
	t.mu.Unlock()
	<-t.lockC
}
