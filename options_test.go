/*
 * Copyright (c) 2021 IBM Corp and others.
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
 *    Måns Ansgariusson
 */

// Portions copyright © 2018 TIBCO Software Inc.
package mqtt

import (
	"fmt"
	"net"
	"net/url"
	"testing"
)

func TestSetCustomConnectionOptions(t *testing.T) {
	var customConnectionFunc OpenConnectionFunc = func(uri *url.URL, options ClientOptions) (net.Conn, error) {
		return nil, fmt.Errorf("not implemented open connection func")
	}
	options := &ClientOptions{}
	options = options.SetCustomOpenConnectionFn(customConnectionFunc)
	if options.CustomOpenConnectionFn == nil {
		t.Error("custom open connection function cannot be set")
	}
}
