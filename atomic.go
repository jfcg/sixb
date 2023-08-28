//	Copyright (c) 2023-present, Serhat Şevki Dinçer.
//	This Source Code Form is subject to the terms of the Mozilla Public
//	License, v. 2.0. If a copy of the MPL was not distributed with this
//	file, You can obtain one at http://mozilla.org/MPL/2.0/.

//go:build 386 || amd64 || arm64

package sixb

// Inc attempts to increment a counter with:
//
//	if ctr == nil || *ctr >= max {
//		return 0
//	}
//	atomic {
//		*ctr++
//		return *ctr
//	}
//
//go:noescape
func Inc(ctr *uint32, max uint32) (new uint32)

// Dec attempts to decrement a counter with:
//
//	if ctr == nil || *ctr == 0 {
//		return ^0 // -1
//	}
//	atomic {
//		*ctr--
//		return *ctr
//	}
//
//go:noescape
func Dec(ctr *uint32) (new uint32)
