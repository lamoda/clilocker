package lock

type Locker interface {
	Lock(name string, quantity int) (*Lock, error)

	Release(lock Lock) error
}
