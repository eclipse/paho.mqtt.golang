/*
 * Copyright (c) 2021 IBM Corp and others.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v2.0
 * and Eclipse Distribution License v1.0 which accompany this distribution.
 *
 * The Eclipse Public License is available at
 *    https://www.eclipse.org/legal/epl-2.0/
 * and the Eclipse Distribution License is available at
 *   http://www.eclipse.org/org/documents/edl-v10.php.
 *
 * Contributors:
 *    Allan Stockdill-Mander
 */

package mqtt

import (
	"errors"
	"sync"
	"testing"
	"time"
)

// Running this test with go test -race will expose a race
// if Result() is modifiable in a thread unsafe way.
// This test illustrates the case where a client receives a SubscribeToken
// abuses it and modifies implementation details for this library
func TestSubscribeToken_Result(t *testing.T) {
	s := SubscribeToken{
		baseToken: baseToken{complete: make(chan struct{})},
		subs:      []string{"mysuv"},
		subResult: map[string]byte{"sd": 0x1},
		messageID: 0,
	}
	w := sync.WaitGroup{}
	w.Add(2)
	go func() {
		s.Result()["sd"] = 0x2
		w.Done()
	}()
	go func() {
		s.Result()["sd"] = 0x3
		w.Done()
	}()
	if s.Result()["sd"] != 0x1 && s.Result()["sd"] != 0x2 && s.Result()["sd"] != 0x3 {
		t.Fatal("Unexpected")
	}
}

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
