package krypto

import (
	"bytes"
	"crypto/rand"
	"io/fs"
	"math/big"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/RyuaNerin/go-krypto/kcdsa"
	"github.com/RyuaNerin/go-krypto/test"
)

func Test_KCDSA(t *testing.T) {
	testDSA(
		t,
		`KCDSA_\((?P<l>[^\)]+)\)\((?P<n>[^\)]+)\)\((?P<h>[^\)]+)\)_[^\.]+\.txt`,
		func(p, q, g, x, y *big.Int, m []byte, l, n, h string) (r *big.Int, s *big.Int, err error) {
			priv := kcdsa.PrivateKey{
				X: x,
				PublicKey: kcdsa.PublicKey{
					Y: y,
					Parameters: kcdsa.Parameters{
						P: p,
						Q: q,
						G: g,
					},
				},
			}

			switch {
			case l == "2048" && n == "224" && h == "SHA-224":
				priv.Sizes = kcdsa.L2048N224SHA224
			case l == "2048" && n == "224" && h == "SHA-256":
				priv.Sizes = kcdsa.L2048N224SHA256
			case l == "2048" && n == "256" && h == "SHA-256":
				priv.Sizes = kcdsa.L2048N256SHA256
			case l == "3072" && n == "256" && h == "SHA-256":
				priv.Sizes = kcdsa.L3072N256SHA256
			default:
				return nil, nil, nil
			}

			return kcdsa.Sign(rand.Reader, &priv, bytes.NewReader(m))
		},
		func(p, q, g, y *big.Int, m []byte, r, s *big.Int, l, n, h string) (ok bool, err error) {
			pub := kcdsa.PublicKey{
				Y: y,
				Parameters: kcdsa.Parameters{
					P: p,
					Q: q,
					G: g,
				},
			}

			switch {
			case l == "2048" && n == "224" && h == "SHA-224":
				pub.Sizes = kcdsa.L2048N224SHA224
			case l == "2048" && n == "224" && h == "SHA-256":
				pub.Sizes = kcdsa.L2048N224SHA256
			case l == "2048" && n == "256" && h == "SHA-256":
				pub.Sizes = kcdsa.L2048N256SHA256
			case l == "3072" && n == "256" && h == "SHA-256":
				pub.Sizes = kcdsa.L3072N256SHA256
			default:
				return false, nil
			}

			return kcdsa.Verify(&pub, bytes.NewReader(m), r, s)
		},
	)
}

func testDSA(
	t *testing.T,
	regexStr string,
	sign func(p, q, g, x, y *big.Int, m []byte, l, n, h string) (r *big.Int, s *big.Int, err error),
	verify func(p, q, g, y *big.Int, m []byte, r, s *big.Int, l, n, h string) (ok bool, err error),
) {
	re := regexp.MustCompile(regexStr)

	filepath.Walk(
		"test",
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			var l, n, h string

			m := re.FindStringSubmatch(filepath.Base(path))
			if m == nil {
				return nil
			}
			for i, name := range re.SubexpNames() {
				switch name {
				case "l":
					l = m[i]
				case "n":
					n = m[i]
				case "h":
					h = m[i]
				}
			}

			test.DSATest(
				t,
				path,
				func(p, q, g, x, y *big.Int, m []byte) (r *big.Int, s *big.Int, err error) {
					return sign(p, q, g, x, y, m, l, n, h)
				},
				func(p, q, g, y *big.Int, m []byte, r, s *big.Int) (ok bool, err error) {
					return verify(p, q, g, y, m, r, s, l, n, h)
				},
			)

			return nil
		},
	)
}
