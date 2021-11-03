/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sixb

// CmpS compares a,b lexicographically and returns:
//  -1 for a < b
//   0 for a = b
//   1 for a > b
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

// CmpB compares a,b lexicographically and returns:
//  -1 for a < b
//   0 for a = b
//   1 for a > b
func CmpB(a, b []byte) int {
	return CmpS(BtoS(a), BtoS(b))
}

// MeanS returns lexicographic average of s1 & s2, up to 31 bytes. It treats ascii
// specially. Results may contain unprintable characters or may not be valid utf8.
func MeanS(s1, s2 string) string {
	if len(s2) < len(s1) {
		s1, s2 = s2, s1
	} else if len(s2) <= 0 {
		return ""
	}
	i := len(s2)
	if i > 31 {
		i = 31
	}
	avg := make([]byte, i+1) // extra byte

	var sum, mask, prsum, prmask, shift uint32
	goto start
loop:
	prmask = mask
	prsum = sum & mask
	shift = 7 + mask>>7
	sum >>= shift // carry
start:
	i--
	c := uint32(s2[i])
	sum += c
	mask = c | 127
	if i < len(s1) {
		c = uint32(s1[i])
		sum += c
		mask |= c // if ascii inputs, consider 7 bits
	}
	avg[i+1] = byte((sum<<(shift-1) ^ prsum>>1) & prmask)
	if i > 0 {
		goto loop
	}
	avg[0] = byte(sum >> 1)

	return string(avg[:len(avg)-1])
}

// MeanU4 returns average of x, y. Mathematically equivalent to floor((x+y)/2).
func MeanU4(x, y uint32) uint32 {
	return x&y + (x^y)>>1
}

// MeanI4 returns average of x, y. Mathematically equivalent to floor((x+y)/2).
func MeanI4(x, y int32) int32 {
	return x&y + (x^y)>>1
}

// MeanU8 returns average of x, y. Mathematically equivalent to floor((x+y)/2).
func MeanU8(x, y uint64) uint64 {
	return x&y + (x^y)>>1
}

// MeanI8 returns average of x, y. Mathematically equivalent to floor((x+y)/2).
func MeanI8(x, y int64) int64 {
	return x&y + (x^y)>>1
}

// MeanU returns average of x, y. Mathematically equivalent to floor((x+y)/2).
func MeanU(x, y uint) uint {
	return x&y + (x^y)>>1
}

// MeanI returns average of x, y. Mathematically equivalent to floor((x+y)/2).
func MeanI(x, y int) int {
	return x&y + (x^y)>>1
}
