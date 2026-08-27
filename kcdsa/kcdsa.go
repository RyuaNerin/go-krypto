// Package kcdsa implements the KCDSA(Korean Certificate-based Digital Signature Algorithm) as defined in TTAK.KO-12.0001/R4
package kcdsa

import (
	"hash"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/kcdsa"
	"github.com/RyuaNerin/go-krypto/internal/randutil"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

// ParameterSizes identifies a supported KCDSA parameter set.
type ParameterSizes int

const (
	// len(P) = 2048, len(Q) = 224, SHA-224, Recommended
	//
	// import "crypto/sha256" before using this parameter size.
	A2048B224SHA224 ParameterSizes = kcdsa.A2048B224SHA224
	// len(P) = 2048, len(Q) = 224, SHA-256
	//
	// import "crypto/sha256" before using this parameter size.
	A2048B224SHA256 ParameterSizes = kcdsa.A2048B224SHA256
	// len(P) = 2048, len(Q) = 256, SHA-256
	//
	// import "crypto/sha256" before using this parameter size.
	A2048B256SHA256 ParameterSizes = kcdsa.A2048B256SHA256
	// len(P) = 3072, len(Q) = 256, SHA-256, Recommended
	//
	// import "crypto/sha256" before using this parameter size.
	A3072B256SHA256 ParameterSizes = kcdsa.A3072B256SHA256
	// Deprecated: unsafe. legacy use only
	//
	// import "github.com/RyuaNerin/go-krypto/has160" before using this parameter size.
	A1024B160HAS160 ParameterSizes = kcdsa.A1024B160HAS160
)

const (
	L2048N224SHA224 ParameterSizes = kcdsa.A2048B224SHA224 // Deprecated: use A2048B224SHA224
	L2048N224SHA256 ParameterSizes = kcdsa.A2048B224SHA256 // Deprecated: use A2048B224SHA256
	L2048N256SHA256 ParameterSizes = kcdsa.A2048B256SHA256 // Deprecated: use A2048B256SHA256
	L3072N256SHA256 ParameterSizes = kcdsa.A3072B256SHA256 // Deprecated: use A3072B256SHA256
)

func (ps ParameterSizes) Hash() hash.Hash {
	domain, ok := kcdsa.GetDomain(int(ps))
	if !ok {
		panic(msgInvalidParameterSizes)
	}

	return domain.NewHash()
}

// GenerateParameters generates KCDSA parameters in a manner similar to crypto/dsa.
//
// The caller must provide a cryptographically secure random source.
//
// It does not generate or use the J, Seed, Count values specified in the TTAK document.
func GenerateParameters(params *Parameters, rand io.Reader, sizes ParameterSizes) (err error) {
	d, ok := kcdsa.GetDomain(int(sizes))
	if !ok {
		return ErrInvalidParameterSizes
	}

	generated, err := kcdsa.GenerateParametersFast(rand, d)
	if err != nil {
		return err
	}

	params.P = generated.P
	params.Q = generated.Q
	params.G = generated.G

	return
}

// GenerateParametersTTAK generates parameters as specified in the TTAK document.
//
// The caller must provide a cryptographically secure random source.
//
// It generates [GenerationParameters]
func GenerateParametersTTAK(params *Parameters, rand io.Reader, sizes ParameterSizes) (err error) {
	domain, ok := kcdsa.GetDomain(int(sizes))
	if !ok {
		return ErrInvalidParameterSizes
	}

	generated, err := kcdsa.GenerateParametersTTAK(rand, domain)
	if err != nil {
		return err
	}

	params.P = generated.P
	params.Q = generated.Q
	params.G = generated.G
	params.GenParameters.J = generated.J
	params.GenParameters.Seed = generated.Seed
	params.GenParameters.Count = generated.Count
	return
}

// RegenerateParameters regenerates parameters (P, Q, G)
// from the stored generation values [GenerationParameters]
//
//   - P, Q : derived deterministically from Seed and Count, so they are always reproduced identically.
//   - G :    computed from rand, so it may differ on each call.
//
// The rand parameter must be a cryptographically secure source of randomness.
func RegenerateParameters(params *Parameters, rand io.Reader, sizes ParameterSizes) error {
	domain, ok := kcdsa.GetDomain(int(sizes))
	if !ok {
		return ErrInvalidParameterSizes
	}

	if params == nil || params.GenParameters.Count <= 0 || params.GenParameters.Count > 1<<24 || len(params.GenParameters.Seed) == 0 {
		return ErrInvalidGenerationParameters
	}
	if params.GenParameters.J == nil || params.GenParameters.J.Sign() <= 0 {
		return ErrInvalidGenerationParameters
	}

	if len(params.GenParameters.Seed) != internal.BitsToBytes(domain.B) {
		return ErrInvalidGenerationParameters
	}

	P, Q := new(big.Int), new(big.Int)
	ok = kcdsa.RegeneratePQ(
		P, Q,
		domain,
		params.GenParameters.J,
		params.GenParameters.Seed,
		params.GenParameters.Count,
	)
	if !ok {
		zeroize.Zeroize(P)
		zeroize.Zeroize(Q)
		return ErrInvalidGenerationParameters
	}

	H, G := new(big.Int), new(big.Int)
	defer zeroize.Zeroize(H)
	_, err := kcdsa.GenerateHG(H, G, nil, rand, P, params.GenParameters.J)
	if err != nil {
		zeroize.Zeroize(P)
		zeroize.Zeroize(Q)
		zeroize.Zeroize(G)
		return err
	}

	params.P = P
	params.Q = Q
	params.G = G

	return nil
}

// GenerateKey generates a KCDSA private key using the provided random source.
//
// The rand parameter must be a cryptographically secure source of randomness.
// Pass crypto/rand.Reader unless a specific source is required.
func GenerateKey(priv *PrivateKey, rand io.Reader) error {
	if priv.P == nil || priv.Q == nil || priv.G == nil {
		return ErrErrorParametersNotSetUp
	}

	X, Y := new(big.Int), new(big.Int)

	XBytes := make([]byte, internal.BitsToBytes(priv.Q.BitLen()))
	defer clear(XBytes)

	for {
		_, err := io.ReadFull(rand, XBytes)
		if err != nil {
			return err
		}
		X.SetBytes(XBytes)
		if X.Sign() > 0 && X.Cmp(priv.Q) < 0 {
			break
		}
	}
	kcdsa.GenerateY(Y, priv.P, priv.Q, priv.G, X)

	priv.Y = Y
	priv.X = X

	return nil
}

// GenerateKeyWithSeed generates a KCDSA private key using the provided seed, along
// with the optional xkey and upri values.
//
// The caller is responsible for the entropy of rand and any supplied xkey and
// upri. Hashing or expanding predictable inputs does not add entropy. Prefer
// GenerateKey when seed-based generation is not required.
//
// If xkey is nil, a random xkey is generated from B bits read from the provided
// rand reader. If len(xkey) is less than B bits, an error is returned.
//
// If upri is nil, a random upri is generated from 64 bytes read from the provided
// rand reader.
//
// The resulting *PrivateKey is a [krypto.PrivateKey] and therefore implements
// Zeroize(); call it to erase the secret X when the key is no longer needed.
func GenerateKeyWithSeed(
	priv *PrivateKey,
	rand io.Reader,
	xkey, upri []byte,
	sizes ParameterSizes,
) (xkeyOut, upriOut []byte, err error) {
	domain, ok := kcdsa.GetDomain(int(sizes))
	if !ok {
		return nil, nil, ErrInvalidParameterSizes
	}

	if priv.P == nil || priv.Q == nil || priv.G == nil {
		return nil, nil, ErrErrorParametersNotSetUp
	}

	if len(xkey) == 0 {
		xkey, err = internal.ReadBits(nil, rand, domain.B)
		if err != nil {
			return nil, nil, err
		}
	} else if len(xkey) < internal.BitsToBytes(domain.B) {
		return nil, nil, ErrErrorShortXKEY
	}
	if len(upri) == 0 {
		upri, err = internal.ReadBytes(nil, rand, 64)
		if err != nil {
			return nil, nil, err
		}
	}

	h := domain.NewHash()
	defer zeroize.Zeroize(h)

	priv.X, priv.Y = new(big.Int), new(big.Int)
	kcdsa.GenerateX(priv.X, priv.Q, upri, xkey, h, domain)
	kcdsa.GenerateY(priv.Y, priv.P, priv.Q, priv.G, priv.X)

	return xkey, upri, nil
}

func Sign(rand io.Reader, priv *PrivateKey, sizes ParameterSizes, data []byte) (r, s *big.Int, err error) {
	domain, ok := kcdsa.GetDomain(int(sizes))
	if !ok {
		return nil, nil, ErrInvalidParameterSizes
	}

	randutil.MaybeReadByte(rand)

	if priv == nil || !kcdsa.IsValidPublicKey(priv.P, priv.Q, priv.G, priv.Y, domain) {
		return nil, nil, ErrInvalidPublicKey
	}
	if priv.X == nil || priv.X.Sign() <= 0 || priv.X.Cmp(priv.Q) >= 0 {
		return nil, nil, ErrInvalidPrivateKey
	}

	// Hedge K: mix the private key and the message into rand, as eckcdsa does,
	// so that a weak or repeating rand does not reveal X.
	h := domain.NewHash()
	h.Write(data)
	digest := h.Sum(nil)
	defer zeroize.Zeroize(digest)

	xBytes := priv.X.Bytes()
	defer zeroize.Zeroize(xBytes)
	csprng, err := randutil.MixedCSPRNG(rand, xBytes, digest)
	if err != nil {
		return nil, nil, err
	}
	defer zeroize.Zeroize(csprng)

	r, s = new(big.Int), new(big.Int)

	qblen := priv.Q.BitLen()

	K := new(big.Int)
	buf := make([]byte, internal.BitsToBytes(qblen))

	tmpInt := new(big.Int)
	var tmpBuf []byte

	defer func() {
		zeroize.Zeroize(K)
		clear(buf)
		zeroize.Zeroize(tmpInt)
		clear(tmpBuf)
	}()

	var attempts int
	for attempts = 10; attempts > 0; attempts-- {
		// step 1. 난수 k를 [1, Q-1]에서 임의로 선택한다.

		for {
			buf, err = internal.ReadBits(buf, csprng, qblen)
			if err != nil {
				return nil, nil, err
			}
			K.SetBytes(buf)
			K.Add(K, internal.One)

			if K.Sign() > 0 && K.Cmp(priv.Q) < 0 {
				break
			}
		}

		tmpBuf, ok = kcdsa.Sign(
			r, s,
			priv.P, priv.Q, priv.G, priv.Y, priv.X,
			domain,
			K, data,
			tmpInt,
			tmpBuf,
		)
		if ok {
			break
		}
	}

	// Only degenerate private keys will require more than a handful of
	// attempts.
	if attempts == 0 {
		return nil, nil, ErrInvalidPrivateKey
	}

	return r, s, nil
}

// Verify checks a signature using a trusted domain. For directly constructed
// keys received from an untrusted source, call pub.Validate before Verify.
// Keys imported by kx509 are validated during parsing.
func Verify(pub *PublicKey, sizes ParameterSizes, data []byte, R, S *big.Int) bool {
	domain, ok := kcdsa.GetDomain(int(sizes))
	if !ok {
		return false
	}

	// A degenerate key such as Y = 0 or G = P makes W' constant, so that one
	// signature would verify for every message.
	if pub == nil || !kcdsa.IsValidPublicKey(pub.P, pub.Q, pub.G, pub.Y, domain) {
		return false
	}

	// step 1. 수신된 서명 {R', S'}에 대해 |R'|=LH, 0 < S' < Q 임을 확인한다.
	if R == nil || S == nil {
		return false
	}
	if R.Sign() < 1 {
		return false
	}
	if R.BitLen() > domain.B {
		return false
	}
	if S.Sign() < 1 || S.Cmp(pub.Q) >= 0 {
		return false
	}

	return kcdsa.Verify(pub.P, pub.Q, pub.G, pub.Y, domain, data, R, S)
}
