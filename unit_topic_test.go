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
	"testing"
)

func Test_ValidateTopicAndQos_qos3(t *testing.T) {
	e := validateTopicAndQos("a", 3)
	if e != ErrInvalidQoS {
		t.Fatalf("invalid error for invalid qos")
	}
}

// func Test_ValidateTopicAndQos_P_0(t *testing.T) {
// 	e := validateTopicAndQos("+", 0)
// 	if e != ErrInvalidTopicNameWildcard {
// 		t.Fatalf("invalid error for topic name with wildcard")
// 	}
// }

// func Test_ValidateTopicAndQos_H_0(t *testing.T) {
// 	e := validateTopicAndQos("#", 0)
// 	if e != ErrInvalidTopicNameWildcard {
// 		t.Fatalf("invalid error for topic name with wildcard")
// 	}
// }

func Test_ValidateTopicAndQos_ES(t *testing.T) {
	e := validateTopicAndQos("", 0)
	if e != ErrInvalidTopicNameEmptyString {
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
	if e != ErrInvalidTopicFilterMultilevel {
		t.Fatalf("invalid error for bad multilevel topic filter")
	}
}
