package mqtt

import (
	"errors"
	"testing"
	"time"
)

func TestWaitTimeout(t *testing.T) {
	b := baseToken{}

	if b.WaitTimeout(time.Second) {
		t.Fatal("Should have failed")
	}

	// Now lets confirm that WaitTimeout returns
	// setError() grabs the mutex which previously caused issues
	// when there is a result (it returns true in this case)
	b = baseToken{complete: make(chan struct{})}
	go func(bt *baseToken) {
		bt.setError(errors.New("test error"))
	}(&b)
	if !b.WaitTimeout(5 * time.Second) {
		t.Fatal("Should have succeeded")
	}
}
