package thread

import (
	"fmt"
	"time"
)

type ExampleRunnable struct {
	Counter int
}

func (t *ExampleRunnable) Run(stop chan bool) error {
	fmt.Println("ExampleRunnable.Run()")
	for {
		select {
		case <-stop:
			fmt.Println(" <-stop")
			return nil
		default:
			fmt.Printf("ExampleRunnable.Counter = %d\n", t.Counter)
			t.Counter++
			time.Sleep(400 * time.Millisecond)
		}
	}
}

func ExampleThread() {
	runnable := &ExampleRunnable{Counter: 0}
	thread := New(runnable)

	fmt.Println("Thread.Start()")
	thread.Start()

	time.Sleep(1 * time.Second)
	fmt.Println("Thread.Stop()")
	thread.Stop()
	thread.Join()
	fmt.Println("Stopped")

	time.Sleep(1 * time.Second)
	fmt.Println("Thread.Start()")
	thread.Start()

	time.Sleep(1 * time.Second)
	fmt.Println("Thread.Stop()")
	thread.Stop()
	thread.Join()
	fmt.Println("Stopped")

	time.Sleep(1 * time.Second)
	fmt.Println("Exit")

	// Output:
	//
	// Thread.Start()
	// ExampleRunnable.Run()
	// ExampleRunnable.Counter = 0
	// ExampleRunnable.Counter = 1
	// ExampleRunnable.Counter = 2
	// Thread.Stop()
	//  <-stop
	// Stopped
	// Thread.Start()
	// ExampleRunnable.Run()
	// ExampleRunnable.Counter = 3
	// ExampleRunnable.Counter = 4
	// ExampleRunnable.Counter = 5
	// Thread.Stop()
	//  <-stop
	// Stopped
	// Exit

}
