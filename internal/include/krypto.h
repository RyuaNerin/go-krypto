#pragma once
#if !defined(__KRYPTO_H__)
#define __KRYPTO_H__

#include <stdint.h>

typedef uint32_t krypto_u32 __attribute__((aligned(1), may_alias));

static inline uint32_t krypto_load_u32(const void* ptr)
{
    return *(const krypto_u32*)ptr;
}

static inline void krypto_store_u32(void* ptr, uint32_t value)
{
    *(krypto_u32*)ptr = value;
}

#define GO_VARIABLE static const volatile


static inline uint32_t krypto_const_u32(const volatile krypto_u32* word)
{
    __asm__ volatile("" : "+r"(word) : : "memory");
    return *word;
}

#define KRYPTO_TABLE_U32(TABLE, INDEX) krypto_const_u32(&(TABLE)[(INDEX)])

static inline void* memcpy(unsigned char* dst, const unsigned char* src, int len)
{
	for (int i = 0; i < len; i++) {
		dst[i] = src[i];
	}
	return dst;
}

static inline void* memset(unsigned char* dst, const unsigned char c, int len)
{
	for (int i = 0; i < len; i++) {
		dst[i] = c;
	}
	return dst;
}

#endif
