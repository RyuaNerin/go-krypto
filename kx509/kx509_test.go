package kx509

import (
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"log"
	"testing"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/RyuaNerin/go-krypto/kcdsa"
)

var (
	curveList = []elliptic.Curve{
		elliptic.P256(),
		elliptic.P224(),
		elliptic.P384(),
		elliptic.P521(),
	}
	sizeList = []kcdsa.ParameterSizes{
		kcdsa.L2048N224SHA256,

		kcdsa.L2048N224SHA224,
		kcdsa.L2048N256SHA256,
		kcdsa.L3072N256SHA256,
	}

	testCases []struct {
		PrivateKey interface{}
		Marshaled  []byte
	}
)

func TestDebug(t *testing.T) {
	b, _ := hex.DecodeString("308202650201003082023a06082a831a8c9a4401013082022c0282010100c837602569203887f994fccb066a0229a60eea6ec380d525aa5839debe75e63d385168a81767068978aa740de9db3a434a07710d5f3a50a412f990d3e3a4e439a51a94c6b6c5b228bfc1fc12ad99e85e951910bd0cb75dc69a8c74de4c16339aa68c4b3a7d002f470fba7fcdd1d5ccced89bd4c1fef153bc2097d58c8303547017c2e739059ed2f484dd6e42d295508c7853f653c72e805958a42b02affb2799539d9a2f6faa4988ed1c5f293a81601307476c9a8c7d09b50f8584394af3de1a63469f0a106256227c537df45b14c2b5ac66168d35b16c460093e6556c59b34657d88141cbbc5be3456beef4e934ba475d29ed2ba278610268ac6cdaa80d6b0f022100f6509561ac9bec64aa798506be26aa332edd5c00a80f6b8265bb592951f6d0f502820100211c59557bd301c7d6c5a0d2c35c7d2523499d98370e690a9a633e25fd32c1dd7ddc6d5697b5d49301b673455426fdd646f4de3234822fb50636661c6a42655abaad05306eb104544c04c46fd8795ea9925453b5df6262d1d534ac9377506a10aa2a03ae863490afa1d3b370744b5dddaf0c5abe3b1499704fb41c0b9e033855300cb83dda07efab98da8cf3ef388405d92bea43b7f33eacb8fdde65f50fbb84e44b4fe01cc69fc51e0e0290d7b7b315f72df7f206554c60c678fbc6768190d76f3bb7015b2e747adf22a5b0a9c7379b9992bfdd893f82c933aa1da673f975d967dbc3e6919029aa93e0b05118eb302e31499cbcdec3a7a26cdd41b0fae6425b04220220267d7b925e395576b282e0a8da5f4b5bf8da5718eb4064497eff3db7226ca03f")

	ParsePKCS8PrivateKey(b)
}

func TestMarshalAndParsePKIXPublicKey(t *testing.T) {
	t.Run("EC-KCDSA", func(t *testing.T) {
		for _, curve := range curveList {
			p1p, _ := eckcdsa.GenerateKey(curve, rand.Reader)
			p1 := &p1p.PublicKey

			der, err := MarshalPKIXPublicKey(p1)
			if err != nil {
				t.Error(err)
				return
			}

			log.Println(hex.EncodeToString(der))

			p2, err := ParsePKIXPublicKey(der)
			if err != nil {
				t.Error(err)
				return
			}

			if !p1.Equal(p2) {
				t.Error("not equals!")
				return
			}
		}
	})
	t.Run("KCDSA", func(t *testing.T) {
		for _, size := range sizeList {
			var p1p kcdsa.PrivateKey

			_ = kcdsa.GenerateParameters(&p1p.Parameters, rand.Reader, size)
			_ = kcdsa.GenerateKey(&p1p, rand.Reader)

			p1 := &p1p.PublicKey

			der, err := MarshalPKIXPublicKey(p1)
			if err != nil {
				t.Error(err)
				return
			}

			p2, err := ParsePKIXPublicKey(der)
			if err != nil {
				t.Error(err)
				return
			}

			if !p1.Equal(p2) {
				t.Error("not equals!")
				return
			}
		}
	})
}

func TestMarshalAndParsePKCS8PrivateKey(t *testing.T) {
	t.Run("EC-KCDSA", func(t *testing.T) {
		for _, curve := range curveList {
			p1, _ := eckcdsa.GenerateKey(curve, rand.Reader)

			der, err := MarshalPKCS8PrivateKey(p1)
			if err != nil {
				t.Error(err)
				return
			}

			p2, err := ParsePKCS8PrivateKey(der)
			if err != nil {
				t.Error(err)
				return
			}

			if !p1.Equal(p2) {
				t.Error("not equals!")
				return
			}
		}
	})

	t.Run("KCDSA", func(t *testing.T) {
		for _, size := range sizeList {
			var p1 kcdsa.PrivateKey

			_ = kcdsa.GenerateParameters(&p1.Parameters, rand.Reader, size)
			_ = kcdsa.GenerateKey(&p1, rand.Reader)

			der, err := MarshalPKCS8PrivateKey(&p1)
			if err != nil {
				t.Error(err)
				return
			}

			p2, err := ParsePKCS8PrivateKey(der)
			if err != nil {
				t.Error(err)
				return
			}

			if !p1.Equal(p2) {
				t.Error("not equals!")
				return
			}
		}
	})
}

func TestMarshalAndParseMarshalECPrivateKey(t *testing.T) {
	for _, curve := range curveList {
		p1, _ := eckcdsa.GenerateKey(curve, rand.Reader)

		der, err := MarshalECPrivateKey(p1)
		if err != nil {
			t.Error(err)
			return
		}

		p2, err := ParseECPrivateKey(der)
		if err != nil {
			t.Error(err)
			return
		}

		if !p1.Equal(p2) {
			t.Error("not equals!")
			return
		}
	}
}
