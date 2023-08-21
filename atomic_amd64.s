//	Copyright (c) 2023-present, Serhat Şevki Dinçer.
//	This Source Code Form is subject to the terms of the Mozilla Public
//	License, v. 2.0. If a copy of the MPL was not distributed with this
//	file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Inc(ctr *uint32, max uint32) (new uint32)
TEXT ·Inc(SB), 4, $0-20
	MOVQ	ctr+0(FP), BX
	XORL	AX, AX
	TESTQ	BX, BX
	JZ		skip
	MOVL	max+8(FP), CX
	CMPL	(BX), CX
	JAE		skip
	INCL	AX
	LOCK
	XADDL	AX, (BX)
	INCL	AX
skip:
	MOVL	AX, new+16(FP)
	RET

// Dec(ctr *uint32) (new uint32)
TEXT ·Dec(SB), 4, $0-12
	MOVQ	ctr+0(FP), BX
	XORL	AX, AX
	TESTQ	BX, BX
	JZ		skip
	CMPL	(BX), AX
	JZ		skip
	DECL	AX
	LOCK
	XADDL	AX, (BX)
skip:
	DECL	AX
	MOVL	AX, new+8(FP)
	RET
