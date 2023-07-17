package kcdsa

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
)

var (
	// samples in TTAK.KO-12.0001_R4
	testCases = []struct {
		Key PrivateKey
		M   []byte
		R   *big.Int
		S   *big.Int
		K   *big.Int
	}{
		// p.30
		// Ⅱ.1 소수 p, q의 길이 (α, β) = (2048, 224), SHA-224 적용 예
		{
			Key: PrivateKey{
				PublicKey: PublicKey{
					Parameters: Parameters{
						P: hi_(`8da8c1b5 c95d11be 46661df5 8c9f803e b729b800 dd92751b 3a4f10c6 a5448e9f
								3bc0e916 f042e399 b34af9be e582ccfc 3ff5000c ff235694 94351cfe a5529ea3
								47dcf43f 302f5894 380709ea 2e1c416b 51a5cdfc 7593b18b 7e3788d5 1b9cc9ae
								828b4f8f b06e0e90 57f7fa0f 93bb0397 031fe7d5 0a6828da 0c1160a0 e66d4e5d
								2a18ad17 a811e70b 14f4f431 1a028260 3233444f 98763c5a 1e829c76 4cf36adb
								56980bd4 c54bbe29 7e790228 4292d75c a3600ff4 59310b09 291cbefb c721528a
								13403b8b 93b711c3 03a2182b 6e6397e0 83380bf2 886af3b9 afcc9f50 55d8b713
								6c0ebd08 c5cf0b38 888cd115 72787f6d f384c97c 91b58c31 dee5655e cbf3fa53`),
						Q: hi_(`864f1884 1ec103cd fd1be7fe e54650f2 2a3bb997 537f32cc 79a51f53`),
						G: hi_(`0e9be1f8 7a414d16 7a9a5a96 8b079e4a d385a357 3edb21aa 67a6f61c 0d00c14a
								7a225044 b6e9eb03 68c1eb57 b24b45cd 854fd93c 1b2dfb0a 3ea302d2 367e4ec7
								2f6e7ee8 ea7f8002 f7704e99 0b954f25 bada8da6 2baeb6f0 6953c0c8 5104ad03
								f36618f7 6c62f4ec f3480183 69850a56 17c999db e68ba17d 5bc72556 74ef4839
								22c6a3f9 9d3c3c6f 358896c4 e63c605e e7db16fc bd9be354 e281f7fe 7813d054
								27ed1912 b5c7653a 167b9434 9147eeaf 85cc9ce2 e81661f3 21512d5d 2c0580b0
								3d1704ee f2317f45 185c8258 387e7ec9 79c04707 ef546241 2784afe4 1a7b45c8
								3b9cbe48 f9127cb4 400be9e9 6ac5de17 f2c9dea3 5e3734e7 9b64673f 85681c4e`),
						Sizes: L2048N224SHA224,
					},
					Y: hi_(`04ede5c6 7ea29297 a8cacb6b de6f4666 aea27d10 3dd1e9e9 582f76a2 f22b8b1b
							32230bc5 8f06b768 f8102b49 fa1cae5e 18921494 7f6239b6 c6ce7c9b c2d230e8
							9a40bee2 c33a8861 fd4f7d35 b788fe95 b2d5885d 8c8faea8 1c90be4c ee2784e3
							3577a71d 3b7f085d 71e9a1d4 7815c73f a087acaa b9fcb565 5ac9570e 6852be7c
							9c0aecea 8bd9aa75 a44fc314 7f733e90 6adb0fd7 6d613561 b1db364b bdc9afd3
							ce8f5f17 e3e71203 4a999350 8059fa52 441fa90d dfe9a0f2 a0b9192f e2220c08
							1bd0c0f0 e07cb5f1 ee4ff405 23591f17 8a4fc7cb 5065f6a3 8216e9a0 99c205b2
							9b8746d8 65e1af6d 903e5a13 8004910b 70eb5b84 eed9760e a60578bf 08852898`),
				},
				X: hi_(`2f1991c1 af401872 8a5a431b 9b5459df b16f6d25 6797fe57 0ec6bc65`),
			},
			M: pb_(`54 68 69 73 20 69 73 20 61 20 74 65 73 74 20 6d 65 73 73 61 67 65 20 66
					6f 72 20 4b 43 44 53 41 20 75 73 61 67 65 21`),
			K: hi_(`49561994 fd2bad5e 410ca1c1 5c3fd3f1 2e70263f 2820ad5c 566ded80`),
			R: hi_(`ed b7 6a 2d 39 f3 d7 fa 16 d0 82 59 41 18 b0 cf 8b a5 76 92 cf 3b aa ec
					6f 6d d9 51`),
			S: hi_(`5260a2df 2e923de8 77b130ac 8b5e8b17 63973b88 d5d4627a dfbacf52`),
		},
		// Ⅱ.2 소수 p, q의 길이 (α, β) = (2048, 224), SHA-256 적용 예
		{
			Key: PrivateKey{
				PublicKey: PublicKey{
					Parameters: Parameters{
						P: hi_(`c3159a30 cdbcc00c e2a99043 9634f7d3 fb16feb1 2c579932 2c14f8b8 a0d9b98e
								35f724bf e14c4afc 475d78f9 3a83f8fb 4636a5de f357bd6f b0c6245c ac4ef29c
								8f7da5e9 b39f3158 f4fd27c8 4088bcbb 6286d964 29c90e82 b7f31bf3 e76e93c6
								8a3163cf b82370e2 75159d66 08f82601 013476d5 50b386ca 34736388 6df337d7
								a54db7e9 8cc2df0d 828c31eb c62f3bc2 3f070c89 9648e276 2b26ffed a9d88ffb
								f684c570 4937fedc 03f60c10 5b69542e d40f910b 4c66fc09 1f5e1c12 47628abc
								e989b74a b0ef6f1a 14e2567f c083991e 1c846242 0bb8fbf9 b3f67b66 b02de042
								0a18d49a 6d4896d0 d1dddbed 24ee1611 8090221f 9fe9a1e1 2194e0d2 b3c61c13`),
						Q: hi_(`bb6a5c40 316bd80e 78246e92 ac9bf881 a9eb0cb9 6c7212eb 1e46ae0d`),
						G: hi_(`487844c0 b67465b7 18f04dbd 453342b7 49076ee1 f4226f18 1db282e1 c51b0f29
								0dae9601 ac73ed1f 1b25adad d50bfb42 1e8a09fa 07689a93 e5fb52a5 f8012956
								b90641f8 45c4b7e4 45cafe2e 3284775b dd70bce4 0ef3274e 52cbc3d5 738da7a8
								61bc46c0 a9693aa8 7e0aae62 bd371fa0 14ffc69f 3625d5a1 fbaaac80 d81c78a5
								9badeae5 fdfea922 ebc330a1 37e7699a 2790e86b db270c21 35eab4e0 bcd28b77
								13a8b241 1534c63f 2edf4e00 5902f6cc 1a155c29 f3eae17f 88acb5c6 70f5cf19
								a5a54e87 6692ab82 08c4a9ef 75a29e74 f08f92ac 1a38592d 46a2557c 3a18c06e
								d6529b40 bc5ecff9 715329a2 c01b4245 874250ed 515537ee 7458f898 6ff920bc`),
						Sizes: L2048N224SHA256,
					},
					Y: hi_(`0712496f cf76ce98 8be97ac0 9f0dbbe6 2d58707a 767d608a 3301115d 479cc871
							4ce3a10b eb152552 46c2623e fe50bfd2 5a83c355 551574e6 e3560e7b d1cd5e7e
							8e1269a4 a6f1976c 84e8fe8e 32e55aed d548fced cc92a6e4 e1bf2d1f 2aa30c0c
							0a991c29 b2595029 f903b634 189aa70c fc429531 93016c1f 7bb6276d f3ebfae7
							c060b987 d89088a0 558fc132 27b86f7a 57dde307 1cc022e0 39be4b68 3858d782
							f52aa730 49d508ef 994a5039 cab5faf2 89bdac07 75efbb51 eb4d5ff9 99b71d59
							c4d833b5 d069202a 968f3ac3 5fa77baf bdd9c096 0752c5da f783929d e2dad916
							f1159e75 a345445d 63c5b422 e0bcd2ba d9379d14 43892ed5 d12f8285 3d51a705`),
				},
				X: hi_(`b55d61ec 0114e020 efc4c9bb 5f2f3d2e 38409e17 d3954174 6d94ff7c`),
			},
			M: pb_(`54 68 69 73 20 69 73 20 61 20 74 65 73 74 20 6d 65 73 73 61 67 65 20 66
					6f 72 20 4b 43 44 53 41 20 75 73 61 67 65 21`),
			K: hi_(`a5c22f64 dde15693 3ad15bcb 928d6a3b 5acf0d7a 2302615c e74ccad6`),
			R: hi_(`53 f7 31 8e 64 b6 1c cc 83 67 ac 08 51 19 a1 cb bb 25 51 0f e1 be c1 24
					c2 99 89 e0`),
			S: hi_(`b750f725 1585204c 236e4204 884166a2 6c6cf08b d281167a 5efadd52`),
		},
		// Ⅱ.3 소수 p, q의 길이 (α, β) = (2048, 256), SHA-256 적용 예
		{
			Key: PrivateKey{
				PublicKey: PublicKey{
					Parameters: Parameters{
						P: hi_(`d06eb9f2 75b3ac7f 2970b578 ad1c3173 2a012684 4776f95c f07b4194 c6def6f4
								16a66751 458b0667 cdbc44af 3f6b5877 0e674a86 1c8febf4 eea0e504 50ec5272
								26b84707 17ee768c f39cfd32 bc2540d2 924e0968 e64d47ee 4cf0ab6c d192284b
								826c7508 2e18840b 67bc4cb1 f1708173 f08825ba 4f6e5fb8 6a357f02 c06f8283
								f3cd58a1 ed4d3062 f4a5c0d2 f26e54c0 fa511b5e d5cfd270 19d4a90d da7aca50
								561397ab eede9cff 45ec6cf3 e22dac5c af454b7b 9b3b5ffe 16128197 768114c9
								cd4be4e9 ecdc431a 0cc0ed54 4fd4da1c 9e98a2c3 cb4297fe 1d1387d8 1c51d492
								5ede6a8b baf660ef 675549b4 aea5267f b5f778d5 308dd691 75de580e c316c4ef`),
						Q: hi_(`cfefed9c 75b5610f db100d91 c4cb8187 a0077917 33128ff1 43ffedf9 7f6ffd65`),
						G: hi_(`023fec34 dfa5e5ce 369dd782 b07034af 037ac187 28d43204 5739b986 1b0df1dc
								aeeb5c9e d3e025d8 3adcdae0 419c158b 09ee35ff 84ab9caa 9ed4e535 f982fb99
								e30d3195 37c05780 a2cf31cf 6bb226c6 6b7b3ed7 6b65dc65 8b216b86 7f186d98
								0d30d1a9 5285a081 c5aba363 939660a5 7596c621 2207e4e3 58b729bc 079778b4
								f385824c 0862cdce 08aeb2c6 58c18559 d3ed865c d6bed194 da447fd4 1789c74d
								352ed26b 56c2d128 f1154f73 3fe71f10 bf676c9f 7e4268c0 53d13152 997a2d9b
								fb73fccb 0dcea4c1 32f68f28 2a6db325 cc467fb7 f1fe2da5 f80fd32c ae781a75
								74845a3d 45712054 3987b348 d5d75b1b 954cba47 3f83951a 8c1be717 b953206c`),
						Sizes: L2048N256SHA256,
					},
					Y: hi_(`44ce4c95 da1ff8bf bc6b7277 ccc6694e 1b1e6dfa cf617533 354da0cf 6966e156
							2124003d b09e3330 9a24f87c 467917ae dfeb911f d5344422 06345275 7c40f0a0
							bb45acc8 e462c5ac 4d8dd0f9 2fcc80f3 3e4160f5 98682bf5 71163c43 bd703c2c
							1827db2e 2336511d 84520afa 97dc4962 40ea4a82 ca2ffc64 6363f822 d037c813
							8f3458a3 e41bd3a0 23b63cc1 13b33ecb 3fcccc5c bed325e7 ec1f07e2 03e9aa8e
							451c96fb dec927d6 ee741540 a90673b4 f2feac07 b6f4eda0 8db28fdf aed8634e
							7ff40582 ae33d8db f377a761 9ad1c006 68633779 2943e6cd 016d5534 e4122bca
							18d12075 79ea4c90 610a1496 b63c23dc 996b686e feb34c36 1f9afdcf 7e8fbf9a`),
				},
				X: hi_(`21e2cf86 8d004318 aca87261 476dfc67 c1983364 82fe1dcb 3cbb5ba0 f081158a`),
			},
			M: pb_(`54 68 69 73 20 69 73 20 61 20 74 65 73 74 20 6d 65 73 73 61 67 65 20 66
					6f 72 20 4b 43 44 53 41 20 75 73 61 67 65 21`),
			K: hi_(`0d30f8f9 2313f7a5 abe0b0de ec219e40 c4640c89 39222aa0 dd6a3329 55778025`),
			R: hi_(`59 49 00 77 f9 8c 21 78 85 09 cb 47 8c cd e7 7a 4f b5 41 4e 13 cf 92 81
					cb 80 97 5b 33 70 d9 7d`),
			S: hi_(`185f21b5 dbf4255b 954a4d62 cf363c32 73211147 cba054e8 3a87da2d d7e0741d`),
		},
		// Ⅱ.4 소수 p, q의 길이 (α, β) = (3072, 256), SHA-256 적용 예
		{
			Key: PrivateKey{
				PublicKey: PublicKey{
					Parameters: Parameters{
						P: hi_(`cbaeace3 677e98ad b2e49c00 2b8b0f43 4143b466 515839bf 813b097d 2d1ee681
								5008c27a 3415bc22 31609874 5e5844f3 3ecc8887 c16dfb1c fb77dc4c 3f3571cc
								eefd4291 8f6c48c3 702ab6ef 0919b7e8 402fc89b 35d09a0e 5040e309 1ee4674b
								e891933c 1007e017 edd40818 7e4114b6 be5548d7 8db58b84 8475a422 62d7eb79
								5f08d161 1055efea 8a6aeb20 eb0f1c22 f002a2e8 195bcbba 830b8461 3531bdd9
								ec71e5a9 7a9dccc6 5d6117b8 5d0ca66c 3fdaa347 6e97adcd 05a1f490 2bd04b92
								f400c42b a0c9940a 32600443 3b6d3001 28bf930f 484eaa63 02cd7a31 9ee5e561
								a12a3625 594020c2 40dba3be bd8a4751 5841f198 ebe43218 2639616f 6a7f9bd7
								434f0534 8f7f1db3 115a9fee ba984a2b 73784334 de7737ee 3704535f ca2f4904
								cb4ad58f 172f2648 e1d62d05 8539ac78 3d032d18 33d2b9aa d96982c9 692e0ddb
								b6615508 83ed66f7 aa8bce8f f0663a0a dda226c7 bd0e06df c72594a3 87c676a3
								ca06a300 62be1d85 f23e3e02 c4d65e06 1b619b04 e83a318e c55eca06 9eb85603`),
						Q: hi_(`c2a8caf4 87180079 66f2ec13 4eaba3cb b07f31a8 f2667acb 5d9b872f a760a401`),
						G: hi_(`17a1c167 af836cc8 5149be43 63f1bb4f 0010848f c9b678b4 e026f1f3 87133749
								a4b1bba4 c23252a4 c86f31e2 1e8acacb 4e33ad89 b7c3d79a 5409268b fba82b45
								814e4352 0c09d631 613fa35d b9caf18f 791c2729 a4b014bc 79a85a90 cd541037
								119eccde 0778863f fcb9c259 31fcd33a 6706e5fe 1f495bb8 bcb3d0ee c9b6d5a9
								373127a2 121e37d9 8a840330 258dbfce e7e06f81 5b69c16c 5d17289c 4cc37e71
								9b856298 d4e1574e 4f4f8515 baf9a850 d11dda09 55bc30fa 5b16792d 673a3b1f
								41512fc3 eb89452d 51509f97 4d878b48 2d2ad2ed 32be1905 6f574504 2bff804f
								b7482796 612b746f e8d70a83 8cc6f496 dd0ffc3d 95c1e0b1 98184d73 523656a0
								6431bc52 5c2bc161 9729e8c0 88f6df91 5645e060 922a4af3 edd63047 c7b6077c
								667c07d8 8eb00f4c fe59d32e 5f545012 c566516b 7874fb3d aed51403 31f29528
								b30fc8b8 a9371c28 18017b09 53a84ffc 9fbff84b 64bf0238 aa7e2af2 ecadc15a
								1c06dadc f1f2e7b1 240a5e64 5a6469c9 b002215d 9a91c2a4 ed2fb547 a942d777`),
						Sizes: L3072N256SHA256,
					},
					Y: hi_(`2574e10e 806f1c42 58f7cf8f a4a6cf2b eb177dbe 60e4ec17 df21dcdb a72073f6
							5565506d a3df98d5 a6c8eee6 1b6b5d88 b98c47c2 b2f6fc6f 504fa4fb c7f411e2
							3eaa3b18 7a353dae d41533a9 558ab932 0a154cae cc544e43 0008889a 2c899373
							ec75a24c ff26247c f297d293 747ecc05 b3483647 a87bcbb8 d4500092 09f5e449
							a00a659b 637ce139 cf6487ac a70f9c00 cb670c7f 3b95bfd7 cf236a0a 6f3c93be
							8d9cf591 c9d30686 9415b1aa 97264b90 4167850a 4794c780 be4527df feb67be6
							e66786c5 cce0378c cb49920d 855558f4 dac4c42f 92dd229b 483b2257 db0ce35d
							c737f980 1a261a02 bdf718c2 fd4d69c5 2e0d9712 b42c4897 bae7c684 d3d35bc5
							726ce899 2696b044 d722afba 78efa858 c4d10f19 72112ce8 ffd39792 49bf14e4
							9d8e0d9a cb1b0a9c a90d0551 1803845d 7c670bcf 1b066497 a7743b08 a219e764
							ea0a3a2a 617661c1 6a372fe0 58b547a2 8b626ecf 442222e1 8eef487c c101dbfb
							715bc33a b85928ec f0bd4dea 30f250a6 a5c86178 83ea0f87 3e7a4651 98c4644b`),
				},
				X: hi_(`7c28569a 94b46fa7 45c8d306 ad7dc189 96ce046e ebe04383 8391c232 078db05a`),
			},
			M: pb_(`54 68 69 73 20 69 73 20 61 20 74 65 73 74 20 6d 65 73 73 61 67 65 20 66
					6f 72 20 4b 43 44 53 41 20 75 73 61 67 65 21`),
			K: hi_(`83f3008f cebae57e c7a64a3a f7ee6ee1 9cc197a6 d5eba3a5 b3ef79b2 f8f3dd53`),
			R: hi_(`54 7a 99 02 07 de dd 6d ff 97 89 c4 78 79 ac d9 60 d7 92 51 4b d9 1c 51
					de c2 a2 4f 90 4c 03 f1`),
			S: hi_(`1668797b 26641e72 94aa68d3 8562eae3 caa842d0 f446949c 4268ae3d 0392434f`),
		},
	}
)

func TestSignAndVerify(t *testing.T) {
	for _, tc := range testCases {
		fmt.Println("Q", hex.EncodeToString(tc.Key.Q.Bytes()))

		R, S, err := sign(tc.K, &tc.Key, bytes.NewReader(tc.M))
		if err != nil {
			t.Error(err)
		}

		if tc.R.Cmp(R) != 0 || tc.S.Cmp(S) != 0 {
			t.Errorf("sign failed")
		}

		ok, _ := Verify(&tc.Key.PublicKey, bytes.NewReader(tc.M), tc.R, tc.S)
		if !ok {
			t.Errorf("verify failed")
		}
	}
}

func TestSignAndVerifyWithBadPublicKey(t *testing.T) {
	for idx, tc := range testCases {
		tc2 := testCases[(idx+1)%len(testCases)]

		ok, _ := Verify(&tc2.Key.PublicKey, bytes.NewReader(tc.M), tc.R, tc.S)
		if ok {
			t.Errorf("Verify unexpected success with non-existent mod inverse of Q")
		}
	}
}

func TestParameterGeneration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping parameter generation test in short mode")
	}

	testParameterGeneration(t, L2048N224SHA224, 2048, 224)
	testParameterGeneration(t, L2048N224SHA256, 2048, 224)
	testParameterGeneration(t, L2048N256SHA256, 2048, 256)
	testParameterGeneration(t, L3072N256SHA256, 3072, 256)
}

func TestSigningWithDegenerateKeys(t *testing.T) {
	// Signing with degenerate private keys should not cause an infinite
	// loop.
	badKeys := []struct {
		p, q, g, y, x string
	}{
		{"00", "01", "00", "00", "00"},
		{"01", "ff", "00", "00", "00"},
	}

	for i, test := range badKeys {
		priv := PrivateKey{
			PublicKey: PublicKey{
				Parameters: Parameters{
					P: hi_(test.p),
					Q: hi_(test.q),
					G: hi_(test.g),
				},
				Y: hi_(test.y),
			},
			X: hi_(test.x),
		}

		data := []byte("testing")
		if _, _, err := Sign(rand.Reader, &priv, bytes.NewReader(data)); err == nil {
			t.Errorf("#%d: unexpected success", i)
		}
	}
}

func testSignAndVerify(t *testing.T, i int, priv *PrivateKey) {
	data := []byte("testing")
	r, s, err := Sign(rand.Reader, priv, bytes.NewReader(data))
	if err != nil {
		t.Errorf("%d: error signing: %s", i, err)
		return
	}

	ok, err := Verify(&priv.PublicKey, bytes.NewReader(data), r, s)
	if err != nil {
		t.Errorf("%d: error verifing: %s", i, err)
		return
	}
	if !ok {
		t.Errorf("%d: Verify failed", i)
	}
}

func testParameterGeneration(t *testing.T, sizes ParameterSizes, L, N int) {
	var priv PrivateKey
	params := &priv.Parameters

	err := GenerateParameters(params, rand.Reader, sizes)
	if err != nil {
		t.Errorf("%d: %s", int(sizes), err)
		return
	}

	if params.P.BitLen() != L {
		t.Errorf("%d: params.BitLen got:%d want:%d", int(sizes), params.P.BitLen(), L)
	}

	if params.Q.BitLen() != N {
		t.Errorf("%d: q.BitLen got:%d want:%d", int(sizes), params.Q.BitLen(), L)
	}

	err = GenerateKey(&priv, rand.Reader)
	if err != nil {
		t.Errorf("error generating key: %s", err)
		return
	}

	testSignAndVerify(t, int(sizes), &priv)
}

func h(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))
	for _, c := range s {
		if '0' <= c && c <= '9' {
			sb.WriteRune(c)
		} else if 'a' <= c && c <= 'f' {
			sb.WriteRune(c)
		} else if 'A' <= c && c <= 'F' {
			sb.WriteRune(c)
		}
	}

	return sb.String()
}

// hex to *big.Int
func hi_(s string) *big.Int {
	s = h(s)
	if len(s)%2 != 0 {
		panic("len(s) must be a multiple of 2")
	}
	result, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic(s)
	}
	return result
}

// hex to byte
func pb_(s string) []byte {
	s = h(s)
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(s)
	}
	return b
}
