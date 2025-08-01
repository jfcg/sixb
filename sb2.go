/*	Copyright (c) 2019-present, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sixb

import "cmp"

// MeanS returns lexicographic average of s1 & s2. It treats ascii specially. The result is
// (len(s1)+len(s2)+1)/2 bytes and may contain unprintable characters or may not be valid utf8.
//
//go:nosplit
func MeanS[T ~string](s1, s2 T) T {
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

	return T(String(avg))
}

// Signed is the set of signed integer types.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is the set of unsigned integer types.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is the set of integer types.
type Integer interface {
	Signed | Unsigned
}

type Float interface {
	~float32 | ~float64
}

// Mean returns average of integers x & y, mathematically equivalent to floor((x+y)/2).
func Mean[T Integer](x, y T) T {
	return x&y + (x^y)>>1
}

// Median3 returns median of three values
func Median3[T cmp.Ordered](a, b, c T) T {
	// almost insertion sort to have a ≤ b ≤ c
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

// None represents empty struct with zero size
type None struct{}

// Set of elements of type K. Create empty set with
//
//	s := Set[K]{}
type Set[K comparable] map[K]None

// NewSet creates a set from given elements
func NewSet[K comparable](els ...K) Set[K] {
	s := Set[K]{}
	s.Add(els...)
	return s
}

// Size returns number of elements in set
func (s Set[K]) Size() int {
	return len(s)
}

// HasAny returns true only if set contains any of given elements.
// If els is empty, it returns false by definition.
func (s Set[K]) HasAny(els ...K) bool {
	for _, e := range els {
		if _, ok := s[e]; ok {
			return true
		}
	}
	return false
}

// HasAll returns true only if set contains all of given elements.
// If els is empty, it returns true by definition.
func (s Set[K]) HasAll(els ...K) bool {
	for _, e := range els {
		if _, ok := s[e]; !ok {
			return false
		}
	}
	return true
}

// Add element(s) to set
func (s Set[K]) Add(els ...K) {
	for _, e := range els {
		s[e] = None{}
	}
}

// Remove element(s) from set
func (s Set[K]) Remove(els ...K) {
	for _, e := range els {
		delete(s, e)
	}
}
