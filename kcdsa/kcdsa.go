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

type ParameterSizes int

const (
	A2048B224SHA224 ParameterSizes = kcdsainternal.A2048B224SHA224 // len(P) = 2048, len(Q) = 224, SHA-224, Recommended
	A2048B224SHA256 ParameterSizes = kcdsainternal.A2048B224SHA256 // len(P) = 2048, len(Q) = 256, SHA-256
	A2048B256SHA256 ParameterSizes = kcdsainternal.A2048B256SHA256 // len(P) = 2048, len(Q) = 256, SHA-256
	A3072B256SHA256 ParameterSizes = kcdsainternal.A3072B256SHA256 // len(P) = 3072, len(Q) = 256, SHA-256, Recommended
	A1024B160HAS160 ParameterSizes = kcdsainternal.A1024B160HAS160 // Deprecated: unsafe. lagacy use only
)

const (
	L2048N224SHA224 ParameterSizes = kcdsainternal.A2048B224SHA224 // Deprecated: use A2048B224SHA224
	L2048N224SHA256 ParameterSizes = kcdsainternal.A2048B224SHA256 // Deprecated: use A2048B224SHA256
	L2048N256SHA256 ParameterSizes = kcdsainternal.A2048B256SHA256 // Deprecated: use A2048B256SHA256
	L3072N256SHA256 ParameterSizes = kcdsainternal.A3072B256SHA256 // Deprecated: use A3072B256SHA256
)

func (ps ParameterSizes) Hash() hash.Hash {
	domain, ok := kcdsainternal.GetDomain(int(ps))
	if !ok {
		panic(msgInvalidParameterSizes)
	}

	return domain.NewHash()
}

// Generate the parameters. without GenParameters
func GenerateParameters(params *Parameters, rand io.Reader, sizes ParameterSizes) (err error) {
	d, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return errors.New(msgInvalidParameterSizes)
	}

	generated, err := kcdsainternal.GenerateParametersFast(rand, d)
	if err != nil {
		return err
	}

	params.P = generated.P
	params.Q = generated.Q
	params.G = generated.G

	return
}

// Generate the parameters. with GenParameters
func GenerateParametersTTAK(params *Parameters, rand io.Reader, sizes ParameterSizes) (err error) {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return errors.New(msgInvalidParameterSizes)
	}

	generated, err := kcdsainternal.GenerateParametersTTAK(rand, domain)
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

// TTAKParameters -> P, Q, G(randomly)
func RegenerateParameters(params *Parameters, rand io.Reader, sizes ParameterSizes) error {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return errors.New(msgInvalidParameterSizes)
	}

	if params.GenParameters.Count == 0 || len(params.GenParameters.Seed) == 0 {
		return errors.New(msgInvalidGenerationParameters)
	}
	if params.GenParameters.J == nil || params.GenParameters.J.Sign() <= 0 {
		return errors.New(msgInvalidGenerationParameters)
	}

	if len(params.GenParameters.Seed) != internal.BitsToBytes(domain.B) {
		return errors.New(msgInvalidGenerationParameters)
	}

	P, Q, ok := kcdsainternal.RegeneratePQ(
		domain,
		params.GenParameters.J,
		params.GenParameters.Seed,
		params.GenParameters.Count,
	)
	if !ok {
		return errors.New(msgInvalidGenerationParameters)
	}

	H, G := new(big.Int), new(big.Int)
	_, err := kcdsainternal.GenerateHG(H, G, nil, rand, P, params.GenParameters.J)
	if err != nil {
		return err
	}

	params.P = P
	params.Q = Q
	params.G = G

	return nil
}

func GenerateKey(priv *PrivateKey, rand io.Reader) error {
	if priv.P == nil || priv.Q == nil || priv.G == nil {
		return errors.New(msgErrorParametersNotSetUp)
	}

	X, Y := new(big.Int), new(big.Int)

	XBytes := make([]byte, internal.BitsToBytes(priv.Q.BitLen()))

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
	kcdsainternal.GenerateY(Y, priv.P, priv.Q, priv.G, X)

	priv.Y = Y
	priv.X = X

	return nil
}

func GenerateKeyWithSeed(priv *PrivateKey, rand io.Reader, xkey, upri []byte, sizes ParameterSizes) (xkeyOut, upriOut []byte, err error) {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return nil, nil, errors.New(msgInvalidParameterSizes)
	}

	if priv.P == nil || priv.Q == nil || priv.G == nil {
		return nil, nil, errors.New(msgErrorParametersNotSetUp)
	}

	if len(xkey) == 0 {
		xkey, err = internal.ReadBits(nil, rand, domain.B)
		if err != nil {
			return nil, nil, err
		}
	} else if len(xkey) < internal.BitsToBytes(domain.B) {
		return nil, nil, errors.New(msgErrorShortXKEY)
	}
	if len(upri) == 0 {
		upri, err = internal.ReadBytes(nil, rand, 64)
		if err != nil {
			return nil, nil, err
		}
	}

	h := domain.NewHash()

	priv.X, priv.Y = new(big.Int), new(big.Int)
	kcdsainternal.GenerateX(priv.X, priv.Q, upri, xkey, h, domain)
	kcdsainternal.GenerateY(priv.Y, priv.P, priv.Q, priv.G, priv.X)

	return xkey, upri, nil
}

func Sign(rand io.Reader, priv *PrivateKey, sizes ParameterSizes, data []byte) (r, s *big.Int, err error) {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return nil, nil, errors.New(msgInvalidParameterSizes)
	}

	randutil.MaybeReadByte(rand)

	if priv.Q.Sign() <= 0 || priv.P.Sign() <= 0 || priv.G.Sign() <= 0 || priv.X.Sign() <= 0 || priv.Q.BitLen()%8 != 0 {
		return nil, nil, errors.New(msgInvalidPublicKey)
	}

	r, s = new(big.Int), new(big.Int)

	qblen := priv.Q.BitLen()

	K := new(big.Int)
	buf := make([]byte, internal.BitsToBytes(qblen))

	tmpInt := new(big.Int)
	var tmpBuf []byte

	var attempts int
	for attempts = 10; attempts > 0; attempts-- {
		// step 1. 난수 k를 [1, Q-1]에서 임의로 선택한다.

		for {
			buf, err = internal.ReadBits(buf, rand, qblen)
			if err != nil {
				return
			}
			K.SetBytes(buf)
			K.Add(K, internal.One)

			if K.Sign() > 0 && K.Cmp(priv.Q) < 0 {
				break
			}
		}

		tmpBuf, ok = kcdsainternal.Sign(
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
		return nil, nil, errors.New(msgInvalidPublicKey)
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
