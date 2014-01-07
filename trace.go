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
	"os"
	"strings"
	"time"
)

type tracelevel byte

type Tracer struct {
	level    tracelevel
	output   *os.File
	clientid string
}

const (
	Off      tracelevel = 0
	Critical tracelevel = 10
	Warn     tracelevel = 20
	Verbose  tracelevel = 30
	// [comp] [timestamp] [clientid] message
	frmt = "%s %s [%s] %s\n"
)

func timestamp() string {
	raw := fmt.Sprintf("%v", time.Now())
	tks := strings.Fields(raw)
	tme := tks[1]
	max := 14
	nzeroes := max - len(tme)
	for i := 0; i < nzeroes; i++ {
		tme += "0"
	}
	tme = tme[0:max]
	return tme
}

func (t *Tracer) Trace_V(cm component, f string, v ...interface{}) {
	if t.level >= Verbose && t.output != nil {
		x := fmt.Sprintf(f, v...)
		m := fmt.Sprintf(frmt, cm, timestamp(), t.clientid, x)
		t.output.WriteString(m)
	}
}

func (c *MqttClient) trace_v(cm component, f string, v ...interface{}) {
	c.t.Trace_V(cm, f, v...)
}

func (t *Tracer) Trace_W(cm component, f string, v ...interface{}) {
	if t.level >= Warn && t.output != nil {
		x := fmt.Sprintf(f, v...)
		m := fmt.Sprintf(frmt, cm, timestamp(), t.clientid, x)
		t.output.WriteString(m)
	}
}

func (c *MqttClient) trace_w(cm component, f string, v ...interface{}) {
	c.t.Trace_W(cm, f, v...)
}

func (t *Tracer) Trace_E(cm component, f string, v ...interface{}) {
	if t.level >= Critical && t.output != nil {
		x := fmt.Sprintf(f, v...)
		m := fmt.Sprintf(frmt, cm, timestamp(), t.clientid, x)
		t.output.WriteString(m)
	}
}

func (c *MqttClient) trace_e(cm component, f string, v ...interface{}) {
	c.t.Trace_E(cm, f, v...)
}
