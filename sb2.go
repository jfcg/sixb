/*	Copyright (c) 2019-present, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sixb

// CmpS compares a, b lexicographically and returns:
//
//	-1 for a < b
//	 0 for a = b
//	 1 for a > b
//
//go:nosplit
func CmpS(a, b string) (r int) {
	n, k := len(a), len(b)
	if n > k {
		n = k
		r++
	} else if n < k {
		r--
	}

	for i := 0; i < n; i++ {
		x, y := a[i], b[i]
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
	}
	return
}

// CmpB compares a, b lexicographically and returns:
//
//	-1 for a < b
//	 0 for a = b
//	 1 for a > b
func CmpB(a, b []byte) int {
	return CmpS(String(a), String(b))
}

// MeanS returns lexicographic average of s1 & s2. It treats ascii specially. The result is
// (len(s1)+len(s2)+1)/2 bytes and may contain unprintable characters or may not be valid utf8.
//
//go:nosplit
func MeanS(s1, s2 string) string {
	if len(s2) < len(s1) {
		s1, s2 = s2, s1
	} else if len(s2) <= 0 {
		return ""
	}
	i := Mean(len(s1)+1, len(s2))
	avg := make([]byte, i)

	i--
	sum := uint32(s2[i])
	mask := sum | 127
	if i < len(s1) {
		c := uint32(s1[i])
		sum += c
		mask |= c // if ascii inputs, consider 7 bits
	}
	for i > 0 {
		prMask := mask
		prSum := sum & mask
		shift := 7 + mask>>7
		sum >>= shift // carry

		i--
		c := uint32(s2[i])
		sum += c
		mask = c | 127
		if i < len(s1) {
			c = uint32(s1[i])
			sum += c
			mask |= c
		}
		avg[i+1] = byte((sum<<(shift-1) ^ prSum>>1) & prMask)
	}
	avg[0] = byte(sum >> 1)

	return String(avg)
}

// Signed is the set of multi-byte signed integer types.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is the set of multi-byte unsigned integer types.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is the set of integer types.
type Integer interface {
	Signed | Unsigned
}

// Mean returns average of integers x & y, mathematically equivalent to floor((x+y)/2).
func Mean[T Integer](x, y T) T {
	return x&y + (x^y)>>1
}

// Median3 returns median of three integers
func Median3[T Integer](a, b, c T) T {
	// almost insertion sort to have a <= b <= c
	if b < a {
		a, b = b, a
	}
	if c < b {
		b = c
		if b < a {
			b = a
		}
	}
	return b
}

// Median4 returns median of four integers
func Median4[T Integer](a, b, c, d T) T {
	if d < b {
		d, b = b, d
	}
	if c < a {
		c, a = a, c
	}
	if d < c {
		c = d
	}
	if b < a {
		b = a
	}
	return Mean(b, c)
}
