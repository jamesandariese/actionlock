package actionlock

import (
	"log"
	"sync"
)

type ActionLockCallback func(*ActionLock)

type ActionLock struct {
	Callback ActionLockCallback
	value    interface{}
	lock     *sync.RWMutex
}

func New(cb ActionLockCallback) *ActionLock {
	return &ActionLock{
		Callback: cb,
		value:    -1,
		lock:     &sync.RWMutex{},
	}
}

func (al *ActionLock) Get() interface{} {
	return al.value
}

func (al *ActionLock) LockValue(value interface{}) {
	// first, we'll take a read lock to determine the current value
	// this way we'll have a chance to short circuit the locking
	// for readers that are using the same value.
	// then lock for write if it's wrong, set the value, and proceed.
	al.lock.RLock()

	if al.value == value {
		// return with the lock asserted for reading
		return
	}
	al.lock.RUnlock() // unassert it instead so we can take the write lock

	func() {
		// protect the lock call 'cause we don't control the callback
		al.lock.Lock()
		defer al.lock.Unlock()
		if al.value != value {
			al.value = value
			al.Callback(al)
		}
	}()

	al.lock.RLock() // now reassert the read lock
}

func (al *ActionLock) UnlockValue(value interface{}) {
	// it is an error for al.lock to not be held for reading upon entry
	defer al.lock.RUnlock()

	if al.value != value {
		log.Fatalf("ActionLock value was not %d (was actually %d)", al.value, value)
	}
}
