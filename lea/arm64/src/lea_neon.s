	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 14, 0
	.globl	_lea_encrypt_4block             ; -- Begin function lea_encrypt_4block
	.p2align	2
_lea_encrypt_4block:                    ; @lea_encrypt_4block
; %bb.0:
	stp	x29, x30, [sp, #-16]!           ; 16-byte Folded Spill
	mov	x29, sp
	add	x8, x0, #16
	ld1r.4s	{ v0 }, [x8]
	ld4.4s	{ v1, v2, v3, v4 }, [x2]
	eor.16b	v0, v0, v3
	add	x8, x0, #20
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v5, v5, v4
	add.4s	v0, v5, v0
	ushr.4s	v5, v0, #3
	sli.4s	v5, v0, #29
	add	x8, x0, #8
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #12
	ld1r.4s	{ v6 }, [x8]
	eor.16b	v6, v6, v3
	add.4s	v0, v6, v0
	ushr.4s	v6, v0, #5
	sli.4s	v6, v0, #27
	mov	x8, x0
	ld1r.4s	{ v0 }, [x8], #4
	eor.16b	v0, v0, v1
	ld1r.4s	{ v7 }, [x8]
	eor.16b	v7, v7, v2
	add.4s	v0, v7, v0
	ushr.4s	v7, v0, #23
	sli.4s	v7, v0, #9
	add	x8, x0, #40
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v5
	add	x8, x0, #44
	ld1r.4s	{ v16 }, [x8]
	eor.16b	v1, v16, v1
	add.4s	v0, v1, v0
	ushr.4s	v1, v0, #3
	sli.4s	v1, v0, #29
	add	x8, x0, #32
	ld1r.4s	{ v0 }, [x8]
	add	x8, x0, #36
	ld1r.4s	{ v2 }, [x8]
	eor.16b	v0, v0, v6
	eor.16b	v2, v2, v5
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #5
	sli.4s	v2, v0, #27
	add	x8, x0, #24
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v7
	add	x8, x0, #28
	ld1r.4s	{ v3 }, [x8]
	eor.16b	v3, v3, v6
	add.4s	v0, v3, v0
	ushr.4s	v3, v0, #23
	add	x8, x0, #64
	ld1r.4s	{ v4 }, [x8]
	sli.4s	v3, v0, #9
	eor.16b	v0, v4, v1
	add	x8, x0, #68
	ld1r.4s	{ v4 }, [x8]
	eor.16b	v4, v4, v7
	add.4s	v0, v4, v0
	ushr.4s	v4, v0, #3
	sli.4s	v4, v0, #29
	add	x8, x0, #56
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #60
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v1, v5, v1
	add.4s	v0, v1, v0
	ushr.4s	v1, v0, #5
	sli.4s	v1, v0, #27
	add	x8, x0, #48
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v3
	add	x8, x0, #52
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v2, v5, v2
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #23
	sli.4s	v2, v0, #9
	add	x8, x0, #88
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v4
	add	x8, x0, #92
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v3, v5, v3
	add.4s	v0, v3, v0
	ushr.4s	v3, v0, #3
	sli.4s	v3, v0, #29
	add	x8, x0, #80
	ld1r.4s	{ v0 }, [x8]
	add	x8, x0, #84
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v0, v0, v1
	eor.16b	v4, v5, v4
	add.4s	v0, v4, v0
	ushr.4s	v4, v0, #5
	sli.4s	v4, v0, #27
	add	x8, x0, #72
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #76
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v1, v5, v1
	add.4s	v0, v1, v0
	ushr.4s	v1, v0, #23
	add	x8, x0, #112
	ld1r.4s	{ v5 }, [x8]
	sli.4s	v1, v0, #9
	eor.16b	v0, v5, v3
	add	x8, x0, #116
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v2, v5, v2
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #3
	sli.4s	v2, v0, #29
	add	x8, x0, #104
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v4
	add	x8, x0, #108
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v3, v5, v3
	add.4s	v0, v3, v0
	ushr.4s	v3, v0, #5
	sli.4s	v3, v0, #27
	add	x8, x0, #96
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v1
	add	x8, x0, #100
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v4, v5, v4
	add.4s	v0, v4, v0
	ushr.4s	v4, v0, #23
	sli.4s	v4, v0, #9
	add	x8, x0, #136
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #140
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v1, v5, v1
	add.4s	v0, v1, v0
	ushr.4s	v1, v0, #3
	sli.4s	v1, v0, #29
	add	x8, x0, #128
	ld1r.4s	{ v0 }, [x8]
	add	x8, x0, #132
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v0, v0, v3
	eor.16b	v2, v5, v2
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #5
	sli.4s	v2, v0, #27
	add	x8, x0, #120
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v4
	add	x8, x0, #124
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v3, v5, v3
	add.4s	v0, v3, v0
	ushr.4s	v3, v0, #23
	add	x8, x0, #160
	ld1r.4s	{ v5 }, [x8]
	sli.4s	v3, v0, #9
	eor.16b	v0, v5, v1
	add	x8, x0, #164
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v4, v5, v4
	add.4s	v0, v4, v0
	ushr.4s	v4, v0, #3
	sli.4s	v4, v0, #29
	add	x8, x0, #152
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #156
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v1, v5, v1
	add.4s	v0, v1, v0
	ushr.4s	v1, v0, #5
	sli.4s	v1, v0, #27
	add	x8, x0, #144
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v3
	add	x8, x0, #148
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v2, v5, v2
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #23
	sli.4s	v2, v0, #9
	add	x8, x0, #184
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v4
	add	x8, x0, #188
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v3, v5, v3
	add.4s	v0, v3, v0
	ushr.4s	v3, v0, #3
	sli.4s	v3, v0, #29
	add	x8, x0, #176
	ld1r.4s	{ v0 }, [x8]
	add	x8, x0, #180
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v0, v0, v1
	eor.16b	v4, v5, v4
	add.4s	v0, v4, v0
	ushr.4s	v4, v0, #5
	sli.4s	v4, v0, #27
	add	x8, x0, #168
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #172
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v1, v5, v1
	add.4s	v0, v1, v0
	ushr.4s	v1, v0, #23
	add	x8, x0, #208
	ld1r.4s	{ v5 }, [x8]
	sli.4s	v1, v0, #9
	eor.16b	v0, v5, v3
	add	x8, x0, #212
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v2, v5, v2
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #3
	sli.4s	v2, v0, #29
	add	x8, x0, #200
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v4
	add	x8, x0, #204
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v3, v5, v3
	add.4s	v0, v3, v0
	ushr.4s	v3, v0, #5
	sli.4s	v3, v0, #27
	add	x8, x0, #192
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v1
	add	x8, x0, #196
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v4, v5, v4
	add.4s	v0, v4, v0
	ushr.4s	v4, v0, #23
	sli.4s	v4, v0, #9
	add	x8, x0, #232
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #236
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v1, v5, v1
	add.4s	v0, v1, v0
	ushr.4s	v1, v0, #3
	sli.4s	v1, v0, #29
	add	x8, x0, #224
	ld1r.4s	{ v0 }, [x8]
	add	x8, x0, #228
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v0, v0, v3
	eor.16b	v2, v5, v2
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #5
	sli.4s	v2, v0, #27
	add	x8, x0, #216
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v4
	add	x8, x0, #220
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v3, v5, v3
	add.4s	v0, v3, v0
	ushr.4s	v3, v0, #23
	add	x8, x0, #256
	ld1r.4s	{ v5 }, [x8]
	sli.4s	v3, v0, #9
	eor.16b	v0, v5, v1
	add	x8, x0, #260
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v4, v5, v4
	add.4s	v0, v4, v0
	ushr.4s	v4, v0, #3
	sli.4s	v4, v0, #29
	add	x8, x0, #248
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #252
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v1, v5, v1
	add.4s	v0, v1, v0
	ushr.4s	v1, v0, #5
	sli.4s	v1, v0, #27
	add	x8, x0, #240
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v3
	add	x8, x0, #244
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v2, v5, v2
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #23
	sli.4s	v2, v0, #9
	add	x8, x0, #280
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v4
	add	x8, x0, #284
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v3, v5, v3
	add.4s	v0, v3, v0
	ushr.4s	v3, v0, #3
	sli.4s	v3, v0, #29
	add	x8, x0, #272
	ld1r.4s	{ v0 }, [x8]
	add	x8, x0, #276
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v0, v0, v1
	eor.16b	v4, v5, v4
	add.4s	v0, v4, v0
	ushr.4s	v4, v0, #5
	sli.4s	v4, v0, #27
	add	x8, x0, #264
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #268
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v1, v5, v1
	add.4s	v0, v1, v0
	ushr.4s	v1, v0, #23
	add	x8, x0, #304
	ld1r.4s	{ v5 }, [x8]
	sli.4s	v1, v0, #9
	eor.16b	v0, v5, v3
	add	x8, x0, #308
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v2, v5, v2
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #3
	sli.4s	v2, v0, #29
	add	x8, x0, #296
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v4
	add	x8, x0, #300
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v3, v5, v3
	add.4s	v0, v3, v0
	ushr.4s	v3, v0, #5
	sli.4s	v3, v0, #27
	add	x8, x0, #288
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v1
	add	x8, x0, #292
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v4, v5, v4
	add.4s	v0, v4, v0
	ushr.4s	v4, v0, #23
	sli.4s	v4, v0, #9
	add	x8, x0, #328
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #332
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v1, v5, v1
	add.4s	v0, v1, v0
	ushr.4s	v1, v0, #3
	sli.4s	v1, v0, #29
	add	x8, x0, #320
	ld1r.4s	{ v0 }, [x8]
	add	x8, x0, #324
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v0, v0, v3
	eor.16b	v2, v5, v2
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #5
	sli.4s	v2, v0, #27
	add	x8, x0, #312
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v4
	add	x8, x0, #316
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v3, v5, v3
	add.4s	v0, v3, v0
	ushr.4s	v3, v0, #23
	add	x8, x0, #352
	ld1r.4s	{ v5 }, [x8]
	sli.4s	v3, v0, #9
	eor.16b	v0, v5, v1
	add	x8, x0, #356
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v4, v5, v4
	add.4s	v0, v4, v0
	ushr.4s	v4, v0, #3
	sli.4s	v4, v0, #29
	add	x8, x0, #344
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #348
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v1, v5, v1
	add.4s	v0, v1, v0
	ushr.4s	v1, v0, #5
	sli.4s	v1, v0, #27
	add	x8, x0, #336
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v3
	add	x8, x0, #340
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v2, v5, v2
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #23
	sli.4s	v2, v0, #9
	add	x8, x0, #376
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v4
	add	x8, x0, #380
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v3, v5, v3
	add.4s	v0, v3, v0
	ushr.4s	v3, v0, #3
	sli.4s	v3, v0, #29
	add	x8, x0, #368
	ld1r.4s	{ v0 }, [x8]
	add	x8, x0, #372
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v0, v0, v1
	eor.16b	v4, v5, v4
	add.4s	v0, v4, v0
	ushr.4s	v4, v0, #5
	sli.4s	v4, v0, #27
	add	x8, x0, #360
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #364
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v1, v5, v1
	add.4s	v0, v1, v0
	ushr.4s	v1, v0, #23
	add	x8, x0, #400
	ld1r.4s	{ v5 }, [x8]
	sli.4s	v1, v0, #9
	eor.16b	v0, v5, v3
	add	x8, x0, #404
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v2, v5, v2
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #3
	sli.4s	v2, v0, #29
	add	x8, x0, #392
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v4
	add	x8, x0, #396
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v3, v5, v3
	add.4s	v0, v3, v0
	ushr.4s	v3, v0, #5
	sli.4s	v3, v0, #27
	add	x8, x0, #384
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v1
	add	x8, x0, #388
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v4, v5, v4
	add.4s	v0, v4, v0
	ushr.4s	v4, v0, #23
	sli.4s	v4, v0, #9
	add	x8, x0, #424
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #428
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v1, v5, v1
	add.4s	v0, v1, v0
	ushr.4s	v1, v0, #3
	sli.4s	v1, v0, #29
	add	x8, x0, #416
	ld1r.4s	{ v0 }, [x8]
	add	x8, x0, #420
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v0, v0, v3
	eor.16b	v2, v5, v2
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #5
	sli.4s	v2, v0, #27
	add	x8, x0, #408
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v4
	add	x8, x0, #412
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v3, v5, v3
	add.4s	v0, v3, v0
	ushr.4s	v3, v0, #23
	add	x8, x0, #448
	ld1r.4s	{ v5 }, [x8]
	sli.4s	v3, v0, #9
	eor.16b	v0, v5, v1
	add	x8, x0, #452
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v4, v5, v4
	add.4s	v0, v4, v0
	ushr.4s	v4, v0, #3
	sli.4s	v4, v0, #29
	add	x8, x0, #440
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #444
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v1, v5, v1
	add.4s	v0, v1, v0
	ushr.4s	v1, v0, #5
	sli.4s	v1, v0, #27
	add	x8, x0, #432
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v3
	add	x8, x0, #436
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v2, v5, v2
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #23
	sli.4s	v2, v0, #9
	add	x8, x0, #472
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v4
	add	x8, x0, #476
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v3, v5, v3
	add.4s	v0, v3, v0
	ushr.4s	v3, v0, #3
	sli.4s	v3, v0, #29
	add	x8, x0, #464
	ld1r.4s	{ v0 }, [x8]
	add	x8, x0, #468
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v0, v0, v1
	eor.16b	v4, v5, v4
	add.4s	v0, v4, v0
	ushr.4s	v4, v0, #5
	sli.4s	v4, v0, #27
	add	x8, x0, #456
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #460
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v1, v5, v1
	add.4s	v0, v1, v0
	ushr.4s	v1, v0, #23
	add	x8, x0, #496
	ld1r.4s	{ v5 }, [x8]
	sli.4s	v1, v0, #9
	eor.16b	v0, v5, v3
	add	x8, x0, #500
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v2, v5, v2
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #3
	sli.4s	v2, v0, #29
	add	x8, x0, #488
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v4
	add	x8, x0, #492
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v3, v5, v3
	add.4s	v0, v3, v0
	ushr.4s	v3, v0, #5
	sli.4s	v3, v0, #27
	add	x8, x0, #480
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v1
	add	x8, x0, #484
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v4, v5, v4
	add.4s	v0, v4, v0
	ushr.4s	v4, v0, #23
	sli.4s	v4, v0, #9
	add	x8, x0, #520
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #524
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v1, v5, v1
	add.4s	v0, v1, v0
	ushr.4s	v1, v0, #3
	sli.4s	v1, v0, #29
	add	x8, x0, #512
	ld1r.4s	{ v0 }, [x8]
	add	x8, x0, #516
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v0, v0, v3
	eor.16b	v2, v5, v2
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #5
	sli.4s	v2, v0, #27
	add	x8, x0, #504
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v4
	add	x8, x0, #508
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v3, v5, v3
	add.4s	v0, v3, v0
	ushr.4s	v5, v0, #23
	add	x8, x0, #544
	ld1r.4s	{ v3 }, [x8]
	sli.4s	v5, v0, #9
	eor.16b	v0, v3, v1
	add	x8, x0, #548
	ld1r.4s	{ v3 }, [x8]
	eor.16b	v3, v3, v4
	add.4s	v0, v3, v0
	ushr.4s	v4, v0, #3
	sli.4s	v4, v0, #29
	add	x8, x0, #536
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #540
	ld1r.4s	{ v3 }, [x8]
	eor.16b	v1, v3, v1
	add.4s	v0, v1, v0
	ushr.4s	v6, v0, #5
	sli.4s	v6, v0, #27
	add	x8, x0, #528
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v5
	add	x8, x0, #532
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v1, v1, v2
	add.4s	v7, v1, v0
	ushr.4s	v3, v7, #23
	sli.4s	v3, v7, #9
	add	x8, x0, #568
	ld1r.4s	{ v7 }, [x8]
	eor.16b	v7, v7, v4
	add	x8, x0, #572
	ld1r.4s	{ v16 }, [x8]
	eor.16b	v5, v16, v5
	add.4s	v5, v5, v7
	ushr.4s	v2, v5, #3
	sli.4s	v2, v5, #29
	add	x8, x0, #560
	ld1r.4s	{ v5 }, [x8]
	add	x8, x0, #564
	ld1r.4s	{ v7 }, [x8]
	eor.16b	v5, v5, v6
	eor.16b	v4, v7, v4
	add.4s	v4, v4, v5
	ushr.4s	v1, v4, #5
	sli.4s	v1, v4, #27
	add	x8, x0, #552
	ld1r.4s	{ v4 }, [x8]
	eor.16b	v4, v4, v3
	add	x8, x0, #556
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v5, v5, v6
	add.4s	v4, v5, v4
	ushr.4s	v0, v4, #23
	sli.4s	v0, v4, #9
	ldr	w8, [x0, #768]
	cmp	w8, #25
	b.lo	LBB0_3
; %bb.1:
	add	x9, x0, #592
	ld1r.4s	{ v4 }, [x9]
	eor.16b	v4, v4, v2
	add	x9, x0, #596
	ld1r.4s	{ v5 }, [x9]
	eor.16b	v5, v5, v3
	add.4s	v4, v5, v4
	ushr.4s	v5, v4, #3
	add	x9, x0, #584
	ld1r.4s	{ v6 }, [x9]
	sli.4s	v5, v4, #29
	eor.16b	v4, v6, v1
	add	x9, x0, #588
	ld1r.4s	{ v6 }, [x9]
	eor.16b	v6, v6, v2
	add.4s	v4, v6, v4
	ushr.4s	v6, v4, #5
	sli.4s	v6, v4, #27
	add	x9, x0, #576
	ld1r.4s	{ v4 }, [x9]
	eor.16b	v4, v4, v0
	add	x9, x0, #580
	ld1r.4s	{ v7 }, [x9]
	eor.16b	v7, v7, v1
	add.4s	v4, v7, v4
	ushr.4s	v7, v4, #23
	sli.4s	v7, v4, #9
	add	x9, x0, #616
	ld1r.4s	{ v4 }, [x9]
	eor.16b	v4, v4, v5
	add	x9, x0, #620
	ld1r.4s	{ v16 }, [x9]
	eor.16b	v0, v16, v0
	add.4s	v0, v0, v4
	ushr.4s	v1, v0, #3
	sli.4s	v1, v0, #29
	add	x9, x0, #608
	ld1r.4s	{ v0 }, [x9]
	eor.16b	v0, v0, v6
	add	x9, x0, #612
	ld1r.4s	{ v2 }, [x9]
	eor.16b	v2, v2, v5
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #5
	sli.4s	v2, v0, #27
	add	x9, x0, #600
	ld1r.4s	{ v0 }, [x9]
	add	x9, x0, #604
	ld1r.4s	{ v3 }, [x9]
	eor.16b	v0, v0, v7
	eor.16b	v3, v3, v6
	add.4s	v0, v3, v0
	ushr.4s	v4, v0, #23
	sli.4s	v4, v0, #9
	add	x9, x0, #640
	ld1r.4s	{ v0 }, [x9]
	eor.16b	v0, v0, v1
	add	x9, x0, #644
	ld1r.4s	{ v3 }, [x9]
	eor.16b	v3, v3, v7
	add.4s	v0, v3, v0
	ushr.4s	v5, v0, #3
	add	x9, x0, #632
	ld1r.4s	{ v3 }, [x9]
	sli.4s	v5, v0, #29
	eor.16b	v0, v3, v2
	add	x9, x0, #636
	ld1r.4s	{ v3 }, [x9]
	eor.16b	v1, v3, v1
	add.4s	v0, v1, v0
	ushr.4s	v6, v0, #5
	sli.4s	v6, v0, #27
	add	x9, x0, #624
	ld1r.4s	{ v0 }, [x9]
	eor.16b	v0, v0, v4
	add	x9, x0, #628
	ld1r.4s	{ v1 }, [x9]
	eor.16b	v1, v1, v2
	add.4s	v7, v1, v0
	ushr.4s	v3, v7, #23
	sli.4s	v3, v7, #9
	add	x9, x0, #664
	ld1r.4s	{ v7 }, [x9]
	eor.16b	v7, v7, v5
	add	x9, x0, #668
	ld1r.4s	{ v16 }, [x9]
	eor.16b	v4, v16, v4
	add.4s	v4, v4, v7
	ushr.4s	v2, v4, #3
	sli.4s	v2, v4, #29
	add	x9, x0, #656
	ld1r.4s	{ v4 }, [x9]
	eor.16b	v4, v4, v6
	add	x9, x0, #660
	ld1r.4s	{ v7 }, [x9]
	eor.16b	v5, v7, v5
	add.4s	v4, v5, v4
	ushr.4s	v1, v4, #5
	sli.4s	v1, v4, #27
	add	x9, x0, #648
	ld1r.4s	{ v4 }, [x9]
	add	x9, x0, #652
	ld1r.4s	{ v5 }, [x9]
	eor.16b	v4, v4, v3
	eor.16b	v5, v5, v6
	add.4s	v4, v5, v4
	ushr.4s	v0, v4, #23
	sli.4s	v0, v4, #9
	cmp	w8, #29
	b.lo	LBB0_3
; %bb.2:
	add	x8, x0, #688
	ld1r.4s	{ v4 }, [x8]
	add	x8, x0, #692
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v4, v4, v2
	eor.16b	v5, v5, v3
	add.4s	v4, v5, v4
	ushr.4s	v5, v4, #3
	sli.4s	v5, v4, #29
	add	x8, x0, #680
	ld1r.4s	{ v4 }, [x8]
	eor.16b	v4, v4, v1
	add	x8, x0, #684
	ld1r.4s	{ v6 }, [x8]
	eor.16b	v6, v6, v2
	add.4s	v4, v6, v4
	ushr.4s	v6, v4, #5
	add	x8, x0, #672
	ld1r.4s	{ v7 }, [x8]
	sli.4s	v6, v4, #27
	eor.16b	v4, v7, v0
	add	x8, x0, #676
	ld1r.4s	{ v7 }, [x8]
	eor.16b	v7, v7, v1
	add.4s	v4, v7, v4
	ushr.4s	v7, v4, #23
	sli.4s	v7, v4, #9
	add	x8, x0, #712
	ld1r.4s	{ v4 }, [x8]
	eor.16b	v4, v4, v5
	add	x8, x0, #716
	ld1r.4s	{ v16 }, [x8]
	eor.16b	v0, v16, v0
	add.4s	v0, v0, v4
	ushr.4s	v1, v0, #3
	sli.4s	v1, v0, #29
	add	x8, x0, #704
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v6
	add	x8, x0, #708
	ld1r.4s	{ v2 }, [x8]
	eor.16b	v2, v2, v5
	add.4s	v0, v2, v0
	ushr.4s	v2, v0, #5
	sli.4s	v2, v0, #27
	add	x8, x0, #696
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v7
	add	x8, x0, #700
	ld1r.4s	{ v3 }, [x8]
	eor.16b	v3, v3, v6
	add.4s	v0, v3, v0
	ushr.4s	v4, v0, #23
	sli.4s	v4, v0, #9
	add	x8, x0, #736
	ld1r.4s	{ v0 }, [x8]
	add	x8, x0, #740
	ld1r.4s	{ v3 }, [x8]
	eor.16b	v0, v0, v1
	eor.16b	v3, v3, v7
	add.4s	v0, v3, v0
	ushr.4s	v5, v0, #3
	sli.4s	v5, v0, #29
	add	x8, x0, #728
	ld1r.4s	{ v0 }, [x8]
	eor.16b	v0, v0, v2
	add	x8, x0, #732
	ld1r.4s	{ v3 }, [x8]
	eor.16b	v1, v3, v1
	add.4s	v0, v1, v0
	ushr.4s	v6, v0, #5
	add	x8, x0, #720
	ld1r.4s	{ v1 }, [x8]
	sli.4s	v6, v0, #27
	eor.16b	v0, v1, v4
	add	x8, x0, #724
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v1, v1, v2
	add.4s	v7, v1, v0
	ushr.4s	v3, v7, #23
	sli.4s	v3, v7, #9
	add	x8, x0, #760
	ld1r.4s	{ v7 }, [x8]
	eor.16b	v7, v7, v5
	add	x8, x0, #764
	ld1r.4s	{ v16 }, [x8]
	eor.16b	v4, v16, v4
	add.4s	v4, v4, v7
	ushr.4s	v2, v4, #3
	sli.4s	v2, v4, #29
	add	x8, x0, #752
	ld1r.4s	{ v4 }, [x8]
	eor.16b	v4, v4, v6
	add	x8, x0, #756
	ld1r.4s	{ v7 }, [x8]
	eor.16b	v5, v7, v5
	add.4s	v4, v5, v4
	ushr.4s	v1, v4, #5
	sli.4s	v1, v4, #27
	add	x8, x0, #744
	ld1r.4s	{ v4 }, [x8]
	eor.16b	v4, v4, v3
	add	x8, x0, #748
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v5, v5, v6
	add.4s	v4, v5, v4
	ushr.4s	v0, v4, #23
	sli.4s	v0, v4, #9
LBB0_3:
	st4.4s	{ v0, v1, v2, v3 }, [x1]
	ldp	x29, x30, [sp], #16             ; 16-byte Folded Reload
	ret
                                        ; -- End function
	.globl	_lea_decrypt_4block             ; -- Begin function lea_decrypt_4block
	.p2align	2
_lea_decrypt_4block:                    ; @lea_decrypt_4block
; %bb.0:
	stp	x29, x30, [sp, #-16]!           ; 16-byte Folded Spill
	mov	x29, sp
	ld4.4s	{ v0, v1, v2, v3 }, [x2]
	ldr	w8, [x0, #768]
	cmp	w8, #29
	b.lo	LBB1_2
; %bb.1:
	ushr.4s	v4, v0, #9
	sli.4s	v4, v0, #23
	add	x8, x0, #744
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v5, v5, v3
	sub.4s	v4, v4, v5
	add	x8, x0, #748
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v6, v4, v5
	ushr.4s	v7, v1, #27
	sli.4s	v7, v1, #5
	add	x8, x0, #752
	ld1r.4s	{ v16 }, [x8]
	eor3.16b	v4, v4, v5, v16
	sub.4s	v4, v7, v4
	add	x8, x0, #756
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v7, v4, v5
	ushr.4s	v16, v2, #29
	sli.4s	v16, v2, #3
	add	x8, x0, #760
	ld1r.4s	{ v17 }, [x8]
	eor3.16b	v4, v4, v5, v17
	sub.4s	v4, v16, v4
	add	x8, x0, #764
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v16, v4, v5
	ushr.4s	v17, v3, #9
	sli.4s	v17, v3, #23
	add	x8, x0, #720
	ld1r.4s	{ v0 }, [x8]
	eor3.16b	v0, v4, v5, v0
	add	x8, x0, #724
	ld1r.4s	{ v1 }, [x8]
	sub.4s	v0, v17, v0
	eor.16b	v2, v0, v1
	ushr.4s	v3, v6, #27
	add	x8, x0, #728
	ld1r.4s	{ v4 }, [x8]
	sli.4s	v3, v6, #5
	eor3.16b	v0, v0, v1, v4
	sub.4s	v0, v3, v0
	add	x8, x0, #732
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v4, v7, #29
	sli.4s	v4, v7, #3
	add	x8, x0, #736
	ld1r.4s	{ v5 }, [x8]
	eor3.16b	v0, v0, v1, v5
	sub.4s	v0, v4, v0
	add	x8, x0, #740
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v4, v0, v1
	ushr.4s	v5, v16, #9
	sli.4s	v5, v16, #23
	add	x8, x0, #696
	ld1r.4s	{ v6 }, [x8]
	eor3.16b	v0, v0, v1, v6
	sub.4s	v0, v5, v0
	add	x8, x0, #700
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #27
	sli.4s	v6, v2, #5
	add	x8, x0, #704
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	sub.4s	v0, v6, v0
	add	x8, x0, #708
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v6, v0, v1
	ushr.4s	v2, v3, #29
	sli.4s	v2, v3, #3
	add	x8, x0, #712
	ld1r.4s	{ v3 }, [x8]
	eor3.16b	v0, v0, v1, v3
	sub.4s	v7, v2, v0
	add	x8, x0, #716
	ld1r.4s	{ v16 }, [x8]
	eor.16b	v0, v7, v16
	ushr.4s	v17, v4, #9
	sli.4s	v17, v4, #23
	add	x8, x0, #672
	ld1r.4s	{ v4 }, [x8]
	eor3.16b	v4, v7, v16, v4
	add	x8, x0, #676
	ld1r.4s	{ v7 }, [x8]
	sub.4s	v4, v17, v4
	eor.16b	v1, v4, v7
	ushr.4s	v16, v5, #27
	add	x8, x0, #680
	ld1r.4s	{ v17 }, [x8]
	sli.4s	v16, v5, #5
	eor3.16b	v4, v4, v7, v17
	sub.4s	v4, v16, v4
	add	x8, x0, #684
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v2, v4, v5
	ushr.4s	v7, v6, #29
	sli.4s	v7, v6, #3
	add	x8, x0, #688
	ld1r.4s	{ v6 }, [x8]
	eor3.16b	v4, v4, v5, v6
	sub.4s	v4, v7, v4
	add	x8, x0, #692
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v3, v4, v5
	b	LBB1_3
LBB1_2:
	cmp	w8, #25
	b.lo	LBB1_4
LBB1_3:
	ushr.4s	v4, v0, #9
	sli.4s	v4, v0, #23
	add	x8, x0, #648
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v5, v5, v3
	sub.4s	v4, v4, v5
	add	x8, x0, #652
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v6, v4, v5
	ushr.4s	v7, v1, #27
	sli.4s	v7, v1, #5
	add	x8, x0, #656
	ld1r.4s	{ v16 }, [x8]
	eor3.16b	v4, v4, v5, v16
	sub.4s	v4, v7, v4
	add	x8, x0, #660
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v7, v4, v5
	ushr.4s	v16, v2, #29
	sli.4s	v16, v2, #3
	add	x8, x0, #664
	ld1r.4s	{ v17 }, [x8]
	eor3.16b	v4, v4, v5, v17
	sub.4s	v4, v16, v4
	add	x8, x0, #668
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v16, v4, v5
	ushr.4s	v17, v3, #9
	sli.4s	v17, v3, #23
	add	x8, x0, #624
	ld1r.4s	{ v0 }, [x8]
	eor3.16b	v0, v4, v5, v0
	add	x8, x0, #628
	ld1r.4s	{ v1 }, [x8]
	sub.4s	v0, v17, v0
	eor.16b	v2, v0, v1
	ushr.4s	v3, v6, #27
	add	x8, x0, #632
	ld1r.4s	{ v4 }, [x8]
	sli.4s	v3, v6, #5
	eor3.16b	v0, v0, v1, v4
	sub.4s	v0, v3, v0
	add	x8, x0, #636
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v4, v7, #29
	sli.4s	v4, v7, #3
	add	x8, x0, #640
	ld1r.4s	{ v5 }, [x8]
	eor3.16b	v0, v0, v1, v5
	sub.4s	v0, v4, v0
	add	x8, x0, #644
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v4, v0, v1
	ushr.4s	v5, v16, #9
	sli.4s	v5, v16, #23
	add	x8, x0, #600
	ld1r.4s	{ v6 }, [x8]
	eor3.16b	v0, v0, v1, v6
	sub.4s	v0, v5, v0
	add	x8, x0, #604
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #27
	sli.4s	v6, v2, #5
	add	x8, x0, #608
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	sub.4s	v0, v6, v0
	add	x8, x0, #612
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v6, v0, v1
	ushr.4s	v2, v3, #29
	sli.4s	v2, v3, #3
	add	x8, x0, #616
	ld1r.4s	{ v3 }, [x8]
	eor3.16b	v0, v0, v1, v3
	sub.4s	v7, v2, v0
	add	x8, x0, #620
	ld1r.4s	{ v16 }, [x8]
	eor.16b	v0, v7, v16
	ushr.4s	v17, v4, #9
	sli.4s	v17, v4, #23
	add	x8, x0, #576
	ld1r.4s	{ v4 }, [x8]
	eor3.16b	v4, v7, v16, v4
	add	x8, x0, #580
	ld1r.4s	{ v7 }, [x8]
	sub.4s	v4, v17, v4
	eor.16b	v1, v4, v7
	ushr.4s	v16, v5, #27
	add	x8, x0, #584
	ld1r.4s	{ v17 }, [x8]
	sli.4s	v16, v5, #5
	eor3.16b	v4, v4, v7, v17
	sub.4s	v4, v16, v4
	add	x8, x0, #588
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v2, v4, v5
	ushr.4s	v7, v6, #29
	sli.4s	v7, v6, #3
	add	x8, x0, #592
	ld1r.4s	{ v6 }, [x8]
	eor3.16b	v4, v4, v5, v6
	sub.4s	v4, v7, v4
	add	x8, x0, #596
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v3, v4, v5
LBB1_4:
	ushr.4s	v4, v0, #9
	sli.4s	v4, v0, #23
	add	x8, x0, #552
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v5, v5, v3
	sub.4s	v4, v4, v5
	add	x8, x0, #556
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v6, v4, v5
	ushr.4s	v7, v1, #27
	sli.4s	v7, v1, #5
	add	x8, x0, #560
	ld1r.4s	{ v16 }, [x8]
	eor3.16b	v4, v4, v5, v16
	sub.4s	v4, v7, v4
	add	x8, x0, #564
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v7, v4, v5
	ushr.4s	v16, v2, #29
	sli.4s	v16, v2, #3
	add	x8, x0, #568
	ld1r.4s	{ v17 }, [x8]
	eor3.16b	v4, v4, v5, v17
	sub.4s	v4, v16, v4
	add	x8, x0, #572
	ld1r.4s	{ v5 }, [x8]
	eor.16b	v16, v4, v5
	ushr.4s	v17, v3, #9
	sli.4s	v17, v3, #23
	add	x8, x0, #528
	ld1r.4s	{ v0 }, [x8]
	eor3.16b	v0, v4, v5, v0
	add	x8, x0, #532
	ld1r.4s	{ v1 }, [x8]
	sub.4s	v0, v17, v0
	eor.16b	v2, v0, v1
	ushr.4s	v3, v6, #27
	add	x8, x0, #536
	ld1r.4s	{ v4 }, [x8]
	sli.4s	v3, v6, #5
	eor3.16b	v0, v0, v1, v4
	sub.4s	v0, v3, v0
	add	x8, x0, #540
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v4, v7, #29
	sli.4s	v4, v7, #3
	add	x8, x0, #544
	ld1r.4s	{ v5 }, [x8]
	eor3.16b	v0, v0, v1, v5
	sub.4s	v0, v4, v0
	add	x8, x0, #548
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v4, v0, v1
	ushr.4s	v5, v16, #9
	sli.4s	v5, v16, #23
	add	x8, x0, #504
	ld1r.4s	{ v6 }, [x8]
	eor3.16b	v0, v0, v1, v6
	sub.4s	v0, v5, v0
	add	x8, x0, #508
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #27
	sli.4s	v6, v2, #5
	add	x8, x0, #512
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	sub.4s	v0, v6, v0
	add	x8, x0, #516
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v2, v0, v1
	ushr.4s	v6, v3, #29
	sli.4s	v6, v3, #3
	add	x8, x0, #520
	ld1r.4s	{ v3 }, [x8]
	eor3.16b	v0, v0, v1, v3
	sub.4s	v0, v6, v0
	add	x8, x0, #524
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v6, v4, #9
	sli.4s	v6, v4, #23
	add	x8, x0, #480
	ld1r.4s	{ v4 }, [x8]
	eor3.16b	v0, v0, v1, v4
	add	x8, x0, #484
	ld1r.4s	{ v1 }, [x8]
	sub.4s	v0, v6, v0
	eor.16b	v4, v0, v1
	ushr.4s	v6, v5, #27
	add	x8, x0, #488
	ld1r.4s	{ v7 }, [x8]
	sli.4s	v6, v5, #5
	eor3.16b	v0, v0, v1, v7
	sub.4s	v0, v6, v0
	add	x8, x0, #492
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #29
	sli.4s	v6, v2, #3
	add	x8, x0, #496
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	sub.4s	v0, v6, v0
	add	x8, x0, #500
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v2, v0, v1
	ushr.4s	v6, v3, #9
	sli.4s	v6, v3, #23
	add	x8, x0, #456
	ld1r.4s	{ v3 }, [x8]
	eor3.16b	v0, v0, v1, v3
	sub.4s	v0, v6, v0
	add	x8, x0, #460
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v6, v4, #27
	sli.4s	v6, v4, #5
	add	x8, x0, #464
	ld1r.4s	{ v4 }, [x8]
	eor3.16b	v0, v0, v1, v4
	sub.4s	v0, v6, v0
	add	x8, x0, #468
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v4, v0, v1
	ushr.4s	v6, v5, #29
	sli.4s	v6, v5, #3
	add	x8, x0, #472
	ld1r.4s	{ v5 }, [x8]
	eor3.16b	v0, v0, v1, v5
	sub.4s	v0, v6, v0
	add	x8, x0, #476
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #9
	sli.4s	v6, v2, #23
	add	x8, x0, #432
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	add	x8, x0, #436
	ld1r.4s	{ v1 }, [x8]
	sub.4s	v0, v6, v0
	eor.16b	v2, v0, v1
	ushr.4s	v6, v3, #27
	add	x8, x0, #440
	ld1r.4s	{ v7 }, [x8]
	sli.4s	v6, v3, #5
	eor3.16b	v0, v0, v1, v7
	sub.4s	v0, v6, v0
	add	x8, x0, #444
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v6, v4, #29
	sli.4s	v6, v4, #3
	add	x8, x0, #448
	ld1r.4s	{ v4 }, [x8]
	eor3.16b	v0, v0, v1, v4
	sub.4s	v0, v6, v0
	add	x8, x0, #452
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v4, v0, v1
	ushr.4s	v6, v5, #9
	sli.4s	v6, v5, #23
	add	x8, x0, #408
	ld1r.4s	{ v5 }, [x8]
	eor3.16b	v0, v0, v1, v5
	sub.4s	v0, v6, v0
	add	x8, x0, #412
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #27
	sli.4s	v6, v2, #5
	add	x8, x0, #416
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	sub.4s	v0, v6, v0
	add	x8, x0, #420
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v2, v0, v1
	ushr.4s	v6, v3, #29
	sli.4s	v6, v3, #3
	add	x8, x0, #424
	ld1r.4s	{ v3 }, [x8]
	eor3.16b	v0, v0, v1, v3
	sub.4s	v0, v6, v0
	add	x8, x0, #428
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v6, v4, #9
	sli.4s	v6, v4, #23
	add	x8, x0, #384
	ld1r.4s	{ v4 }, [x8]
	eor3.16b	v0, v0, v1, v4
	add	x8, x0, #388
	ld1r.4s	{ v1 }, [x8]
	sub.4s	v0, v6, v0
	eor.16b	v4, v0, v1
	ushr.4s	v6, v5, #27
	add	x8, x0, #392
	ld1r.4s	{ v7 }, [x8]
	sli.4s	v6, v5, #5
	eor3.16b	v0, v0, v1, v7
	sub.4s	v0, v6, v0
	add	x8, x0, #396
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #29
	sli.4s	v6, v2, #3
	add	x8, x0, #400
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	sub.4s	v0, v6, v0
	add	x8, x0, #404
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v2, v0, v1
	ushr.4s	v6, v3, #9
	sli.4s	v6, v3, #23
	add	x8, x0, #360
	ld1r.4s	{ v3 }, [x8]
	eor3.16b	v0, v0, v1, v3
	sub.4s	v0, v6, v0
	add	x8, x0, #364
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v6, v4, #27
	sli.4s	v6, v4, #5
	add	x8, x0, #368
	ld1r.4s	{ v4 }, [x8]
	eor3.16b	v0, v0, v1, v4
	sub.4s	v0, v6, v0
	add	x8, x0, #372
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v4, v0, v1
	ushr.4s	v6, v5, #29
	sli.4s	v6, v5, #3
	add	x8, x0, #376
	ld1r.4s	{ v5 }, [x8]
	eor3.16b	v0, v0, v1, v5
	sub.4s	v0, v6, v0
	add	x8, x0, #380
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #9
	sli.4s	v6, v2, #23
	add	x8, x0, #336
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	add	x8, x0, #340
	ld1r.4s	{ v1 }, [x8]
	sub.4s	v0, v6, v0
	eor.16b	v2, v0, v1
	ushr.4s	v6, v3, #27
	add	x8, x0, #344
	ld1r.4s	{ v7 }, [x8]
	sli.4s	v6, v3, #5
	eor3.16b	v0, v0, v1, v7
	sub.4s	v0, v6, v0
	add	x8, x0, #348
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v6, v4, #29
	sli.4s	v6, v4, #3
	add	x8, x0, #352
	ld1r.4s	{ v4 }, [x8]
	eor3.16b	v0, v0, v1, v4
	sub.4s	v0, v6, v0
	add	x8, x0, #356
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v4, v0, v1
	ushr.4s	v6, v5, #9
	sli.4s	v6, v5, #23
	add	x8, x0, #312
	ld1r.4s	{ v5 }, [x8]
	eor3.16b	v0, v0, v1, v5
	sub.4s	v0, v6, v0
	add	x8, x0, #316
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #27
	sli.4s	v6, v2, #5
	add	x8, x0, #320
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	sub.4s	v0, v6, v0
	add	x8, x0, #324
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v2, v0, v1
	ushr.4s	v6, v3, #29
	sli.4s	v6, v3, #3
	add	x8, x0, #328
	ld1r.4s	{ v3 }, [x8]
	eor3.16b	v0, v0, v1, v3
	sub.4s	v0, v6, v0
	add	x8, x0, #332
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v6, v4, #9
	sli.4s	v6, v4, #23
	add	x8, x0, #288
	ld1r.4s	{ v4 }, [x8]
	eor3.16b	v0, v0, v1, v4
	add	x8, x0, #292
	ld1r.4s	{ v1 }, [x8]
	sub.4s	v0, v6, v0
	eor.16b	v4, v0, v1
	ushr.4s	v6, v5, #27
	add	x8, x0, #296
	ld1r.4s	{ v7 }, [x8]
	sli.4s	v6, v5, #5
	eor3.16b	v0, v0, v1, v7
	sub.4s	v0, v6, v0
	add	x8, x0, #300
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #29
	sli.4s	v6, v2, #3
	add	x8, x0, #304
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	sub.4s	v0, v6, v0
	add	x8, x0, #308
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v2, v0, v1
	ushr.4s	v6, v3, #9
	sli.4s	v6, v3, #23
	add	x8, x0, #264
	ld1r.4s	{ v3 }, [x8]
	eor3.16b	v0, v0, v1, v3
	sub.4s	v0, v6, v0
	add	x8, x0, #268
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v6, v4, #27
	sli.4s	v6, v4, #5
	add	x8, x0, #272
	ld1r.4s	{ v4 }, [x8]
	eor3.16b	v0, v0, v1, v4
	sub.4s	v0, v6, v0
	add	x8, x0, #276
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v4, v0, v1
	ushr.4s	v6, v5, #29
	sli.4s	v6, v5, #3
	add	x8, x0, #280
	ld1r.4s	{ v5 }, [x8]
	eor3.16b	v0, v0, v1, v5
	sub.4s	v0, v6, v0
	add	x8, x0, #284
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #9
	sli.4s	v6, v2, #23
	add	x8, x0, #240
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	add	x8, x0, #244
	ld1r.4s	{ v1 }, [x8]
	sub.4s	v0, v6, v0
	eor.16b	v2, v0, v1
	ushr.4s	v6, v3, #27
	add	x8, x0, #248
	ld1r.4s	{ v7 }, [x8]
	sli.4s	v6, v3, #5
	eor3.16b	v0, v0, v1, v7
	sub.4s	v0, v6, v0
	add	x8, x0, #252
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v6, v4, #29
	sli.4s	v6, v4, #3
	add	x8, x0, #256
	ld1r.4s	{ v4 }, [x8]
	eor3.16b	v0, v0, v1, v4
	sub.4s	v0, v6, v0
	add	x8, x0, #260
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v4, v0, v1
	ushr.4s	v6, v5, #9
	sli.4s	v6, v5, #23
	add	x8, x0, #216
	ld1r.4s	{ v5 }, [x8]
	eor3.16b	v0, v0, v1, v5
	sub.4s	v0, v6, v0
	add	x8, x0, #220
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #27
	sli.4s	v6, v2, #5
	add	x8, x0, #224
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	sub.4s	v0, v6, v0
	add	x8, x0, #228
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v2, v0, v1
	ushr.4s	v6, v3, #29
	sli.4s	v6, v3, #3
	add	x8, x0, #232
	ld1r.4s	{ v3 }, [x8]
	eor3.16b	v0, v0, v1, v3
	sub.4s	v0, v6, v0
	add	x8, x0, #236
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v6, v4, #9
	sli.4s	v6, v4, #23
	add	x8, x0, #192
	ld1r.4s	{ v4 }, [x8]
	eor3.16b	v0, v0, v1, v4
	add	x8, x0, #196
	ld1r.4s	{ v1 }, [x8]
	sub.4s	v0, v6, v0
	eor.16b	v4, v0, v1
	ushr.4s	v6, v5, #27
	add	x8, x0, #200
	ld1r.4s	{ v7 }, [x8]
	sli.4s	v6, v5, #5
	eor3.16b	v0, v0, v1, v7
	sub.4s	v0, v6, v0
	add	x8, x0, #204
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #29
	sli.4s	v6, v2, #3
	add	x8, x0, #208
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	sub.4s	v0, v6, v0
	add	x8, x0, #212
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v2, v0, v1
	ushr.4s	v6, v3, #9
	sli.4s	v6, v3, #23
	add	x8, x0, #168
	ld1r.4s	{ v3 }, [x8]
	eor3.16b	v0, v0, v1, v3
	sub.4s	v0, v6, v0
	add	x8, x0, #172
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v6, v4, #27
	sli.4s	v6, v4, #5
	add	x8, x0, #176
	ld1r.4s	{ v4 }, [x8]
	eor3.16b	v0, v0, v1, v4
	sub.4s	v0, v6, v0
	add	x8, x0, #180
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v4, v0, v1
	ushr.4s	v6, v5, #29
	sli.4s	v6, v5, #3
	add	x8, x0, #184
	ld1r.4s	{ v5 }, [x8]
	eor3.16b	v0, v0, v1, v5
	sub.4s	v0, v6, v0
	add	x8, x0, #188
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #9
	sli.4s	v6, v2, #23
	add	x8, x0, #144
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	add	x8, x0, #148
	ld1r.4s	{ v1 }, [x8]
	sub.4s	v0, v6, v0
	eor.16b	v2, v0, v1
	ushr.4s	v6, v3, #27
	add	x8, x0, #152
	ld1r.4s	{ v7 }, [x8]
	sli.4s	v6, v3, #5
	eor3.16b	v0, v0, v1, v7
	sub.4s	v0, v6, v0
	add	x8, x0, #156
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v6, v4, #29
	sli.4s	v6, v4, #3
	add	x8, x0, #160
	ld1r.4s	{ v4 }, [x8]
	eor3.16b	v0, v0, v1, v4
	sub.4s	v0, v6, v0
	add	x8, x0, #164
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v4, v0, v1
	ushr.4s	v6, v5, #9
	sli.4s	v6, v5, #23
	add	x8, x0, #120
	ld1r.4s	{ v5 }, [x8]
	eor3.16b	v0, v0, v1, v5
	sub.4s	v0, v6, v0
	add	x8, x0, #124
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #27
	sli.4s	v6, v2, #5
	add	x8, x0, #128
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	sub.4s	v0, v6, v0
	add	x8, x0, #132
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v2, v0, v1
	ushr.4s	v6, v3, #29
	sli.4s	v6, v3, #3
	add	x8, x0, #136
	ld1r.4s	{ v3 }, [x8]
	eor3.16b	v0, v0, v1, v3
	sub.4s	v0, v6, v0
	add	x8, x0, #140
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v6, v4, #9
	sli.4s	v6, v4, #23
	add	x8, x0, #96
	ld1r.4s	{ v4 }, [x8]
	eor3.16b	v0, v0, v1, v4
	add	x8, x0, #100
	ld1r.4s	{ v1 }, [x8]
	sub.4s	v0, v6, v0
	eor.16b	v4, v0, v1
	ushr.4s	v6, v5, #27
	add	x8, x0, #104
	ld1r.4s	{ v7 }, [x8]
	sli.4s	v6, v5, #5
	eor3.16b	v0, v0, v1, v7
	sub.4s	v0, v6, v0
	add	x8, x0, #108
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #29
	sli.4s	v6, v2, #3
	add	x8, x0, #112
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	sub.4s	v0, v6, v0
	add	x8, x0, #116
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v2, v0, v1
	ushr.4s	v6, v3, #9
	sli.4s	v6, v3, #23
	add	x8, x0, #72
	ld1r.4s	{ v3 }, [x8]
	eor3.16b	v0, v0, v1, v3
	sub.4s	v0, v6, v0
	add	x8, x0, #76
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v6, v4, #27
	sli.4s	v6, v4, #5
	add	x8, x0, #80
	ld1r.4s	{ v4 }, [x8]
	eor3.16b	v0, v0, v1, v4
	sub.4s	v0, v6, v0
	add	x8, x0, #84
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v4, v0, v1
	ushr.4s	v6, v5, #29
	sli.4s	v6, v5, #3
	add	x8, x0, #88
	ld1r.4s	{ v5 }, [x8]
	eor3.16b	v0, v0, v1, v5
	sub.4s	v0, v6, v0
	add	x8, x0, #92
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #9
	sli.4s	v6, v2, #23
	add	x8, x0, #48
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	add	x8, x0, #52
	ld1r.4s	{ v1 }, [x8]
	sub.4s	v0, v6, v0
	eor.16b	v2, v0, v1
	ushr.4s	v6, v3, #27
	add	x8, x0, #56
	ld1r.4s	{ v7 }, [x8]
	sli.4s	v6, v3, #5
	eor3.16b	v0, v0, v1, v7
	sub.4s	v0, v6, v0
	add	x8, x0, #60
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v3, v0, v1
	ushr.4s	v6, v4, #29
	sli.4s	v6, v4, #3
	add	x8, x0, #64
	ld1r.4s	{ v4 }, [x8]
	eor3.16b	v0, v0, v1, v4
	sub.4s	v0, v6, v0
	add	x8, x0, #68
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v4, v0, v1
	ushr.4s	v6, v5, #9
	sli.4s	v6, v5, #23
	add	x8, x0, #24
	ld1r.4s	{ v5 }, [x8]
	eor3.16b	v0, v0, v1, v5
	sub.4s	v0, v6, v0
	add	x8, x0, #28
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v5, v0, v1
	ushr.4s	v6, v2, #27
	sli.4s	v6, v2, #5
	add	x8, x0, #32
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	sub.4s	v0, v6, v0
	add	x8, x0, #36
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v2, v0, v1
	ushr.4s	v6, v3, #29
	sli.4s	v6, v3, #3
	add	x8, x0, #40
	ld1r.4s	{ v3 }, [x8]
	eor3.16b	v0, v0, v1, v3
	sub.4s	v0, v6, v0
	add	x8, x0, #44
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v16, v0, v1
	ushr.4s	v3, v4, #9
	sli.4s	v3, v4, #23
	mov	x8, x0
	ld1r.4s	{ v4 }, [x8], #4
	eor3.16b	v0, v0, v1, v4
	ld1r.4s	{ v1 }, [x8]
	sub.4s	v0, v3, v0
	eor.16b	v17, v0, v1
	ushr.4s	v3, v5, #27
	add	x8, x0, #8
	ld1r.4s	{ v4 }, [x8]
	sli.4s	v3, v5, #5
	eor3.16b	v0, v0, v1, v4
	sub.4s	v0, v3, v0
	add	x8, x0, #12
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v18, v0, v1
	ushr.4s	v3, v2, #29
	sli.4s	v3, v2, #3
	add	x8, x0, #16
	ld1r.4s	{ v2 }, [x8]
	eor3.16b	v0, v0, v1, v2
	sub.4s	v0, v3, v0
	add	x8, x0, #20
	ld1r.4s	{ v1 }, [x8]
	eor.16b	v19, v0, v1
	st4.4s	{ v16, v17, v18, v19 }, [x1]
	ldp	x29, x30, [sp], #16             ; 16-byte Folded Reload
	ret
                                        ; -- End function
.subsections_via_symbols
