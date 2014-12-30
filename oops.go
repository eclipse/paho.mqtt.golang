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
	"os"

	"git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git/packets"
)

/*
 * Connect Errors
 */
var connErrors = map[byte]error{
	packets.Accepted:                     nil,
	packets.RefusedBadProtocolVersion:    errors.New("Unnacceptable protocol version"),
	packets.RefusedIDRejected:            errors.New("Identifier rejected"),
	packets.RefusedServerUnavailable:     errors.New("Server Unavailable"),
	packets.RefusedBadUsernameOrPassword: errors.New("Bad user name or password"),
	packets.RefusedNotAuthorised:         errors.New("Not Authorized"),
	packets.NetworkError:                 errors.New("Network Error"),
	packets.ProtocolViolation:            errors.New("Protocol Violation"),
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

func DefaultErrorHandler(client *Client, reason error) {
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
