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
	"errors"
	"fmt"
	. "github.com/alsm/hrotti/packets"
	"os"
)

/*
 * Connect Errors
 */
var connErrors = map[byte]error{
	CONN_ACCEPTED:           nil,
	CONN_REF_BAD_PROTO_VER:  errors.New("Unnacceptable protocol version"),
	CONN_REF_ID_REJ:         errors.New("Identifier rejected"),
	CONN_REF_SERV_UNAVAIL:   errors.New("Server Unavailable"),
	CONN_REF_BAD_USER_PASS:  errors.New("Bad user name or password"),
	CONN_REF_NOT_AUTH:       errors.New("Not Authorized"),
	CONN_NETWORK_ERROR:      errors.New("Network Error"),
	CONN_PROTOCOL_VIOLATION: errors.New("Protocol Violation"),
}

var ErrNotConnected = errors.New("Not Connected")

/*
 * Topic Errors
 */
var ErrInvalidTopicNameEmptyString = errors.New("Invalid TopicName - may not be empty string")
var ErrInvalidTopicNameWildcard = errors.New("Invalid TopicName - may not contain wild card")
var ErrInvalidTopicFilterEmptyString = errors.New("Invalid TopicFilter - may not be empty string")
var ErrInvalidTopicFilterMultilevel = errors.New("Invalid TopicFilter - multi-level wildcard must be last level")

/*
 * QoS Errors
 */
var ErrInvalidQoS = errors.New("Invalid QoS")

func DefaultErrorHandler(client *MqttClient, reason error) {
	fmt.Fprintf(os.Stderr, "%s go-mqtt suffered fatal error %v", ERR, reason)
	os.Exit(1)
}

func chkerr(e error) {
	if e != nil {
		panic(e)
	}
}

func chkcond(b bool) {
	if !b {
		panic("oops")
	}
}
