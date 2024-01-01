#include <emmintrin.h>
#include <immintrin.h>

typedef struct lea_key_st
{
	unsigned int rk[192];
	unsigned int round;
} LEA_KEY;

#define XAR3_AVX2(cur, pre, tmp, rk1, rk2)																						\
	tmp = _mm256_add_epi32(_mm256_xor_si256(pre, _mm256_set1_epi32(rk1)), _mm256_xor_si256(cur, _mm256_set1_epi32(rk2)));	\
	cur = _mm256_xor_si256(_mm256_srli_epi32(tmp, 3), _mm256_slli_epi32(tmp, 29));
#define XAR5_AVX2(cur, pre, tmp, rk1, rk2)																						\
	tmp = _mm256_add_epi32(_mm256_xor_si256(pre, _mm256_set1_epi32(rk1)), _mm256_xor_si256(cur, _mm256_set1_epi32(rk2)));	\
	cur = _mm256_xor_si256(_mm256_srli_epi32(tmp, 5), _mm256_slli_epi32(tmp, 27));
#define XAR9_AVX2(cur, pre, tmp, rk1, rk2)																						\
	tmp = _mm256_add_epi32(_mm256_xor_si256(pre, _mm256_set1_epi32(rk1)), _mm256_xor_si256(cur, _mm256_set1_epi32(rk2)));	\
	cur = _mm256_xor_si256(_mm256_srli_epi32(tmp, 23), _mm256_slli_epi32(tmp, 9));

#define XSR9_AVX2(cur, pre, rk1, rk2)																																		\
	cur = _mm256_xor_si256(_mm256_sub_epi32(_mm256_xor_si256(_mm256_srli_epi32(cur, 9), _mm256_slli_epi32(cur, 23)), _mm256_xor_si256(pre, _mm256_set1_epi32(rk1))),	\
		_mm256_set1_epi32(rk2));
#define XSR5_AVX2(cur, pre, rk1, rk2)																																		\
	cur = _mm256_xor_si256(_mm256_sub_epi32(_mm256_xor_si256(_mm256_srli_epi32(cur, 27), _mm256_slli_epi32(cur, 5)), _mm256_xor_si256(pre, _mm256_set1_epi32(rk1))),	\
		_mm256_set1_epi32(rk2));
#define XSR3_AVX2(cur, pre, rk1, rk2)																																		\
	cur = _mm256_xor_si256(_mm256_sub_epi32(_mm256_xor_si256(_mm256_srli_epi32(cur, 29), _mm256_slli_epi32(cur, 3)), _mm256_xor_si256(pre, _mm256_set1_epi32(rk1))),	\
		_mm256_set1_epi32(rk2));

void lea_encrypt_8block(const LEA_KEY *key, unsigned char *dst, const unsigned char *src)
{
	__m256i x0, x1, x2, x3, tmp;
	__m128i tmp128;

	x0 = _mm256_setr_epi32(*((unsigned int *)src), *((unsigned int *)src + 0x04), *((unsigned int *)src + 0x08), *((unsigned int *)src + 0x0c),
		*((unsigned int *)src + 0x10), *((unsigned int *)src + 0x14), *((unsigned int *)src + 0x18), *((unsigned int *)src + 0x1c));
	x1 = _mm256_setr_epi32(*((unsigned int *)src + 0x01), *((unsigned int *)src + 0x05), *((unsigned int *)src + 0x09), *((unsigned int *)src + 0x0d),
		*((unsigned int *)src + 0x11), *((unsigned int *)src + 0x15), *((unsigned int *)src + 0x19), *((unsigned int *)src + 0x1d));
	x2 = _mm256_setr_epi32(*((unsigned int *)src + 0x02), *((unsigned int *)src + 0x06), *((unsigned int *)src + 0x0a), *((unsigned int *)src + 0x0e),
		*((unsigned int *)src + 0x12), *((unsigned int *)src + 0x16), *((unsigned int *)src + 0x1a), *((unsigned int *)src + 0x1e));
	x3 = _mm256_setr_epi32(*((unsigned int *)src + 0x03), *((unsigned int *)src + 0x07), *((unsigned int *)src + 0x0b), *((unsigned int *)src + 0x0f),
		*((unsigned int *)src + 0x13), *((unsigned int *)src + 0x17), *((unsigned int *)src + 0x1b), *((unsigned int *)src + 0x1f));

	XAR3_AVX2(x3, x2, tmp, key->rk[4], key->rk[5]);
	XAR5_AVX2(x2, x1, tmp, key->rk[2], key->rk[3]);
	XAR9_AVX2(x1, x0, tmp, key->rk[0], key->rk[1]);
	XAR3_AVX2(x0, x3, tmp, key->rk[10], key->rk[11]);
	XAR5_AVX2(x3, x2, tmp, key->rk[8], key->rk[9]);
	XAR9_AVX2(x2, x1, tmp, key->rk[6], key->rk[7]);
	XAR3_AVX2(x1, x0, tmp, key->rk[16], key->rk[17]);
	XAR5_AVX2(x0, x3, tmp, key->rk[14], key->rk[15]);
	XAR9_AVX2(x3, x2, tmp, key->rk[12], key->rk[13]);
	XAR3_AVX2(x2, x1, tmp, key->rk[22], key->rk[23]);
	XAR5_AVX2(x1, x0, tmp, key->rk[20], key->rk[21]);
	XAR9_AVX2(x0, x3, tmp, key->rk[18], key->rk[19]);

	XAR3_AVX2(x3, x2, tmp, key->rk[28], key->rk[29]);
	XAR5_AVX2(x2, x1, tmp, key->rk[26], key->rk[27]);
	XAR9_AVX2(x1, x0, tmp, key->rk[24], key->rk[25]);
	XAR3_AVX2(x0, x3, tmp, key->rk[34], key->rk[35]);
	XAR5_AVX2(x3, x2, tmp, key->rk[32], key->rk[33]);
	XAR9_AVX2(x2, x1, tmp, key->rk[30], key->rk[31]);
	XAR3_AVX2(x1, x0, tmp, key->rk[40], key->rk[41]);
	XAR5_AVX2(x0, x3, tmp, key->rk[38], key->rk[39]);
	XAR9_AVX2(x3, x2, tmp, key->rk[36], key->rk[37]);
	XAR3_AVX2(x2, x1, tmp, key->rk[46], key->rk[47]);
	XAR5_AVX2(x1, x0, tmp, key->rk[44], key->rk[45]);
	XAR9_AVX2(x0, x3, tmp, key->rk[42], key->rk[43]);

	XAR3_AVX2(x3, x2, tmp, key->rk[52], key->rk[53]);
	XAR5_AVX2(x2, x1, tmp, key->rk[50], key->rk[51]);
	XAR9_AVX2(x1, x0, tmp, key->rk[48], key->rk[49]);
	XAR3_AVX2(x0, x3, tmp, key->rk[58], key->rk[59]);
	XAR5_AVX2(x3, x2, tmp, key->rk[56], key->rk[57]);
	XAR9_AVX2(x2, x1, tmp, key->rk[54], key->rk[55]);
	XAR3_AVX2(x1, x0, tmp, key->rk[64], key->rk[65]);
	XAR5_AVX2(x0, x3, tmp, key->rk[62], key->rk[63]);
	XAR9_AVX2(x3, x2, tmp, key->rk[60], key->rk[61]);
	XAR3_AVX2(x2, x1, tmp, key->rk[70], key->rk[71]);
	XAR5_AVX2(x1, x0, tmp, key->rk[68], key->rk[69]);
	XAR9_AVX2(x0, x3, tmp, key->rk[66], key->rk[67]);

	XAR3_AVX2(x3, x2, tmp, key->rk[76], key->rk[77]);
	XAR5_AVX2(x2, x1, tmp, key->rk[74], key->rk[75]);
	XAR9_AVX2(x1, x0, tmp, key->rk[72], key->rk[73]);
	XAR3_AVX2(x0, x3, tmp, key->rk[82], key->rk[83]);
	XAR5_AVX2(x3, x2, tmp, key->rk[80], key->rk[81]);
	XAR9_AVX2(x2, x1, tmp, key->rk[78], key->rk[79]);
	XAR3_AVX2(x1, x0, tmp, key->rk[88], key->rk[89]);
	XAR5_AVX2(x0, x3, tmp, key->rk[86], key->rk[87]);
	XAR9_AVX2(x3, x2, tmp, key->rk[84], key->rk[85]);
	XAR3_AVX2(x2, x1, tmp, key->rk[94], key->rk[95]);
	XAR5_AVX2(x1, x0, tmp, key->rk[92], key->rk[93]);
	XAR9_AVX2(x0, x3, tmp, key->rk[90], key->rk[91]);

	XAR3_AVX2(x3, x2, tmp, key->rk[100], key->rk[101]);
	XAR5_AVX2(x2, x1, tmp, key->rk[98], key->rk[99]);
	XAR9_AVX2(x1, x0, tmp, key->rk[96], key->rk[97]);
	XAR3_AVX2(x0, x3, tmp, key->rk[106], key->rk[107]);
	XAR5_AVX2(x3, x2, tmp, key->rk[104], key->rk[105]);
	XAR9_AVX2(x2, x1, tmp, key->rk[102], key->rk[103]);
	XAR3_AVX2(x1, x0, tmp, key->rk[112], key->rk[113]);
	XAR5_AVX2(x0, x3, tmp, key->rk[110], key->rk[111]);
	XAR9_AVX2(x3, x2, tmp, key->rk[108], key->rk[109]);
	XAR3_AVX2(x2, x1, tmp, key->rk[118], key->rk[119]);
	XAR5_AVX2(x1, x0, tmp, key->rk[116], key->rk[117]);
	XAR9_AVX2(x0, x3, tmp, key->rk[114], key->rk[115]);

	XAR3_AVX2(x3, x2, tmp, key->rk[124], key->rk[125]);
	XAR5_AVX2(x2, x1, tmp, key->rk[122], key->rk[123]);
	XAR9_AVX2(x1, x0, tmp, key->rk[120], key->rk[121]);
	XAR3_AVX2(x0, x3, tmp, key->rk[130], key->rk[131]);
	XAR5_AVX2(x3, x2, tmp, key->rk[128], key->rk[129]);
	XAR9_AVX2(x2, x1, tmp, key->rk[126], key->rk[127]);
	XAR3_AVX2(x1, x0, tmp, key->rk[136], key->rk[137]);
	XAR5_AVX2(x0, x3, tmp, key->rk[134], key->rk[135]);
	XAR9_AVX2(x3, x2, tmp, key->rk[132], key->rk[133]);
	XAR3_AVX2(x2, x1, tmp, key->rk[142], key->rk[143]);
	XAR5_AVX2(x1, x0, tmp, key->rk[140], key->rk[141]);
	XAR9_AVX2(x0, x3, tmp, key->rk[138], key->rk[139]);

	if (key->round > 24)
	{
		XAR3_AVX2(x3, x2, tmp, key->rk[148], key->rk[149]);
		XAR5_AVX2(x2, x1, tmp, key->rk[146], key->rk[147]);
		XAR9_AVX2(x1, x0, tmp, key->rk[144], key->rk[145]);
		XAR3_AVX2(x0, x3, tmp, key->rk[154], key->rk[155]);
		XAR5_AVX2(x3, x2, tmp, key->rk[152], key->rk[153]);
		XAR9_AVX2(x2, x1, tmp, key->rk[150], key->rk[151]);
		XAR3_AVX2(x1, x0, tmp, key->rk[160], key->rk[161]);
		XAR5_AVX2(x0, x3, tmp, key->rk[158], key->rk[159]);
		XAR9_AVX2(x3, x2, tmp, key->rk[156], key->rk[157]);
		XAR3_AVX2(x2, x1, tmp, key->rk[166], key->rk[167]);
		XAR5_AVX2(x1, x0, tmp, key->rk[164], key->rk[165]);
		XAR9_AVX2(x0, x3, tmp, key->rk[162], key->rk[163]);
	}

	if (key->round > 28)
	{
		XAR3_AVX2(x3, x2, tmp, key->rk[172], key->rk[173]);
		XAR5_AVX2(x2, x1, tmp, key->rk[170], key->rk[171]);
		XAR9_AVX2(x1, x0, tmp, key->rk[168], key->rk[169]);
		XAR3_AVX2(x0, x3, tmp, key->rk[178], key->rk[179]);
		XAR5_AVX2(x3, x2, tmp, key->rk[176], key->rk[177]);
		XAR9_AVX2(x2, x1, tmp, key->rk[174], key->rk[175]);
		XAR3_AVX2(x1, x0, tmp, key->rk[184], key->rk[185]);
		XAR5_AVX2(x0, x3, tmp, key->rk[182], key->rk[183]);
		XAR9_AVX2(x3, x2, tmp, key->rk[180], key->rk[181]);
		XAR3_AVX2(x2, x1, tmp, key->rk[190], key->rk[191]);
		XAR5_AVX2(x1, x0, tmp, key->rk[188], key->rk[189]);
		XAR9_AVX2(x0, x3, tmp, key->rk[186], key->rk[187]);
	}

	tmp128 = _mm256_extractf128_si256(x0, 0);
	*((unsigned int *)dst) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)dst + 0x04) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)dst + 0x08) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)dst + 0x0c) = _mm_extract_epi32(tmp128, 3);
	tmp128 = _mm256_extractf128_si256(x0, 1);
	*((unsigned int *)dst + 0x10) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)dst + 0x14) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)dst + 0x18) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)dst + 0x1c) = _mm_extract_epi32(tmp128, 3);

	tmp128 = _mm256_extractf128_si256(x1, 0);
	*((unsigned int *)dst + 0x01) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)dst + 0x05) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)dst + 0x09) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)dst + 0x0d) = _mm_extract_epi32(tmp128, 3);
	tmp128 = _mm256_extractf128_si256(x1, 1);
	*((unsigned int *)dst + 0x11) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)dst + 0x15) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)dst + 0x19) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)dst + 0x1d) = _mm_extract_epi32(tmp128, 3);

	tmp128 = _mm256_extractf128_si256(x2, 0);
	*((unsigned int *)dst + 0x02) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)dst + 0x06) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)dst + 0x0a) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)dst + 0x0e) = _mm_extract_epi32(tmp128, 3);
	tmp128 = _mm256_extractf128_si256(x2, 1);
	*((unsigned int *)dst + 0x12) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)dst + 0x16) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)dst + 0x1a) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)dst + 0x1e) = _mm_extract_epi32(tmp128, 3);

	tmp128 = _mm256_extractf128_si256(x3, 0);
	*((unsigned int *)dst + 0x03) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)dst + 0x07) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)dst + 0x0b) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)dst + 0x0f) = _mm_extract_epi32(tmp128, 3);
	tmp128 = _mm256_extractf128_si256(x3, 1);
	*((unsigned int *)dst + 0x13) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)dst + 0x17) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)dst + 0x1b) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)dst + 0x1f) = _mm_extract_epi32(tmp128, 3);
}

void lea_decrypt_8block(const LEA_KEY *key, unsigned char *dst, const unsigned char *src)
{
	__m256i x0, x1, x2, x3;
	__m128i tmp128;

	x0 = _mm256_setr_epi32(*((unsigned int *)src), *((unsigned int *)src + 0x04), *((unsigned int *)src + 0x08), *((unsigned int *)src + 0x0c),
		*((unsigned int *)src + 0x10), *((unsigned int *)src + 0x14), *((unsigned int *)src + 0x18), *((unsigned int *)src + 0x1c));
	x1 = _mm256_setr_epi32(*((unsigned int *)src + 0x01), *((unsigned int *)src + 0x05), *((unsigned int *)src + 0x09), *((unsigned int *)src + 0x0d),
		*((unsigned int *)src + 0x11), *((unsigned int *)src + 0x15), *((unsigned int *)src + 0x19), *((unsigned int *)src + 0x1d));
	x2 = _mm256_setr_epi32(*((unsigned int *)src + 0x02), *((unsigned int *)src + 0x06), *((unsigned int *)src + 0x0a), *((unsigned int *)src + 0x0e),
		*((unsigned int *)src + 0x12), *((unsigned int *)src + 0x16), *((unsigned int *)src + 0x1a), *((unsigned int *)src + 0x1e));
	x3 = _mm256_setr_epi32(*((unsigned int *)src + 0x03), *((unsigned int *)src + 0x07), *((unsigned int *)src + 0x0b), *((unsigned int *)src + 0x0f),
		*((unsigned int *)src + 0x13), *((unsigned int *)src + 0x17), *((unsigned int *)src + 0x1b), *((unsigned int *)src + 0x1f));

	if (key->round > 28)
	{
		XSR9_AVX2(x0, x3, key->rk[186], key->rk[187]);
		XSR5_AVX2(x1, x0, key->rk[188], key->rk[189]);
		XSR3_AVX2(x2, x1, key->rk[190], key->rk[191]);
		XSR9_AVX2(x3, x2, key->rk[180], key->rk[181]);
		XSR5_AVX2(x0, x3, key->rk[182], key->rk[183]);
		XSR3_AVX2(x1, x0, key->rk[184], key->rk[185]);
		XSR9_AVX2(x2, x1, key->rk[174], key->rk[175]);
		XSR5_AVX2(x3, x2, key->rk[176], key->rk[177]);
		XSR3_AVX2(x0, x3, key->rk[178], key->rk[179]);
		XSR9_AVX2(x1, x0, key->rk[168], key->rk[169]);
		XSR5_AVX2(x2, x1, key->rk[170], key->rk[171]);
		XSR3_AVX2(x3, x2, key->rk[172], key->rk[173]);
	}

	if (key->round > 24)
	{
		XSR9_AVX2(x0, x3, key->rk[162], key->rk[163]);
		XSR5_AVX2(x1, x0, key->rk[164], key->rk[165]);
		XSR3_AVX2(x2, x1, key->rk[166], key->rk[167]);
		XSR9_AVX2(x3, x2, key->rk[156], key->rk[157]);
		XSR5_AVX2(x0, x3, key->rk[158], key->rk[159]);
		XSR3_AVX2(x1, x0, key->rk[160], key->rk[161]);
		XSR9_AVX2(x2, x1, key->rk[150], key->rk[151]);
		XSR5_AVX2(x3, x2, key->rk[152], key->rk[153]);
		XSR3_AVX2(x0, x3, key->rk[154], key->rk[155]);
		XSR9_AVX2(x1, x0, key->rk[144], key->rk[145]);
		XSR5_AVX2(x2, x1, key->rk[146], key->rk[147]);
		XSR3_AVX2(x3, x2, key->rk[148], key->rk[149]);
	}

	XSR9_AVX2(x0, x3, key->rk[138], key->rk[139]);
	XSR5_AVX2(x1, x0, key->rk[140], key->rk[141]);
	XSR3_AVX2(x2, x1, key->rk[142], key->rk[143]);
	XSR9_AVX2(x3, x2, key->rk[132], key->rk[133]);
	XSR5_AVX2(x0, x3, key->rk[134], key->rk[135]);
	XSR3_AVX2(x1, x0, key->rk[136], key->rk[137]);
	XSR9_AVX2(x2, x1, key->rk[126], key->rk[127]);
	XSR5_AVX2(x3, x2, key->rk[128], key->rk[129]);
	XSR3_AVX2(x0, x3, key->rk[130], key->rk[131]);
	XSR9_AVX2(x1, x0, key->rk[120], key->rk[121]);
	XSR5_AVX2(x2, x1, key->rk[122], key->rk[123]);
	XSR3_AVX2(x3, x2, key->rk[124], key->rk[125]);

	XSR9_AVX2(x0, x3, key->rk[114], key->rk[115]);
	XSR5_AVX2(x1, x0, key->rk[116], key->rk[117]);
	XSR3_AVX2(x2, x1, key->rk[118], key->rk[119]);
	XSR9_AVX2(x3, x2, key->rk[108], key->rk[109]);
	XSR5_AVX2(x0, x3, key->rk[110], key->rk[111]);
	XSR3_AVX2(x1, x0, key->rk[112], key->rk[113]);
	XSR9_AVX2(x2, x1, key->rk[102], key->rk[103]);
	XSR5_AVX2(x3, x2, key->rk[104], key->rk[105]);
	XSR3_AVX2(x0, x3, key->rk[106], key->rk[107]);
	XSR9_AVX2(x1, x0, key->rk[96], key->rk[97]);
	XSR5_AVX2(x2, x1, key->rk[98], key->rk[99]);
	XSR3_AVX2(x3, x2, key->rk[100], key->rk[101]);

	XSR9_AVX2(x0, x3, key->rk[90], key->rk[91]);
	XSR5_AVX2(x1, x0, key->rk[92], key->rk[93]);
	XSR3_AVX2(x2, x1, key->rk[94], key->rk[95]);
	XSR9_AVX2(x3, x2, key->rk[84], key->rk[85]);
	XSR5_AVX2(x0, x3, key->rk[86], key->rk[87]);
	XSR3_AVX2(x1, x0, key->rk[88], key->rk[89]);
	XSR9_AVX2(x2, x1, key->rk[78], key->rk[79]);
	XSR5_AVX2(x3, x2, key->rk[80], key->rk[81]);
	XSR3_AVX2(x0, x3, key->rk[82], key->rk[83]);
	XSR9_AVX2(x1, x0, key->rk[72], key->rk[73]);
	XSR5_AVX2(x2, x1, key->rk[74], key->rk[75]);
	XSR3_AVX2(x3, x2, key->rk[76], key->rk[77]);

	XSR9_AVX2(x0, x3, key->rk[66], key->rk[67]);
	XSR5_AVX2(x1, x0, key->rk[68], key->rk[69]);
	XSR3_AVX2(x2, x1, key->rk[70], key->rk[71]);
	XSR9_AVX2(x3, x2, key->rk[60], key->rk[61]);
	XSR5_AVX2(x0, x3, key->rk[62], key->rk[63]);
	XSR3_AVX2(x1, x0, key->rk[64], key->rk[65]);
	XSR9_AVX2(x2, x1, key->rk[54], key->rk[55]);
	XSR5_AVX2(x3, x2, key->rk[56], key->rk[57]);
	XSR3_AVX2(x0, x3, key->rk[58], key->rk[59]);
	XSR9_AVX2(x1, x0, key->rk[48], key->rk[49]);
	XSR5_AVX2(x2, x1, key->rk[50], key->rk[51]);
	XSR3_AVX2(x3, x2, key->rk[52], key->rk[53]);

	XSR9_AVX2(x0, x3, key->rk[42], key->rk[43]);
	XSR5_AVX2(x1, x0, key->rk[44], key->rk[45]);
	XSR3_AVX2(x2, x1, key->rk[46], key->rk[47]);
	XSR9_AVX2(x3, x2, key->rk[36], key->rk[37]);
	XSR5_AVX2(x0, x3, key->rk[38], key->rk[39]);
	XSR3_AVX2(x1, x0, key->rk[40], key->rk[41]);
	XSR9_AVX2(x2, x1, key->rk[30], key->rk[31]);
	XSR5_AVX2(x3, x2, key->rk[32], key->rk[33]);
	XSR3_AVX2(x0, x3, key->rk[34], key->rk[35]);
	XSR9_AVX2(x1, x0, key->rk[24], key->rk[25]);
	XSR5_AVX2(x2, x1, key->rk[26], key->rk[27]);
	XSR3_AVX2(x3, x2, key->rk[28], key->rk[29]);

	XSR9_AVX2(x0, x3, key->rk[18], key->rk[19]);
	XSR5_AVX2(x1, x0, key->rk[20], key->rk[21]);
	XSR3_AVX2(x2, x1, key->rk[22], key->rk[23]);
	XSR9_AVX2(x3, x2, key->rk[12], key->rk[13]);
	XSR5_AVX2(x0, x3, key->rk[14], key->rk[15]);
	XSR3_AVX2(x1, x0, key->rk[16], key->rk[17]);
	XSR9_AVX2(x2, x1, key->rk[6], key->rk[7]);
	XSR5_AVX2(x3, x2, key->rk[8], key->rk[9]);
	XSR3_AVX2(x0, x3, key->rk[10], key->rk[11]);
	XSR9_AVX2(x1, x0, key->rk[0], key->rk[1]);
	XSR5_AVX2(x2, x1, key->rk[2], key->rk[3]);
	XSR3_AVX2(x3, x2, key->rk[4], key->rk[5]);


	tmp128 = _mm256_extractf128_si256(x0, 0);
	*((unsigned int *)dst) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)dst + 0x04) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)dst + 0x08) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)dst + 0x0c) = _mm_extract_epi32(tmp128, 3);
	tmp128 = _mm256_extractf128_si256(x0, 1);
	*((unsigned int *)dst + 0x10) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)dst + 0x14) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)dst + 0x18) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)dst + 0x1c) = _mm_extract_epi32(tmp128, 3);

	tmp128 = _mm256_extractf128_si256(x1, 0);
	*((unsigned int *)dst + 0x01) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)dst + 0x05) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)dst + 0x09) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)dst + 0x0d) = _mm_extract_epi32(tmp128, 3);
	tmp128 = _mm256_extractf128_si256(x1, 1);
	*((unsigned int *)dst + 0x11) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)dst + 0x15) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)dst + 0x19) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)dst + 0x1d) = _mm_extract_epi32(tmp128, 3);

	tmp128 = _mm256_extractf128_si256(x2, 0);
	*((unsigned int *)dst + 0x02) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)dst + 0x06) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)dst + 0x0a) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)dst + 0x0e) = _mm_extract_epi32(tmp128, 3);
	tmp128 = _mm256_extractf128_si256(x2, 1);
	*((unsigned int *)dst + 0x12) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)dst + 0x16) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)dst + 0x1a) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)dst + 0x1e) = _mm_extract_epi32(tmp128, 3);

	tmp128 = _mm256_extractf128_si256(x3, 0);
	*((unsigned int *)dst + 0x03) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)dst + 0x07) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)dst + 0x0b) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)dst + 0x0f) = _mm_extract_epi32(tmp128, 3);
	tmp128 = _mm256_extractf128_si256(x3, 1);
	*((unsigned int *)dst + 0x13) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)dst + 0x17) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)dst + 0x1b) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)dst + 0x1f) = _mm_extract_epi32(tmp128, 3);
}
