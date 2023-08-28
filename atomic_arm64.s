//	Copyright (c) 2023-present, Serhat Şevki Dinçer.
//	This Source Code Form is subject to the terms of the Mozilla Public
//	License, v. 2.0. If a copy of the MPL was not distributed with this
//	file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Inc(ctr *uint32, max uint32) (new uint32)
TEXT ·Inc(SB), 4, $0-20
	MOVD	ctr+0(FP), R3
	MOVW	ZR, R0
	CBZ		R3, skip
	MOVW	(R3), R2
	MOVW	max+8(FP), R1
	CMPW	R1, R2
	BHS		skip
	MOVW	$1, R2
	LDADDALW	R2, (R3), R0
	ADDW	$1, R0
skip:
	MOVW	R0, new+16(FP)
	RET

// Dec(ctr *uint32) (new uint32)
TEXT ·Dec(SB), 4, $0-12
	MOVD	ctr+0(FP), R2
	MOVW	ZR, R0
	CBZ		R2, skip
	MOVW	(R2), R1
	CBZW	R1, skip
	MOVW	$-1, R1
	LDADDALW	R1, (R2), R0
skip:
	SUBW	$1, R0
	MOVW	R0, new+8(FP)
	RET
