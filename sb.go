/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

// Package sixb provides some string utility functions
package sixb

import "unsafe"

// An2sb is a bijection (without fixed points, single cycle and inverse of Sb2an) that
// maps 0-9: @A-Z a-z onto 6-bits
var An2sb = [...]byte{208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220,
	221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237,
	238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254,
	255, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 26, 27, 28, 29,
	30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
	51, 52, 69, 70, 71, 72, 73, 74, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
	16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85,
	86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105,
	106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122,
	123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139,
	140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 207, 154, 155, 156,
	157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173,
	174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190,
	191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 153}

// Sb2an is a bijection (without fixed points, single cycle and inverse of An2sb) that
// maps 6-bits onto 0-9: @A-Z a-z
var Sb2an = [...]byte{97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
	111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 64, 65, 66, 67, 68, 69,
	70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90,
	48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 91, 92, 93, 94, 95,
	96, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138,
	139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155,
	156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172,
	173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189,
	190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 255, 202, 203, 204, 205, 206,
	207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223,
	224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240,
	241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 201, 0, 1, 2, 3,
	4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
	27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47}

// Copy creates an identical copy of x
func Copy(x []byte) []byte {
	r := make([]byte, len(x))
	copy(r, x)
	return r
}

// String internals from reflect
type String struct {
	Data unsafe.Pointer
	Len  int
}

// Slice internals from reflect
type Slice struct {
	String
	Cap int
}

// BtI4 converts byte slice to integer slice
func BtI4(b []byte) (i []uint32) {
	I := (*Slice)(unsafe.Pointer(&i))
	B := (*Slice)(unsafe.Pointer(&b))
	I.Data = B.Data
	I.Len = B.Len >> 2
	I.Cap = I.Len
	return
}

// I4tB converts integer slice to byte slice
func I4tB(i []uint32) (b []byte) {
	I := (*Slice)(unsafe.Pointer(&i))
	B := (*Slice)(unsafe.Pointer(&b))
	B.Data = I.Data
	B.Len = I.Len << 2
	B.Cap = B.Len
	return
}

// I4tI8 converts uint32 slice to uint64 slice
func I4tI8(i []uint32) (k []uint64) {
	I := (*Slice)(unsafe.Pointer(&i))
	K := (*Slice)(unsafe.Pointer(&k))
	K.Data = I.Data
	K.Len = I.Len >> 1
	K.Cap = K.Len
	return
}

// I8tI4 converts uint64 slice to uint32 slice
func I8tI4(i []uint64) (k []uint32) {
	I := (*Slice)(unsafe.Pointer(&i))
	K := (*Slice)(unsafe.Pointer(&k))
	K.Data = I.Data
	K.Len = I.Len << 1
	K.Cap = K.Len
	return
}

// BtI8 converts byte slice to integer slice
func BtI8(b []byte) (i []uint64) {
	I := (*Slice)(unsafe.Pointer(&i))
	B := (*Slice)(unsafe.Pointer(&b))
	I.Data = B.Data
	I.Len = B.Len >> 3
	I.Cap = I.Len
	return
}

// I8tB converts integer slice to byte slice
func I8tB(i []uint64) (b []byte) {
	I := (*Slice)(unsafe.Pointer(&i))
	B := (*Slice)(unsafe.Pointer(&b))
	B.Data = I.Data
	B.Len = I.Len << 3
	B.Cap = B.Len
	return
}

// BtS converts byte slice to string
func BtS(b []byte) (s string) {
	B := (*Slice)(unsafe.Pointer(&b))
	S := (*String)(unsafe.Pointer(&s))
	S.Data = B.Data
	S.Len = B.Len
	return
}

// StB converts string to byte slice
func StB(s string) (b []byte) {
	B := (*Slice)(unsafe.Pointer(&b))
	S := (*String)(unsafe.Pointer(&s))
	B.Data = S.Data
	B.Len = S.Len
	B.Cap = B.Len
	return
}

// StI4 converts string to integer slice
func StI4(s string) (i []uint32) {
	I := (*Slice)(unsafe.Pointer(&i))
	S := (*String)(unsafe.Pointer(&s))
	I.Data = S.Data
	I.Len = S.Len >> 2
	I.Cap = I.Len
	return
}

// I4tS converts integer slice to string
func I4tS(i []uint32) (s string) {
	I := (*Slice)(unsafe.Pointer(&i))
	S := (*String)(unsafe.Pointer(&s))
	S.Data = I.Data
	S.Len = I.Len << 2
	return
}

// StI8 converts string to integer slice
func StI8(s string) (i []uint64) {
	I := (*Slice)(unsafe.Pointer(&i))
	S := (*String)(unsafe.Pointer(&s))
	I.Data = S.Data
	I.Len = S.Len >> 3
	I.Cap = I.Len
	return
}

// I8tS converts integer slice to string
func I8tS(i []uint64) (s string) {
	I := (*Slice)(unsafe.Pointer(&i))
	S := (*String)(unsafe.Pointer(&s))
	S.Data = I.Data
	S.Len = I.Len << 3
	return
}

// string size
const stsz = int(unsafe.Sizeof(""))

// BtSs converts byte slice to String slice
func BtSs(b []byte) (ss []String) {
	B := (*Slice)(unsafe.Pointer(&b))
	S := (*Slice)(unsafe.Pointer(&ss))
	S.Data = B.Data
	S.Len = B.Len / stsz
	S.Cap = S.Len
	return
}

// I4tSs converts integer slice to String slice
func I4tSs(i []uint32) (ss []String) {
	I := (*Slice)(unsafe.Pointer(&i))
	S := (*Slice)(unsafe.Pointer(&ss))
	S.Data = I.Data
	S.Len = 4 * I.Len / stsz
	S.Cap = S.Len
	return
}

// I8tSs converts integer slice to String slice
func I8tSs(i []uint64) (ss []String) {
	I := (*Slice)(unsafe.Pointer(&i))
	S := (*Slice)(unsafe.Pointer(&ss))
	S.Data = I.Data
	S.Len = 8 * I.Len / stsz
	S.Cap = S.Len
	return
}

// CmpS returns -1 for a < b, 0 for a = b, and 1 for a > b lexicographically
func CmpS(a, b string) int32 {
	n, k, r := len(a), len(b), int32(0)
	if n > k {
		n = k
		r++
	} else if n < k {
		r--
	}

	for i := 0; i < n; i++ {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}
	return r
}

// CmpB returns -1 for a < b, 0 for a = b, and 1 for a > b lexicographically
func CmpB(a, b []byte) int32 {
	return CmpS(BtS(a), BtS(b))
}
