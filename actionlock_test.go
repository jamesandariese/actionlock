package actionlock

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func Example() {
	al := New(func(al *ActionLock) {
		fmt.Println("Setting to", al.Get())
	})

	wg := &sync.WaitGroup{}
	wg.Add(4)
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		fmt.Println("A")
		al.LockValue(1)
		time.Sleep(900 * time.Millisecond)
		al.UnlockValue(1)
	}()
	go func() {
		defer wg.Done()
		time.Sleep(200 * time.Millisecond)
		fmt.Println("B")
		al.LockValue(1)
		time.Sleep(100 * time.Millisecond)
		al.UnlockValue(1)
	}()
	go func() {
		// because this goroutine will request a write lock, it will block
		// further read locks.  goroutine D will run last even though it could
		// have run during A's time.  This is ensures that other values aren't
		// starved.
		defer wg.Done()
		time.Sleep(300 * time.Millisecond)
		fmt.Println("C")
		al.LockValue(2)
		time.Sleep(100 * time.Millisecond)
		al.UnlockValue(2)
	}()
	go func() {
		defer wg.Done()
		time.Sleep(400 * time.Millisecond)
		fmt.Println("D")
		al.LockValue(1)
		time.Sleep(100 * time.Millisecond)
		al.UnlockValue(1)
	}()
	wg.Wait()
	// Output:
	// A
	// Setting to 1
	// B
	// C
	// D
	// Setting to 2
	// Setting to 1
}

func TestRace(t *testing.T) {
	al := New(func(al *ActionLock) {
		// fmt.Println("sup", al.Get())
	})
	wg := sync.WaitGroup{}
	rounds := 100000
	wg.Add(rounds)
	for i := 0; i < rounds; i++ {
		go func(v, ms int) {
			defer wg.Done()
			al.LockValue(v)
			defer al.UnlockValue(v)
			runtime.Gosched()
		}(i, 1)
	}
	wg.Wait()
}
