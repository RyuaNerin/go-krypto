package lsh256

import (
	"bytes"
	"crypto/rand"
	"encoding"
	"hash"
	"testing"
)

type h interface {
	hash.Hash
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

func TestMarshal(t *testing.T) {
	src := make([]byte, 4096)
	rand.Read(src)

	part1 := src[:2048]
	part2 := src[2048:]

	h := New().(h)
	h.Write(part1)
	data, _ := h.MarshalBinary()

	h.Write(part2)
	sum1 := h.Sum(nil)

	if err := h.UnmarshalBinary(data); err != nil {
		t.Errorf("UnmarshalBinary failed: %v", err)
		return
	}
	h.Write(part2)
	if !bytes.Equal(sum1, h.Sum(nil)) {
		t.Errorf("hash mismatch")
	}

	if err := h.UnmarshalBinary(data); err != nil {
		t.Errorf("UnmarshalBinary failed: %v", err)
		return
	}
	h.Write(part2)
	h.Write([]byte{0})
	if bytes.Equal(sum1, h.Sum(nil)) {
		t.Errorf("unexpected hash match")
	}

	if err := h.UnmarshalBinary(data[:len(data)-1]); err == nil {
		t.Errorf("unexpected success")
		return
	}

	data[0] = 'X'
	if err := h.UnmarshalBinary(data); err == nil {
		t.Errorf("unexpected success")
		return
	}
}
