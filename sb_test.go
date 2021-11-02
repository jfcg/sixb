/*	Copyright (c) 2019, Serhat Şevki Dinçer.
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
func Test1(t *testing.T) {
	if len(AnumToSixb) != 256 || len(SixbToAnum) != 256 {
		t.Fatal("invalid lengths")
	}

	for i := 255; i >= 0; i-- {
		c := byte(i)
		d := AnumToSixb[c]
		if c == d {
			t.Fatal("fixed point", i)
		}
		if c != SixbToAnum[d] {
			t.Fatal("inverse does not work", i)
		}
	}

	n := 0 // cycle length
	for d := AnumToSixb[0]; d != 0; n++ {
		d = AnumToSixb[d]
	}
	if n != 255 {
		t.Fatal("multiple cycles")
	}

	l := "0:@Zaz"
	for i := 4; i >= 0; i -= 2 {
		for c := l[i]; c <= l[i+1]; c++ {
			if AnumToSixb[c] > 63 {
				t.Fatal("domain error", c)
			}
		}
	}
}

var (
	big = "qwert123qwert123qwert123"
	str = big[:9]
	buf = []byte(str)
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
	if &b[0] == &buf[0] || string(b) != string(buf) {
		t.Fatal("Copy does not work")
	}
}

// bad byte slice?
func badb(q []byte) bool {
	return len(q) != 8 || cap(q) != 8 || &q[0] != &buf[0]
}

// slice conversions
func Test2(t *testing.T) {
	y := BtoU8(buf)
	z := U8toB(y)
	p := BtoU4(buf)
	q := U4toB(p)

	if unsafe.Sizeof(buf) != unsafe.Sizeof(Slice{}) ||
		len(y) != 1 || cap(y) != 1 || y[0] != cn2 ||
		len(p) != 2 || cap(p) != 2 || p[0] != cn0 || p[1] != cn1 ||
		unsafe.Pointer(&y[0]) != unsafe.Pointer(&p[0]) ||
		unsafe.Pointer(&y[0]) != unsafe.Pointer(&buf[0]) ||
		badb(z) || badb(q) {
		t.Fatal("slice conversion error")
	}
}

// slice conversions
func Test2a(t *testing.T) {
	p := BtoU4(buf)
	y := U4toU8(p)
	z := U8toU4(y)
	a := unsafe.Pointer(&y[0])

	if len(y) != 1 || cap(y) != 1 || y[0] != cn2 ||
		len(p) != 2 || cap(p) != 2 || p[0] != cn0 || p[1] != cn1 ||
		len(z) != 2 || cap(z) != 2 || z[0] != cn0 || z[1] != cn1 ||
		a != unsafe.Pointer(&p[0]) || a != unsafe.Pointer(&z[0]) ||
		a != unsafe.Pointer(&buf[0]) {
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

	if StoB(s) != nil || StoU4(s) != nil || StoU8(s) != nil ||
		U4toS(b) != "" || U8toS(c) != "" ||
		BtoU4(a) != nil || BtoU8(a) != nil || U4toB(b) != nil || U8toB(c) != nil ||
		BtoStrs(a) != nil || U4toStrs(b) != nil || U8toStrs(c) != nil ||
		BtoSlcs(a) != nil || U4toSlcs(b) != nil || U8toSlcs(c) != nil {
		t.Fatal("nil string/slice conversion error")
	}
}

// string conversions
func Test3(t *testing.T) {
	a := StoU8(str)
	b := U8toS(a)
	r := StoU4(str)
	s := U4toS(r)

	if SliceSize != int(unsafe.Sizeof(Slice{})) ||
		StrSize != int(unsafe.Sizeof(String{})) ||
		len(a) != 1 || cap(a) != 1 || a[0] != cn2 ||
		len(r) != 2 || cap(r) != 2 || r[0] != cn0 || r[1] != cn1 ||
		unsafe.Pointer(&a[0]) == unsafe.Pointer(&buf[0]) ||
		unsafe.Pointer(&a[0]) != unsafe.Pointer(&r[0]) || b != str[:8] || s != str[:8] {
		t.Fatal("string conversion error")
	}
}

// string conversions
func Test3b(t *testing.T) {
	a := StoB(big)
	b := StoB(str)
	c := BtoS(buf)

	if len(a) != len(big) || len(b) != len(str) ||
		len(a) != cap(a) || len(b) != cap(b) ||
		len(c) != len(buf) || c != str ||
		&a[0] != &b[0] || &a[0] == &buf[0] {
		t.Fatal("string conversion error")
	}
}

// != "qwert123"
func badx(x int) bool {
	return uint32(x) != cn0 || uint32(x>>32) != cn1
}

// bad String slice?
func badStr(a []String) bool {
	if len(a) != cap(a) {
		return true
	}
	d := uintptr(a[0].Data)
	if StrSize == 8 {
		return len(a) != 3 || d != cn0 || a[0].Len != cn1
	}
	return len(a) != 1 || badx(int(d)) || badx(a[0].Len)
}

// []String conversions
func Test4(t *testing.T) {
	buf := []byte(big)
	a := BtoStrs(buf)
	b := U8toStrs(BtoU8(buf))
	c := U4toStrs(BtoU4(buf))

	if badStr(a) || badStr(b) || badStr(c) {
		t.Fatal("String slice conversion error")
	}
}

// bad Slice slice?
func badSlc(a []Slice) bool {
	if len(a) != cap(a) {
		return true
	}
	d := uintptr(a[0].Data)
	if SliceSize == 12 {
		return len(a) != 2 || d != cn0 || a[0].Len != cn1 || a[0].Cap != cn0
	}
	return len(a) != 1 || badx(int(d)) || badx(a[0].Len) || badx(a[0].Cap)
}

// []String conversions
func Test4b(t *testing.T) {
	buf := []byte(big)
	a := BtoSlcs(buf)
	b := U8toSlcs(BtoU8(buf))
	c := U4toSlcs(BtoU4(buf))

	if badSlc(a) || badSlc(b) || badSlc(c) {
		t.Fatal("Slice slice conversion error")
	}
}

var cmp = [...]string{"", "A", "AA", "AB", "B", "BA", "BB"}

func TestCmp(t *testing.T) {
	for i := len(cmp) - 1; i >= 0; i-- {
		s := cmp[i]
		b := []byte(s)
		if CmpS(s, s) != 0 || CmpB(b, b) != 0 {
			t.Fatal(s, "must be equal to itself")
		}

		for k := i - 1; k >= 0; k-- {
			r := cmp[k]
			a := []byte(r)
			if CmpS(r, s) != -1 || CmpS(s, r) != 1 ||
				CmpB(a, b) != -1 || CmpB(b, a) != 1 {
				t.Fatal(r, "must be less than", s)
			}
		}
	}
}

var mnt = [...]string{
	"", "B", "!", // in1 < in2, out
	"abc", "cde", "bcd",
	"abc", "abd", "abc",
	"SeRhat", "Tansu ", "T#`n+J",
	"SeRgat", "Tantu ", "T#`n+J",
	"Sergat", "TaNtu ", "T#`n+J",
	"Serhat", "TaNsu ", "T#`n+J",
	"NİreÇ", "eŞVkü", "ZE'dhá",
	"Golang", "Python", "L4pe/*",
	"JAVA", "RUST", "NKU\n",
	"致命的", "警告abc", "蚭呤$~s",
}

func TestMean(t *testing.T) {

	for i := len(mnt) - 1; i > 0; i -= 3 {
		// MeanStr(in1,in2) = out ?
		m := MeanStr(mnt[i-2], mnt[i-1])
		if m != mnt[i] {
			t.Fatal("MeanStr: expected:", mnt[i], "got:", m)
		}
		// MeanStr(in2,in1) = out ?
		m2 := MeanStr(mnt[i-1], mnt[i-2])
		if m != m2 {
			t.Fatal("MeanStr: different means:", m, m2)
		}
		// in1 <= MeanStr(in1,in2) < in2 ?
		if !(mnt[i-2] <= m && m < mnt[i-1]) {
			t.Fatal("MeanStr: bad order:", mnt[i-2], m, mnt[i-1])
		}
	}
	for i := len(mnt) - 1; i >= 0; i-- {
		// MeanStr(in,in) = in ?
		m := MeanStr(mnt[i], mnt[i])
		if m != mnt[i] {
			t.Fatal("MeanStr: same expected:", mnt[i], "got:", m)
		}
	}
}
