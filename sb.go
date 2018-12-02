//	Some string utility functions
package sixb

import "unsafe"

//	A quick & dirty Text-to-Int function for short text inputs.
func Txt2int(s string) uint64 {
	x := uint64(len(s))
	for i := len(s) - 1; i >= 0; i-- {
		x = x<<15 ^ x>>49
		x += uint64(s[i])
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

//	Creates an identical copy.
func Copy(x []byte) []byte {
	r := make([]byte, len(x))
	copy(r, x)
	return r
}

type str struct { // not worth importing reflect
	Data uintptr
	Len  int
}

type slice struct {
	str
	Cap int
}

//	Converts byte slice to int slice.
func BtI4(b []byte) (i []uint32) {
	I := (*slice)(unsafe.Pointer(&i))
	B := (*slice)(unsafe.Pointer(&b))
	I.Data = B.Data
	I.Len = B.Len >> 2
	I.Cap = I.Len
	return
}

//	Converts int slice to byte slice.
func I4tB(i []uint32) (b []byte) {
	I := (*slice)(unsafe.Pointer(&i))
	B := (*slice)(unsafe.Pointer(&b))
	B.Data = I.Data
	B.Len = I.Len << 2
	B.Cap = B.Len
	return
}

//	Converts byte slice to int slice.
func BtI8(b []byte) (i []uint64) {
	I := (*slice)(unsafe.Pointer(&i))
	B := (*slice)(unsafe.Pointer(&b))
	I.Data = B.Data
	I.Len = B.Len >> 3
	I.Cap = I.Len
	return
}

//	Converts int slice to byte slice.
func I8tB(i []uint64) (b []byte) {
	I := (*slice)(unsafe.Pointer(&i))
	B := (*slice)(unsafe.Pointer(&b))
	B.Data = I.Data
	B.Len = I.Len << 3
	B.Cap = B.Len
	return
}

//	Converts string to int slice.
func StI4(s string) (i []uint32) {
	I := (*slice)(unsafe.Pointer(&i))
	S := (*str)(unsafe.Pointer(&s))
	I.Data = S.Data
	I.Len = S.Len >> 2
	I.Cap = I.Len
	return
}

//	Converts int slice to string.
func I4tS(i []uint32) (s string) {
	I := (*slice)(unsafe.Pointer(&i))
	S := (*str)(unsafe.Pointer(&s))
	S.Data = I.Data
	S.Len = I.Len << 2
	return
}

//	Converts string to int slice.
func StI8(s string) (i []uint64) {
	I := (*slice)(unsafe.Pointer(&i))
	S := (*str)(unsafe.Pointer(&s))
	I.Data = S.Data
	I.Len = S.Len >> 3
	I.Cap = I.Len
	return
}

//	Converts int slice to string.
func I8tS(i []uint64) (s string) {
	I := (*slice)(unsafe.Pointer(&i))
	S := (*str)(unsafe.Pointer(&s))
	S.Data = I.Data
	S.Len = I.Len << 3
	return
}
