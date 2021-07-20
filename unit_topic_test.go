/*
 * Copyright (c) 2013 IBM Corp.
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
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package mqtt

import (
	"testing"
)

func Test_ValidateTopicAndQos_qos3(t *testing.T) {
	e := validateTopicAndQos("a", 3)
	if e != ErrInvalidQos {
		t.Fatalf("invalid error for invalid qos")
	}
}

func Test_ValidateTopicAndQos_ES(t *testing.T) {
	e := validateTopicAndQos("", 0)
	if e != ErrInvalidTopicEmptyString {
		t.Fatalf("invalid error for empty topic name")
	}
}

func Test_ValidateTopicAndQos_a_0(t *testing.T) {
	e := validateTopicAndQos("a", 0)
	if e != nil {
		t.Fatalf("error from valid NewTopicFilter")
	}
}

func Test_ValidateTopicAndQos_H(t *testing.T) {
	e := validateTopicAndQos("a/#/c", 0)
	if e != ErrInvalidTopicMultilevel {
		t.Fatalf("invalid error for bad multilevel topic filter")
	}
}
