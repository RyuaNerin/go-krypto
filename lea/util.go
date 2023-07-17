package lea

func bitsRotateRight32(W, i uint32) uint32 {
	return (((W) >> (i)) | ((W) << (32 - (i))))
}
