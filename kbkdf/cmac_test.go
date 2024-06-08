package kbkdf

// TTAK.KO-12.0272

import (
	"bytes"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/seed"
)

// I.1. 카운터 모드를 이용한 키 유도 함수 (KDF in Counter Mode)
func TestCMAC_Counter_SEED_R1(t *testing.T) {
	answer := CounterMode(
		nil,
		NewCMACPRF(seed.NewCipher),
		internal.HB(`d899a261ac0bf40b04f110681b8859e4`),
		internal.HB(`8ed363214d5ca03328964b636a889f83850bc41af69ce426d080490cf57465514ee2768d47fde795fad82ca919ca1c9553122d8017549d431ddd6b0a`),
		internal.HB(`b78e81454183e6fea2de1e842e0b782433fdecc33f70b60279b9284fe02f6df2517faffd9ac52c1ef442fc1145deee1edc53f55ffcba9ed9fe1087ca`),
		1,
		0x0200/8,
	)
	expect := internal.HB(`762dfe66668c7807611e8115156277e595597320b22c2e25881f1bcc292077166a39dbe9aa75e05d57d1f5d69f286eccd3d7e0fef68b3e1e88d5cd26f85ab7cb`)

	if !bytes.Equal(answer, expect) {
		t.Fail()
	}
}

// 2. 피드백 모드를 이용한 키 유도 함수 (KDF in Feedback Mode)
func TestCMAC_Feedback_SEED_R1(t *testing.T) {
	answer := FeedbackMode(
		nil,
		NewCMACPRF(seed.NewCipher),
		internal.HB(`c0ec404a828088e5f79e2840b4dfda41`),
		internal.HB(`67bd99e57f8b859fd13b89a0250b945a9ddb7408d85e8008daea9cc1a57305e639f1550f5e2a8cbd9dcc368c4391759e66bb332977670a237cb60753`),
		internal.HB(`ad7d03f1d618fe349992a3be3a3c8641799da4f937169e010fcf738614a6ca2db04312c1f592d9fbf409f789227ad5c5b4f6c587baca990fe28b8c9a`),
		internal.HB(`9e3ef8ab426d06cffccf9e7b016d4981`),
		1,
		0x0200/8,
	)
	expect := internal.HB(`3c6dba1fee26203c0d221826148013bafe5ab7886c59420e10870ff241e6a8e0525b31d62d6907a4233104c03beb06b9511249ba4b9e327d0f56e546180524fe`)

	if !bytes.Equal(answer, expect) {
		t.Fail()
	}
}

// 3. 더블파이프라인 반복 모드를 이용한 키 유도 함수 (KDF in Double-Pipeline Iteration Mode)
func TestCMAC_DoublePipeline_SEED_R1(t *testing.T) {
	answer := DoublePipeMode(
		nil,
		NewCMACPRF(seed.NewCipher),
		internal.HB(`5b898cdfdd9b46a8e3f8e5edb79247c4`),
		internal.HB(`50db604f57ac1864fb19657a60cef868c8e5c04025b08d78222cdeed819093a6bd5117b30d98b3bc69989be2f84e0043761b2686c2060db81c4cc1a3`),
		internal.HB(`c499e583a24617d0b741677665b75ab65bf95ea9971802d7e73d476b8267bce40697b89efa7d9efbf1a3fcd404f2e58f539df3326ae69adedf17f8fb`),
		1,
		0x0200/8,
	)
	expect := internal.HB(`13604558349b96745683036e8232ec5c78b923854d3b3d0e144befe773a406d1183b62cdaa05f7e560769fd286137a80b23faad6309f90bfcbd83f05c4f788a3`)

	if !bytes.Equal(answer, expect) {
		t.Fail()
	}
}
