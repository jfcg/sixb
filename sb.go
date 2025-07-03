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

// Copy creates a copy of s.
func Copy[S ~[]T, T any](s S) []T {
	r := make([]T, len(s))
	copy(r, s)
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

// PtrToInt converts a pointer value to an integer.
func PtrToInt[P ~*T, T any](p P) uint {
	return uint(uintptr(unsafe.Pointer(p)))
}

// InString represents internals of a Go string
type InString struct {
	Data unsafe.Pointer
	Len  uint
}

// InSlice represents internals of a Go slice
type InSlice struct {
	Data unsafe.Pointer
	Len  uint
	Cap  uint
}

func toStr(s *string) *InString {
	return (*InString)(unsafe.Pointer(s))
}

func toSlc[S ~[]T, T any](s *S) *InSlice {
	return (*InSlice)(unsafe.Pointer(s))
}

// Size of x in bytes.
func Size[T any](x T) uint {
	return uint(unsafe.Sizeof(x))
}

// Cast s to an actual slice type.
func Cast[T any](s InSlice) []T {
	return *(*[]T)(unsafe.Pointer(&s))
}

// Slice converts a slice to another slice type, considering
// element type sizes. Be careful with types that contain pointers.
func Slice[U any, S ~[]T, T any](in S) (out []U) {
	src := toSlc(&in)
	dst := toSlc(&out)
	dst.Data = src.Data
	var s T
	var d U
	l, ns, nd := src.Len, Size(s), Size(d)
	if ns != nd {
		l = ns * l / nd
	}
	dst.Len = l
	dst.Cap = l
	return
}

// String converts integer slice (including []byte) to string.
func String[S ~[]T, T Integer](in S) (out string) {
	src := toSlc(&in)
	dst := toStr(&out)
	dst.Data = src.Data
	dst.Len = src.Len * Size(T(0))
	return
}

// Integers converts string to integer slice (including []byte).
func Integers[T Integer](in string) (out []T) {
	src := toStr(&in)
	dst := toSlc(&out)
	dst.Data = src.Data
	n := src.Len / Size(T(0))
	dst.Len = n
	dst.Cap = n
	return
}

// Bytes converts string to byte slice.
func Bytes(s string) []byte { // alias for common case
	return Integers[byte](s)
}
