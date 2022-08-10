/*
 * Copyright (c) 2022 IBM Corp and others.
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
 *    Matt Brittan
 */

package mqtt

import (
	"fmt"
	"testing"
	"time"
)

func Test_BasicStatusOperations(t *testing.T) {
	t.Parallel()
	s := connectionStatus{}
	if s.ConnectionStatus() != disconnected {
		t.Fatalf("Expected disconnected; got: %v", s.ConnectionStatus())
	}

	// Normal connection and disconnection
	cf, err := s.Connecting()
	if err != nil {
		t.Fatalf("Error connecting: %v", err)
	}
	if s.ConnectionStatus() != connecting {
		t.Fatalf("Expected connecting; got: %v", s.ConnectionStatus())
	}
	if err = cf(true); err != nil {
		t.Fatalf("Error completing connection: %v", err)
	}
	if s.ConnectionStatus() != connected {
		t.Fatalf("Expected connected; got: %v", s.ConnectionStatus())
	}

	// reconnect so we test all statuses
	rf, err := s.ConnectionLost(true)
	if err != nil {
		t.Fatalf("Error calling connection lost: %v", err)
	}
	if s.ConnectionStatus() != disconnecting {
		t.Fatalf("Expected disconnecting; got: %v", s.ConnectionStatus())
	}
	if cf, err = rf(true); err != nil {
		t.Fatalf("Error completing disconnection portion of reconnect: %v", err)
	}
	if s.ConnectionStatus() != reconnecting {
		t.Fatalf("Expected reconnecting; got: %v", s.ConnectionStatus())
	}
	if err = cf(true); err != nil {
		t.Fatalf("Error completing reconnection: %v", err)
	}
	if s.ConnectionStatus() != connected {
		t.Fatalf("Expected connected(2); got: %v", s.ConnectionStatus())
	}

	// And disconnect
	df, err := s.Disconnecting()
	if err != nil {
		t.Fatalf("Error disconnecting: %v", err)
	}
	if s.ConnectionStatus() != disconnecting {
		t.Fatalf("Expected disconnecting; got: %v", s.ConnectionStatus())
	}
	df()
	if s.ConnectionStatus() != disconnected {
		t.Fatalf("Expected disconnected; got: %v", s.ConnectionStatus())
	}
}

// Test_AdvancedStatusOperations checks a few of the more unusual transitions
func Test_AdvancedStatusOperations(t *testing.T) {
	t.Parallel()

	// Aborted connection (i.e. user triggered)
	s := connectionStatus{}
	if s.ConnectionStatus() != disconnected {
		t.Fatalf("Expected disconnected; got: %v", s.ConnectionStatus())
	}

	// Normal connection and disconnection
	cf, err := s.Connecting()
	if err != nil {
		t.Fatalf("Error connecting: %v", err)
	}
	if s.ConnectionStatus() != connecting {
		t.Fatalf("Expected connecting; got: %v", s.ConnectionStatus())
	}
	if err = cf(false); err != nil { // Unsuccessful connection (e.g. user aborted connection)
		t.Fatalf("Error completing connection: %v", err)
	}
	if s.ConnectionStatus() != disconnected {
		t.Fatalf("Expected disconnected; got: %v", s.ConnectionStatus())
	}

	// Connection lost - no reconnection requested
	s = connectionStatus{status: connected}
	rf, err := s.ConnectionLost(false)
	if err != nil {
		t.Fatalf("Error connecting: %v", err)
	}
	if s.ConnectionStatus() != disconnecting {
		t.Fatalf("Expected disconnecting; got: %v", s.ConnectionStatus())
	}
	cf, err = rf(true) // argument should be ignored as no reconnect was requested
	if cf != nil {
		t.Fatalf("Function to complete reconnection should not be returned (as reconnection not requested)")
	}
	if err != nil {
		t.Fatalf("Error completing connection lost operation: %v", err)
	}
	if s.ConnectionStatus() != disconnected {
		t.Fatalf("Expected disconnected; got: %v", s.ConnectionStatus())
	}

	// Aborted reconnection - stage 1 (i.e. user triggered whist disconnect in progress)
	s = connectionStatus{status: connected}
	rf, err = s.ConnectionLost(true)
	if err != nil {
		t.Fatalf("Error connecting: %v", err)
	}
	if s.ConnectionStatus() != disconnecting {
		t.Fatalf("Expected disconnecting; got: %v", s.ConnectionStatus())
	}
	cf, err = rf(false)
	if cf != nil {
		t.Fatalf("Function to complete reconnection should not be returned (as reconnection not requested)")
	}
	if err != nil {
		t.Fatalf("Error completing connection lost operation: %v", err)
	}
	if s.ConnectionStatus() != disconnected {
		t.Fatalf("Expected disconnected; got: %v", s.ConnectionStatus())
	}

	// Aborted reconnection - stage 2 (i.e. user triggered whist disconnect in progress)
	s = connectionStatus{status: connected}
	rf, err = s.ConnectionLost(true)
	if err != nil {
		t.Fatalf("Error connecting: %v", err)
	}
	if s.ConnectionStatus() != disconnecting {
		t.Fatalf("Expected disconnecting; got: %v", s.ConnectionStatus())
	}
	cf, err = rf(true)
	if err != nil {
		t.Fatalf("Error completing connection lost operation: %v", err)
	}
	if cf == nil {
		t.Fatalf("Function to complete reconnection should be returned (as reconnection requested)")
	}
	if s.ConnectionStatus() != reconnecting {
		t.Fatalf("Expected reconnecting; got: %v", s.ConnectionStatus())
	}
	if err = cf(false); err != nil {
		t.Fatalf("Error completing reconnection: %v", err)
	}
	if s.ConnectionStatus() != disconnected {
		t.Fatalf("Expected disconnected; got: %v", s.ConnectionStatus())
	}
}

func Test_AbortedConnection(t *testing.T) {
	t.Parallel()
	s := connectionStatus{}
	if s.ConnectionStatus() != disconnected {
		t.Fatalf("Expected disconnected; got: %v", s.ConnectionStatus())
	}

	// Start Connection
	cf, err := s.Connecting()
	if err != nil {
		t.Fatalf("Error connecting: %v", err)
	}
	if s.ConnectionStatus() != connecting {
		t.Fatalf("Expected connecting; got: %v", s.ConnectionStatus())
	}

	// Another goroutine calls Disconnect
	discErr := make(chan error)
	go func() {
		dfFn, err := s.Disconnecting()
		discErr <- err
		dfFn()
		close(discErr)
	}()
	time.Sleep(time.Millisecond) // Provide time for Disconnect call to run
	if s.ConnectionStatus() != disconnecting {
		t.Fatalf("Expected disconnecting; got: %v", s.ConnectionStatus())
	}
	select {
	case err = <-discErr:
		t.Fatalf("Disconnecting must block until connection attempt terminates: %v", err)
	default:
	}

	err = cf(true) // status should not matter
	if err != errAbortConnection {
		t.Fatalf("Expected errAbortConnection got: %v", err)
	}
	if s.ConnectionStatus() != disconnecting {
		t.Fatalf("Expected disconnecting; got: %v", s.ConnectionStatus())
	}

	select {
	case err = <-discErr:
		if err != nil {
			t.Fatalf("Did not expect an error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatalf("Timeout waiting for goroutine to complete")
	}

	time.Sleep(time.Millisecond) // Provide time for other goroutine to complete
	if s.ConnectionStatus() != disconnected {
		t.Fatalf("Expected disconnected; got: %v", s.ConnectionStatus())
	}
	select {
	case <-discErr: // channel should be closed
	default:
		t.Fatalf("Completion of connect should unblock Disconnecting call")
	}
}

func Test_AbortedReConnection(t *testing.T) {
	t.Parallel()
	s := connectionStatus{status: connected} // start in connected state
	if s.ConnectionStatus() != connected {
		t.Fatalf("Expected connected; got: %v", s.ConnectionStatus())
	}

	// Connection is lost but we want to reconnect
	lhf, err := s.ConnectionLost(true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Another goroutine calls Disconnect
	discErr := make(chan error)
	go func() {
		dfFn, err := s.Disconnecting()
		if dfFn != nil {
			discErr <- fmt.Errorf("should not get a functiuon back from s.Disconnecting in this case")
			return
		}
		discErr <- err
	}()
	time.Sleep(time.Millisecond) // Provide time for Disconnect call to run
	if s.ConnectionStatus() != disconnecting {
		t.Fatalf("Expected disconnecting; got: %v", s.ConnectionStatus())
	}
	select {
	case err = <-discErr:
		t.Fatalf("Disconnecting must block until reconnection attempt terminates: %v", err)
	default:
	}

	cf, err := lhf(true) // status should not matter
	if cf != nil {
		t.Fatalf("As Disconnect has been called we should not have any ability to continue")
	}
	if err != errDisconnectionRequested {
		t.Fatalf("Expected errDisconnectionRequested got: %v", err)
	}
	if s.ConnectionStatus() != disconnected {
		t.Fatalf("Expected disconnected; got: %v", s.ConnectionStatus())
	}

	select {
	case err = <-discErr:
		if err != errAlreadyDisconnected {
			t.Fatalf("Expected errAlreadyDisconnected got: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatalf("Timeout waiting for goroutine to complete")
	}
}

// Test_ConnectionLostDuringConnect don't really expect this to happen due to connMu
// If it does happen and reconnect is true the results would not be great
func Test_ConnectionLostDuringConnect(t *testing.T) {
	t.Parallel()
	s := connectionStatus{}
	if s.ConnectionStatus() != disconnected {
		t.Fatalf("Expected disconnected; got: %v", s.ConnectionStatus())
	}

	// Start Connection
	cf, err := s.Connecting()
	if err != nil {
		t.Fatalf("Error connecting: %v", err)
	}
	if s.ConnectionStatus() != connecting {
		t.Fatalf("Expected connecting; got: %v", s.ConnectionStatus())
	}

	// Another goroutine calls ConnectionLost (don't expect this to every actually happen but...)
	clErr := make(chan error)
	go func() {
		_, err := s.ConnectionLost(false)
		clErr <- err
	}()
	time.Sleep(time.Millisecond) // Provide time for Disconnect call to run
	if s.ConnectionStatus() != disconnecting {
		t.Fatalf("Expected disconnecting; got: %v", s.ConnectionStatus())
	}
	select {
	case err = <-clErr:
		t.Fatalf("ConnectionLost must block until connection attempt terminates: %v", err)
	default:
	}

	err = cf(true) // status should not matter
	if err != errAbortConnection {
		t.Fatalf("Expected errAbortConnection got: %v", err)
	}
	if s.ConnectionStatus() != disconnecting {
		t.Fatalf("Expected disconnecting; got: %v", s.ConnectionStatus())
	}

	select {
	case err = <-clErr:
		if err != errAlreadyDisconnected {
			t.Fatalf("Expected errAlreadyDisconnected got: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatalf("Timeout waiting for goroutine to complete")
	}
}

/*
clErr := make(chan error)
	go func() {
		rf, err := s.ConnectionLost(false)
		clErr <- err
		cf, err := rf(false)
		if err != errAlreadyDisconnected {
			clErr <- fmt.Errorf("expected errAlreadyDisconnected got %v", err)
		}
		if cf != nil {
			clErr <- fmt.Errorf("cf is not nil")
		}
		close(clErr)

	}()
	time.Sleep(time.Millisecond) // Provide time for Disconnect call to run
	if s.ConnectionStatus() != disconnecting {
		t.Fatalf("Expected disconnecting; got: %v", s.ConnectionStatus())
	}
	select {
	case err = <-clErr:
		t.Fatalf("ConnectionLost must block until connection attempt terminates: %v", err)
	default:
	}

	err = cf(true) // status should not matter
	if err != errAbortConnection {
		t.Fatalf("Expected errAbortConnection got: %v", err)
	}
	if s.ConnectionStatus() != disconnecting {
		t.Fatalf("Expected disconnecting; got: %v", s.ConnectionStatus())
	}

	select {
	case err = <-clErr:
		if err != nil {
			t.Fatalf("Did not expect an error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatalf("Timeout waiting for goroutine to complete")
	}

	time.Sleep(time.Millisecond) // Provide time for other goroutine to complete
	if s.ConnectionStatus() != disconnected {
		t.Fatalf("Expected disconnected; got: %v", s.ConnectionStatus())
	}
	select {
	case <-clErr: // channel should be closed
	default:
		t.Fatalf("Completion of connect should unblock Disconnecting call")
	}
*/

/*
// TODO - Test aborting functions etc

disconnected -> `Connecting()` -> connecting -> `connCompletedFn(true)` -> connected
connected -> `Disconnecting()` -> disconnecting -> `disconnectCompletedFn()` -> disconnected
connected -> `ConnectionLost(false)` -> disconnecting -> `connectionLostHandledFn(true/false)` -> disconnected
connected -> `ConnectionLost(true)` -> disconnecting -> `connectionLostHandledFn(true)` -> connected

Unfortunately the above workflows are complicated by the fact that `Disconnecting()` or `ConnectionLost()` may,
potentially, be called at any time (i.e.whilst in the middle of transitioning between states).If this happens:

* The state will be set to disconnecting (which will prevent any request to move the status to connected)
* The call to `Disconnecting()`/`ConnectionLost()` will block until the previously active call completes and then
handle the disconnection.

Reading the tests (unit_client_test.go ) might help understand these rules.
*/
