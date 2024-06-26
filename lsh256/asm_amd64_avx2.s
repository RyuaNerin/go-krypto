//go:build amd64 && !purego && (!gccgo || go1.18)
// +build amd64
// +build !purego
// +build !gccgo go1.18

// GENERATED BY C2GOASM
// EDITED BY RYUANERIN
// -- DO NOT EDIT

#include "textflag.h"

TEXT ·__lsh256_avx2_init(SB), NOSPLIT, $0-16
	MOVQ	ctx+0(FP), DI
	MOVQ	algtype+8(FP), SI

	MOVQ	SI, 0(DI)
	MOVQ	$0, 16(DI)

	CMPQ	SI, $32
	JEQ 	LSH256_INIT_256

LSH256_INIT_224:
	MOVD	·iv224(SB), BP
	JMP		LSH256_INIT_RET

LSH256_INIT_256:
	MOVD	·iv256(SB), BP

LSH256_INIT_RET:
	VMOVAPS	 0(BP), Y0
	VMOVAPS	32(BP), Y1

	VMOVUPS	Y0, 32(DI)
	VMOVUPS	Y1, 64(DI)

	VZEROUPPER
	RET

DATA LCDATA2<>+0x000(SB)/8, $0x6c1b10a2917caf90
DATA LCDATA2<>+0x008(SB)/8, $0xcf7782436f352943
DATA LCDATA2<>+0x010(SB)/8, $0x29e96ff22ceb7472
DATA LCDATA2<>+0x018(SB)/8, $0x2eeb26428a9ba428
DATA LCDATA2<>+0x020(SB)/8, $0x0c0f0e0d03020100
DATA LCDATA2<>+0x028(SB)/8, $0x0605040709080b0a
DATA LCDATA2<>+0x030(SB)/8, $0x0f0e0d0c00030201
DATA LCDATA2<>+0x038(SB)/8, $0x050407060a09080b
DATA LCDATA2<>+0x040(SB)/8, $0x872bb30e0e2c4021
DATA LCDATA2<>+0x048(SB)/8, $0x46f9c612a45e6cb2
DATA LCDATA2<>+0x050(SB)/8, $0x1359621b185fe69e
DATA LCDATA2<>+0x058(SB)/8, $0x1a116870263fccb2
DATA LCDATA2<>+0x060(SB)/8, $0x0000000200000003
DATA LCDATA2<>+0x068(SB)/8, $0x0000000100000000
DATA LCDATA2<>+0x070(SB)/8, $0x0000000400000007
DATA LCDATA2<>+0x078(SB)/8, $0x0000000600000005
GLOBL LCDATA2<>(SB), RODATA|NOPTR, $128

TEXT ·__lsh256_avx2_update(SB), NOSPLIT, $0-32
	MOVQ ctx+0(FP), DI
	MOVQ data_base+8(FP), SI
	MOVQ data_len+16(FP), DX
	//   data_cap+24

	LEAQ LCDATA2<>(SB), BP

	WORD $0x478b; BYTE $0x10                   // mov    eax, dword [rdi + 16]
	LONG $0x100c8d48                           // lea    rcx, [rax + rdx]
	LONG $0x7ff98348                           // cmp    rcx, 127
	JA   LBB1_4
	WORD $0xd285                               // test    edx, edx
	JLE  LBB1_28
	LONG $0x07048d4c                           // lea    r8, [rdi + rax]
	LONG $0x60c08349                           // add    r8, 96
	WORD $0xd189                               // mov    ecx, edx
	LONG $0x10f98348                           // cmp    rcx, 16
	JB   LBB1_3
	LONG $0x070c8d4c                           // lea    r9, [rdi + rax]
	WORD $0x2949; BYTE $0xf1                   // sub    r9, rsi
	LONG $0x60c18349                           // add    r9, 96
	LONG $0x80f98149; WORD $0x0000; BYTE $0x00 // cmp    r9, 128
	JAE  LBB1_12

LBB1_3:
	WORD $0x3145; BYTE $0xc9 // xor    r9d, r9d

LBB1_23:
	WORD $0x8941; BYTE $0xd3 // mov    r11d, edx
	WORD $0x2945; BYTE $0xcb // sub    r11d, r9d
	WORD $0x894d; BYTE $0xca // mov    r10, r9
	WORD $0xf749; BYTE $0xd2 // not    r10
	WORD $0x0149; BYTE $0xca // add    r10, rcx
	LONG $0x03e38341         // and    r11d, 3
	JE   LBB1_25

LBB1_24:
	LONG $0x1cb60f42; BYTE $0x0e // movzx    ebx, byte [rsi + r9]
	LONG $0x081c8843             // mov    byte [r8 + r9], bl
	WORD $0xff49; BYTE $0xc1     // inc    r9
	WORD $0xff49; BYTE $0xcb     // dec    r11
	JNE  LBB1_24

LBB1_25:
	LONG $0x03fa8349 // cmp    r10, 3
	JB   LBB1_28
	LONG $0x38048d4c // lea    r8, [rax + rdi]
	LONG $0x63c08349 // add    r8, 99

LBB1_27:
	LONG $0x14b60f46; BYTE $0x0e   // movzx    r10d, byte [rsi + r9]
	LONG $0x08548847; BYTE $0xfd   // mov    byte [r8 + r9 - 3], r10b
	LONG $0x54b60f46; WORD $0x010e // movzx    r10d, byte [rsi + r9 + 1]
	LONG $0x08548847; BYTE $0xfe   // mov    byte [r8 + r9 - 2], r10b
	LONG $0x54b60f46; WORD $0x020e // movzx    r10d, byte [rsi + r9 + 2]
	LONG $0x08548847; BYTE $0xff   // mov    byte [r8 + r9 - 1], r10b
	LONG $0x54b60f46; WORD $0x030e // movzx    r10d, byte [rsi + r9 + 3]
	LONG $0x08148847               // mov    byte [r8 + r9], r10b
	LONG $0x04c18349               // add    r9, 4
	WORD $0x394c; BYTE $0xc9       // cmp    rcx, r9
	JNE  LBB1_27

LBB1_28:
	WORD $0xd001 // add    eax, edx
	JMP  LBB1_62

LBB1_4:
	LONG $0x476ffec5; BYTE $0x20               // vmovdqu    ymm0, yword [rdi + 32]
	LONG $0x4f6ffec5; BYTE $0x40               // vmovdqu    ymm1, yword [rdi + 64]
	WORD $0x8548; BYTE $0xc0                   // test    rax, rax
	JE   LBB1_44
	LONG $0x0080b841; WORD $0x0000             // mov    r8d, 128
	WORD $0x2941; BYTE $0xc0                   // sub    r8d, eax
	WORD $0x8944; BYTE $0xc1                   // mov    ecx, r8d
	WORD $0x8545; BYTE $0xc0                   // test    r8d, r8d
	JLE  LBB1_41
	LONG $0x070c8d4c                           // lea    r9, [rdi + rax]
	LONG $0x60c18349                           // add    r9, 96
	LONG $0x10f88341                           // cmp    r8d, 16
	JB   LBB1_7
	LONG $0x07048d4c                           // lea    r8, [rdi + rax]
	WORD $0x2949; BYTE $0xf0                   // sub    r8, rsi
	LONG $0x60c08349                           // add    r8, 96
	LONG $0x80f88149; WORD $0x0000; BYTE $0x00 // cmp    r8, 128
	JAE  LBB1_14

LBB1_7:
	WORD $0x3145; BYTE $0xc0 // xor    r8d, r8d

LBB1_36:
	WORD $0x894d; BYTE $0xc2 // mov    r10, r8
	WORD $0xf749; BYTE $0xd2 // not    r10
	WORD $0x0149; BYTE $0xca // add    r10, rcx
	WORD $0x8949; BYTE $0xcb // mov    r11, rcx
	LONG $0x03e38349         // and    r11, 3
	JE   LBB1_38

LBB1_37:
	LONG $0x1cb60f42; BYTE $0x06 // movzx    ebx, byte [rsi + r8]
	LONG $0x011c8843             // mov    byte [r9 + r8], bl
	WORD $0xff49; BYTE $0xc0     // inc    r8
	WORD $0xff49; BYTE $0xcb     // dec    r11
	JNE  LBB1_37

LBB1_38:
	LONG $0x03fa8349         // cmp    r10, 3
	JB   LBB1_41
	WORD $0x0148; BYTE $0xf8 // add    rax, rdi
	LONG $0x63c08348         // add    rax, 99

LBB1_40:
	LONG $0x0cb60f46; BYTE $0x06   // movzx    r9d, byte [rsi + r8]
	LONG $0x004c8846; BYTE $0xfd   // mov    byte [rax + r8 - 3], r9b
	LONG $0x4cb60f46; WORD $0x0106 // movzx    r9d, byte [rsi + r8 + 1]
	LONG $0x004c8846; BYTE $0xfe   // mov    byte [rax + r8 - 2], r9b
	LONG $0x4cb60f46; WORD $0x0206 // movzx    r9d, byte [rsi + r8 + 2]
	LONG $0x004c8846; BYTE $0xff   // mov    byte [rax + r8 - 1], r9b
	LONG $0x4cb60f46; WORD $0x0306 // movzx    r9d, byte [rsi + r8 + 3]
	LONG $0x000c8846               // mov    byte [rax + r8], r9b
	LONG $0x04c08349               // add    r8, 4
	WORD $0x394c; BYTE $0xc1       // cmp    rcx, r8
	JNE  LBB1_40

LBB1_41:
	LONG $0x5f6ffec5; BYTE $0x60   // vmovdqu    ymm3, yword [rdi + 96]
	QUAD $0x00000080976ffec5       // vmovdqu    ymm2, yword [rdi + 128]
	QUAD $0x000000a0af6ffec5       // vmovdqu    ymm5, yword [rdi + 160]
	QUAD $0x000000c0a76ffec5       // vmovdqu    ymm4, yword [rdi + 192]
	LONG $0xc0efe5c5               // vpxor    ymm0, ymm3, ymm0
	LONG $0xc9efedc5               // vpxor    ymm1, ymm2, ymm1
	LONG $0xc0fef5c5               // vpaddd    ymm0, ymm1, ymm0
	LONG $0xd072cdc5; BYTE $0x03   // vpsrld    ymm6, ymm0, 3
	LONG $0xf072fdc5; BYTE $0x1d   // vpslld    ymm0, ymm0, 29
	LONG $0xc6ebfdc5               // vpor    ymm0, ymm0, ymm6
	LONG $0x45effdc5; BYTE $0x00   // vpxor    ymm0, ymm0, yword 0[rbp]
	LONG $0xc9fefdc5               // vpaddd    ymm1, ymm0, ymm1
	LONG $0xd172cdc5; BYTE $0x1f   // vpsrld    ymm6, ymm1, 31
	LONG $0xc9fef5c5               // vpaddd    ymm1, ymm1, ymm1
	LONG $0xceebf5c5               // vpor    ymm1, ymm1, ymm6
	LONG $0xc0fef5c5               // vpaddd    ymm0, ymm1, ymm0
	LONG $0xf070fdc5; BYTE $0xd2   // vpshufd    ymm6, ymm0, 210
	LONG $0x456ffdc5; BYTE $0x20   // vmovdqa    ymm0, yword 32[rbp]
	LONG $0x0075e2c4; BYTE $0xc8   // vpshufb    ymm1, ymm1, ymm0
	LONG $0x464de3c4; WORD $0x31f9 // vperm2i128    ymm7, ymm6, ymm1, 49
	LONG $0x384de3c4; WORD $0x01c9 // vinserti128    ymm1, ymm6, xmm1, 1
	LONG $0xf5efc5c5               // vpxor    ymm6, ymm7, ymm5
	LONG $0xcceff5c5               // vpxor    ymm1, ymm1, ymm4
	LONG $0xf1fecdc5               // vpaddd    ymm6, ymm6, ymm1
	LONG $0xd672c5c5; BYTE $0x1b   // vpsrld    ymm7, ymm6, 27
	LONG $0xf672cdc5; BYTE $0x05   // vpslld    ymm6, ymm6, 5
	LONG $0xf7ebcdc5               // vpor    ymm6, ymm6, ymm7
	LONG $0x75efcdc5; BYTE $0x40   // vpxor    ymm6, ymm6, yword 64[rbp]
	LONG $0xc9fecdc5               // vpaddd    ymm1, ymm6, ymm1
	LONG $0xd172c5c5; BYTE $0x0f   // vpsrld    ymm7, ymm1, 15
	LONG $0xf172f5c5; BYTE $0x11   // vpslld    ymm1, ymm1, 17
	LONG $0xcfebf5c5               // vpor    ymm1, ymm1, ymm7
	LONG $0xf6fef5c5               // vpaddd    ymm6, ymm1, ymm6
	LONG $0xfe70fdc5; BYTE $0xd2   // vpshufd    ymm7, ymm6, 210
	LONG $0x0075e2c4; BYTE $0xc8   // vpshufb    ymm1, ymm1, ymm0
	LONG $0x4645e3c4; WORD $0x31f1 // vperm2i128    ymm6, ymm7, ymm1, 49
	LONG $0x3845e3c4; WORD $0x01f9 // vinserti128    ymm7, ymm7, xmm1, 1
	LONG $0x000060b8; BYTE $0x00   // mov    eax, 96
	LONG $0x4d6ffdc5; BYTE $0x60   // vmovdqa    ymm1, yword 96[rbp]
	MOVQ ·step(SB), R8             // lea    r8, [rip + _g_StepConstants]
LBB1_42:
	LONG $0x3675e2c4; BYTE $0xdb               // vpermd    ymm3, ymm1, ymm3
	LONG $0xddfee5c5                           // vpaddd    ymm3, ymm3, ymm5
	LONG $0x3675e2c4; BYTE $0xd2               // vpermd    ymm2, ymm1, ymm2
	LONG $0xd4feedc5                           // vpaddd    ymm2, ymm2, ymm4
	LONG $0xffefedc5                           // vpxor    ymm7, ymm2, ymm7
	LONG $0xf6efe5c5                           // vpxor    ymm6, ymm3, ymm6
	LONG $0xf6fec5c5                           // vpaddd    ymm6, ymm7, ymm6
	LONG $0xd672bdc5; BYTE $0x03               // vpsrld    ymm8, ymm6, 3
	LONG $0xf672cdc5; BYTE $0x1d               // vpslld    ymm6, ymm6, 29
	LONG $0xf6ebbdc5                           // vpor    ymm6, ymm8, ymm6
	LONG $0xef4da1c4; WORD $0x0074; BYTE $0xe0 // vpxor    ymm6, ymm6, yword [rax + r8 - 32]
	LONG $0xfffecdc5                           // vpaddd    ymm7, ymm6, ymm7
	LONG $0xd772bdc5; BYTE $0x1f               // vpsrld    ymm8, ymm7, 31
	LONG $0xfffec5c5                           // vpaddd    ymm7, ymm7, ymm7
	LONG $0xffebbdc5                           // vpor    ymm7, ymm8, ymm7
	LONG $0xf6fec5c5                           // vpaddd    ymm6, ymm7, ymm6
	LONG $0xf670fdc5; BYTE $0xd2               // vpshufd    ymm6, ymm6, 210
	LONG $0x0045e2c4; BYTE $0xf8               // vpshufb    ymm7, ymm7, ymm0
	LONG $0x3675e2c4; BYTE $0xed               // vpermd    ymm5, ymm1, ymm5
	LONG $0xedfee5c5                           // vpaddd    ymm5, ymm3, ymm5
	LONG $0x3675e2c4; BYTE $0xe4               // vpermd    ymm4, ymm1, ymm4
	LONG $0xe4feedc5                           // vpaddd    ymm4, ymm2, ymm4
	LONG $0x384d63c4; WORD $0x01c7             // vinserti128    ymm8, ymm6, xmm7, 1
	LONG $0xc4ef3dc5                           // vpxor    ymm8, ymm8, ymm4
	LONG $0x464de3c4; WORD $0x31f7             // vperm2i128    ymm6, ymm6, ymm7, 49
	LONG $0xf6efd5c5                           // vpxor    ymm6, ymm5, ymm6
	LONG $0xf6febdc5                           // vpaddd    ymm6, ymm8, ymm6
	LONG $0xd672c5c5; BYTE $0x1b               // vpsrld    ymm7, ymm6, 27
	LONG $0xf672cdc5; BYTE $0x05               // vpslld    ymm6, ymm6, 5
	LONG $0xf7ebcdc5                           // vpor    ymm6, ymm6, ymm7
	LONG $0xef4da1c4; WORD $0x0034             // vpxor    ymm6, ymm6, yword [rax + r8]
	LONG $0xfefebdc5                           // vpaddd    ymm7, ymm8, ymm6
	LONG $0xd772bdc5; BYTE $0x0f               // vpsrld    ymm8, ymm7, 15
	LONG $0xf772c5c5; BYTE $0x11               // vpslld    ymm7, ymm7, 17
	LONG $0xffebbdc5                           // vpor    ymm7, ymm8, ymm7
	LONG $0xf6fec5c5                           // vpaddd    ymm6, ymm7, ymm6
	LONG $0xc6707dc5; BYTE $0xd2               // vpshufd    ymm8, ymm6, 210
	LONG $0x0045e2c4; BYTE $0xf8               // vpshufb    ymm7, ymm7, ymm0
	LONG $0x463de3c4; WORD $0x31f7             // vperm2i128    ymm6, ymm8, ymm7, 49
	LONG $0x383de3c4; WORD $0x01ff             // vinserti128    ymm7, ymm8, xmm7, 1
	LONG $0x40c08348                           // add    rax, 64
	LONG $0x03603d48; WORD $0x0000             // cmp    rax, 864
	JNE  LBB1_42
	LONG $0x3675e2c4; BYTE $0xc3               // vpermd    ymm0, ymm1, ymm3
	LONG $0xc0fed5c5                           // vpaddd    ymm0, ymm5, ymm0
	LONG $0x3675e2c4; BYTE $0xca               // vpermd    ymm1, ymm1, ymm2
	LONG $0xc9feddc5                           // vpaddd    ymm1, ymm4, ymm1
	LONG $0xc0efcdc5                           // vpxor    ymm0, ymm6, ymm0
	LONG $0xc9efc5c5                           // vpxor    ymm1, ymm7, ymm1
	WORD $0x0148; BYTE $0xce                   // add    rsi, rcx
	WORD $0x2948; BYTE $0xca                   // sub    rdx, rcx
	LONG $0x001047c7; WORD $0x0000; BYTE $0x00 // mov    dword [rdi + 16], 0

LBB1_44:
	LONG $0x80fa8148; WORD $0x0000; BYTE $0x00 // cmp    rdx, 128
	JB   LBB1_49
	LONG $0x556ffdc5; BYTE $0x00               // vmovdqa    ymm2, yword 0[rbp]
	LONG $0x5d6ffdc5; BYTE $0x20               // vmovdqa    ymm3, yword 32[rbp]
	LONG $0x656ffdc5; BYTE $0x40               // vmovdqa    ymm4, yword 64[rbp]
	LONG $0x6d6ffdc5; BYTE $0x60               // vmovdqa    ymm5, yword 96[rbp]
	MOVQ ·step(SB), AX                         //  lea    rax, [rip + _g_StepConstants]
LBB1_46:
	LONG $0x3e6ffec5               // vmovdqu    ymm7, yword [rsi]
	LONG $0x766ffec5; BYTE $0x20   // vmovdqu    ymm6, yword [rsi + 32]
	LONG $0x4e6f7ec5; BYTE $0x40   // vmovdqu    ymm9, yword [rsi + 64]
	LONG $0x466f7ec5; BYTE $0x60   // vmovdqu    ymm8, yword [rsi + 96]
	LONG $0xc0efc5c5               // vpxor    ymm0, ymm7, ymm0
	LONG $0xc9efcdc5               // vpxor    ymm1, ymm6, ymm1
	LONG $0xc0fef5c5               // vpaddd    ymm0, ymm1, ymm0
	LONG $0xd072adc5; BYTE $0x03   // vpsrld    ymm10, ymm0, 3
	LONG $0xf072fdc5; BYTE $0x1d   // vpslld    ymm0, ymm0, 29
	LONG $0xc0ebadc5               // vpor    ymm0, ymm10, ymm0
	LONG $0xc2effdc5               // vpxor    ymm0, ymm0, ymm2
	LONG $0xc9fefdc5               // vpaddd    ymm1, ymm0, ymm1
	LONG $0xd172adc5; BYTE $0x1f   // vpsrld    ymm10, ymm1, 31
	LONG $0xc9fef5c5               // vpaddd    ymm1, ymm1, ymm1
	LONG $0xc9ebadc5               // vpor    ymm1, ymm10, ymm1
	LONG $0xc0fef5c5               // vpaddd    ymm0, ymm1, ymm0
	LONG $0xc070fdc5; BYTE $0xd2   // vpshufd    ymm0, ymm0, 210
	LONG $0x0075e2c4; BYTE $0xcb   // vpshufb    ymm1, ymm1, ymm3
	LONG $0x467d63c4; WORD $0x31d1 // vperm2i128    ymm10, ymm0, ymm1, 49
	LONG $0x387de3c4; WORD $0x01c1 // vinserti128    ymm0, ymm0, xmm1, 1
	LONG $0xef2dc1c4; BYTE $0xc9   // vpxor    ymm1, ymm10, ymm9
	LONG $0xc0efbdc5               // vpxor    ymm0, ymm8, ymm0
	LONG $0xc8fef5c5               // vpaddd    ymm1, ymm1, ymm0
	LONG $0xd172adc5; BYTE $0x1b   // vpsrld    ymm10, ymm1, 27
	LONG $0xf172f5c5; BYTE $0x05   // vpslld    ymm1, ymm1, 5
	LONG $0xc9ebadc5               // vpor    ymm1, ymm10, ymm1
	LONG $0xcceff5c5               // vpxor    ymm1, ymm1, ymm4
	LONG $0xc0fef5c5               // vpaddd    ymm0, ymm1, ymm0
	LONG $0xd072adc5; BYTE $0x0f   // vpsrld    ymm10, ymm0, 15
	LONG $0xf072fdc5; BYTE $0x11   // vpslld    ymm0, ymm0, 17
	LONG $0xc0ebadc5               // vpor    ymm0, ymm10, ymm0
	LONG $0xc9fefdc5               // vpaddd    ymm1, ymm0, ymm1
	LONG $0xc970fdc5; BYTE $0xd2   // vpshufd    ymm1, ymm1, 210
	LONG $0x007d62c4; BYTE $0xd3   // vpshufb    ymm10, ymm0, ymm3
	LONG $0x4675c3c4; WORD $0x31c2 // vperm2i128    ymm0, ymm1, ymm10, 49
	LONG $0x3875c3c4; WORD $0x01ca // vinserti128    ymm1, ymm1, xmm10, 1
	LONG $0x000060b9; BYTE $0x00   // mov    ecx, 96

LBB1_47:
	LONG $0x3655e2c4; BYTE $0xff               // vpermd    ymm7, ymm5, ymm7
	LONG $0xfffeb5c5                           // vpaddd    ymm7, ymm9, ymm7
	LONG $0x3655e2c4; BYTE $0xf6               // vpermd    ymm6, ymm5, ymm6
	LONG $0xf6febdc5                           // vpaddd    ymm6, ymm8, ymm6
	LONG $0xc7effdc5                           // vpxor    ymm0, ymm0, ymm7
	LONG $0xceeff5c5                           // vpxor    ymm1, ymm1, ymm6
	LONG $0xc0fef5c5                           // vpaddd    ymm0, ymm1, ymm0
	LONG $0xd072adc5; BYTE $0x03               // vpsrld    ymm10, ymm0, 3
	LONG $0xf072fdc5; BYTE $0x1d               // vpslld    ymm0, ymm0, 29
	LONG $0xc0ebadc5                           // vpor    ymm0, ymm10, ymm0
	LONG $0x44effdc5; WORD $0xe001             // vpxor    ymm0, ymm0, yword [rcx + rax - 32]
	LONG $0xc9fefdc5                           // vpaddd    ymm1, ymm0, ymm1
	LONG $0xd172adc5; BYTE $0x1f               // vpsrld    ymm10, ymm1, 31
	LONG $0xc9fef5c5                           // vpaddd    ymm1, ymm1, ymm1
	LONG $0xc9ebadc5                           // vpor    ymm1, ymm10, ymm1
	LONG $0xc0fef5c5                           // vpaddd    ymm0, ymm1, ymm0
	LONG $0xc070fdc5; BYTE $0xd2               // vpshufd    ymm0, ymm0, 210
	LONG $0x0075e2c4; BYTE $0xcb               // vpshufb    ymm1, ymm1, ymm3
	LONG $0x467d63c4; WORD $0x31d1             // vperm2i128    ymm10, ymm0, ymm1, 49
	LONG $0x387de3c4; WORD $0x01c1             // vinserti128    ymm0, ymm0, xmm1, 1
	LONG $0x3655c2c4; BYTE $0xc9               // vpermd    ymm1, ymm5, ymm9
	LONG $0xc9fe45c5                           // vpaddd    ymm9, ymm7, ymm1
	LONG $0x3655c2c4; BYTE $0xc8               // vpermd    ymm1, ymm5, ymm8
	LONG $0xc1fe4dc5                           // vpaddd    ymm8, ymm6, ymm1
	LONG $0xef2dc1c4; BYTE $0xc9               // vpxor    ymm1, ymm10, ymm9
	LONG $0xc0efbdc5                           // vpxor    ymm0, ymm8, ymm0
	LONG $0xc8fef5c5                           // vpaddd    ymm1, ymm1, ymm0
	LONG $0xd172adc5; BYTE $0x1b               // vpsrld    ymm10, ymm1, 27
	LONG $0xf172f5c5; BYTE $0x05               // vpslld    ymm1, ymm1, 5
	LONG $0xc9ebadc5                           // vpor    ymm1, ymm10, ymm1
	LONG $0x0ceff5c5; BYTE $0x01               // vpxor    ymm1, ymm1, yword [rcx + rax]
	LONG $0xc0fef5c5                           // vpaddd    ymm0, ymm1, ymm0
	LONG $0xd072adc5; BYTE $0x0f               // vpsrld    ymm10, ymm0, 15
	LONG $0xf072fdc5; BYTE $0x11               // vpslld    ymm0, ymm0, 17
	LONG $0xc0ebadc5                           // vpor    ymm0, ymm10, ymm0
	LONG $0xc9fefdc5                           // vpaddd    ymm1, ymm0, ymm1
	LONG $0xc970fdc5; BYTE $0xd2               // vpshufd    ymm1, ymm1, 210
	LONG $0x007d62c4; BYTE $0xd3               // vpshufb    ymm10, ymm0, ymm3
	LONG $0x4675c3c4; WORD $0x31c2             // vperm2i128    ymm0, ymm1, ymm10, 49
	LONG $0x3875c3c4; WORD $0x01ca             // vinserti128    ymm1, ymm1, xmm10, 1
	LONG $0x40c18348                           // add    rcx, 64
	LONG $0x60f98148; WORD $0x0003; BYTE $0x00 // cmp    rcx, 864
	JNE  LBB1_47
	LONG $0x3655e2c4; BYTE $0xff               // vpermd    ymm7, ymm5, ymm7
	LONG $0xfffeb5c5                           // vpaddd    ymm7, ymm9, ymm7
	LONG $0x3655e2c4; BYTE $0xf6               // vpermd    ymm6, ymm5, ymm6
	LONG $0xf6febdc5                           // vpaddd    ymm6, ymm8, ymm6
	LONG $0xc7effdc5                           // vpxor    ymm0, ymm0, ymm7
	LONG $0xceeff5c5                           // vpxor    ymm1, ymm1, ymm6
	LONG $0x80ee8348                           // sub    rsi, -128
	LONG $0x80c28348                           // add    rdx, -128
	LONG $0x7ffa8348                           // cmp    rdx, 127
	JA   LBB1_46

LBB1_49:
	LONG $0x477ffec5; BYTE $0x20   // vmovdqu    yword [rdi + 32], ymm0
	LONG $0x4f7ffec5; BYTE $0x40   // vmovdqu    yword [rdi + 64], ymm1
	WORD $0x8548; BYTE $0xd2       // test    rdx, rdx
	JE   LBB1_63
	LONG $0x10fa8348               // cmp    rdx, 16
	JB   LBB1_51
	WORD $0x8948; BYTE $0xf8       // mov    rax, rdi
	WORD $0x2948; BYTE $0xf0       // sub    rax, rsi
	LONG $0x60c08348               // add    rax, 96
	LONG $0x00803d48; WORD $0x0000 // cmp    rax, 128
	JAE  LBB1_54

LBB1_51:
	WORD $0xc031 // xor    eax, eax

LBB1_57:
	WORD $0x8941; BYTE $0xd0 // mov    r8d, edx
	WORD $0x2941; BYTE $0xc0 // sub    r8d, eax
	WORD $0x8948; BYTE $0xc1 // mov    rcx, rax
	WORD $0xf748; BYTE $0xd1 // not    rcx
	WORD $0x0148; BYTE $0xd1 // add    rcx, rdx
	LONG $0x03e08341         // and    r8d, 3
	JE   LBB1_59

LBB1_58:
	LONG $0x0cb60f44; BYTE $0x06 // movzx    r9d, byte [rsi + rax]
	LONG $0x074c8844; BYTE $0x60 // mov    byte [rdi + rax + 96], r9b
	WORD $0xff48; BYTE $0xc0     // inc    rax
	WORD $0xff49; BYTE $0xc8     // dec    r8
	JNE  LBB1_58

LBB1_59:
	LONG $0x03f98348 // cmp    rcx, 3
	JB   LBB1_61

LBB1_60:
	LONG $0x060cb60f             // movzx    ecx, byte [rsi + rax]
	LONG $0x60074c88             // mov    byte [rdi + rax + 96], cl
	LONG $0x064cb60f; BYTE $0x01 // movzx    ecx, byte [rsi + rax + 1]
	LONG $0x61074c88             // mov    byte [rdi + rax + 97], cl
	LONG $0x064cb60f; BYTE $0x02 // movzx    ecx, byte [rsi + rax + 2]
	LONG $0x62074c88             // mov    byte [rdi + rax + 98], cl
	LONG $0x064cb60f; BYTE $0x03 // movzx    ecx, byte [rsi + rax + 3]
	LONG $0x63074c88             // mov    byte [rdi + rax + 99], cl
	LONG $0x04c08348             // add    rax, 4
	WORD $0x3948; BYTE $0xc2     // cmp    rdx, rax
	JNE  LBB1_60
	JMP  LBB1_61

LBB1_12:
	LONG $0x0080f981; WORD $0x0000 // cmp    ecx, 128
	JAE  LBB1_16
	WORD $0x3145; BYTE $0xc9       // xor    r9d, r9d
	JMP  LBB1_20

LBB1_54:
	WORD $0xd089             // mov    eax, edx
	WORD $0xe083; BYTE $0x70 // and    eax, 112
	WORD $0xc931             // xor    ecx, ecx

LBB1_55:
	LONG $0x046ffac5; BYTE $0x0e   // vmovdqu    xmm0, oword [rsi + rcx]
	LONG $0x447ffac5; WORD $0x600f // vmovdqu    oword [rdi + rcx + 96], xmm0
	LONG $0x10c18348               // add    rcx, 16
	WORD $0x3948; BYTE $0xc8       // cmp    rax, rcx
	JNE  LBB1_55
	WORD $0x3948; BYTE $0xc2       // cmp    rdx, rax
	JNE  LBB1_57

LBB1_61:
	WORD $0xd089 // mov    eax, edx

LBB1_62:
	WORD $0x4789; BYTE $0x10 // mov    dword [rdi + 16], eax

LBB1_63:
	VZEROUPPER
	RET

LBB1_14:
	LONG $0x0000813d; BYTE $0x00 // cmp    eax, 129
	JAE  LBB1_29
	WORD $0x3145; BYTE $0xc0     // xor    r8d, r8d
	JMP  LBB1_33

LBB1_16:
	WORD $0x8941; BYTE $0xd2                   // mov    r10d, edx
	LONG $0x7fe28341                           // and    r10d, 127
	WORD $0x8949; BYTE $0xc9                   // mov    r9, rcx
	WORD $0x294d; BYTE $0xd1                   // sub    r9, r10
	LONG $0x381c8d4c                           // lea    r11, [rax + rdi]
	LONG $0xc0c38149; WORD $0x0000; BYTE $0x00 // add    r11, 192
	WORD $0xdb31                               // xor    ebx, ebx

LBB1_17:
	LONG $0x046ffec5; BYTE $0x1e               // vmovdqu    ymm0, yword [rsi + rbx]
	LONG $0x4c6ffec5; WORD $0x201e             // vmovdqu    ymm1, yword [rsi + rbx + 32]
	LONG $0x546ffec5; WORD $0x401e             // vmovdqu    ymm2, yword [rsi + rbx + 64]
	LONG $0x5c6ffec5; WORD $0x601e             // vmovdqu    ymm3, yword [rsi + rbx + 96]
	LONG $0x7f7ec1c4; WORD $0x1b44; BYTE $0xa0 // vmovdqu    yword [r11 + rbx - 96], ymm0
	LONG $0x7f7ec1c4; WORD $0x1b4c; BYTE $0xc0 // vmovdqu    yword [r11 + rbx - 64], ymm1
	LONG $0x7f7ec1c4; WORD $0x1b54; BYTE $0xe0 // vmovdqu    yword [r11 + rbx - 32], ymm2
	LONG $0x7f7ec1c4; WORD $0x1b1c             // vmovdqu    yword [r11 + rbx], ymm3
	LONG $0x80eb8348                           // sub    rbx, -128
	WORD $0x3949; BYTE $0xd9                   // cmp    r9, rbx
	JNE  LBB1_17
	WORD $0x854d; BYTE $0xd2                   // test    r10, r10
	JE   LBB1_28
	LONG $0x10fa8341                           // cmp    r10d, 16
	JB   LBB1_23

LBB1_20:
	WORD $0x894d; BYTE $0xca // mov    r10, r9
	WORD $0x8941; BYTE $0xd3 // mov    r11d, edx
	LONG $0x0fe38341         // and    r11d, 15
	WORD $0x8949; BYTE $0xc9 // mov    r9, rcx
	WORD $0x294d; BYTE $0xd9 // sub    r9, r11

LBB1_21:
	LONG $0x6f7aa1c4; WORD $0x1604 // vmovdqu    xmm0, oword [rsi + r10]
	LONG $0x7f7a81c4; WORD $0x1004 // vmovdqu    oword [r8 + r10], xmm0
	LONG $0x10c28349               // add    r10, 16
	WORD $0x394d; BYTE $0xd1       // cmp    r9, r10
	JNE  LBB1_21
	WORD $0x854d; BYTE $0xdb       // test    r11, r11
	JNE  LBB1_23
	JMP  LBB1_28

LBB1_29:
	WORD $0x8941; BYTE $0xc8                   // mov    r8d, ecx
	LONG $0x80e08341                           // and    r8d, -128
	LONG $0x38148d4c                           // lea    r10, [rax + rdi]
	LONG $0xc0c28149; WORD $0x0000; BYTE $0x00 // add    r10, 192
	WORD $0x3145; BYTE $0xdb                   // xor    r11d, r11d

LBB1_30:
	LONG $0x107ca1c4; WORD $0x1e14             // vmovups    ymm2, yword [rsi + r11]
	LONG $0x107ca1c4; WORD $0x1e5c; BYTE $0x20 // vmovups    ymm3, yword [rsi + r11 + 32]
	LONG $0x107ca1c4; WORD $0x1e64; BYTE $0x40 // vmovups    ymm4, yword [rsi + r11 + 64]
	LONG $0x107ca1c4; WORD $0x1e6c; BYTE $0x60 // vmovups    ymm5, yword [rsi + r11 + 96]
	LONG $0x117c81c4; WORD $0x1a54; BYTE $0xa0 // vmovups    yword [r10 + r11 - 96], ymm2
	LONG $0x117c81c4; WORD $0x1a5c; BYTE $0xc0 // vmovups    yword [r10 + r11 - 64], ymm3
	LONG $0x117c81c4; WORD $0x1a64; BYTE $0xe0 // vmovups    yword [r10 + r11 - 32], ymm4
	LONG $0x117c81c4; WORD $0x1a2c             // vmovups    yword [r10 + r11], ymm5
	LONG $0x80eb8349                           // sub    r11, -128
	WORD $0x394d; BYTE $0xd8                   // cmp    r8, r11
	JNE  LBB1_30
	WORD $0x3949; BYTE $0xc8                   // cmp    r8, rcx
	JE   LBB1_41
	WORD $0xc1f6; BYTE $0x70                   // test    cl, 112
	JE   LBB1_36

LBB1_33:
	WORD $0x894d; BYTE $0xc2 // mov    r10, r8
	WORD $0x8941; BYTE $0xc8 // mov    r8d, ecx
	LONG $0xf0e08341         // and    r8d, -16

LBB1_34:
	LONG $0x1078a1c4; WORD $0x1614 // vmovups    xmm2, oword [rsi + r10]
	LONG $0x117881c4; WORD $0x1114 // vmovups    oword [r9 + r10], xmm2
	LONG $0x10c28349               // add    r10, 16
	WORD $0x394d; BYTE $0xd0       // cmp    r8, r10
	JNE  LBB1_34
	WORD $0x3949; BYTE $0xc8       // cmp    r8, rcx
	JE   LBB1_41
	JMP  LBB1_36

TEXT ·__lsh256_avx2_final(SB), NOSPLIT, $0-16
	MOVQ ctx+0(FP), DI
	MOVQ hashval+8(FP), SI

	LEAQ LCDATA2<>(SB), BP

	LONG $0x10478b44               // mov    r8d, dword [rdi + 16]
	LONG $0x0744c642; WORD $0x8060 // mov    byte [rdi + r8 + 96], -128
	LONG $0x00007fba; BYTE $0x00   // mov    edx, 127
	WORD $0x2944; BYTE $0xc2       // sub    edx, r8d
	WORD $0xd285                   // test    edx, edx
	JLE  LBB2_13
	LONG $0x07048d4a               // lea    rax, [rdi + r8]
	LONG $0x61c08348               // add    rax, 97
	WORD $0xd189                   // mov    ecx, edx
	WORD $0xfa83; BYTE $0x10       // cmp    edx, 16
	JAE  LBB2_3
	WORD $0xd231                   // xor    edx, edx
	JMP  LBB2_12

LBB2_3:
	LONG $0x0080fa81; WORD $0x0000 // cmp    edx, 128
	JAE  LBB2_5
	WORD $0xd231                   // xor    edx, edx
	JMP  LBB2_9

LBB2_5:
	WORD $0xca89                               // mov    edx, ecx
	WORD $0xe283; BYTE $0x80                   // and    edx, -128
	WORD $0x0149; BYTE $0xf8                   // add    r8, rdi
	LONG $0xc1c08149; WORD $0x0000; BYTE $0x00 // add    r8, 193
	WORD $0x3145; BYTE $0xc9                   // xor    r9d, r9d
	LONG $0xc057f8c5                           // vxorps    xmm0, xmm0, xmm0

LBB2_6:
	LONG $0x117c81c4; WORD $0x0844; BYTE $0xa0 // vmovups    yword [r8 + r9 - 96], ymm0
	LONG $0x117c81c4; WORD $0x0844; BYTE $0xc0 // vmovups    yword [r8 + r9 - 64], ymm0
	LONG $0x117c81c4; WORD $0x0844; BYTE $0xe0 // vmovups    yword [r8 + r9 - 32], ymm0
	LONG $0x117c81c4; WORD $0x0804             // vmovups    yword [r8 + r9], ymm0
	LONG $0x80e98349                           // sub    r9, -128
	WORD $0x394c; BYTE $0xca                   // cmp    rdx, r9
	JNE  LBB2_6
	WORD $0x3948; BYTE $0xca                   // cmp    rdx, rcx
	JE   LBB2_13
	WORD $0xc1f6; BYTE $0x70                   // test    cl, 112
	JE   LBB2_12

LBB2_9:
	WORD $0x8949; BYTE $0xd0 // mov    r8, rdx
	WORD $0xca89             // mov    edx, ecx
	WORD $0xe283; BYTE $0xf0 // and    edx, -16
	LONG $0xc057f8c5         // vxorps    xmm0, xmm0, xmm0

LBB2_10:
	LONG $0x1178a1c4; WORD $0x0004 // vmovups    oword [rax + r8], xmm0
	LONG $0x10c08349               // add    r8, 16
	WORD $0x394c; BYTE $0xc2       // cmp    rdx, r8
	JNE  LBB2_10
	WORD $0x3948; BYTE $0xca       // cmp    rdx, rcx
	JE   LBB2_13

LBB2_12:
	LONG $0x001004c6         // mov    byte [rax + rdx], 0
	WORD $0xff48; BYTE $0xc2 // inc    rdx
	WORD $0x3948; BYTE $0xd1 // cmp    rcx, rdx
	JNE  LBB2_12

LBB2_13:
	LONG $0x4f6ffec5; BYTE $0x60   // vmovdqu    ymm1, yword [rdi + 96]
	QUAD $0x00000080876ffec5       // vmovdqu    ymm0, yword [rdi + 128]
	QUAD $0x000000a09f6ffec5       // vmovdqu    ymm3, yword [rdi + 160]
	QUAD $0x000000c0976ffec5       // vmovdqu    ymm2, yword [rdi + 192]
	LONG $0x67eff5c5; BYTE $0x20   // vpxor    ymm4, ymm1, yword [rdi + 32]
	LONG $0x6feffdc5; BYTE $0x40   // vpxor    ymm5, ymm0, yword [rdi + 64]
	LONG $0xe4fed5c5               // vpaddd    ymm4, ymm5, ymm4
	LONG $0xd472cdc5; BYTE $0x03   // vpsrld    ymm6, ymm4, 3
	LONG $0xf472ddc5; BYTE $0x1d   // vpslld    ymm4, ymm4, 29
	LONG $0xe6ebddc5               // vpor    ymm4, ymm4, ymm6
	LONG $0x65efddc5; BYTE $0x00   // vpxor    ymm4, ymm4, yword 0[rbp]
	LONG $0xedfeddc5               // vpaddd    ymm5, ymm4, ymm5
	LONG $0xd572cdc5; BYTE $0x1f   // vpsrld    ymm6, ymm5, 31
	LONG $0xedfed5c5               // vpaddd    ymm5, ymm5, ymm5
	LONG $0xeeebd5c5               // vpor    ymm5, ymm5, ymm6
	LONG $0xe4fed5c5               // vpaddd    ymm4, ymm5, ymm4
	LONG $0xf470fdc5; BYTE $0xd2   // vpshufd    ymm6, ymm4, 210
	LONG $0x656ffdc5; BYTE $0x20   // vmovdqa    ymm4, yword 32[rbp]
	LONG $0x0055e2c4; BYTE $0xec   // vpshufb    ymm5, ymm5, ymm4
	LONG $0x464de3c4; WORD $0x31fd // vperm2i128    ymm7, ymm6, ymm5, 49
	LONG $0x384de3c4; WORD $0x01ed // vinserti128    ymm5, ymm6, xmm5, 1
	LONG $0xf3efc5c5               // vpxor    ymm6, ymm7, ymm3
	LONG $0xeaefd5c5               // vpxor    ymm5, ymm5, ymm2
	LONG $0xf5fecdc5               // vpaddd    ymm6, ymm6, ymm5
	LONG $0xd672c5c5; BYTE $0x1b   // vpsrld    ymm7, ymm6, 27
	LONG $0xf672cdc5; BYTE $0x05   // vpslld    ymm6, ymm6, 5
	LONG $0xf7ebcdc5               // vpor    ymm6, ymm6, ymm7
	LONG $0x75efcdc5; BYTE $0x40   // vpxor    ymm6, ymm6, yword 64[rbp]
	LONG $0xedfecdc5               // vpaddd    ymm5, ymm6, ymm5
	LONG $0xd572c5c5; BYTE $0x0f   // vpsrld    ymm7, ymm5, 15
	LONG $0xf572d5c5; BYTE $0x11   // vpslld    ymm5, ymm5, 17
	LONG $0xefebd5c5               // vpor    ymm5, ymm5, ymm7
	LONG $0xf6fed5c5               // vpaddd    ymm6, ymm5, ymm6
	LONG $0xfe70fdc5; BYTE $0xd2   // vpshufd    ymm7, ymm6, 210
	LONG $0x0055e2c4; BYTE $0xec   // vpshufb    ymm5, ymm5, ymm4
	LONG $0x4645e3c4; WORD $0x31f5 // vperm2i128    ymm6, ymm7, ymm5, 49
	LONG $0x3845e3c4; WORD $0x01fd // vinserti128    ymm7, ymm7, xmm5, 1
	LONG $0x000060b8; BYTE $0x00   // mov    eax, 96
	LONG $0x6d6ffdc5; BYTE $0x60   // vmovdqa    ymm5, yword 96[rbp]
	MOVQ ·step(SB), CX             // lea    rcx, [rip + _g_StepConstants]

LBB2_14:
	LONG $0x3655e2c4; BYTE $0xc9   // vpermd    ymm1, ymm5, ymm1
	LONG $0xcbfef5c5               // vpaddd    ymm1, ymm1, ymm3
	LONG $0x3655e2c4; BYTE $0xc0   // vpermd    ymm0, ymm5, ymm0
	LONG $0xc2fefdc5               // vpaddd    ymm0, ymm0, ymm2
	LONG $0xffeffdc5               // vpxor    ymm7, ymm0, ymm7
	LONG $0xf6eff5c5               // vpxor    ymm6, ymm1, ymm6
	LONG $0xf6fec5c5               // vpaddd    ymm6, ymm7, ymm6
	LONG $0xd672bdc5; BYTE $0x03   // vpsrld    ymm8, ymm6, 3
	LONG $0xf672cdc5; BYTE $0x1d   // vpslld    ymm6, ymm6, 29
	LONG $0xf6ebbdc5               // vpor    ymm6, ymm8, ymm6
	LONG $0x74efcdc5; WORD $0xe008 // vpxor    ymm6, ymm6, yword [rax + rcx - 32]
	LONG $0xfffecdc5               // vpaddd    ymm7, ymm6, ymm7
	LONG $0xd772bdc5; BYTE $0x1f   // vpsrld    ymm8, ymm7, 31
	LONG $0xfffec5c5               // vpaddd    ymm7, ymm7, ymm7
	LONG $0xffebbdc5               // vpor    ymm7, ymm8, ymm7
	LONG $0xf6fec5c5               // vpaddd    ymm6, ymm7, ymm6
	LONG $0xf670fdc5; BYTE $0xd2   // vpshufd    ymm6, ymm6, 210
	LONG $0x0045e2c4; BYTE $0xfc   // vpshufb    ymm7, ymm7, ymm4
	LONG $0x3655e2c4; BYTE $0xdb   // vpermd    ymm3, ymm5, ymm3
	LONG $0xdbfef5c5               // vpaddd    ymm3, ymm1, ymm3
	LONG $0x3655e2c4; BYTE $0xd2   // vpermd    ymm2, ymm5, ymm2
	LONG $0xd2fefdc5               // vpaddd    ymm2, ymm0, ymm2
	LONG $0x384d63c4; WORD $0x01c7 // vinserti128    ymm8, ymm6, xmm7, 1
	LONG $0xc2ef3dc5               // vpxor    ymm8, ymm8, ymm2
	LONG $0x464de3c4; WORD $0x31f7 // vperm2i128    ymm6, ymm6, ymm7, 49
	LONG $0xf6efe5c5               // vpxor    ymm6, ymm3, ymm6
	LONG $0xf6febdc5               // vpaddd    ymm6, ymm8, ymm6
	LONG $0xd672c5c5; BYTE $0x1b   // vpsrld    ymm7, ymm6, 27
	LONG $0xf672cdc5; BYTE $0x05   // vpslld    ymm6, ymm6, 5
	LONG $0xf7ebcdc5               // vpor    ymm6, ymm6, ymm7
	LONG $0x34efcdc5; BYTE $0x08   // vpxor    ymm6, ymm6, yword [rax + rcx]
	LONG $0xfefebdc5               // vpaddd    ymm7, ymm8, ymm6
	LONG $0xd772bdc5; BYTE $0x0f   // vpsrld    ymm8, ymm7, 15
	LONG $0xf772c5c5; BYTE $0x11   // vpslld    ymm7, ymm7, 17
	LONG $0xffebbdc5               // vpor    ymm7, ymm8, ymm7
	LONG $0xf6fec5c5               // vpaddd    ymm6, ymm7, ymm6
	LONG $0xc6707dc5; BYTE $0xd2   // vpshufd    ymm8, ymm6, 210
	LONG $0x0045e2c4; BYTE $0xfc   // vpshufb    ymm7, ymm7, ymm4
	LONG $0x463de3c4; WORD $0x31f7 // vperm2i128    ymm6, ymm8, ymm7, 49
	LONG $0x383de3c4; WORD $0x01ff // vinserti128    ymm7, ymm8, xmm7, 1
	LONG $0x40c08348               // add    rax, 64
	LONG $0x03603d48; WORD $0x0000 // cmp    rax, 864
	JNE  LBB2_14
	LONG $0x3655e2c4; BYTE $0xc9   // vpermd    ymm1, ymm5, ymm1
	LONG $0xc9fee5c5               // vpaddd    ymm1, ymm3, ymm1
	LONG $0x3655e2c4; BYTE $0xc0   // vpermd    ymm0, ymm5, ymm0
	LONG $0xc0feedc5               // vpaddd    ymm0, ymm2, ymm0
	LONG $0xc1effdc5               // vpxor    ymm0, ymm0, ymm1
	LONG $0xc0efc5c5               // vpxor    ymm0, ymm7, ymm0
	LONG $0xc6effdc5               // vpxor    ymm0, ymm0, ymm6
	LONG $0x067ffec5               // vmovdqu    yword [rsi], ymm0
	VZEROUPPER
	RET
