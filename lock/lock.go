package lock

import "fmt"
import "sync"

type WriteLock struct{ lock *sync.Map }

func New() *WriteLock {
	lock := &WriteLock{lock: &sync.Map{}}
	return lock
}

func (lock *WriteLock) GetWrite(name string) *sync.RWMutex {
	mutex, _ := lock.lock.LoadOrStore(name, sync.RWMutex{})
	switch m := mutex.(type) {
	case sync.RWMutex:
		m.Lock()
		return &m
	default:
		fmt.Println("sfsdb: Lock store poison")
		var mutex sync.RWMutex
		mutex.Lock()
		return &mutex
	}
}

func (lock *WriteLock) DoneWrite(name string, mutex *sync.RWMutex) {
	mutex.Unlock()
	lock.lock.Delete(name)
}

func (lock *WriteLock) GetRead(name string) *sync.RWMutex {
	mutex, _ := lock.lock.LoadOrStore(name, sync.RWMutex{})
	switch m := mutex.(type) {
	case sync.RWMutex:
		m.RLock()
		return &m
	default:
		fmt.Println("sfsdb: Lock store poison")
		var mutex sync.RWMutex
		mutex.RLock()
		return &mutex
	}
}

func (lock *WriteLock) DoneRead(name string, mutex *sync.RWMutex) {
	mutex.RUnlock()
	lock.lock.Delete(name)
}
