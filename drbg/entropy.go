package drbg

import (
	"bytes"
	"crypto/rand"
	"io"
	"sync"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

// Entropy supplies the randomness used to instantiate and reseed a DRBG.
type Entropy interface {
	// Get returns a fresh entropy input as a byte slice.
	//
	// The returned byte slice must be unique for each call and not reused across different DRBG instantiations.
	Get() ([]byte, error)
}

// NewEntropy returns an Entropy source that draws from rand.
func NewEntropy(rand io.Reader, length int) (Entropy, error) {
	if length < 0 {
		return nil, ErrEntropyLengthInvalid
	}
	return newEntropy(rand, length), nil
}

type entropy struct {
	z      zeroize.ZeroizedFlag
	m      sync.Mutex
	rand   io.Reader
	length int
}

var _ krypto.Zeroizable = (*entropy)(nil)

func newEntropy(rand io.Reader, length int) *entropy {
	return &entropy{rand: rand, length: length}
}

func (e *entropy) Zeroize() {
	e.z.SetInvalid()

	e.m.Lock()
	defer e.m.Unlock()

	e.rand = nil
	e.length = 0
}

func (e *entropy) Get() (p []byte, err error) {
	if err := e.z.ReturnIfZeroized(); err != nil {
		return nil, err
	}

	e.m.Lock()
	defer e.m.Unlock()

	if e.rand == nil {
		return nil, krypto.ErrStateCleared
	}

	p = make([]byte, e.length)
	_, err = io.ReadFull(e.rand, p)
	if err != nil {
		clear(p)
		return nil, err
	}

	return p, nil
}

// readEntropy gets an entropy input from src and checks its length.
func readEntropy(src Entropy, minLen int, maxLen uint64) ([]byte, error) {
	entropyInput, err := src.Get()
	if err != nil {
		return nil, err
	}

	// NewEntropy transfers ownership; custom sources retain their buffers.
	_, ownedBuffer := src.(*entropy)
	if err := checkEntropy(entropyInput, minLen, maxLen); err != nil {
		if ownedBuffer {
			clear(entropyInput)
		}
		return nil, err
	}
	if !ownedBuffer {
		entropyInput = bytes.Clone(entropyInput)
	}
	return entropyInput, nil
}

func instantiateNonce(args args, strength, entropyLen int) ([]byte, error) {
	nonceLen := (strength + 1) / 2
	if !args.nonceIsSet {
		nonce := make([]byte, nonceLen)
		if _, err := rand.Read(nonce); err != nil {
			clear(nonce)
			return nil, err
		}
		return nonce, nil
	}
	minNonceLen := max(0, nonceLen-(entropyLen-strength))
	if len(args.nonce) < minNonceLen {
		return nil, InvalidNonceError{Min: minNonceLen}
	}
	return args.nonce, nil
}

func checkEntropy(entropyInput []byte, minLen int, maxLen uint64) error {
	if len(entropyInput) < minLen || uint64(len(entropyInput)) > maxLen {
		return InvalidEntropyError{Min: uint64(minLen), Max: maxLen}
	}
	return nil
}

func maxHashSecurityStrength(hashOutLen int) int {
	switch {
	case hashOutLen <= 20: // SHA-1, HAS-160
		return 16
	case hashOutLen <= 28: // SHA-224, SHA-512/224, LSH-256-224, LSH-512-224
		return 24
	default:
		return 32
	}
}
