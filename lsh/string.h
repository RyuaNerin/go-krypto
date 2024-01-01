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
