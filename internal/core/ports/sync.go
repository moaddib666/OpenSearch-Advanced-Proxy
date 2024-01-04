package ports

type TryLocker interface {
	TryLock() bool
	Unlock()
}
