package lock

type Lock struct {
	id     string
	locker Locker
	taken  bool
}

func NewTakenLock(id string, locker Locker) *Lock {
	return &Lock{
		id:     id,
		locker: locker,
		taken:  true,
	}
}

func NewUntakenLock() *Lock {
	return &Lock{
		id:     "",
		locker: nil,
		taken:  false,
	}
}

func (l Lock) Id() string {
	return l.id
}

func (l Lock) Taken() bool {
	return l.taken
}

func (l Lock) Release() error {
	if !l.taken {
		return nil
	}
	err := l.locker.Release(l)
	if err != nil {
		return err
	}

	l.taken = false
	l.locker = nil

	return nil
}
