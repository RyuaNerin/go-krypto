#include <arm_neon.h>


#define VROTL(x,i) vsriq_n_u32(vshlq_n_u32(x,i),x,32-i)
#define VROTR(x,i) vsliq_n_u32(vshrq_n_u32(x,i),x,32-i)

#define XAR3_NEON(cur, pre, tmp, rk1, rk2)																			\
	tmp = vaddq_u32(veorq_u32(pre, vdupq_n_u32(rk1)), veorq_u32(cur, vdupq_n_u32(rk2)));		\
	cur = VROTR(tmp, 3);
#define XAR5_NEON(cur, pre, tmp, rk1, rk2)																			\
	tmp = vaddq_u32(veorq_u32(pre, vdupq_n_u32(rk1)), veorq_u32(cur, vdupq_n_u32(rk2)));		\
	cur = VROTR(tmp, 5);
#define XAR9_NEON(cur, pre, tmp, rk1, rk2)																			\
	tmp = vaddq_u32(veorq_u32(pre, vdupq_n_u32(rk1)), veorq_u32(cur, vdupq_n_u32(rk2)));		\
	cur = VROTR(tmp, 23);

#define XSR9_NEON(cur, pre, rk1, rk2)																																		\
	cur = veorq_u32(vsubq_u32(VROTR(cur,9), veorq_u32(pre, vdupq_n_u32(rk1))), vdupq_n_u32(rk2));
#define XSR5_NEON(cur, pre, rk1, rk2)																																		\
	cur = veorq_u32(vsubq_u32(VROTR(cur, 27), veorq_u32(pre, vdupq_n_u32(rk1))), vdupq_n_u32(rk2));
#define XSR3_NEON(cur, pre, rk1, rk2)																																		\
	cur = veorq_u32(vsubq_u32(VROTR(cur, 29), veorq_u32(pre, vdupq_n_u32(rk1))), vdupq_n_u32(rk2));

#define x3 x.val[3]
#define x2 x.val[2]
#define x1 x.val[1]
#define x0 x.val[0]

void lea_encrypt_4block(char *ct, const char *pt, const unsigned int *rk, const unsigned long round)
{
	uint32x4x4_t x;
	uint32x4_t tmp;
	
	x = vld4q_u32((unsigned int*)pt);

	XAR3_NEON(x3, x2, tmp, rk[  4], rk[  5]);
	XAR5_NEON(x2, x1, tmp, rk[  2], rk[  3]);
	XAR9_NEON(x1, x0, tmp, rk[  0], rk[  1]);
	XAR3_NEON(x0, x3, tmp, rk[ 10], rk[ 11]);
	XAR5_NEON(x3, x2, tmp, rk[  8], rk[  9]);
	XAR9_NEON(x2, x1, tmp, rk[  6], rk[  7]);
	XAR3_NEON(x1, x0, tmp, rk[ 16], rk[ 17]);
	XAR5_NEON(x0, x3, tmp, rk[ 14], rk[ 15]);
	XAR9_NEON(x3, x2, tmp, rk[ 12], rk[ 13]);
	XAR3_NEON(x2, x1, tmp, rk[ 22], rk[ 23]);
	XAR5_NEON(x1, x0, tmp, rk[ 20], rk[ 21]);
	XAR9_NEON(x0, x3, tmp, rk[ 18], rk[ 19]);

	XAR3_NEON(x3, x2, tmp, rk[ 28], rk[ 29]);
	XAR5_NEON(x2, x1, tmp, rk[ 26], rk[ 27]);
	XAR9_NEON(x1, x0, tmp, rk[ 24], rk[ 25]);
	XAR3_NEON(x0, x3, tmp, rk[ 34], rk[ 35]);
	XAR5_NEON(x3, x2, tmp, rk[ 32], rk[ 33]);
	XAR9_NEON(x2, x1, tmp, rk[ 30], rk[ 31]);
	XAR3_NEON(x1, x0, tmp, rk[ 40], rk[ 41]);
	XAR5_NEON(x0, x3, tmp, rk[ 38], rk[ 39]);
	XAR9_NEON(x3, x2, tmp, rk[ 36], rk[ 37]);
	XAR3_NEON(x2, x1, tmp, rk[ 46], rk[ 47]);
	XAR5_NEON(x1, x0, tmp, rk[ 44], rk[ 45]);
	XAR9_NEON(x0, x3, tmp, rk[ 42], rk[ 43]);

	XAR3_NEON(x3, x2, tmp, rk[ 52], rk[ 53]);
	XAR5_NEON(x2, x1, tmp, rk[ 50], rk[ 51]);
	XAR9_NEON(x1, x0, tmp, rk[ 48], rk[ 49]);
	XAR3_NEON(x0, x3, tmp, rk[ 58], rk[ 59]);
	XAR5_NEON(x3, x2, tmp, rk[ 56], rk[ 57]);
	XAR9_NEON(x2, x1, tmp, rk[ 54], rk[ 55]);
	XAR3_NEON(x1, x0, tmp, rk[ 64], rk[ 65]);
	XAR5_NEON(x0, x3, tmp, rk[ 62], rk[ 63]);
	XAR9_NEON(x3, x2, tmp, rk[ 60], rk[ 61]);
	XAR3_NEON(x2, x1, tmp, rk[ 70], rk[ 71]);
	XAR5_NEON(x1, x0, tmp, rk[ 68], rk[ 69]);
	XAR9_NEON(x0, x3, tmp, rk[ 66], rk[ 67]);

	XAR3_NEON(x3, x2, tmp, rk[ 76], rk[ 77]);
	XAR5_NEON(x2, x1, tmp, rk[ 74], rk[ 75]);
	XAR9_NEON(x1, x0, tmp, rk[ 72], rk[ 73]);
	XAR3_NEON(x0, x3, tmp, rk[ 82], rk[ 83]);
	XAR5_NEON(x3, x2, tmp, rk[ 80], rk[ 81]);
	XAR9_NEON(x2, x1, tmp, rk[ 78], rk[ 79]);
	XAR3_NEON(x1, x0, tmp, rk[ 88], rk[ 89]);
	XAR5_NEON(x0, x3, tmp, rk[ 86], rk[ 87]);
	XAR9_NEON(x3, x2, tmp, rk[ 84], rk[ 85]);
	XAR3_NEON(x2, x1, tmp, rk[ 94], rk[ 95]);
	XAR5_NEON(x1, x0, tmp, rk[ 92], rk[ 93]);
	XAR9_NEON(x0, x3, tmp, rk[ 90], rk[ 91]);

	XAR3_NEON(x3, x2, tmp, rk[100], rk[101]);
	XAR5_NEON(x2, x1, tmp, rk[ 98], rk[ 99]);
	XAR9_NEON(x1, x0, tmp, rk[ 96], rk[ 97]);
	XAR3_NEON(x0, x3, tmp, rk[106], rk[107]);
	XAR5_NEON(x3, x2, tmp, rk[104], rk[105]);
	XAR9_NEON(x2, x1, tmp, rk[102], rk[103]);
	XAR3_NEON(x1, x0, tmp, rk[112], rk[113]);
	XAR5_NEON(x0, x3, tmp, rk[110], rk[111]);
	XAR9_NEON(x3, x2, tmp, rk[108], rk[109]);
	XAR3_NEON(x2, x1, tmp, rk[118], rk[119]);
	XAR5_NEON(x1, x0, tmp, rk[116], rk[117]);
	XAR9_NEON(x0, x3, tmp, rk[114], rk[115]);

	XAR3_NEON(x3, x2, tmp, rk[124], rk[125]);
	XAR5_NEON(x2, x1, tmp, rk[122], rk[123]);
	XAR9_NEON(x1, x0, tmp, rk[120], rk[121]);
	XAR3_NEON(x0, x3, tmp, rk[130], rk[131]);
	XAR5_NEON(x3, x2, tmp, rk[128], rk[129]);
	XAR9_NEON(x2, x1, tmp, rk[126], rk[127]);
	XAR3_NEON(x1, x0, tmp, rk[136], rk[137]);
	XAR5_NEON(x0, x3, tmp, rk[134], rk[135]);
	XAR9_NEON(x3, x2, tmp, rk[132], rk[133]);
	XAR3_NEON(x2, x1, tmp, rk[142], rk[143]);
	XAR5_NEON(x1, x0, tmp, rk[140], rk[141]);
	XAR9_NEON(x0, x3, tmp, rk[138], rk[139]);

	if(round > 24)
	{
		XAR3_NEON(x3, x2, tmp, rk[148], rk[149]);
		XAR5_NEON(x2, x1, tmp, rk[146], rk[147]);
		XAR9_NEON(x1, x0, tmp, rk[144], rk[145]);
		XAR3_NEON(x0, x3, tmp, rk[154], rk[155]);
		XAR5_NEON(x3, x2, tmp, rk[152], rk[153]);
		XAR9_NEON(x2, x1, tmp, rk[150], rk[151]);
		XAR3_NEON(x1, x0, tmp, rk[160], rk[161]);
		XAR5_NEON(x0, x3, tmp, rk[158], rk[159]);
		XAR9_NEON(x3, x2, tmp, rk[156], rk[157]);
		XAR3_NEON(x2, x1, tmp, rk[166], rk[167]);
		XAR5_NEON(x1, x0, tmp, rk[164], rk[165]);
		XAR9_NEON(x0, x3, tmp, rk[162], rk[163]);
	}

	if(round > 28)
	{
		XAR3_NEON(x3, x2, tmp, rk[172], rk[173]);
		XAR5_NEON(x2, x1, tmp, rk[170], rk[171]);
		XAR9_NEON(x1, x0, tmp, rk[168], rk[169]);
		XAR3_NEON(x0, x3, tmp, rk[178], rk[179]);
		XAR5_NEON(x3, x2, tmp, rk[176], rk[177]);
		XAR9_NEON(x2, x1, tmp, rk[174], rk[175]);
		XAR3_NEON(x1, x0, tmp, rk[184], rk[185]);
		XAR5_NEON(x0, x3, tmp, rk[182], rk[183]);
		XAR9_NEON(x3, x2, tmp, rk[180], rk[181]);
		XAR3_NEON(x2, x1, tmp, rk[190], rk[191]);
		XAR5_NEON(x1, x0, tmp, rk[188], rk[189]);
		XAR9_NEON(x0, x3, tmp, rk[186], rk[187]);
	}

	vst4q_u32((unsigned int *)ct, x);
}

void lea_decrypt_4block(char *pt, const char *ct, const unsigned int *rk, const unsigned long round)
{
	uint32x4x4_t x;
	unsigned int t[4], itmp;

	x = vld4q_u32((unsigned int*)ct);

	if(round > 28)
	{
		XSR9_NEON(x0, x3, rk[186], rk[187]);
		XSR5_NEON(x1, x0, rk[188], rk[189]);
		XSR3_NEON(x2, x1, rk[190], rk[191]);
		XSR9_NEON(x3, x2, rk[180], rk[181]);
		XSR5_NEON(x0, x3, rk[182], rk[183]);
		XSR3_NEON(x1, x0, rk[184], rk[185]);
		XSR9_NEON(x2, x1, rk[174], rk[175]);
		XSR5_NEON(x3, x2, rk[176], rk[177]);
		XSR3_NEON(x0, x3, rk[178], rk[179]);
		XSR9_NEON(x1, x0, rk[168], rk[169]);
		XSR5_NEON(x2, x1, rk[170], rk[171]);
		XSR3_NEON(x3, x2, rk[172], rk[173]);
	}

	if(round > 24)
	{
		XSR9_NEON(x0, x3, rk[162], rk[163]);
		XSR5_NEON(x1, x0, rk[164], rk[165]);
		XSR3_NEON(x2, x1, rk[166], rk[167]);
		XSR9_NEON(x3, x2, rk[156], rk[157]);
		XSR5_NEON(x0, x3, rk[158], rk[159]);
		XSR3_NEON(x1, x0, rk[160], rk[161]);
		XSR9_NEON(x2, x1, rk[150], rk[151]);
		XSR5_NEON(x3, x2, rk[152], rk[153]);
		XSR3_NEON(x0, x3, rk[154], rk[155]);
		XSR9_NEON(x1, x0, rk[144], rk[145]);
		XSR5_NEON(x2, x1, rk[146], rk[147]);
		XSR3_NEON(x3, x2, rk[148], rk[149]);
	}

	XSR9_NEON(x0, x3, rk[138], rk[139]);
	XSR5_NEON(x1, x0, rk[140], rk[141]);
	XSR3_NEON(x2, x1, rk[142], rk[143]);
	XSR9_NEON(x3, x2, rk[132], rk[133]);
	XSR5_NEON(x0, x3, rk[134], rk[135]);
	XSR3_NEON(x1, x0, rk[136], rk[137]);
	XSR9_NEON(x2, x1, rk[126], rk[127]);
	XSR5_NEON(x3, x2, rk[128], rk[129]);
	XSR3_NEON(x0, x3, rk[130], rk[131]);
	XSR9_NEON(x1, x0, rk[120], rk[121]);
	XSR5_NEON(x2, x1, rk[122], rk[123]);
	XSR3_NEON(x3, x2, rk[124], rk[125]);

	XSR9_NEON(x0, x3, rk[114], rk[115]);
	XSR5_NEON(x1, x0, rk[116], rk[117]);
	XSR3_NEON(x2, x1, rk[118], rk[119]);
	XSR9_NEON(x3, x2, rk[108], rk[109]);
	XSR5_NEON(x0, x3, rk[110], rk[111]);
	XSR3_NEON(x1, x0, rk[112], rk[113]);
	XSR9_NEON(x2, x1, rk[102], rk[103]);
	XSR5_NEON(x3, x2, rk[104], rk[105]);
	XSR3_NEON(x0, x3, rk[106], rk[107]);
	XSR9_NEON(x1, x0, rk[ 96], rk[ 97]);
	XSR5_NEON(x2, x1, rk[ 98], rk[ 99]);
	XSR3_NEON(x3, x2, rk[100], rk[101]);

	XSR9_NEON(x0, x3, rk[ 90], rk[ 91]);
	XSR5_NEON(x1, x0, rk[ 92], rk[ 93]);
	XSR3_NEON(x2, x1, rk[ 94], rk[ 95]);
	XSR9_NEON(x3, x2, rk[ 84], rk[ 85]);
	XSR5_NEON(x0, x3, rk[ 86], rk[ 87]);
	XSR3_NEON(x1, x0, rk[ 88], rk[ 89]);
	XSR9_NEON(x2, x1, rk[ 78], rk[ 79]);
	XSR5_NEON(x3, x2, rk[ 80], rk[ 81]);
	XSR3_NEON(x0, x3, rk[ 82], rk[ 83]);
	XSR9_NEON(x1, x0, rk[ 72], rk[ 73]);
	XSR5_NEON(x2, x1, rk[ 74], rk[ 75]);
	XSR3_NEON(x3, x2, rk[ 76], rk[ 77]);

	XSR9_NEON(x0, x3, rk[ 66], rk[ 67]);
	XSR5_NEON(x1, x0, rk[ 68], rk[ 69]);
	XSR3_NEON(x2, x1, rk[ 70], rk[ 71]);
	XSR9_NEON(x3, x2, rk[ 60], rk[ 61]);
	XSR5_NEON(x0, x3, rk[ 62], rk[ 63]);
	XSR3_NEON(x1, x0, rk[ 64], rk[ 65]);
	XSR9_NEON(x2, x1, rk[ 54], rk[ 55]);
	XSR5_NEON(x3, x2, rk[ 56], rk[ 57]);
	XSR3_NEON(x0, x3, rk[ 58], rk[ 59]);
	XSR9_NEON(x1, x0, rk[ 48], rk[ 49]);
	XSR5_NEON(x2, x1, rk[ 50], rk[ 51]);
	XSR3_NEON(x3, x2, rk[ 52], rk[ 53]);

	XSR9_NEON(x0, x3, rk[ 42], rk[ 43]);
	XSR5_NEON(x1, x0, rk[ 44], rk[ 45]);
	XSR3_NEON(x2, x1, rk[ 46], rk[ 47]);
	XSR9_NEON(x3, x2, rk[ 36], rk[ 37]);
	XSR5_NEON(x0, x3, rk[ 38], rk[ 39]);
	XSR3_NEON(x1, x0, rk[ 40], rk[ 41]);
	XSR9_NEON(x2, x1, rk[ 30], rk[ 31]);
	XSR5_NEON(x3, x2, rk[ 32], rk[ 33]);
	XSR3_NEON(x0, x3, rk[ 34], rk[ 35]);
	XSR9_NEON(x1, x0, rk[ 24], rk[ 25]);
	XSR5_NEON(x2, x1, rk[ 26], rk[ 27]);
	XSR3_NEON(x3, x2, rk[ 28], rk[ 29]);

	XSR9_NEON(x0, x3, rk[ 18], rk[ 19]);
	XSR5_NEON(x1, x0, rk[ 20], rk[ 21]);
	XSR3_NEON(x2, x1, rk[ 22], rk[ 23]);
	XSR9_NEON(x3, x2, rk[ 12], rk[ 13]);
	XSR5_NEON(x0, x3, rk[ 14], rk[ 15]);
	XSR3_NEON(x1, x0, rk[ 16], rk[ 17]);
	XSR9_NEON(x2, x1, rk[  6], rk[  7]);
	XSR5_NEON(x3, x2, rk[  8], rk[  9]);
	XSR3_NEON(x0, x3, rk[ 10], rk[ 11]);
	XSR9_NEON(x1, x0, rk[  0], rk[  1]);
	XSR5_NEON(x2, x1, rk[  2], rk[  3]);
	XSR3_NEON(x3, x2, rk[  4], rk[  5]);

	vst4q_u32((unsigned int *)pt, x);
}
