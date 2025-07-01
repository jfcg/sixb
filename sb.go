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
func Copy[T any](s []T) []T {
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
func SamePtr[T, U any](a *T, b *U) bool {
	return unsafe.Pointer(a) == unsafe.Pointer(b)
}

// PtrToInt converts a pointer value to an integer.
func PtrToInt[T any](p *T) uint {
	return uint(uintptr(unsafe.Pointer(p)))
}

type str struct {
	Data unsafe.Pointer
	Len  uint
}

type slc struct {
	Data unsafe.Pointer
	Len  uint
	Cap  uint
}

func toStr(s *string) *str {
	return (*str)(unsafe.Pointer(s))
}

func toSlc[T any](s *[]T) *slc {
	return (*slc)(unsafe.Pointer(s))
}

// Size of x in bytes.
func Size[T any](x T) uint {
	return uint(unsafe.Sizeof(x))
}

// Slice converts a slice to another slice.
func Slice[O, I any](in []I) (out []O) {
	src := toSlc(&in)
	dst := toSlc(&out)
	dst.Data = src.Data
	var s I
	var d O
	n := Size(s) * src.Len / Size(d)
	dst.Len = n
	dst.Cap = n
	return
}

// String converts integer slice (including []byte) to string.
func String[T Integer](i []T) (s string) {
	src := toSlc(&i)
	dst := toStr(&s)
	dst.Data = src.Data
	dst.Len = src.Len * Size(T(0))
	return
}

// Integers converts string to integer slice (including []byte).
func Integers[T Integer](s string) (i []T) {
	src := toStr(&s)
	dst := toSlc(&i)
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
