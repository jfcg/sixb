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
	for b.Loop() {
		res = MeanS(strTable[l-2], strTable[l-1])
	}
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
