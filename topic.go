/*
 * Copyright (c) 2014 IBM Corp.
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
	"strings"
)

// Topic Names and Topic Filters
// The MQTT v3.1.1 spec clarifies a number of ambiguities with regard
// to the validity of Topic strings.
// - A Topic must be between 1 and 65535 bytes.
// - A Topic is case sensitive.
// - A Topic may contain whitespace.
// - A Topic containing a leading forward slash is different than a Topic without.
// - A Topic may be "/" (two levels, both empty string).
// - A Topic must be UTF-8 encoded.
// - A Topic may contain any number of levels.
// - A Topic may contain an empty level (two forward slashes in a row).
// - A TopicName may not contain a wildcard.
// - A TopicFilter may only have a # (multi-level) wildcard as the last level.
// - A TopicFilter may contain any number of + (single-level) wildcards.
// - A TopicFilter with a # will match the absense of a level
//     Example:  a subscription to "foo/#" will match messages published to "foo".

type TopicName struct {
	QoS
	string
}

func NewTopicName(topic string, qos byte) (*TopicName, error) {
	if qos < 0 || qos > 2 {
		return nil, ErrInvalidQoS
	}
	tn := &TopicName{
		QoS(qos),
		topic,
	}
	if e := validateTopicName(topic); e != nil {
		return nil, e
	}
	return tn, nil
}

func validateTopicName(topic string) error {
	if len(topic) == 0 {
		return ErrInvalidTopicNameEmptyString
	}

	levels := strings.Split(topic, "/")
	for _, level := range levels {
		if level == "#" || level == "+" {
			return ErrInvalidTopicNameWildcard
		}
	}
	return nil
}

type TopicFilter struct {
	QoS
	string
}

func NewTopicFilter(topic string, qos byte) (*TopicFilter, error) {
	if qos < 0 || qos > 2 {
		return nil, ErrInvalidQoS
	}
	tf := &TopicFilter{
		QoS(qos),
		topic,
	}
	if e := validateTopicFilter(topic); e != nil {
		return nil, e
	}
	return tf, nil
}

func validateTopicFilter(topic string) error {
	if len(topic) == 0 {
		return ErrInvalidTopicFilterEmptyString
	}

	levels := strings.Split(topic, "/")
	for i, level := range levels {
		if level == "#" && i != len(levels)-1 {
			return ErrInvalidTopicFilterMultilevel
		}
	}
	return nil
}
