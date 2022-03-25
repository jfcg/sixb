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
	if &b[0] == &buf[0] || BtoS(b) != BtoS(buf) {
		t.Fatal("Copy does not work")
	}
}

// bad byte slice?
func badb(q []byte) bool {
	return len(q) != 8 || cap(q) != 8 || &q[0] != &buf[0]
}

// slice conversions
func TestSlice(t *testing.T) {
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
func TestSlice2(t *testing.T) {
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
func TestSlice3(t *testing.T) {
	var s string
	var a []byte

	if StoB(s) != nil || StoU4(s) != nil || StoU8(s) != nil || BtoU4(a) != nil ||
		BtoU8(a) != nil || BtoStrs(a) != nil || BtoSlcs(a) != nil {
		t.Fatal("nil string/slice conversion error")
	}
}

func TestSlice3b(t *testing.T) {
	var b []uint32
	var c []uint64

	if U4toS(b) != "" || U8toS(c) != "" || U4toB(b) != nil ||
		U8toB(c) != nil || U4toStrs(b) != nil || U8toStrs(c) != nil ||
		U4toSlcs(b) != nil || U8toSlcs(c) != nil {
		t.Fatal("nil string/slice conversion error")
	}
}

// string conversions
func TestString(t *testing.T) {
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
func TestString2(t *testing.T) {
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
func TestString3(t *testing.T) {
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
func TestString4(t *testing.T) {
	buf := []byte(big)
	a := BtoSlcs(buf)
	b := U8toSlcs(BtoU8(buf))
	c := U4toSlcs(BtoU4(buf))

	if badSlc(a) || badSlc(b) || badSlc(c) {
		t.Fatal("Slice slice conversion error")
	}
}

var cmpArr = [...]string{"", "A", "AA", "AB", "B", "BA", "BB"}

func TestCmp(t *testing.T) {
	for i := len(cmpArr) - 1; i >= 0; i-- {
		s := cmpArr[i]
		b := []byte(s)
		if CmpS(s, s) != 0 || CmpB(b, b) != 0 {
			t.Fatal(s, "must be equal to itself")
		}

		for k := i - 1; k >= 0; k-- {
			r := cmpArr[k]
			a := []byte(r)
			if CmpS(r, s) != -1 || CmpS(s, r) != 1 ||
				CmpB(a, b) != -1 || CmpB(b, a) != 1 {
				t.Fatal(r, "must be less than", s)
			}
		}
	}
}

var strTable = [...]string{
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

func BenchmarkMeanS(b *testing.B) {
	var res string
	for i := 0; i < b.N; i++ {
		res = MeanS(strTable[9], strTable[10])
	}
	if res != strTable[11] {
		b.Fatal("MeanS error")
	}
}

func TestMeanS(t *testing.T) {
	for i := len(strTable) - 1; i > 0; i -= 3 {
		// MeanS(in1,in2) = out ?
		m := MeanS(strTable[i-2], strTable[i-1])
		if m != strTable[i] {
			t.Fatal("MeanS: expected:", strTable[i], "got:", m)
		}
		// MeanS(in2,in1) = out ?
		m2 := MeanS(strTable[i-1], strTable[i-2])
		if m != m2 {
			t.Fatal("MeanS: different means:", m, m2)
		}
		// in1 <= MeanS(in1,in2) < in2 ?
		if !(strTable[i-2] <= m && m < strTable[i-1]) {
			t.Fatal("MeanS: bad order:", strTable[i-2], m, strTable[i-1])
		}
	}
	for i := len(strTable) - 1; i >= 0; i-- {
		// MeanS(in,in) = in ?
		m := MeanS(strTable[i], strTable[i])
		if m != strTable[i] {
			t.Fatal("MeanS: same expected:", strTable[i], "got:", m)
		}
	}
}

var u4Table = [...]uint32{
	0, 1, 0,
	100, 200, 150, // in1 < in2, out
	101, 200, 150,

	1<<31 - 200, 1<<31 - 100, 1<<31 - 150,
	1<<31 - 200, 1<<31 - 101, 1<<31 - 151,
	1<<32 - 200, 1<<32 - 100, 1<<32 - 150,
	1<<32 - 200, 1<<32 - 101, 1<<32 - 151,

	1<<31 - 200, 1<<32 - 100, 3<<30 - 150,
	1<<31 - 200, 1<<32 - 101, 3<<30 - 151,
	1<<31 - 100, 1<<32 - 200, 3<<30 - 150,
	1<<31 - 100, 1<<32 - 201, 3<<30 - 151,

	1<<31 - 1, 1 << 31, 1<<31 - 1,
	1 << 31, 1<<31 + 1, 1 << 31,
	1<<31 - 1, 1<<32 - 1, 3<<30 - 1,
	1<<32 - 2, 1<<32 - 1, 1<<32 - 2,
}

func TestMeanU4(t *testing.T) {
	for i := len(u4Table) - 1; i > 0; i -= 3 {
		// MeanU4(in1,in2) = out ?
		m := MeanU4(u4Table[i-2], u4Table[i-1])
		if m != u4Table[i] {
			t.Fatal("MeanU4: expected:", u4Table[i], "got:", m)
		}
		// MeanU4(in2,in1) = out ?
		m2 := MeanU4(u4Table[i-1], u4Table[i-2])
		if m != m2 {
			t.Fatal("MeanU4: different means:", m, m2)
		}
		// in1 <= MeanU4(in1,in2) < in2 ?
		if !(u4Table[i-2] <= m && m < u4Table[i-1]) {
			t.Fatal("MeanU4: bad order:", u4Table[i-2], m, u4Table[i-1])
		}
	}
	for i := len(u4Table) - 1; i >= 0; i-- {
		// MeanU4(in,in) = in ?
		m := MeanU4(u4Table[i], u4Table[i])
		if m != u4Table[i] {
			t.Fatal("MeanU4: same expected:", u4Table[i], "got:", m)
		}
	}
}

var i4Table = [...]int32{
	100, 200, 150, // in1 < in2, out
	101, 200, 150,
	-200, -100, -150,
	-200, -101, -151,

	-100, 100, 0,
	-101, 101, 0,
	100 - 1<<31, 1<<31 - 100, 0,
	101 - 1<<31, 1<<31 - 101, 0,

	1<<31 - 200, 1<<31 - 100, 1<<31 - 150,
	1<<31 - 200, 1<<31 - 101, 1<<31 - 151,
	100 - 1<<31, 200 - 1<<31, 150 - 1<<31,
	101 - 1<<31, 200 - 1<<31, 150 - 1<<31,

	100 - 1<<31, 1<<31 - 200, -50,
	100 - 1<<31, 1<<31 - 201, -51,
	200 - 1<<31, 1<<31 - 100, 50,
	201 - 1<<31, 1<<31 - 100, 50,

	1 - 1<<31, 1<<31 - 1, 0,
	-1 << 31, 1<<31 - 1, -1,
	-1 << 31, 1 - 1<<31, -1 << 31,
}

func TestMeanI4(t *testing.T) {
	for i := len(i4Table) - 1; i > 0; i -= 3 {
		// MeanI4(in1,in2) = out ?
		m := MeanI4(i4Table[i-2], i4Table[i-1])
		if m != i4Table[i] {
			t.Fatal("MeanI4: expected:", i4Table[i], "got:", m)
		}
		// MeanI4(in2,in1) = out ?
		m2 := MeanI4(i4Table[i-1], i4Table[i-2])
		if m != m2 {
			t.Fatal("MeanI4: different means:", m, m2)
		}
		// in1 <= MeanI4(in1,in2) < in2 ?
		if !(i4Table[i-2] <= m && m < i4Table[i-1]) {
			t.Fatal("MeanI4: bad order:", i4Table[i-2], m, i4Table[i-1])
		}
	}
	for i := len(i4Table) - 1; i >= 0; i-- {
		// MeanI4(in,in) = in ?
		m := MeanI4(i4Table[i], i4Table[i])
		if m != i4Table[i] {
			t.Fatal("MeanI4: same expected:", i4Table[i], "got:", m)
		}
	}
}

var u8Table = [...]uint64{
	0, 1, 0,
	100, 200, 150, // in1 < in2, out
	101, 200, 150,

	1<<63 - 200, 1<<63 - 100, 1<<63 - 150,
	1<<63 - 200, 1<<63 - 101, 1<<63 - 151,
	1<<64 - 200, 1<<64 - 100, 1<<64 - 150,
	1<<64 - 200, 1<<64 - 101, 1<<64 - 151,

	1<<63 - 200, 1<<64 - 100, 3<<62 - 150,
	1<<63 - 200, 1<<64 - 101, 3<<62 - 151,
	1<<63 - 100, 1<<64 - 200, 3<<62 - 150,
	1<<63 - 100, 1<<64 - 201, 3<<62 - 151,

	1<<63 - 1, 1 << 63, 1<<63 - 1,
	1 << 63, 1<<63 + 1, 1 << 63,
	1<<63 - 1, 1<<64 - 1, 3<<62 - 1,
	1<<64 - 2, 1<<64 - 1, 1<<64 - 2,
}

func TestMeanU8(t *testing.T) {
	for i := len(u8Table) - 1; i > 0; i -= 3 {
		// MeanU8(in1,in2) = out ?
		m := MeanU8(u8Table[i-2], u8Table[i-1])
		if m != u8Table[i] {
			t.Fatal("MeanU8: expected:", u8Table[i], "got:", m)
		}
		// MeanU8(in2,in1) = out ?
		m2 := MeanU8(u8Table[i-1], u8Table[i-2])
		if m != m2 {
			t.Fatal("MeanU8: different means:", m, m2)
		}
		// in1 <= MeanU8(in1,in2) < in2 ?
		if !(u8Table[i-2] <= m && m < u8Table[i-1]) {
			t.Fatal("MeanU8: bad order:", u8Table[i-2], m, u8Table[i-1])
		}
	}
	for i := len(u8Table) - 1; i >= 0; i-- {
		// MeanU8(in,in) = in ?
		m := MeanU8(u8Table[i], u8Table[i])
		if m != u8Table[i] {
			t.Fatal("MeanU8: same expected:", u8Table[i], "got:", m)
		}
	}
}

var i8Table = [...]int64{
	100, 200, 150, // in1 < in2, out
	101, 200, 150,
	-200, -100, -150,
	-200, -101, -151,

	-100, 100, 0,
	-101, 101, 0,
	100 - 1<<63, 1<<63 - 100, 0,
	101 - 1<<63, 1<<63 - 101, 0,

	1<<63 - 200, 1<<63 - 100, 1<<63 - 150,
	1<<63 - 200, 1<<63 - 101, 1<<63 - 151,
	100 - 1<<63, 200 - 1<<63, 150 - 1<<63,
	101 - 1<<63, 200 - 1<<63, 150 - 1<<63,

	100 - 1<<63, 1<<63 - 200, -50,
	100 - 1<<63, 1<<63 - 201, -51,
	200 - 1<<63, 1<<63 - 100, 50,
	201 - 1<<63, 1<<63 - 100, 50,

	1 - 1<<63, 1<<63 - 1, 0,
	-1 << 63, 1<<63 - 1, -1,
	-1 << 63, 1 - 1<<63, -1 << 63,
}

func TestMeanI8(t *testing.T) {
	for i := len(i8Table) - 1; i > 0; i -= 3 {
		// MeanI8(in1,in2) = out ?
		m := MeanI8(i8Table[i-2], i8Table[i-1])
		if m != i8Table[i] {
			t.Fatal("MeanI8: expected:", i8Table[i], "got:", m)
		}
		// MeanI8(in2,in1) = out ?
		m2 := MeanI8(i8Table[i-1], i8Table[i-2])
		if m != m2 {
			t.Fatal("MeanI8: different means:", m, m2)
		}
		// in1 <= MeanI8(in1,in2) < in2 ?
		if !(i8Table[i-2] <= m && m < i8Table[i-1]) {
			t.Fatal("MeanI8: bad order:", i8Table[i-2], m, i8Table[i-1])
		}
	}
	for i := len(i8Table) - 1; i >= 0; i-- {
		// MeanI8(in,in) = in ?
		m := MeanI8(i8Table[i], i8Table[i])
		if m != i8Table[i] {
			t.Fatal("MeanI8: same expected:", i8Table[i], "got:", m)
		}
	}
}
