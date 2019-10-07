/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sixb

import (
	"github.com/jfcg/sorty"
	"testing"
	"unsafe"
)

// Txt2int collision test
func Test0(t *testing.T) {
	const N = 5 << 26
	hl := make([]uint64, N)     // hash list: 2.5 GiB ram
	bf := [...]byte{0, 0, 0, 0} // input buffer

	// fill hl with hashes of short utf8 text
	// hl[N-1] = Txt2int("") = 0
	for i, l := N-2, 0; i >= 0; i-- {

		// next utf8-ish input
		for k := 0; ; k++ {
			if bf[k] == 0 {
				l++ // increase input length
			}
			bf[k]++

			if bf[k] != 0 {
				break
			}
			bf[k]++ // skip zero digit, continue with carry
		}
		hl[i] = Txt2int(string(bf[:l]))
	}

	sorty.SortU8(hl)

	k := 0 // count collisions
	for i := N - 1; i > 0; i-- {
		if hl[i] == hl[i-1] {
			k++
		}
	}
	if k > 0 {
		t.Fatal("Txt2int has at least", k, "collisions for short utf8 inputs")
	}
}

// An2sb & Sb2an bijection & domain
func Test1(t *testing.T) {
	if len(An2sb) != 256 || len(Sb2an) != 256 {
		t.Fatal("invalid lengths")
	}

	for i := 255; i >= 0; i-- {
		c := byte(i)
		d := An2sb[c]
		if c == d {
			t.Fatal("fixed point", i)
		}
		if c != Sb2an[d] {
			t.Fatal("inverse does not work", i)
		}
	}

	l := "0:@Zaz"
	for i := 4; i >= 0; i -= 2 {
		for c := l[i]; c <= l[i+1]; c++ {
			if An2sb[c] > 63 {
				t.Fatal("domain error", c)
			}
		}
	}
}

var (
	big = "qwert123qwert12345"
	str = big[:9]
	buf = []byte(str)
)

const (
	cn0 = 1919252337
	cn1 = 858927476
	cn2 = cn0 + cn1<<32
)

// bad byte slice?
func badb(q []byte) bool {
	return len(q) != 8 || cap(q) != 8 || &q[0] != &buf[0]
}

// slice conversions
func Test2(t *testing.T) {
	y := BtI8(buf)
	z := I8tB(y)
	p := BtI4(buf)
	q := I4tB(p)

	if unsafe.Sizeof(buf) != unsafe.Sizeof(Slice{}) ||
		len(y) != 1 || cap(y) != 1 || y[0] != cn2 ||
		len(p) != 2 || cap(p) != 2 || p[0] != cn0 || p[1] != cn1 ||
		unsafe.Pointer(&y[0]) != unsafe.Pointer(&p[0]) ||
		unsafe.Pointer(&y[0]) != unsafe.Pointer(&buf[0]) ||
		badb(z) || badb(q) {
		t.Fatal("slice conversion error")
	}
}

// nil string/slice conversions
func Test2b(t *testing.T) {
	var (
		s string
		a []byte
		b []uint32
		c []uint64
	)

	if StB(s) != nil || StI4(s) != nil || StI8(s) != nil ||
		I4tS(b) != "" || I8tS(c) != "" ||
		BtI4(a) != nil || BtI8(a) != nil || I4tB(b) != nil || I8tB(c) != nil ||
		BtSs(a) != nil || I4tSs(b) != nil || I8tSs(c) != nil {
		t.Fatal("nil string/slice conversion error")
	}
}

// string conversions
func Test3(t *testing.T) {
	a := StI8(str)
	b := I8tS(a)
	r := StI4(str)
	s := I4tS(r)

	if stsz != 8 && stsz != 16 || stsz != int(unsafe.Sizeof(String{})) ||
		len(a) != 1 || cap(a) != 1 || a[0] != cn2 ||
		len(r) != 2 || cap(r) != 2 || r[0] != cn0 || r[1] != cn1 ||
		unsafe.Pointer(&a[0]) == unsafe.Pointer(&buf[0]) ||
		unsafe.Pointer(&a[0]) != unsafe.Pointer(&r[0]) || b != str[:8] || s != str[:8] {
		t.Fatal("string conversion error")
	}
}

// string conversions
func Test3b(t *testing.T) {
	a := StB(big)
	b := StB(str)

	if len(a) != len(big) || len(b) != len(str) ||
		len(a) != cap(a) || len(b) != cap(b) ||
		&a[0] != &b[0] || &a[0] == &buf[0] {
		t.Fatal("string conversion error")
	}
}

// bad slice?
func bad(a []String) bool {
	const (
		s = 16 / stsz
		r = uint(2-s) << 5
	)
	return len(a) != s || cap(a) != s ||
		uint32(uintptr(a[0].Data)) != cn0 || uint32(a[0].Len>>r) != cn1
}

// []String conversions
func Test4(t *testing.T) {
	buf := []byte(big)
	a := BtSs(buf)
	b := I8tSs(BtI8(buf))
	c := I4tSs(BtI4(buf))

	if bad(a) || bad(b) || bad(c) {
		t.Fatal("String slice conversion error")
	}
}
