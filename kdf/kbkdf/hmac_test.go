package kbkdf

// TTAK.KO-12.0333-Part2

import (
	"bytes"
	"crypto/sha256"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
)

// 5 카운터 모드를 이용한 키 유도 함수 참조구현값
// 5.1 HMAC-SHA-224의 단계별 참조구현값
func TestHMAC_Counter_SHA224(t *testing.T) {
	answer := CounterMode(
		NewHMACPRF(sha256.New224),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
		1,
		0x230/8,
	)
	expect := internal.HB(`e50dccef026d0d446f34c90f919bbff924cb57e09caa9ec30f05eaba00619c4998eb35f28295dac8ef49efd6865c7e5f847d62daee89339dfc04d77a29a3edf406ffe053219d`)

	if !bytes.Equal(answer, expect) {
		t.Fail()
	}
}

// 6 피드백 모드를 이용한 키 유도 함수 참조구현값
// 6.1.1 카운터 입력을 사용하고 IV ≠ ∅는 경우
func TestHMAC_Feedback_SHA224_Counter_IV(t *testing.T) {
	answer := FeedbackMode(
		NewHMACPRF(sha256.New224),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
		1,
		0x230/8,
	)
	expect := internal.HB(`fc211d4bbb65e20544617963fdab1f5217837ec583d45b4fb45208ffd56179c61dc5b54e04dbab27481706941cafe4565f89d40f16a6ce96d2d2f9d71b3aad54135b04c69c5e`)
	if !bytes.Equal(answer, expect) {
		t.Fail()
	}
}

// 6 피드백 모드를 이용한 키 유도 함수 참조구현값
// 6.1.2 카운터 입력을 사용하고 IV = ∅인 경우
func TestHMAC_Feedback_SHA224_Counter(t *testing.T) {
	answer := FeedbackMode(
		NewHMACPRF(sha256.New224),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(``),
		1,
		0x230/8,
	)
	expect := internal.HB(`e50dccef026d0d446f34c90f919bbff924cb57e09caa9ec30f05eaba66227939ddbc98fb6b066429c0b694399a21f42bad3445917ec5068b87b6f6c97145d096958594fb4bd4`)
	if !bytes.Equal(answer, expect) {
		t.Fail()
	}
}

// 6 피드백 모드를 이용한 키 유도 함수 참조구현값
// 6.1.3 카운터 입력을 사용하지 않고 IV ≠ ∅는 경우
func TestHMAC_Feedback_SHA224_IV(t *testing.T) {
	answer := FeedbackMode(
		NewHMACPRF(sha256.New224),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
		0,
		0x230/8,
	)
	expect := internal.HB(`08d1602aef799313203f31dee8296a72ada6da376a6433bc05981b40e80e72f8923290b616bc30fd96c3af7f133bbc9660bafbec2bcc656d2efaefb8c7b4da4d79d4d9f85c72`)
	if !bytes.Equal(answer, expect) {
		t.Fail()
	}
}

// 6 피드백 모드를 이용한 키 유도 함수 참조구현값
// 6.1.4 카운터 입력을 사용하지 않고 IV = ∅인 경우
func TestHMAC_Feedback_SHA224(t *testing.T) {
	answer := FeedbackMode(
		NewHMACPRF(sha256.New224),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(``),
		0,
		0x230/8,
	)
	expect := internal.HB(`d1ce68d860217fc020000127de74b7d0b8a83fd6212a014344683abf6e6f3681408fa8a8c0271b451d460487e645ee1cc3882c514b7e2c92cb3ec40f4a51538bac3d12905f56`)
	if !bytes.Equal(answer, expect) {
		t.Fail()
	}
}

// 7 더블-파이프라인 반복 모드를 이용한 키 유도 함수 참조구현값
// 7.1.1 카운터 입력을 사용한 경우
func TestHMAC_DoublePipeline_SHA224_Counter(t *testing.T) {
	answer := DoublePipeMode(
		NewHMACPRF(sha256.New224),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
		1,
		0x230/8,
	)
	expect := internal.HB(`91e89383285b9691d38d0f0833cac5a030cc6a51741753e26b5f042f2687b69e972feaeaf6d4732353fd6095e56e0aa338ad361c24df78b37f6503437debf14ceda8efa84c9a`)
	if !bytes.Equal(answer, expect) {
		t.Fail()
	}
}

// 7 더블-파이프라인 반복 모드를 이용한 키 유도 함수 참조구현값
// 7.1.2 카운터 입력을 사용하지 않는 경우
func TestHMAC_DoublePipeline_SHA224(t *testing.T) {
	answer := DoublePipeMode(
		NewHMACPRF(sha256.New224),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
		internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
		0,
		0x230/8,
	)
	expect := internal.HB(`6e6f3681408fa8a8c0271b451d460487e645ee1cc3882c514b7e2c92fede6faf07ccf39fba75fc19f5a5b6035b6bd49bebc8fb5caf39e4c37df1788b63485a1956d5c395c14e`)
	if !bytes.Equal(answer, expect) {
		t.Fail()
	}
}
