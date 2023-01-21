package debouncer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/DanielRenne/GoCore/core/debouncer"
)

func Example() {
	wait := 5 * time.Second
	d := debouncer.New(wait).WithTriggered(func() {
		fmt.Println("Trigger") // Triggered func will be called after 5 seconds from last SendSignal().
	})

	fmt.Println("Action 1")
	d.SendSignal()

	time.Sleep(1 * time.Second)

	fmt.Println("Action 2")
	d.SendSignal()

	// After 5 seconds, the trigger will be called.
	//Previous `SendSignal()` will be ignore to trigger the triggered function.
	<-d.Done()
}

func createIncrementCount(counter int) (*int, func()) {
	return &counter, func() {
		fmt.Println("Triggered")
		counter++
	}
}

func TestDebounceDoBeforeExpired(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	d := debouncer.New(200 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(1)

	d.Do(func() {
		fmt.Println("Action 1")
	})

	time.Sleep(50 * time.Millisecond)

	d.Do(func() {
		fmt.Println("Action 2")
	})

	<-d.Done()

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDebounceDoAfterExpired(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	d := debouncer.New(200 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(2)

	d.Do(func() {
		fmt.Println("Action 1")
	})

	<-d.Done()

	d.Do(func() {
		fmt.Println("Action 2")
	})

	<-d.Done()

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDebounceMixed(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	d := debouncer.New(200 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(2)

	d.Do(func() {
		fmt.Println("Action 1")
	})

	d.Do(func() {
		fmt.Println("Action 2")
	})

	<-d.Done()

	d.Do(func() {
		fmt.Println("Action 3")
	})

	<-d.Done()

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDebounceWithoutTriggeredFunc(t *testing.T) {
	d := debouncer.New(200 * time.Millisecond)

	d.Do(func() {
		fmt.Println("Action 1")
	})
	<-d.Done()

	fmt.Println("debouncer.Do() finished successfully!")
}

func TestDebounceSendSignal(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	d := debouncer.New(200 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(1)

	d.SendSignal()
	<-d.Done()

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDebounceUpdateTriggeredFuncBeforeDuration(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	d := debouncer.New(200 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(2)

	d.SendSignal()
	time.Sleep(50 * time.Millisecond)

	d.UpdateTriggeredFunc(func() {
		*countPtr += 2
	})
	<-d.Done()

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDebounceUpdateTriggeredFuncAfterDuration(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	d := debouncer.New(200 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(3)

	d.SendSignal()
	<-d.Done()

	d.UpdateTriggeredFunc(func() {
		*countPtr += 2
	})
	d.SendSignal()
	<-d.Done()

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDebounceCancel(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	d := debouncer.New(200 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(0)

	d.SendSignal()
	time.Sleep(50 * time.Millisecond)

	d.Cancel()
	time.Sleep(400 * time.Millisecond)

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDebounceUpdateDuration(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	d := debouncer.New(600 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(1)

	d.UpdateTimeDuration(200 * time.Millisecond)
	d.SendSignal()
	<-d.Done()

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDebounceUpdateDurationAfterSendSignal(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	d := debouncer.New(400 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(1)

	d.SendSignal()
	time.Sleep(200 * time.Millisecond)

	d.UpdateTimeDuration(600 * time.Millisecond)
	<-d.Done()

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDone(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	d := debouncer.New(200 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(2)

	d.SendSignal()
	<-d.Done()

	d.SendSignal()
	<-d.Done()

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDoneInGoroutine(t *testing.T) {
	countPtr, incrementCount := createIncrementCount(0)
	d := debouncer.New(200 * time.Millisecond).WithTriggered(incrementCount)
	expectedCounter := int(3)

	d.SendSignal()
	go func() {
		<-d.Done() // awaits for the second send signal to complete
		*countPtr += 2
	}()

	d.SendSignal() // after 1 milliseconds, unblock done channel in 2 goroutines
	<-d.Done()

	time.Sleep(200 * time.Millisecond)

	if *countPtr != expectedCounter {
		t.Errorf("Expected count %d, was %d", expectedCounter, *countPtr)
	}
}

func TestDoneHangBeforeSendSignal(t *testing.T) {
	d := debouncer.New(200 * time.Millisecond).WithTriggered(func() {})
	select {
	case <-d.Done():
		t.Error("Done() must hang when being called before SendSignal()")
	case <-time.After(time.Second):
	}
}

func TestDoneHangIfBeingCalledTwice(t *testing.T) {
	d := debouncer.New(200 * time.Millisecond).WithTriggered(func() {})
	d.SendSignal()
	<-d.Done()

	select {
	case <-d.Done():
		t.Error("Done() must hang if being called twice")
	case <-time.After(time.Second):
	}
}
