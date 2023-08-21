//	Copyright (c) 2023-present, Serhat Şevki Dinçer.
//	This Source Code Form is subject to the terms of the Mozilla Public
//	License, v. 2.0. If a copy of the MPL was not distributed with this
//	file, You can obtain one at http://mozilla.org/MPL/2.0/.

//go:build 386 || amd64

package sixb

import "testing"

func TestInc(t *testing.T) {
	var counter uint32

	for max := ^uint32(0) - 10; max != 10; max++ {
		if Inc(nil, max) != 0 {
			t.Fatal("unexpected return")
		}
		for c := ^uint32(0) - 20; c != 20; c++ {
			counter = c
			exp := c
			ret := uint32(0)
			if c < max {
				exp++
				ret = exp
			}
			if Inc(&counter, max) != ret {
				t.Fatal("unexpected return")
			}
			if counter != exp {
				t.Fatal("unexpected value")
			}
		}
	}
}

func TestDec(t *testing.T) {
	var counter uint32

	if ^Dec(nil) != 0 {
		t.Fatal("unexpected return")
	}
	for c := ^uint32(0) - 20; c != 20; c++ {
		counter = c
		exp := c
		ret := ^uint32(0)
		if c != 0 {
			exp--
			ret = exp
		}
		if Dec(&counter) != ret {
			t.Fatal("unexpected return")
		}
		if counter != exp {
			t.Fatal("unexpected value")
		}
	}
}
