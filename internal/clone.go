package internal

// Clone returns a copy of b[:len(b)].
// The result may have additional unused capacity.
// Clone(nil) returns nil.
func BytesClone(b []byte) []byte {
	// bytes.Clone (go1.20)
	if b == nil {
		return nil
	}
	return append([]byte{}, b...)
}
