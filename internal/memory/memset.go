package memory

func Memset[S ~[]E, E comparable](a S, v E) {
	if len(a) == 0 {
		return
	}
	var zero E
	if v == zero {
		clear(a)
		return
	}

	a[0] = v
	for bp := 1; bp < len(a); bp *= 2 {
		copy(a[bp:], a[:bp])
	}
}
