package lock

import "testing"

type testLocker struct {

}

func (t testLocker) Lock(name string, quantity int) (*Lock, error) {
	return NewTakenLock("test", t), nil
}

func (t testLocker) Release(lock Lock) error {
	return nil
}

func TestNewTakenLock(t *testing.T) {
	locker := &testLocker{}

	lock := NewTakenLock("test", locker)

	if lock.Id() != "test" {
		t.Errorf("lock.Id expected to be %s, got %s", "test", lock.Id())
	}

	if lock.Taken() != true {
		t.Errorf("lock.Taken expected to be true, got %v", lock.Taken())
	}
}

func TestNewUntakenLock(t *testing.T) {
	lock := NewUntakenLock()

	if lock.Id() != "" {
		t.Errorf("lock.Id expected to be %s, got %s", "", lock.Id())
	}

	if lock.Taken() != false {
		t.Errorf("lock.Taken expected to be false, got %v", lock.Taken())
	}
}