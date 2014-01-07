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
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

const (
	_MSG_EXT = ".msg"
	_BKP_EXT = ".bkp"
)

// FileStore implements the store interface using the filesystem to provide
// true persistence, even across client failure. This is designed to use a
// single directory per running client. If you are running multiple clients
// on the same filesystem, you will need to be careful to specify unique
// store directories for each.
type FileStore struct {
	sync.RWMutex
	directory string
	opened    bool
	t         *Tracer
}

// NewFileStore will create a new FileStore which stores its messages in the
// directory provided.
func NewFileStore(directory string) *FileStore {
	store := &FileStore{
		directory: directory,
		opened:    false,
		t:         nil,
	}
	return store
}

// Open will allow the FileStore to be used.
func (store *FileStore) Open() {
	defer store.Unlock()
	store.Lock()
	// if no store directory was specified in ClientOpts, by default use the
	// current working directory
	if store.directory == "" {
		store.directory, _ = os.Getwd()
	}

	// if store dir exists, great, otherwise, create it
	if !exists(store.directory) {
		perms := os.FileMode(0770)
		merr := os.MkdirAll(store.directory, perms)
		chkerr(merr)
	}
	store.opened = true
	store.t.Trace_V(STR, "store is opened at %s", store.directory)
}

// Close will disallow the FileStore from being used.
func (store *FileStore) Close() {
	defer store.Unlock()
	store.Lock()
	store.opened = false
	store.t.Trace_V(STR, "store is not open")
}

// Put will put a message into the store, associated with the provided
// key value.
func (store *FileStore) Put(key string, m *Message) {
	defer store.Unlock()
	store.Lock()
	chkcond(store.opened)
	full := fullpath(store.directory, key)
	if exists(full) {
		backup(store.directory, key) // make a copy of what already exists
		defer unbackup(store.directory, key)
	}
	write(store.directory, key, m)
	chkcond(exists(full))
}

// Get will retrieve a message from the store, the one associated with
// the provided key value.
func (store *FileStore) Get(key string) (m *Message) {
	defer store.RUnlock()
	store.RLock()
	chkcond(store.opened)
	filepath := fullpath(store.directory, key)
	if !exists(filepath) {
		return nil
	}
	mfile, oerr := os.Open(filepath)
	chkerr(oerr)
	all, rerr := ioutil.ReadAll(mfile)
	chkerr(rerr)
	msg := decode(all)
	cerr := mfile.Close()
	chkerr(cerr)
	return msg
}

// All will provide a list of all of the keys associated with messages
// currenly residing in the FileStore.
func (store *FileStore) All() []string {
	defer store.RUnlock()
	store.RLock()
	return store.all()
}

// Del will remove the persisted message associated with the provided
// key from the FileStore.
func (store *FileStore) Del(key string) {
	defer store.Unlock()
	store.Lock()
	store.del(key)
}

// Reset will remove all persisted messages from the FileStore.
func (store *FileStore) Reset() {
	defer store.Unlock()
	store.Lock()
	store.t.Trace_W(STR, "FileStore Reset")
	for _, key := range store.all() {
		store.del(key)
	}
}

// lockless
func (store *FileStore) all() []string {
	chkcond(store.opened)
	keys := []string{}
	files, rderr := ioutil.ReadDir(store.directory)
	chkerr(rderr)
	for _, f := range files {
		store.t.Trace_V(STR, "file in All(): %s", f.Name())
		key := f.Name()[0 : len(f.Name())-4] // remove file extension
		keys = append(keys, key)
	}
	return keys
}

// lockless
func (store *FileStore) del(key string) {
	chkcond(store.opened)
	store.t.Trace_V(STR, "store del filepath: %s", store.directory)
	store.t.Trace_V(STR, "store delete key: %v", key)
	filepath := fullpath(store.directory, key)
	store.t.Trace_V(STR, "path of deletion: `%s`", filepath)
	if !exists(filepath) {
		store.t.Trace_E(STR, "store could not delete key: %v", key)
		return
	}
	rerr := os.Remove(filepath)
	chkerr(rerr)
	store.t.Trace_V(STR, "del msg: %v", key)
	chkcond(!exists(filepath))
}

func (store *FileStore) SetTracer(trace *Tracer) {
	store.t = trace
}

func fullpath(store string, key string) string {
	p := path.Join(store, key+_MSG_EXT)
	return p
}

func bkppath(store string, key string) string {
	p := path.Join(store, key+_BKP_EXT)
	return p
}

// create file called "X.[messageid].msg" located in the store
// the contents of the file is the bytes of the message
// if a message with m's message id already exists, it will
// be overwritten
// X will be 'i' for inbound messages, and O for outbound messages
func write(store, key string, m *Message) {
	filepath := fullpath(store, key)
	f, err := os.Create(filepath)
	chkerr(err)
	_, werr := f.Write(m.Bytes())
	chkerr(werr)
	cerr := f.Close()
	chkerr(cerr)
}

func exists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		chkerr(err)
	}
	return true
}

func backup(store, key string) {
	bkpp := bkppath(store, key)
	fulp := fullpath(store, key)
	backup, err := os.Create(bkpp)
	chkerr(err)
	mfile, oerr := os.Open(fulp)
	chkerr(oerr)
	_, cerr := io.Copy(backup, mfile)
	chkerr(cerr)
	clberr := backup.Close()
	chkerr(clberr)
	clmerr := mfile.Close()
	chkerr(clmerr)
}

// Identify .bkp files in the store and turn them into .msg files,
// whether or not it overwrites an existing file. This is safe because
// I'm copying the Paho Java client and they say it is.
func restore(store string) {
	files, rderr := ioutil.ReadDir(store)
	chkerr(rderr)
	for _, f := range files {
		fname := f.Name()
		if len(fname) > 4 {
			if fname[len(fname)-4:] == _BKP_EXT {
				key := fname[0 : len(fname)-4]
				fulp := fullpath(store, key)
				msg, cerr := os.Create(fulp)
				chkerr(cerr)
				bkpp := path.Join(store, fname)
				bkp, oerr := os.Open(bkpp)
				chkerr(oerr)
				n, cerr := io.Copy(msg, bkp)
				chkerr(cerr)
				chkcond(n > 0)
				clmerr := msg.Close()
				chkerr(clmerr)
				clberr := bkp.Close()
				chkerr(clberr)
				remerr := os.Remove(bkpp)
				chkerr(remerr)
			}
		}
	}
}

func unbackup(store, key string) {
	bkpp := bkppath(store, key)
	remerr := os.Remove(bkpp)
	chkerr(remerr)
}
