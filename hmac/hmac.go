// Package hmac implements the HMAC (Hash-based Message Authentication Code) as defined in KS X ISO/IEC 9797-2 and NIST SP 800-224.
package hmac

import (
	"crypto/subtle"
	"hash"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

const (
	// the standard tag size of HMAC, in bytes.
	StandardTagSize = 16
)

const (
	maxBlockSize = 256
	maxTagSize   = 64
)

// Equal reports whether mac1 and mac2 are identical, comparing them in constant time.
func Equal(mac1, mac2 []byte) bool {
	return subtle.ConstantTimeCompare(mac1, mac2) == 1
}

// [New] returns a new hash.Hash computing the HMAC checksum.
//
// The returned value is a [krypto.Hash] and therefore implements Zeroize(): call it
// to wipe the padded key and hash state when the MAC is no longer needed. After
// Zeroize() Sum panics and Write returns krypto.ErrStateCleared. NewWithSize is the
// same.
func New(h func() hash.Hash, key []byte) hash.Hash {
	return newHMAC(h, key, h().Size())
}

// [NewWithSize] returns a new hash.Hash computing the HMAC checksum
// with the specified tag size.
func NewWithSize(h func() hash.Hash, key []byte, tagSize int) krypto.Hash {
	return newHMAC(h, key, tagSize)
}

func newHMAC(h func() hash.Hash, key []byte, tagSize int) krypto.Hash {
	if tagSize <= 0 || tagSize > maxTagSize {
		panic(msgInvalidTagSize)
	}

	m := &mac{
		tagSize: tagSize,
		inner:   h(),
		outer:   h(),
	}

	unique := true
	func() {
		defer func() {
			// The comparison might panic if the underlying types are not comparable.
			_ = recover()
		}()
		if m.outer == m.inner {
			unique = false
		}
	}()
	if !unique {
		panic(msgHashGenerationNotUnique)
	}

	bs := m.inner.BlockSize()
	if bs <= 0 || bs > maxBlockSize {
		panic(msgUnsupportedHashCipher)
	}
	if tagSize > m.inner.Size() {
		panic(msgInvalidTagSize)
	}
	m.blockSize = bs

	if len(key) > bs {
		m.outer.Write(key)
		key = m.outer.Sum(m.sum[:0])
		defer clear(key)
	}
	copy(m.ipad[:], key)
	copy(m.opad[:], key)

	for i := range m.ipad {
		m.ipad[i] ^= 0x36
		m.opad[i] ^= 0x5c
	}

	m.z.SetValid()
	m.Reset()
	return m
}

type mac struct {
	z zeroize.ZeroizedFlag

	inner, outer hash.Hash

	ipad, opad [maxBlockSize]byte
	sum        [maxTagSize]byte
	blockSize  int
	tagSize    int
}

var _ krypto.Hash = (*mac)(nil)

func (m *mac) Zeroize() {
	m.z.Zeroize()

	zeroize.Zeroize(m.inner)
	zeroize.Zeroize(m.outer)
	m.inner, m.outer = nil, nil

	clear(m.ipad[:])
	clear(m.opad[:])
	clear(m.sum[:])
	m.blockSize = 0
	m.tagSize = 0
}

func (m *mac) Size() int      { return m.tagSize }
func (m *mac) BlockSize() int { return m.blockSize }

func (m *mac) Reset() {
	m.z.ThrowIfZeroized()

	m.inner.Reset()
	m.inner.Write(m.ipad[:m.blockSize])
}

func (m *mac) Write(p []byte) (int, error) {
	if err := m.z.ReturnIfZeroized(); err != nil {
		return 0, err
	}

	return m.inner.Write(p)
}

func (m *mac) Sum(dst []byte) []byte {
	m.z.ThrowIfZeroized()

	sum := m.inner.Sum(m.sum[:0])
	defer clear(sum)
	m.outer.Reset()
	m.outer.Write(m.opad[:m.blockSize])
	m.outer.Write(sum)

	return m.outer.Sum(dst)[:len(dst)+m.tagSize]
}
