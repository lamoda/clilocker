package lock

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path"
)
import "github.com/gofrs/flock"

type LocalFileLocker struct {
}

func NewLocalFileLocker() Locker {
	return &LocalFileLocker{}
}

func (l *LocalFileLocker) Lock(name string, quantity int) (*Lock, error) {
	for i := 0; i < quantity; i++ {
		lockId := createLockId(name, i)
		fileLock := flock.New(lockFileFromId(lockId))
		locked, err := fileLock.TryLock()

		if err != nil {
			return nil, err
		}

		if locked {
			return NewTakenLock(lockId, l), nil
		}
	}

	return NewUntakenLock(), nil
}

func (l *LocalFileLocker) Release(lock Lock) error {
	lockName := lock.Id()
	fileLock := flock.New(lockFileFromId(lockName))
	return fileLock.Unlock()
}

func createLockId(name string, index int) string {
	return fmt.Sprintf("lock_%s_%d", name, index)
}

func lockFileFromId(lockId string) string {
	return path.Join(os.TempDir(), hex.EncodeToString(md5.New().Sum([]byte(lockId))))
}
