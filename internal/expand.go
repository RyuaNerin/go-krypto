package internal

func Expand(arr []byte, bytes int) []byte {
	if bytes < cap(arr) {
		return arr[:bytes]
	} else {
		return make([]byte, bytes)
	}
}
