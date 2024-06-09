//nolint:typecheck
//go:build go1.16
// +build go1.16

package ctrdrbg_test

import (
	"bytes"
	"crypto/cipher"
	_ "embed"
	"encoding/csv"
	"errors"
	"io"
	"testing"

	"github.com/RyuaNerin/go-krypto/aria"
	"github.com/RyuaNerin/go-krypto/hight"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/drbg/ctrdrbg"
	"github.com/RyuaNerin/go-krypto/lea"
	"github.com/RyuaNerin/go-krypto/seed"
)

/**
B.2 기반 블록 암호별 참조 구현값

B.2.1 개요
아래는 CTR_DRBG의 기반 블록 암호 알고리즘으로 SEED, ARIA, LEA, HIGHT를 사용하는 경우의 참조 구현값을 정리한다.
CTR_DRBG의 다양한 설정을 반영하여, 참조 구현값은 특정 입력 값에 대한 각 설정의 512비트 출력 난수열로 구성한다.
단, 카운터 블록 길이(ctr_len )는 블록 길이 전체로 고정한다.
기반 블록 암호 알고리즘의 키 길이와 블록 길이에 따른 입력값을 정리하면 <표 B.2-1>와 같다.

|---------+-----------+-------------------------+-------------|
| 키 길이 | 블록 길이 | 블록 암호 알고리즘      | 입력값 정보 |
|---------+-----------+-------------------------+-------------|
|   128   |     64    | HIGHT                   | <표 B.2-5>  |
|   128   |    128    | SEED, ARIA-128, LEA-128 | <표 B.2-2>  |
|   192   |    128    | ARIA-192, LEA-192       | <표 B.2-3>  |
|   256   |    128    | ARIA-256, LEA-256       | <표 B.2-4>  |
|---------+-----------+-------------------------+-------------|
*/

type testcaseInput struct {
	KeySize   int
	BlockSize int

	entropy    [][]byte
	nonce      []byte
	personal   []byte
	additional [][]byte
}

var testCaseInputs = []testcaseInput{
	// <표 B.2-2> 난수 생성 입력 값(키 길이: 128 비트, 블록 길이: 128 비트)
	{
		KeySize:   128 / 8,
		BlockSize: 128 / 8,
		// SEED, ARIA128, LEA128

		entropy: [][]byte{
			internal.HB(`00010203 04050607 08090A0B 0C0D0E0F 10111213 14151617 18191A1B 1C1D1E1F`),
			internal.HB(`80818283 84858687 88898A8B 8C8D8E8F 90919293 94959697 98999A9B 9C9D9E9F`),
			internal.HB(`C0C1C2C3 C4C5C6C7 C8C9CACB CCCDCECF D0D1D2D3 D4D5D6D7 D8D9DADB DCDDDEDF`),
		},
		nonce:    internal.HB(`20212223 24252627`),
		personal: internal.HB(`40414243 44454647 48494A4B 4C4D4E4F 50515253 54555657 58595A5B 5C5D5E5F`),
		additional: [][]byte{
			internal.HB(`60616263 64656667 68696A6B 6C6D6E6F 70717273 74757677 78797A7B 7C7D7E7F`),
			internal.HB(`A0A1A2A3 A4A5A6A7 A8A9AAAB ACADAEAF B0B1B2B3 B4B5B6B7 B8B9BABB BCBDBEBF`),
		},
	},
	// <표 B.2-3> 난수 생성 입력 값(키 길이: 192 비트, 블록 길이: 128 비트)
	{
		KeySize:   192 / 8,
		BlockSize: 128 / 8,
		// ARIA-192, LEA-192

		entropy: [][]byte{
			internal.HB(`00010203 04050607 08090A0B 0C0D0E0F 10111213 14151617 18191A1B 1C1D1E1F 20212223 24252627`),
			internal.HB(`80818283 84858687 88898A8B 8C8D8E8F 90919293 94959697 98999A9B 9C9D9E9F A0A1A2A3 A4A5A6A7`),
			internal.HB(`C0C1C2C3 C4C5C6C7 C8C9CACB CCCDCECF D0D1D2D3 D4D5D6D7 D8D9DADB DCDDDEDF E0E1E2E3 E4E5E6E7`),
		},
		nonce:    internal.HB(`20212223 24252627 28292A2B`),
		personal: internal.HB(`40414243 44454647 48494A4B 4C4D4E4F 50515253 54555657 58595A5B 5C5D5E5F 60616263 64656667`),
		additional: [][]byte{
			internal.HB(`60616263 64656667 68696A6B 6C6D6E6F 70717273 74757677 78797A7B 7C7D7E7F 80818283 84858687`),
			internal.HB(`A0A1A2A3 A4A5A6A7 A8A9AAAB ACADAEAF B0B1B2B3 B4B5B6B7 B8B9BABB BCBDBEBF C0C1C2C3 C4C5C6C7`),
		},
	},
	// <표 B.2-4> 난수 생성 입력 값(키 길이: 256 비트, 블록 길이: 128 비트)
	{
		KeySize:   256 / 8,
		BlockSize: 128 / 8,
		// ARIA-256, LEA-256

		entropy: [][]byte{
			internal.HB(`00010203 04050607 08090A0B 0C0D0E0F 10111213 14151617 18191A1B 1C1D1E1F 20212223 24252627 28292A2B 2C2D2E2F`),
			internal.HB(`80818283 84858687 88898A8B 8C8D8E8F 90919293 94959697 98999A9B 9C9D9E9F A0A1A2A3 A4A5A6A7 A8A9AAAB ACADAEAF`),
			internal.HB(`C0C1C2C3 C4C5C6C7 C8C9CACB CCCDCECF D0D1D2D3 D4D5D6D7 D8D9DADB DCDDDEDF E0E1E2E3 E4E5E6E7 E8E9EAEB ECEDEEEF`),
		},
		nonce:    internal.HB(`20212223 24252627 28292A2B 2C2D2E2F`),
		personal: internal.HB(`40414243 44454647 48494A4B 4C4D4E4F 50515253 54555657 58595A5B 5C5D5E5F 60616263 64656667 68696A6B 6C6D6E6F`),
		additional: [][]byte{
			internal.HB(`60616263 64656667 68696A6B 6C6D6E6F 70717273 74757677 78797A7B 7C7D7E7F 80818283 84858687 88898A8B 8C8D8E8F`),
			internal.HB(`A0A1A2A3 A4A5A6A7 A8A9AAAB ACADAEAF B0B1B2B3 B4B5B6B7 B8B9BABB BCBDBEBF C0C1C2C3 C4C5C6C7 C8C9CACB CCCDCECF`),
		},
	},
	// <표 B.2-5> 난수 생성 입력 값(키 길이: 128 비트, 블록 길이: 64 비트)
	{
		KeySize:   128 / 8,
		BlockSize: 64 / 8,
		// HIGHT

		entropy: [][]byte{
			internal.HB(`00010203 04050607 08090A0B 0C0D0E0F 10111213 14151617`),
			internal.HB(`80818283 84858687 88898A8B 8C8D8E8F 90919293 94959697`),
			internal.HB(`C0C1C2C3 C4C5C6C7 C8C9CACB CCCDCECF D0D1D2D3 D4D5D6D7`),
		},
		nonce:    internal.HB(`20212223 24252627`),
		personal: internal.HB(`40414243 44454647 48494A4B 4C4D4E4F 50515253 54555657`),
		additional: [][]byte{
			internal.HB(`60616263 64656667 68696A6B 6C6D6E6F 70717273 74757677`),
			internal.HB(`A0A1A2A3 A4A5A6A7 A8A9AAAB ACADAEAF B0B1B2B3 B4B5B6B7`),
		},
	},
}

//go:embed testvector.csv
var testvectorCSV []byte

func TestCTRDRBG_B2(t *testing.T) {
	dst := make([]byte, 256)

	r := csv.NewReader(bytes.NewReader(testvectorCSV))
	r.Comma = '\t'
	r.ReuseRecord = true
	r.Read()

	testCaseIdx := 0
	for {
		testCaseIdx++
		records, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			panic(err)
		}
		var (
			algorithm               = records[0]              // A
			useDerivationFunction   = records[1] == "1"       // B 유도 함수 사용
			description             = records[2]              // C
			usePredictionResistance = records[3] == "1"       // D 예측 내성 사용
			usePersonal             = records[4] == "1"       // E 개별화 문자열 사용
			useAdditional           = records[5] == "1"       // F 추가 입력 사용
			refreshInterval1        = records[6] == "1"       // G 갱신 주기 1
			refreshInterval2        = records[7] == "1"       // H 갱신 주기 2
			output1                 = internal.HB(records[8]) // I
			output2                 = internal.HB(records[9]) // J
		)

		var input testcaseInput
		var newCipher func(key []byte) (cipher.Block, error)
		switch algorithm {
		case "SEED":
			input = testCaseInputs[0]
			newCipher = seed.NewCipher
		case "ARIA128":
			input = testCaseInputs[0]
			newCipher = aria.NewCipher
		case "LEA128":
			input = testCaseInputs[0]
			newCipher = lea.NewCipher

		case "ARIA192":
			input = testCaseInputs[1]
			newCipher = aria.NewCipher
		case "LEA192":
			input = testCaseInputs[1]
			newCipher = lea.NewCipher

		case "ARIA256":
			input = testCaseInputs[2]
			newCipher = aria.NewCipher
		case "LEA256":
			input = testCaseInputs[2]
			newCipher = lea.NewCipher

		case "HIGHT":
			input = testCaseInputs[3]
			newCipher = hight.NewCipher
		}

		var refreshInterval uint64 = 0
		if refreshInterval1 {
			refreshInterval = 1
		} else if refreshInterval2 {
			refreshInterval = 2
		}

		var personalizationString []byte
		additionalData := [][]byte{nil, nil}
		if usePersonal {
			personalizationString = input.personal
		}
		if useAdditional {
			additionalData = input.additional
		}

		state := ctrdrbg.Instantiate_CTR_DRBG(
			newCipher,
			input.KeySize,
			refreshInterval,
			0,
			input.entropy[0],
			input.nonce,
			personalizationString,
			useDerivationFunction,
			usePredictionResistance,
		)

		entropyIdx := 0

		dst = internal.Grow(dst, len(output1))
		err = state.Generate_CTR_DRBG(dst, func() ([]byte, error) { entropyIdx++; return input.entropy[entropyIdx], nil }, additionalData[0])
		if err != nil {
			t.Error(err)
			return
		}
		if !bytes.Equal(dst, output1) {
			t.Errorf(
				"case %3d FAILED - output1 - %7s df:%5v  pr:%5v  p:%5v  a:%5v  i1:%5v  i2:%5v - %s\n",
				testCaseIdx, algorithm,
				useDerivationFunction, usePredictionResistance, usePersonal, useAdditional, refreshInterval1, refreshInterval2,
				description,
			)
			return
		} else {
			dst = internal.Grow(dst, len(output2))
			err = state.Generate_CTR_DRBG(dst, func() ([]byte, error) { entropyIdx++; return input.entropy[entropyIdx], nil }, additionalData[1])
			if err != nil {
				t.Error(err)
				return
			}
			if !bytes.Equal(dst, output2) {
				t.Errorf(
					"case %3d FAILED - output2 - %7s df:%5v  pr:%5v  p:%5v  a:%5v  i1:%5v  i2:%5v - %s\n",
					testCaseIdx, algorithm,
					useDerivationFunction, usePredictionResistance, usePersonal, useAdditional, refreshInterval1, refreshInterval2,
					description,
				)
				return
			}
		}
	}
}
