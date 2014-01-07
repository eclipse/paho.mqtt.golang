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

var ErrInvalidProtocolVersion = errors.New("Unnacceptable protocol version")
var ErrInvalidClientID = errors.New("Identifier rejected")
var ErrServerUnavailable = errors.New("Server Unavailable")
var ErrBadCredentials = errors.New("Bad user name or password")
var ErrNotAuthorized = errors.New("Not Authorized")
var ErrUnknownReason = errors.New("Unknown RC")
var ErrNotConnected = errors.New("Not Connected")

func DefaultErrorHandler(reason error) {
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
