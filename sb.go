//	Some string utility functions
package sixb

import "unsafe"

// A quick & dirty Text-to-Int function for short text inputs.
func Txt2int(s string) (x uint64) {
	for i := 0; i < len(s); i++ {
		x = x<<7 ^ x>>57
		x += uint64(s[i] - ' ')
	}
	return x
}

//	Accepts 0..:, @..Z, a..z & maps it onto 6-bits. This is actually a bijection & inverse of Sb2an().
func An2sb(x byte) byte {
	x = 122 - x
	if x > 133 || x < 26 {
		return x
	}
	if x > 63 {
		return x - 11
	}
	if x > 58 {
		return x + 70
	}
	if x > 31 {
		return x - 6
	}
	return x + 97
}

//	Accepts 6-bits & maps it onto 0..:, @..Z, a..z. This is actually a bijection & inverse of An2sb().
func Sb2an(x byte) byte {
	x = 122 - x
	if x < 70 {
		return x - 11
	}
	if x < 97 {
		return x - 6
	}
	if x < 245 {
		return x
	}
	if x < 250 {
		return x + 70
	}
	return x + 97
}

type slice struct { // not worth importing reflect
	Data     uintptr
	Len, Cap int
}

//	Convert byte slice to int slice
func Bs2is(x []byte) (y []uint64) {
	s := (*slice)(unsafe.Pointer(&y))
	t := (*slice)(unsafe.Pointer(&x))
	s.Data = t.Data
	s.Len = t.Len >> 3
	s.Cap = s.Len
	return
}

//	Convert int slice to byte slice
func Is2bs(y []uint64) (x []byte) {
	s := (*slice)(unsafe.Pointer(&y))
	t := (*slice)(unsafe.Pointer(&x))
	t.Data = s.Data
	t.Len = s.Len << 3
	t.Cap = t.Len
	return
}
