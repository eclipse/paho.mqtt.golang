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
)

/*
 * Connect Errors
 */
var ErrInvalidProtocolVersion = errors.New("Unnacceptable protocol version")
var ErrInvalidClientID = errors.New("Identifier rejected")
var ErrServerUnavailable = errors.New("Server Unavailable")
var ErrBadCredentials = errors.New("Bad user name or password")
var ErrNotAuthorized = errors.New("Not Authorized")
var ErrUnknownReason = errors.New("Unknown RC")
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

func chkrc(rc ConnRC) error {
	if rc != CONN_ACCEPTED {
		switch rc {
		case CONN_REF_BAD_PROTO_VER:
			return ErrInvalidProtocolVersion
		case CONN_REF_ID_REJ:
			return ErrInvalidClientID
		case CONN_REF_SERV_UNAVAIL:
			return ErrServerUnavailable
		case CONN_REF_BAD_USER_PASS:
			return ErrBadCredentials
		case CONN_REF_NOT_AUTH:
			return ErrNotAuthorized
		default:
			return ErrUnknownReason
		}
	}
	return nil
}

func rc2str(rc ConnRC) string {
	switch rc {
	case CONN_ACCEPTED:
		return "CONN_ACCEPTED"
	case CONN_REF_BAD_PROTO_VER:
		return "CONN_REF_BAD_PROTO_VER"
	case CONN_REF_ID_REJ:
		return "CONN_REF_ID_REJ"
	case CONN_REF_SERV_UNAVAIL:
		return "CONN_REF_SERV_UNAVAIL"
	case CONN_REF_BAD_USER_PASS:
		return "CONN_REF_BAD_USER_PASS"
	case CONN_REF_NOT_AUTH:
		return "CONN_REF_NOT_AUTH"
	default:
		return "UNKNOWN"
	}
}
