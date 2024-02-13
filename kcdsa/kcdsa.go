// Package kcdsa implements the KCDSA(Korean Certificate-based Digital Signature Algorithm) as defined in TTAK.KO-12.0001/R4
package kcdsa

import (
	"errors"
	"hash"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
	kcdsainternal "github.com/RyuaNerin/go-krypto/internal/kcdsa"
	"github.com/RyuaNerin/go-krypto/internal/randutil"
)

var (
	ErrInvalidPublicKey       = errors.New("krypto/kcdsa: invalid public key")
	ErrInvalidTTAKParameters  = errors.New("krypto/kcdsa: invalid domain parameters")
	ErrInvalidParameterSizes  = errors.New("krypto/kcdsa: invalid ParameterSizes")
	ErrParametersNotSetUp     = errors.New("krypto/kcdsa: parameters not set up before generating key")
	ErrTTAKParametersNotSetUp = errors.New("krypto/kcdsa: ttakparameters not set up before generating key")
	ErrShortXKEY              = errors.New("krypto/kcdsa: XKEY is too small.")
)

type ParameterSizes int

const (
	L2048N224SHA224 ParameterSizes = kcdsainternal.L2048N224SHA224
	L2048N224SHA256 ParameterSizes = kcdsainternal.L2048N224SHA256
	L2048N256SHA256 ParameterSizes = kcdsainternal.L2048N256SHA256
	L3072N256SHA256 ParameterSizes = kcdsainternal.L3072N256SHA256
)

func (ps ParameterSizes) Hash() hash.Hash {
	domain, ok := kcdsainternal.GetDomain(int(ps))
	if !ok {
		panic(ErrInvalidParameterSizes.Error())
	}
	return domain.NewHash()
}

// Generate the paramters
func GenerateParameters(params *Parameters, rand io.Reader, sizes ParameterSizes) (err error) {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return ErrInvalidParameterSizes
	}

	generated, err := kcdsainternal.GenerateParameters(rand, domain)
	if err != nil {
		return err
	}

	params.P = generated.P
	params.Q = generated.Q
	params.G = generated.G
	params.TTAKParams.J = generated.J
	params.TTAKParams.Seed = generated.Seed
	params.TTAKParams.Count = generated.Count
	return
}

// TTAKParameters -> P, Q, G(randomly)
func RegenerateParameters(params *Parameters, rand io.Reader, sizes ParameterSizes) error {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return ErrInvalidParameterSizes
	}

	if params.TTAKParams.Count == 0 || params.TTAKParams.J == nil || params.TTAKParams.Seed == nil || params.TTAKParams.J.Sign() <= 0 {
		return ErrInvalidTTAKParameters
	}
	if params.TTAKParams.J.Sign() <= 0 {
		return ErrInvalidTTAKParameters
	}

	if len(params.TTAKParams.Seed) != internal.Bytes(domain.B) {
		return ErrInvalidTTAKParameters
	}

	P, Q, G, err := kcdsainternal.RegenerateParameters(
		rand,
		domain,
		params.TTAKParams.J,
		params.TTAKParams.Seed,
		params.TTAKParams.Count,
	)
	if err == kcdsainternal.ErrInvalidTTAKParameters {
		return ErrInvalidTTAKParameters
	}

	params.P = P
	params.Q = Q
	params.G = G

	return nil
}

func GenerateKey(priv *PrivateKey, rand io.Reader) error {
	if priv.P == nil || priv.Q == nil || priv.G == nil {
		return ErrParametersNotSetUp
	}

	X := new(big.Int)
	XBytes := make([]byte, priv.Q.BitLen()/8)

	for {
		_, err := io.ReadFull(rand, XBytes)
		if err != nil {
			return err
		}
		X.SetBytes(XBytes)
		if X.Sign() != 0 && X.Cmp(priv.Q) < 0 {
			break
		}
	}

	priv.X = X
	priv.Y = kcdsainternal.GenerateY(priv.P, priv.Q, priv.G, priv.X)

	return nil
}

func GenerateKeyWithSeed(priv *PrivateKey, rand io.Reader, xkey, upri []byte, sizes ParameterSizes) (xkeyOut, upriOut []byte, err error) {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return nil, nil, ErrInvalidParameterSizes
	}

	if priv.P == nil || priv.Q == nil || priv.G == nil {
		return nil, nil, ErrParametersNotSetUp
	}

	if len(xkey) == 0 {
		xkey, err = internal.ReadBits(rand, xkey, domain.B)
		if err != nil {
			return nil, nil, err
		}
	} else if len(xkey) < internal.Bytes(domain.B) {
		return nil, nil, ErrShortXKEY
	}
	if len(upri) == 0 {
		upri, err = internal.ReadBytes(rand, upri, 64)
		if err != nil {
			return nil, nil, err
		}
	}

	h := domain.NewHash()

	priv.X = kcdsainternal.GenerateX(priv.Q, upri, xkey, h, domain)
	priv.Y = kcdsainternal.GenerateY(priv.P, priv.Q, priv.G, priv.X)

	return xkey, upri, nil
}

func Sign(rand io.Reader, priv *PrivateKey, sizes ParameterSizes, data []byte) (r, s *big.Int, err error) {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return nil, nil, ErrInvalidParameterSizes
	}

	randutil.MaybeReadByte(rand)

	if priv.Q.Sign() <= 0 || priv.P.Sign() <= 0 || priv.G.Sign() <= 0 || priv.X.Sign() <= 0 || priv.Q.BitLen()%8 != 0 {
		err = ErrInvalidPublicKey
		return
	}

	qblen := priv.Q.BitLen()

	K := new(big.Int)
	buf := make([]byte, internal.Bytes(qblen))

	var attempts int
	for attempts = 10; attempts > 0; attempts-- {
		// step 1. 난수 k를 [1, Q-1]에서 임의로 선택한다.

		for {
			buf, err = internal.ReadBits(rand, buf, qblen)
			if err != nil {
				return
			}
			K.SetBytes(buf)
			K.Add(K, internal.One)

			if K.Sign() > 0 && K.Cmp(priv.Q) < 0 {
				break
			}
		}

		r, s, ok = kcdsainternal.Sign(
			priv.P, priv.Q, priv.G, priv.Y, priv.X,
			domain,
			K, data,
		)
		if ok {
			break
		}
	}

	// Only degenerate private keys will require more than a handful of
	// attempts.
	if attempts == 0 {
		return nil, nil, ErrInvalidPublicKey
	}

	return
}

func Verify(pub *PublicKey, sizes ParameterSizes, data []byte, R, S *big.Int) bool {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return false
	}

	// step 1. 수신된 서명 {R', S'}에 대해 |R'|=LH, 0 < S' < Q 임을 확인한다.
	if pub.P.Sign() <= 0 {
		return false
	}

	if R.Sign() < 1 {
		return false
	}
	if S.Sign() < 1 || S.Cmp(pub.Q) >= 0 {
		return false
	}

	return kcdsainternal.Verify(pub.P, pub.Q, pub.G, pub.Y, domain, data, R, S)
}
