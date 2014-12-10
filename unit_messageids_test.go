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
	"fmt"
	"testing"
)

func Test_getId(t *testing.T) {
	mids := &messageIds{index: make(map[MId]bool)}
	mids.generateMsgIds()

	i1, err := mids.getId()
	if err != nil {
		t.Fatalf("Failed to get id. %s", err.Error())
	}

	if i1 != MId(1) {
		t.Fatalf("i1 was wrong: %v", i1)
	}

	i2, err := mids.getId()
	if err != nil {
		t.Fatalf("Failed to get id. %s", err.Error())
	}

	if i2 != MId(2) {
		t.Fatalf("i2 was wrong: %v", i2)
	}

	for i := 3; i < 100; i++ {
		id, err := mids.getId()
		if err != nil {
			t.Fatalf("Failed to get id. %s", err.Error())
		}
		if id != MId(i) {
			t.Fatalf("id was wrong expected %v got %v", i, id)
		}
	}
}

func Test_freeId(t *testing.T) {
	mids := &messageIds{index: make(map[MId]bool)}
	mids.generateMsgIds()

	i1, err := mids.getId()
	if err != nil {
		t.Fatalf("Failed to get id. %s", err.Error())
	}
	mids.freeId(i1)

	if i1 != MId(1) {
		t.Fatalf("i1 was wrong: %v", i1)
	}

	i2, err := mids.getId()
	if err != nil {
		t.Fatalf("Failed to get id. %s", err.Error())
	}
	fmt.Printf("i2: %v\n", i2)
}

func Test_messageids_mix(t *testing.T) {
	mids := &messageIds{index: make(map[MId]bool)}
	mids.generateMsgIds()

	done := make(chan bool)
	a := make(chan MId, 3)
	b := make(chan MId, 20)
	c := make(chan MId, 100)

	go func() {
		for i := 0; i < 10000; i++ {
			id, err := mids.getId()
			if err != nil {
				t.Fatalf("Failed to get id. %s", err.Error())
			}
			a <- id

			mids.freeId(<-b)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			id, err := mids.getId()
			if err != nil {
				t.Fatalf("Failed to get id. %s", err.Error())
			}
			b <- id

			mids.freeId(<-c)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			id, err := mids.getId()
			if err != nil {
				t.Fatalf("Failed to get id. %s", err.Error())
			}
			c <- id

			mids.freeId(<-a)
		}
		done <- true
	}()

	<-done
	<-done
	<-done
}
