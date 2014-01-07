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

import "testing"

func same(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func checkOrder(exp, res []string, t *testing.T) {
	if !same(exp, res) {
		t.Fatalf("failed\nexp: %v\nres: %v\n", exp, res)
	}
}

func Test_reorder_simple(t *testing.T) {
	keys := []string{"i.4", "i.9", "i.3", "i.6"}
	exp := []string{"i.3", "i.4", "i.6", "i.9"}
	res := reorder(keys)
	checkOrder(exp, res, t)
}

func Test_reorder_typical(t *testing.T) {
	keys := []string{"o.7", "o.3", "o.6", "o.12100", "o.12300", "o.12200"}
	exp := []string{"o.3", "o.6", "o.7", "o.12100", "o.12200", "o.12300"}
	res := reorder(keys)
	checkOrder(exp, res, t)
}

func Test_reorder_looped_gap(t *testing.T) {
	keys := []string{"i.123", "i.110", "i.120", "i.65534", "i.65531", "i.65532"}
	exp := []string{"i.65531", "i.65532", "i.65534", "i.110", "i.120", "i.123"}
	res := reorder(keys)
	checkOrder(exp, res, t)
}

func Test_reorder_expose_ascii_sort_problem(t *testing.T) {
	keys := []string{"i.1", "i.112", "i.12"}
	exp := []string{"i.1", "i.12", "i.112"}
	res := reorder(keys)
	checkOrder(exp, res, t)
}
