	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 14, 0
	.globl	_EncKeySetup                    ; -- Begin function EncKeySetup
	.p2align	2
_EncKeySetup:                           ; @EncKeySetup
; %bb.0:
	stp	x28, x27, [sp, #-96]!           ; 16-byte Folded Spill
	stp	x26, x25, [sp, #16]             ; 16-byte Folded Spill
	stp	x24, x23, [sp, #32]             ; 16-byte Folded Spill
	stp	x22, x21, [sp, #48]             ; 16-byte Folded Spill
	stp	x20, x19, [sp, #64]             ; 16-byte Folded Spill
	stp	x29, x30, [sp, #80]             ; 16-byte Folded Spill
	add	x29, sp, #80
	mov	w8, #-16711936                  ; =0xff00ff00
	ldp	w9, w10, [x1]
	rev	w15, w9
	rev	w13, w10
	ldp	w9, w10, [x1, #8]
	rev	w16, w9
	rev	w14, w10
	fmov	s0, w15
	mov.s	v0[1], w13
	mov	x9, #34359738352                ; =0x7fffffff0
	add	x3, x2, x9
	lsl	x9, x3, #29
	asr	x9, x9, #32
Lloh0:
	adrp	x4, _KRK@PAGE
Lloh1:
	add	x4, x4, _KRK@PAGEOFF
	add	x9, x4, x9, lsl #4
	ldp	w10, w11, [x9]
	eor	w17, w10, w15
	eor	w5, w11, w13
	ldp	w10, w9, [x9, #8]
	eor	w6, w10, w16
	eor	w7, w14, w9
	mov.s	v0[2], w16
Lloh2:
	adrp	x9, _S1@PAGE
Lloh3:
	add	x9, x9, _S1@PAGEOFF
	lsr	w10, w17, #24
	ldr	w11, [x9, w10, uxtw #2]
	mov.s	v0[3], w14
	ubfx	w12, w17, #16, #8
Lloh4:
	adrp	x10, _S2@PAGE
Lloh5:
	add	x10, x10, _S2@PAGEOFF
	ldr	w12, [x10, w12, uxtw #2]
	eor	w19, w12, w11
	ubfx	w12, w17, #8, #8
Lloh6:
	adrp	x11, _X1@PAGE
Lloh7:
	add	x11, x11, _X1@PAGEOFF
	ldr	w20, [x11, w12, uxtw #2]
Lloh8:
	adrp	x12, _X2@PAGE
Lloh9:
	add	x12, x12, _X2@PAGEOFF
	and	w17, w17, #0xff
	ldr	w17, [x12, w17, uxtw #2]
	eor	w17, w20, w17
	eor	w19, w19, w17
	lsr	w17, w5, #24
	ldr	w17, [x9, w17, uxtw #2]
	ubfx	w20, w5, #16, #8
	ldr	w20, [x10, w20, uxtw #2]
	eor	w17, w20, w17
	ubfx	w20, w5, #8, #8
	ldr	w20, [x11, w20, uxtw #2]
	and	w5, w5, #0xff
	ldr	w5, [x12, w5, uxtw #2]
	eor	w5, w20, w5
	eor	w5, w17, w5
	lsr	w17, w6, #24
	ldr	w17, [x9, w17, uxtw #2]
	ubfx	w20, w6, #16, #8
	ldr	w20, [x10, w20, uxtw #2]
	eor	w17, w20, w17
	ubfx	w20, w6, #8, #8
	ldr	w20, [x11, w20, uxtw #2]
	and	w6, w6, #0xff
	ldr	w6, [x12, w6, uxtw #2]
	eor	w6, w20, w6
	eor	w6, w17, w6
	lsr	w17, w7, #24
	ldr	w17, [x9, w17, uxtw #2]
	ubfx	w20, w7, #16, #8
	ldr	w20, [x10, w20, uxtw #2]
	eor	w20, w20, w17
	ubfx	w17, w7, #8, #8
	ldr	w21, [x11, w17, uxtw #2]
	and	w17, w7, #0xff
	ldr	w7, [x12, w17, uxtw #2]
	mov	w17, #16711935                  ; =0xff00ff
	eor	w7, w21, w7
	eor	w7, w20, w7
	eor	w5, w5, w6
	eor	w6, w7, w6
	eor	w20, w5, w19
	eor	w5, w7, w5
	eor	w7, w6, w20
	eor	w6, w6, w19
	rev	w6, w6
	ror	w21, w7, #16
	rev	w7, w5
	eor	w19, w21, w7
	eor	w5, w21, w6, ror #16
	eor	w6, w5, w20
	movi.2d	v1, #0000000000000000
	cmp	x2, #17
	b.lo	LBB0_4
; %bb.1:
	ldp	w21, w22, [x1, #16]
	rev	w21, w21
	rev	w22, w22
	fmov	s1, w21
	mov.s	v1[1], w22
	cmp	x2, #25
	b.lo	LBB0_3
; %bb.2:
	ldp	w21, w1, [x1, #24]
	rev	w21, w21
	mov.s	v1[2], w21
	rev	w1, w1
	mov.s	v1[3], w1
	b	LBB0_4
LBB0_3:
	movi.2d	v2, #0000000000000000
	mov.d	v1[1], v2[1]
LBB0_4:
	lsr	x21, x3, #3
	fmov	w1, s1
	mov.s	w22, v1[1]
	eor	w3, w1, w6
	eor	w1, w20, w22
	eor	w1, w1, w19
	cmp	w21, #2
	csinc	w20, wzr, w21, eq
	add	x21, x4, w20, sxtw #4
	ldp	w22, w23, [x21]
	eor	w22, w3, w22
	lsr	w24, w22, #24
	ldr	w24, [x11, w24, uxtw #2]
	eor	w25, w1, w23
	ubfx	w23, w22, #16, #8
	ldr	w23, [x12, w23, uxtw #2]
	eor	w23, w23, w24
	ubfx	w24, w22, #8, #8
	ldr	w24, [x9, w24, uxtw #2]
	and	w22, w22, #0xff
	ldr	w22, [x10, w22, uxtw #2]
	eor	w22, w24, w22
	lsr	w24, w25, #24
	ldr	w24, [x11, w24, uxtw #2]
	ubfx	w26, w25, #16, #8
	ldr	w26, [x12, w26, uxtw #2]
	eor	w23, w23, w22
	eor	w22, w26, w24
	ubfx	w24, w25, #8, #8
	ldr	w24, [x9, w24, uxtw #2]
	ldp	w26, w27, [x21, #8]
	and	w21, w25, #0xff
	ldr	w21, [x10, w21, uxtw #2]
	eor	w21, w24, w21
	eor	w24, w22, w21
	cmp	w20, #2
	csinc	w20, wzr, w20, eq
	add	x22, x4, w20, sxtw #4
	ldp	w4, w20, [x22]
	ldp	w21, w22, [x22, #8]
	fmov	s2, w19
	mov.s	v2[1], w7
	ext.16b	v1, v1, v1, #8
	eor.8b	v1, v2, v1
	fmov	s2, w6
	mov.s	v2[1], w5
	eor.8b	v1, v1, v2
	fmov	w5, s1
	eor	w5, w5, w26
	mov.s	w6, v1[1]
	eor	w6, w27, w6
	lsr	w7, w5, #24
	ldr	w7, [x11, w7, uxtw #2]
	ubfx	w19, w5, #16, #8
	ldr	w19, [x12, w19, uxtw #2]
	ubfx	w25, w5, #8, #8
	ldr	w25, [x9, w25, uxtw #2]
	eor	w7, w19, w7
	and	w5, w5, #0xff
	ldr	w5, [x10, w5, uxtw #2]
	eor	w5, w25, w5
	eor	w5, w7, w5
	lsr	w7, w6, #24
	ldr	w7, [x11, w7, uxtw #2]
	ubfx	w19, w6, #16, #8
	ldr	w19, [x12, w19, uxtw #2]
	ubfx	w25, w6, #8, #8
	ldr	w25, [x9, w25, uxtw #2]
	and	w6, w6, #0xff
	ldr	w6, [x10, w6, uxtw #2]
	eor	w7, w19, w7
	eor	w6, w25, w6
	eor	w6, w7, w6
	eor	w7, w24, w5
	eor	w5, w6, w5
	eor	w19, w7, w23
	eor	w6, w6, w7
	eor	w7, w5, w19
	eor	w5, w5, w23
	and	w23, w8, w6, lsl #8
	and	w6, w17, w6, lsr #8
	orr	w6, w23, w6
	ror	w19, w19, #16
	rev	w5, w5
	eor	w5, w5, w7
	eor	w7, w6, w7
	eor	w23, w5, w19
	eor	w15, w23, w15
	eor	w13, w19, w13
	eor	w13, w13, w7
	eor	w16, w23, w16
	eor	w16, w16, w7
	eor	w14, w5, w14
	eor	w14, w14, w6
	eor	w4, w15, w4
	eor	w5, w13, w20
	eor	w6, w16, w21
	eor	w7, w14, w22
	lsr	w19, w4, #24
	ldr	w19, [x9, w19, uxtw #2]
	ubfx	w20, w4, #16, #8
	ldr	w20, [x10, w20, uxtw #2]
	eor	w19, w20, w19
	ubfx	w20, w4, #8, #8
	ldr	w20, [x11, w20, uxtw #2]
	and	w4, w4, #0xff
	ldr	w4, [x12, w4, uxtw #2]
	eor	w4, w20, w4
	eor	w4, w19, w4
	lsr	w19, w5, #24
	ldr	w19, [x9, w19, uxtw #2]
	ubfx	w20, w5, #16, #8
	ldr	w20, [x10, w20, uxtw #2]
	eor	w19, w20, w19
	ubfx	w20, w5, #8, #8
	ldr	w20, [x11, w20, uxtw #2]
	and	w5, w5, #0xff
	ldr	w5, [x12, w5, uxtw #2]
	eor	w5, w20, w5
	eor	w5, w19, w5
	lsr	w19, w6, #24
	ldr	w19, [x9, w19, uxtw #2]
	ubfx	w20, w6, #16, #8
	ldr	w20, [x10, w20, uxtw #2]
	eor	w19, w20, w19
	ubfx	w20, w6, #8, #8
	ldr	w20, [x11, w20, uxtw #2]
	and	w6, w6, #0xff
	ldr	w6, [x12, w6, uxtw #2]
	eor	w6, w20, w6
	eor	w6, w19, w6
	lsr	w19, w7, #24
	ldr	w9, [x9, w19, uxtw #2]
	ubfx	w19, w7, #16, #8
	ldr	w10, [x10, w19, uxtw #2]
	eor	w9, w10, w9
	ubfx	w10, w7, #8, #8
	ldr	w10, [x11, w10, uxtw #2]
	and	w11, w7, #0xff
	ldr	w11, [x12, w11, uxtw #2]
	eor	w10, w10, w11
	eor	w9, w9, w10
	eor	w10, w5, w6
	eor	w11, w9, w6
	eor	w12, w10, w4
	eor	w5, w11, w12
	eor	w11, w11, w4
	and	w8, w8, w11, lsl #8
	and	w11, w17, w11, lsr #8
	orr	w8, w8, w11
	ror	w11, w5, #16
	eor	w8, w8, w11
	eor	w17, w8, w12
	eor	w4, w17, w3
	fmov	s2, w3
	mov.s	v2[1], w1
	mov.d	v2[1], v1[0]
	fmov	s3, w15
	mov.s	v3[1], w13
	mov.s	v3[2], w16
	mov.s	v3[3], w14
	eor	w9, w9, w10
	rev	w9, w9
	eor	w10, w11, w9
	eor	w11, w12, w1
	eor	w11, w11, w10
	fmov	s4, w10
	mov.s	v4[1], w9
	fmov	s5, w17
	mov.s	v5[1], w8
	eor.8b	v1, v4, v1
	eor.8b	v4, v1, v5
	fmov	s1, w4
	mov.s	v1[1], w11
	mov.d	v1[1], v4[0]
	ext.16b	v4, v2, v2, #12
	ushr.4s	v5, v2, #19
	shl.4s	v6, v4, #13
	orr.16b	v5, v6, v5
	ext.16b	v6, v3, v3, #12
	ushr.4s	v7, v3, #19
	shl.4s	v16, v6, #13
	orr.16b	v7, v16, v7
	ext.16b	v16, v1, v1, #12
	ushr.4s	v17, v1, #19
	shl.4s	v18, v16, #13
	orr.16b	v17, v18, v17
	ext.16b	v18, v0, v0, #12
	ushr.4s	v19, v0, #19
	shl.4s	v20, v18, #13
	orr.16b	v19, v20, v19
	ushr.4s	v20, v2, #31
	add.4s	v4, v4, v4
	orr.16b	v4, v4, v20
	ushr.4s	v20, v3, #31
	add.4s	v6, v6, v6
	orr.16b	v6, v6, v20
	ushr.4s	v20, v1, #31
	add.4s	v16, v16, v16
	orr.16b	v16, v16, v20
	ushr.4s	v20, v0, #31
	add.4s	v18, v18, v18
	orr.16b	v18, v18, v20
	eor.16b	v5, v5, v0
	eor.16b	v7, v7, v2
	eor.16b	v17, v17, v3
	eor.16b	v19, v1, v19
	eor.16b	v4, v4, v0
	eor.16b	v6, v6, v2
	eor.16b	v16, v16, v3
	eor.16b	v18, v1, v18
	stp	q5, q7, [x0]
	ushr.4s	v5, v2, #3
	stp	q17, q19, [x0, #32]
	stp	q4, q6, [x0, #64]
	stp	q16, q18, [x0, #96]
	ext.16b	v5, v5, v5, #8
	ext.16b	v4, v2, v2, #4
	shl.4s	v6, v4, #29
	orr.16b	v5, v6, v5
	eor.16b	v5, v5, v0
	ushr.4s	v6, v3, #3
	ext.16b	v7, v6, v6, #8
	ext.16b	v6, v3, v3, #4
	shl.4s	v16, v6, #29
	orr.16b	v7, v16, v7
	eor.16b	v7, v7, v2
	stp	q5, q7, [x0, #128]
	ushr.4s	v5, v1, #3
	ext.16b	v5, v5, v5, #8
	ext.16b	v7, v1, v1, #4
	shl.4s	v16, v7, #29
	orr.16b	v16, v16, v5
	ushr.4s	v5, v0, #3
	ext.16b	v17, v5, v5, #8
	ext.16b	v5, v0, v0, #4
	shl.4s	v18, v5, #29
	orr.16b	v17, v18, v17
	eor.16b	v16, v16, v3
	eor.16b	v17, v1, v17
	stp	q16, q17, [x0, #160]
	shl.4s	v16, v2, #31
	ushr.4s	v17, v4, #1
	orr.16b	v16, v16, v17
	eor.16b	v16, v16, v0
	str	q16, [x0, #192]
	cmp	x2, #17
	b.lo	LBB0_7
; %bb.5:
	ushr.4s	v6, v6, #1
	shl.4s	v16, v3, #31
	orr.16b	v6, v16, v6
	eor.16b	v6, v6, v2
	ushr.4s	v7, v7, #1
	shl.4s	v16, v1, #31
	orr.16b	v7, v16, v7
	eor.16b	v3, v7, v3
	stp	q6, q3, [x0, #208]
	cmp	x2, #25
	b.lo	LBB0_7
; %bb.6:
	ushr.4s	v3, v5, #1
	shl.4s	v5, v0, #31
	orr.16b	v3, v5, v3
	eor.16b	v1, v1, v3
	ushr.4s	v3, v4, #13
	shl.4s	v2, v2, #19
	orr.16b	v2, v2, v3
	eor.16b	v0, v2, v0
	stp	q1, q0, [x0, #240]
LBB0_7:
	ldp	x29, x30, [sp, #80]             ; 16-byte Folded Reload
	ldp	x20, x19, [sp, #64]             ; 16-byte Folded Reload
	ldp	x22, x21, [sp, #48]             ; 16-byte Folded Reload
	ldp	x24, x23, [sp, #32]             ; 16-byte Folded Reload
	ldp	x26, x25, [sp, #16]             ; 16-byte Folded Reload
	ldp	x28, x27, [sp], #96             ; 16-byte Folded Reload
	ret
	.loh AdrpAdd	Lloh8, Lloh9
	.loh AdrpAdd	Lloh6, Lloh7
	.loh AdrpAdd	Lloh4, Lloh5
	.loh AdrpAdd	Lloh2, Lloh3
	.loh AdrpAdd	Lloh0, Lloh1
                                        ; -- End function
	.globl	_DecKeySetup                    ; -- Begin function DecKeySetup
	.p2align	2
_DecKeySetup:                           ; @DecKeySetup
; %bb.0:
	mov	w8, #-16711936                  ; =0xff00ff00
	mov	w9, #16711935                   ; =0xff00ff
	lsl	w11, w1, #2
	add	x12, x0, w11, sxtw #2
	ldr	q0, [x0]
	ldr	q1, [x12]
	str	q1, [x0]
	mov	x13, x12
	str	q0, [x13], #-16
	cmp	w1, #3
	b.lt	LBB1_4
; %bb.1:
	mov	x10, #0                         ; =0x0
	sxtw	x11, w11
	lsl	x11, x11, #2
LBB1_2:                                 ; =>This Inner Loop Header: Depth=1
	add	x12, x0, x10
	ldp	w14, w15, [x12, #16]
	lsl	w16, w14, #8
	add	x13, x0, x11
	eor	w16, w16, w14, lsr #8
	eor	w16, w16, w14, lsl #16
	eor	w16, w16, w14, lsr #16
	eor	w16, w16, w14, lsl #24
	eor	w14, w16, w14, lsr #24
	lsl	w16, w15, #8
	eor	w16, w16, w15, lsr #8
	eor	w16, w16, w15, lsl #16
	eor	w16, w16, w15, lsr #16
	eor	w16, w16, w15, lsl #24
	eor	w15, w16, w15, lsr #24
	ldp	w16, w17, [x12, #24]
	lsl	w1, w16, #8
	eor	w1, w1, w16, lsr #8
	eor	w1, w1, w16, lsl #16
	eor	w1, w1, w16, lsr #16
	eor	w1, w1, w16, lsl #24
	eor	w16, w1, w16, lsr #24
	lsl	w1, w17, #8
	eor	w1, w1, w17, lsr #8
	eor	w1, w1, w17, lsl #16
	eor	w1, w1, w17, lsr #16
	eor	w1, w1, w17, lsl #24
	eor	w17, w1, w17, lsr #24
	eor	w15, w15, w16
	eor	w16, w17, w16
	eor	w1, w15, w14
	eor	w15, w17, w15
	eor	w17, w16, w1
	eor	w14, w16, w14
	and	w16, w8, w14, lsl #8
	and	w14, w9, w14, lsr #8
	orr	w14, w16, w14
	ror	w16, w17, #16
	rev	w15, w15
	eor	w14, w14, w16
	eor	w16, w16, w15
	eor	w17, w14, w1
	eor	w14, w14, w15
	eor	w15, w17, w16
	eor	w16, w16, w1
	ldp	w1, w2, [x13, #-16]
	lsl	w3, w2, #8
	eor	w3, w3, w2, lsr #8
	eor	w3, w3, w2, lsl #16
	eor	w3, w3, w2, lsr #16
	eor	w3, w3, w2, lsl #24
	eor	w2, w3, w2, lsr #24
	ldp	w3, w4, [x13, #-8]
	lsl	w5, w3, #8
	eor	w5, w5, w3, lsr #8
	eor	w5, w5, w3, lsl #16
	eor	w5, w5, w3, lsr #16
	eor	w5, w5, w3, lsl #24
	eor	w3, w5, w3, lsr #24
	lsl	w5, w4, #8
	eor	w5, w5, w4, lsr #8
	eor	w5, w5, w4, lsl #16
	eor	w5, w5, w4, lsr #16
	eor	w2, w2, w3
	eor	w5, w5, w4, lsl #24
	eor	w4, w5, w4, lsr #24
	eor	w3, w4, w3
	lsl	w5, w1, #8
	eor	w5, w5, w1, lsr #8
	eor	w5, w5, w1, lsl #16
	eor	w5, w5, w1, lsr #16
	eor	w5, w5, w1, lsl #24
	eor	w1, w5, w1, lsr #24
	eor	w5, w2, w1
	eor	w2, w4, w2
	eor	w4, w3, w5
	eor	w1, w3, w1
	and	w3, w8, w1, lsl #8
	and	w1, w9, w1, lsr #8
	orr	w1, w3, w1
	rev	w2, w2
	ror	w3, w4, #16
	eor	w1, w1, w3
	eor	w3, w3, w2
	eor	w4, w1, w5
	eor	w1, w1, w2
	eor	w2, w4, w3
	eor	w3, w3, w5
	stp	w4, w3, [x12, #16]
	stp	w2, w1, [x12, #24]
	stp	w17, w16, [x13, #-16]
	stp	w15, w14, [x13, #-8]
	add	x14, x12, #32
	sub	x11, x11, #16
	add	x10, x10, #16
	add	x12, x0, x11
	sub	x13, x12, #16
	cmp	x14, x13
	b.lo	LBB1_2
; %bb.3:
	add	x0, x0, x10
LBB1_4:
	add	x10, x0, #16
	stp	x29, x30, [sp, #-16]!           ; 16-byte Folded Spill
	mov	x29, sp
	ldp	w11, w14, [x0, #20]
	lsl	w15, w11, #8
	ldr	w10, [x10]
	eor	w15, w15, w11, lsr #8
	eor	w15, w15, w11, lsl #16
	eor	w15, w15, w11, lsr #16
	eor	w15, w15, w11, lsl #24
	eor	w11, w15, w11, lsr #24
	lsl	w15, w14, #8
	eor	w15, w15, w14, lsr #8
	eor	w15, w15, w14, lsl #16
	eor	w15, w15, w14, lsr #16
	eor	w15, w15, w14, lsl #24
	eor	w14, w15, w14, lsr #24
	ldr	w15, [x0, #28]
	lsl	w16, w15, #8
	eor	w16, w16, w15, lsr #8
	eor	w16, w16, w15, lsl #16
	eor	w16, w16, w15, lsr #16
	eor	w16, w16, w15, lsl #24
	eor	w15, w16, w15, lsr #24
	lsl	w16, w10, #8
	eor	w16, w16, w10, lsr #8
	eor	w16, w16, w10, lsl #16
	eor	w16, w16, w10, lsr #16
	eor	w16, w16, w10, lsl #24
	eor	w11, w11, w14
	eor	w14, w15, w14
	eor	w10, w16, w10, lsr #24
	eor	w16, w11, w10
	eor	w11, w15, w11
	eor	w15, w14, w16
	eor	w10, w14, w10
	and	w8, w8, w10, lsl #8
	and	w9, w9, w10, lsr #8
	orr	w8, w8, w9
	rev	w9, w11
	ror	w10, w15, #16
	eor	w8, w8, w10
	eor	w10, w10, w9
	eor	w11, w8, w16
	eor	w8, w8, w9
	eor	w9, w11, w10
	str	w11, [x13]
	eor	w10, w10, w16
	stp	w10, w9, [x12, #-12]
	stur	w8, [x12, #-4]
	ldp	x29, x30, [sp], #16             ; 16-byte Folded Reload
	ret
                                        ; -- End function
	.globl	_Crypt                          ; -- Begin function Crypt
	.p2align	2
_Crypt:                                 ; @Crypt
; %bb.0:
	stp	x29, x30, [sp, #-16]!           ; 16-byte Folded Spill
	mov	x29, sp
	mov	w12, #-16711936                 ; =0xff00ff00
	mov	w13, #16711935                  ; =0xff00ff
	ldp	w8, w9, [x1]
	rev	w16, w8
	rev	w17, w9
	ldp	w8, w10, [x1, #8]
	rev	w15, w8
Lloh10:
	adrp	x9, _S1@PAGE
Lloh11:
	add	x9, x9, _S1@PAGEOFF
Lloh12:
	adrp	x8, _S2@PAGE
Lloh13:
	add	x8, x8, _S2@PAGEOFF
	rev	w14, w10
Lloh14:
	adrp	x11, _X1@PAGE
Lloh15:
	add	x11, x11, _X1@PAGEOFF
Lloh16:
	adrp	x10, _X2@PAGE
Lloh17:
	add	x10, x10, _X2@PAGEOFF
	cmp	x3, #13
	b.lo	LBB2_4
; %bb.1:
	ldp	w1, w4, [x2]
	eor	w16, w1, w16
	eor	w17, w4, w17
	lsr	w1, w16, #24
	ldr	w1, [x9, w1, uxtw #2]
	ldp	w4, w5, [x2, #8]
	eor	w15, w4, w15
	eor	w14, w5, w14
	ubfx	w4, w16, #16, #8
	ldr	w4, [x8, w4, uxtw #2]
	ubfx	w5, w16, #8, #8
	ldr	w5, [x11, w5, uxtw #2]
	eor	w1, w4, w1
	and	w16, w16, #0xff
	ldr	w16, [x10, w16, uxtw #2]
	eor	w16, w5, w16
	eor	w16, w1, w16
	lsr	w1, w17, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w4, w17, #16, #8
	ldr	w4, [x8, w4, uxtw #2]
	ubfx	w5, w17, #8, #8
	ldr	w5, [x11, w5, uxtw #2]
	and	w17, w17, #0xff
	ldr	w17, [x10, w17, uxtw #2]
	eor	w1, w4, w1
	eor	w17, w5, w17
	eor	w17, w1, w17
	lsr	w1, w15, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w4, w15, #16, #8
	ldr	w4, [x8, w4, uxtw #2]
	ubfx	w5, w15, #8, #8
	ldr	w5, [x11, w5, uxtw #2]
	eor	w1, w4, w1
	and	w15, w15, #0xff
	ldr	w15, [x10, w15, uxtw #2]
	eor	w15, w5, w15
	eor	w15, w1, w15
	lsr	w1, w14, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w4, w14, #16, #8
	ldr	w4, [x8, w4, uxtw #2]
	ubfx	w5, w14, #8, #8
	ldr	w5, [x11, w5, uxtw #2]
	and	w14, w14, #0xff
	ldr	w14, [x10, w14, uxtw #2]
	eor	w1, w4, w1
	eor	w14, w5, w14
	eor	w14, w1, w14
	eor	w17, w17, w15
	eor	w15, w14, w15
	eor	w1, w17, w16
	eor	w14, w14, w17
	eor	w17, w15, w1
	eor	w15, w15, w16
	and	w16, w12, w15, lsl #8
	and	w15, w13, w15, lsr #8
	orr	w15, w16, w15
	ror	w16, w17, #16
	rev	w14, w14
	eor	w15, w15, w16
	eor	w16, w16, w14
	ldp	w17, w4, [x2, #16]
	eor	w4, w4, w1
	eor	w4, w4, w16
	ldp	w5, w6, [x2, #24]
	eor	w16, w16, w5
	eor	w1, w15, w1
	eor	w17, w1, w17
	eor	w16, w16, w1
	lsr	w1, w17, #24
	ldr	w1, [x11, w1, uxtw #2]
	eor	w14, w6, w14
	eor	w14, w14, w15
	ubfx	w15, w17, #16, #8
	ldr	w15, [x10, w15, uxtw #2]
	eor	w15, w15, w1
	ubfx	w1, w17, #8, #8
	ldr	w1, [x9, w1, uxtw #2]
	and	w17, w17, #0xff
	ldr	w17, [x8, w17, uxtw #2]
	eor	w17, w1, w17
	eor	w15, w15, w17
	lsr	w17, w4, #24
	ldr	w17, [x11, w17, uxtw #2]
	ubfx	w1, w4, #16, #8
	ldr	w1, [x10, w1, uxtw #2]
	eor	w17, w1, w17
	ubfx	w1, w4, #8, #8
	ldr	w1, [x9, w1, uxtw #2]
	and	w4, w4, #0xff
	ldr	w4, [x8, w4, uxtw #2]
	eor	w1, w1, w4
	eor	w17, w17, w1
	lsr	w1, w16, #24
	ldr	w1, [x11, w1, uxtw #2]
	ubfx	w4, w16, #16, #8
	ldr	w4, [x10, w4, uxtw #2]
	eor	w1, w4, w1
	ubfx	w4, w16, #8, #8
	ldr	w4, [x9, w4, uxtw #2]
	and	w16, w16, #0xff
	ldr	w16, [x8, w16, uxtw #2]
	eor	w16, w4, w16
	eor	w16, w1, w16
	lsr	w1, w14, #24
	ldr	w1, [x11, w1, uxtw #2]
	ubfx	w4, w14, #16, #8
	ldr	w4, [x10, w4, uxtw #2]
	eor	w1, w4, w1
	ubfx	w4, w14, #8, #8
	ldr	w4, [x9, w4, uxtw #2]
	and	w14, w14, #0xff
	ldr	w14, [x8, w14, uxtw #2]
	eor	w14, w4, w14
	eor	w14, w1, w14
	eor	w17, w17, w16
	eor	w16, w14, w16
	eor	w1, w17, w15
	eor	w14, w14, w17
	eor	w17, w16, w1
	eor	w15, w16, w15
	and	w16, w12, w14, lsl #8
	and	w14, w13, w14, lsr #8
	orr	w14, w16, w14
	ror	w1, w1, #16
	rev	w15, w15
	eor	w15, w15, w17
	eor	w17, w14, w17
	eor	w16, w15, w1
	eor	w14, w14, w15
	eor	w15, w17, w16
	eor	w17, w17, w1
	cmp	x3, #15
	b.lo	LBB2_3
; %bb.2:
	ldp	w1, w3, [x2, #32]
	eor	w16, w1, w16
	eor	w17, w3, w17
	lsr	w1, w16, #24
	ldr	w1, [x9, w1, uxtw #2]
	ldp	w3, w4, [x2, #40]
	eor	w15, w3, w15
	eor	w14, w4, w14
	ubfx	w3, w16, #16, #8
	ldr	w3, [x8, w3, uxtw #2]
	ubfx	w4, w16, #8, #8
	ldr	w4, [x11, w4, uxtw #2]
	eor	w1, w3, w1
	and	w16, w16, #0xff
	ldr	w16, [x10, w16, uxtw #2]
	eor	w16, w4, w16
	eor	w16, w1, w16
	lsr	w1, w17, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w3, w17, #16, #8
	ldr	w3, [x8, w3, uxtw #2]
	ubfx	w4, w17, #8, #8
	ldr	w4, [x11, w4, uxtw #2]
	and	w17, w17, #0xff
	ldr	w17, [x10, w17, uxtw #2]
	eor	w1, w3, w1
	eor	w17, w4, w17
	eor	w17, w1, w17
	lsr	w1, w15, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w3, w15, #16, #8
	ldr	w3, [x8, w3, uxtw #2]
	ubfx	w4, w15, #8, #8
	ldr	w4, [x11, w4, uxtw #2]
	eor	w1, w3, w1
	and	w15, w15, #0xff
	ldr	w15, [x10, w15, uxtw #2]
	eor	w15, w4, w15
	eor	w15, w1, w15
	lsr	w1, w14, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w3, w14, #16, #8
	ldr	w3, [x8, w3, uxtw #2]
	ubfx	w4, w14, #8, #8
	ldr	w4, [x11, w4, uxtw #2]
	and	w14, w14, #0xff
	ldr	w14, [x10, w14, uxtw #2]
	eor	w1, w3, w1
	eor	w14, w4, w14
	eor	w14, w1, w14
	eor	w17, w17, w15
	eor	w15, w14, w15
	eor	w1, w17, w16
	eor	w14, w14, w17
	eor	w17, w15, w1
	eor	w15, w15, w16
	and	w16, w12, w15, lsl #8
	and	w15, w13, w15, lsr #8
	orr	w15, w16, w15
	ror	w16, w17, #16
	rev	w14, w14
	eor	w15, w15, w16
	eor	w16, w16, w14
	ldp	w17, w3, [x2, #48]
	eor	w3, w3, w1
	eor	w3, w3, w16
	ldp	w4, w5, [x2, #56]
	eor	w16, w16, w4
	eor	w1, w15, w1
	eor	w17, w1, w17
	eor	w16, w16, w1
	lsr	w1, w17, #24
	ldr	w1, [x11, w1, uxtw #2]
	eor	w14, w5, w14
	ubfx	w4, w17, #16, #8
	eor	w14, w14, w15
	add	x2, x2, #64
	ldr	w15, [x10, w4, uxtw #2]
	ubfx	w4, w17, #8, #8
	ldr	w4, [x9, w4, uxtw #2]
	and	w17, w17, #0xff
	ldr	w17, [x8, w17, uxtw #2]
	eor	w15, w15, w1
	eor	w17, w4, w17
	eor	w15, w15, w17
	lsr	w17, w3, #24
	ldr	w17, [x11, w17, uxtw #2]
	ubfx	w1, w3, #16, #8
	ldr	w1, [x10, w1, uxtw #2]
	ubfx	w4, w3, #8, #8
	ldr	w4, [x9, w4, uxtw #2]
	eor	w17, w1, w17
	and	w1, w3, #0xff
	ldr	w1, [x8, w1, uxtw #2]
	eor	w1, w4, w1
	eor	w17, w17, w1
	lsr	w1, w16, #24
	ldr	w1, [x11, w1, uxtw #2]
	ubfx	w3, w16, #16, #8
	ldr	w3, [x10, w3, uxtw #2]
	ubfx	w4, w16, #8, #8
	ldr	w4, [x9, w4, uxtw #2]
	and	w16, w16, #0xff
	ldr	w16, [x8, w16, uxtw #2]
	eor	w1, w3, w1
	eor	w16, w4, w16
	eor	w16, w1, w16
	lsr	w1, w14, #24
	ldr	w1, [x11, w1, uxtw #2]
	ubfx	w3, w14, #16, #8
	ldr	w3, [x10, w3, uxtw #2]
	ubfx	w4, w14, #8, #8
	ldr	w4, [x9, w4, uxtw #2]
	eor	w1, w3, w1
	and	w14, w14, #0xff
	ldr	w14, [x8, w14, uxtw #2]
	eor	w14, w4, w14
	eor	w14, w1, w14
	eor	w17, w17, w16
	eor	w16, w14, w16
	eor	w1, w17, w15
	eor	w14, w14, w17
	eor	w17, w16, w1
	eor	w15, w16, w15
	and	w16, w12, w14, lsl #8
	and	w14, w13, w14, lsr #8
	orr	w14, w16, w14
	ror	w1, w1, #16
	rev	w15, w15
	eor	w15, w15, w17
	eor	w17, w14, w17
	eor	w16, w15, w1
	eor	w14, w14, w15
	eor	w15, w17, w16
	eor	w17, w17, w1
	b	LBB2_4
LBB2_3:
	add	x2, x2, #32
LBB2_4:
	ldp	w1, w3, [x2]
	eor	w16, w1, w16
	eor	w17, w3, w17
	ldp	w1, w3, [x2, #8]
	eor	w15, w1, w15
	eor	w14, w3, w14
	lsr	w1, w16, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w3, w16, #16, #8
	ldr	w3, [x8, w3, uxtw #2]
	eor	w1, w3, w1
	ubfx	w3, w16, #8, #8
	ldr	w3, [x11, w3, uxtw #2]
	and	w16, w16, #0xff
	ldr	w16, [x10, w16, uxtw #2]
	eor	w16, w3, w16
	eor	w16, w1, w16
	lsr	w1, w17, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w3, w17, #16, #8
	ldr	w3, [x8, w3, uxtw #2]
	eor	w1, w3, w1
	ubfx	w3, w17, #8, #8
	ldr	w3, [x11, w3, uxtw #2]
	and	w17, w17, #0xff
	ldr	w17, [x10, w17, uxtw #2]
	eor	w17, w3, w17
	eor	w17, w1, w17
	lsr	w1, w15, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w3, w15, #16, #8
	ldr	w3, [x8, w3, uxtw #2]
	eor	w1, w3, w1
	ubfx	w3, w15, #8, #8
	ldr	w3, [x11, w3, uxtw #2]
	and	w15, w15, #0xff
	ldr	w15, [x10, w15, uxtw #2]
	eor	w15, w3, w15
	eor	w15, w1, w15
	lsr	w1, w14, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w3, w14, #16, #8
	ldr	w3, [x8, w3, uxtw #2]
	eor	w1, w3, w1
	ubfx	w3, w14, #8, #8
	ldr	w3, [x11, w3, uxtw #2]
	and	w14, w14, #0xff
	ldr	w14, [x10, w14, uxtw #2]
	eor	w14, w3, w14
	eor	w14, w1, w14
	eor	w17, w17, w15
	eor	w15, w14, w15
	eor	w1, w17, w16
	eor	w14, w14, w17
	eor	w17, w15, w1
	eor	w15, w15, w16
	and	w16, w12, w15, lsl #8
	and	w15, w13, w15, lsr #8
	ror	w17, w17, #16
	orr	w15, w16, w15
	rev	w14, w14
	eor	w15, w15, w17
	eor	w16, w17, w14
	eor	w17, w15, w1
	ldp	w3, w4, [x2, #16]
	eor	w3, w17, w3
	eor	w1, w4, w1
	eor	w1, w1, w16
	ldp	w4, w5, [x2, #24]
	eor	w16, w16, w4
	eor	w16, w16, w17
	eor	w14, w5, w14
	lsr	w17, w3, #24
	ldr	w17, [x11, w17, uxtw #2]
	eor	w14, w14, w15
	ubfx	w15, w3, #16, #8
	ldr	w15, [x10, w15, uxtw #2]
	eor	w15, w15, w17
	ubfx	w17, w3, #8, #8
	ldr	w17, [x9, w17, uxtw #2]
	and	w3, w3, #0xff
	ldr	w3, [x8, w3, uxtw #2]
	eor	w17, w17, w3
	lsr	w3, w1, #24
	ldr	w3, [x11, w3, uxtw #2]
	ubfx	w4, w1, #16, #8
	ldr	w4, [x10, w4, uxtw #2]
	eor	w15, w15, w17
	eor	w17, w4, w3
	ubfx	w3, w1, #8, #8
	ldr	w3, [x9, w3, uxtw #2]
	and	w1, w1, #0xff
	ldr	w1, [x8, w1, uxtw #2]
	eor	w1, w3, w1
	lsr	w3, w16, #24
	ldr	w3, [x11, w3, uxtw #2]
	eor	w17, w17, w1
	ubfx	w1, w16, #16, #8
	ldr	w1, [x10, w1, uxtw #2]
	eor	w1, w1, w3
	ubfx	w3, w16, #8, #8
	ldr	w3, [x9, w3, uxtw #2]
	and	w16, w16, #0xff
	ldr	w16, [x8, w16, uxtw #2]
	eor	w16, w3, w16
	lsr	w3, w14, #24
	ldr	w3, [x11, w3, uxtw #2]
	ubfx	w4, w14, #16, #8
	ldr	w4, [x10, w4, uxtw #2]
	eor	w16, w1, w16
	eor	w1, w4, w3
	ubfx	w3, w14, #8, #8
	ldr	w3, [x9, w3, uxtw #2]
	and	w14, w14, #0xff
	ldr	w14, [x8, w14, uxtw #2]
	eor	w14, w3, w14
	eor	w14, w1, w14
	eor	w17, w17, w16
	eor	w16, w14, w16
	eor	w1, w17, w15
	eor	w14, w14, w17
	eor	w17, w16, w1
	eor	w15, w16, w15
	and	w16, w12, w14, lsl #8
	and	w14, w13, w14, lsr #8
	orr	w14, w16, w14
	ror	w16, w1, #16
	rev	w15, w15
	eor	w15, w15, w17
	eor	w17, w14, w17
	eor	w1, w15, w16
	ldp	w3, w4, [x2, #32]
	eor	w3, w1, w3
	eor	w16, w4, w16
	eor	w16, w16, w17
	ldp	w4, w5, [x2, #40]
	eor	w1, w1, w4
	eor	w17, w1, w17
	eor	w15, w15, w5
	eor	w14, w15, w14
	lsr	w15, w3, #24
	ldr	w15, [x9, w15, uxtw #2]
	ubfx	w1, w3, #16, #8
	ldr	w1, [x8, w1, uxtw #2]
	ubfx	w4, w3, #8, #8
	ldr	w4, [x11, w4, uxtw #2]
	and	w3, w3, #0xff
	ldr	w3, [x10, w3, uxtw #2]
	eor	w15, w1, w15
	eor	w1, w4, w3
	eor	w15, w15, w1
	lsr	w1, w16, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w3, w16, #16, #8
	ldr	w3, [x8, w3, uxtw #2]
	ubfx	w4, w16, #8, #8
	ldr	w4, [x11, w4, uxtw #2]
	eor	w1, w3, w1
	and	w16, w16, #0xff
	ldr	w16, [x10, w16, uxtw #2]
	eor	w16, w4, w16
	eor	w16, w1, w16
	lsr	w1, w17, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w3, w17, #16, #8
	ldr	w3, [x8, w3, uxtw #2]
	ubfx	w4, w17, #8, #8
	ldr	w4, [x11, w4, uxtw #2]
	and	w17, w17, #0xff
	ldr	w17, [x10, w17, uxtw #2]
	eor	w1, w3, w1
	eor	w17, w4, w17
	eor	w17, w1, w17
	lsr	w1, w14, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w3, w14, #16, #8
	ldr	w3, [x8, w3, uxtw #2]
	ubfx	w4, w14, #8, #8
	ldr	w4, [x11, w4, uxtw #2]
	eor	w1, w3, w1
	and	w14, w14, #0xff
	ldr	w14, [x10, w14, uxtw #2]
	eor	w14, w4, w14
	eor	w14, w1, w14
	eor	w16, w16, w17
	eor	w17, w14, w17
	eor	w1, w16, w15
	eor	w14, w14, w16
	eor	w16, w17, w1
	eor	w15, w17, w15
	and	w17, w12, w15, lsl #8
	and	w15, w13, w15, lsr #8
	orr	w15, w17, w15
	ror	w16, w16, #16
	rev	w14, w14
	eor	w15, w15, w16
	eor	w16, w16, w14
	eor	w17, w15, w1
	ldp	w3, w4, [x2, #48]
	eor	w3, w17, w3
	eor	w1, w4, w1
	eor	w1, w1, w16
	ldp	w4, w5, [x2, #56]
	eor	w16, w16, w4
	eor	w16, w16, w17
	eor	w14, w5, w14
	eor	w14, w14, w15
	lsr	w15, w3, #24
	ldr	w15, [x11, w15, uxtw #2]
	ubfx	w17, w3, #16, #8
	ldr	w17, [x10, w17, uxtw #2]
	eor	w15, w17, w15
	ubfx	w17, w3, #8, #8
	ldr	w17, [x9, w17, uxtw #2]
	and	w3, w3, #0xff
	ldr	w3, [x8, w3, uxtw #2]
	eor	w17, w17, w3
	eor	w15, w15, w17
	lsr	w17, w1, #24
	ldr	w17, [x11, w17, uxtw #2]
	ubfx	w3, w1, #16, #8
	ldr	w3, [x10, w3, uxtw #2]
	eor	w17, w3, w17
	ubfx	w3, w1, #8, #8
	ldr	w3, [x9, w3, uxtw #2]
	and	w1, w1, #0xff
	ldr	w1, [x8, w1, uxtw #2]
	eor	w1, w3, w1
	eor	w17, w17, w1
	lsr	w1, w16, #24
	ldr	w1, [x11, w1, uxtw #2]
	ubfx	w3, w16, #16, #8
	ldr	w3, [x10, w3, uxtw #2]
	eor	w1, w3, w1
	ubfx	w3, w16, #8, #8
	ldr	w3, [x9, w3, uxtw #2]
	and	w16, w16, #0xff
	ldr	w16, [x8, w16, uxtw #2]
	eor	w16, w3, w16
	eor	w16, w1, w16
	lsr	w1, w14, #24
	ldr	w1, [x11, w1, uxtw #2]
	ubfx	w3, w14, #16, #8
	ldr	w3, [x10, w3, uxtw #2]
	eor	w1, w3, w1
	ubfx	w3, w14, #8, #8
	ldr	w3, [x9, w3, uxtw #2]
	and	w14, w14, #0xff
	ldr	w14, [x8, w14, uxtw #2]
	eor	w14, w3, w14
	eor	w14, w1, w14
	eor	w17, w17, w16
	eor	w16, w14, w16
	eor	w1, w17, w15
	eor	w14, w14, w17
	eor	w17, w16, w1
	eor	w15, w16, w15
	and	w16, w12, w14, lsl #8
	and	w14, w13, w14, lsr #8
	ror	w1, w1, #16
	orr	w14, w16, w14
	rev	w15, w15
	eor	w15, w15, w17
	eor	w16, w14, w17
	eor	w17, w15, w1
	ldp	w3, w4, [x2, #64]
	eor	w3, w17, w3
	eor	w1, w4, w1
	eor	w1, w1, w16
	ldp	w4, w5, [x2, #72]
	eor	w17, w17, w4
	eor	w16, w17, w16
	eor	w15, w15, w5
	lsr	w17, w3, #24
	ldr	w17, [x9, w17, uxtw #2]
	eor	w14, w15, w14
	ubfx	w15, w3, #16, #8
	ldr	w15, [x8, w15, uxtw #2]
	eor	w15, w15, w17
	ubfx	w17, w3, #8, #8
	ldr	w17, [x11, w17, uxtw #2]
	and	w3, w3, #0xff
	ldr	w3, [x10, w3, uxtw #2]
	eor	w17, w17, w3
	lsr	w3, w1, #24
	ldr	w3, [x9, w3, uxtw #2]
	ubfx	w4, w1, #16, #8
	ldr	w4, [x8, w4, uxtw #2]
	eor	w15, w15, w17
	eor	w17, w4, w3
	ubfx	w3, w1, #8, #8
	ldr	w3, [x11, w3, uxtw #2]
	and	w1, w1, #0xff
	ldr	w1, [x10, w1, uxtw #2]
	eor	w1, w3, w1
	lsr	w3, w16, #24
	ldr	w3, [x9, w3, uxtw #2]
	eor	w17, w17, w1
	ubfx	w1, w16, #16, #8
	ldr	w1, [x8, w1, uxtw #2]
	eor	w1, w1, w3
	ubfx	w3, w16, #8, #8
	ldr	w3, [x11, w3, uxtw #2]
	and	w16, w16, #0xff
	ldr	w16, [x10, w16, uxtw #2]
	eor	w16, w3, w16
	lsr	w3, w14, #24
	ldr	w3, [x9, w3, uxtw #2]
	ubfx	w4, w14, #16, #8
	ldr	w4, [x8, w4, uxtw #2]
	eor	w16, w1, w16
	eor	w1, w4, w3
	ubfx	w3, w14, #8, #8
	ldr	w3, [x11, w3, uxtw #2]
	and	w14, w14, #0xff
	ldr	w14, [x10, w14, uxtw #2]
	eor	w14, w3, w14
	eor	w14, w1, w14
	eor	w17, w17, w16
	eor	w16, w14, w16
	eor	w1, w17, w15
	eor	w14, w14, w17
	eor	w17, w16, w1
	eor	w15, w16, w15
	and	w16, w12, w15, lsl #8
	and	w15, w13, w15, lsr #8
	orr	w15, w16, w15
	ror	w16, w17, #16
	rev	w14, w14
	eor	w15, w15, w16
	eor	w16, w16, w14
	eor	w17, w15, w1
	ldp	w3, w4, [x2, #80]
	eor	w3, w17, w3
	eor	w1, w4, w1
	eor	w1, w1, w16
	ldp	w4, w5, [x2, #88]
	eor	w16, w16, w4
	eor	w16, w16, w17
	eor	w14, w5, w14
	eor	w14, w14, w15
	lsr	w15, w3, #24
	ldr	w15, [x11, w15, uxtw #2]
	ubfx	w17, w3, #16, #8
	ldr	w17, [x10, w17, uxtw #2]
	ubfx	w4, w3, #8, #8
	ldr	w4, [x9, w4, uxtw #2]
	and	w3, w3, #0xff
	ldr	w3, [x8, w3, uxtw #2]
	eor	w15, w17, w15
	eor	w17, w4, w3
	eor	w15, w15, w17
	lsr	w17, w1, #24
	ldr	w17, [x11, w17, uxtw #2]
	ubfx	w3, w1, #16, #8
	ldr	w3, [x10, w3, uxtw #2]
	ubfx	w4, w1, #8, #8
	ldr	w4, [x9, w4, uxtw #2]
	eor	w17, w3, w17
	and	w1, w1, #0xff
	ldr	w1, [x8, w1, uxtw #2]
	eor	w1, w4, w1
	eor	w17, w17, w1
	lsr	w1, w16, #24
	ldr	w1, [x11, w1, uxtw #2]
	ubfx	w3, w16, #16, #8
	ldr	w3, [x10, w3, uxtw #2]
	ubfx	w4, w16, #8, #8
	ldr	w4, [x9, w4, uxtw #2]
	and	w16, w16, #0xff
	ldr	w16, [x8, w16, uxtw #2]
	eor	w1, w3, w1
	eor	w16, w4, w16
	eor	w16, w1, w16
	lsr	w1, w14, #24
	ldr	w1, [x11, w1, uxtw #2]
	ubfx	w3, w14, #16, #8
	ldr	w3, [x10, w3, uxtw #2]
	ubfx	w4, w14, #8, #8
	ldr	w4, [x9, w4, uxtw #2]
	eor	w1, w3, w1
	and	w14, w14, #0xff
	ldr	w14, [x8, w14, uxtw #2]
	eor	w14, w4, w14
	eor	w14, w1, w14
	eor	w17, w17, w16
	eor	w16, w14, w16
	eor	w1, w17, w15
	eor	w14, w14, w17
	eor	w17, w16, w1
	eor	w15, w16, w15
	and	w16, w12, w14, lsl #8
	and	w14, w13, w14, lsr #8
	orr	w14, w16, w14
	ror	w16, w1, #16
	rev	w15, w15
	eor	w15, w15, w17
	eor	w17, w14, w17
	eor	w1, w15, w16
	ldp	w3, w4, [x2, #96]
	eor	w3, w1, w3
	eor	w16, w4, w16
	eor	w16, w16, w17
	ldp	w4, w5, [x2, #104]
	eor	w1, w1, w4
	eor	w17, w1, w17
	eor	w15, w15, w5
	eor	w14, w15, w14
	lsr	w15, w3, #24
	ldr	w15, [x9, w15, uxtw #2]
	ubfx	w1, w3, #16, #8
	ldr	w1, [x8, w1, uxtw #2]
	eor	w15, w1, w15
	ubfx	w1, w3, #8, #8
	ldr	w1, [x11, w1, uxtw #2]
	and	w3, w3, #0xff
	ldr	w3, [x10, w3, uxtw #2]
	eor	w1, w1, w3
	eor	w15, w15, w1
	lsr	w1, w16, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w3, w16, #16, #8
	ldr	w3, [x8, w3, uxtw #2]
	eor	w1, w3, w1
	ubfx	w3, w16, #8, #8
	ldr	w3, [x11, w3, uxtw #2]
	and	w16, w16, #0xff
	ldr	w16, [x10, w16, uxtw #2]
	eor	w16, w3, w16
	eor	w16, w1, w16
	lsr	w1, w17, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w3, w17, #16, #8
	ldr	w3, [x8, w3, uxtw #2]
	eor	w1, w3, w1
	ubfx	w3, w17, #8, #8
	ldr	w3, [x11, w3, uxtw #2]
	and	w17, w17, #0xff
	ldr	w17, [x10, w17, uxtw #2]
	eor	w17, w3, w17
	eor	w17, w1, w17
	lsr	w1, w14, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w3, w14, #16, #8
	ldr	w3, [x8, w3, uxtw #2]
	eor	w1, w3, w1
	ubfx	w3, w14, #8, #8
	ldr	w3, [x11, w3, uxtw #2]
	and	w14, w14, #0xff
	ldr	w14, [x10, w14, uxtw #2]
	eor	w14, w3, w14
	eor	w14, w1, w14
	eor	w16, w16, w17
	eor	w17, w14, w17
	eor	w1, w16, w15
	eor	w14, w14, w16
	eor	w16, w17, w1
	eor	w15, w17, w15
	and	w17, w12, w15, lsl #8
	and	w15, w13, w15, lsr #8
	ror	w16, w16, #16
	orr	w15, w17, w15
	rev	w14, w14
	eor	w15, w15, w16
	eor	w16, w16, w14
	eor	w17, w15, w1
	ldp	w3, w4, [x2, #112]
	eor	w3, w17, w3
	eor	w1, w4, w1
	eor	w1, w1, w16
	ldp	w4, w5, [x2, #120]
	eor	w16, w16, w4
	eor	w16, w16, w17
	eor	w14, w5, w14
	lsr	w17, w3, #24
	ldr	w17, [x11, w17, uxtw #2]
	eor	w14, w14, w15
	ubfx	w15, w3, #16, #8
	ldr	w15, [x10, w15, uxtw #2]
	eor	w15, w15, w17
	ubfx	w17, w3, #8, #8
	ldr	w17, [x9, w17, uxtw #2]
	and	w3, w3, #0xff
	ldr	w3, [x8, w3, uxtw #2]
	eor	w17, w17, w3
	lsr	w3, w1, #24
	ldr	w3, [x11, w3, uxtw #2]
	ubfx	w4, w1, #16, #8
	ldr	w4, [x10, w4, uxtw #2]
	eor	w15, w15, w17
	eor	w17, w4, w3
	ubfx	w3, w1, #8, #8
	ldr	w3, [x9, w3, uxtw #2]
	and	w1, w1, #0xff
	ldr	w1, [x8, w1, uxtw #2]
	eor	w1, w3, w1
	lsr	w3, w16, #24
	ldr	w3, [x11, w3, uxtw #2]
	eor	w17, w17, w1
	ubfx	w1, w16, #16, #8
	ldr	w1, [x10, w1, uxtw #2]
	eor	w1, w1, w3
	ubfx	w3, w16, #8, #8
	ldr	w3, [x9, w3, uxtw #2]
	and	w16, w16, #0xff
	ldr	w16, [x8, w16, uxtw #2]
	eor	w16, w3, w16
	lsr	w3, w14, #24
	ldr	w3, [x11, w3, uxtw #2]
	ubfx	w4, w14, #16, #8
	ldr	w4, [x10, w4, uxtw #2]
	eor	w16, w1, w16
	eor	w1, w4, w3
	ubfx	w3, w14, #8, #8
	ldr	w3, [x9, w3, uxtw #2]
	and	w14, w14, #0xff
	ldr	w14, [x8, w14, uxtw #2]
	eor	w14, w3, w14
	eor	w14, w1, w14
	eor	w17, w17, w16
	eor	w16, w14, w16
	eor	w1, w17, w15
	eor	w14, w14, w17
	eor	w17, w16, w1
	eor	w15, w16, w15
	and	w16, w12, w14, lsl #8
	and	w14, w13, w14, lsr #8
	orr	w14, w16, w14
	ror	w16, w1, #16
	rev	w15, w15
	eor	w15, w15, w17
	eor	w17, w14, w17
	eor	w1, w15, w16
	ldp	w3, w4, [x2, #128]
	eor	w3, w1, w3
	eor	w16, w4, w16
	eor	w16, w16, w17
	ldp	w4, w5, [x2, #136]
	eor	w1, w1, w4
	eor	w17, w1, w17
	eor	w15, w15, w5
	eor	w14, w15, w14
	lsr	w15, w3, #24
	ldr	w15, [x9, w15, uxtw #2]
	ubfx	w1, w3, #16, #8
	ldr	w1, [x8, w1, uxtw #2]
	ubfx	w4, w3, #8, #8
	ldr	w4, [x11, w4, uxtw #2]
	and	w3, w3, #0xff
	ldr	w3, [x10, w3, uxtw #2]
	eor	w15, w1, w15
	eor	w1, w4, w3
	eor	w15, w15, w1
	lsr	w1, w16, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w3, w16, #16, #8
	ldr	w3, [x8, w3, uxtw #2]
	ubfx	w4, w16, #8, #8
	ldr	w4, [x11, w4, uxtw #2]
	eor	w1, w3, w1
	and	w16, w16, #0xff
	ldr	w16, [x10, w16, uxtw #2]
	eor	w16, w4, w16
	eor	w16, w1, w16
	lsr	w1, w17, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w3, w17, #16, #8
	ldr	w3, [x8, w3, uxtw #2]
	ubfx	w4, w17, #8, #8
	ldr	w4, [x11, w4, uxtw #2]
	and	w17, w17, #0xff
	ldr	w17, [x10, w17, uxtw #2]
	eor	w1, w3, w1
	eor	w17, w4, w17
	eor	w17, w1, w17
	lsr	w1, w14, #24
	ldr	w1, [x9, w1, uxtw #2]
	ubfx	w3, w14, #16, #8
	ldr	w3, [x8, w3, uxtw #2]
	ubfx	w4, w14, #8, #8
	ldr	w4, [x11, w4, uxtw #2]
	eor	w1, w3, w1
	and	w14, w14, #0xff
	ldr	w14, [x10, w14, uxtw #2]
	eor	w14, w4, w14
	eor	w14, w1, w14
	eor	w16, w16, w17
	eor	w17, w14, w17
	eor	w1, w16, w15
	eor	w14, w14, w16
	eor	w16, w17, w1
	eor	w15, w17, w15
	and	w17, w12, w15, lsl #8
	and	w15, w13, w15, lsr #8
	orr	w15, w17, w15
	ror	w16, w16, #16
	rev	w14, w14
	eor	w15, w15, w16
	eor	w16, w16, w14
	eor	w17, w15, w1
	ldp	w3, w4, [x2, #144]
	eor	w3, w17, w3
	eor	w1, w4, w1
	eor	w1, w1, w16
	ldp	w4, w5, [x2, #152]
	eor	w16, w16, w4
	eor	w16, w16, w17
	eor	w14, w5, w14
	eor	w14, w14, w15
	lsr	w15, w3, #24
	ldr	w15, [x11, w15, uxtw #2]
	ubfx	w17, w3, #16, #8
	ldr	w17, [x10, w17, uxtw #2]
	eor	w15, w17, w15
	ubfx	w17, w3, #8, #8
	ldr	w17, [x9, w17, uxtw #2]
	and	w3, w3, #0xff
	ldr	w3, [x8, w3, uxtw #2]
	eor	w17, w17, w3
	eor	w15, w15, w17
	lsr	w17, w1, #24
	ldr	w17, [x11, w17, uxtw #2]
	ubfx	w3, w1, #16, #8
	ldr	w3, [x10, w3, uxtw #2]
	eor	w17, w3, w17
	ubfx	w3, w1, #8, #8
	ldr	w3, [x9, w3, uxtw #2]
	and	w1, w1, #0xff
	ldr	w1, [x8, w1, uxtw #2]
	eor	w1, w3, w1
	eor	w17, w17, w1
	lsr	w1, w16, #24
	ldr	w1, [x11, w1, uxtw #2]
	ubfx	w3, w16, #16, #8
	ldr	w3, [x10, w3, uxtw #2]
	eor	w1, w3, w1
	ubfx	w3, w16, #8, #8
	ldr	w3, [x9, w3, uxtw #2]
	and	w16, w16, #0xff
	ldr	w16, [x8, w16, uxtw #2]
	eor	w16, w3, w16
	eor	w16, w1, w16
	lsr	w1, w14, #24
	ldr	w1, [x11, w1, uxtw #2]
	ubfx	w3, w14, #16, #8
	ldr	w3, [x10, w3, uxtw #2]
	eor	w1, w3, w1
	ubfx	w3, w14, #8, #8
	ldr	w3, [x9, w3, uxtw #2]
	and	w14, w14, #0xff
	ldr	w14, [x8, w14, uxtw #2]
	eor	w14, w3, w14
	eor	w14, w1, w14
	eor	w17, w17, w16
	eor	w16, w14, w16
	eor	w1, w17, w15
	eor	w14, w14, w17
	eor	w17, w16, w1
	eor	w15, w16, w15
	and	w16, w12, w14, lsl #8
	and	w14, w13, w14, lsr #8
	ror	w1, w1, #16
	orr	w14, w16, w14
	rev	w15, w15
	eor	w15, w15, w17
	eor	w16, w15, w1
	ldp	w3, w4, [x2, #160]
	eor	w3, w16, w3
	eor	w1, w4, w1
	ldp	w4, w5, [x2, #168]
	eor	w16, w16, w4
	eor	w17, w14, w17
	eor	w1, w1, w17
	eor	w16, w16, w17
	eor	w15, w15, w5
	lsr	w17, w3, #24
	ldr	w17, [x9, w17, uxtw #2]
	eor	w14, w15, w14
	ubfx	w15, w3, #16, #8
	ldr	w15, [x8, w15, uxtw #2]
	eor	w15, w15, w17
	ubfx	w17, w3, #8, #8
	ldr	w17, [x11, w17, uxtw #2]
	and	w3, w3, #0xff
	ldr	w3, [x10, w3, uxtw #2]
	eor	w17, w17, w3
	lsr	w3, w1, #24
	ldr	w3, [x9, w3, uxtw #2]
	ubfx	w4, w1, #16, #8
	ldr	w4, [x8, w4, uxtw #2]
	eor	w15, w15, w17
	eor	w17, w4, w3
	ubfx	w3, w1, #8, #8
	ldr	w3, [x11, w3, uxtw #2]
	and	w1, w1, #0xff
	ldr	w1, [x10, w1, uxtw #2]
	eor	w1, w3, w1
	lsr	w3, w16, #24
	ldr	w3, [x9, w3, uxtw #2]
	eor	w17, w17, w1
	ubfx	w1, w16, #16, #8
	ldr	w1, [x8, w1, uxtw #2]
	eor	w1, w1, w3
	ubfx	w3, w16, #8, #8
	ldr	w3, [x11, w3, uxtw #2]
	and	w16, w16, #0xff
	ldr	w16, [x10, w16, uxtw #2]
	eor	w16, w3, w16
	lsr	w3, w14, #24
	ldr	w3, [x9, w3, uxtw #2]
	ubfx	w4, w14, #16, #8
	ldr	w4, [x8, w4, uxtw #2]
	eor	w16, w1, w16
	eor	w1, w4, w3
	ubfx	w3, w14, #8, #8
	ldr	w3, [x11, w3, uxtw #2]
	and	w14, w14, #0xff
	ldr	w14, [x10, w14, uxtw #2]
	eor	w14, w3, w14
	eor	w14, w1, w14
	eor	w17, w17, w16
	eor	w16, w14, w16
	eor	w1, w17, w15
	eor	w14, w14, w17
	eor	w17, w16, w1
	eor	w15, w16, w15
	and	w12, w12, w15, lsl #8
	and	w13, w13, w15, lsr #8
	orr	w12, w12, w13
	lsr	w13, w17, #16
	rev	w14, w14
	eor	w12, w12, w13
	eor	w13, w13, w14
	eor	w15, w12, w1
	ldp	w16, w17, [x2, #176]
	eor	w16, w15, w16
	eor	w17, w17, w1
	eor	w17, w17, w13
	ldp	w1, w3, [x2, #184]
	eor	w13, w13, w1
	ubfx	w1, w16, #3, #8
	ldr	w1, [x11, w1, uxtw #2]
	strb	w1, [x0]
	ubfx	w1, w16, #2, #8
	ldr	w1, [x10, w1, uxtw #2]
	lsr	w1, w1, #8
	strb	w1, [x0, #1]
	ubfx	w1, w16, #1, #8
	ldr	w1, [x9, w1, uxtw #2]
	strb	w1, [x0, #2]
	and	w16, w16, #0xff
	ldr	w16, [x8, w16, uxtw #2]
	eor	w13, w13, w15
	strb	w16, [x0, #3]
	ubfx	w15, w17, #3, #8
	ldr	w15, [x11, w15, uxtw #2]
	strb	w15, [x0, #4]
	ubfx	w15, w17, #2, #8
	ldr	w15, [x10, w15, uxtw #2]
	lsr	w15, w15, #8
	strb	w15, [x0, #5]
	ubfx	w15, w17, #1, #8
	ldr	w15, [x9, w15, uxtw #2]
	strb	w15, [x0, #6]
	and	w15, w17, #0xff
	ldr	w15, [x8, w15, uxtw #2]
	eor	w14, w3, w14
	strb	w15, [x0, #7]
	ubfx	w15, w13, #3, #8
	ldr	w15, [x11, w15, uxtw #2]
	strb	w15, [x0, #8]
	ubfx	w15, w13, #2, #8
	ldr	w15, [x10, w15, uxtw #2]
	lsr	w15, w15, #8
	strb	w15, [x0, #9]
	ubfx	w15, w13, #1, #8
	ldr	w15, [x9, w15, uxtw #2]
	strb	w15, [x0, #10]
	and	w13, w13, #0xff
	ldr	w13, [x8, w13, uxtw #2]
	eor	w12, w14, w12
	strb	w13, [x0, #11]
	ubfx	w13, w12, #3, #8
	ldr	w11, [x11, w13, uxtw #2]
	strb	w11, [x0, #12]
	ubfx	w11, w12, #2, #8
	ldr	w10, [x10, w11, uxtw #2]
	lsr	w10, w10, #8
	strb	w10, [x0, #13]
	ubfx	w10, w12, #1, #8
	ldr	w9, [x9, w10, uxtw #2]
	strb	w9, [x0, #14]
	and	w9, w12, #0xff
	ldr	w8, [x8, w9, uxtw #2]
	strb	w8, [x0, #15]
	ldr	q0, [x0]
	ldr	q1, [x2, #192]
	rev32.16b	v1, v1
	eor.16b	v0, v1, v0
	str	q0, [x0]
	ldp	x29, x30, [sp], #16             ; 16-byte Folded Reload
	ret
	.loh AdrpAdd	Lloh16, Lloh17
	.loh AdrpAdd	Lloh14, Lloh15
	.loh AdrpAdd	Lloh12, Lloh13
	.loh AdrpAdd	Lloh10, Lloh11
                                        ; -- End function
	.section	__TEXT,__const
	.globl	_KRK                            ; @KRK
	.p2align	2, 0x0
_KRK:
	.long	1367130551                      ; 0x517cc1b7
	.long	656542356                       ; 0x27220a94
	.long	4262702056                      ; 0xfe13abe8
	.long	4204424928                      ; 0xfa9a6ee0
	.long	1840335564                      ; 0x6db14acc
	.long	2653014048                      ; 0x9e21c820
	.long	4280857045                      ; 0xff28b1d5
	.long	4015907504                      ; 0xef5de2b0
	.long	3683792669                      ; 0xdb92371d
	.long	556198256                       ; 0x2126e970
	.long	52729717                        ; 0x3249775
	.long	82364686                        ; 0x4e8c90e

	.globl	_S1                             ; @S1
	.p2align	2, 0x0
_S1:
	.long	6513507                         ; 0x636363
	.long	8158332                         ; 0x7c7c7c
	.long	7829367                         ; 0x777777
	.long	8092539                         ; 0x7b7b7b
	.long	15921906                        ; 0xf2f2f2
	.long	7039851                         ; 0x6b6b6b
	.long	7303023                         ; 0x6f6f6f
	.long	12961221                        ; 0xc5c5c5
	.long	3158064                         ; 0x303030
	.long	65793                           ; 0x10101
	.long	6776679                         ; 0x676767
	.long	2829099                         ; 0x2b2b2b
	.long	16711422                        ; 0xfefefe
	.long	14145495                        ; 0xd7d7d7
	.long	11250603                        ; 0xababab
	.long	7763574                         ; 0x767676
	.long	13290186                        ; 0xcacaca
	.long	8553090                         ; 0x828282
	.long	13224393                        ; 0xc9c9c9
	.long	8224125                         ; 0x7d7d7d
	.long	16448250                        ; 0xfafafa
	.long	5855577                         ; 0x595959
	.long	4671303                         ; 0x474747
	.long	15790320                        ; 0xf0f0f0
	.long	11382189                        ; 0xadadad
	.long	13948116                        ; 0xd4d4d4
	.long	10658466                        ; 0xa2a2a2
	.long	11513775                        ; 0xafafaf
	.long	10263708                        ; 0x9c9c9c
	.long	10790052                        ; 0xa4a4a4
	.long	7500402                         ; 0x727272
	.long	12632256                        ; 0xc0c0c0
	.long	12040119                        ; 0xb7b7b7
	.long	16645629                        ; 0xfdfdfd
	.long	9671571                         ; 0x939393
	.long	2500134                         ; 0x262626
	.long	3552822                         ; 0x363636
	.long	4144959                         ; 0x3f3f3f
	.long	16250871                        ; 0xf7f7f7
	.long	13421772                        ; 0xcccccc
	.long	3421236                         ; 0x343434
	.long	10855845                        ; 0xa5a5a5
	.long	15066597                        ; 0xe5e5e5
	.long	15856113                        ; 0xf1f1f1
	.long	7434609                         ; 0x717171
	.long	14211288                        ; 0xd8d8d8
	.long	3223857                         ; 0x313131
	.long	1381653                         ; 0x151515
	.long	263172                          ; 0x40404
	.long	13092807                        ; 0xc7c7c7
	.long	2302755                         ; 0x232323
	.long	12829635                        ; 0xc3c3c3
	.long	1579032                         ; 0x181818
	.long	9868950                         ; 0x969696
	.long	328965                          ; 0x50505
	.long	10132122                        ; 0x9a9a9a
	.long	460551                          ; 0x70707
	.long	1184274                         ; 0x121212
	.long	8421504                         ; 0x808080
	.long	14869218                        ; 0xe2e2e2
	.long	15461355                        ; 0xebebeb
	.long	2565927                         ; 0x272727
	.long	11711154                        ; 0xb2b2b2
	.long	7697781                         ; 0x757575
	.long	592137                          ; 0x90909
	.long	8618883                         ; 0x838383
	.long	2894892                         ; 0x2c2c2c
	.long	1710618                         ; 0x1a1a1a
	.long	1776411                         ; 0x1b1b1b
	.long	7237230                         ; 0x6e6e6e
	.long	5921370                         ; 0x5a5a5a
	.long	10526880                        ; 0xa0a0a0
	.long	5395026                         ; 0x525252
	.long	3881787                         ; 0x3b3b3b
	.long	14079702                        ; 0xd6d6d6
	.long	11776947                        ; 0xb3b3b3
	.long	2697513                         ; 0x292929
	.long	14935011                        ; 0xe3e3e3
	.long	3092271                         ; 0x2f2f2f
	.long	8684676                         ; 0x848484
	.long	5460819                         ; 0x535353
	.long	13750737                        ; 0xd1d1d1
	.long	0                               ; 0x0
	.long	15592941                        ; 0xededed
	.long	2105376                         ; 0x202020
	.long	16579836                        ; 0xfcfcfc
	.long	11645361                        ; 0xb1b1b1
	.long	5987163                         ; 0x5b5b5b
	.long	6974058                         ; 0x6a6a6a
	.long	13355979                        ; 0xcbcbcb
	.long	12500670                        ; 0xbebebe
	.long	3750201                         ; 0x393939
	.long	4868682                         ; 0x4a4a4a
	.long	5000268                         ; 0x4c4c4c
	.long	5789784                         ; 0x585858
	.long	13619151                        ; 0xcfcfcf
	.long	13684944                        ; 0xd0d0d0
	.long	15724527                        ; 0xefefef
	.long	11184810                        ; 0xaaaaaa
	.long	16514043                        ; 0xfbfbfb
	.long	4408131                         ; 0x434343
	.long	5066061                         ; 0x4d4d4d
	.long	3355443                         ; 0x333333
	.long	8750469                         ; 0x858585
	.long	4539717                         ; 0x454545
	.long	16382457                        ; 0xf9f9f9
	.long	131586                          ; 0x20202
	.long	8355711                         ; 0x7f7f7f
	.long	5263440                         ; 0x505050
	.long	3947580                         ; 0x3c3c3c
	.long	10461087                        ; 0x9f9f9f
	.long	11053224                        ; 0xa8a8a8
	.long	5329233                         ; 0x515151
	.long	10724259                        ; 0xa3a3a3
	.long	4210752                         ; 0x404040
	.long	9408399                         ; 0x8f8f8f
	.long	9605778                         ; 0x929292
	.long	10329501                        ; 0x9d9d9d
	.long	3684408                         ; 0x383838
	.long	16119285                        ; 0xf5f5f5
	.long	12369084                        ; 0xbcbcbc
	.long	11974326                        ; 0xb6b6b6
	.long	14342874                        ; 0xdadada
	.long	2171169                         ; 0x212121
	.long	1052688                         ; 0x101010
	.long	16777215                        ; 0xffffff
	.long	15987699                        ; 0xf3f3f3
	.long	13816530                        ; 0xd2d2d2
	.long	13487565                        ; 0xcdcdcd
	.long	789516                          ; 0xc0c0c
	.long	1250067                         ; 0x131313
	.long	15527148                        ; 0xececec
	.long	6250335                         ; 0x5f5f5f
	.long	9934743                         ; 0x979797
	.long	4473924                         ; 0x444444
	.long	1513239                         ; 0x171717
	.long	12895428                        ; 0xc4c4c4
	.long	10987431                        ; 0xa7a7a7
	.long	8289918                         ; 0x7e7e7e
	.long	4013373                         ; 0x3d3d3d
	.long	6579300                         ; 0x646464
	.long	6118749                         ; 0x5d5d5d
	.long	1644825                         ; 0x191919
	.long	7566195                         ; 0x737373
	.long	6316128                         ; 0x606060
	.long	8487297                         ; 0x818181
	.long	5197647                         ; 0x4f4f4f
	.long	14474460                        ; 0xdcdcdc
	.long	2236962                         ; 0x222222
	.long	2763306                         ; 0x2a2a2a
	.long	9474192                         ; 0x909090
	.long	8947848                         ; 0x888888
	.long	4605510                         ; 0x464646
	.long	15658734                        ; 0xeeeeee
	.long	12105912                        ; 0xb8b8b8
	.long	1315860                         ; 0x141414
	.long	14606046                        ; 0xdedede
	.long	6184542                         ; 0x5e5e5e
	.long	723723                          ; 0xb0b0b
	.long	14408667                        ; 0xdbdbdb
	.long	14737632                        ; 0xe0e0e0
	.long	3289650                         ; 0x323232
	.long	3815994                         ; 0x3a3a3a
	.long	657930                          ; 0xa0a0a
	.long	4802889                         ; 0x494949
	.long	394758                          ; 0x60606
	.long	2368548                         ; 0x242424
	.long	6052956                         ; 0x5c5c5c
	.long	12763842                        ; 0xc2c2c2
	.long	13882323                        ; 0xd3d3d3
	.long	11316396                        ; 0xacacac
	.long	6447714                         ; 0x626262
	.long	9539985                         ; 0x919191
	.long	9803157                         ; 0x959595
	.long	15000804                        ; 0xe4e4e4
	.long	7960953                         ; 0x797979
	.long	15198183                        ; 0xe7e7e7
	.long	13158600                        ; 0xc8c8c8
	.long	3618615                         ; 0x373737
	.long	7171437                         ; 0x6d6d6d
	.long	9276813                         ; 0x8d8d8d
	.long	14013909                        ; 0xd5d5d5
	.long	5131854                         ; 0x4e4e4e
	.long	11119017                        ; 0xa9a9a9
	.long	7105644                         ; 0x6c6c6c
	.long	5658198                         ; 0x565656
	.long	16053492                        ; 0xf4f4f4
	.long	15395562                        ; 0xeaeaea
	.long	6645093                         ; 0x656565
	.long	8026746                         ; 0x7a7a7a
	.long	11447982                        ; 0xaeaeae
	.long	526344                          ; 0x80808
	.long	12237498                        ; 0xbababa
	.long	7895160                         ; 0x787878
	.long	2434341                         ; 0x252525
	.long	3026478                         ; 0x2e2e2e
	.long	1842204                         ; 0x1c1c1c
	.long	10921638                        ; 0xa6a6a6
	.long	11842740                        ; 0xb4b4b4
	.long	13027014                        ; 0xc6c6c6
	.long	15263976                        ; 0xe8e8e8
	.long	14540253                        ; 0xdddddd
	.long	7631988                         ; 0x747474
	.long	2039583                         ; 0x1f1f1f
	.long	4934475                         ; 0x4b4b4b
	.long	12434877                        ; 0xbdbdbd
	.long	9145227                         ; 0x8b8b8b
	.long	9079434                         ; 0x8a8a8a
	.long	7368816                         ; 0x707070
	.long	4079166                         ; 0x3e3e3e
	.long	11908533                        ; 0xb5b5b5
	.long	6710886                         ; 0x666666
	.long	4737096                         ; 0x484848
	.long	197379                          ; 0x30303
	.long	16185078                        ; 0xf6f6f6
	.long	921102                          ; 0xe0e0e
	.long	6381921                         ; 0x616161
	.long	3487029                         ; 0x353535
	.long	5723991                         ; 0x575757
	.long	12171705                        ; 0xb9b9b9
	.long	8816262                         ; 0x868686
	.long	12698049                        ; 0xc1c1c1
	.long	1907997                         ; 0x1d1d1d
	.long	10395294                        ; 0x9e9e9e
	.long	14803425                        ; 0xe1e1e1
	.long	16316664                        ; 0xf8f8f8
	.long	10000536                        ; 0x989898
	.long	1118481                         ; 0x111111
	.long	6908265                         ; 0x696969
	.long	14277081                        ; 0xd9d9d9
	.long	9342606                         ; 0x8e8e8e
	.long	9737364                         ; 0x949494
	.long	10197915                        ; 0x9b9b9b
	.long	1973790                         ; 0x1e1e1e
	.long	8882055                         ; 0x878787
	.long	15329769                        ; 0xe9e9e9
	.long	13553358                        ; 0xcecece
	.long	5592405                         ; 0x555555
	.long	2631720                         ; 0x282828
	.long	14671839                        ; 0xdfdfdf
	.long	9211020                         ; 0x8c8c8c
	.long	10592673                        ; 0xa1a1a1
	.long	9013641                         ; 0x898989
	.long	855309                          ; 0xd0d0d
	.long	12566463                        ; 0xbfbfbf
	.long	15132390                        ; 0xe6e6e6
	.long	4342338                         ; 0x424242
	.long	6842472                         ; 0x686868
	.long	4276545                         ; 0x414141
	.long	10066329                        ; 0x999999
	.long	2960685                         ; 0x2d2d2d
	.long	986895                          ; 0xf0f0f
	.long	11579568                        ; 0xb0b0b0
	.long	5526612                         ; 0x545454
	.long	12303291                        ; 0xbbbbbb
	.long	1447446                         ; 0x161616

	.globl	_S2                             ; @S2
	.p2align	2, 0x0
_S2:
	.long	3791708898                      ; 0xe200e2e2
	.long	1308642894                      ; 0x4e004e4e
	.long	1409307732                      ; 0x54005454
	.long	4227923196                      ; 0xfc00fcfc
	.long	2483066004                      ; 0x94009494
	.long	3254829762                      ; 0xc200c2c2
	.long	1241533002                      ; 0x4a004a4a
	.long	3422604492                      ; 0xcc00cccc
	.long	1644192354                      ; 0x62006262
	.long	218107149                       ; 0xd000d0d
	.long	1778412138                      ; 0x6a006a6a
	.long	1174423110                      ; 0x46004646
	.long	1006648380                      ; 0x3c003c3c
	.long	1291865421                      ; 0x4d004d4d
	.long	2332068747                      ; 0x8b008b8b
	.long	3506491857                      ; 0xd100d1d1
	.long	1577082462                      ; 0x5e005e5e
	.long	4194368250                      ; 0xfa00fafa
	.long	1677747300                      ; 0x64006464
	.long	3405827019                      ; 0xcb00cbcb
	.long	3019945140                      ; 0xb400b4b4
	.long	2533398423                      ; 0x97009797
	.long	3187719870                      ; 0xbe00bebe
	.long	721431339                       ; 0x2b002b2b
	.long	3154164924                      ; 0xbc00bcbc
	.long	1996519287                      ; 0x77007777
	.long	771763758                       ; 0x2e002e2e
	.long	50332419                        ; 0x3000303
	.long	3540046803                      ; 0xd300d3d3
	.long	419436825                       ; 0x19001919
	.long	1493195097                      ; 0x59005959
	.long	3238052289                      ; 0xc100c1c1
	.long	486546717                       ; 0x1d001d1d
	.long	100664838                       ; 0x6000606
	.long	1090535745                      ; 0x41004141
	.long	1795189611                      ; 0x6b006b6b
	.long	1426085205                      ; 0x55005555
	.long	4026593520                      ; 0xf000f0f0
	.long	2566953369                      ; 0x99009999
	.long	1761634665                      ; 0x69006969
	.long	3925928682                      ; 0xea00eaea
	.long	2617285788                      ; 0x9c009c9c
	.long	402659352                       ; 0x18001818
	.long	2919280302                      ; 0xae00aeae
	.long	1660969827                      ; 0x63006363
	.long	3741376479                      ; 0xdf00dfdf
	.long	3875596263                      ; 0xe700e7e7
	.long	3137387451                      ; 0xbb00bbbb
	.long	0                               ; 0x0
	.long	1929409395                      ; 0x73007373
	.long	1711302246                      ; 0x66006666
	.long	4211145723                      ; 0xfb00fbfb
	.long	2516620950                      ; 0x96009696
	.long	1275087948                      ; 0x4c004c4c
	.long	2231403909                      ; 0x85008585
	.long	3825263844                      ; 0xe400e4e4
	.long	973093434                       ; 0x3a003a3a
	.long	150997257                       ; 0x9000909
	.long	1157645637                      ; 0x45004545
	.long	2852170410                      ; 0xaa00aaaa
	.long	251662095                       ; 0xf000f0f
	.long	3993038574                      ; 0xee00eeee
	.long	268439568                       ; 0x10001010
	.long	3942706155                      ; 0xeb00ebeb
	.long	754986285                       ; 0x2d002d2d
	.long	2130739071                      ; 0x7f007f7f
	.long	4093703412                      ; 0xf400f4f4
	.long	687876393                       ; 0x29002929
	.long	2885725356                      ; 0xac00acac
	.long	3472936911                      ; 0xcf00cfcf
	.long	2902502829                      ; 0xad00adad
	.long	2432733585                      ; 0x91009191
	.long	2365623693                      ; 0x8d008d8d
	.long	2013296760                      ; 0x78007878
	.long	3355494600                      ; 0xc800c8c8
	.long	2499843477                      ; 0x95009595
	.long	4177590777                      ; 0xf900f9f9
	.long	788541231                       ; 0x2f002f2f
	.long	3456159438                      ; 0xce00cece
	.long	3439381965                      ; 0xcd00cdcd
	.long	134219784                       ; 0x8000808
	.long	2046851706                      ; 0x7a007a7a
	.long	2281736328                      ; 0x88008888
	.long	939538488                       ; 0x38003838
	.long	1543527516                      ; 0x5c005c5c
	.long	2197848963                      ; 0x83008383
	.long	704653866                       ; 0x2a002a2a
	.long	671098920                       ; 0x28002828
	.long	1191200583                      ; 0x47004747
	.long	3674266587                      ; 0xdb00dbdb
	.long	3087055032                      ; 0xb800b8b8
	.long	3338717127                      ; 0xc700c7c7
	.long	2466288531                      ; 0x93009393
	.long	2751505572                      ; 0xa400a4a4
	.long	301994514                       ; 0x12001212
	.long	1392530259                      ; 0x53005353
	.long	4278255615                      ; 0xff00ffff
	.long	2264958855                      ; 0x87008787
	.long	234884622                       ; 0xe000e0e
	.long	822096177                       ; 0x31003131
	.long	905983542                       ; 0x36003636
	.long	553656609                       ; 0x21002121
	.long	1476417624                      ; 0x58005858
	.long	1207978056                      ; 0x48004848
	.long	16777473                        ; 0x1000101
	.long	2382401166                      ; 0x8e008e8e
	.long	922761015                       ; 0x37003737
	.long	1946186868                      ; 0x74007474
	.long	838873650                       ; 0x32003232
	.long	3389049546                      ; 0xca00caca
	.long	3909151209                      ; 0xe900e9e9
	.long	2969612721                      ; 0xb100b1b1
	.long	3070277559                      ; 0xb700b7b7
	.long	2868947883                      ; 0xab00abab
	.long	201329676                       ; 0xc000c0c
	.long	3607156695                      ; 0xd700d7d7
	.long	3288384708                      ; 0xc400c4c4
	.long	1442862678                      ; 0x56005656
	.long	1107313218                      ; 0x42004242
	.long	637543974                       ; 0x26002626
	.long	117442311                       ; 0x7000707
	.long	2550175896                      ; 0x98009898
	.long	1610637408                      ; 0x60006060
	.long	3640711641                      ; 0xd900d9d9
	.long	3053500086                      ; 0xb600b6b6
	.long	3103832505                      ; 0xb900b9b9
	.long	285217041                       ; 0x11001111
	.long	1073758272                      ; 0x40004040
	.long	3959483628                      ; 0xec00ecec
	.long	536879136                       ; 0x20002020
	.long	2348846220                      ; 0x8c008c8c
	.long	3170942397                      ; 0xbd00bdbd
	.long	2684395680                      ; 0xa000a0a0
	.long	3372272073                      ; 0xc900c9c9
	.long	2214626436                      ; 0x84008484
	.long	67109892                        ; 0x4000404
	.long	1224755529                      ; 0x49004949
	.long	587211555                       ; 0x23002323
	.long	4043370993                      ; 0xf100f1f1
	.long	1325420367                      ; 0x4f004f4f
	.long	1342197840                      ; 0x50005050
	.long	520101663                       ; 0x1f001f1f
	.long	318771987                       ; 0x13001313
	.long	3691044060                      ; 0xdc00dcdc
	.long	3623934168                      ; 0xd800d8d8
	.long	3221274816                      ; 0xc000c0c0
	.long	2650840734                      ; 0x9e009e9e
	.long	1459640151                      ; 0x57005757
	.long	3808486371                      ; 0xe300e3e3
	.long	3271607235                      ; 0xc300c3c3
	.long	2063629179                      ; 0x7b007b7b
	.long	1694524773                      ; 0x65006565
	.long	989870907                       ; 0x3b003b3b
	.long	33554946                        ; 0x2000202
	.long	2399178639                      ; 0x8f008f8f
	.long	1040203326                      ; 0x3e003e3e
	.long	3892373736                      ; 0xe800e8e8
	.long	620766501                       ; 0x25002525
	.long	2449511058                      ; 0x92009292
	.long	3842041317                      ; 0xe500e5e5
	.long	352326933                       ; 0x15001515
	.long	3707821533                      ; 0xdd00dddd
	.long	4244700669                      ; 0xfd00fdfd
	.long	385881879                       ; 0x17001717
	.long	2835392937                      ; 0xa900a9a9
	.long	3204497343                      ; 0xbf00bfbf
	.long	3556824276                      ; 0xd400d4d4
	.long	2583730842                      ; 0x9a009a9a
	.long	2113961598                      ; 0x7e007e7e
	.long	3305162181                      ; 0xc500c5c5
	.long	956315961                       ; 0x39003939
	.long	1728079719                      ; 0x67006767
	.long	4261478142                      ; 0xfe00fefe
	.long	1979741814                      ; 0x76007676
	.long	2634063261                      ; 0x9d009d9d
	.long	1124090691                      ; 0x43004343
	.long	2801837991                      ; 0xa700a7a7
	.long	3774931425                      ; 0xe100e1e1
	.long	3489714384                      ; 0xd000d0d0
	.long	4110480885                      ; 0xf500f5f5
	.long	1744857192                      ; 0x68006868
	.long	4060148466                      ; 0xf200f2f2
	.long	452991771                       ; 0x1b001b1b
	.long	872428596                       ; 0x34003434
	.long	1879076976                      ; 0x70007070
	.long	83887365                        ; 0x5000505
	.long	2734728099                      ; 0xa300a3a3
	.long	2315291274                      ; 0x8a008a8a
	.long	3573601749                      ; 0xd500d5d5
	.long	2030074233                      ; 0x79007979
	.long	2248181382                      ; 0x86008686
	.long	2818615464                      ; 0xa800a8a8
	.long	805318704                       ; 0x30003030
	.long	3321939654                      ; 0xc600c6c6
	.long	1358975313                      ; 0x51005151
	.long	1258310475                      ; 0x4b004b4b
	.long	503324190                       ; 0x1e001e1e
	.long	2785060518                      ; 0xa600a6a6
	.long	654321447                       ; 0x27002727
	.long	4127258358                      ; 0xf600f6f6
	.long	889206069                       ; 0x35003535
	.long	3523269330                      ; 0xd200d2d2
	.long	1845522030                      ; 0x6e006e6e
	.long	603989028                       ; 0x24002424
	.long	369104406                       ; 0x16001616
	.long	2181071490                      ; 0x82008282
	.long	1593859935                      ; 0x5f005f5f
	.long	3657489114                      ; 0xda00dada
	.long	3858818790                      ; 0xe600e6e6
	.long	1962964341                      ; 0x75007575
	.long	2717950626                      ; 0xa200a2a2
	.long	4009816047                      ; 0xef00efef
	.long	738208812                       ; 0x2c002c2c
	.long	2986390194                      ; 0xb200b2b2
	.long	469769244                       ; 0x1c001c1c
	.long	2667618207                      ; 0x9f009f9f
	.long	1560304989                      ; 0x5d005d5d
	.long	1862299503                      ; 0x6f006f6f
	.long	2147516544                      ; 0x80008080
	.long	167774730                       ; 0xa000a0a
	.long	1912631922                      ; 0x72007272
	.long	1140868164                      ; 0x44004444
	.long	2600508315                      ; 0x9b009b9b
	.long	1811967084                      ; 0x6c006c6c
	.long	2415956112                      ; 0x90009090
	.long	184552203                       ; 0xb000b0b
	.long	1526750043                      ; 0x5b005b5b
	.long	855651123                       ; 0x33003333
	.long	2097184125                      ; 0x7d007d7d
	.long	1509972570                      ; 0x5a005a5a
	.long	1375752786                      ; 0x52005252
	.long	4076925939                      ; 0xf300f3f3
	.long	1627414881                      ; 0x61006161
	.long	2701173153                      ; 0xa100a1a1
	.long	4144035831                      ; 0xf700f7f7
	.long	2952835248                      ; 0xb000b0b0
	.long	3590379222                      ; 0xd600d6d6
	.long	1056980799                      ; 0x3f003f3f
	.long	2080406652                      ; 0x7c007c7c
	.long	1828744557                      ; 0x6d006d6d
	.long	3976261101                      ; 0xed00eded
	.long	335549460                       ; 0x14001414
	.long	3758153952                      ; 0xe000e0e0
	.long	2768283045                      ; 0xa500a5a5
	.long	1023425853                      ; 0x3d003d3d
	.long	570434082                       ; 0x22002222
	.long	3003167667                      ; 0xb300b3b3
	.long	4160813304                      ; 0xf800f8f8
	.long	2298513801                      ; 0x89008989
	.long	3724599006                      ; 0xde00dede
	.long	1895854449                      ; 0x71007171
	.long	436214298                       ; 0x1a001a1a
	.long	2936057775                      ; 0xaf00afaf
	.long	3120609978                      ; 0xba00baba
	.long	3036722613                      ; 0xb500b5b5
	.long	2164294017                      ; 0x81008181

	.globl	_X1                             ; @X1
	.p2align	2, 0x0
_X1:
	.long	1381105746                      ; 0x52520052
	.long	151584777                       ; 0x9090009
	.long	1785331818                      ; 0x6a6a006a
	.long	3587506389                      ; 0xd5d500d5
	.long	808452144                       ; 0x30300030
	.long	909508662                       ; 0x36360036
	.long	2779054245                      ; 0xa5a500a5
	.long	943194168                       ; 0x38380038
	.long	3216965823                      ; 0xbfbf00bf
	.long	1077936192                      ; 0x40400040
	.long	2745368739                      ; 0xa3a300a3
	.long	2661154974                      ; 0x9e9e009e
	.long	2172715137                      ; 0x81810081
	.long	4092788979                      ; 0xf3f300f3
	.long	3621191895                      ; 0xd7d700d7
	.long	4227531003                      ; 0xfbfb00fb
	.long	2088501372                      ; 0x7c7c007c
	.long	3823304931                      ; 0xe3e300e3
	.long	960036921                       ; 0x39390039
	.long	2189557890                      ; 0x82820082
	.long	2610626715                      ; 0x9b9b009b
	.long	791609391                       ; 0x2f2f002f
	.long	4294902015                      ; 0xffff00ff
	.long	2273771655                      ; 0x87870087
	.long	875823156                       ; 0x34340034
	.long	2391670926                      ; 0x8e8e008e
	.long	1128464451                      ; 0x43430043
	.long	1145307204                      ; 0x44440044
	.long	3301179588                      ; 0xc4c400c4
	.long	3739091166                      ; 0xdede00de
	.long	3924361449                      ; 0xe9e900e9
	.long	3419078859                      ; 0xcbcb00cb
	.long	1414791252                      ; 0x54540054
	.long	2071658619                      ; 0x7b7b007b
	.long	2492727444                      ; 0x94940094
	.long	842137650                       ; 0x32320032
	.long	2795896998                      ; 0xa6a600a6
	.long	3267494082                      ; 0xc2c200c2
	.long	589496355                       ; 0x23230023
	.long	1027407933                      ; 0x3d3d003d
	.long	4008575214                      ; 0xeeee00ee
	.long	1280049228                      ; 0x4c4c004c
	.long	2509570197                      ; 0x95950095
	.long	185270283                       ; 0xb0b000b
	.long	1111621698                      ; 0x42420042
	.long	4210688250                      ; 0xfafa00fa
	.long	3284336835                      ; 0xc3c300c3
	.long	1313734734                      ; 0x4e4e004e
	.long	134742024                       ; 0x8080008
	.long	774766638                       ; 0x2e2e002e
	.long	2711683233                      ; 0xa1a100a1
	.long	1717960806                      ; 0x66660066
	.long	673710120                       ; 0x28280028
	.long	3654877401                      ; 0xd9d900d9
	.long	606339108                       ; 0x24240024
	.long	2998010034                      ; 0xb2b200b2
	.long	1987444854                      ; 0x76760076
	.long	1532690523                      ; 0x5b5b005b
	.long	2728525986                      ; 0xa2a200a2
	.long	1229520969                      ; 0x49490049
	.long	1835860077                      ; 0x6d6d006d
	.long	2341142667                      ; 0x8b8b008b
	.long	3520135377                      ; 0xd1d100d1
	.long	623181861                       ; 0x25250025
	.long	1920073842                      ; 0x72720072
	.long	4177002744                      ; 0xf8f800f8
	.long	4143317238                      ; 0xf6f600f6
	.long	1684275300                      ; 0x64640064
	.long	2256928902                      ; 0x86860086
	.long	1751646312                      ; 0x68680068
	.long	2560098456                      ; 0x98980098
	.long	370540566                       ; 0x16160016
	.long	3570663636                      ; 0xd4d400d4
	.long	2762211492                      ; 0xa4a400a4
	.long	1549533276                      ; 0x5c5c005c
	.long	3435921612                      ; 0xcccc00cc
	.long	1566376029                      ; 0x5d5d005d
	.long	1701118053                      ; 0x65650065
	.long	3065381046                      ; 0xb6b600b6
	.long	2459041938                      ; 0x92920092
	.long	1819017324                      ; 0x6c6c006c
	.long	1886388336                      ; 0x70700070
	.long	1212678216                      ; 0x48480048
	.long	1347420240                      ; 0x50500050
	.long	4261216509                      ; 0xfdfd00fd
	.long	3991732461                      ; 0xeded00ed
	.long	3115909305                      ; 0xb9b900b9
	.long	3671720154                      ; 0xdada00da
	.long	1583218782                      ; 0x5e5e005e
	.long	353697813                       ; 0x15150015
	.long	1178992710                      ; 0x46460046
	.long	1465319511                      ; 0x57570057
	.long	2812739751                      ; 0xa7a700a7
	.long	2374828173                      ; 0x8d8d008d
	.long	2644312221                      ; 0x9d9d009d
	.long	2223243396                      ; 0x84840084
	.long	2425356432                      ; 0x90900090
	.long	3638034648                      ; 0xd8d800d8
	.long	2880110763                      ; 0xabab00ab
	.long	0                               ; 0x0
	.long	2357985420                      ; 0x8c8c008c
	.long	3166437564                      ; 0xbcbc00bc
	.long	3553820883                      ; 0xd3d300d3
	.long	168427530                       ; 0xa0a000a
	.long	4160159991                      ; 0xf7f700f7
	.long	3840147684                      ; 0xe4e400e4
	.long	1482162264                      ; 0x58580058
	.long	84213765                        ; 0x5050005
	.long	3099066552                      ; 0xb8b800b8
	.long	3014852787                      ; 0xb3b300b3
	.long	1162149957                      ; 0x45450045
	.long	101056518                       ; 0x6060006
	.long	3503292624                      ; 0xd0d000d0
	.long	741081132                       ; 0x2c2c002c
	.long	505282590                       ; 0x1e1e001e
	.long	2408513679                      ; 0x8f8f008f
	.long	3402236106                      ; 0xcaca00ca
	.long	1061093439                      ; 0x3f3f003f
	.long	252641295                       ; 0xf0f000f
	.long	33685506                        ; 0x2020002
	.long	3250651329                      ; 0xc1c100c1
	.long	2947481775                      ; 0xafaf00af
	.long	3183280317                      ; 0xbdbd00bd
	.long	50528259                        ; 0x3030003
	.long	16842753                        ; 0x1010001
	.long	320012307                       ; 0x13130013
	.long	2324299914                      ; 0x8a8a008a
	.long	1802174571                      ; 0x6b6b006b
	.long	976879674                       ; 0x3a3a003a
	.long	2442199185                      ; 0x91910091
	.long	286326801                       ; 0x11110011
	.long	1094778945                      ; 0x41410041
	.long	1330577487                      ; 0x4f4f004f
	.long	1734803559                      ; 0x67670067
	.long	3705405660                      ; 0xdcdc00dc
	.long	3941204202                      ; 0xeaea00ea
	.long	2543255703                      ; 0x97970097
	.long	4075946226                      ; 0xf2f200f2
	.long	3486449871                      ; 0xcfcf00cf
	.long	3469607118                      ; 0xcece00ce
	.long	4042260720                      ; 0xf0f000f0
	.long	3031695540                      ; 0xb4b400b4
	.long	3873833190                      ; 0xe6e600e6
	.long	1936916595                      ; 0x73730073
	.long	2526412950                      ; 0x96960096
	.long	2896953516                      ; 0xacac00ac
	.long	1953759348                      ; 0x74740074
	.long	572653602                       ; 0x22220022
	.long	3890675943                      ; 0xe7e700e7
	.long	2913796269                      ; 0xadad00ad
	.long	892665909                       ; 0x35350035
	.long	2240086149                      ; 0x85850085
	.long	3806462178                      ; 0xe2e200e2
	.long	4193845497                      ; 0xf9f900f9
	.long	926351415                       ; 0x37370037
	.long	3907518696                      ; 0xe8e800e8
	.long	471597084                       ; 0x1c1c001c
	.long	1970602101                      ; 0x75750075
	.long	3755933919                      ; 0xdfdf00df
	.long	1852702830                      ; 0x6e6e006e
	.long	1195835463                      ; 0x47470047
	.long	4059103473                      ; 0xf1f100f1
	.long	437911578                       ; 0x1a1a001a
	.long	1903231089                      ; 0x71710071
	.long	488439837                       ; 0x1d1d001d
	.long	690552873                       ; 0x29290029
	.long	3318022341                      ; 0xc5c500c5
	.long	2307457161                      ; 0x89890089
	.long	1869545583                      ; 0x6f6f006f
	.long	3082223799                      ; 0xb7b700b7
	.long	1650589794                      ; 0x62620062
	.long	235798542                       ; 0xe0e000e
	.long	2863268010                      ; 0xaaaa00aa
	.long	404226072                       ; 0x18180018
	.long	3200123070                      ; 0xbebe00be
	.long	454754331                       ; 0x1b1b001b
	.long	4244373756                      ; 0xfcfc00fc
	.long	1448476758                      ; 0x56560056
	.long	1044250686                      ; 0x3e3e003e
	.long	1263206475                      ; 0x4b4b004b
	.long	3334865094                      ; 0xc6c600c6
	.long	3536978130                      ; 0xd2d200d2
	.long	2037973113                      ; 0x79790079
	.long	538968096                       ; 0x20200020
	.long	2593783962                      ; 0x9a9a009a
	.long	3688562907                      ; 0xdbdb00db
	.long	3233808576                      ; 0xc0c000c0
	.long	4278059262                      ; 0xfefe00fe
	.long	2021130360                      ; 0x78780078
	.long	3452764365                      ; 0xcdcd00cd
	.long	1515847770                      ; 0x5a5a005a
	.long	4109631732                      ; 0xf4f400f4
	.long	522125343                       ; 0x1f1f001f
	.long	3722248413                      ; 0xdddd00dd
	.long	2829582504                      ; 0xa8a800a8
	.long	858980403                       ; 0x33330033
	.long	2290614408                      ; 0x88880088
	.long	117899271                       ; 0x7070007
	.long	3351707847                      ; 0xc7c700c7
	.long	825294897                       ; 0x31310031
	.long	2981167281                      ; 0xb1b100b1
	.long	303169554                       ; 0x12120012
	.long	269484048                       ; 0x10100010
	.long	1499005017                      ; 0x59590059
	.long	656867367                       ; 0x27270027
	.long	2155872384                      ; 0x80800080
	.long	3974889708                      ; 0xecec00ec
	.long	1600061535                      ; 0x5f5f005f
	.long	1616904288                      ; 0x60600060
	.long	1364262993                      ; 0x51510051
	.long	2139029631                      ; 0x7f7f007f
	.long	2846425257                      ; 0xa9a900a9
	.long	421068825                       ; 0x19190019
	.long	3048538293                      ; 0xb5b500b5
	.long	1246363722                      ; 0x4a4a004a
	.long	218955789                       ; 0xd0d000d
	.long	757923885                       ; 0x2d2d002d
	.long	3856990437                      ; 0xe5e500e5
	.long	2054815866                      ; 0x7a7a007a
	.long	2677997727                      ; 0x9f9f009f
	.long	2475884691                      ; 0x93930093
	.long	3385393353                      ; 0xc9c900c9
	.long	2627469468                      ; 0x9c9c009c
	.long	4025417967                      ; 0xefef00ef
	.long	2694840480                      ; 0xa0a000a0
	.long	3772776672                      ; 0xe0e000e0
	.long	993722427                       ; 0x3b3b003b
	.long	1296891981                      ; 0x4d4d004d
	.long	2930639022                      ; 0xaeae00ae
	.long	707395626                       ; 0x2a2a002a
	.long	4126474485                      ; 0xf5f500f5
	.long	2964324528                      ; 0xb0b000b0
	.long	3368550600                      ; 0xc8c800c8
	.long	3958046955                      ; 0xebeb00eb
	.long	3149594811                      ; 0xbbbb00bb
	.long	1010565180                      ; 0x3c3c003c
	.long	2206400643                      ; 0x83830083
	.long	1397948499                      ; 0x53530053
	.long	2576941209                      ; 0x99990099
	.long	1633747041                      ; 0x61610061
	.long	387383319                       ; 0x17170017
	.long	724238379                       ; 0x2b2b002b
	.long	67371012                        ; 0x4040004
	.long	2122186878                      ; 0x7e7e007e
	.long	3132752058                      ; 0xbaba00ba
	.long	2004287607                      ; 0x77770077
	.long	3604349142                      ; 0xd6d600d6
	.long	640024614                       ; 0x26260026
	.long	3789619425                      ; 0xe1e100e1
	.long	1768489065                      ; 0x69690069
	.long	336855060                       ; 0x14140014
	.long	1667432547                      ; 0x63630063
	.long	1431634005                      ; 0x55550055
	.long	555810849                       ; 0x21210021
	.long	202113036                       ; 0xc0c000c
	.long	2105344125                      ; 0x7d7d007d

	.globl	_X2                             ; @X2
	.p2align	2, 0x0
_X2:
	.long	808464384                       ; 0x30303000
	.long	1751672832                      ; 0x68686800
	.long	2576980224                      ; 0x99999900
	.long	454761216                       ; 0x1b1b1b00
	.long	2273806080                      ; 0x87878700
	.long	3115956480                      ; 0xb9b9b900
	.long	555819264                       ; 0x21212100
	.long	2021160960                      ; 0x78787800
	.long	1347440640                      ; 0x50505000
	.long	960051456                       ; 0x39393900
	.long	3688618752                      ; 0xdbdbdb00
	.long	3789676800                      ; 0xe1e1e100
	.long	1920102912                      ; 0x72727200
	.long	151587072                       ; 0x9090900
	.long	1650614784                      ; 0x62626200
	.long	1010580480                      ; 0x3c3c3c00
	.long	1044266496                      ; 0x3e3e3e00
	.long	2122219008                      ; 0x7e7e7e00
	.long	1583242752                      ; 0x5e5e5e00
	.long	2391707136                      ; 0x8e8e8e00
	.long	4059164928                      ; 0xf1f1f100
	.long	2694881280                      ; 0xa0a0a000
	.long	3435973632                      ; 0xcccccc00
	.long	2745410304                      ; 0xa3a3a300
	.long	707406336                       ; 0x2a2a2a00
	.long	488447232                       ; 0x1d1d1d00
	.long	4227595008                      ; 0xfbfbfb00
	.long	3065427456                      ; 0xb6b6b600
	.long	3604403712                      ; 0xd6d6d600
	.long	538976256                       ; 0x20202000
	.long	3301229568                      ; 0xc4c4c400
	.long	2374864128                      ; 0x8d8d8d00
	.long	2172748032                      ; 0x81818100
	.long	1701143808                      ; 0x65656500
	.long	4126536960                      ; 0xf5f5f500
	.long	2307492096                      ; 0x89898900
	.long	3419130624                      ; 0xcbcbcb00
	.long	2644352256                      ; 0x9d9d9d00
	.long	2004317952                      ; 0x77777700
	.long	3334915584                      ; 0xc6c6c600
	.long	1465341696                      ; 0x57575700
	.long	1128481536                      ; 0x43434300
	.long	1448498688                      ; 0x56565600
	.long	387389184                       ; 0x17171700
	.long	3570717696                      ; 0xd4d4d400
	.long	1077952512                      ; 0x40404000
	.long	437918208                       ; 0x1a1a1a00
	.long	1296911616                      ; 0x4d4d4d00
	.long	3233857536                      ; 0xc0c0c000
	.long	1667457792                      ; 0x63636300
	.long	1819044864                      ; 0x6c6c6c00
	.long	3823362816                      ; 0xe3e3e300
	.long	3082270464                      ; 0xb7b7b700
	.long	3368601600                      ; 0xc8c8c800
	.long	1684300800                      ; 0x64646400
	.long	1785358848                      ; 0x6a6a6a00
	.long	1397969664                      ; 0x53535300
	.long	2863311360                      ; 0xaaaaaa00
	.long	943208448                       ; 0x38383800
	.long	2560137216                      ; 0x98989800
	.long	202116096                       ; 0xc0c0c00
	.long	4109693952                      ; 0xf4f4f400
	.long	2610666240                      ; 0x9b9b9b00
	.long	3991792896                      ; 0xededed00
	.long	2139062016                      ; 0x7f7f7f00
	.long	572662272                       ; 0x22222200
	.long	1987474944                      ; 0x76767600
	.long	2947526400                      ; 0xafafaf00
	.long	3722304768                      ; 0xdddddd00
	.long	976894464                       ; 0x3a3a3a00
	.long	185273088                       ; 0xb0b0b00
	.long	1482184704                      ; 0x58585800
	.long	1734829824                      ; 0x67676700
	.long	2290649088                      ; 0x88888800
	.long	101058048                       ; 0x6060600
	.long	3284386560                      ; 0xc3c3c300
	.long	892679424                       ; 0x35353500
	.long	218959104                       ; 0xd0d0d00
	.long	16843008                        ; 0x1010100
	.long	2341178112                      ; 0x8b8b8b00
	.long	2358021120                      ; 0x8c8c8c00
	.long	3267543552                      ; 0xc2c2c200
	.long	3873891840                      ; 0xe6e6e600
	.long	1600085760                      ; 0x5f5f5f00
	.long	33686016                        ; 0x2020200
	.long	606348288                       ; 0x24242400
	.long	1970631936                      ; 0x75757500
	.long	2475922176                      ; 0x93939300
	.long	1717986816                      ; 0x66666600
	.long	505290240                       ; 0x1e1e1e00
	.long	3857048832                      ; 0xe5e5e500
	.long	3806519808                      ; 0xe2e2e200
	.long	1414812672                      ; 0x54545400
	.long	3638089728                      ; 0xd8d8d800
	.long	269488128                       ; 0x10101000
	.long	3469659648                      ; 0xcecece00
	.long	2054846976                      ; 0x7a7a7a00
	.long	3907577856                      ; 0xe8e8e800
	.long	134744064                       ; 0x8080800
	.long	741092352                       ; 0x2c2c2c00
	.long	303174144                       ; 0x12121200
	.long	2543294208                      ; 0x97979700
	.long	842150400                       ; 0x32323200
	.long	2880154368                      ; 0xababab00
	.long	3031741440                      ; 0xb4b4b400
	.long	656877312                       ; 0x27272700
	.long	168430080                       ; 0xa0a0a00
	.long	589505280                       ; 0x23232300
	.long	3755990784                      ; 0xdfdfdf00
	.long	4025478912                      ; 0xefefef00
	.long	3402287616                      ; 0xcacaca00
	.long	3654932736                      ; 0xd9d9d900
	.long	3099113472                      ; 0xb8b8b800
	.long	4210752000                      ; 0xfafafa00
	.long	3705461760                      ; 0xdcdcdc00
	.long	825307392                       ; 0x31313100
	.long	1802201856                      ; 0x6b6b6b00
	.long	3520188672                      ; 0xd1d1d100
	.long	2913840384                      ; 0xadadad00
	.long	421075200                       ; 0x19191900
	.long	1229539584                      ; 0x49494900
	.long	3183328512                      ; 0xbdbdbd00
	.long	1364283648                      ; 0x51515100
	.long	2526451200                      ; 0x96969600
	.long	4008635904                      ; 0xeeeeee00
	.long	3840205824                      ; 0xe4e4e400
	.long	2829625344                      ; 0xa8a8a800
	.long	1094795520                      ; 0x41414100
	.long	3671775744                      ; 0xdadada00
	.long	4294967040                      ; 0xffffff00
	.long	3452816640                      ; 0xcdcdcd00
	.long	1431655680                      ; 0x55555500
	.long	2256963072                      ; 0x86868600
	.long	909522432                       ; 0x36363600
	.long	3200171520                      ; 0xbebebe00
	.long	1633771776                      ; 0x61616100
	.long	1381126656                      ; 0x52525200
	.long	4177065984                      ; 0xf8f8f800
	.long	3149642496                      ; 0xbbbbbb00
	.long	235802112                       ; 0xe0e0e00
	.long	2189591040                      ; 0x82828200
	.long	1212696576                      ; 0x48484800
	.long	1768515840                      ; 0x69696900
	.long	2593823232                      ; 0x9a9a9a00
	.long	3772833792                      ; 0xe0e0e000
	.long	1195853568                      ; 0x47474700
	.long	2661195264                      ; 0x9e9e9e00
	.long	1549556736                      ; 0x5c5c5c00
	.long	67372032                        ; 0x4040400
	.long	1263225600                      ; 0x4b4b4b00
	.long	875836416                       ; 0x34343400
	.long	353703168                       ; 0x15151500
	.long	2038003968                      ; 0x79797900
	.long	640034304                       ; 0x26262600
	.long	2812782336                      ; 0xa7a7a700
	.long	3739147776                      ; 0xdedede00
	.long	690563328                       ; 0x29292900
	.long	2930683392                      ; 0xaeaeae00
	.long	2459079168                      ; 0x92929200
	.long	3621246720                      ; 0xd7d7d700
	.long	2223277056                      ; 0x84848400
	.long	3924420864                      ; 0xe9e9e900
	.long	3537031680                      ; 0xd2d2d200
	.long	3132799488                      ; 0xbababa00
	.long	1566399744                      ; 0x5d5d5d00
	.long	4092850944                      ; 0xf3f3f300
	.long	3318072576                      ; 0xc5c5c500
	.long	2964369408                      ; 0xb0b0b000
	.long	3217014528                      ; 0xbfbfbf00
	.long	2762253312                      ; 0xa4a4a400
	.long	993737472                       ; 0x3b3b3b00
	.long	1903259904                      ; 0x71717100
	.long	1145324544                      ; 0x44444400
	.long	1179010560                      ; 0x46464600
	.long	724249344                       ; 0x2b2b2b00
	.long	4244438016                      ; 0xfcfcfc00
	.long	3958106880                      ; 0xebebeb00
	.long	1869573888                      ; 0x6f6f6f00
	.long	3587560704                      ; 0xd5d5d500
	.long	4143379968                      ; 0xf6f6f600
	.long	336860160                       ; 0x14141400
	.long	4278124032                      ; 0xfefefe00
	.long	2088532992                      ; 0x7c7c7c00
	.long	1886416896                      ; 0x70707000
	.long	1515870720                      ; 0x5a5a5a00
	.long	2105376000                      ; 0x7d7d7d00
	.long	4261281024                      ; 0xfdfdfd00
	.long	791621376                       ; 0x2f2f2f00
	.long	404232192                       ; 0x18181800
	.long	2206434048                      ; 0x83838300
	.long	370546176                       ; 0x16161600
	.long	2779096320                      ; 0xa5a5a500
	.long	2442236160                      ; 0x91919100
	.long	522133248                       ; 0x1f1f1f00
	.long	84215040                        ; 0x5050500
	.long	2509608192                      ; 0x95959500
	.long	1953788928                      ; 0x74747400
	.long	2846468352                      ; 0xa9a9a900
	.long	3250700544                      ; 0xc1c1c100
	.long	1532713728                      ; 0x5b5b5b00
	.long	1246382592                      ; 0x4a4a4a00
	.long	2240120064                      ; 0x85858500
	.long	1835887872                      ; 0x6d6d6d00
	.long	320017152                       ; 0x13131300
	.long	117901056                       ; 0x7070700
	.long	1330597632                      ; 0x4f4f4f00
	.long	1313754624                      ; 0x4e4e4e00
	.long	1162167552                      ; 0x45454500
	.long	2998055424                      ; 0xb2b2b200
	.long	252645120                       ; 0xf0f0f00
	.long	3385444608                      ; 0xc9c9c900
	.long	471604224                       ; 0x1c1c1c00
	.long	2795939328                      ; 0xa6a6a600
	.long	3166485504                      ; 0xbcbcbc00
	.long	3974949888                      ; 0xececec00
	.long	1936945920                      ; 0x73737300
	.long	2425393152                      ; 0x90909000
	.long	2071689984                      ; 0x7b7b7b00
	.long	3486502656                      ; 0xcfcfcf00
	.long	1499027712                      ; 0x59595900
	.long	2408550144                      ; 0x8f8f8f00
	.long	2711724288                      ; 0xa1a1a100
	.long	4193908992                      ; 0xf9f9f900
	.long	757935360                       ; 0x2d2d2d00
	.long	4076007936                      ; 0xf2f2f200
	.long	2981212416                      ; 0xb1b1b100
	.long	0                               ; 0x0
	.long	2492765184                      ; 0x94949400
	.long	926365440                       ; 0x37373700
	.long	2678038272                      ; 0x9f9f9f00
	.long	3503345664                      ; 0xd0d0d000
	.long	774778368                       ; 0x2e2e2e00
	.long	2627509248                      ; 0x9c9c9c00
	.long	1852730880                      ; 0x6e6e6e00
	.long	673720320                       ; 0x28282800
	.long	1061109504                      ; 0x3f3f3f00
	.long	2155905024                      ; 0x80808000
	.long	4042321920                      ; 0xf0f0f000
	.long	1027423488                      ; 0x3d3d3d00
	.long	3553874688                      ; 0xd3d3d300
	.long	623191296                       ; 0x25252500
	.long	2324335104                      ; 0x8a8a8a00
	.long	3048584448                      ; 0xb5b5b500
	.long	3890734848                      ; 0xe7e7e700
	.long	1111638528                      ; 0x42424200
	.long	3014898432                      ; 0xb3b3b300
	.long	3351758592                      ; 0xc7c7c700
	.long	3941263872                      ; 0xeaeaea00
	.long	4160222976                      ; 0xf7f7f700
	.long	1280068608                      ; 0x4c4c4c00
	.long	286331136                       ; 0x11111100
	.long	858993408                       ; 0x33333300
	.long	50529024                        ; 0x3030300
	.long	2728567296                      ; 0xa2a2a200
	.long	2896997376                      ; 0xacacac00
	.long	1616928768                      ; 0x60606000

.subsections_via_symbols
