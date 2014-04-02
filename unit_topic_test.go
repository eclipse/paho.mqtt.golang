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

func Test_NewTopicName_a_0(t *testing.T) {
	tn, e := NewTopicName("a", 0)
	if e != nil {
		t.Fatalf("error from valid NewTopicName")
	}
	if tn.QoS != QOS_ZERO {
		t.Fatalf("wrong qos from NewTopicName")
	}
	if tn.string != "a" {
		t.Fatalf("wrong name from NewTopicName")
	}
}

func Test_NewTopicName_qos3(t *testing.T) {
	_, e := NewTopicName("a", 3)
	if e != ErrInvalidQoS {
		t.Fatalf("invalid error for invalid qos")
	}
}

func Test_NewTopicName_P_0(t *testing.T) {
	_, e := NewTopicName("+", 0)
	if e != ErrInvalidTopicNameWildcard {
		t.Fatalf("invalid error for topic name with wildcard")
	}
}

func Test_NewTopicName_H_0(t *testing.T) {
	_, e := NewTopicName("#", 0)
	if e != ErrInvalidTopicNameWildcard {
		t.Fatalf("invalid error for topic name with wildcard")
	}
}

func Test_NewTopicName_ES(t *testing.T) {
	_, e := NewTopicName("", 0)
	if e != ErrInvalidTopicNameEmptyString {
		t.Fatalf("invalid error for empty topic name")
	}
}

func Test_NewTopicFilter_a_0(t *testing.T) {
	tf, e := NewTopicFilter("a", 0)
	if e != nil {
		t.Fatalf("error from valid NewTopicFilter")
	}
	if tf.QoS != QOS_ZERO {
		t.Fatalf("wrong qos from NewTopicFilter")
	}
	if tf.string != "a" {
		t.Fatalf("wrong filter from NewTopicFilter")
	}
}

func Test_NewTopicFilter_ES(t *testing.T) {
	_, e := NewTopicFilter("", 0)
	if e != ErrInvalidTopicFilterEmptyString {
		t.Fatalf("invalid error for empty topic filter")
	}
}

func Test_NewTopicFilter_H(t *testing.T) {
	_, e := NewTopicFilter("a/#/c", 0)
	if e != ErrInvalidTopicFilterMultilevel {
		t.Fatalf("invalid error for bad multilevel topic filter")
	}
}

func Test_NewTopicFilter_qos3(t *testing.T) {
	_, e := NewTopicFilter("a", 3)
	if e != ErrInvalidQoS {
		t.Fatalf("invalid error for invalid qos")
	}
}
