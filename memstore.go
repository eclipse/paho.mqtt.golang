/*
 * Copyright (c) 2013 IBM Corp.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package mqtt

import (
	"sync"
)

// MemoryStore implements the store interface to provide a "persistence"
// mechanism wholly stored in memory. This is only useful for
// as long as the client instance exists.
type MemoryStore struct {
	sync.RWMutex
	messages map[string]*Message
	opened   bool
	t        *Tracer
}

// NewMemoryStore returns a pointer to a new instance of
// MemoryStore, the instance is not initialized and ready to
// use until Open() has been called on it.
func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		messages: nil,
		opened:   false,
		t:        nil,
	}
	return store
}

// Open initializes a MemoryStore instance.
func (store *MemoryStore) Open() {
	defer store.Unlock()
	store.Lock()
	store.messages = make(map[string]*Message)
	store.opened = true
	store.t.Trace_V(STR, "memorystore initialized")
}

// Put takes a key and a pointer to a Message and stores the
// message.
func (store *MemoryStore) Put(key string, message *Message) {
	defer store.Unlock()
	store.Lock()
	chkcond(store.opened)
	store.messages[key] = message
}

// Get takes a key and looks in the store for a matching Message
// returning either the Message pointer or nil.
func (store *MemoryStore) Get(key string) *Message {
	defer store.RUnlock()
	store.RLock()
	chkcond(store.opened)
	mid := key2mid(key)
	m := store.messages[key]
	if m == nil {
		store.t.Trace_V(STR, "memorystore get: message %v not found", mid)
	} else {
		store.t.Trace_V(STR, "memorystore get: message %v found", mid)
	}
	return m
}

// All returns a slice of strings containing all the keys currently
// in the MemoryStore.
func (store *MemoryStore) All() []string {
	defer store.RUnlock()
	store.RLock()
	chkcond(store.opened)
	keys := []string{}
	for k, _ := range store.messages {
		keys = append(keys, k)
	}
	return keys
}

// Del takes a key, searches the MemoryStore and if the key is found
// deletes the Message pointer associated with it.
func (store *MemoryStore) Del(key string) {
	defer store.Unlock()
	store.Lock()
	mid := key2mid(key)
	m := store.messages[key]
	if m == nil {
		store.t.Trace_V(STR, "memorystore del: message %v not found", mid)
	} else {
		store.messages[key] = nil
		store.t.Trace_V(STR, "memorystore del: message %v was deleted", mid)
	}
}

// Close will disallow modifications to the state of the store.
func (store *MemoryStore) Close() {
	defer store.Unlock()
	store.Lock()
	chkcond(store.opened)
	store.opened = false
	store.t.Trace_V(STR, "memorystore closed")
}

// Reset eliminates all persisted message data in the store.
func (store *MemoryStore) Reset() {
	defer store.Unlock()
	store.Lock()
	chkcond(store.opened)
	store.messages = make(map[string]*Message)
	store.t.Trace_W(STR, "memorystore wiped")
}

func (store *MemoryStore) SetTracer(tracer *Tracer) {
	defer store.Unlock()
	store.Lock()
	store.t = tracer
}
