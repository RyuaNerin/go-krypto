	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 14, 0
	.section	__TEXT,__literal16,16byte_literals
	.p2align	4, 0x0                          ; -- Begin function lsh256_neon_init
lCPI0_0:
	.long	3475206790                      ; 0xcf237286
	.long	3993835303                      ; 0xee0d1727
	.long	862152085                       ; 0x33636595
	.long	2344144991                      ; 0x8bb8d05f
lCPI0_1:
	.long	1887352057                      ; 0x707eb4f9
	.long	4133173346                      ; 0xf65b3862
	.long	1795893950                      ; 0x6b0b2abe
	.long	1454959626                      ; 0x56b8ec0a
lCPI0_2:
	.long	3183741608                      ; 0xbdc40aa8
	.long	516557672                       ; 0x1eca0b68
	.long	3659172286                      ; 0xda1a89be
	.long	826790740                       ; 0x3147d354
lCPI0_3:
	.long	109447379                       ; 0x68608d3
	.long	1658386343                      ; 0x62d8f7a7
	.long	3613807275                      ; 0xd76652ab
	.long	1281362499                      ; 0x4c600a43
lCPI0_4:
	.long	2152805754                      ; 0x8051357a
	.long	327575752                       ; 0x138668c8
	.long	1202340996                      ; 0x47aa4484
	.long	3759864641                      ; 0xe01afb41
lCPI0_5:
	.long	274551672                       ; 0x105d5378
	.long	796188244                       ; 0x2f74de54
	.long	1546595733                      ; 0x5c2f2d95
	.long	4065673150                      ; 0xf2553fbe
lCPI0_6:
	.long	855914637                       ; 0x3304388d
	.long	2968888263                      ; 0xb0f5a3c7
	.long	3009438148                      ; 0xb36061c4
	.long	2061227347                      ; 0x7adbd553
lCPI0_7:
	.long	1184960287                      ; 0x46a10f1f
	.long	4259112070                      ; 0xfddce486
	.long	3021226920                      ; 0xb41443a8
	.long	428764061                       ; 0x198e6b9d
	.section	__TEXT,__text,regular,pure_instructions
	.globl	_lsh256_neon_init
	.p2align	2
_lsh256_neon_init:                      ; @lsh256_neon_init
; %bb.0:
	stp	x29, x30, [sp, #-16]!           ; 16-byte Folded Spill
	mov	x29, sp
	str	w1, [x0]
	str	wzr, [x0, #16]
	cmp	x1, #28
	b.ne	LBB0_2
; %bb.1:
Lloh0:
	adrp	x8, lCPI0_0@PAGE
Lloh1:
	ldr	q0, [x8, lCPI0_0@PAGEOFF]
Lloh2:
	adrp	x8, lCPI0_1@PAGE
Lloh3:
	ldr	q1, [x8, lCPI0_1@PAGEOFF]
Lloh4:
	adrp	x8, lCPI0_2@PAGE
Lloh5:
	ldr	q2, [x8, lCPI0_2@PAGEOFF]
Lloh6:
	adrp	x8, lCPI0_3@PAGE
Lloh7:
	ldr	q3, [x8, lCPI0_3@PAGEOFF]
	stp	q3, q2, [x0, #32]
	stp	q1, q0, [x0, #64]
	ldp	x29, x30, [sp], #16             ; 16-byte Folded Reload
	ret
LBB0_2:
Lloh8:
	adrp	x8, lCPI0_4@PAGE
Lloh9:
	ldr	q0, [x8, lCPI0_4@PAGEOFF]
Lloh10:
	adrp	x8, lCPI0_5@PAGE
Lloh11:
	ldr	q1, [x8, lCPI0_5@PAGEOFF]
Lloh12:
	adrp	x8, lCPI0_6@PAGE
Lloh13:
	ldr	q2, [x8, lCPI0_6@PAGEOFF]
Lloh14:
	adrp	x8, lCPI0_7@PAGE
Lloh15:
	ldr	q3, [x8, lCPI0_7@PAGEOFF]
	stp	q3, q2, [x0, #32]
	stp	q1, q0, [x0, #64]
	ldp	x29, x30, [sp], #16             ; 16-byte Folded Reload
	ret
	.loh AdrpLdr	Lloh6, Lloh7
	.loh AdrpAdrp	Lloh4, Lloh6
	.loh AdrpLdr	Lloh4, Lloh5
	.loh AdrpAdrp	Lloh2, Lloh4
	.loh AdrpLdr	Lloh2, Lloh3
	.loh AdrpAdrp	Lloh0, Lloh2
	.loh AdrpLdr	Lloh0, Lloh1
	.loh AdrpLdr	Lloh14, Lloh15
	.loh AdrpAdrp	Lloh12, Lloh14
	.loh AdrpLdr	Lloh12, Lloh13
	.loh AdrpAdrp	Lloh10, Lloh12
	.loh AdrpLdr	Lloh10, Lloh11
	.loh AdrpAdrp	Lloh8, Lloh10
	.loh AdrpLdr	Lloh8, Lloh9
                                        ; -- End function
	.section	__TEXT,__literal16,16byte_literals
	.p2align	4, 0x0                          ; -- Begin function lsh256_neon_update
lCPI1_0:
	.long	0                               ; 0x0
	.long	8                               ; 0x8
	.long	16                              ; 0x10
	.long	24                              ; 0x18
lCPI1_1:
	.long	4294967264                      ; 0xffffffe0
	.long	4294967272                      ; 0xffffffe8
	.long	4294967280                      ; 0xfffffff0
	.long	4294967288                      ; 0xfffffff8
lCPI1_2:
	.long	24                              ; 0x18
	.long	16                              ; 0x10
	.long	8                               ; 0x8
	.long	0                               ; 0x0
lCPI1_3:
	.long	4294967288                      ; 0xfffffff8
	.long	4294967280                      ; 0xfffffff0
	.long	4294967272                      ; 0xffffffe8
	.long	4294967264                      ; 0xffffffe0
	.section	__TEXT,__text,regular,pure_instructions
	.globl	_lsh256_neon_update
	.p2align	2
_lsh256_neon_update:                    ; @lsh256_neon_update
; %bb.0:
	stp	x29, x30, [sp, #-16]!           ; 16-byte Folded Spill
	mov	x29, sp
	ldr	w14, [x0, #16]
	add	x8, x14, x2
	cmp	x8, #127
	b.hi	LBB1_8
; %bb.1:
	cmp	w2, #1
	b.lt	LBB1_7
; %bb.2:
	and	x8, x2, #0xffffffff
	cmp	x8, #8
	b.lo	LBB1_4
; %bb.3:
	add	x9, x0, x14
	sub	x9, x9, x1
	add	x9, x9, #96
	cmp	x9, #64
	b.hs	LBB1_31
LBB1_4:
	mov	x9, #0                          ; =0x0
LBB1_5:
	sub	x8, x8, x9
	add	x10, x9, x14
	add	x10, x10, x0
	add	x10, x10, #96
	add	x9, x1, x9
LBB1_6:                                 ; =>This Inner Loop Header: Depth=1
	ldrb	w11, [x9], #1
	strb	w11, [x10], #1
	subs	x8, x8, #1
	b.ne	LBB1_6
LBB1_7:
	add	w2, w14, w2
	b	LBB1_29
LBB1_8:
	adrp	x9, lCPI1_0@PAGE
	adrp	x10, lCPI1_1@PAGE
	adrp	x11, lCPI1_2@PAGE
	adrp	x12, lCPI1_3@PAGE
Lloh16:
	adrp	x8, _STEP256@PAGE
Lloh17:
	add	x8, x8, _STEP256@PAGEOFF
	cbz	w14, LBB1_18
; %bb.9:
	mov	w13, #128                       ; =0x80
	sub	x13, x13, x14
	cmp	w13, #1
	b.lt	LBB1_15
; %bb.10:
	and	x15, x13, #0xffffffff
	cmp	x15, #8
	b.lo	LBB1_12
; %bb.11:
	add	x16, x0, x14
	sub	x16, x16, x1
	add	x16, x16, #96
	cmp	x16, #64
	b.hs	LBB1_35
LBB1_12:
	mov	x16, #0                         ; =0x0
LBB1_13:
	sub	x15, x15, x16
	add	x14, x16, x14
	add	x14, x14, x0
	add	x14, x14, #96
	add	x16, x1, x16
LBB1_14:                                ; =>This Inner Loop Header: Depth=1
	ldrb	w17, [x16], #1
	strb	w17, [x14], #1
	subs	x15, x15, #1
	b.ne	LBB1_14
LBB1_15:
	mov	x14, #0                         ; =0x0
	ldp	q3, q1, [x0, #96]
	ldp	q2, q0, [x0, #128]
	ldp	q16, q18, [x0, #160]
	ldp	q4, q5, [x0, #192]
	ldp	q22, q23, [x0, #32]
	ldp	q20, q21, [x0, #64]
	ldr	q6, [x9, lCPI1_0@PAGEOFF]
	ldr	q7, [x10, lCPI1_1@PAGEOFF]
	ldr	q17, [x11, lCPI1_2@PAGEOFF]
	ldr	q19, [x12, lCPI1_3@PAGEOFF]
LBB1_16:                                ; =>This Inner Loop Header: Depth=1
	add	x15, x8, x14
	ldp	q24, q25, [x15]
	eor.16b	v22, v22, v3
	eor.16b	v23, v23, v1
	eor.16b	v20, v20, v2
	eor.16b	v21, v21, v0
	add.4s	v22, v20, v22
	add.4s	v23, v21, v23
	shl.4s	v26, v22, #29
	sri.4s	v26, v22, #3
	shl.4s	v22, v23, #29
	sri.4s	v22, v23, #3
	eor.16b	v23, v26, v24
	eor.16b	v22, v22, v25
	add.4s	v20, v23, v20
	add.4s	v21, v22, v21
	add.4s	v24, v20, v20
	sri.4s	v24, v20, #31
	add.4s	v20, v21, v21
	sri.4s	v20, v21, #31
	add.4s	v21, v24, v23
	add.4s	v22, v20, v22
	ushl.4s	v23, v24, v6
	ushl.4s	v24, v24, v7
	eor.16b	v23, v24, v23
	ushl.4s	v24, v20, v17
	ushl.4s	v20, v20, v19
	eor.16b	v20, v20, v24
	ext.16b	v24, v21, v21, #8
	rev64.2s	v24, v24
	ext.8b	v25, v21, v24, #4
	ext.8b	v21, v24, v21, #4
	mov.d	v21[1], v25[0]
	ext.16b	v24, v22, v22, #8
	rev64.2s	v24, v24
	ext.8b	v25, v22, v24, #4
	ext.8b	v22, v24, v22, #4
	mov.d	v22[1], v25[0]
	ext.16b	v23, v23, v23, #12
	rev64.4s	v23, v23
	ext.16b	v20, v20, v20, #12
	rev64.4s	v20, v20
	ext.16b	v24, v3, v3, #8
	rev64.2s	v24, v24
	mov.d	v24[1], v3[0]
	ext.16b	v3, v2, v2, #8
	rev64.2s	v25, v3
	mov.d	v25[1], v2[0]
	ext.16b	v1, v1, v1, #12
	ext.16b	v0, v0, v0, #12
	add.4s	v3, v16, v24
	add.4s	v1, v18, v1
	add.4s	v2, v4, v25
	add.4s	v0, v5, v0
	ldp	q24, q25, [x15, #32]
	eor.16b	v22, v22, v16
	eor.16b	v20, v20, v18
	eor.16b	v21, v21, v4
	eor.16b	v23, v23, v5
	add.4s	v22, v22, v21
	add.4s	v20, v20, v23
	shl.4s	v26, v22, #5
	sri.4s	v26, v22, #27
	shl.4s	v22, v20, #5
	sri.4s	v22, v20, #27
	eor.16b	v20, v26, v24
	eor.16b	v22, v22, v25
	add.4s	v21, v20, v21
	add.4s	v23, v22, v23
	shl.4s	v24, v21, #17
	sri.4s	v24, v21, #15
	shl.4s	v21, v23, #17
	sri.4s	v21, v23, #15
	add.4s	v20, v24, v20
	add.4s	v22, v21, v22
	ushl.4s	v23, v24, v6
	ushl.4s	v24, v24, v7
	eor.16b	v23, v24, v23
	ushl.4s	v24, v21, v17
	ushl.4s	v21, v21, v19
	eor.16b	v24, v21, v24
	ext.16b	v21, v20, v20, #8
	rev64.2s	v21, v21
	ext.8b	v25, v20, v21, #4
	ext.8b	v20, v21, v20, #4
	mov.d	v20[1], v25[0]
	ext.16b	v21, v22, v22, #8
	rev64.2s	v21, v21
	ext.8b	v25, v22, v21, #4
	ext.8b	v22, v21, v22, #4
	mov.d	v22[1], v25[0]
	ext.16b	v21, v23, v23, #12
	rev64.4s	v21, v21
	ext.16b	v23, v24, v24, #12
	rev64.4s	v23, v23
	ext.16b	v24, v16, v16, #8
	rev64.2s	v24, v24
	mov.d	v24[1], v16[0]
	ext.16b	v18, v18, v18, #12
	ext.16b	v16, v4, v4, #8
	rev64.2s	v25, v16
	mov.d	v25[1], v4[0]
	ext.16b	v5, v5, v5, #12
	add.4s	v16, v24, v3
	add.4s	v18, v1, v18
	add.4s	v4, v25, v2
	add.4s	v5, v5, v0
	add	x14, x14, #64
	cmp	w14, #832
	b.ne	LBB1_16
; %bb.17:
	eor.16b	v3, v22, v3
	eor.16b	v1, v23, v1
	stp	q3, q1, [x0, #32]
	eor.16b	v1, v20, v2
	eor.16b	v0, v21, v0
	stp	q1, q0, [x0, #64]
	add	x1, x1, x13
	sub	x2, x2, x13
	str	wzr, [x0, #16]
LBB1_18:
	cmp	x2, #128
	b.lo	LBB1_23
; %bb.19:
	ldp	q21, q23, [x0, #32]
	ldp	q20, q22, [x0, #64]
	ldr	q0, [x9, lCPI1_0@PAGEOFF]
	ldr	q1, [x10, lCPI1_1@PAGEOFF]
	ldr	q2, [x11, lCPI1_2@PAGEOFF]
	ldr	q3, [x12, lCPI1_3@PAGEOFF]
LBB1_20:                                ; =>This Loop Header: Depth=1
                                        ;     Child Loop BB1_21 Depth 2
	mov	x9, #0                          ; =0x0
	ldp	q7, q5, [x1]
	ldp	q6, q4, [x1, #32]
	ldp	q18, q19, [x1, #64]
	ldp	q16, q17, [x1, #96]
LBB1_21:                                ;   Parent Loop BB1_20 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	add	x10, x8, x9
	ldp	q24, q25, [x10]
	eor.16b	v21, v21, v7
	eor.16b	v23, v23, v5
	eor.16b	v20, v20, v6
	eor.16b	v22, v22, v4
	add.4s	v21, v20, v21
	add.4s	v23, v22, v23
	shl.4s	v26, v21, #29
	sri.4s	v26, v21, #3
	shl.4s	v21, v23, #29
	sri.4s	v21, v23, #3
	eor.16b	v23, v26, v24
	eor.16b	v21, v21, v25
	add.4s	v20, v23, v20
	add.4s	v22, v21, v22
	add.4s	v24, v20, v20
	sri.4s	v24, v20, #31
	add.4s	v20, v22, v22
	sri.4s	v20, v22, #31
	add.4s	v22, v24, v23
	add.4s	v21, v20, v21
	ushl.4s	v23, v24, v0
	ushl.4s	v24, v24, v1
	eor.16b	v23, v24, v23
	ushl.4s	v24, v20, v2
	ushl.4s	v20, v20, v3
	eor.16b	v20, v20, v24
	ext.16b	v24, v22, v22, #8
	rev64.2s	v24, v24
	ext.8b	v25, v22, v24, #4
	ext.8b	v22, v24, v22, #4
	mov.d	v22[1], v25[0]
	ext.16b	v24, v21, v21, #8
	rev64.2s	v24, v24
	ext.8b	v25, v21, v24, #4
	ext.8b	v21, v24, v21, #4
	mov.d	v21[1], v25[0]
	ext.16b	v23, v23, v23, #12
	rev64.4s	v23, v23
	ext.16b	v20, v20, v20, #12
	rev64.4s	v20, v20
	ext.16b	v24, v7, v7, #8
	rev64.2s	v24, v24
	mov.d	v24[1], v7[0]
	ext.16b	v7, v6, v6, #8
	rev64.2s	v25, v7
	mov.d	v25[1], v6[0]
	ext.16b	v5, v5, v5, #12
	ext.16b	v4, v4, v4, #12
	add.4s	v7, v18, v24
	add.4s	v5, v19, v5
	add.4s	v6, v16, v25
	add.4s	v4, v17, v4
	ldp	q24, q25, [x10, #32]
	eor.16b	v21, v21, v18
	eor.16b	v20, v20, v19
	eor.16b	v22, v22, v16
	eor.16b	v23, v23, v17
	add.4s	v21, v21, v22
	add.4s	v20, v20, v23
	shl.4s	v26, v21, #5
	sri.4s	v26, v21, #27
	shl.4s	v21, v20, #5
	sri.4s	v21, v20, #27
	eor.16b	v20, v26, v24
	eor.16b	v21, v21, v25
	add.4s	v22, v20, v22
	add.4s	v23, v21, v23
	shl.4s	v24, v22, #17
	sri.4s	v24, v22, #15
	shl.4s	v22, v23, #17
	sri.4s	v22, v23, #15
	add.4s	v20, v24, v20
	add.4s	v21, v22, v21
	ushl.4s	v23, v24, v0
	ushl.4s	v24, v24, v1
	eor.16b	v23, v24, v23
	ushl.4s	v24, v22, v2
	ushl.4s	v22, v22, v3
	eor.16b	v24, v22, v24
	ext.16b	v22, v20, v20, #8
	rev64.2s	v22, v22
	ext.8b	v25, v20, v22, #4
	ext.8b	v20, v22, v20, #4
	mov.d	v20[1], v25[0]
	ext.16b	v22, v21, v21, #8
	rev64.2s	v22, v22
	ext.8b	v25, v21, v22, #4
	ext.8b	v21, v22, v21, #4
	mov.d	v21[1], v25[0]
	ext.16b	v22, v23, v23, #12
	rev64.4s	v22, v22
	ext.16b	v23, v24, v24, #12
	rev64.4s	v23, v23
	ext.16b	v24, v18, v18, #8
	rev64.2s	v24, v24
	mov.d	v24[1], v18[0]
	ext.16b	v19, v19, v19, #12
	ext.16b	v18, v16, v16, #8
	rev64.2s	v25, v18
	mov.d	v25[1], v16[0]
	ext.16b	v17, v17, v17, #12
	add.4s	v18, v24, v7
	add.4s	v19, v5, v19
	add.4s	v16, v25, v6
	add.4s	v17, v17, v4
	add	x9, x9, #64
	cmp	w9, #832
	b.ne	LBB1_21
; %bb.22:                               ;   in Loop: Header=BB1_20 Depth=1
	eor.16b	v21, v21, v7
	eor.16b	v23, v23, v5
	stp	q21, q23, [x0, #32]
	eor.16b	v20, v20, v6
	eor.16b	v22, v22, v4
	stp	q20, q22, [x0, #64]
	add	x1, x1, #128
	sub	x2, x2, #128
	cmp	x2, #127
	b.hi	LBB1_20
LBB1_23:
	cbz	x2, LBB1_30
; %bb.24:
	cmp	x2, #8
	b.lo	LBB1_26
; %bb.25:
	sub	x8, x0, x1
	add	x8, x8, #96
	cmp	x8, #64
	b.hs	LBB1_33
LBB1_26:
	mov	x8, #0                          ; =0x0
LBB1_27:
	sub	x9, x2, x8
	add	x10, x8, x0
	add	x10, x10, #96
	add	x8, x1, x8
LBB1_28:                                ; =>This Inner Loop Header: Depth=1
	ldrb	w11, [x8], #1
	strb	w11, [x10], #1
	subs	x9, x9, #1
	b.ne	LBB1_28
LBB1_29:
	str	w2, [x0, #16]
LBB1_30:
	ldp	x29, x30, [sp], #16             ; 16-byte Folded Reload
	ret
LBB1_31:
	cmp	x8, #64
	b.hs	LBB1_37
; %bb.32:
	mov	x9, #0                          ; =0x0
	b	LBB1_41
LBB1_33:
	cmp	x2, #64
	b.hs	LBB1_44
; %bb.34:
	mov	x8, #0                          ; =0x0
	b	LBB1_48
LBB1_35:
	cmp	x15, #64
	b.hs	LBB1_51
; %bb.36:
	mov	x16, #0                         ; =0x0
	b	LBB1_55
LBB1_37:
	and	x10, x2, #0x3f
	sub	x9, x8, x10
	add	x11, x14, x0
	add	x11, x11, #144
	add	x12, x1, #32
	mov	x13, x9
LBB1_38:                                ; =>This Inner Loop Header: Depth=1
	ldp	q0, q1, [x12, #-32]
	ldp	q2, q3, [x12], #64
	stp	q0, q1, [x11, #-48]
	stp	q2, q3, [x11, #-16]
	add	x11, x11, #64
	subs	x13, x13, #64
	b.ne	LBB1_38
; %bb.39:
	cbz	x10, LBB1_7
; %bb.40:
	cmp	x10, #8
	b.lo	LBB1_5
LBB1_41:
	mov	x13, x9
	and	x10, x2, #0x7
	add	x11, x1, x9
	sub	x9, x8, x10
	add	x12, x13, x14
	add	x12, x12, x0
	add	x12, x12, #96
	add	x13, x13, x10
	sub	x13, x13, x8
LBB1_42:                                ; =>This Inner Loop Header: Depth=1
	ldr	d0, [x11], #8
	str	d0, [x12], #8
	adds	x13, x13, #8
	b.ne	LBB1_42
; %bb.43:
	cbnz	x10, LBB1_5
	b	LBB1_7
LBB1_44:
	and	x8, x2, #0x40
	add	x9, x1, #32
	add	x10, x0, #144
	mov	x11, x8
LBB1_45:                                ; =>This Inner Loop Header: Depth=1
	ldp	q0, q1, [x9, #-32]
	ldp	q2, q3, [x9], #64
	stp	q0, q1, [x10, #-48]
	stp	q2, q3, [x10, #-16]
	add	x10, x10, #64
	subs	x11, x11, #64
	b.ne	LBB1_45
; %bb.46:
	cmp	x2, x8
	b.eq	LBB1_29
; %bb.47:
	tst	x2, #0x38
	b.eq	LBB1_27
LBB1_48:
	mov	x11, x8
	and	x8, x2, #0x78
	add	x9, x1, x11
	add	x10, x11, x0
	add	x10, x10, #96
	sub	x11, x11, x8
LBB1_49:                                ; =>This Inner Loop Header: Depth=1
	ldr	d0, [x9], #8
	str	d0, [x10], #8
	adds	x11, x11, #8
	b.ne	LBB1_49
; %bb.50:
	cmp	x2, x8
	b.eq	LBB1_29
	b	LBB1_27
LBB1_51:
	and	x17, x13, #0x3f
	sub	x16, x15, x17
	add	x3, x14, x0
	add	x3, x3, #144
	add	x4, x1, #32
	mov	x5, x16
LBB1_52:                                ; =>This Inner Loop Header: Depth=1
	ldp	q0, q1, [x4, #-32]
	ldp	q2, q3, [x4], #64
	stp	q0, q1, [x3, #-48]
	stp	q2, q3, [x3, #-16]
	add	x3, x3, #64
	subs	x5, x5, #64
	b.ne	LBB1_52
; %bb.53:
	cbz	x17, LBB1_15
; %bb.54:
	cmp	x17, #8
	b.lo	LBB1_13
LBB1_55:
	mov	x5, x16
	and	x17, x13, #0x7
	add	x3, x1, x16
	sub	x16, x15, x17
	add	x4, x5, x14
	add	x4, x4, x0
	add	x4, x4, #96
	add	x5, x5, x17
	sub	x5, x5, x15
LBB1_56:                                ; =>This Inner Loop Header: Depth=1
	ldr	d0, [x3], #8
	str	d0, [x4], #8
	adds	x5, x5, #8
	b.ne	LBB1_56
; %bb.57:
	cbnz	x17, LBB1_13
	b	LBB1_15
	.loh AdrpAdd	Lloh16, Lloh17
                                        ; -- End function
	.section	__TEXT,__literal16,16byte_literals
	.p2align	4, 0x0                          ; -- Begin function lsh256_neon_final
lCPI2_0:
	.long	0                               ; 0x0
	.long	8                               ; 0x8
	.long	16                              ; 0x10
	.long	24                              ; 0x18
lCPI2_1:
	.long	4294967264                      ; 0xffffffe0
	.long	4294967272                      ; 0xffffffe8
	.long	4294967280                      ; 0xfffffff0
	.long	4294967288                      ; 0xfffffff8
lCPI2_2:
	.long	24                              ; 0x18
	.long	16                              ; 0x10
	.long	8                               ; 0x8
	.long	0                               ; 0x0
lCPI2_3:
	.long	4294967288                      ; 0xfffffff8
	.long	4294967280                      ; 0xfffffff0
	.long	4294967272                      ; 0xffffffe8
	.long	4294967264                      ; 0xffffffe0
	.section	__TEXT,__text,regular,pure_instructions
	.globl	_lsh256_neon_final
	.p2align	2
_lsh256_neon_final:                     ; @lsh256_neon_final
; %bb.0:
	stp	x29, x30, [sp, #-16]!           ; 16-byte Folded Spill
	mov	x29, sp
	ldr	w9, [x0, #16]
	add	x8, x0, x9
	mov	w10, #128                       ; =0x80
	strb	w10, [x8, #96]
	add	w8, w9, #1
	cmp	w8, #127
	b.hi	LBB2_14
; %bb.1:
	mov	w10, #127                       ; =0x7f
	sub	w9, w10, w9
	cmp	w9, #8
	b.hs	LBB2_3
; %bb.2:
	mov	x10, #0                         ; =0x0
	b	LBB2_12
LBB2_3:
	cmp	w9, #64
	b.hs	LBB2_5
; %bb.4:
	mov	x10, #0                         ; =0x0
	b	LBB2_9
LBB2_5:
	and	x10, x9, #0xffffffc0
	add	x11, x8, x0
	add	x11, x11, #144
	movi.2d	v0, #0000000000000000
	mov	x12, x10
LBB2_6:                                 ; =>This Inner Loop Header: Depth=1
	stp	q0, q0, [x11, #-48]
	stp	q0, q0, [x11, #-16]
	add	x11, x11, #64
	subs	x12, x12, #64
	b.ne	LBB2_6
; %bb.7:
	cmp	x10, x9
	b.eq	LBB2_14
; %bb.8:
	tst	x9, #0x38
	b.eq	LBB2_12
LBB2_9:
	mov	x12, x10
	and	x10, x9, #0xfffffff8
	add	x11, x12, x8
	add	x11, x11, x0
	add	x11, x11, #96
	sub	x12, x12, x10
	movi.2d	v0, #0000000000000000
LBB2_10:                                ; =>This Inner Loop Header: Depth=1
	str	d0, [x11], #8
	adds	x12, x12, #8
	b.ne	LBB2_10
; %bb.11:
	cmp	x10, x9
	b.eq	LBB2_14
LBB2_12:
	sub	x9, x9, x10
	add	x8, x10, x8
	add	x8, x8, x0
	add	x8, x8, #96
LBB2_13:                                ; =>This Inner Loop Header: Depth=1
	strb	wzr, [x8], #1
	subs	x9, x9, #1
	b.ne	LBB2_13
LBB2_14:
	mov	x8, #0                          ; =0x0
	ldp	q3, q0, [x0, #96]
	ldp	q2, q1, [x0, #128]
	ldp	q6, q16, [x0, #160]
	ldp	q4, q5, [x0, #192]
	ldp	q20, q21, [x0, #32]
Lloh18:
	adrp	x9, lCPI2_0@PAGE
Lloh19:
	ldr	q7, [x9, lCPI2_0@PAGEOFF]
Lloh20:
	adrp	x9, lCPI2_1@PAGE
Lloh21:
	ldr	q17, [x9, lCPI2_1@PAGEOFF]
Lloh22:
	adrp	x9, lCPI2_2@PAGE
Lloh23:
	ldr	q18, [x9, lCPI2_2@PAGEOFF]
Lloh24:
	adrp	x9, lCPI2_3@PAGE
Lloh25:
	ldr	q19, [x9, lCPI2_3@PAGEOFF]
Lloh26:
	adrp	x9, _STEP256@PAGE
Lloh27:
	add	x9, x9, _STEP256@PAGEOFF
	ldp	q22, q23, [x0, #64]
LBB2_15:                                ; =>This Inner Loop Header: Depth=1
	add	x10, x9, x8
	ldp	q24, q25, [x10]
	eor.16b	v20, v20, v3
	eor.16b	v21, v21, v0
	eor.16b	v22, v22, v2
	eor.16b	v23, v23, v1
	add.4s	v20, v22, v20
	add.4s	v21, v23, v21
	shl.4s	v26, v20, #29
	sri.4s	v26, v20, #3
	shl.4s	v20, v21, #29
	sri.4s	v20, v21, #3
	eor.16b	v21, v26, v24
	eor.16b	v20, v20, v25
	add.4s	v22, v21, v22
	add.4s	v23, v20, v23
	add.4s	v24, v22, v22
	sri.4s	v24, v22, #31
	add.4s	v22, v23, v23
	sri.4s	v22, v23, #31
	add.4s	v21, v24, v21
	add.4s	v20, v22, v20
	ushl.4s	v23, v24, v7
	ushl.4s	v24, v24, v17
	eor.16b	v23, v24, v23
	ushl.4s	v24, v22, v18
	ushl.4s	v22, v22, v19
	eor.16b	v22, v22, v24
	ext.16b	v24, v21, v21, #8
	rev64.2s	v24, v24
	ext.8b	v25, v21, v24, #4
	ext.8b	v21, v24, v21, #4
	mov.d	v21[1], v25[0]
	ext.16b	v24, v20, v20, #8
	rev64.2s	v24, v24
	ext.8b	v25, v20, v24, #4
	ext.8b	v20, v24, v20, #4
	mov.d	v20[1], v25[0]
	ext.16b	v23, v23, v23, #12
	rev64.4s	v23, v23
	ext.16b	v22, v22, v22, #12
	rev64.4s	v22, v22
	ext.16b	v24, v3, v3, #8
	rev64.2s	v24, v24
	mov.d	v24[1], v3[0]
	ext.16b	v3, v2, v2, #8
	rev64.2s	v25, v3
	mov.d	v25[1], v2[0]
	ext.16b	v0, v0, v0, #12
	ext.16b	v1, v1, v1, #12
	add.4s	v3, v6, v24
	add.4s	v0, v16, v0
	add.4s	v2, v4, v25
	add.4s	v1, v5, v1
	ldp	q24, q25, [x10, #32]
	eor.16b	v20, v20, v6
	eor.16b	v22, v22, v16
	eor.16b	v21, v21, v4
	eor.16b	v23, v23, v5
	add.4s	v20, v20, v21
	add.4s	v22, v22, v23
	shl.4s	v26, v20, #5
	sri.4s	v26, v20, #27
	shl.4s	v20, v22, #5
	sri.4s	v20, v22, #27
	eor.16b	v22, v26, v24
	eor.16b	v20, v20, v25
	add.4s	v21, v22, v21
	add.4s	v23, v20, v23
	shl.4s	v24, v21, #17
	sri.4s	v24, v21, #15
	shl.4s	v21, v23, #17
	sri.4s	v21, v23, #15
	add.4s	v22, v24, v22
	add.4s	v20, v21, v20
	ushl.4s	v23, v24, v7
	ushl.4s	v24, v24, v17
	eor.16b	v23, v24, v23
	ushl.4s	v24, v21, v18
	ushl.4s	v21, v21, v19
	eor.16b	v21, v21, v24
	ext.16b	v24, v22, v22, #8
	rev64.2s	v24, v24
	ext.8b	v25, v22, v24, #4
	ext.8b	v22, v24, v22, #4
	mov.d	v22[1], v25[0]
	ext.16b	v24, v20, v20, #8
	rev64.2s	v24, v24
	ext.8b	v25, v20, v24, #4
	ext.8b	v20, v24, v20, #4
	mov.d	v20[1], v25[0]
	ext.16b	v23, v23, v23, #12
	rev64.4s	v23, v23
	ext.16b	v21, v21, v21, #12
	rev64.4s	v21, v21
	ext.16b	v24, v6, v6, #8
	rev64.2s	v24, v24
	mov.d	v24[1], v6[0]
	ext.16b	v16, v16, v16, #12
	ext.16b	v6, v4, v4, #8
	rev64.2s	v25, v6
	mov.d	v25[1], v4[0]
	ext.16b	v5, v5, v5, #12
	add.4s	v6, v24, v3
	add.4s	v16, v0, v16
	add.4s	v4, v25, v2
	add.4s	v5, v5, v1
	add	x8, x8, #64
	cmp	w8, #832
	b.ne	LBB2_15
; %bb.16:
	eor.16b	v2, v22, v2
	eor.16b	v1, v23, v1
	stp	q2, q1, [x0, #64]
	eor3.16b	v2, v20, v3, v2
	eor3.16b	v0, v21, v0, v1
	stp	q2, q0, [x0, #32]
	str	q2, [x1]
	ldr	q0, [x0, #48]
	str	q0, [x1, #16]
	ldp	x29, x30, [sp], #16             ; 16-byte Folded Reload
	ret
	.loh AdrpAdd	Lloh26, Lloh27
	.loh AdrpAdrp	Lloh24, Lloh26
	.loh AdrpLdr	Lloh24, Lloh25
	.loh AdrpAdrp	Lloh22, Lloh24
	.loh AdrpLdr	Lloh22, Lloh23
	.loh AdrpAdrp	Lloh20, Lloh22
	.loh AdrpLdr	Lloh20, Lloh21
	.loh AdrpAdrp	Lloh18, Lloh20
	.loh AdrpLdr	Lloh18, Lloh19
                                        ; -- End function
	.section	__TEXT,__const
	.p2align	5, 0x0                          ; @STEP256
_STEP256:
	.long	2440867728                      ; 0x917caf90
	.long	1813713058                      ; 0x6c1b10a2
	.long	1865754947                      ; 0x6f352943
	.long	3480715843                      ; 0xcf778243
	.long	753628274                       ; 0x2ceb7472
	.long	703164402                       ; 0x29e96ff2
	.long	2325455912                      ; 0x8a9ba428
	.long	787162690                       ; 0x2eeb2642
	.long	237781025                       ; 0xe2c4021
	.long	2267788046                      ; 0x872bb30e
	.long	2757651634                      ; 0xa45e6cb2
	.long	1190774290                      ; 0x46f9c612
	.long	408938142                       ; 0x185fe69e
	.long	324624923                       ; 0x1359621b
	.long	641715378                       ; 0x263fccb2
	.long	437348464                       ; 0x1a116870
	.long	980181295                       ; 0x3a6c612f
	.long	3000942997                      ; 0xb2dec195
	.long	46866262                        ; 0x2cb1f56
	.long	1086314584                      ; 0x40bfd858
	.long	2017887414                      ; 0x784684b6
	.long	1824226606                      ; 0x6cbb7d2e
	.long	1712094936                      ; 0x660c7ed8
	.long	729405578                       ; 0x2b79d88a
	.long	2798489705                      ; 0xa6cd9069
	.long	2443204423                      ; 0x91a05747
	.long	3454694744                      ; 0xcdea7558
	.long	9973912                         ; 0x983098
	.long	3200989998                      ; 0xbecb3b2e
	.long	674802586                       ; 0x2838ab9a
	.long	1921734462                      ; 0x728b573e
	.long	2773639861                      ; 0xa55262b5
	.long	1952315919                      ; 0x745dfa0f
	.long	838311640                       ; 0x31f79ed8
	.long	3093286437                      ; 0xb85fce25
	.long	2563295384                      ; 0x98c8c898
	.long	2315676140                      ; 0x8a0669ec
	.long	1625572802                      ; 0x60e445c2
	.long	4259485104                      ; 0xfde295b0
	.long	4155840602                      ; 0xf7b5185a
	.long	3528984963                      ; 0xd2580983
	.long	697726729                       ; 0x29967709
	.long	405664733                       ; 0x182df3dd
	.long	1636917552                      ; 0x61916130
	.long	2423281270                      ; 0x90705676
	.long	1160382498                      ; 0x452a0822
	.long	3765978797                      ; 0xe07846ad
	.long	2899145553                      ; 0xaccd7351
	.long	711036245                       ; 0x2a618d55
	.long	3222110258                      ; 0xc00d8032
	.long	1176621301                      ; 0x4621d0f5
	.long	4075983249                      ; 0xf2f29191
	.long	13028614                        ; 0xc6cd06
	.long	1865558631                      ; 0x6f322a67
	.long	1488909453                      ; 0x58bef48d
	.long	2051065085                      ; 0x7a40c4fd
	.long	2347688575                      ; 0x8beee27f
	.long	3448615666                      ; 0xcd8db2f2
	.long	1743963707                      ; 0x67f2c63b
	.long	3850642307                      ; 0xe5842383
	.long	3348353798                      ; 0xc793d306
	.long	2707198422                      ; 0xa15c91d6
	.long	397640165                       ; 0x17b381e5
	.long	3137716855                      ; 0xbb05c277
	.long	2060542474                      ; 0x7ad1620a
	.long	1530963391                      ; 0x5b40a5bf
	.long	1522074018                      ; 0x5ab901a2
	.long	1772595048                      ; 0x69a7a768
	.long	1533467085                      ; 0x5b66d9cd
	.long	4260259959                      ; 0xfdee6877
	.long	3409274620                      ; 0xcb3566fc
	.long	3234347570                      ; 0xc0c83a32
	.long	1278438532                      ; 0x4c336c84
	.long	2615567642                      ; 0x9be6651a
	.long	330998780                       ; 0x13baa3fc
	.long	290394065                       ; 0x114f0fd1
	.long	3259017000                      ; 0xc240a728
	.long	3965116532                      ; 0xec56e074
	.long	10249159                        ; 0x9c63c7
	.long	2298637554                      ; 0x89026cf2
	.long	2141188304                      ; 0x7f9ff0d0
	.long	2185985973                      ; 0x824b7fb5
	.long	3462307855                      ; 0xce5ea00f
	.long	1616830690                      ; 0x605ee0e2
	.long	48746474                        ; 0x2e7cfea
	.long	1127699808                      ; 0x43375560
	.long	2634033863                      ; 0x9d002ac7
	.long	2339331963                      ; 0x8b6f5f7b
	.long	529580367                       ; 0x1f90c14f
	.long	3452646711                      ; 0xcdcb3537
	.long	754888669                       ; 0x2cfeafdd
	.long	3208627010                      ; 0xbf3fc342
	.long	3937909228                      ; 0xeab7b9ec
	.long	2056041891                      ; 0x7a8cb5a3
	.long	2636837476                      ; 0x9d2af264
	.long	4207860486                      ; 0xfacedb06
	.long	2958168174                      ; 0xb052106e
	.long	2566941956                      ; 0x99006d04
	.long	732859657                       ; 0x2bae8d09
	.long	4278388225                      ; 0xff030601
	.long	2725357270                      ; 0xa271a6d6
	.long	121788701                       ; 0x742591d
	.long	3357366017                      ; 0xc81d5701
	.long	3383353856                      ; 0xc9a9e200
	.long	40009502                        ; 0x2627f1e
	.long	2574086557                      ; 0x996d719d
	.long	3661338164                      ; 0xda3b9634
	.long	34146304                        ; 0x2090800
	.long	337149304                       ; 0x14187d78
	.long	1234925092                      ; 0x499b7624
	.long	3849607369                      ; 0xe57458c9
	.long	1938547401                      ; 0x738be2c9
	.long	1692507424                      ; 0x64e19d20
	.long	115281718                       ; 0x6df0f36
	.long	366070542                       ; 0x15d1cb0e
	.long	185665538                       ; 0xb110802
	.long	748025228                       ; 0x2c95f58c
	.long	3843136109                      ; 0xe5119a6d
	.long	1506615982                      ; 0x59cd22ae
	.long	4285443132                      ; 0xff6eac3c
	.long	1182711172                      ; 0x467ebd84
	.long	3857597756                      ; 0xe5ee453c
	.long	3885816099                      ; 0xe79cd923
	.long	471403021                       ; 0x1c190a0d
	.long	3263922616                      ; 0xc28b81b8
	.long	4138469458                      ; 0xf6ac0852
	.long	653250823                       ; 0x26efd107
	.long	1847257403                      ; 0x6e1ae93b
	.long	3309060554                      ; 0xc53c41ca
	.long	3560145441                      ; 0xd4338221
	.long	2222325002                      ; 0x8475fd0a
	.long	891492137                       ; 0x35231729
	.long	1309489786                      ; 0x4e0d3a7a
	.long	2729728840                      ; 0xa2b45b48
	.long	381737005                       ; 0x16c0d82d
	.long	2298750121                      ; 0x890424a9
	.long	25037967                        ; 0x17e0c8f
	.long	129344501                       ; 0x7b5a3f5
	.long	4201842574                      ; 0xfa73078e
	.long	1480212574                      ; 0x583a405e
	.long	1531425992                      ; 0x5b47b4c8
	.long	1460642794                      ; 0x570fa3ea
	.long	3617129795                      ; 0xd7990543
	.long	2368261682                      ; 0x8d28ce32
	.long	2139790224                      ; 0x7f8a9b90
	.long	3176765692                      ; 0xbd5998fc
	.long	1836750472                      ; 0x6d7a9688
	.long	2457509558                      ; 0x927a9eb6
	.long	2734456099                      ; 0xa2fc7d23
	.long	1723043393                      ; 0x66b38e41
	.long	1889421594                      ; 0x709e491a
	.long	3052863679                      ; 0xb5f700bf
	.long	170273807                       ; 0xa262c0f
	.long	384996793                       ; 0x16f295b9
	.long	3893436149                      ; 0xe8111ef5
	.long	219764040                       ; 0xd195548
	.long	2675548357                      ; 0x9f79a0c5
	.long	440520615                       ; 0x1a41cfa7
	.long	250045322                       ; 0xee7638a
	.long	2901917812                      ; 0xacf7c074
	.long	810695449                       ; 0x30523b19
	.long	159928015                       ; 0x9884ecf
	.long	4180677853                      ; 0xf93014dd
	.long	644783445                       ; 0x266e9d55
	.long	421160548                       ; 0x191a6664
	.long	1544648385                      ; 0x5c1176c1
	.long	4132105624                      ; 0xf64aed98
	.long	2763535648                      ; 0xa4b83520
	.long	2190300233                      ; 0x828d5449
	.long	2446794200                      ; 0x91d71dd8
	.long	692384470                       ; 0x2944f2d6
	.long	2500588155                      ; 0x950bf27b
	.long	864078461                       ; 0x3380ca7d
	.long	1837643805                      ; 0x6d88381d
	.long	1094223502                      ; 0x4138868e
	.long	1559057860                      ; 0x5ced55c4
	.long	266444235                       ; 0xfe19dcb
	.long	1760884329                      ; 0x68f4f669
	.long	1849149695                      ; 0x6e37c8ff
	.long	2701028880                      ; 0xa0fe6e10
	.long	3024832432                      ; 0xb44b47b0
	.long	4123022730                      ; 0xf5c0558a
	.long	2042565839                      ; 0x79bf14cf
	.long	1245911584                      ; 0x4a431a20
	.long	4051658970                      ; 0xf17f68da
	.long	1575706577                      ; 0x5deb5fd1
	.long	2785069165                      ; 0xa600c86d
	.long	2674687664                      ; 0x9f6c7eb0
	.long	4287821924                      ; 0xff92f864
	.long	3054887039                      ; 0xb615e07f
	.long	953410632                       ; 0x38d3e448
	.long	2371697258                      ; 0x8d5d3a6a
	.long	1894269899                      ; 0x70e843cb
	.long	1229664558                      ; 0x494b312e
	.long	2798204435                      ; 0xa6c93613
	.long	199962447                       ; 0xbeb2f4f
	.long	2458606947                      ; 0x928b5d63
	.long	3421921333                      ; 0xcbf66035
	.long	213396608                       ; 0xcb82c80
	.long	3935806711                      ; 0xea97a4f7
	.long	1496059707                      ; 0x592c0f3b
	.long	2491178871                      ; 0x947c5f77
	.long	1879001529                      ; 0x6fff49b9
	.long	4145708634                      ; 0xf71a7e5a
	.long	501793013                       ; 0x1de8c0f5
	.long	3260454400                      ; 0xc2569600
	.long	3303320716                      ; 0xc4e4ac8c
	.long	2185010401                      ; 0x823c9ce1

.subsections_via_symbols
