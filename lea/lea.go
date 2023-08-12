package lea

type leaContext struct {
	round uint8
	rk    [192]uint32
	ecb   bool
}
