package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

// argument align
// arg0 ---- rdi
// arg1 ---- esi
// arg2 ---- rdx

var (
	gLC0  Mem = global64("LC0", -153997021875728609, 1841527622716375976)
	gLC1  Mem = global64("LC1", -5695466077790390131, 8852904047995281860)
	gLC2  Mem = global64("LC2", 3419602469714219896, -984810856687653483)
	gLC3  Mem = global64("LC3", 1406927143955412346, -2298248402025429884)
	gLC4  Mem = global64("LC4", 7122715107427485907, 5503410031139639979)
	gLC5  Mem = global64("LC5", 2218598310921636520, 3551039192594811326)
	gLC6  Mem = global64("LC6", -694899722053307143, 6249004012466285246)
	gLC7  Mem = global64("LC7", -1293352058239094138, -8378717999420185195)
	gLC8  Mem = global64("LC8", 7789838270879018896, -3497183359489726141)
	gLC9  Mem = global64("LC9", 3020068111055025266, 3380838012506842152)
	gLC10 Mem = global64("LC10", -4294967296, -1)
	gLC11 Mem = global64("LC11", 4294967295, 0)
	gLC12 Mem = global64("LC12", 0, -1)
	gLC13 Mem = global64("LC13", -1, 0)
	gLC14 Mem = global64("LC14", 0, -4294967296)
	gLC15 Mem = global64("LC15", -1, 4294967295)
	gLC16 Mem = global64("LC16", -8706668581642026975, 5114336635225271474)
	gLC17 Mem = global64("LC17", 1394253428160456350, 1878397350477548722)

	g_StepConstants = global32(
		"g_StepConstants",
		-1854099568, 1813713058, 1865754947, -814251453,
		753628274, 703164402, -1969511384, 787162690,
		237781025, -2027179250, -1537315662, 1190774290,
		408938142, 324624923, 641715378, 437348464,
		980181295, -1294024299, 46866262, 1086314584,
		2017887414, 1824226606, 1712094936, 729405578,
		-1496477591, -1851762873, -840272552, 9973912,
		-1093977298, 674802586, 1921734462, -1521327435,
		1952315919, 838311640, -1201680859, -1731671912,
		-1979291156, 1625572802, -35482192, -139126694,
		-765982333, 697726729, 405664733, 1636917552,
		-1871686026, 1160382498, -528988499, -1395821743,
		711036245, -1072857038, 1176621301, -218984047,
		13028614, 1865558631, 1488909453, 2051065085,
		-1947278721, -846351630, 1743963707, -444324989,
		-946613498, -1587768874, 397640165, -1157250441,
		2060542474, 1530963391, 1522074018, 1772595048,
		1533467085, -34707337, -885692676, -1060619726,
		1278438532, -1679399654, 330998780, 290394065,
		-1035950296, -329850764, 10249159, -1996329742,
		2141188304, -2108981323, -832659441, 1616830690,
		48746474, 1127699808, -1660933433, -1955635333,
		529580367, -842320585, 754888669, -1086340286,
		-357058068, 2056041891, -1658129820, -87106810,
		-1336799122, -1728025340, 732859657, -16579071,
		-1569610026, 121788701, -937601279, -911613440,
		40009502, -1720880739, -633629132, 34146304,
		337149304, 1234925092, -445359927, 1938547401,
		1692507424, 115281718, 366070542, 185665538,
		748025228, -451831187, 1506615982, -9524164,
		1182711172, -437369540, -409151197, 471403021,
		-1031044680, -156497838, 653250823, 1847257403,
		-985906742, -734821855, -2072642294, 891492137,
		1309489786, -1565238456, 381737005, -1996217175,
		25037967, 129344501, -93124722, 1480212574,
		1531425992, 1460642794, -677837501, -1926705614,
		2139790224, -1118201604, 1836750472, -1837457738,
		-1560511197, 1723043393, 1889421594, -1242103617,
		170273807, 384996793, -401531147, 219764040,
		-1619418939, 440520615, 250045322, -1393049484,
		810695449, 159928015, -114289443, 644783445,
		421160548, 1544648385, -162861672, -1531431648,
		-2104667063, -1848173096, 692384470, -1794379141,
		864078461, 1837643805, 1094223502, 1559057860,
		266444235, 1760884329, 1849149695, -1593938416,
		-1270134864, -171944566, 2042565839, 1245911584,
		-243308326, 1575706577, -1509898131, -1620279632,
		-7145372, -1240080257, 953410632, -1923270038,
		1894269899, 1229664558, -1496762861, 199962447,
		-1836360349, -873045963, 213396608, -359160585,
		1496059707, -1803788425, 1879001529, -149258662,
		501793013, -1034512896, -991646580, -2109956895,
	)
)

func global32(name string, values ...int64) Mem {
	mem := GLOBL(name, NOPTR|RODATA)

	for i, v := range values {
		DATA(i*4, I32(v))
	}

	return mem
}
func global64(name string, values ...int64) Mem {
	mem := GLOBL(name, NOPTR|RODATA)

	for i, v := range values {
		DATA(i*8, I64(v))
	}

	return mem
}

var (
	MOVDQA = MOVOA
	MOVDQU = MOVOU
	PSLLD  = PSLLDQ
	PSRLD  = PSRLDQ
)

func lsh256_sse2_init() {
	TEXT("lsh256_sse2_init", NOSPLIT, "func(ctx *lsh256ContextAsmData)")

	ctxMem, err := Param("ctx").Resolve()
	if err != nil {
		panic(err)
	}

	// *lsh256ContextAsmData
	rdi := RDI
	LEAQ(ctxMem.Addr, rdi)

	xmm0 := X0
	eax := EAX

	//lsh256_sse2_init(LSH256_Context*):
	MOVL(Mem{Base: rdi}, eax)              //        mov     eax, DWORD PTR [rdi]
	MOVL(U32(0), Mem{Base: rdi, Disp: 16}) //        mov     DWORD PTR [rdi+16], 0
	CMPL(eax, U32(28))                     //        cmp     eax, 28
	JE(LabelRef("L2"))                     //        je      .L2
	CMPL(eax, U32(32))                     //        cmp     eax, 32
	JE(LabelRef("L5"))                     //        je      .L5
	XORL(eax, eax)                         //        xor     eax, eax
	RET()                                  //        ret

	Label("L5")                            //.L5:
	MOVDQA(gLC0, xmm0)                     //        movdqa  xmm0, XMMWORD PTR .LC0[rip]
	XORL(eax, eax)                         //        xor     eax, eax
	MOVAPS(xmm0, Mem{Base: rdi, Disp: 32}) //        movaps  XMMWORD PTR [rdi+32], xmm0
	MOVDQA(gLC1, xmm0)                     //        movdqa  xmm0, XMMWORD PTR .LC1[rip]
	MOVAPS(xmm0, Mem{Base: rdi, Disp: 48}) //        movaps  XMMWORD PTR [rdi+48], xmm0
	MOVDQA(gLC2, xmm0)                     //        movdqa  xmm0, XMMWORD PTR .LC2[rip]
	MOVAPS(xmm0, Mem{Base: rdi, Disp: 64}) //        movaps  XMMWORD PTR [rdi+64], xmm0
	MOVDQA(gLC3, xmm0)                     //        movdqa  xmm0, XMMWORD PTR .LC3[rip]
	MOVAPS(xmm0, Mem{Base: rdi, Disp: 80}) //        movaps  XMMWORD PTR [rdi+80], xmm0
	RET()                                  //        ret

	Label("L2")                            //.L2:
	MOVDQA(gLC4, xmm0)                     //        movdqa  xmm0, XMMWORD PTR .LC4[rip]
	XORL(eax, eax)                         //        xor     eax, eax
	MOVAPS(xmm0, Mem{Base: rdi, Disp: 32}) //        movaps  XMMWORD PTR [rdi+32], xmm0
	MOVDQA(gLC5, xmm0)                     //        movdqa  xmm0, XMMWORD PTR .LC5[rip]
	MOVAPS(xmm0, Mem{Base: rdi, Disp: 48}) //        movaps  XMMWORD PTR [rdi+48], xmm0
	MOVDQA(gLC6, xmm0)                     //        movdqa  xmm0, XMMWORD PTR .LC6[rip]
	MOVAPS(xmm0, Mem{Base: rdi, Disp: 64}) //        movaps  XMMWORD PTR [rdi+64], xmm0
	MOVDQA(gLC7, xmm0)                     //        movdqa  xmm0, XMMWORD PTR .LC7[rip]
	MOVAPS(xmm0, Mem{Base: rdi, Disp: 80}) //        movaps  XMMWORD PTR [rdi+80], xmm0
	RET()                                  //        ret
}

func lsh256_sse2_update() {
	TEXT("lsh256_sse2_update", NOSPLIT, "func(ctx *lsh256ContextAsmData, data []byte, databitlen uint32)")

	ctxMem, err := Param("ctx").Resolve()
	if err != nil {
		panic(err)
	}
	dataMem, err := Param("data").Base().Resolve()
	if err != nil {
		panic(err)
	}

	//lsh256_sse2_update(LSH256_Context*, unsigned char const*, unsigned long):
	rdi := RDI
	LEAQ(ctxMem.Addr, rdi)
	rsi := RSI
	LEAQ(dataMem.Addr, rsi)
	rdx := Load(Param("databitlen"), RDX)

	r13 := R13
	r12 := R12
	rbp := RBP
	rbx := RBX
	rcx := RCX
	eax := EAX
	rax := RAX
	rsp := RSP
	edx := EDX
	esi := ESI
	dx := DX
	xmm0 := X0
	xmm1 := X1
	xmm2 := X2
	xmm3 := X3
	xmm4 := X4
	xmm5 := X5
	xmm6 := X6
	xmm7 := X7
	xmm8 := X8
	xmm9 := X9
	xmm10 := X10
	xmm11 := X11
	xmm12 := X12
	xmm13 := X13
	xmm14 := X14
	xmm15 := X15
	r12d := R12L
	r13d := R13L
	ecx := ECX

	PUSHQ(r13)                                      //        push    r13
	MOVQ(rdx, r13)                                  //        mov     r13, rdx
	PUSHQ(r12)                                      //        push    r12
	MOVQ(rdx, r12)                                  //        mov     r12, rdx
	PUSHQ(rbp)                                      //        push    rbp
	SHRQ(U8(3), r12)                                //        shr     r12, 3
	MOVQ(rdi, rbp)                                  //        mov     rbp, rdi
	PUSHQ(rbx)                                      //        push    rbx
	MOVQ(rsi, rbx)                                  //        mov     rbx, rsi
	SUBQ(U32(328), rsp)                             //        sub     rsp, 328
	MOVL(Mem{Base: rdi, Disp: 16}, eax)             //        mov     eax, DWORD PTR [rdi+16]
	SHRL(U8(3), eax)                                //        shr     eax, 3
	MOVL(eax, edx)                                  //        mov     edx, eax
	LEAQ(Mem{Base: rcx, Index: r12, Scale: 1}, rcx) //        lea     rcx, [rdx+r12]
	CMPQ(rcx, U32(127))                             //        cmp     rcx, 127
	JBE(LabelRef("L35"))                            //        jbe     .L35
	MOVOU(Mem{Base: rdi, Disp: 32}, xmm5)           //        movdqu  xmm5, XMMWORD PTR [rdi+32]
	MOVOU(Mem{Base: rdi, Disp: 48}, xmm1)           //        movdqu  xmm1, XMMWORD PTR [rdi+48]
	MOVOU(Mem{Base: rdi, Disp: 64}, xmm2)           //        movdqu  xmm2, XMMWORD PTR [rdi+64]
	MOVOU(Mem{Base: rdi, Disp: 80}, xmm7)           //        movdqu  xmm7, XMMWORD PTR [rdi+80]
	MOVDQA(xmm5, xmm4)                              //        movdqa  xmm4, xmm5
	MOVDQA(xmm1, xmm3)                              //        movdqa  xmm3, xmm1
	MOVDQA(xmm2, xmm0)                              //        movdqa  xmm0, xmm2
	MOVDQA(xmm7, xmm6)                              //        movdqa  xmm6, xmm7
	TESTL(eax, eax)                                 //        test    eax, eax
	JNE(LabelRef("L36"))                            //        jne     .L36

	Label("L9")                                                //.L9:
	CMPQ(r12, U32(127))                                        //        cmp     r12, 127
	JBE(LabelRef("L18"))                                       //        jbe     .L18
	MOVDQA(gLC8, xmm5)                                         //        movdqa  xmm5, XMMWORD PTR .LC8[rip]
	LEAQ(Mem{Base: R12, Disp: -128}, RAX)                      //        lea     rax, [r12-128]
	MOVDQA(xmm6, xmm10)                                        //        movdqa  xmm10, xmm6
	MOVDQA(gLC15, xmm9)                                        //        movdqa  xmm9, XMMWORD PTR .LC15[rip]
	ANDQ(I8(-128), rax)                                        //        and     rax, -128
	MOVAPS(xmm5, Mem{Base: RSP, Disp: 256})                    //        movaps  XMMWORD PTR [rsp+256], xmm5
	MOVDQA(gLC9, xmm5)                                         //        movdqa  xmm5, XMMWORD PTR .LC9[rip]
	LEAQ(Mem{Base: RSP, Disp: 256, Index: rax, Scale: 1}, rdx) //        lea     rdx, [rbx+128+rax]
	MOVAPS(xmm9, Mem{Base: RSP, Disp: 32})                     //        movaps  XMMWORD PTR [rsp+32], xmm9
	MOVDQA(xmm0, xmm9)                                         //        movdqa  xmm9, xmm0
	MOVAPS(xmm5, Mem{Base: RSP, Disp: 272})                    //        movaps  XMMWORD PTR [rsp+272], xmm5
	MOVDQA(gLC10, xmm5)                                        //        movdqa  xmm5, XMMWORD PTR .LC10[rip]
	MOVAPS(xmm5, Mem{Base: RSP, Disp: 48})                     //        movaps  XMMWORD PTR [rsp+48], xmm5
	MOVDQA(gLC11, xmm5)                                        //        movdqa  xmm5, XMMWORD PTR .LC11[rip]
	MOVAPS(xmm5, Mem{Base: RSP, Disp: 64})                     //        movaps  XMMWORD PTR [rsp+64], xmm5
	MOVDQA(gLC12, xmm5)                                        //        movdqa  xmm5, XMMWORD PTR .LC12[rip]
	MOVAPS(xmm5, Mem{Base: RSP, Disp: 16})                     //        movaps  XMMWORD PTR [rsp+16], xmm5
	MOVDQA(gLC13, xmm5)                                        //        movdqa  xmm5, XMMWORD PTR .LC13[rip]
	MOVAPS(xmm5, Mem{Base: RSP})                               //        movaps  XMMWORD PTR [rsp], xmm5
	MOVDQA(gLC14, xmm5)                                        //        movdqa  xmm5, XMMWORD PTR .LC14[rip]
	MOVAPS(xmm5, Mem{Base: RSP, Disp: 80})                     //        movaps  XMMWORD PTR [rsp+80], xmm5
	MOVDQA(gLC16, xmm5)                                        //        movdqa  xmm5, XMMWORD PTR .LC16[rip]
	MOVAPS(xmm5, Mem{Base: RSP, Disp: 288})                    //        movaps  XMMWORD PTR [rsp+288], xmm5
	MOVDQA(gLC17, xmm5)                                        //        movdqa  xmm5, XMMWORD PTR .LC17[rip]
	MOVAPS(xmm5, Mem{Base: RSP, Disp: 304})                    //        movaps  XMMWORD PTR [rsp+304], xmm5

	Label("L13")                             //.L13:
	MOVDQU(Mem{Base: rbx}, xmm5)             //        movdqu  xmm5, XMMWORD PTR [rbx]
	MOVDQU(Mem{Base: rbx, Disp: 32}, xmm6)   //        movdqu  xmm6, XMMWORD PTR [rbx+32]
	MOVL(g_StepConstants.Offset(64), eax)    //        mov     eax, OFFSET FLAT:g_StepConstants+64
	MOVDQU(Mem{Base: rbx, Disp: 16}, xmm1)   //        movdqu  xmm1, XMMWORD PTR [rbx+16]
	MOVDQU(Mem{Base: rbx, Disp: 48}, xmm2)   //        movdqu  xmm2, XMMWORD PTR [rbx+48]
	PXOR(xmm6, xmm9)                         //        pxor    xmm9, xmm6
	PXOR(xmm5, xmm4)                         //        pxor    xmm4, xmm5
	MOVDQU(Mem{Base: rbx, Disp: 96}, xmm13)  //        movdqu  xmm13, XMMWORD PTR [rbx+96]
	MOVDQU(Mem{Base: rbx, Disp: 80}, xmm0)   //        movdqu  xmm0, XMMWORD PTR [rbx+80]
	PXOR(xmm2, xmm10)                        //        pxor    xmm10, xmm2
	PADDD(xmm9, xmm4)                        //        paddd   xmm4, xmm9
	PXOR(xmm1, xmm3)                         //        pxor    xmm3, xmm1
	MOVAPS(xmm2, Mem{Base: rsp, Disp: 240})  //        movaps  XMMWORD PTR [rsp+240], xmm2
	PADDD(xmm10, xmm3)                       //        paddd   xmm3, xmm10
	MOVDQA(xmm4, xmm2)                       //        movdqa  xmm2, xmm4
	MOVAPS(xmm1, Mem{Base: rsp, Disp: 208})  //        movaps  XMMWORD PTR [rsp+208], xmm1
	MOVDQA(Mem{Base: rsp, Disp: 16}, xmm15)  //        movdqa  xmm15, XMMWORD PTR [rsp+16]
	PSRLDQ(I8(3), xmm2)                      //        psrld   xmm2, 3
	MOVDQA(xmm3, xmm1)                       //        movdqa  xmm1, xmm3
	MOVAPS(xmm13, Mem{Base: rsp, Disp: 128}) //        movaps  XMMWORD PTR [rsp+128], xmm13
	MOVDQA(Mem{Base: rsp, Disp: 48}, xmm13)  //        movdqa  xmm13, XMMWORD PTR [rsp+48]
	PSLLD(I8(29), xmm4)                      //        pslld   xmm4, 29
	PSRLD(I8(3), xmm1)                       //        psrld   xmm1, 3
	MOVAPS(xmm0, Mem{Base: rsp, Disp: 112})  //        movaps  XMMWORD PTR [rsp+112], xmm0
	MOVDQU(Mem{Base: rbx, Disp: 112}, xmm14) //        movdqu  xmm14, XMMWORD PTR [rbx+112]
	POR(xmm2, xmm4)                          //        por     xmm4, xmm2
	PSLLD(I8(29), xmm3)                      //        pslld   xmm3, 29
	PXOR(Mem{Base: rsp, Disp: 256}, xmm4)    //        pxor    xmm4, XMMWORD PTR [rsp+256]
	MOVDQA(Mem{Base: rsp}, xmm12)            //        movdqa  xmm12, XMMWORD PTR [rsp]
	POR(xmm1, xmm3)                          //        por     xmm3, xmm1
	MOVDQA(xmm13, xmm1)                      //        movdqa  xmm1, xmm13
	PXOR(Mem{Base: rsp, Disp: 272}, xmm3)    //        pxor    xmm3, XMMWORD PTR [rsp+272]
	MOVAPS(xmm14, Mem{Base: rsp, Disp: 144}) //        movaps  XMMWORD PTR [rsp+144], xmm14
	PADDD(xmm4, xmm9)                        //        paddd   xmm9, xmm4
	MOVDQU(Mem{Base: rbx, Disp: 64}, xmm7)   //        movdqu  xmm7, XMMWORD PTR [rbx+64]
	MOVDQA(Mem{Base: rsp, Disp: 80}, xmm14)  //        movdqa  xmm14, XMMWORD PTR [rsp+80]
	MOVDQA(xmm4, xmm8)                       //        movdqa  xmm8, xmm4
	MOVDQA(xmm3, xmm11)                      //        movdqa  xmm11, xmm3
	MOVDQA(xmm9, xmm3)                       //        movdqa  xmm3, xmm9
	MOVAPS(xmm5, Mem{Base: rsp, Disp: 192})  //        movaps  XMMWORD PTR [rsp+192], xmm5
	PADDD(xmm11, xmm10)                      //        paddd   xmm10, xmm11
	PSRLD(I8(31), xmm3)                      //        psrld   xmm3, 31
	MOVAPS(xmm7, Mem{Base: rsp, Disp: 96})   //        movaps  XMMWORD PTR [rsp+96], xmm7
	PSLLD(I8(1), xmm9)                       //        pslld   xmm9, 1
	MOVAPS(xmm6, Mem{Base: rsp, Disp: 224})  //        movaps  XMMWORD PTR [rsp+224], xmm6
	POR(xmm3, xmm9)                          //        por     xmm9, xmm3
	MOVDQA(xmm10, xmm3)                      //        movdqa  xmm3, xmm10
	PSRLD(I8(31), xmm3)                      //        psrld   xmm3, 31
	PSLLD(I8(1), xmm10)                      //        pslld   xmm10, 1
	PAND(xmm9, xmm1)                         //        pand    xmm1, xmm9
	MOVDQA(xmm3, xmm0)                       //        movdqa  xmm0, xmm3
	MOVDQA(xmm10, xmm3)                      //        movdqa  xmm3, xmm10
	MOVDQA(Mem{Base: rsp, Disp: 64}, xmm10)  //        movdqa  xmm10, XMMWORD PTR [rsp+64]
	MOVDQA(xmm1, xmm2)                       //        movdqa  xmm2, xmm1
	PSRLD(I8(24), xmm2)                      //        psrld   xmm2, 24
	PSLLD(I8(8), xmm1)                       //        pslld   xmm1, 8
	POR(xmm0, xmm3)                          //        por     xmm3, xmm0
	MOVDQA(xmm10, xmm0)                      //        movdqa  xmm0, xmm10
	PXOR(xmm2, xmm1)                         //        pxor    xmm1, xmm2
	PAND(xmm9, xmm0)                         //        pand    xmm0, xmm9
	PXOR(xmm1, xmm0)                         //        pxor    xmm0, xmm1
	MOVDQA(xmm15, xmm1)                      //        movdqa  xmm1, xmm15
	PAND(xmm0, xmm1)                         //        pand    xmm1, xmm0
	PAND(xmm12, xmm0)                        //        pand    xmm0, xmm12
	MOVDQA(xmm1, xmm2)                       //        movdqa  xmm2, xmm1
	PSLLD(I8(8), xmm1)                       //        pslld   xmm1, 8
	PSRLD(I8(24), xmm2)                      //        psrld   xmm2, 24
	PXOR(xmm2, xmm1)                         //        pxor    xmm1, xmm2
	PXOR(xmm1, xmm0)                         //        pxor    xmm0, xmm1
	MOVDQA(xmm14, xmm1)                      //        movdqa  xmm1, xmm14
	PAND(xmm0, xmm1)                         //        pand    xmm1, xmm0
	MOVDQA(xmm1, xmm2)                       //        movdqa  xmm2, xmm1
	MOVDQA(xmm1, xmm5)                       //        movdqa  xmm5, xmm1
	MOVDQA(xmm14, xmm1)                      //        movdqa  xmm1, xmm14
	PSRLD(I8(24), xmm2)                      //        psrld   xmm2, 24
	PSLLD(I8(8), xmm5)                       //        pslld   xmm5, 8
	PAND(xmm3, xmm1)                         //        pand    xmm1, xmm3
	MOVDQA(xmm2, xmm7)                       //        movdqa  xmm7, xmm2
	MOVDQA(Mem{Base: rsp, Disp: 32}, xmm2)   //        movdqa  xmm2, XMMWORD PTR [rsp+32]
	PXOR(xmm7, xmm5)                         //        pxor    xmm5, xmm7
	PAND(xmm3, xmm2)                         //        pand    xmm2, xmm3
	PADDD(xmm11, xmm3)                       //        paddd   xmm3, xmm11
	MOVDQA(Mem{Base: rsp, Disp: 32}, xmm11)  //        movdqa  xmm11, XMMWORD PTR [rsp+32]
	MOVDQA(xmm2, xmm4)                       //        movdqa  xmm4, xmm2
	PSLLD(I8(8), xmm2)                       //        pslld   xmm2, 8
	PSHUFD(U8(210), xmm3, xmm3)              //        pshufd  xmm3, xmm3, 210
	PSRLD(I8(24), xmm4)                      //        psrld   xmm4, 24
	PAND(xmm11, xmm0)                        //        pand    xmm0, xmm11
	PXOR(Mem{Base: rsp, Disp: 96}, xmm3)     //        pxor    xmm3, XMMWORD PTR [rsp+96]
	PXOR(xmm4, xmm2)                         //        pxor    xmm2, xmm4
	PXOR(xmm5, xmm0)                         //        pxor    xmm0, xmm5
	PXOR(xmm2, xmm1)                         //        pxor    xmm1, xmm2
	MOVDQA(xmm12, xmm2)                      //        movdqa  xmm2, xmm12
	PSHUFD(I8(108), xmm0, xmm0)              //        pshufd  xmm0, xmm0, 108
	PXOR(Mem{Base: rsp, Disp: 144}, xmm0)    //        pxor    xmm0, XMMWORD PTR [rsp+144]
	PAND(xmm1, xmm2)                         //        pand    xmm2, xmm1
	PAND(xmm15, xmm1)                        //        pand    xmm1, xmm15
	MOVDQA(xmm2, xmm4)                       //        movdqa  xmm4, xmm2
	PSLLD(I8(8), xmm2)                       //        pslld   xmm2, 8
	PSRLD(I8(24), xmm4)                      //        psrld   xmm4, 24
	PXOR(xmm4, xmm2)                         //        pxor    xmm2, xmm4
	MOVDQA(xmm10, xmm4)                      //        movdqa  xmm4, xmm10
	PXOR(xmm2, xmm1)                         //        pxor    xmm1, xmm2
	MOVDQA(xmm8, xmm2)                       //        movdqa  xmm2, xmm8
	PAND(xmm1, xmm4)                         //        pand    xmm4, xmm1
	PADDD(xmm9, xmm2)                        //        paddd   xmm2, xmm9
	PAND(xmm13, xmm1)                        //        pand    xmm1, xmm13
	MOVDQA(xmm4, xmm6)                       //        movdqa  xmm6, xmm4
	PSLLD(I8(8), xmm4)                       //        pslld   xmm4, 8
	PSHUFD(U8(210), xmm2, xmm2)              //        pshufd  xmm2, xmm2, 210
	PSRLD(I8(24), xmm6)                      //        psrld   xmm6, 24
	MOVDQA(xmm13, xmm9)                      //        movdqa  xmm9, xmm13
	PXOR(Mem{Base: rsp, Disp: 128}, xmm2)    //        pxor    xmm2, XMMWORD PTR [rsp+128]
	PXOR(xmm6, xmm4)                         //        pxor    xmm4, xmm6
	MOVDQA(xmm10, xmm6)                      //        movdqa  xmm6, xmm10
	PXOR(xmm4, xmm1)                         //        pxor    xmm1, xmm4
	PADDD(xmm2, xmm3)                        //        paddd   xmm3, xmm2
	PSHUFD(I8(108), xmm1, xmm1)              //        pshufd  xmm1, xmm1, 108
	PXOR(Mem{Base: rsp, Disp: 112}, xmm1)    //        pxor    xmm1, XMMWORD PTR [rsp+112]
	MOVDQA(xmm3, xmm4)                       //        movdqa  xmm4, xmm3
	PSRLD(I8(27), xmm4)                      //        psrld   xmm4, 27
	PSLLD(I8(5), xmm3)                       //        pslld   xmm3, 5
	PADDD(xmm0, xmm1)                        //        paddd   xmm1, xmm0
	POR(xmm4, xmm3)                          //        por     xmm3, xmm4
	PXOR(Mem{Base: rsp, Disp: 288}, xmm3)    //        pxor    xmm3, XMMWORD PTR [rsp+288]
	MOVDQA(xmm1, xmm5)                       //        movdqa  xmm5, xmm1
	PSLLD(I8(5), xmm1)                       //        pslld   xmm1, 5
	PSRLD(I8(27), xmm5)                      //        psrld   xmm5, 27
	PADDD(xmm3, xmm2)                        //        paddd   xmm2, xmm3
	MOVDQA(xmm3, xmm4)                       //        movdqa  xmm4, xmm3
	POR(xmm5, xmm1)                          //        por     xmm1, xmm5
	PXOR(Mem{Base: rsp, Disp: 304}, xmm1)    //        pxor    xmm1, XMMWORD PTR [rsp+304]
	MOVDQA(xmm1, xmm8)                       //        movdqa  xmm8, xmm1
	MOVDQA(xmm2, xmm1)                       //        movdqa  xmm1, xmm2
	PSRLD(I8(15), xmm1)                      //        psrld   xmm1, 15
	PSLLD(I8(17), xmm2)                      //        pslld   xmm2, 17
	PADDD(xmm8, xmm0)                        //        paddd   xmm0, xmm8
	POR(xmm1, xmm2)                          //        por     xmm2, xmm1
	MOVDQA(xmm0, xmm1)                       //        movdqa  xmm1, xmm0
	PSLLD(I8(17), xmm0)                      //        pslld   xmm0, 17
	PSRLD(I8(15), xmm1)                      //        psrld   xmm1, 15
	PAND(xmm2, xmm6)                         //        pand    xmm6, xmm2
	MOVDQA(xmm0, xmm5)                       //        movdqa  xmm5, xmm0
	MOVDQA(xmm13, xmm0)                      //        movdqa  xmm0, xmm13
	MOVDQA(xmm12, xmm13)                     //        movdqa  xmm13, xmm12
	PAND(xmm2, xmm0)                         //        pand    xmm0, xmm2
	POR(xmm1, xmm5)                          //        por     xmm5, xmm1
	PADDD(xmm2, xmm4)                        //        paddd   xmm4, xmm2
	MOVDQA(xmm0, xmm1)                       //        movdqa  xmm1, xmm0
	PSLLD(I8(8), xmm0)                       //        pslld   xmm0, 8
	PSHUFD(U8(210), xmm4, xmm4)              //        pshufd  xmm4, xmm4, 210
	PSRLD(I8(24), xmm1)                      //        psrld   xmm1, 24
	PXOR(xmm1, xmm0)                         //        pxor    xmm0, xmm1
	PXOR(xmm0, xmm6)                         //        pxor    xmm6, xmm0
	MOVDQA(xmm15, xmm0)                      //        movdqa  xmm0, xmm15
	PAND(xmm6, xmm0)                         //        pand    xmm0, xmm6
	PAND(xmm12, xmm6)                        //        pand    xmm6, xmm12
	MOVDQA(xmm0, xmm1)                       //        movdqa  xmm1, xmm0
	PSLLD(I8(8), xmm0)                       //        pslld   xmm0, 8
	PSRLD(I8(24), xmm1)                      //        psrld   xmm1, 24
	PXOR(xmm1, xmm0)                         //        pxor    xmm0, xmm1
	PXOR(xmm0, xmm6)                         //        pxor    xmm6, xmm0
	MOVDQA(xmm14, xmm0)                      //        movdqa  xmm0, xmm14
	PAND(xmm5, xmm14)                        //        pand    xmm14, xmm5
	PAND(xmm6, xmm0)                         //        pand    xmm0, xmm6
	MOVDQA(xmm14, xmm12)                     //        movdqa  xmm12, xmm14
	PAND(xmm11, xmm6)                        //        pand    xmm6, xmm11
	MOVDQA(xmm0, xmm7)                       //        movdqa  xmm7, xmm0
	MOVDQA(xmm0, xmm3)                       //        movdqa  xmm3, xmm0
	MOVDQA(xmm11, xmm0)                      //        movdqa  xmm0, xmm11
	PSLLD(I8(8), xmm3)                       //        pslld   xmm3, 8
	PSRLD(I8(24), xmm7)                      //        psrld   xmm7, 24
	PAND(xmm5, xmm0)                         //        pand    xmm0, xmm5
	MOVDQA(xmm0, xmm1)                       //        movdqa  xmm1, xmm0
	PSLLD(I8(8), xmm0)                       //        pslld   xmm0, 8
	MOVDQA(xmm3, xmm2)                       //        movdqa  xmm2, xmm3
	PSRLD(I8(24), xmm1)                      //        psrld   xmm1, 24
	PXOR(xmm7, xmm2)                         //        pxor    xmm2, xmm7
	PADDD(xmm8, xmm5)                        //        paddd   xmm5, xmm8
	PXOR(xmm1, xmm0)                         //        pxor    xmm0, xmm1
	PXOR(xmm2, xmm6)                         //        pxor    xmm6, xmm2
	MOVDQA(xmm4, xmm3)                       //        movdqa  xmm3, xmm4
	PXOR(xmm0, xmm12)                        //        pxor    xmm12, xmm0
	MOVDQA(xmm13, xmm0)                      //        movdqa  xmm0, xmm13
	PSHUFD(I8(108), xmm6, xmm6)              //        pshufd  xmm6, xmm6, 108
	PAND(xmm12, xmm0)                        //        pand    xmm0, xmm12
	PAND(xmm15, xmm12)                       //        pand    xmm12, xmm15
	PSHUFD(U8(210), xmm5, xmm5)              //        pshufd  xmm5, xmm5, 210
	MOVDQA(xmm0, xmm1)                       //        movdqa  xmm1, xmm0
	PSLLD(I8(8), xmm0)                       //        pslld   xmm0, 8
	MOVDQA(xmm6, xmm4)                       //        movdqa  xmm4, xmm6
	PSRLD(I8(24), xmm1)                      //        psrld   xmm1, 24
	PXOR(xmm1, xmm0)                         //        pxor    xmm0, xmm1
	PXOR(xmm0, xmm12)                        //        pxor    xmm12, xmm0
	PAND(xmm12, xmm10)                       //        pand    xmm10, xmm12
	PAND(xmm9, xmm12)                        //        pand    xmm12, xmm9
	MOVDQA(xmm10, xmm0)                      //        movdqa  xmm0, xmm10
	MOVDQA(xmm10, xmm1)                      //        movdqa  xmm1, xmm10
	PSRLD(I8(24), xmm1)                      //        psrld   xmm1, 24
	PSLLD(I8(8), xmm0)                       //        pslld   xmm0, 8
	PXOR(xmm1, xmm0)                         //        pxor    xmm0, xmm1
	PXOR(xmm0, xmm12)                        //        pxor    xmm12, xmm0
	PSHUFD(I8(108), xmm12, xmm12)            //        pshufd  xmm12, xmm12, 108

	Label("L12") //.L12:
	/**
	MOVDQA(xmm12, xmm1) //        movdqa  xmm1, xmm12
	MOVDQA(XMMWORD PTR [rsp+48], xmm14) //        movdqa  xmm14, XMMWORD PTR [rsp+48]
	ADD(64, rax) //        add     rax, 64
	PSHUFD(75, xmm0, XMMWORD PTR [rsp+192]) //        pshufd  xmm0, XMMWORD PTR [rsp+192], 75
	PADDD(XMMWORD PTR [rsp+96], xmm0) //        paddd   xmm0, XMMWORD PTR [rsp+96]
	MOVDQA(XMMWORD PTR [rsp+64], xmm12) //        movdqa  xmm12, XMMWORD PTR [rsp+64]
	PSHUFD(75, xmm9, XMMWORD PTR [rsp+224]) //        pshufd  xmm9, XMMWORD PTR [rsp+224], 75
	PSHUFD(147, xmm10, XMMWORD PTR [rsp+240]) //        pshufd  xmm10, XMMWORD PTR [rsp+240], 147
	PADDD(XMMWORD PTR [rsp+128], xmm9) //        paddd   xmm9, XMMWORD PTR [rsp+128]
	MOVDQA(XMMWORD PTR [rsp+16], xmm13) //        movdqa  xmm13, XMMWORD PTR [rsp+16]
	MOVDQA(xmm0, xmm2) //        movdqa  xmm2, xmm0
	MOVAPS(xmm0, XMMWORD PTR [rsp+192]) //        movaps  XMMWORD PTR [rsp+192], xmm0
	PSHUFD(147, xmm0, XMMWORD PTR [rsp+208]) //        pshufd  xmm0, XMMWORD PTR [rsp+208], 147
	PADDD(XMMWORD PTR [rsp+112], xmm0) //        paddd   xmm0, XMMWORD PTR [rsp+112]
	MOVAPS(xmm2, XMMWORD PTR [rsp+160]) //        movaps  XMMWORD PTR [rsp+160], xmm2
	PXOR(xmm9, xmm3) //        pxor    xmm3, xmm9
	PXOR(XMMWORD PTR [rsp+160], xmm5) //        pxor    xmm5, XMMWORD PTR [rsp+160]
	PADDD(XMMWORD PTR [rsp+144], xmm10) //        paddd   xmm10, XMMWORD PTR [rsp+144]
	MOVAPS(xmm0, XMMWORD PTR [rsp+176]) //        movaps  XMMWORD PTR [rsp+176], xmm0
	MOVDQA(XMMWORD PTR [rsp], xmm15) //        movdqa  xmm15, XMMWORD PTR [rsp]
	PXOR(XMMWORD PTR [rsp+176], xmm1) //        pxor    xmm1, XMMWORD PTR [rsp+176]
	MOVDQA(xmm5, xmm6) //        movdqa  xmm6, xmm5
	PXOR(xmm10, xmm4) //        pxor    xmm4, xmm10
	MOVAPS(xmm0, XMMWORD PTR [rsp+208]) //        movaps  XMMWORD PTR [rsp+208], xmm0
	MOVDQA(XMMWORD PTR [rsp+32], xmm7) //        movdqa  xmm7, XMMWORD PTR [rsp+32]
	PADDD(xmm3, xmm6) //        paddd   xmm6, xmm3
	PADDD(xmm4, xmm1) //        paddd   xmm1, xmm4
	MOVAPS(xmm9, XMMWORD PTR [rsp+224]) //        movaps  XMMWORD PTR [rsp+224], xmm9
	MOVDQA(xmm6, xmm2) //        movdqa  xmm2, xmm6
	PSLLD(29, xmm6) //        pslld   xmm6, 29
	MOVDQA(xmm1, xmm0) //        movdqa  xmm0, xmm1
	MOVAPS(xmm10, XMMWORD PTR [rsp+240]) //        movaps  XMMWORD PTR [rsp+240], xmm10
	PSRLD(3, xmm2) //        psrld   xmm2, 3
	PSRLD(3, xmm0) //        psrld   xmm0, 3
	POR(xmm2, xmm6) //        por     xmm6, xmm2
	PSLLD(29, xmm1) //        pslld   xmm1, 29
	PXOR(XMMWORD PTR [rax-64], xmm6) //        pxor    xmm6, XMMWORD PTR [rax-64]
	POR(xmm0, xmm1) //        por     xmm1, xmm0
	PXOR(XMMWORD PTR [rax-48], xmm1) //        pxor    xmm1, XMMWORD PTR [rax-48]
	PADDD(xmm6, xmm3) //        paddd   xmm3, xmm6
	MOVDQA(xmm3, xmm0) //        movdqa  xmm0, xmm3
	PSLLD(1, xmm3) //        pslld   xmm3, 1
	PADDD(xmm1, xmm4) //        paddd   xmm4, xmm1
	PSRLD(31, xmm0) //        psrld   xmm0, 31
	POR(xmm0, xmm3) //        por     xmm3, xmm0
	MOVDQA(xmm4, xmm0) //        movdqa  xmm0, xmm4
	PSRLD(31, xmm0) //        psrld   xmm0, 31
	PSLLD(1, xmm4) //        pslld   xmm4, 1
	PADDD(xmm3, xmm6) //        paddd   xmm6, xmm3
	POR(xmm0, xmm4) //        por     xmm4, xmm0
	MOVDQA(xmm14, xmm0) //        movdqa  xmm0, xmm14
	PAND(xmm3, xmm0) //        pand    xmm0, xmm3
	PAND(xmm4, xmm7) //        pand    xmm7, xmm4
	PADDD(xmm4, xmm1) //        paddd   xmm1, xmm4
	MOVDQA(xmm0, xmm5) //        movdqa  xmm5, xmm0
	PSLLD(8, xmm0) //        pslld   xmm0, 8
	MOVDQA(xmm7, xmm11) //        movdqa  xmm11, xmm7
	PSRLD(24, xmm5) //        psrld   xmm5, 24
	PSRLD(24, xmm11) //        psrld   xmm11, 24
	PSHUFD(210, xmm1, xmm1) //        pshufd  xmm1, xmm1, 210
	MOVDQA(xmm5, xmm2) //        movdqa  xmm2, xmm5
	PSLLD(8, xmm7) //        pslld   xmm7, 8
	MOVDQA(xmm12, xmm5) //        movdqa  xmm5, xmm12
	PXOR(xmm2, xmm0) //        pxor    xmm0, xmm2
	PAND(xmm3, xmm5) //        pand    xmm5, xmm3
	PXOR(xmm11, xmm7) //        pxor    xmm7, xmm11
	PXOR(xmm0, xmm5) //        pxor    xmm5, xmm0
	MOVDQA(xmm13, xmm0) //        movdqa  xmm0, xmm13
	PSHUFD(210, xmm3, xmm6) //        pshufd  xmm3, xmm6, 210
	PAND(xmm5, xmm0) //        pand    xmm0, xmm5
	PAND(xmm15, xmm5) //        pand    xmm5, xmm15
	MOVDQA(xmm0, xmm2) //        movdqa  xmm2, xmm0
	PSLLD(8, xmm0) //        pslld   xmm0, 8
	PSRLD(24, xmm2) //        psrld   xmm2, 24
	PXOR(xmm2, xmm0) //        pxor    xmm0, xmm2
	PXOR(xmm0, xmm5) //        pxor    xmm5, xmm0
	MOVDQA(XMMWORD PTR [rsp+80], xmm0) //        movdqa  xmm0, XMMWORD PTR [rsp+80]
	MOVDQA(xmm0, xmm8) //        movdqa  xmm8, xmm0
	PAND(xmm4, xmm0) //        pand    xmm0, xmm4
	PXOR(xmm7, xmm0) //        pxor    xmm0, xmm7
	MOVDQA(xmm15, xmm7) //        movdqa  xmm7, xmm15
	PAND(xmm5, xmm8) //        pand    xmm8, xmm5
	PSHUFD(75, xmm15, XMMWORD PTR [rsp+96]) //        pshufd  xmm15, XMMWORD PTR [rsp+96], 75
	PAND(xmm0, xmm7) //        pand    xmm7, xmm0
	PAND(xmm13, xmm0) //        pand    xmm0, xmm13
	MOVDQA(xmm8, xmm2) //        movdqa  xmm2, xmm8
	MOVDQA(xmm7, xmm11) //        movdqa  xmm11, xmm7
	PSLLD(8, xmm7) //        pslld   xmm7, 8
	PAND(XMMWORD PTR [rsp+32], xmm5) //        pand    xmm5, XMMWORD PTR [rsp+32]
	PSRLD(24, xmm11) //        psrld   xmm11, 24
	PSRLD(24, xmm2) //        psrld   xmm2, 24
	PADDD(XMMWORD PTR [rsp+160], xmm15) //        paddd   xmm15, XMMWORD PTR [rsp+160]
	PXOR(xmm11, xmm7) //        pxor    xmm7, xmm11
	PSLLD(8, xmm8) //        pslld   xmm8, 8
	PXOR(xmm7, xmm0) //        pxor    xmm0, xmm7
	MOVDQA(xmm12, xmm7) //        movdqa  xmm7, xmm12
	PXOR(xmm2, xmm8) //        pxor    xmm8, xmm2
	MOVAPS(xmm15, XMMWORD PTR [rsp+96]) //        movaps  XMMWORD PTR [rsp+96], xmm15
	PAND(xmm0, xmm7) //        pand    xmm7, xmm0
	PAND(xmm14, xmm0) //        pand    xmm0, xmm14
	PXOR(xmm8, xmm5) //        pxor    xmm5, xmm8
	MOVDQA(xmm7, xmm11) //        movdqa  xmm11, xmm7
	PSLLD(8, xmm7) //        pslld   xmm7, 8
	MOVDQA(xmm14, xmm8) //        movdqa  xmm8, xmm14
	PSHUFD(147, xmm14, XMMWORD PTR [rsp+112]) //        pshufd  xmm14, XMMWORD PTR [rsp+112], 147
	PSRLD(24, xmm11) //        psrld   xmm11, 24
	PXOR(xmm15, xmm1) //        pxor    xmm1, xmm15
	PSHUFD(108, xmm5, xmm5) //        pshufd  xmm5, xmm5, 108
	PADDD(XMMWORD PTR [rsp+176], xmm14) //        paddd   xmm14, XMMWORD PTR [rsp+176]
	PXOR(xmm11, xmm7) //        pxor    xmm7, xmm11
	PXOR(xmm7, xmm0) //        pxor    xmm0, xmm7
	PSHUFD(108, xmm2, xmm0) //        pshufd  xmm2, xmm0, 108
	MOVAPS(xmm14, XMMWORD PTR [rsp+112]) //        movaps  XMMWORD PTR [rsp+112], xmm14
	PSHUFD(75, xmm0, XMMWORD PTR [rsp+128]) //        pshufd  xmm0, XMMWORD PTR [rsp+128], 75
	PADDD(xmm9, xmm0) //        paddd   xmm0, xmm9
	PXOR(xmm0, xmm3) //        pxor    xmm3, xmm0
	MOVAPS(xmm0, XMMWORD PTR [rsp+128]) //        movaps  XMMWORD PTR [rsp+128], xmm0
	PXOR(xmm14, xmm2) //        pxor    xmm2, xmm14
	PSHUFD(147, xmm13, XMMWORD PTR [rsp+144]) //        pshufd  xmm13, XMMWORD PTR [rsp+144], 147
	PADDD(xmm3, xmm1) //        paddd   xmm1, xmm3
	PADDD(xmm10, xmm13) //        paddd   xmm13, xmm10
	MOVDQA(XMMWORD PTR [rsp+16], xmm7) //        movdqa  xmm7, XMMWORD PTR [rsp+16]
	PXOR(xmm13, xmm5) //        pxor    xmm5, xmm13
	MOVDQA(xmm1, xmm6) //        movdqa  xmm6, xmm1
	MOVAPS(xmm13, XMMWORD PTR [rsp+144]) //        movaps  XMMWORD PTR [rsp+144], xmm13
	PSRLD(27, xmm6) //        psrld   xmm6, 27
	PSLLD(5, xmm1) //        pslld   xmm1, 5
	PADDD(xmm5, xmm2) //        paddd   xmm2, xmm5
	MOVDQA(xmm2, xmm4) //        movdqa  xmm4, xmm2
	PSLLD(5, xmm2) //        pslld   xmm2, 5
	POR(xmm6, xmm1) //        por     xmm1, xmm6
	PXOR(XMMWORD PTR [rax-32], xmm1) //        pxor    xmm1, XMMWORD PTR [rax-32]
	PSRLD(27, xmm4) //        psrld   xmm4, 27
	MOVDQA(xmm12, xmm6) //        movdqa  xmm6, xmm12
	MOVDQA(XMMWORD PTR [rsp+32], xmm12) //        movdqa  xmm12, XMMWORD PTR [rsp+32]
	POR(xmm4, xmm2) //        por     xmm2, xmm4
	PXOR(XMMWORD PTR [rax-16], xmm2) //        pxor    xmm2, XMMWORD PTR [rax-16]
	PADDD(xmm1, xmm3) //        paddd   xmm3, xmm1
	MOVDQA(xmm3, xmm4) //        movdqa  xmm4, xmm3
	PSLLD(17, xmm3) //        pslld   xmm3, 17
	PADDD(xmm2, xmm5) //        paddd   xmm5, xmm2
	PSRLD(15, xmm4) //        psrld   xmm4, 15
	POR(xmm4, xmm3) //        por     xmm3, xmm4
	MOVDQA(xmm5, xmm4) //        movdqa  xmm4, xmm5
	PSRLD(15, xmm4) //        psrld   xmm4, 15
	PSLLD(17, xmm5) //        pslld   xmm5, 17
	PAND(xmm3, xmm6) //        pand    xmm6, xmm3
	POR(xmm4, xmm5) //        por     xmm5, xmm4
	MOVDQA(xmm8, xmm4) //        movdqa  xmm4, xmm8
	MOVDQA(XMMWORD PTR [rsp+80], xmm8) //        movdqa  xmm8, XMMWORD PTR [rsp+80]
	PADDD(xmm3, xmm1) //        paddd   xmm1, xmm3
	PAND(xmm3, xmm4) //        pand    xmm4, xmm3
	PAND(xmm5, xmm12) //        pand    xmm12, xmm5
	PADDD(xmm5, xmm2) //        paddd   xmm2, xmm5
	MOVDQA(xmm4, xmm11) //        movdqa  xmm11, xmm4
	PSLLD(8, xmm4) //        pslld   xmm4, 8
	PSHUFD(210, xmm1, xmm1) //        pshufd  xmm1, xmm1, 210
	PSRLD(24, xmm11) //        psrld   xmm11, 24
	PSHUFD(210, xmm2, xmm2) //        pshufd  xmm2, xmm2, 210
	MOVDQA(xmm1, xmm3) //        movdqa  xmm3, xmm1
	PXOR(xmm11, xmm4) //        pxor    xmm4, xmm11
	PXOR(xmm4, xmm6) //        pxor    xmm6, xmm4
	PAND(xmm6, xmm7) //        pand    xmm7, xmm6
	PAND(XMMWORD PTR [rsp], xmm6) //        pand    xmm6, XMMWORD PTR [rsp]
	MOVDQA(xmm7, xmm4) //        movdqa  xmm4, xmm7
	MOVDQA(xmm7, xmm11) //        movdqa  xmm11, xmm7
	MOVDQA(xmm8, xmm7) //        movdqa  xmm7, xmm8
	PSRLD(24, xmm11) //        psrld   xmm11, 24
	PSLLD(8, xmm4) //        pslld   xmm4, 8
	PAND(xmm5, xmm8) //        pand    xmm8, xmm5
	PXOR(xmm11, xmm4) //        pxor    xmm4, xmm11
	MOVDQA(xmm12, xmm11) //        movdqa  xmm11, xmm12
	MOVDQA(xmm2, xmm5) //        movdqa  xmm5, xmm2
	PSRLD(24, xmm12) //        psrld   xmm12, 24
	PSLLD(8, xmm11) //        pslld   xmm11, 8
	PXOR(xmm4, xmm6) //        pxor    xmm6, xmm4
	PXOR(xmm12, xmm11) //        pxor    xmm11, xmm12
	PAND(xmm6, xmm7) //        pand    xmm7, xmm6
	PAND(XMMWORD PTR [rsp+32], xmm6) //        pand    xmm6, XMMWORD PTR [rsp+32]
	PXOR(xmm11, xmm8) //        pxor    xmm8, xmm11
	MOVDQA(XMMWORD PTR [rsp], xmm11) //        movdqa  xmm11, XMMWORD PTR [rsp]
	MOVDQA(xmm7, xmm4) //        movdqa  xmm4, xmm7
	PSRLD(24, xmm4) //        psrld   xmm4, 24
	PSLLD(8, xmm7) //        pslld   xmm7, 8
	PAND(xmm8, xmm11) //        pand    xmm11, xmm8
	PAND(XMMWORD PTR [rsp+16], xmm8) //        pand    xmm8, XMMWORD PTR [rsp+16]
	PXOR(xmm4, xmm7) //        pxor    xmm7, xmm4
	MOVDQA(xmm11, xmm12) //        movdqa  xmm12, xmm11
	PSLLD(8, xmm11) //        pslld   xmm11, 8
	PXOR(xmm7, xmm6) //        pxor    xmm6, xmm7
	PSRLD(24, xmm12) //        psrld   xmm12, 24
	PSHUFD(108, xmm6, xmm6) //        pshufd  xmm6, xmm6, 108
	PXOR(xmm12, xmm11) //        pxor    xmm11, xmm12
	MOVDQA(xmm6, xmm4) //        movdqa  xmm4, xmm6
	PXOR(xmm11, xmm8) //        pxor    xmm8, xmm11
	MOVDQA(XMMWORD PTR [rsp+64], xmm11) //        movdqa  xmm11, XMMWORD PTR [rsp+64]
	PAND(xmm8, xmm11) //        pand    xmm11, xmm8
	PAND(XMMWORD PTR [rsp+48], xmm8) //        pand    xmm8, XMMWORD PTR [rsp+48]
	MOVDQA(xmm11, xmm12) //        movdqa  xmm12, xmm11
	PSLLD(8, xmm11) //        pslld   xmm11, 8
	PSRLD(24, xmm12) //        psrld   xmm12, 24
	PXOR(xmm12, xmm11) //        pxor    xmm11, xmm12
	PXOR(xmm11, xmm8) //        pxor    xmm8, xmm11
	PSHUFD(108, xmm8, xmm8) //        pshufd  xmm8, xmm8, 108
	MOVDQA(xmm8, xmm12) //        movdqa  xmm12, xmm8
	CMP(OFFSET FLAT:g_StepConstants+832, rax) //        cmp     rax, OFFSET FLAT:g_StepConstants+832
	JNE(.L12) //        jne     .L12
	PSHUFD(75, xmm9, xmm9) //        pshufd  xmm9, xmm9, 75
	PSHUFD(147, xmm10, xmm10) //        pshufd  xmm10, xmm10, 147
	PSHUFD(75, xmm4, XMMWORD PTR [rsp+160]) //        pshufd  xmm4, XMMWORD PTR [rsp+160], 75
	PSHUFD(147, xmm3, XMMWORD PTR [rsp+176]) //        pshufd  xmm3, XMMWORD PTR [rsp+176], 147
	PADDD(xmm15, xmm4) //        paddd   xmm4, xmm15
	PADDD(xmm0, xmm9) //        paddd   xmm9, xmm0
	SUB(-128, rbx) //        sub     rbx, -128
	PADDD(xmm14, xmm3) //        paddd   xmm3, xmm14
	PADDD(xmm13, xmm10) //        paddd   xmm10, xmm13
	PXOR(xmm2, xmm4) //        pxor    xmm4, xmm2
	PXOR(xmm1, xmm9) //        pxor    xmm9, xmm1
	PXOR(xmm8, xmm3) //        pxor    xmm3, xmm8
	PXOR(xmm6, xmm10) //        pxor    xmm10, xmm6
	CMP(rdx, rbx) //        cmp     rbx, rdx
	JNE(.L13) //        jne     .L13
	MOVDQA(xmm9, xmm0) //        movdqa  xmm0, xmm9
	MOVDQA(xmm10, xmm6) //        movdqa  xmm6, xmm10
	AND(127, r12d) //        and     r12d, 127
	*/

	Label("L11")                           //.L11:
	MOVUPS(xmm4, Mem{Base: rbp, Disp: 32}) //        movups  XMMWORD PTR [rbp+32], xmm4
	MOVUPS(xmm3, Mem{Base: rbp, Disp: 48}) //        movups  XMMWORD PTR [rbp+48], xmm3
	MOVUPS(xmm0, Mem{Base: rbp, Disp: 64}) //        movups  XMMWORD PTR [rbp+64], xmm0
	MOVUPS(xmm6, Mem{Base: rbp, Disp: 80}) //        movups  XMMWORD PTR [rbp+80], xmm6
	TESTQ(r12, r12)                        //        test    r12, r12
	JNE(LabelRef("L37"))                   //        jne     .L37
	ADDQ(I32(328), rsp)                    //        add     rsp, 328
	XORL(eax, eax)                         //        xor     eax, eax
	POPQ(rbx)                              //        pop     rbx
	POPQ(rbp)                              //        pop     rbp
	POPQ(r12)                              //        pop     r12
	POPQ(r13)                              //        pop     r13
	RET()                                  //        ret

	Label("L36") //.L36:
	/**
	MOV(128, r13d) //        mov     r13d, 128
	LEA([rdi+96+rdx], rdi) //        lea     rdi, [rdi+96+rdx]
	MOVAPS(xmm7, XMMWORD PTR [rsp+48]) //        movaps  XMMWORD PTR [rsp+48], xmm7
	SUB(eax, r13d) //        sub     r13d, eax
	MOVAPS(xmm2, XMMWORD PTR [rsp+32]) //        movaps  XMMWORD PTR [rsp+32], xmm2
	MOV(r13, rdx) //        mov     rdx, r13
	MOVAPS(xmm1, XMMWORD PTR [rsp+16]) //        movaps  XMMWORD PTR [rsp+16], xmm1
	MOVAPS(xmm5, XMMWORD PTR [rsp]) //        movaps  XMMWORD PTR [rsp], xmm5
	CALL(memcpy) //        call    memcpy
	MOVDQU(XMMWORD PTR [rbp+96], xmm3) //        movdqu  xmm3, XMMWORD PTR [rbp+96]
	MOVDQU(XMMWORD PTR [rbp+192], xmm2) //        movdqu  xmm2, XMMWORD PTR [rbp+192]
	MOV(OFFSET FLAT:g_StepConstants+64, eax) //        mov     eax, OFFSET FLAT:g_StepConstants+64
	MOVDQU(XMMWORD PTR [rbp+128], xmm5) //        movdqu  xmm5, XMMWORD PTR [rbp+128]
	MOVDQU(XMMWORD PTR [rbp+208], xmm7) //        movdqu  xmm7, XMMWORD PTR [rbp+208]
	MOVAPS(xmm2, XMMWORD PTR [rsp+96]) //        movaps  XMMWORD PTR [rsp+96], xmm2
	MOVDQA(XMMWORD PTR [rsp+32], xmm2) //        movdqa  xmm2, XMMWORD PTR [rsp+32]
	MOVDQU(XMMWORD PTR [rbp+112], xmm6) //        movdqu  xmm6, XMMWORD PTR [rbp+112]
	MOVAPS(xmm5, XMMWORD PTR [rsp+224]) //        movaps  XMMWORD PTR [rsp+224], xmm5
	MOVDQU(XMMWORD PTR [rbp+144], xmm1) //        movdqu  xmm1, XMMWORD PTR [rbp+144]
	MOVDQA(XMMWORD PTR .LC11[rip], xmm14) //        movdqa  xmm14, XMMWORD PTR .LC11[rip]
	PXOR(xmm5, xmm2) //        pxor    xmm2, xmm5
	MOVDQA(XMMWORD PTR [rsp], xmm5) //        movdqa  xmm5, XMMWORD PTR [rsp]
	MOVAPS(xmm7, XMMWORD PTR [rsp+112]) //        movaps  XMMWORD PTR [rsp+112], xmm7
	MOVDQA(XMMWORD PTR [rsp+48], xmm7) //        movdqa  xmm7, XMMWORD PTR [rsp+48]
	MOVDQA(XMMWORD PTR .LC14[rip], xmm15) //        movdqa  xmm15, XMMWORD PTR .LC14[rip]
	MOVAPS(xmm14, XMMWORD PTR [rsp+64]) //        movaps  XMMWORD PTR [rsp+64], xmm14
	MOVDQU(XMMWORD PTR [rbp+160], xmm13) //        movdqu  xmm13, XMMWORD PTR [rbp+160]
	PXOR(xmm3, xmm5) //        pxor    xmm5, xmm3
	MOVAPS(xmm3, XMMWORD PTR [rsp+192]) //        movaps  XMMWORD PTR [rsp+192], xmm3
	MOVDQA(xmm7, xmm0) //        movdqa  xmm0, xmm7
	MOVDQU(XMMWORD PTR [rbp+176], xmm12) //        movdqu  xmm12, XMMWORD PTR [rbp+176]
	PADDD(xmm2, xmm5) //        paddd   xmm5, xmm2
	PXOR(xmm1, xmm0) //        pxor    xmm0, xmm1
	MOVAPS(xmm1, XMMWORD PTR [rsp+240]) //        movaps  XMMWORD PTR [rsp+240], xmm1
	MOVDQA(XMMWORD PTR [rsp+16], xmm1) //        movdqa  xmm1, XMMWORD PTR [rsp+16]
	MOVDQA(xmm5, xmm3) //        movdqa  xmm3, xmm5
	PSLLD(29, xmm5) //        pslld   xmm5, 29
	MOVAPS(xmm6, XMMWORD PTR [rsp+208]) //        movaps  XMMWORD PTR [rsp+208], xmm6
	PSRLD(3, xmm3) //        psrld   xmm3, 3
	PXOR(xmm6, xmm1) //        pxor    xmm1, xmm6
	MOVAPS(xmm15, XMMWORD PTR [rsp+80]) //        movaps  XMMWORD PTR [rsp+80], xmm15
	MOVDQA(XMMWORD PTR .LC12[rip], xmm6) //        movdqa  xmm6, XMMWORD PTR .LC12[rip]
	POR(xmm3, xmm5) //        por     xmm5, xmm3
	MOVDQA(XMMWORD PTR .LC8[rip], xmm3) //        movdqa  xmm3, XMMWORD PTR .LC8[rip]
	PADDD(xmm0, xmm1) //        paddd   xmm1, xmm0
	MOVDQA(xmm1, xmm4) //        movdqa  xmm4, xmm1
	PSLLD(29, xmm1) //        pslld   xmm1, 29
	MOVAPS(xmm6, XMMWORD PTR [rsp+16]) //        movaps  XMMWORD PTR [rsp+16], xmm6
	PXOR(xmm5, xmm3) //        pxor    xmm3, xmm5
	PSRLD(3, xmm4) //        psrld   xmm4, 3
	MOVDQA(xmm3, xmm7) //        movdqa  xmm7, xmm3
	MOVDQA(XMMWORD PTR .LC9[rip], xmm3) //        movdqa  xmm3, XMMWORD PTR .LC9[rip]
	POR(xmm4, xmm1) //        por     xmm1, xmm4
	PADDD(xmm7, xmm2) //        paddd   xmm2, xmm7
	MOVDQA(xmm3, xmm8) //        movdqa  xmm8, xmm3
	PXOR(xmm1, xmm8) //        pxor    xmm8, xmm1
	MOVDQA(xmm2, xmm1) //        movdqa  xmm1, xmm2
	PSLLD(1, xmm2) //        pslld   xmm2, 1
	PSRLD(31, xmm1) //        psrld   xmm1, 31
	PADDD(xmm8, xmm0) //        paddd   xmm0, xmm8
	MOVDQA(xmm2, xmm3) //        movdqa  xmm3, xmm2
	POR(xmm1, xmm3) //        por     xmm3, xmm1
	MOVDQA(xmm0, xmm1) //        movdqa  xmm1, xmm0
	PSLLD(1, xmm0) //        pslld   xmm0, 1
	PSRLD(31, xmm1) //        psrld   xmm1, 31
	MOVDQA(xmm0, xmm5) //        movdqa  xmm5, xmm0
	POR(xmm1, xmm5) //        por     xmm5, xmm1
	MOVDQA(xmm14, xmm1) //        movdqa  xmm1, xmm14
	MOVDQA(xmm5, xmm2) //        movdqa  xmm2, xmm5
	MOVDQA(XMMWORD PTR .LC10[rip], xmm5) //        movdqa  xmm5, XMMWORD PTR .LC10[rip]
	PAND(xmm3, xmm1) //        pand    xmm1, xmm3
	MOVDQA(xmm5, xmm0) //        movdqa  xmm0, xmm5
	MOVAPS(xmm5, XMMWORD PTR [rsp+48]) //        movaps  XMMWORD PTR [rsp+48], xmm5
	MOVDQA(xmm6, xmm5) //        movdqa  xmm5, xmm6
	PAND(xmm3, xmm0) //        pand    xmm0, xmm3
	PADDD(xmm7, xmm3) //        paddd   xmm3, xmm7
	MOVDQA(xmm0, xmm4) //        movdqa  xmm4, xmm0
	PSLLD(8, xmm0) //        pslld   xmm0, 8
	PSHUFD(210, xmm3, xmm3) //        pshufd  xmm3, xmm3, 210
	PSRLD(24, xmm4) //        psrld   xmm4, 24
	PXOR(xmm4, xmm0) //        pxor    xmm0, xmm4
	PXOR(xmm0, xmm1) //        pxor    xmm1, xmm0
	PAND(xmm1, xmm5) //        pand    xmm5, xmm1
	MOVDQA(xmm5, xmm0) //        movdqa  xmm0, xmm5
	MOVDQA(xmm5, xmm4) //        movdqa  xmm4, xmm5
	MOVDQA(XMMWORD PTR .LC13[rip], xmm5) //        movdqa  xmm5, XMMWORD PTR .LC13[rip]
	PSRLD(24, xmm4) //        psrld   xmm4, 24
	PSLLD(8, xmm0) //        pslld   xmm0, 8
	PXOR(xmm4, xmm0) //        pxor    xmm0, xmm4
	PAND(xmm5, xmm1) //        pand    xmm1, xmm5
	MOVDQA(xmm2, xmm4) //        movdqa  xmm4, xmm2
	MOVAPS(xmm5, XMMWORD PTR [rsp]) //        movaps  XMMWORD PTR [rsp], xmm5
	PXOR(xmm0, xmm1) //        pxor    xmm1, xmm0
	MOVDQA(xmm15, xmm0) //        movdqa  xmm0, xmm15
	PAND(xmm1, xmm0) //        pand    xmm0, xmm1
	MOVDQA(xmm0, xmm9) //        movdqa  xmm9, xmm0
	MOVDQA(xmm0, xmm11) //        movdqa  xmm11, xmm0
	PSRLD(24, xmm9) //        psrld   xmm9, 24
	PSLLD(8, xmm11) //        pslld   xmm11, 8
	MOVDQA(xmm9, xmm10) //        movdqa  xmm10, xmm9
	MOVDQA(XMMWORD PTR .LC15[rip], xmm9) //        movdqa  xmm9, XMMWORD PTR .LC15[rip]
	PAND(xmm9, xmm4) //        pand    xmm4, xmm9
	PAND(xmm9, xmm1) //        pand    xmm1, xmm9
	MOVDQA(xmm4, xmm0) //        movdqa  xmm0, xmm4
	PSLLD(8, xmm4) //        pslld   xmm4, 8
	PSRLD(24, xmm0) //        psrld   xmm0, 24
	MOVDQA(xmm0, xmm5) //        movdqa  xmm5, xmm0
	MOVDQA(xmm15, xmm0) //        movdqa  xmm0, xmm15
	PXOR(xmm5, xmm4) //        pxor    xmm4, xmm5
	PAND(xmm2, xmm0) //        pand    xmm0, xmm2
	PADDD(xmm8, xmm2) //        paddd   xmm2, xmm8
	PXOR(xmm4, xmm0) //        pxor    xmm0, xmm4
	MOVDQA(XMMWORD PTR [rsp], xmm4) //        movdqa  xmm4, XMMWORD PTR [rsp]
	MOVDQA(XMMWORD PTR [rsp+48], xmm7) //        movdqa  xmm7, XMMWORD PTR [rsp+48]
	MOVDQA(xmm11, xmm8) //        movdqa  xmm8, xmm11
	PXOR(xmm10, xmm8) //        pxor    xmm8, xmm10
	PXOR(XMMWORD PTR [rsp+96], xmm3) //        pxor    xmm3, XMMWORD PTR [rsp+96]
	PSHUFD(210, xmm2, xmm2) //        pshufd  xmm2, xmm2, 210
	MOVAPS(xmm13, XMMWORD PTR [rsp+128]) //        movaps  XMMWORD PTR [rsp+128], xmm13
	PAND(xmm0, xmm4) //        pand    xmm4, xmm0
	PAND(xmm6, xmm0) //        pand    xmm0, xmm6
	PXOR(xmm8, xmm1) //        pxor    xmm1, xmm8
	MOVAPS(xmm12, XMMWORD PTR [rsp+144]) //        movaps  XMMWORD PTR [rsp+144], xmm12
	MOVDQA(xmm4, xmm5) //        movdqa  xmm5, xmm4
	PSLLD(8, xmm4) //        pslld   xmm4, 8
	PSHUFD(108, xmm1, xmm1) //        pshufd  xmm1, xmm1, 108
	MOVAPS(xmm9, XMMWORD PTR [rsp+32]) //        movaps  XMMWORD PTR [rsp+32], xmm9
	PSRLD(24, xmm5) //        psrld   xmm5, 24
	PXOR(XMMWORD PTR [rsp+112], xmm1) //        pxor    xmm1, XMMWORD PTR [rsp+112]
	PXOR(xmm13, xmm2) //        pxor    xmm2, xmm13
	PXOR(xmm5, xmm4) //        pxor    xmm4, xmm5
	PADDD(xmm3, xmm2) //        paddd   xmm2, xmm3
	PXOR(xmm4, xmm0) //        pxor    xmm0, xmm4
	MOVDQA(xmm14, xmm4) //        movdqa  xmm4, xmm14
	PAND(xmm0, xmm4) //        pand    xmm4, xmm0
	PAND(xmm7, xmm0) //        pand    xmm0, xmm7
	MOVDQA(xmm4, xmm6) //        movdqa  xmm6, xmm4
	PSLLD(8, xmm4) //        pslld   xmm4, 8
	PSRLD(24, xmm6) //        psrld   xmm6, 24
	PXOR(xmm6, xmm4) //        pxor    xmm4, xmm6
	PXOR(xmm4, xmm0) //        pxor    xmm0, xmm4
	MOVDQA(xmm2, xmm4) //        movdqa  xmm4, xmm2
	PSRLD(27, xmm4) //        psrld   xmm4, 27
	PSLLD(5, xmm2) //        pslld   xmm2, 5
	PSHUFD(108, xmm0, xmm0) //        pshufd  xmm0, xmm0, 108
	PXOR(xmm12, xmm0) //        pxor    xmm0, xmm12
	POR(xmm4, xmm2) //        por     xmm2, xmm4
	PADDD(xmm1, xmm0) //        paddd   xmm0, xmm1
	MOVDQA(xmm0, xmm10) //        movdqa  xmm10, xmm0
	PSLLD(5, xmm0) //        pslld   xmm0, 5
	PSRLD(27, xmm10) //        psrld   xmm10, 27
	MOVDQA(xmm10, xmm5) //        movdqa  xmm5, xmm10
	MOVDQA(XMMWORD PTR .LC16[rip], xmm10) //        movdqa  xmm10, XMMWORD PTR .LC16[rip]
	POR(xmm5, xmm0) //        por     xmm0, xmm5
	MOVDQA(XMMWORD PTR [rsp+16], xmm5) //        movdqa  xmm5, XMMWORD PTR [rsp+16]
	PXOR(xmm2, xmm10) //        pxor    xmm10, xmm2
	MOVDQA(XMMWORD PTR .LC17[rip], xmm2) //        movdqa  xmm2, XMMWORD PTR .LC17[rip]
	PADDD(xmm10, xmm3) //        paddd   xmm3, xmm10
	MOVDQA(xmm10, xmm4) //        movdqa  xmm4, xmm10
	MOVDQA(xmm2, xmm8) //        movdqa  xmm8, xmm2
	PXOR(xmm0, xmm8) //        pxor    xmm8, xmm0
	MOVDQA(xmm3, xmm0) //        movdqa  xmm0, xmm3
	PSRLD(15, xmm0) //        psrld   xmm0, 15
	PSLLD(17, xmm3) //        pslld   xmm3, 17
	PADDD(xmm8, xmm1) //        paddd   xmm1, xmm8
	POR(xmm0, xmm3) //        por     xmm3, xmm0
	MOVDQA(xmm1, xmm0) //        movdqa  xmm0, xmm1
	PSRLD(15, xmm0) //        psrld   xmm0, 15
	PSLLD(17, xmm1) //        pslld   xmm1, 17
	PADDD(xmm3, xmm4) //        paddd   xmm4, xmm3
	POR(xmm0, xmm1) //        por     xmm1, xmm0
	MOVDQA(xmm7, xmm0) //        movdqa  xmm0, xmm7
	MOVDQA(xmm14, xmm7) //        movdqa  xmm7, xmm14
	PAND(xmm3, xmm0) //        pand    xmm0, xmm3
	MOVDQA(xmm1, xmm6) //        movdqa  xmm6, xmm1
	PAND(xmm3, xmm7) //        pand    xmm7, xmm3
	MOVDQA(xmm0, xmm1) //        movdqa  xmm1, xmm0
	PSLLD(8, xmm0) //        pslld   xmm0, 8
	MOVDQA(xmm8, xmm3) //        movdqa  xmm3, xmm8
	PSRLD(24, xmm1) //        psrld   xmm1, 24
	PADDD(xmm6, xmm3) //        paddd   xmm3, xmm6
	PSHUFD(210, xmm4, xmm4) //        pshufd  xmm4, xmm4, 210
	PXOR(xmm1, xmm0) //        pxor    xmm0, xmm1
	MOVDQA(xmm5, xmm1) //        movdqa  xmm1, xmm5
	PXOR(xmm0, xmm7) //        pxor    xmm7, xmm0
	PAND(xmm7, xmm1) //        pand    xmm1, xmm7
	PAND(XMMWORD PTR [rsp], xmm7) //        pand    xmm7, XMMWORD PTR [rsp]
	MOVDQA(xmm1, xmm0) //        movdqa  xmm0, xmm1
	PSRLD(24, xmm1) //        psrld   xmm1, 24
	PSLLD(8, xmm0) //        pslld   xmm0, 8
	PXOR(xmm1, xmm0) //        pxor    xmm0, xmm1
	MOVDQA(xmm6, xmm1) //        movdqa  xmm1, xmm6
	PXOR(xmm0, xmm7) //        pxor    xmm7, xmm0
	MOVDQA(xmm15, xmm0) //        movdqa  xmm0, xmm15
	PAND(xmm9, xmm1) //        pand    xmm1, xmm9
	PAND(xmm7, xmm0) //        pand    xmm0, xmm7
	PAND(xmm6, xmm15) //        pand    xmm15, xmm6
	PAND(xmm9, xmm7) //        pand    xmm7, xmm9
	MOVDQA(xmm0, xmm10) //        movdqa  xmm10, xmm0
	MOVDQA(xmm0, xmm11) //        movdqa  xmm11, xmm0
	MOVDQA(xmm1, xmm0) //        movdqa  xmm0, xmm1
	PSRLD(24, xmm1) //        psrld   xmm1, 24
	PSLLD(8, xmm0) //        pslld   xmm0, 8
	PSLLD(8, xmm11) //        pslld   xmm11, 8
	PSRLD(24, xmm10) //        psrld   xmm10, 24
	MOVDQA(xmm1, xmm2) //        movdqa  xmm2, xmm1
	PXOR(xmm2, xmm0) //        pxor    xmm0, xmm2
	MOVDQA(xmm15, xmm1) //        movdqa  xmm1, xmm15
	MOVDQA(XMMWORD PTR [rsp], xmm2) //        movdqa  xmm2, XMMWORD PTR [rsp]
	PXOR(xmm0, xmm1) //        pxor    xmm1, xmm0
	PAND(xmm1, xmm2) //        pand    xmm2, xmm1
	PAND(xmm5, xmm1) //        pand    xmm1, xmm5
	MOVDQA(xmm2, xmm0) //        movdqa  xmm0, xmm2
	PSRLD(24, xmm2) //        psrld   xmm2, 24
	PSLLD(8, xmm0) //        pslld   xmm0, 8
	PXOR(xmm2, xmm0) //        pxor    xmm0, xmm2
	PXOR(xmm0, xmm1) //        pxor    xmm1, xmm0
	PSHUFD(210, xmm0, xmm3) //        pshufd  xmm0, xmm3, 210
	MOVDQA(xmm11, xmm3) //        movdqa  xmm3, xmm11
	PAND(xmm1, xmm14) //        pand    xmm14, xmm1
	PAND(XMMWORD PTR [rsp+48], xmm1) //        pand    xmm1, XMMWORD PTR [rsp+48]
	PXOR(xmm10, xmm3) //        pxor    xmm3, xmm10
	MOVDQA(xmm14, xmm2) //        movdqa  xmm2, xmm14
	MOVDQA(xmm14, xmm5) //        movdqa  xmm5, xmm14
	PXOR(xmm3, xmm7) //        pxor    xmm7, xmm3
	PSRLD(24, xmm5) //        psrld   xmm5, 24
	PSLLD(8, xmm2) //        pslld   xmm2, 8
	PSHUFD(108, xmm7, xmm7) //        pshufd  xmm7, xmm7, 108
	PXOR(xmm5, xmm2) //        pxor    xmm2, xmm5
	PXOR(xmm2, xmm1) //        pxor    xmm1, xmm2
	MOVDQA(xmm4, xmm2) //        movdqa  xmm2, xmm4
	PSHUFD(108, xmm3, xmm1) //        pshufd  xmm3, xmm1, 108
	MOVDQA(xmm3, xmm4) //        movdqa  xmm4, xmm3
	*/
	Label("L10") //.L10:
	/**
	MOVDQA(xmm0, xmm5) //        movdqa  xmm5, xmm0
	ADD(64, rax) //        add     rax, 64
	PSHUFD(75, xmm1, XMMWORD PTR [rsp+192]) //        pshufd  xmm1, XMMWORD PTR [rsp+192], 75
	PADDD(XMMWORD PTR [rsp+128], xmm1) //        paddd   xmm1, XMMWORD PTR [rsp+128]
	MOVDQA(XMMWORD PTR [rsp+48], xmm13) //        movdqa  xmm13, XMMWORD PTR [rsp+48]
	MOVDQA(XMMWORD PTR [rsp+64], xmm15) //        movdqa  xmm15, XMMWORD PTR [rsp+64]
	PSHUFD(75, xmm11, XMMWORD PTR [rsp+224]) //        pshufd  xmm11, XMMWORD PTR [rsp+224], 75
	PSHUFD(147, xmm10, XMMWORD PTR [rsp+240]) //        pshufd  xmm10, XMMWORD PTR [rsp+240], 147
	PADDD(XMMWORD PTR [rsp+96], xmm11) //        paddd   xmm11, XMMWORD PTR [rsp+96]
	MOVDQA(xmm1, xmm3) //        movdqa  xmm3, xmm1
	MOVAPS(xmm1, XMMWORD PTR [rsp+192]) //        movaps  XMMWORD PTR [rsp+192], xmm1
	PSHUFD(147, xmm1, XMMWORD PTR [rsp+208]) //        pshufd  xmm1, XMMWORD PTR [rsp+208], 147
	PADDD(XMMWORD PTR [rsp+144], xmm1) //        paddd   xmm1, XMMWORD PTR [rsp+144]
	MOVAPS(xmm3, XMMWORD PTR [rsp+160]) //        movaps  XMMWORD PTR [rsp+160], xmm3
	PXOR(XMMWORD PTR [rsp+160], xmm5) //        pxor    xmm5, XMMWORD PTR [rsp+160]
	PADDD(XMMWORD PTR [rsp+112], xmm10) //        paddd   xmm10, XMMWORD PTR [rsp+112]
	PXOR(xmm11, xmm2) //        pxor    xmm2, xmm11
	MOVDQA(XMMWORD PTR [rsp+16], xmm14) //        movdqa  xmm14, XMMWORD PTR [rsp+16]
	MOVAPS(xmm11, XMMWORD PTR [rsp+224]) //        movaps  XMMWORD PTR [rsp+224], xmm11
	MOVDQA(xmm1, xmm9) //        movdqa  xmm9, xmm1
	PADDD(xmm2, xmm5) //        paddd   xmm5, xmm2
	MOVAPS(xmm1, XMMWORD PTR [rsp+208]) //        movaps  XMMWORD PTR [rsp+208], xmm1
	MOVDQA(xmm4, xmm1) //        movdqa  xmm1, xmm4
	PXOR(xmm10, xmm7) //        pxor    xmm7, xmm10
	MOVDQA(xmm5, xmm6) //        movdqa  xmm6, xmm5
	MOVDQA(XMMWORD PTR [rsp], xmm12) //        movdqa  xmm12, XMMWORD PTR [rsp]
	MOVAPS(xmm9, XMMWORD PTR [rsp+176]) //        movaps  XMMWORD PTR [rsp+176], xmm9
	PXOR(XMMWORD PTR [rsp+176], xmm1) //        pxor    xmm1, XMMWORD PTR [rsp+176]
	PSRLD(3, xmm6) //        psrld   xmm6, 3
	MOVAPS(xmm10, XMMWORD PTR [rsp+240]) //        movaps  XMMWORD PTR [rsp+240], xmm10
	PSLLD(29, xmm5) //        pslld   xmm5, 29
	PADDD(xmm7, xmm1) //        paddd   xmm1, xmm7
	POR(xmm6, xmm5) //        por     xmm5, xmm6
	PXOR(XMMWORD PTR [rax-64], xmm5) //        pxor    xmm5, XMMWORD PTR [rax-64]
	MOVDQA(xmm1, xmm3) //        movdqa  xmm3, xmm1
	PSLLD(29, xmm1) //        pslld   xmm1, 29
	PSRLD(3, xmm3) //        psrld   xmm3, 3
	PADDD(xmm5, xmm2) //        paddd   xmm2, xmm5
	POR(xmm3, xmm1) //        por     xmm1, xmm3
	PXOR(XMMWORD PTR [rax-48], xmm1) //        pxor    xmm1, XMMWORD PTR [rax-48]
	MOVDQA(xmm2, xmm0) //        movdqa  xmm0, xmm2
	PSRLD(31, xmm0) //        psrld   xmm0, 31
	PSLLD(1, xmm2) //        pslld   xmm2, 1
	MOVDQA(xmm15, xmm3) //        movdqa  xmm3, xmm15
	PADDD(xmm1, xmm7) //        paddd   xmm7, xmm1
	POR(xmm0, xmm2) //        por     xmm2, xmm0
	MOVDQA(xmm7, xmm0) //        movdqa  xmm0, xmm7
	PSLLD(1, xmm7) //        pslld   xmm7, 1
	PAND(xmm2, xmm3) //        pand    xmm3, xmm2
	PSRLD(31, xmm0) //        psrld   xmm0, 31
	PADDD(xmm2, xmm5) //        paddd   xmm5, xmm2
	POR(xmm0, xmm7) //        por     xmm7, xmm0
	MOVDQA(xmm13, xmm0) //        movdqa  xmm0, xmm13
	PAND(xmm2, xmm0) //        pand    xmm0, xmm2
	PADDD(xmm7, xmm1) //        paddd   xmm1, xmm7
	PSHUFD(210, xmm2, xmm5) //        pshufd  xmm2, xmm5, 210
	MOVDQA(xmm0, xmm4) //        movdqa  xmm4, xmm0
	PSLLD(8, xmm0) //        pslld   xmm0, 8
	PSHUFD(210, xmm1, xmm1) //        pshufd  xmm1, xmm1, 210
	PSRLD(24, xmm4) //        psrld   xmm4, 24
	PXOR(xmm4, xmm0) //        pxor    xmm0, xmm4
	PXOR(xmm0, xmm3) //        pxor    xmm3, xmm0
	MOVDQA(xmm14, xmm0) //        movdqa  xmm0, xmm14
	PAND(xmm3, xmm0) //        pand    xmm0, xmm3
	PAND(xmm12, xmm3) //        pand    xmm3, xmm12
	MOVDQA(xmm0, xmm4) //        movdqa  xmm4, xmm0
	PSLLD(8, xmm0) //        pslld   xmm0, 8
	PSRLD(24, xmm4) //        psrld   xmm4, 24
	PXOR(xmm4, xmm0) //        pxor    xmm0, xmm4
	PXOR(xmm0, xmm3) //        pxor    xmm3, xmm0
	MOVDQA(XMMWORD PTR [rsp+80], xmm0) //        movdqa  xmm0, XMMWORD PTR [rsp+80]
	MOVDQA(xmm0, xmm6) //        movdqa  xmm6, xmm0
	PAND(xmm7, xmm0) //        pand    xmm0, xmm7
	PAND(xmm3, xmm6) //        pand    xmm6, xmm3
	PAND(XMMWORD PTR [rsp+32], xmm3) //        pand    xmm3, XMMWORD PTR [rsp+32]
	MOVDQA(xmm6, xmm4) //        movdqa  xmm4, xmm6
	PSLLD(8, xmm6) //        pslld   xmm6, 8
	PSRLD(24, xmm4) //        psrld   xmm4, 24
	MOVDQA(xmm4, xmm9) //        movdqa  xmm9, xmm4
	MOVDQA(XMMWORD PTR [rsp+32], xmm4) //        movdqa  xmm4, XMMWORD PTR [rsp+32]
	PXOR(xmm9, xmm6) //        pxor    xmm6, xmm9
	PSHUFD(75, xmm9, XMMWORD PTR [rsp+128]) //        pshufd  xmm9, XMMWORD PTR [rsp+128], 75
	PADDD(XMMWORD PTR [rsp+160], xmm9) //        paddd   xmm9, XMMWORD PTR [rsp+160]
	PAND(xmm7, xmm4) //        pand    xmm4, xmm7
	MOVDQA(xmm13, xmm7) //        movdqa  xmm7, xmm13
	PXOR(xmm6, xmm3) //        pxor    xmm3, xmm6
	MOVDQA(xmm4, xmm8) //        movdqa  xmm8, xmm4
	PSLLD(8, xmm4) //        pslld   xmm4, 8
	PXOR(xmm9, xmm1) //        pxor    xmm1, xmm9
	MOVAPS(xmm9, XMMWORD PTR [rsp+128]) //        movaps  XMMWORD PTR [rsp+128], xmm9
	PSRLD(24, xmm8) //        psrld   xmm8, 24
	PSHUFD(108, xmm3, xmm3) //        pshufd  xmm3, xmm3, 108
	PXOR(xmm8, xmm4) //        pxor    xmm4, xmm8
	PXOR(xmm4, xmm0) //        pxor    xmm0, xmm4
	MOVDQA(xmm12, xmm4) //        movdqa  xmm4, xmm12
	PSHUFD(147, xmm12, XMMWORD PTR [rsp+144]) //        pshufd  xmm12, XMMWORD PTR [rsp+144], 147
	PADDD(XMMWORD PTR [rsp+176], xmm12) //        paddd   xmm12, XMMWORD PTR [rsp+176]
	PAND(xmm0, xmm4) //        pand    xmm4, xmm0
	PAND(xmm14, xmm0) //        pand    xmm0, xmm14
	MOVAPS(xmm12, XMMWORD PTR [rsp+144]) //        movaps  XMMWORD PTR [rsp+144], xmm12
	MOVDQA(xmm4, xmm8) //        movdqa  xmm8, xmm4
	PSLLD(8, xmm4) //        pslld   xmm4, 8
	PSRLD(24, xmm8) //        psrld   xmm8, 24
	PXOR(xmm8, xmm4) //        pxor    xmm4, xmm8
	PXOR(xmm4, xmm0) //        pxor    xmm0, xmm4
	MOVDQA(xmm15, xmm4) //        movdqa  xmm4, xmm15
	PAND(xmm0, xmm4) //        pand    xmm4, xmm0
	PAND(xmm13, xmm0) //        pand    xmm0, xmm13
	PSHUFD(75, xmm13, XMMWORD PTR [rsp+96]) //        pshufd  xmm13, XMMWORD PTR [rsp+96], 75
	PADDD(xmm11, xmm13) //        paddd   xmm13, xmm11
	MOVDQA(xmm4, xmm8) //        movdqa  xmm8, xmm4
	PSLLD(8, xmm4) //        pslld   xmm4, 8
	PXOR(xmm13, xmm2) //        pxor    xmm2, xmm13
	MOVAPS(xmm13, XMMWORD PTR [rsp+96]) //        movaps  XMMWORD PTR [rsp+96], xmm13
	PSRLD(24, xmm8) //        psrld   xmm8, 24
	PADDD(xmm2, xmm1) //        paddd   xmm1, xmm2
	PSHUFD(147, xmm14, XMMWORD PTR [rsp+112]) //        pshufd  xmm14, XMMWORD PTR [rsp+112], 147
	PADDD(xmm10, xmm14) //        paddd   xmm14, xmm10
	PXOR(xmm8, xmm4) //        pxor    xmm4, xmm8
	PXOR(xmm14, xmm3) //        pxor    xmm3, xmm14
	MOVDQA(xmm1, xmm5) //        movdqa  xmm5, xmm1
	MOVDQA(XMMWORD PTR [rsp+16], xmm8) //        movdqa  xmm8, XMMWORD PTR [rsp+16]
	PXOR(xmm4, xmm0) //        pxor    xmm0, xmm4
	PSRLD(27, xmm5) //        psrld   xmm5, 27
	MOVAPS(xmm14, XMMWORD PTR [rsp+112]) //        movaps  XMMWORD PTR [rsp+112], xmm14
	PSHUFD(108, xmm0, xmm0) //        pshufd  xmm0, xmm0, 108
	PSLLD(5, xmm1) //        pslld   xmm1, 5
	PXOR(xmm12, xmm0) //        pxor    xmm0, xmm12
	POR(xmm5, xmm1) //        por     xmm1, xmm5
	PXOR(XMMWORD PTR [rax-32], xmm1) //        pxor    xmm1, XMMWORD PTR [rax-32]
	PADDD(xmm3, xmm0) //        paddd   xmm0, xmm3
	MOVDQA(xmm15, xmm5) //        movdqa  xmm5, xmm15
	MOVDQA(xmm0, xmm4) //        movdqa  xmm4, xmm0
	PSLLD(5, xmm0) //        pslld   xmm0, 5
	PADDD(xmm1, xmm2) //        paddd   xmm2, xmm1
	PSRLD(27, xmm4) //        psrld   xmm4, 27
	POR(xmm4, xmm0) //        por     xmm0, xmm4
	PXOR(XMMWORD PTR [rax-16], xmm0) //        pxor    xmm0, XMMWORD PTR [rax-16]
	MOVDQA(xmm2, xmm4) //        movdqa  xmm4, xmm2
	PSRLD(15, xmm4) //        psrld   xmm4, 15
	PSLLD(17, xmm2) //        pslld   xmm2, 17
	PADDD(xmm0, xmm3) //        paddd   xmm3, xmm0
	POR(xmm4, xmm2) //        por     xmm2, xmm4
	MOVDQA(xmm3, xmm4) //        movdqa  xmm4, xmm3
	PSLLD(17, xmm3) //        pslld   xmm3, 17
	PAND(xmm2, xmm5) //        pand    xmm5, xmm2
	PSRLD(15, xmm4) //        psrld   xmm4, 15
	PADDD(xmm2, xmm1) //        paddd   xmm1, xmm2
	POR(xmm4, xmm3) //        por     xmm3, xmm4
	MOVDQA(xmm7, xmm4) //        movdqa  xmm4, xmm7
	PSHUFD(210, xmm1, xmm1) //        pshufd  xmm1, xmm1, 210
	PAND(xmm2, xmm4) //        pand    xmm4, xmm2
	PADDD(xmm3, xmm0) //        paddd   xmm0, xmm3
	MOVDQA(xmm1, xmm2) //        movdqa  xmm2, xmm1
	MOVDQA(xmm4, xmm6) //        movdqa  xmm6, xmm4
	PSLLD(8, xmm4) //        pslld   xmm4, 8
	PSRLD(24, xmm6) //        psrld   xmm6, 24
	PXOR(xmm6, xmm4) //        pxor    xmm4, xmm6
	PXOR(xmm4, xmm5) //        pxor    xmm5, xmm4
	PAND(xmm5, xmm8) //        pand    xmm8, xmm5
	PAND(XMMWORD PTR [rsp], xmm5) //        pand    xmm5, XMMWORD PTR [rsp]
	MOVDQA(xmm8, xmm4) //        movdqa  xmm4, xmm8
	MOVDQA(xmm8, xmm6) //        movdqa  xmm6, xmm8
	PSRLD(24, xmm6) //        psrld   xmm6, 24
	PSLLD(8, xmm4) //        pslld   xmm4, 8
	PXOR(xmm6, xmm4) //        pxor    xmm4, xmm6
	MOVDQA(XMMWORD PTR [rsp+80], xmm6) //        movdqa  xmm6, XMMWORD PTR [rsp+80]
	PXOR(xmm4, xmm5) //        pxor    xmm5, xmm4
	MOVDQA(xmm6, xmm7) //        movdqa  xmm7, xmm6
	PAND(xmm3, xmm6) //        pand    xmm6, xmm3
	PAND(xmm5, xmm7) //        pand    xmm7, xmm5
	PAND(XMMWORD PTR [rsp+32], xmm5) //        pand    xmm5, XMMWORD PTR [rsp+32]
	MOVDQA(xmm7, xmm15) //        movdqa  xmm15, xmm7
	PSLLD(8, xmm7) //        pslld   xmm7, 8
	PSRLD(24, xmm15) //        psrld   xmm15, 24
	MOVDQA(xmm15, xmm4) //        movdqa  xmm4, xmm15
	MOVDQA(XMMWORD PTR [rsp+32], xmm15) //        movdqa  xmm15, XMMWORD PTR [rsp+32]
	PXOR(xmm4, xmm7) //        pxor    xmm7, xmm4
	PAND(xmm3, xmm15) //        pand    xmm15, xmm3
	PXOR(xmm7, xmm5) //        pxor    xmm5, xmm7
	PSHUFD(210, xmm3, xmm0) //        pshufd  xmm3, xmm0, 210
	MOVDQA(xmm15, xmm8) //        movdqa  xmm8, xmm15
	PSRLD(24, xmm15) //        psrld   xmm15, 24
	PSHUFD(108, xmm5, xmm5) //        pshufd  xmm5, xmm5, 108
	PSLLD(8, xmm8) //        pslld   xmm8, 8
	MOVDQA(xmm3, xmm0) //        movdqa  xmm0, xmm3
	MOVDQA(xmm5, xmm7) //        movdqa  xmm7, xmm5
	PXOR(xmm15, xmm8) //        pxor    xmm8, xmm15
	PXOR(xmm8, xmm6) //        pxor    xmm6, xmm8
	MOVDQA(XMMWORD PTR [rsp], xmm8) //        movdqa  xmm8, XMMWORD PTR [rsp]
	PAND(xmm6, xmm8) //        pand    xmm8, xmm6
	PAND(XMMWORD PTR [rsp+16], xmm6) //        pand    xmm6, XMMWORD PTR [rsp+16]
	MOVDQA(xmm8, xmm15) //        movdqa  xmm15, xmm8
	PSLLD(8, xmm8) //        pslld   xmm8, 8
	PSRLD(24, xmm15) //        psrld   xmm15, 24
	PXOR(xmm15, xmm8) //        pxor    xmm8, xmm15
	PXOR(xmm8, xmm6) //        pxor    xmm6, xmm8
	MOVDQA(XMMWORD PTR [rsp+64], xmm8) //        movdqa  xmm8, XMMWORD PTR [rsp+64]
	PAND(xmm6, xmm8) //        pand    xmm8, xmm6
	PAND(XMMWORD PTR [rsp+48], xmm6) //        pand    xmm6, XMMWORD PTR [rsp+48]
	MOVDQA(xmm8, xmm15) //        movdqa  xmm15, xmm8
	PSLLD(8, xmm8) //        pslld   xmm8, 8
	PSRLD(24, xmm15) //        psrld   xmm15, 24
	PXOR(xmm15, xmm8) //        pxor    xmm8, xmm15
	PXOR(xmm8, xmm6) //        pxor    xmm6, xmm8
	PSHUFD(108, xmm8, xmm6) //        pshufd  xmm8, xmm6, 108
	MOVDQA(xmm8, xmm4) //        movdqa  xmm4, xmm8
	CMP(OFFSET FLAT:g_StepConstants+832, rax) //        cmp     rax, OFFSET FLAT:g_StepConstants+832
	JNE(.L10) //        jne     .L10
	MOVDQA(xmm14, xmm7) //        movdqa  xmm7, xmm14
	MOVDQA(xmm3, xmm2) //        movdqa  xmm2, xmm3
	ADD(r13, rbx) //        add     rbx, r13
	SUB(r13, r12) //        sub     r12, r13
	MOVDQA(XMMWORD PTR [rsp+176], xmm15) //        movdqa  xmm15, XMMWORD PTR [rsp+176]
	PSHUFD(75, xmm0, xmm11) //        pshufd  xmm0, xmm11, 75
	MOVDQA(XMMWORD PTR [rsp+160], xmm14) //        movdqa  xmm14, XMMWORD PTR [rsp+160]
	PSHUFD(147, xmm6, xmm10) //        pshufd  xmm6, xmm10, 147
	PADDD(xmm13, xmm0) //        paddd   xmm0, xmm13
	PADDD(xmm7, xmm6) //        paddd   xmm6, xmm7
	MOV(0, DWORD PTR [rbp+16]) //        mov     DWORD PTR [rbp+16], 0
	PSHUFD(75, xmm4, xmm14) //        pshufd  xmm4, xmm14, 75
	PSHUFD(147, xmm3, xmm15) //        pshufd  xmm3, xmm15, 147
	PXOR(xmm1, xmm0) //        pxor    xmm0, xmm1
	PADDD(xmm9, xmm4) //        paddd   xmm4, xmm9
	PADDD(xmm12, xmm3) //        paddd   xmm3, xmm12
	PXOR(xmm5, xmm6) //        pxor    xmm6, xmm5
	PXOR(xmm2, xmm4) //        pxor    xmm4, xmm2
	PXOR(xmm8, xmm3) //        pxor    xmm3, xmm8
	JMP(.L9) //        jmp     .L9
	*/
	Label("L37") //.L37:
	/**
	LEA([rbp+96], rcx) //        lea     rcx, [rbp+96]
	MOV(r12d, eax) //        mov     eax, r12d
	CMP(8, r12d) //        cmp     r12d, 8
	JNB(.L14) //        jnb     .L14
	TEST(4, r12b) //        test    r12b, 4
	JNE(.L38) //        jne     .L38
	TEST(eax, eax) //        test    eax, eax
	JE(.L15) //        je      .L15
	MOVZX(BYTE PTR [rdx], esi) //        movzx   esi, BYTE PTR [rdx]
	MOV(sil, BYTE PTR [rbp+96]) //        mov     BYTE PTR [rbp+96], sil
	TEST(2, al) //        test    al, 2
	JNE(.L39) //        jne     .L39
	*/

	Label("L15")                         //.L15:
	SALL(U8(3), r12d)                    //        sal     r12d, 3
	XORL(eax, eax)                       //        xor     eax, eax
	MOVL(r12d, Mem{Base: rbp, Disp: 16}) //        mov     DWORD PTR [rbp+16], r12d
	ADDQ(U32(328), rsp)                  //        add     rsp, 328
	POPQ(rbx)                            //        pop     rbx
	POPQ(rbp)                            //        pop     rbp
	POPQ(r12)                            //        pop     r12
	POPQ(r13)                            //        pop     r13
	RET()                                //        ret

	Label("L35")                                              //.L35:
	LEAQ(Mem{Base: rdi, Disp: 96, Index: rdx, Scale: 1}, rdi) //        lea     rdi, [rdi+96+rdx]
	MOVQ(r12, rdx)                                            //        mov     rdx, r12
	//CALL(memcpy)                                             //        call    memcpy
	ADDQ(r13d, Mem{Base: rbp, Disp: 16}) //        add     DWORD PTR [rbp+16], r13d
	ADDL(I32(328), rsp)                  //        add     rsp, 328
	XORL(eax, eax)                       //        xor     eax, eax
	POPQ(rbx)                            //        pop     rbx
	POPQ(rbp)                            //        pop     rbp
	POPQ(r12)                            //        pop     r12
	POPQ(r13)                            //        pop     r13
	RET()                                //        ret

	Label("L14")                                              //.L14:
	MOVQ(Mem{Base: rdx}, rax)                                 //        mov     rax, QWORD PTR [rdx]
	LEAQ(Mem{Base: rbp, Disp: 104}, rdi)                      //        lea     rdi, [rbp+104]
	ANDQ(I8(-8), rdi)                                         //        and     rdi, -8
	MOVQ(rax, Mem{Base: rbp, Disp: 96})                       //        mov     QWORD PTR [rbp+96], rax
	MOVL(r12d, eax)                                           //        mov     eax, r12d
	MOVQ(Mem{Base: rdx, Disp: -8, Index: rax, Scale: 1}, rsi) //        mov     rsi, QWORD PTR [rdx-8+rax]
	MOVQ(rsi, Mem{Base: rcx, Disp: -8, Index: rax, Scale: 1}) //        mov     QWORD PTR [rcx-8+rax], rsi
	SUBQ(rdi, rcx)                                            //        sub     rcx, rdi
	MOVQ(rdx, rsi)                                            //        mov     rsi, rdx
	LEAL(Mem{Base: r12, Index: rcx, Scale: 1}, eax)           //        lea     eax, [r12+rcx]
	SUBQ(rcx, rsi)                                            //        sub     rsi, rcx
	SHRL(I8(3), eax)                                          //        shr     eax, 3
	MOVL(eax, ecx)                                            //        mov     ecx, eax
	//REPQ(movsq)                                               //        rep movsq
	JMP(LabelRef("L15")) //        jmp     .L15

	Label("L18")         //.L18:
	MOVQ(rbx, rdx)       //        mov     rdx, rbx
	JMP(LabelRef("L11")) //        jmp     .L11

	Label("L38")                                              //.L38:
	MOVL(Mem{Base: rdx}, esi)                                 //        mov     esi, DWORD PTR [rdx]
	MOVL(esi, Mem{Base: rbp, Disp: 96})                       //        mov     DWORD PTR [rbp+96], esi
	MOVL(Mem{Base: rdx, Index: rax, Scale: 1, Disp: -4}, edx) //        mov     edx, DWORD PTR [rdx-4+rax]
	MOVL(edx, Mem{Base: rcx, Index: rax, Scale: 1, Disp: -4}) //        mov     DWORD PTR [rcx-4+rax], edx
	JMP(LabelRef("L15"))                                      //        jmp     .L15

	Label("L39")                                                 //.L39:
	MOVWLZX(Mem{Base: rdx, Index: rax, Scale: 1, Disp: -2}, edx) //        movzx   edx, WORD PTR [rdx-2+rax]
	MOVW(dx, Mem{Base: rcx, Index: rax, Scale: 1, Disp: -2})     //        mov     WORD PTR [rcx-2+rax], dx
	JMP(LabelRef("L15"))                                         //        jmp     .L15
}

//lsh256_sse2_final(LSH256_Context*, unsigned char*):
//        push    rbp
//        mov     rbp, rsi
//        xor     esi, esi
//        push    rbx
//        mov     rbx, rdi
//        sub     rsp, 168
//        mov     eax, DWORD PTR [rdi+16]
//        shr     eax, 3
//        mov     edx, eax
//        mov     BYTE PTR [rdi+96+rdx], -128
//        mov     edx, 127
//        sub     edx, eax
//        add     eax, 1
//        lea     rdi, [rdi+96+rax]
//        call    memset
//        movdqu  xmm2, XMMWORD PTR [rbx+64]
//        movdqu  xmm7, XMMWORD PTR [rbx+96]
//        mov     eax, OFFSET FLAT:g_StepConstants+64
//        movdqu  xmm3, XMMWORD PTR [rbx+128]
//        movdqu  xmm4, XMMWORD PTR [rbx+112]
//        mov     edx, OFFSET FLAT:g_StepConstants+832
//        movdqu  xmm0, XMMWORD PTR [rbx+144]
//        movaps  XMMWORD PTR [rsp+48], xmm7
//        movdqu  xmm15, XMMWORD PTR [rbx+192]
//        pxor    xmm2, xmm3
//        movaps  XMMWORD PTR [rsp+80], xmm3
//        movdqu  xmm3, XMMWORD PTR [rbx+80]
//        movdqu  xmm12, XMMWORD PTR [rbx+160]
//        movdqu  xmm11, XMMWORD PTR [rbx+176]
//        movaps  XMMWORD PTR [rsp+96], xmm0
//        movdqu  xmm13, XMMWORD PTR [rbx+208]
//        pxor    xmm3, xmm0
//        movdqu  xmm0, XMMWORD PTR [rbx+32]
//        movaps  XMMWORD PTR [rsp+64], xmm4
//        movaps  XMMWORD PTR [rsp+16], xmm15
//        pxor    xmm0, xmm7
//        movdqa  xmm7, XMMWORD PTR .LC8[rip]
//        movaps  XMMWORD PTR [rsp+112], xmm12
//        paddd   xmm0, xmm2
//        movaps  XMMWORD PTR [rsp], xmm11
//        movdqa  xmm5, xmm0
//        pslld   xmm0, 29
//        movaps  XMMWORD PTR [rsp+32], xmm13
//        movdqa  xmm1, xmm0
//        movdqu  xmm0, XMMWORD PTR [rbx+48]
//        psrld   xmm5, 3
//        por     xmm1, xmm5
//        pxor    xmm0, xmm4
//        pxor    xmm7, xmm1
//        movdqa  xmm1, XMMWORD PTR .LC10[rip]
//        paddd   xmm0, xmm3
//        paddd   xmm2, xmm7
//        movdqa  xmm4, xmm0
//        pslld   xmm0, 29
//        psrld   xmm4, 3
//        por     xmm0, xmm4
//        pxor    xmm0, XMMWORD PTR .LC9[rip]
//        movdqa  xmm8, xmm0
//        movdqa  xmm0, xmm2
//        psrld   xmm0, 31
//        pslld   xmm2, 1
//        paddd   xmm3, xmm8
//        por     xmm2, xmm0
//        movdqa  xmm0, xmm3
//        pand    xmm1, xmm2
//        psrld   xmm0, 31
//        paddd   xmm7, xmm2
//        pslld   xmm3, 1
//        movdqa  xmm4, xmm1
//        psrld   xmm4, 24
//        pslld   xmm1, 8
//        por     xmm3, xmm0
//        movdqa  xmm0, XMMWORD PTR .LC11[rip]
//        pxor    xmm1, xmm4
//        pand    xmm0, xmm2
//        pshufd  xmm2, xmm7, 210
//        pxor    xmm0, xmm1
//        movdqa  xmm1, XMMWORD PTR .LC12[rip]
//        pxor    xmm2, xmm15
//        pand    xmm1, xmm0
//        pand    xmm0, XMMWORD PTR .LC13[rip]
//        movdqa  xmm4, xmm1
//        pslld   xmm1, 8
//        psrld   xmm4, 24
//        pxor    xmm1, xmm4
//        movdqa  xmm4, XMMWORD PTR .LC15[rip]
//        pxor    xmm0, xmm1
//        movdqa  xmm1, XMMWORD PTR .LC14[rip]
//        pand    xmm4, xmm3
//        pand    xmm1, xmm0
//        movdqa  xmm5, xmm4
//        movdqa  xmm10, xmm1
//        psrld   xmm5, 24
//        movdqa  xmm9, xmm1
//        movdqa  xmm1, XMMWORD PTR .LC14[rip]
//        pslld   xmm4, 8
//        pslld   xmm9, 8
//        pxor    xmm4, xmm5
//        psrld   xmm10, 24
//        pand    xmm1, xmm3
//        pxor    xmm1, xmm4
//        movdqa  xmm4, XMMWORD PTR .LC13[rip]
//        paddd   xmm3, xmm8
//        pshufd  xmm3, xmm3, 210
//        pand    xmm4, xmm1
//        pand    xmm1, XMMWORD PTR .LC12[rip]
//        pxor    xmm3, xmm12
//        movdqa  xmm5, xmm4
//        pslld   xmm4, 8
//        paddd   xmm3, xmm2
//        psrld   xmm5, 24
//        pxor    xmm4, xmm5
//        movdqa  xmm5, xmm9
//        pxor    xmm1, xmm4
//        movdqa  xmm4, XMMWORD PTR .LC11[rip]
//        pxor    xmm5, xmm10
//        pand    xmm0, XMMWORD PTR .LC15[rip]
//        movdqa  xmm7, XMMWORD PTR .LC17[rip]
//        movdqa  xmm14, XMMWORD PTR .LC11[rip]
//        pand    xmm4, xmm1
//        pand    xmm1, XMMWORD PTR .LC10[rip]
//        pxor    xmm0, xmm5
//        movdqa  xmm6, xmm4
//        pslld   xmm4, 8
//        pshufd  xmm0, xmm0, 108
//        psrld   xmm6, 24
//        pxor    xmm0, xmm13
//        pxor    xmm4, xmm6
//        movdqa  xmm6, XMMWORD PTR .LC14[rip]
//        pxor    xmm1, xmm4
//        movdqa  xmm4, xmm3
//        psrld   xmm4, 27
//        pslld   xmm3, 5
//        pshufd  xmm1, xmm1, 108
//        pxor    xmm1, xmm11
//        por     xmm3, xmm4
//        pxor    xmm3, XMMWORD PTR .LC16[rip]
//        paddd   xmm1, xmm0
//        movdqa  xmm5, xmm1
//        pslld   xmm1, 5
//        paddd   xmm2, xmm3
//        psrld   xmm5, 27
//        movdqa  xmm15, xmm3
//        por     xmm1, xmm5
//        pxor    xmm7, xmm1
//        movdqa  xmm1, xmm2
//        psrld   xmm1, 15
//        pslld   xmm2, 17
//        paddd   xmm0, xmm7
//        por     xmm2, xmm1
//        movdqa  xmm1, xmm0
//        pslld   xmm0, 17
//        psrld   xmm1, 15
//        pand    xmm14, xmm2
//        movdqa  xmm3, xmm0
//        movdqa  xmm0, XMMWORD PTR .LC10[rip]
//        paddd   xmm15, xmm2
//        por     xmm3, xmm1
//        pshufd  xmm15, xmm15, 210
//        pand    xmm0, xmm2
//        pand    xmm6, xmm3
//        movdqa  xmm2, xmm7
//        movdqa  xmm1, xmm0
//        pslld   xmm0, 8
//        paddd   xmm2, xmm3
//        psrld   xmm1, 24
//        pshufd  xmm2, xmm2, 210
//        pxor    xmm0, xmm1
//        pxor    xmm14, xmm0
//        movdqa  xmm0, XMMWORD PTR .LC12[rip]
//        pand    xmm0, xmm14
//        pand    xmm14, XMMWORD PTR .LC13[rip]
//        movdqa  xmm1, xmm0
//        pslld   xmm0, 8
//        psrld   xmm1, 24
//        pxor    xmm0, xmm1
//        pxor    xmm14, xmm0
//        movdqa  xmm0, XMMWORD PTR .LC14[rip]
//        pand    xmm0, xmm14
//        pand    xmm14, XMMWORD PTR .LC15[rip]
//        movdqa  xmm5, xmm0
//        movdqa  xmm4, xmm0
//        movdqa  xmm0, XMMWORD PTR .LC15[rip]
//        pslld   xmm4, 8
//        psrld   xmm5, 24
//        pand    xmm0, xmm3
//        movdqa  xmm3, xmm4
//        movdqa  xmm1, xmm0
//        pslld   xmm0, 8
//        pxor    xmm3, xmm5
//        psrld   xmm1, 24
//        pxor    xmm14, xmm3
//        movdqa  xmm3, xmm2
//        pxor    xmm0, xmm1
//        pshufd  xmm14, xmm14, 108
//        movdqa  xmm2, xmm15
//        pxor    xmm6, xmm0
//        movdqa  xmm0, XMMWORD PTR .LC13[rip]
//        movdqa  xmm15, xmm14
//        pand    xmm0, xmm6
//        pand    xmm6, XMMWORD PTR .LC12[rip]
//        movdqa  xmm1, xmm0
//        pslld   xmm0, 8
//        psrld   xmm1, 24
//        pxor    xmm0, xmm1
//        pxor    xmm6, xmm0
//        movdqa  xmm0, XMMWORD PTR .LC11[rip]
//        pand    xmm0, xmm6
//        pand    xmm6, XMMWORD PTR .LC10[rip]
//        movdqa  xmm1, xmm0
//        pslld   xmm0, 8
//        psrld   xmm1, 24
//        pxor    xmm0, xmm1
//        pxor    xmm6, xmm0
//        pshufd  xmm6, xmm6, 108
//.L41:
//        movdqa  xmm12, XMMWORD PTR [rsp+112]
//        pshufd  xmm0, XMMWORD PTR [rsp+48], 75
//        add     rax, 64
//        pshufd  xmm8, XMMWORD PTR [rsp+80], 75
//        paddd   xmm8, XMMWORD PTR [rsp+16]
//        pshufd  xmm7, XMMWORD PTR [rsp+96], 147
//        paddd   xmm7, XMMWORD PTR [rsp+32]
//        movdqa  xmm9, XMMWORD PTR .LC14[rip]
//        paddd   xmm0, xmm12
//        pshufd  xmm12, xmm12, 75
//        movdqa  xmm13, xmm0
//        movaps  XMMWORD PTR [rsp+48], xmm0
//        pxor    xmm2, xmm8
//        pxor    xmm15, xmm7
//        pxor    xmm3, xmm13
//        paddd   xmm12, xmm13
//        movaps  XMMWORD PTR [rsp+96], xmm7
//        movaps  XMMWORD PTR [rsp+128], xmm0
//        pshufd  xmm0, XMMWORD PTR [rsp+64], 147
//        paddd   xmm0, XMMWORD PTR [rsp]
//        movdqa  xmm4, xmm3
//        paddd   xmm4, xmm2
//        movdqa  xmm3, XMMWORD PTR .LC11[rip]
//        movaps  XMMWORD PTR [rsp+80], xmm8
//        pxor    xmm6, xmm0
//        movdqa  xmm14, xmm0
//        movaps  XMMWORD PTR [rsp+64], xmm0
//        movdqa  xmm5, xmm4
//        psrld   xmm5, 3
//        movaps  XMMWORD PTR [rsp+144], xmm0
//        movdqa  xmm0, xmm6
//        paddd   xmm0, xmm15
//        pslld   xmm4, 29
//        movaps  XMMWORD PTR [rsp+112], xmm12
//        movdqa  xmm1, xmm0
//        pslld   xmm0, 29
//        por     xmm4, xmm5
//        pxor    xmm4, XMMWORD PTR [rax-64]
//        psrld   xmm1, 3
//        por     xmm0, xmm1
//        pxor    xmm0, XMMWORD PTR [rax-48]
//        paddd   xmm2, xmm4
//        movdqa  xmm1, xmm2
//        pslld   xmm2, 1
//        paddd   xmm15, xmm0
//        psrld   xmm1, 31
//        por     xmm2, xmm1
//        movdqa  xmm1, xmm15
//        psrld   xmm1, 31
//        pslld   xmm15, 1
//        pand    xmm3, xmm2
//        por     xmm15, xmm1
//        movdqa  xmm1, XMMWORD PTR .LC10[rip]
//        paddd   xmm4, xmm2
//        paddd   xmm0, xmm15
//        pand    xmm1, xmm2
//        pshufd  xmm0, xmm0, 210
//        pshufd  xmm2, xmm4, 210
//        movdqa  xmm5, xmm1
//        pslld   xmm1, 8
//        pxor    xmm0, xmm12
//        psrld   xmm5, 24
//        pxor    xmm1, xmm5
//        pxor    xmm3, xmm1
//        movdqa  xmm1, XMMWORD PTR .LC12[rip]
//        pand    xmm1, xmm3
//        pand    xmm3, XMMWORD PTR .LC13[rip]
//        movdqa  xmm5, xmm1
//        pslld   xmm1, 8
//        psrld   xmm5, 24
//        pxor    xmm1, xmm5
//        movdqa  xmm5, XMMWORD PTR .LC15[rip]
//        pxor    xmm3, xmm1
//        movdqa  xmm1, XMMWORD PTR .LC14[rip]
//        pand    xmm5, xmm15
//        pand    xmm9, xmm3
//        pand    xmm3, XMMWORD PTR .LC15[rip]
//        movdqa  xmm6, xmm5
//        pslld   xmm5, 8
//        pand    xmm1, xmm15
//        psrld   xmm6, 24
//        movdqa  xmm11, xmm9
//        psrld   xmm11, 24
//        pslld   xmm9, 8
//        pxor    xmm5, xmm6
//        pxor    xmm1, xmm5
//        movdqa  xmm5, XMMWORD PTR .LC13[rip]
//        pxor    xmm9, xmm11
//        pshufd  xmm11, XMMWORD PTR [rsp+16], 75
//        paddd   xmm11, xmm8
//        pxor    xmm3, xmm9
//        pand    xmm5, xmm1
//        pxor    xmm2, xmm11
//        pshufd  xmm3, xmm3, 108
//        movdqa  xmm6, xmm5
//        pand    xmm1, XMMWORD PTR .LC12[rip]
//        paddd   xmm0, xmm2
//        psrld   xmm6, 24
//        pslld   xmm5, 8
//        pxor    xmm5, xmm6
//        pxor    xmm1, xmm5
//        movdqa  xmm5, XMMWORD PTR .LC11[rip]
//        pand    xmm5, xmm1
//        pand    xmm1, XMMWORD PTR .LC10[rip]
//        movdqa  xmm10, xmm5
//        pslld   xmm5, 8
//        psrld   xmm10, 24
//        pxor    xmm5, xmm10
//        pshufd  xmm10, XMMWORD PTR [rsp], 147
//        paddd   xmm10, xmm14
//        pxor    xmm1, xmm5
//        movaps  XMMWORD PTR [rsp], xmm10
//        movdqa  xmm5, xmm0
//        movaps  XMMWORD PTR [rsp+16], xmm11
//        pshufd  xmm1, xmm1, 108
//        psrld   xmm5, 27
//        pshufd  xmm9, XMMWORD PTR [rsp+32], 147
//        paddd   xmm9, xmm7
//        pxor    xmm1, xmm10
//        movdqa  xmm13, XMMWORD PTR .LC14[rip]
//        pxor    xmm3, xmm9
//        pslld   xmm0, 5
//        movaps  XMMWORD PTR [rsp+32], xmm9
//        paddd   xmm1, xmm3
//        por     xmm0, xmm5
//        pxor    xmm0, XMMWORD PTR [rax-32]
//        movdqa  xmm5, XMMWORD PTR .LC10[rip]
//        movdqa  xmm4, xmm1
//        pslld   xmm1, 5
//        psrld   xmm4, 27
//        paddd   xmm2, xmm0
//        por     xmm1, xmm4
//        pxor    xmm1, XMMWORD PTR [rax-16]
//        movdqa  xmm4, xmm2
//        psrld   xmm4, 15
//        pslld   xmm2, 17
//        paddd   xmm3, xmm1
//        por     xmm2, xmm4
//        movdqa  xmm4, xmm3
//        pand    xmm5, xmm2
//        paddd   xmm0, xmm2
//        psrld   xmm4, 15
//        pslld   xmm3, 17
//        movdqa  xmm6, xmm5
//        por     xmm3, xmm4
//        psrld   xmm6, 24
//        pshufd  xmm0, xmm0, 210
//        movdqa  xmm4, XMMWORD PTR .LC11[rip]
//        pslld   xmm5, 8
//        paddd   xmm1, xmm3
//        pxor    xmm5, xmm6
//        pand    xmm4, xmm2
//        pshufd  xmm1, xmm1, 210
//        pxor    xmm4, xmm5
//        movdqa  xmm5, XMMWORD PTR .LC12[rip]
//        movdqa  xmm2, xmm0
//        pand    xmm5, xmm4
//        pand    xmm4, XMMWORD PTR .LC13[rip]
//        movdqa  xmm6, xmm5
//        pslld   xmm5, 8
//        psrld   xmm6, 24
//        pxor    xmm5, xmm6
//        movdqa  xmm6, XMMWORD PTR .LC15[rip]
//        pxor    xmm4, xmm5
//        movdqa  xmm5, XMMWORD PTR .LC14[rip]
//        pand    xmm6, xmm3
//        pand    xmm13, xmm4
//        pand    xmm4, XMMWORD PTR .LC15[rip]
//        movdqa  xmm14, xmm6
//        pslld   xmm6, 8
//        pand    xmm5, xmm3
//        psrld   xmm14, 24
//        movdqa  xmm15, xmm13
//        movdqa  xmm3, xmm1
//        pxor    xmm6, xmm14
//        psrld   xmm15, 24
//        pxor    xmm5, xmm6
//        movdqa  xmm6, XMMWORD PTR .LC13[rip]
//        pslld   xmm13, 8
//        pxor    xmm13, xmm15
//        pand    xmm6, xmm5
//        pand    xmm5, XMMWORD PTR .LC12[rip]
//        pxor    xmm4, xmm13
//        movdqa  xmm14, xmm6
//        pslld   xmm6, 8
//        pshufd  xmm4, xmm4, 108
//        psrld   xmm14, 24
//        movdqa  xmm15, xmm4
//        pxor    xmm6, xmm14
//        pxor    xmm5, xmm6
//        movdqa  xmm6, XMMWORD PTR .LC11[rip]
//        pand    xmm6, xmm5
//        pand    xmm5, XMMWORD PTR .LC10[rip]
//        movdqa  xmm14, xmm6
//        pslld   xmm6, 8
//        psrld   xmm14, 24
//        pxor    xmm6, xmm14
//        pxor    xmm5, xmm6
//        pshufd  xmm5, xmm5, 108
//        movdqa  xmm6, xmm5
//        cmp     rdx, rax
//        jne     .L41
//        pshufd  xmm6, xmm8, 75
//        movdqa  xmm3, xmm5
//        lea     rdi, [rbx+8]
//        xor     eax, eax
//        pshufd  xmm7, xmm7, 147
//        pxor    xmm1, xmm0
//        pxor    xmm3, xmm4
//        and     rdi, -8
//        pshufd  xmm2, XMMWORD PTR [rsp+128], 75
//        paddd   xmm2, xmm12
//        paddd   xmm7, xmm9
//        pshufd  xmm5, XMMWORD PTR [rsp+144], 147
//        movdqa  xmm8, xmm2
//        movdqa  xmm2, xmm6
//        paddd   xmm5, xmm10
//        paddd   xmm2, xmm11
//        pxor    xmm5, xmm7
//        pxor    xmm8, xmm2
//        pxor    xmm5, xmm3
//        pxor    xmm1, xmm8
//        movups  XMMWORD PTR [rbp+16], xmm5
//        movups  XMMWORD PTR [rbp+0], xmm1
//        mov     QWORD PTR [rbx], 0
//        mov     QWORD PTR [rbx+216], 0
//        sub     rbx, rdi
//        lea     ecx, [rbx+224]
//        shr     ecx, 3
//        rep stosq
//        add     rsp, 168
//        pop     rbx
//        pop     rbp
//        ret
//g_StepConstants:
//        .long   -1854099568
//        .long   1813713058
//        .long   1865754947
//        .long   -814251453
//        .long   753628274
//        .long   703164402
//        .long   -1969511384
//        .long   787162690
//        .long   237781025
//        .long   -2027179250
//        .long   -1537315662
//        .long   1190774290
//        .long   408938142
//        .long   324624923
//        .long   641715378
//        .long   437348464
//        .long   980181295
//        .long   -1294024299
//        .long   46866262
//        .long   1086314584
//        .long   2017887414
//        .long   1824226606
//        .long   1712094936
//        .long   729405578
//        .long   -1496477591
//        .long   -1851762873
//        .long   -840272552
//        .long   9973912
//        .long   -1093977298
//        .long   674802586
//        .long   1921734462
//        .long   -1521327435
//        .long   1952315919
//        .long   838311640
//        .long   -1201680859
//        .long   -1731671912
//        .long   -1979291156
//        .long   1625572802
//        .long   -35482192
//        .long   -139126694
//        .long   -765982333
//        .long   697726729
//        .long   405664733
//        .long   1636917552
//        .long   -1871686026
//        .long   1160382498
//        .long   -528988499
//        .long   -1395821743
//        .long   711036245
//        .long   -1072857038
//        .long   1176621301
//        .long   -218984047
//        .long   13028614
//        .long   1865558631
//        .long   1488909453
//        .long   2051065085
//        .long   -1947278721
//        .long   -846351630
//        .long   1743963707
//        .long   -444324989
//        .long   -946613498
//        .long   -1587768874
//        .long   397640165
//        .long   -1157250441
//        .long   2060542474
//        .long   1530963391
//        .long   1522074018
//        .long   1772595048
//        .long   1533467085
//        .long   -34707337
//        .long   -885692676
//        .long   -1060619726
//        .long   1278438532
//        .long   -1679399654
//        .long   330998780
//        .long   290394065
//        .long   -1035950296
//        .long   -329850764
//        .long   10249159
//        .long   -1996329742
//        .long   2141188304
//        .long   -2108981323
//        .long   -832659441
//        .long   1616830690
//        .long   48746474
//        .long   1127699808
//        .long   -1660933433
//        .long   -1955635333
//        .long   529580367
//        .long   -842320585
//        .long   754888669
//        .long   -1086340286
//        .long   -357058068
//        .long   2056041891
//        .long   -1658129820
//        .long   -87106810
//        .long   -1336799122
//        .long   -1728025340
//        .long   732859657
//        .long   -16579071
//        .long   -1569610026
//        .long   121788701
//        .long   -937601279
//        .long   -911613440
//        .long   40009502
//        .long   -1720880739
//        .long   -633629132
//        .long   34146304
//        .long   337149304
//        .long   1234925092
//        .long   -445359927
//        .long   1938547401
//        .long   1692507424
//        .long   115281718
//        .long   366070542
//        .long   185665538
//        .long   748025228
//        .long   -451831187
//        .long   1506615982
//        .long   -9524164
//        .long   1182711172
//        .long   -437369540
//        .long   -409151197
//        .long   471403021
//        .long   -1031044680
//        .long   -156497838
//        .long   653250823
//        .long   1847257403
//        .long   -985906742
//        .long   -734821855
//        .long   -2072642294
//        .long   891492137
//        .long   1309489786
//        .long   -1565238456
//        .long   381737005
//        .long   -1996217175
//        .long   25037967
//        .long   129344501
//        .long   -93124722
//        .long   1480212574
//        .long   1531425992
//        .long   1460642794
//        .long   -677837501
//        .long   -1926705614
//        .long   2139790224
//        .long   -1118201604
//        .long   1836750472
//        .long   -1837457738
//        .long   -1560511197
//        .long   1723043393
//        .long   1889421594
//        .long   -1242103617
//        .long   170273807
//        .long   384996793
//        .long   -401531147
//        .long   219764040
//        .long   -1619418939
//        .long   440520615
//        .long   250045322
//        .long   -1393049484
//        .long   810695449
//        .long   159928015
//        .long   -114289443
//        .long   644783445
//        .long   421160548
//        .long   1544648385
//        .long   -162861672
//        .long   -1531431648
//        .long   -2104667063
//        .long   -1848173096
//        .long   692384470
//        .long   -1794379141
//        .long   864078461
//        .long   1837643805
//        .long   1094223502
//        .long   1559057860
//        .long   266444235
//        .long   1760884329
//        .long   1849149695
//        .long   -1593938416
//        .long   -1270134864
//        .long   -171944566
//        .long   2042565839
//        .long   1245911584
//        .long   -243308326
//        .long   1575706577
//        .long   -1509898131
//        .long   -1620279632
//        .long   -7145372
//        .long   -1240080257
//        .long   953410632
//        .long   -1923270038
//        .long   1894269899
//        .long   1229664558
//        .long   -1496762861
//        .long   199962447
//        .long   -1836360349
//        .long   -873045963
//        .long   213396608
//        .long   -359160585
//        .long   1496059707
//        .long   -1803788425
//        .long   1879001529
//        .long   -149258662
//        .long   501793013
//        .long   -1034512896
//        .long   -991646580
//        .long   -2109956895
//.LC0:
//        .quad   -153997021875728609
//        .quad   1841527622716375976
//.LC1:
//        .quad   -5695466077790390131
//        .quad   8852904047995281860
//.LC2:
//        .quad   3419602469714219896
//        .quad   -984810856687653483
//.LC3:
//        .quad   1406927143955412346
//        .quad   -2298248402025429884
//.LC4:
//        .quad   7122715107427485907
//        .quad   5503410031139639979
//.LC5:
//        .quad   2218598310921636520
//        .quad   3551039192594811326
//.LC6:
//        .quad   -694899722053307143
//        .quad   6249004012466285246
//.LC7:
//        .quad   -1293352058239094138
//        .quad   -8378717999420185195
//.LC8:
//        .quad   7789838270879018896
//        .quad   -3497183359489726141
//.LC9:
//        .quad   3020068111055025266
//        .quad   3380838012506842152
//.LC10:
//        .quad   -4294967296
//        .quad   -1
//.LC11:
//        .quad   4294967295
//        .quad   0
//.LC12:
//        .quad   0
//        .quad   -1
//.LC13:
//        .quad   -1
//        .quad   0
//.LC14:
//        .quad   0
//        .quad   -4294967296
//.LC15:
//        .quad   -1
//        .quad   4294967295
//.LC16:
//        .quad   -8706668581642026975
//        .quad   5114336635225271474
//.LC17:
//        .quad   1394253428160456350
//        .quad   1878397350477548722
