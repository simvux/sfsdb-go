package lock

import "sync"

type WriteLock struct{ lock sync.Map }

func New() *WriteLock {
	lock := &WriteLock{sync.Map{}}
	return lock
}

func (lock *WriteLock) Get(name string) *sync.Mutex {
	mutex := &sync.Mutex{}
	mutex.Lock()
	lock.lock.LoadOrStore(name, mutex)
	return mutex
}

func (lock *WriteLock) Done(name string, mutex *sync.Mutex) {
	mutex.Unlock()
	lock.lock.Delete(name)
}
