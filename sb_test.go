/*	Copyright (c) 2019-present, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sixb

import (
	"testing"
	"unsafe"
)

// AnumToSixb & SixbToAnum bijection & domain
func TestSixb(t *testing.T) {
	if len(anumToSixb) != 256 || len(sixbToAnum) != 256 {
		t.Fatal("invalid lengths")
	}

	for i := 255; i >= 0; i-- {
		c := byte(i)
		d := AnumToSixb(c)
		if c == d {
			t.Fatal("fixed point", i)
		}
		if c != SixbToAnum(d) {
			t.Fatal("inverse does not work", i)
		}
	}

	n := 0 // cycle length
	for d := AnumToSixb(0); d != 0; n++ {
		d = AnumToSixb(d)
	}
	if n != 255 {
		t.Fatal("multiple cycles")
	}

	l := "0:@Zaz"
	for i := 4; i >= 0; i -= 2 {
		for c := l[i]; c <= l[i+1]; c++ {
			if AnumToSixb(c) > 63 {
				t.Fatal("domain error", c)
			}
		}
	}
}

var (
	big = "qwert123qwert123qwert123"
	sml = big[:9]
	buf = []byte(sml)
)

const (
	cn0 = 1919252337 // "qwer"
	cn1 = 858927476  // "t123"
	cn2 = cn0 + cn1<<32
)

func TestCopy(t *testing.T) {
	if !InsideTest() {
		t.Fatal("InsideTest does not work")
	}
	b := Copy(buf)
	if &b[0] == &buf[0] || String(b) != String(buf) {
		t.Fatal("Copy does not work")
	}
}

// bad byte slice?
func badb(q []byte) bool {
	return len(q) != 8 || cap(q) != 8 || &q[0] != &buf[0]
}

// slice conversions
func TestSlice(t *testing.T) {
	y := Slice[uint64](buf)
	z := Slice[byte](y)
	p := Slice[uint32](buf)
	q := Slice[byte](p)

	if len(y) != 1 || cap(y) != 1 || y[0] != cn2 ||
		len(p) != 2 || cap(p) != 2 || p[0] != cn0 || p[1] != cn1 ||
		!SamePtr(&y[0], &p[0]) ||
		!SamePtr(&y[0], &buf[0]) ||
		badb(z) || badb(q) {
		t.Fatal("slice conversion error")
	}
}

// slice conversions
func TestSlice2(t *testing.T) {
	p := Slice[uint32](buf)
	y := Slice[uint64](p)
	z := Slice[uint32](y)
	u := Slice[uint16](z)

	if len(y) != 1 || cap(y) != 1 || y[0] != cn2 ||
		len(p) != 2 || cap(p) != 2 || p[0] != cn0 || p[1] != cn1 ||
		len(z) != 2 || cap(z) != 2 || z[0] != cn0 || z[1] != cn1 ||
		!SamePtr(&y[0], &p[0]) || !SamePtr(&y[0], &z[0]) ||
		!SamePtr(&y[0], &buf[0]) || !SamePtr(&y[0], &u[0]) {
		t.Fatal("slice conversion error")
	}
}

// nil string/slice conversions
func TestSlice3(t *testing.T) {
	var s string
	var a []byte

	if Bytes(s) != nil || Integers[uint32](s) != nil ||
		Integers[uint64](s) != nil || Slice[uint32](a) != nil ||
		Slice[uint64](a) != nil || Slice[string](a) != nil || Slice[[]byte](a) != nil {
		t.Fatal("nil string/slice conversion error")
	}
}

func TestSlice3b(t *testing.T) {
	var b []uint32
	var c []uint64

	if String(b) != "" || String(c) != "" || Slice[byte](b) != nil ||
		Slice[byte](c) != nil || Slice[string](b) != nil || Slice[string](c) != nil ||
		Slice[[]byte](b) != nil || Slice[[]byte](c) != nil {
		t.Fatal("nil string/slice conversion error")
	}
}

// string conversions
func TestString(t *testing.T) {
	a := Integers[uint64](sml)
	b := String(a)
	r := Integers[uint32](sml)
	s := String(r)

	if len(a) != 1 || cap(a) != 1 || a[0] != cn2 ||
		len(r) != 2 || cap(r) != 2 || r[0] != cn0 || r[1] != cn1 ||
		SamePtr(&a[0], &buf[0]) || !SamePtr(&a[0], &r[0]) ||
		b != sml[:8] || s != sml[:8] {
		t.Fatal("string conversion error")
	}
}

// string conversions
func TestString2(t *testing.T) {
	a := Bytes(big)
	b := Bytes(sml)
	c := String(buf)

	if len(a) != len(big) || len(b) != len(sml) ||
		len(a) != cap(a) || len(b) != cap(b) ||
		len(c) != len(buf) || c != sml ||
		&a[0] != &b[0] || &a[0] == &buf[0] {
		t.Fatal("string conversion error")
	}
}

// != "qwert123"
func badx(x uint) bool {
	return uint32(x) != cn0 || uint32(x>>32) != cn1
}

// bad string slice?
func badStr(a []string) bool {
	if len(a) == 0 || len(a) != cap(a) {
		return true
	}
	d, l := PtrToInt(unsafe.StringData(a[0])), len(a[0])
	if unsafe.Sizeof("") == 8 {
		return len(a) != 3 || d != cn0 || l != cn1
	}
	return len(a) != 1 || badx(d) || badx(uint(l))
}

// []string conversions
func TestSlice4(t *testing.T) {
	buf := []byte(big)
	a := Slice[string](buf)
	b := Slice[string](Slice[uint64](buf))
	c := Slice[string](Slice[uint32](buf))

	if badStr(a) || badStr(b) || badStr(c) {
		t.Fatal("string slice conversion error")
	}
}

// bad []byte slice?
func badSlc(a [][]byte) bool {
	if len(a) == 0 || len(a) != cap(a) {
		return true
	}
	d, l, c := PtrToInt(unsafe.SliceData(a[0])), len(a[0]), cap(a[0])
	if unsafe.Sizeof([]byte{}) == 12 {
		return len(a) != 2 || d != cn0 || l != cn1 || c != cn0
	}
	return len(a) != 1 || badx(d) || badx(uint(l)) || badx(uint(c))
}

// [][]byte conversions
func TestSlice5(t *testing.T) {
	buf := []byte(big)
	a := Slice[[]byte](buf)
	b := Slice[[]byte](Slice[uint64](buf))
	c := Slice[[]byte](Slice[uint32](buf))

	if badSlc(a) || badSlc(b) || badSlc(c) {
		t.Fatal("Slice slice conversion error")
	}
}

func TestPtrToInt(t *testing.T) {
	var p *int
	if PtrToInt(p) != 0 {
		t.Fatal("Nil pointer must convert to to zero")
	}
	var arr [2]int
	v1 := PtrToInt(&arr[0])
	v2 := PtrToInt(&arr[1])
	if PtrToInt(t) == 0 || v1 == 0 || v2 == 0 ||
		v1+uint(unsafe.Sizeof(arr[0])) != v2 {
		t.Fatal("Pointer to unsigned integer conversion error")
	}
}
