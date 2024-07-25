package kbkdf_test

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/kbkdf"
)

func TestHMAC_CounterMode(t *testing.T) {
	var answer []byte
	for idx, tc := range testcasesHMACCtr {
		expect := tc.K0
		answer = kbkdf.CounterMode(answer[:0], kbkdf.NewHMACPRF(tc.Hash), tc.KI, tc.Label, tc.Context, tc.CounterSize, tc.L/8)

		if !bytes.Equal(answer, expect) {
			t.Errorf("failed test case %d\nexpect %x\nanswer %x", idx, expect, answer)
		}
	}
}

func TestHMAC_FeedbackMode(t *testing.T) {
	var answer []byte
	for idx, tc := range testcasesHMACFB {
		expect := tc.K0
		answer = kbkdf.FeedbackMode(answer[:0], kbkdf.NewHMACPRF(tc.Hash), tc.KI, tc.Label, tc.Context, tc.IV, tc.CounterSize, tc.L/8)

		if !bytes.Equal(answer, expect) {
			t.Errorf("failed test case %d\nexpect %x\nanswer %x", idx, expect, answer)
		}
	}
}

func TestHMAC_DoublePipeMode(t *testing.T) {
	var answer []byte
	for idx, tc := range testcasesHMACDP {
		expect := tc.K0
		answer = kbkdf.DoublePipeMode(answer[:0], kbkdf.NewHMACPRF(tc.Hash), tc.KI, tc.Label, tc.Context, tc.CounterSize, tc.L/8)

		if !bytes.Equal(answer, expect) {
			t.Errorf("failed test case %d\nexpect %x\nanswer %x", idx, expect, answer)
		}
	}
}

/**
TTAK.KO-12.0333-Part2

HMAC 기반 키 유도 함수
- 제2부: 해시 함수 SHA-2

HMAC-based Key Derivation Functions
- Part2: Hash Function SHA-2
*/

type testVectorHMAC struct {
	Hash                   func() hash.Hash
	KI, Label, Context, IV []byte
	L                      int // bits
	CounterSize            int
	K0                     []byte
}

var (

	// 5 카운터 모드를 이용한 키 유도 함수 참조구현값
	testcasesHMACCtr = []testVectorHMAC{
		{
			// 5.1 HMAC-SHA-224의 단계별 참조구현값
			Hash:        sha256.New224,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			L:           0x0230,
			CounterSize: 1,
			K0:          internal.HB(`e50dccef026d0d446f34c90f919bbff924cb57e09caa9ec30f05eaba00619c4998eb35f28295dac8ef49efd6865c7e5f847d62daee89339dfc04d77a29a3edf406ffe053219d`),
		},
		{
			// 5.2 HMAC-SHA-256의 단계별 참조구현값
			Hash:        sha256.New,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			L:           0x0280,
			CounterSize: 1,
			K0:          internal.HB(`8da6974e8e8683042d73b24239658e7e2cef712e9335059cd34d5b2f8a25e9c94a6e8b64bddc77688aad7390ce08e9107bf1d165edf0c41dad8e9195549da3b2ef33b2af0c803c6d36c392f745753ab9`),
		},
		{
			// 5.3 HMAC-SHA-384의 단계별 참조구현값
			Hash:        sha512.New384,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			L:           0x03C0,
			CounterSize: 1,
			K0:          internal.HB(`aa66fcf02050a4a0ae12eb9a7be4263374fcd08127634e05d6afe4237d10a6e75042ab595d9b08d72bcbd11d5cf727e9420f220a19aad1e5aaf156aebf70e6dfdd0d9f350548cbc175e29c533138f99272ad7767952258d0d30b2d749d41cafd0234aaa094dd2189c49f6568b4c001dd52d8ad6b63c724f8`),
		},
		{
			// 5.4 HMAC-SHA-512의 단계별 참조구현값
			Hash:        sha512.New,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			L:           0x0500,
			CounterSize: 1,
			K0:          internal.HB(`19bd6999e03d0250ee5ae90a78429897fcaf7498c6a2fc44245ee4bd9455b19e343ab44e98d97d8a75460126647579a1c8b4d9ce796f38688acb1d03613f7dd359a3df3109ce3dc95df83cba0e13d991ea48abd145025e942f1cf78dad4b62a9f7aa0b532c64b5ddc233dd6db910c83664cb373b6ba3dbd5d4751c5284c9aa2f73aa82e6233b6a8a51b08441127e838700e4480ee8b5300cbad5fe48ba2a4ab1`),
		},
	}

	// 6 피드백 모드를 이용한 키 유도 함수 참조구현값
	testcasesHMACFB = []testVectorHMAC{
		// 6.1 HMAC-SHA-224의 단계별 참조구현값
		{
			// 6.1.1 카운터 입력을 사용하고 IV ≠ ∅인 경우
			Hash:        sha256.New224,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			IV:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
			L:           0x0230,
			CounterSize: 1,
			K0:          internal.HB(`fc211d4bbb65e20544617963fdab1f5217837ec583d45b4fb45208ffd56179c61dc5b54e04dbab27481706941cafe4565f89d40f16a6ce96d2d2f9d71b3aad54135b04c69c5e`),
		},
		{
			// 6.1.2 카운터 입력을 사용하고 IV = ∅인 경우
			Hash:        sha256.New224,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			IV:          nil,
			L:           0x0230,
			CounterSize: 1,
			K0:          internal.HB(`e50dccef026d0d446f34c90f919bbff924cb57e09caa9ec30f05eaba66227939ddbc98fb6b066429c0b694399a21f42bad3445917ec5068b87b6f6c97145d096958594fb4bd4`),
		},
		{
			// 6.1.3 카운터 입력을 사용하지 않고 IV ≠ ∅인 경우
			Hash:        sha256.New224,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			IV:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
			L:           0x0230,
			CounterSize: 0,
			K0:          internal.HB(`08d1602aef799313203f31dee8296a72ada6da376a6433bc05981b40e80e72f8923290b616bc30fd96c3af7f133bbc9660bafbec2bcc656d2efaefb8c7b4da4d79d4d9f85c72`),
		},
		{
			// 6.1.4 카운터 입력을 사용하지 않고 IV = ∅인 경우
			Hash:        sha256.New224,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			IV:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
			L:           0x0230,
			CounterSize: 0,
			K0:          internal.HB(`08d1602aef799313203f31dee8296a72ada6da376a6433bc05981b40e80e72f8923290b616bc30fd96c3af7f133bbc9660bafbec2bcc656d2efaefb8c7b4da4d79d4d9f85c72`),
		},
		// 6.2 HMAC-SHA-256의 단계별 참조구현값
		{
			// 6.2.1 카운터 입력을 사용하고 IV ≠ ∅인 경우
			Hash:        sha256.New,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			IV:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			L:           0x0280,
			CounterSize: 1,
			K0:          internal.HB(`c976ededeabac660620a1bfb10f8f1471c445600d4b639de34cff39d7796f0a4147eb3426930257c2dd9990bdf95441f49f78bc03c6e63a2c306c32b46fac3afbae4ccf905d7fcc6127167fa77506709`),
		},
		{
			// 6.2.2 카운터 입력을 사용하고 IV = ∅인 경우
			Hash:        sha256.New,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			IV:          nil,
			L:           0x0280,
			CounterSize: 1,
			K0:          internal.HB(`8da6974e8e8683042d73b24239658e7e2cef712e9335059cd34d5b2f8a25e9c90e905156c6141095efe48d6d797c51a5fe3a51202a99390697294d3f784c8388e3da32cb27d3ac2b32eaca99b971c2f0`),
		},
		{
			// 6.2.3 카운터 입력을 사용하지 않고 IV ≠ ∅인 경우
			Hash:        sha256.New,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			IV:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			L:           0x0280,
			CounterSize: 0,
			K0:          internal.HB(`0dec972f54944b77c2bc3d54ba89e2d513c4a1226a738631f14fec1f919ed7f12e98d548a0e293238a7dc423d3bc8f92bc5253b59904709732817730939e6e1225ca8a6c45ee39c2b65a146d7bc8d36a`),
		},
		{
			// 6.2.4 카운터 입력을 사용하지 않고 IV = ∅인 경우
			Hash:        sha256.New,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			IV:          nil,
			L:           0x0280,
			CounterSize: 0,
			K0:          internal.HB(`8005439e7f5d28f81d1fe0f640b5d6dfafa8e9266bcc248ad6197280c8a65f1bd22a07b187c697ff56190d7bf61ea5c6d2217abd3a114a643d6ffe9259c50db2fcb17fb72b29463293a580f394ca3c9d`),
		},
		// 6.3 HMAC-SHA-384의 단계별 참조구현값
		{
			// 6.3.1 카운터 입력을 사용하고 IV ≠ ∅인 경우
			Hash:        sha512.New384,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			IV:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			L:           0x03C0,
			CounterSize: 1,
			K0:          internal.HB(`69ae00a7556ecec8769493b4937a473e25a322a0f9e9a68326ddb0ce512c5cf0b4dd7823ede1d4ee45cdee7bc5e5a34769aee0fb17194af45d2edfd6b9f6e87f3f13419efe31a965c20284b5c956b4a0ac14f24b479d589556f5e10a3dfd6fa8bd2d93fcdc02585eafc055f3db4446873add76fc465cf4ea`),
		},
		{
			// 6.3.2 카운터 입력을 사용하고 IV = ∅인 경우
			Hash:        sha512.New384,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			IV:          nil,
			L:           0x03C0,
			CounterSize: 1,
			K0:          internal.HB(`aa66fcf02050a4a0ae12eb9a7be4263374fcd08127634e05d6afe4237d10a6e75042ab595d9b08d72bcbd11d5cf727e9b70d5096bd3fa298448355184bea1198b8213bd9a8fb485910209cd630115e4bd87c16898ee0d0b5f20da76fe1d82bf9e05118a3f8d5a1dc74ebbeaf401aac96394ccfdd425a2eeb`),
		},
		{
			// 6.3.3 카운터 입력을 사용하지 않고 IV ≠ ∅인 경우
			Hash:        sha512.New384,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			IV:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			L:           0x03C0,
			CounterSize: 0,
			K0:          internal.HB(`612ca163396ec96b6ae094eaed548e1f649adce06300ab252705f65124ca63e9eac39e5bbb3e440f62a5d2595cccebed80ebd68a337581817a4e1bd0c320ffe36ab83c29a415660e0e859fcf3e4c6eb2107425a179263dcdd79cb97dbb810afa3b09a0cd2e1045f88666581da70294b0fc69032cd4c98fb3`),
		},
		{
			// 6.3.4 카운터 입력을 사용하지 않고 IV = ∅인 경우
			Hash:        sha512.New384,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			IV:          nil,
			L:           0x03C0,
			CounterSize: 0,
			K0:          internal.HB(`4011a8396cf4caa2aedc8e7ad400641e94af29cc98d6be574203db5b1276dab20cccedca5a6f4678eeefdf5effe0d05ff622a6303e4eeea0fc3e9c6dc3349fc7ecbaab0de9b069e3448f80d49f9cd0d707ccd6fa5a8acdb84b4a5b47f5adf387698c3ac63e07d18c589a880093c4bd28355a7a656b2f6493`),
		},
		// 6.4 HMAC-SHA-512의 단계별 참조구현값
		{
			// 6.4.1 카운터 입력을 사용하고 IV ≠ ∅인 경우
			Hash:        sha512.New,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			IV:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			L:           0x0500,
			CounterSize: 1,
			K0:          internal.HB(`833fa5ecd02726858d6486bf3a93429b0d2f3882d791b26434121e7de50a66414707625fcd0067390f4ead1ab4e511d62f479aaef7aa22e3876ed3ee9866e7405b6cf435bbd313dbea90752d495859f40209e38dac39bfa2d9bdd28a70db4db0829304d75c7c62c792e87a728832a91eb67d0b5d56e3c7a972a6cb5388610899306c36630ee1a80c886e45eb8131884b0c02b914b8ab847917b8407b8d029408`),
		},
		{
			// 6.4.2 카운터 입력을 사용하고 IV = ∅인 경우
			Hash:        sha512.New,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			IV:          nil,
			L:           0x0500,
			CounterSize: 1,
			K0:          internal.HB(`19bd6999e03d0250ee5ae90a78429897fcaf7498c6a2fc44245ee4bd9455b19e343ab44e98d97d8a75460126647579a1c8b4d9ce796f38688acb1d03613f7dd3a91e9a8953e9593a2ea6d8a6311e0e80bcf3d6dd81fc1b9230e657464df9997d032a76725d868b6ded739a70e273e5c29fe922a03200371d3e587c2dcfd7a02083746b74c09591995ed75e917ebd77532e5d0da75416a668d76ee6b877d05bc5`),
		},
		{
			// 6.4.3 카운터 입력을 사용하지 않고 IV ≠ ∅인 경우
			Hash:        sha512.New,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			IV:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			L:           0x0500,
			CounterSize: 0,
			K0:          internal.HB(`2ddd5c4ae498189e38318472306347206e4ef9cb5447c605a0104b054922b13d2e2e8b98a10eed328a208b3e4bc27a68f397d5a28e8d10d8b7c1f402c90f250489179d8617a963b6e4f1499888beb7d794535c2366376f7ecd86e7d779528aefca7f926befd98e5d1951af5a157c25f7c9d3a0c08bbb85ac30995732dd13ae4ab9fd332c98a7908ccffec6e9574b45a5284bff4f68a96618c4fa88e95c29e39d`),
		},
		{
			// 6.4.4 카운터 입력을 사용하지 않고 IV = ∅인 경우
			Hash:        sha512.New,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			IV:          nil,
			L:           0x0500,
			CounterSize: 0,
			K0:          internal.HB(`b78b6c62a332158347e6afb027d1d74ff47a35b5e75aa0a52f5bce12aedc372091063e61e4fa73c05db78211b572fc7821abc8e473dc1cbed34471ec32a4ee9cad519e1bcb7dad930fe966ac6e3a0d410afeadcf0ee812639043027db76723e9d7a2695653c61283fe33bafbbec85d9263d50c6911e1a16932c0d65afbbc6001e760fd3c41090f06aa2259924f5da09a16f6bb9cfd86e2c27d406b1a64c233c2`),
		},
	}

	// 7 더블-파이프라인 반복 모드를 이용한 키 유도 함수 참조구현값
	testcasesHMACDP = []testVectorHMAC{
		// 7.1 HMAC-SHA-224의 단계별 참조구현값
		{
			// 7.1.1 카운터 입력을 사용한 경우
			Hash:        sha256.New224,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			L:           0x0230,
			CounterSize: 1,
			K0:          internal.HB(`91e89383285b9691d38d0f0833cac5a030cc6a51741753e26b5f042f2687b69e972feaeaf6d4732353fd6095e56e0aa338ad361c24df78b37f6503437debf14ceda8efa84c9a`),
		},
		{
			// 7.1.2 카운터 입력을 사용하지 않는 경우
			Hash:        sha256.New224,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			L:           0x0230,
			CounterSize: 0,
			K0:          internal.HB(`6e6f3681408fa8a8c0271b451d460487e645ee1cc3882c514b7e2c92fede6faf07ccf39fba75fc19f5a5b6035b6bd49bebc8fb5caf39e4c37df1788b63485a1956d5c395c14e`),
		},
		// 7.2 HMAC-SHA-256의 단계별 참조구현값
		{
			// 7.1.1 카운터 입력을 사용한 경우
			Hash:        sha256.New,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			L:           0x0280,
			CounterSize: 1,
			K0:          internal.HB(`87fbdd5fa8238f98deee8715a6987b33db68ee2807b289a0f42fe99350784f9c7899174abc221c4d2dcd21d850332223ca01b827324ae506c8ef58fa165a23534f998ed2dfea25a935f964e14831cabc`),
		},
		{
			// 7.1.2 카운터 입력을 사용하지 않는 경우
			Hash:        sha256.New,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			L:           0x0280,
			CounterSize: 0,
			K0:          internal.HB(`d22a07b187c697ff56190d7bf61ea5c6d2217abd3a114a643d6ffe9259c50db2a62327c9f79a5b723b3ed727aadc58470fa473d3171d02ecd2c8017b4ba48fd6edfcc08ca9a3197aad5f19591f0c4812`),
		},
		// 7.3 HMAC-SHA-384의 단계별 참조구현값
		{
			// 7.3.1 카운터 입력을 사용한 경우
			Hash:        sha512.New384,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			L:           0x03C0,
			CounterSize: 1,
			K0:          internal.HB(`d4cf1159dbf331b40985a160531d73f2166303ca0bb46ad717e9207cd36a29f42eb748490e94e07596ae3a1f099f928a78e27a472f7d9bd0965220846f2367af5c8e16c9d680841a5b02e5fcbf65d70dc82cfd3c6d7589536a24ac5342529359c13536d6b440a652d0c62d13a3411f26495dfc680ac55fd4`),
		},
		{
			// 7.3.2 카운터 입력을 사용하지 않는 경우
			Hash:        sha512.New384,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			L:           0x03C0,
			CounterSize: 0,
			K0:          internal.HB(`f622a6303e4eeea0fc3e9c6dc3349fc7ecbaab0de9b069e3448f80d49f9cd0d707ccd6fa5a8acdb84b4a5b47f5adf387a048db4855c63f8a8d6a1b5d7c4c730535cbc6a647f291cc0c66870c1e01cc7a3cac9b07553da694b54ab2363d7972fdbf7ce73488ebf9045d5e978ceb1490e433f82dbe915fe223`),
		},
		// 7.4 HMAC-SHA-512의 단계별 참조구현값
		{
			// 7.3.1 카운터 입력을 사용한 경우
			Hash:        sha512.New,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			L:           0x0500,
			CounterSize: 1,
			K0:          internal.HB(`d3a6779fd531ff6f327ec70f26d1f54798c867e86eae28ae6ced773ba2382190a8c72fb5c912c9f8a7f4a9cdf65928d133e28b71764c841959d7a1133c08fd65df7ab5c457be9ce477ba71d3b54eb5c4cd824216e5f3267f300c071045cd368937f0e730328b364a6ae0a8e9a80ddac6430dc02247408e3b97a6a063502be9b0f52e2742021aab1f651b7063b25876310c0027bcf36ea9b4d577d9fc51f88b64`),
		},
		{
			// 7.3.2 카운터 입력을 사용하지 않는 경우
			Hash:        sha512.New,
			KI:          internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`),
			Label:       internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			Context:     internal.HB(`00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabb`),
			L:           0x0500,
			CounterSize: 0,
			K0:          internal.HB(`ad519e1bcb7dad930fe966ac6e3a0d410afeadcf0ee812639043027db76723e9d7a2695653c61283fe33bafbbec85d9263d50c6911e1a16932c0d65afbbc60017dbd12e3ff6b1aa617894cca40b5d3d168daed59d4c41a2464cbfc2d3b9408a171b5e7e14128ffdd9f9f4cc760907e8ed7a8cfdd3515f7d90ce5cb2c57e0fe19eaeb52a4ebb7a4ab25ac273dfb078f413dda8f65f81b8c7fb9aaebad8e36cc3d`),
		},
	}
)
