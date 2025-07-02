/*	Copyright (c) 2019-present, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sixb

import (
	"cmp"
	"testing"
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

	if Size(buf) != Size(slc{}) ||
		len(y) != 1 || cap(y) != 1 || y[0] != cn2 ||
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

	if len(y) != 1 || cap(y) != 1 || y[0] != cn2 ||
		len(p) != 2 || cap(p) != 2 || p[0] != cn0 || p[1] != cn1 ||
		len(z) != 2 || cap(z) != 2 || z[0] != cn0 || z[1] != cn1 ||
		!SamePtr(&y[0], &p[0]) || !SamePtr(&y[0], &z[0]) ||
		!SamePtr(&y[0], &buf[0]) {
		t.Fatal("slice conversion error")
	}
}

// nil string/slice conversions
func TestSlice3(t *testing.T) {
	var s string
	var a []byte

	if Bytes(s) != nil || Integers[uint32](s) != nil ||
		Integers[uint64](s) != nil || Slice[uint32](a) != nil ||
		Slice[uint64](a) != nil || Slice[str](a) != nil || Slice[slc](a) != nil {
		t.Fatal("nil string/slice conversion error")
	}
}

func TestSlice3b(t *testing.T) {
	var b []uint32
	var c []uint64

	if String(b) != "" || String(c) != "" || Slice[byte](b) != nil ||
		Slice[byte](c) != nil || Slice[str](b) != nil || Slice[str](c) != nil ||
		Slice[slc](b) != nil || Slice[slc](c) != nil {
		t.Fatal("nil string/slice conversion error")
	}
}

// string conversions
func TestString(t *testing.T) {
	a := Integers[uint64](sml)
	b := String(a)
	r := Integers[uint32](sml)
	s := String(r)

	if Size([]byte{}) != Size(slc{}) ||
		Size("") != Size(str{}) ||
		len(a) != 1 || cap(a) != 1 || a[0] != cn2 ||
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

// bad String slice?
func badStr(a []str) bool {
	if len(a) != cap(a) {
		return true
	}
	d := uint(uintptr(a[0].Data))
	if Size("") == 8 {
		return len(a) != 3 || d != cn0 || a[0].Len != cn1
	}
	return len(a) != 1 || badx(d) || badx(a[0].Len)
}

// []String conversions
func TestString3(t *testing.T) {
	buf := []byte(big)
	a := Slice[str](buf)
	b := Slice[str](Slice[uint64](buf))
	c := Slice[str](Slice[uint32](buf))

	if badStr(a) || badStr(b) || badStr(c) {
		t.Fatal("String slice conversion error")
	}
}

// bad Slice slice?
func badSlc(a []slc) bool {
	if len(a) != cap(a) {
		return true
	}
	d := uint(uintptr(a[0].Data))
	if Size([]byte{}) == 12 {
		return len(a) != 2 || d != cn0 || a[0].Len != cn1 || a[0].Cap != cn0
	}
	return len(a) != 1 || badx(d) || badx(a[0].Len) || badx(a[0].Cap)
}

// []String conversions
func TestString4(t *testing.T) {
	buf := []byte(big)
	a := Slice[slc](buf)
	b := Slice[slc](Slice[uint64](buf))
	c := Slice[slc](Slice[uint32](buf))

	if badSlc(a) || badSlc(b) || badSlc(c) {
		t.Fatal("Slice slice conversion error")
	}
}

func TestPtoU8(t *testing.T) {
	var p *int
	if PtrToInt(p) != 0 {
		t.Fatal("Nil pointer must convert to to zero")
	}
	var arr [2]int
	v1 := PtrToInt(&arr[0])
	v2 := PtrToInt(&arr[1])
	if PtrToInt(t) == 0 || v1 == 0 || v2 == 0 ||
		v1+Size(arr[0]) != v2 {
		t.Fatal("Pointer to unsigned integer conversion error")
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

func isort(a, b, c int) (int, int, int) {
	// insertion sort to have a >= b >= c
	if b > a {
		a, b = b, a
	}
	if c > b {
		b, c = c, b
		if b > a {
			a, b = b, a
		}
	}
	return a, b, c
}

func TestMedian(tst *testing.T) {
	const N = 5
	for i := -N; i <= N; i++ {
		for k := -N; k <= N; k++ {
			for r := -N; r <= N; r++ {

				a, b, c := isort(i, k, r)
				if m := Median3(i, k, r); m != b {
					tst.Fatal("expected:", b, "got:", m)
				}

				for p := -N; p <= N; p++ {

					b, c := b, c
					// almost insertion sort to have a >= b >= c >= p
					if p > c {
						c = p
						if c > b {
							b, c = c, b
							if b > a {
								b = a
							}
						}
					}
					if m, e := Median4(i, k, r, p), Mean(b, c); m != e {
						tst.Fatal("expected:", e, "got:", m)
					}
				}
			}
		}
	}
}

func BenchmarkMeanS(b *testing.B) {
	res, l := "", len(strTable)-1
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res = MeanS(strTable[l-2], strTable[l-1])
	}
	b.StopTimer()
	if res != strTable[l] {
		b.Fatal("MeanS error")
	}
}

func meanTest[O cmp.Ordered](t *testing.T, mean func(O, O) O, table []O) {
	for i := len(table) - 1; i > 1; i -= 3 {
		// mean(in1,in2) = out ?
		m := mean(table[i-2], table[i-1])
		if m != table[i] {
			t.Fatal("expected:", table[i], "got:", m)
		}
		// mean(in2,in1) = out ?
		m2 := mean(table[i-1], table[i-2])
		if m != m2 {
			t.Fatal("different means:", m, m2)
		}
		// in1 <= mean(in1,in2) < in2 ?
		if !(table[i-2] <= m && m < table[i-1]) {
			t.Fatal("bad order:", table[i-2], m, table[i-1])
		}
	}
	for i := len(table) - 1; i >= 0; i-- {
		// mean(in,in) = in ?
		m := mean(table[i], table[i])
		if m != table[i] {
			t.Fatal("expected:", table[i], "got:", m)
		}
	}
}

func TestMeans(t *testing.T) {
	meanTest(t, MeanS, strTable)
	meanTest(t, Mean, u4Table)
	meanTest(t, Mean, i4Table)
	meanTest(t, Mean, u8Table)
	meanTest(t, Mean, i8Table)
}

var strTable = []string{
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

var u4Table = []uint32{
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

var i4Table = []int32{
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

var u8Table = []uint64{
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

var i8Table = []int64{
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
