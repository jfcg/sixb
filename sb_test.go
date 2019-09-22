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

// slice conversions
func Test2(t *testing.T) {
	y := BtI8(buf)
	z := I8tB(y)
	p := BtI4(buf)
	q := I4tB(p)

	if unsafe.Sizeof(buf) != unsafe.Sizeof(Slice{}) ||
		len(y) != 1 || cap(y) != 1 || y[0] != cn2 ||
		len(p) != 2 || cap(p) != 2 || p[0] != cn0 || p[1] != cn1 ||
		len(z) != 8 || cap(z) != 8 || &z[0] != &buf[0] ||
		len(q) != 8 || cap(q) != 8 || &q[0] != &buf[0] {
		t.Fatal("slice conversion error")
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
		b != str[:8] || s != str[:8] {
		t.Fatal("string conversion error")
	}
}

// []String conversions
func Test4(t *testing.T) {
	buf := []byte(big)
	a := BtSs(buf)
	b := I8tSs(BtI8(buf))
	c := I4tSs(BtI4(buf))

	s := 16 / stsz
	r := uint(2-s) << 5
	if len(a) != s || cap(a) != s || uint32(a[0].Data) != cn0 || uint32(a[0].Len>>r) != cn1 ||
		len(b) != s || cap(b) != s || uint32(b[0].Data) != cn0 || uint32(b[0].Len>>r) != cn1 ||
		len(c) != s || cap(c) != s || uint32(c[0].Data) != cn0 || uint32(c[0].Len>>r) != cn1 {
		t.Fatal("String slice conversion error")
	}
}
