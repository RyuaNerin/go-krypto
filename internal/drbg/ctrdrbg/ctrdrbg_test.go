package ctrdrbg_test

import (
	"bytes"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/drbg/ctrdrbg"
	"github.com/RyuaNerin/go-krypto/seed"
)

// TTAK.KO-12.0189/R2

func TestCTRDRBG_B1(t *testing.T) {
	/**
	B.1 세부 함수 참조 구현값

	아래는 CTR_DRBG를 구성하는 세부 함수에 대한 참조 구현값을 정리한다.
	이를 위해 CTR_DRBG를 <표 B.1-1>과 같은 설정으로 사용하는 것을 가정한다.

	|-----------------+----------|
	| 설정 옵션       | 설정값   |
	|-----------------+----------|
	| 유도 함수       | 사용     |
	| 예측 내성       | 미지원   |
	| reseed_interval | 1        |
	| 개별화 문자열   | 사용     |
	| 추가 입력       | 사용     |
	| ctr_len         | blocklen |
	|-----------------+----------|

	이러한 설정을 도시하면 (그림 B.1-1)과 같다.
	그리고 CTR_DRBG의 기반 블록 암호 알고리즘은 SEED,
	출력 난수열의 비트 길이는 512로 설정하며, CTR_DRBG 입력값은 <표 B.2-1>와 같다.
	*/

	dst := make([]byte, 1024)

	////////////////////////////////////////////////////////////////////////////////////////////////////
	// a) 초기화 함수
	state := ctrdrbg.Instantiate_CTR_DRBG(
		seed.NewCipher,
		128/8,
		1,
		0,
		internal.HB(`00010203 04050607 08090A0B 0C0D0E0F 10111213 14151617 18191A1B 1C1D1E1F`),
		internal.HB(`20212223 24252627`),
		internal.HB(`40414243 44454647 48494A4B 4C4D4E4F 50515253 54555657 58595A5B 5C5D5E5F`),
		true,
		false,
	)

	keySize := state.KeyLenByte
	blockSize := state.BlockLenByte

	if !bytes.Equal(state.Key[:keySize], internal.HB(`6EEB8C7D A6047793 F4554882 35A9A6B2`)) ||
		!bytes.Equal(state.V[:blockSize], internal.HB(`96E22D00 A182179A CDDBE18A 7D1A8559`)) {
		t.Fail()
		return
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////
	// b) 첫 번째 출력 생성
	dst = internal.Grow(dst, 256/8)

	// 1) 유도 함수
	err := state.Generate_CTR_DRBG(
		dst,
		nil,
		internal.HB(`60616263 64656667 68696A6B 6C6D6E6F 70717273 74757677 78797A7B 7C7D7E7F`),
	)
	if err != nil {
		t.Error(err)
		return
	}
	// 3) 첫 번째 출력
	if !bytes.Equal(dst, internal.HB(`E6150AB3 97A9E74B B4235730 FFE40CFB 7B6A1E97 A9B934A0 D6FBC0CA 83371B72`)) {
		t.Fail()
		return
	}
	// 4) 갱신 함수
	if !bytes.Equal(state.Key[:keySize], internal.HB(`DCB95524 A9264141 7D1880C1 688680C0`)) ||
		!bytes.Equal(state.V[:blockSize], internal.HB(`9F1DA519 64CE39E7 E59A6C6A 60BACB95`)) {
		t.Fail()
		return
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////
	// c) 두 번째 출력 생성
	dst = internal.Grow(dst, 256/8)

	// 1) 재초기화 함수
	err = state.Generate_CTR_DRBG(
		dst,
		func() ([]byte, error) {
			return internal.HB(`80818283 84858687 88898A8B 8C8D8E8F 90919293 94959697 98999A9B 9C9D9E9F`), nil
		},
		internal.HB(`A0A1A2A3 A4A5A6A7 A8A9AAAB ACADAEAF B0B1B2B3 B4B5B6B7 B8B9BABB BCBDBEBF`),
	)
	if err != nil {
		t.Error(err)
		return
	}
	// 2) 두 번째 출력
	if !bytes.Equal(dst, internal.HB(`F09F73FB C00D7EED C8DB54A6 1314DE17 0D167246 D93FE1D4 AAA075E2 B0F73703`)) {
		t.Fail()
		return
	}
	// 3) 갱신 함수
	if !bytes.Equal(state.Key[:keySize], internal.HB(`40EC068E B8A93E05 C920861F 5C6ACB98`)) ||
		!bytes.Equal(state.V[:blockSize], internal.HB(`67A8DE4B 8B171A62 8B19504F 170775D9`)) {
		t.Fail()
		return
	}
}
