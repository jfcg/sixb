/*	Copyright (c) 2019-present, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

// Package sixb provides string, slice & integer utility functions.
// string/slice functions help avoid redundant memory allocations.
package sixb

import (
	"os"
	"runtime"
	"strings"
	"unsafe"
)

// AnumToSixb is a bijection (without fixed points, single cycle and inverse of SixbToAnum)
// that maps 0-9: @A-Z a-z onto 6-bits.
var AnumToSixb = [...]byte{208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220,
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

// SixbToAnum is a bijection (without fixed points, single cycle and inverse of AnumToSixb)
// that maps 6-bits onto 0-9: @A-Z a-z.
var SixbToAnum = [...]byte{97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109,
	110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 64, 65, 66, 67, 68,
	69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90,
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

// Copy creates a copy of a slice.
func Copy[S ~[]T, T any](slc S) S {
	r := make(S, len(slc))
	copy(r, slc)
	return r
}

// InsideTest returns true inside a Go test.
func InsideTest() bool {
	suffix := ".test"
	if runtime.GOOS == "windows" {
		suffix += ".exe"
	}
	return len(os.Args) > 1 && strings.HasSuffix(os.Args[0], suffix) &&
		strings.HasPrefix(os.Args[1], "-test.")
}

// SamePtr returns true if pointers a & b are same addresses in memory.
func SamePtr[P1 ~*T, P2 ~*U, T, U any](a P1, b P2) bool {
	return unsafe.Pointer(a) == unsafe.Pointer(b)
}

// PtrToInt converts a pointer to an integer.
func PtrToInt[P ~*T, T any](ptr P) uint {
	return uint(uintptr(unsafe.Pointer(ptr)))
}

// Slice converts a slice to another slice type, considering
// element type sizes. Be careful with types that contain pointers.
func Slice[U any, S ~[]T, T any](slc S) []U {
	var t T
	var u U
	l, st, su := len(slc), int(unsafe.Sizeof(t)), int(unsafe.Sizeof(u))
	if st != su {
		l = st * l / su
	}
	p := unsafe.Pointer(unsafe.SliceData(slc))
	return unsafe.Slice((*U)(p), l)
}

// String converts integer slice (including []byte) to string.
func String[S ~[]T, T Integer](slc S) string {
	p := unsafe.Pointer(unsafe.SliceData(slc))
	return unsafe.String((*byte)(p), len(slc)*int(unsafe.Sizeof(T(0))))
}

// Integers converts string to integer slice (including []byte).
func Integers[U Integer, T ~string](str T) []U {
	p := unsafe.Pointer(unsafe.StringData(string(str)))
	return unsafe.Slice((*U)(p), len(str)/int(unsafe.Sizeof(U(0))))
}

// Bytes converts string to byte slice.
func Bytes[T ~string](str T) []byte { // alias for common case
	return Integers[byte](str)
}
