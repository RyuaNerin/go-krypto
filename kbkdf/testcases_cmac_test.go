package kbkdf_test

import (
	"bytes"
	"crypto/cipher"
	"testing"

	"github.com/RyuaNerin/go-krypto/aria"
	"github.com/RyuaNerin/go-krypto/hight"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/kbkdf"
	"github.com/RyuaNerin/go-krypto/lea"
	"github.com/RyuaNerin/go-krypto/seed"
)

func TestCMAC_CounterMode(t *testing.T) {
	var answer []byte
	for idx, tc := range testcasesCMACCtr {
		expect := tc.K0
		answer = kbkdf.CounterMode(answer[:0], kbkdf.NewCMACPRF(tc.NewCipher), tc.KI, tc.Label, tc.Context, tc.CounterSize/8, tc.L/8)

		if !bytes.Equal(answer, expect) {
			t.Errorf("failed test case %d\nexpect %x\nanswer %x", idx, expect, answer)
		}
	}
}

func TestCMAC_FeedbackMode(t *testing.T) {
	var answer []byte
	for idx, tc := range testcasesCMACFB {
		expect := tc.K0
		answer = kbkdf.FeedbackMode(answer[:0], kbkdf.NewCMACPRF(tc.NewCipher), tc.KI, tc.Label, tc.Context, tc.IV, tc.CounterSize/8, tc.L/8)

		if !bytes.Equal(answer, expect) {
			t.Errorf("failed test case %d\nexpect %x\nanswer %x", idx, expect, answer)
		}
	}
}

func TestCMAC_DoublePipeMode(t *testing.T) {
	var answer []byte
	for idx, tc := range testcasesCMACDP {
		expect := tc.K0
		answer = kbkdf.DoublePipeMode(answer[:0], kbkdf.NewCMACPRF(tc.NewCipher), tc.KI, tc.Label, tc.Context, tc.CounterSize/8, tc.L/8)

		if !bytes.Equal(answer, expect) {
			t.Errorf("failed test case %d\nexpect %x\nanswer %x", idx, expect, answer)
		}
	}
}

// TTAK.KO-12.0272

type testVectorCMAC struct {
	NewCipher              func(key []byte) (cipher.Block, error)
	CounterSize            int // bits
	KI, Label, Context, IV []byte
	L                      int // bits
	K0                     []byte
}

var (
	testcasesCMACCtr = []testVectorCMAC{
		// I.1. 카운터 모드를 이용한 키 유도 함수 (KDF in Counter Mode)
		// □ 알고리즘: ARIA-128
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`fa0317cb4ec773e1b5bb021b6dd892f7`),
			Label:       internal.HB(`43231f09aa4d88b6233d6717132c6a273f0032b81f0a07194ec652c02736bf7e163b05f3b5dc9fa29aedec2186224fca9500f3872b0752a97cdcf052`),
			Context:     internal.HB(`31b212c8928abc3c455ea68c5ace97a82c92aed394333a74463e888460df2b410d9a361764f261d6582c337eb6bd3941102d8b7d9d582f4a57243399`),
			L:           0x200,
			K0:          internal.HB(`bd6b11cd9c272b32e1febce4e7caac3ebd7fed3ae53a77ac7f4a097656357d484a2fa490e67d49557256cac9ff9b9b60f7d544985f053a9ba2dc5b187cb142f7`),
		},
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`cda2ab4a91c1caba26e755a83d8ddf20`),
			Label:       internal.HB(`129358c99a1f6144eff0462f72f3dd4bf27d9d0b43578c3286533c249c883279d8ca4d7d12d329f66c85908eac9e52a25fcadb6a0f47b65e32c7b42b`),
			Context:     internal.HB(`8a32472f8c9edcf399b896bea7ee7727b3e8870f8b4f708bdcd3e20738c1e63f6505b9848084c026643f357e2e1ac8c13a5178e5b76de7638c13c2e4`),
			L:           0x200,
			K0:          internal.HB(`9809d6df77529f7e15e61a87d88362c525d0685206eb4a0dc2cda5eecff7b18ca514561b9e6fc03dd36a8928d4c672f86d7ce8ed5cca62ea9fa1550288097949`),
		},
		// □ 알고리즘: ARIA-192
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`79b90534f31516fafeafa9db830b8a3b15781056f22713a2`),
			Label:       internal.HB(`767ab85a5c910130cf2c82cbb92d84ca22173b9eddcb9d10e9ff02b8fe4588544d356071bb410bfd831ae13cb53d2590419695dd84a6551d55beeb84`),
			Context:     internal.HB(`462a69a5f744f3f89ca81726a056d2e6cda5df668f4470825d2681073d64388a8b18a780aeb2b6633c5f8db29e2642a34f0eb74d055673f69e69cafd`),
			L:           0x200,
			K0:          internal.HB(`8ea637b83d368a34be920374810e971605122b088fa5dadbac04dd9b7a60fcbe0f92783446fce8937afac473b6c1447be626e2e34f06c6202f27122f12352859`),
		},
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`b2d776cb982fa2310039c4d2297f4529452ca6c7e28720ef`),
			Label:       internal.HB(`a97a4fc189073beefcbd7406c08d4148b8bd4f303b96fd076ba031a819db3f3088c70e2ee5462f5c99a01502f22f599228be2390fe5c22dffe83f977`),
			Context:     internal.HB(`e49e1cc48d2527ca679a918fcf0c61a128ae68de97a705f0a47a2f7b667815d1fc64bd9a5710a6b3236dfa558a55852bee1353709ffdb180a1997e08`),
			L:           0x200,
			K0:          internal.HB(`514db53eadc31980063cc381390ea5b3cd598525bb333449afc3b6b6d8c2def41a625a2173ed32e8b88907f499ec3886d21fc45a2199b00dc3a0c14c720cf00a`),
		},
		// □ 알고리즘: ARIA-256
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`b73bb3912b3be6192d632889e2f9fb4468ffe750d9a947a4ee4749c4f4953b46`),
			Label:       internal.HB(`25bbc3a5dd65a0b34726216415341e538761f11f2d0166fec127eebd61b0d94ec6032621198c75490b86566bd9a6b4e11985ceadd6b4e40c1932cb52`),
			Context:     internal.HB(`400ee1a7d91e24af85cda45afbecee21cf01a8831ce7d5317d227a4be182d83ed9f800d76465f73b53c4e1ce719f60e7fdec70a54f4e313dc1a8503b`),
			L:           0x200,
			K0:          internal.HB(`265fd0179e3279e879abc0d1f7706e325f7f21aa68abf2ec7dc29bc7362dcec6a72de4788ad0a99a2ac05094ef81dbeb687e3c526553908a838289f3cd8583a4`),
		},
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`b073c6a26c12298d151e5b224c77faf6884081f325469745594e1ff547a1b765`),
			Label:       internal.HB(`92698ad4e787b3e0acf3b47df530985f08e36caac56b6f8734318f6bcb60b9de519b0fc7394609aeb6e7eaaccffce40631c0fa9d11918dbd3ea5f9b1`),
			Context:     internal.HB(`8b03612b4431e5b0eea25e12cafed2de4af08b5701a14128348f1fc70ecedd7e43ab97ca0745e081e6ef9e2ddd852affaba2f64ab8d66b0f8f7ff45b`),
			L:           0x200,
			K0:          internal.HB(`85321fd2bf37b5334048795e8a590dc6f8f5f3a07b2cce630754a25fc610b2f14be3bdb108232852236aff23170c56dad0d2d60aa3741c71e3820f87812eb85b`),
		},
		// □ 알고리즘: HIGHT
		{
			NewCipher:   hight.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`bb0e90646aad5ba088fda90a52b2d304`),
			Label:       internal.HB(`6e4236aeea446820ed0284b0addce602595aa4b1fcd13e3983af76ca7b4ac275ce5f1e9680258c3a96628c936421c134c1c38adeb80f75da8023ba01`),
			Context:     internal.HB(`496b7b4717b02012bdc242cf5e487eff3e68b4f1e53913e22d203cad0b73de596a629a8779d5b9be060d929c87e1236610beb2ce2463521b95f2f75e`),
			L:           0x200,
			K0:          internal.HB(`36f977f7c8619c7d7dc5a9f8ce35d1ac945d12b6d9931f4566972f6d5395304bc5cf6d9841e9380c71fce4708ae5200aeee1f125ea6106ed69e21b80f902116d`),
		},
		{
			NewCipher:   hight.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`c13518447632cd6a3f2b2945c38ae57c`),
			Label:       internal.HB(`133d79a8d4524c1dba38626fbe622f27bd38727bb2b3af80794fda1b8047fe595419c3ae8f21a20952982cef4121ce7c3b6c5cff17a7bfe6be3dc70b`),
			Context:     internal.HB(`73fa89538915a31270d30e0e0be5ddfab04ebe0187008db6d7109872ee6b145d2892739d83e6d378012f4be13e11c5b62bf18fe2138f87dc79c0af16`),
			L:           0x200,
			K0:          internal.HB(`b392e26af25522e817b5b8fc09e112c3e0f92009f439285d1ef3ce1b4cd3fe8cbc6c8fff1eb49e2aad36d54ecd1235a534104906d78469c1443408ef1bc589ab`),
		},
		// □ 알고리즘: LEA-128
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`67214e1567a57e91b5b6086a94681964`),
			Label:       internal.HB(`af76ea779b393da394677000b6e21ff28bb723ba134ed07e5157e030463a0e1a2b1523578990152707bc22127a530193e22ed44094b12c3f0b6826ba`),
			Context:     internal.HB(`89a41feaa9390b92d7aa0de5346ceb50c3bba013d3e9df0b56a3e65023f3a82d78ae12b3f68dbbf0eb5bdd136c9dda4ac6a26987c1364c36aa439da7`),
			L:           0x200,
			K0:          internal.HB(`cf066e7f854c4f6e4552eaf3a09860699ca954bd55e27084094162935a5af4740637297063153d34f4157ac516f7b3a1771f40b585ef8da8f5fe85ee6216438d`),
		},
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`95478594634aa5a1cb48c60efb907a46`),
			Label:       internal.HB(`ce31e323fcacccea5afe398fab1889adb60deb2e354c8ba21918d7cf7edfdf20c285c4a4fe8734692031602b4ae3e5dddb3e2de201da7392c47f8660`),
			Context:     internal.HB(`cdf105bbebe6800eacdc52e73a9027ff6cec0c97960f3e09f3a0712eb64c7d8403f36c28957516ad0f1d8832ad165bd6bf7789b853c649978448390d`),
			L:           0x200,
			K0:          internal.HB(`61a1311f855ca8a28241a183a7e9483e78f9a759848d9be90d8bf44543e23277b8b40bdb9258a013491aea65d770d2aeafe7d0d3c5c0577837c72df9f3bac344`),
		},
		// □ 알고리즘: LEA-192
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`a9b09cce95b618184054cbe3adc370fabf74ccbd5aefc238`),
			Label:       internal.HB(`974cfa0085a4da361c4253f7302076a10738238e4312c90fe59b70543ae9d4200a1226eb01caba20eeba4b0bd4977da6b4985eb75a81e5a909434163`),
			Context:     internal.HB(`d663d846c6322c68cc0ab639891b0aeb90d78f00c6441d6564554b69c3f33b688caa6ffa1a90e6a71770b9b206313349ad5d4d3c9fe020c8dcb409f6`),
			L:           0x200,
			K0:          internal.HB(`1084f8fdebff7b4595830a6840c54056a00dea6d7383ae50e7b6327e882098e56e3e9df16dd4d25fad2b73bbcfeee5f166a4e43f8dbe3d7b5eebab86f742ebc6`),
		},
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`9aae0dcddab56a5bfb62c39de2cf19a0629070c150fae74a`),
			Label:       internal.HB(`74008c528f671fcc3a04db0ca990be68abd3d0b8267b47f3357dbdbb52a9c9b973f7b8146e8b9e83c4e209080904b920eaaa4be0c7aa53953240c0a1`),
			Context:     internal.HB(`7162e585a90eedadcdc81fc92c597e3934a7fdd724794188d30c8c3781ffbbabb403e61e32bd71360d3dce7f3fef150f504132a1cc9b4a0b54c750f4`),
			L:           0x200,
			K0:          internal.HB(`c2a70646c920ff6afd867c5c0268c0a4ead2bf8484aef8f7b76bd8c6ec9a60f4cd1d836060fb36a25befacad6135346e47d5bd62a874ceab6f8f0ade23046a0b`),
		},
		// □ 알고리즘: LEA-256
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`1f7b2cc87f9f070895e13ea8e293a8e513929271fded1edab1c9a0424ee739ca`),
			Label:       internal.HB(`58b15aaee24de7ac4d6150e1a7584b985bbb65fb36327da51bffa1919c7299b22733a9f35f7faff59147b25c8b5ecd11ffb1ecf62b77d1f84654bf0e`),
			Context:     internal.HB(`d2615f9fd80f0d20c44033585b8cd51a3d7b8bd36cc6129a4c4340389aa065f2391076b8bb56b50f3f431ef60291740316529e113c017e15a0b3e6db`),
			L:           0x200,
			K0:          internal.HB(`4b22fcfa09af9d556016a12bc35deba4c44e776ee5848cbd1969c28518d3288a1e4e30bc0d5da7652f088c3aa14309fdfbe0cb35d97c3002ea2aeaa1ad530e9b`),
		},
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`e4b1befa37be1c550e835ea6c6a50ac25857c7a1b14e31ff660ac284d3519ba6`),
			Label:       internal.HB(`1bf503cc9301fbb548f05dea0ce3c83c230321e4d9a4af408faac1622dde2d19cbe653558a8ed0698bb06a1849e6639cde66b76a1556b8aa6cb9c672`),
			Context:     internal.HB(`742f4e4740f9770ff0430059bc7995bfa46544e68243d02a875c022b623084ffc5559254bf38ad5de16c15192741eafbe5f6257bdccb736129801b50`),
			L:           0x200,
			K0:          internal.HB(`eb02635f1af1b4a22413e408a5b058adc993f215a9bfa0a3e76d339588fa04828b395331fd7523f73f438c9bffc368dff3d8e2ba0f334ff1453f5f3472d9e3b2`),
		},
		// □ 알고리즘: SEED
		{
			NewCipher:   seed.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`d899a261ac0bf40b04f110681b8859e4`),
			Label:       internal.HB(`8ed363214d5ca03328964b636a889f83850bc41af69ce426d080490cf57465514ee2768d47fde795fad82ca919ca1c9553122d8017549d431ddd6b0a`),
			Context:     internal.HB(`b78e81454183e6fea2de1e842e0b782433fdecc33f70b60279b9284fe02f6df2517faffd9ac52c1ef442fc1145deee1edc53f55ffcba9ed9fe1087ca`),
			L:           0x200,
			K0:          internal.HB(`762dfe66668c7807611e8115156277e595597320b22c2e25881f1bcc292077166a39dbe9aa75e05d57d1f5d69f286eccd3d7e0fef68b3e1e88d5cd26f85ab7cb`),
		},
		{
			NewCipher:   seed.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`54445b756402f067c3955d70d3bfafdc`),
			Label:       internal.HB(`f2f63d0be2fc8a4b26f1f723569f4181f3a1322d4b43a2cef004cac5f03672321302f4779291e2c98cfadea0780fcb6c62d0bb8e32ebfdcdc3181c11`),
			Context:     internal.HB(`a26665d4b8b6ccce0dccc79851861052cf78449ea724d53a3a780e8eb35101c8708bbfcb47938d7a1806508b4a4960a781eb6d89fd730869b8fd262b`),
			L:           0x200,
			K0:          internal.HB(`699e7450003e2ab7500330eb854e09c4ecd2d87cf54cd6b780eb50b36d1756c01391a1b839e090c960c183427abca3954f8d137d08cf4a7ef8cd7d5da7e2a7ca`),
		},
	}
	testcasesCMACFB = []testVectorCMAC{
		// □ 알고리즘: ARIA-128
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`33b6fd27d8744df0db928827fa6677b9`),
			IV:          internal.HB(`d42bcc8bd0866bc5e6758979753e54e9`),
			Label:       internal.HB(`b9f5de570139e6eae8b22cd888b90973009fc4be979934d3840f726a09e5a8689b0c18b5603b619cf61cddcfd02383e8132bc938023b90dc5cc2d3d4`),
			Context:     internal.HB(`cd6424945d9af9657447474f11d0d13a461b2c81151ac26c087376785965556dd005a171908b0bb86645d34a9cc5296dc67cc656102257587b8ecc58`),
			L:           0x200,
			K0:          internal.HB(`d05c4264067051fe53c6381639f281d0ba124e0a4ce8fd49db399ef498a1cee996cf553f84bc37bafa3ae5b789347c447e3b9eef6511692c361581a9d9713110`),
		},
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`c33df77be674bd658abf9ce0d17071ed`),
			IV:          internal.HB(`c73a79be0fa9ab5b0dd3d0016b998ab1`),
			Label:       internal.HB(`f029d930b7902bcbdb5acdd310642b47edf530761e53cd29001c1557869734bd531957233482ca435f3fc09deacb9556f4078f8a3e51a84779512183`),
			Context:     internal.HB(`b1ee0ab44fc173e13d9a2c68eaaaf40e0be587343823f6847db0026ef7125b51bc6fbae6d075c02fb8bfc90a0262d5bcb7c777ac1362fa68f8aa7bc0`),
			L:           0x200,
			K0:          internal.HB(`50f65b5c84d7c9f92bc31b2be82b0c3f4d0c02bd049f7f8677d243742dd54154f26c45a4df88f0e211737411aa19248272b51b02cfff8e1d87b498035053bcbb`),
		},
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 0,
			KI:          internal.HB(`e8a5af1eaa3be23ce1c2374ac6ec1768`),
			IV:          internal.HB(`7b19e15ee649877a35cd0778f4fc8059`),
			Label:       internal.HB(`f6f418a655d582e5038888a6d8c06e83ff9dbcf21abf46340979ebd412a9ba3911b46973fedc6487a0ca80248b34c6367a3e8027a694924ca3370f9d`),
			Context:     internal.HB(`63c883639176ffde77cfbf50b76ea65abfbe0474ee358ec51b402bf0fd5a055ffc7abe48563dd338f3d5ea32ff07a9299c537df8c0d478d21be1c317`),
			L:           0x200,
			K0:          internal.HB(`528d12498323c182e8c3a2b8876aeb7cf6882f892efb64113c746238b68874369f4145a50151ca0bac1ca86c7833cc6b293250e5f3a6bba7990036ec5c7edfa2`),
		},
		// □ 알고리즘: ARIA-192
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`a6f3ccef5e8445953a585592896d6bab35190dd2e57f90fb`),
			IV:          internal.HB(`cf0af591d90434b68f0b68b012147b24`),
			Label:       internal.HB(`857af3e3819d6ece674e2f1915886acf59073f58b7c09042546591f6be7912185c7d4bae2f29e68d4967944a7e18841470767721d753188ce28a655f`),
			Context:     internal.HB(`2787be4dff233ab3cf65ac982da7613c10ac3d243e1896949ca92169b3b4d0bb3733449402296451b149f81d44f2f1653f9402545e61b354c5097f6e`),
			L:           0x200,
			K0:          internal.HB(`ea66baea7288e1d1441482eb6c9bd155d9e38c3c2b5618b44726c468d563350d795ba522b66c2dcf62fe3bf24520903e6bff5ddc346cc0a9597c13a8ca7929d5`),
		},
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`7f6630090c078be6d2bfecdcd88182918bd0fd6753e1d570`),
			IV:          internal.HB(`90201cf67a2ba19b3ee9fc8f0d692342`),
			Label:       internal.HB(`3705136adee8c83705509495f9082540ed7619a472e8081bbc5eeabe5e03c55bafa866aaff1eb6934393555b307207a739e82aa0918b87e369679d86`),
			Context:     internal.HB(`df34b79b89ec3e694c66c1285998bd5effe8a0bebc364d99f5ed7b42cb99a737f76babc3572c0ef924631fdf7af2cd658765e217f565eae823433471`),
			L:           0x200,
			K0:          internal.HB(`be8456a62298fc16bbb4a75b614ff95261da45c6a8f003f0f2cad01d2bac34e3275ed0f1bd9c83ac797ffd3f01dafd73071fa4580408a239d83adccdcb470178`),
		},
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 0,
			KI:          internal.HB(`401dd9923cea5b6378fd04ce04007d3474a11f9d69e9acfe`),
			IV:          internal.HB(`873250ae15d36f6d583966680c76cdf4`),
			Label:       internal.HB(`a2b2ca637fade4fbc3359c4ee80888524d40e189d456f507d66da7651bd519346786e52e0c49f9e7fa47763355b90c60230d6fe4ed4c681a23f07f96`),
			Context:     internal.HB(`2b7d460075a11839e918a0d298b62b73a20b08b106f858fc36bd65818e30d33b69e74f99f8d3aea2b128bc5bf2c4f2c9334dda9c2e7a7e8b092572a3`),
			L:           0x200,
			K0:          internal.HB(`9921509cb4397223d33f82bcb9ff3245b5be1f599b483831053268b98443f7a689c2a7d19ba22b0ad885af28e2d4f7f1e2e9dc9f16e1618c837b5caee393884d`),
		},
		// □ 알고리즘: ARIA-256
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`6f365b5f220a854f65a285a29ecd045e5aefe0b90e019b09fb47252cf753d6b1`),
			IV:          internal.HB(`44bba33cc2055a79cd4db3537d20f417`),
			Label:       internal.HB(`4cb4c8baefc3e1fe6e8ffccb6891610a087a58bee2deb43beb0a31cd1a1721736a8bfbaaf9a2c5fe0dcc5fbbdfe3bd40322ecf5082a770c3cc4f6d0c`),
			Context:     internal.HB(`35694629b75284f27001b52dec900980e4c7655bcc3fe67cf2fa537b6f31109bce9bb5b98dffdad4745d7a094a8ebdd4d049209b6b5f2eaa61a23dad`),
			L:           0x200,
			K0:          internal.HB(`16e1dd83a9dc46f0207962e2226e1e392dcfed6ca1ad45c71b32db22c7bc4399a4007eb976d54297f7a43465cb821e94de290d0a6809276f1e544c85aebf4d92`),
		},
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`a4f30aeb37a6304f9db28d3631884873babac49dd24887d9e21a8b65fc73a8e9`),
			IV:          internal.HB(`7aa1c9be2b93a1b49a1691925efe63ba`),
			Label:       internal.HB(`6391b40f0e805abb8f3137fea4076dfbebb76f3fb0f14d9a52a11dab4ead116ee6f7f58857025de4eab43fbd6e2da8a9d00f45f777de3d10595c7c01`),
			Context:     internal.HB(`8bc92aa3ead3d9fedcafb7165ea967256876d6cf3c22d0244128a7991dd9bbe42d5492b0d39a44aa2dce4ad64da007615de18a56541078075b081f67`),
			L:           0x200,
			K0:          internal.HB(`6e98abe64e5653120f1c47503f6ce0200c9b3fd878f26caaf11af372b155c9fe54ad7850c21c1488bffd9b92df780d23ee43978818019bcfa84735e80749a151`),
		},
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 0,
			KI:          internal.HB(`539d74308654a7f42866f174cfc5019072b75e261f200bc178c02b96077c47c6`),
			IV:          internal.HB(`6f100ea4530477afa27f9cd7838ab324`),
			Label:       internal.HB(`2e49fcc88af112e847d3083316c7d749d9f1b612ed74803ee787ad128927795b61da6ec3b5a19356e5768bc6ecbc7254cc0889a2e983b648d8197a8a`),
			Context:     internal.HB(`7c32f0fac2bd47bc63db832b0adff6e92abfe5fa6e084db2daf08ed8a8bcdc917bc8d49848204b015cf8016c1ed15f89085da1932fc4c60437aac26d`),
			L:           0x200,
			K0:          internal.HB(`7497a82466256ec4e79d9b0f5713fbd8a9f198e434f8e20b5d607d957e57ea50feac1414c9cca4e176691c7eadc45a88bc77f52323f513911c3e248f15be21d1`),
		},
		// □ 알고리즘: HIGHT
		{
			NewCipher:   hight.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`d11f43e61e5dc7208b9e463b0bf1dae8`),
			IV:          internal.HB(`1d57ef5d80ad66ff`),
			Label:       internal.HB(`a228b86372682476c6b246c264e8ce5d946a4346d2a92d2ad05c623304c5134e4d0790150ba79140ee0e4a25f9f359fdc4ea52023ca75e7ea1815092`),
			Context:     internal.HB(`52e59ac33b9bd8ca39fd5d22b0631ff5048ab69c21aad1246286d7463ed260600932fe8683edd51bf0fe41f8e82878b6a3c3a9d6580c883fe7619544`),
			L:           0x200,
			K0:          internal.HB(`9a4c467868fdecb0976622e7f24ae3e7604c90de037a038572b4c890dc8fda1fa6fee5590c5410d8a3df921430bfbac0ea917298b51e2f0b4a12b6eb8e9bd17f`),
		},
		{
			NewCipher:   hight.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`a9176aa3f756ad8db62b67cc20eacf06`),
			IV:          internal.HB(`5ebe826bcfc1e940`),
			Label:       internal.HB(`c46a4a1440f70f676b6e637f74d6659389a98f8a9c7931aaf8bd72f789e549ef1a89ed040bbad5cb3489570183987a0be4bd3df8079d5cd27a813f21`),
			Context:     internal.HB(`50563a8091d297ce71cd61fa9f47992c7bc3bc1287a2e42b3448c9501782f3a24b8a3503f77a1928b08d2f2bc07c99cd55e25891b9baf7dcc3597401`),
			L:           0x200,
			K0:          internal.HB(`a6a66d9c63f6f02d4742fe8e578b73096b6c6b43c75fbd1a7a0a7c758436f2d3cf5ade6c9aa11edc9c0486cd162052ed3807669553d82d57c866e8f2596bade4`),
		},
		{
			NewCipher:   hight.NewCipher,
			CounterSize: 0,
			KI:          internal.HB(`e829c6877cd71ac0a2d8f3d33c956c96`),
			IV:          internal.HB(`dde6850e4ee9e2ac`),
			Label:       internal.HB(`4595366224dd4a17b27466a6b570f9681d9889bd922b605b98301e66a87f714ab35f8ac8d1b6a064ff9e18bf24b42eb6d47658e333db6e24574fd30f`),
			Context:     internal.HB(`a5480131327522e777cc7b0ace19af866e092b22dc3d061d4989e7cc3a8968b7d659f91d1380c74c82dae4d27ce26deb83bf9b32e2cd263133889339`),
			L:           0x200,
			K0:          internal.HB(`2ae235eae847fda1b53dad269ab1392a0874867852d55bbae8cf3ef1d8852742369ed44f46cc22e3c7938acaad0341e0c97d349d41c54d756f52c84f16ddb5c9`),
		},
		// □ 알고리즘: LEA-128
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`cfdf648e23ce688d7bf1075bda9783c3`),
			IV:          internal.HB(`a9722a6eb6321a45bc5f84944db1daea`),
			Label:       internal.HB(`cd50d19fe31cde5ca2e031948864aecd9fcef452f229a621f949952200e0416b413c14652a79dcab21b31021c603678553d78be06679082b8f597b9f`),
			Context:     internal.HB(`dea4c83f1e225dd2a9bd1f308c6eb4394ccb42ddca5aa6089d68444df40de4619569c41c03f2f916b1a085c2e33f0427a801b2c0d68027398f0aa0fb`),
			L:           0x200,
			K0:          internal.HB(`21adcd885f0a06b4d7427c865f9b3c9d2946e1691c1d4d087bad107d8ff169ebee26e7bccc128728b09c66cf55ba65fce3ba09717e0295ee3a1f0d6c88f2d6e1`),
		},
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`32ee82084eee5b25be0fa3ea4eb56bd7`),
			IV:          internal.HB(`80b88108080d5f79e9d53a34b89a9333`),
			Label:       internal.HB(`1b8f9aced687d6f6ae06abce680a151a89429be3fea44cb34d0088f8248b0eac8256ee4d4749537fb9dfa6b42939e3390d197a00396fbf74e5ad4fe5`),
			Context:     internal.HB(`fe9e0efa056ba5a59b076fc5747ad936de64a32fa400d71d421f974f62363bd78fa38f7cce6833cc4f6035db5ba745d2e70427e06ba6bcdd03f63c29`),
			L:           0x200,
			K0:          internal.HB(`ea96dc1380ae1a0bac70f9136198ee84cdbb12fcec4e11f2f5e368a039049a5d8215e460121b7b4d0ec1dc4e69a048abc92901991afa12fbf9ff35f6fd7e93cb`),
		},
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 0,
			KI:          internal.HB(`1ecbebec87f1f754d9416eb4a2dc7702`),
			IV:          internal.HB(`87b1b70df2da3ee17e7608190424b910`),
			Label:       internal.HB(`da3d0abe8cf518a2bd9d98630086d0e14ee2bd2fc7239b500eb1ade46c36793600792935e39a881ac2b9a9aeab15021bb4226adf18bcbe924454938a`),
			Context:     internal.HB(`510cf6919cae28d3c263365a2ed79276c75e7be269508897f24fba2b44c85d35c1065c3f084c455df532b0a0a8d1813700946d2f16472e47eb2a32cf`),
			L:           0x200,
			K0:          internal.HB(`8e1c647a5115d0622160b64a09f6eef90d7814da22aec17d1be917f9ad736967a7b6cdf89ce82e0d29b94d1da642164b82c221484c9bc52e13473af29c46b906`),
		},
		// □ 알고리즘: ARIA-192
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`6fa3e9017b2d899fdf37e5d9902814eac095e7db28bc503a`),
			IV:          internal.HB(`9213b1dfad5eba06db1a7418219afeca`),
			Label:       internal.HB(`1dc9c06484aa38ecced4e08c52574136a31ab4bba89be4a2ab890e8bf9fd7ee4f32e3d31f378acd1ae4aeaada60d2172b66161b60808f19bbe6b4d35`),
			Context:     internal.HB(`dce9274162de89a748259e767cb82bb385a2f8af40d632be9705f8d4ba53d1b1c43b4a98923dbd1fd4c5a7b162317a3faecd99f078d0b855da20e0e0`),
			L:           0x200,
			K0:          internal.HB(`b88de1e00f43670352cd6e32467f171165f14602d13d7f27151ee314a270c13adb63b2983850b32d9fd76cbdcdfb582f7426ae43e35006265042b03d189a44ab`),
		},
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`ff848eef864f1696d2e5dca374480a3ae6fdec96f755cec5`),
			IV:          internal.HB(`9b1849c351510ba1121214564fae7226`),
			Label:       internal.HB(`5e9803d12ab771fad285d451f729d889fe174ad220569110851664e24f3dd23f60ac7624e5ecc68beb3abaeea5fc94fb3fcdadcb533132425f55b6a0`),
			Context:     internal.HB(`c553618204b332a382e522ee21633a9304c87ef5f5fce810581e10335c1ab062f3f75e9ca8a11ff43ba97e465f6d791cd04014c024ca14fcbe20c043`),
			L:           0x200,
			K0:          internal.HB(`397da6eb0fc5e2ea863b4ab23fb5916403c28cd58fc29879a550b108e9685ef2b3420dd9695dc3d42cf65610d8786c4f14e894adae3e6e26a85a4e685f1fd0e2`),
		},
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 0,
			KI:          internal.HB(`f7e75b0f4a57681e06834f581ccac0263f3aa297e52c153f`),
			IV:          internal.HB(`d9f11f46c1c49f72bef0a660f8fa1da4`),
			Label:       internal.HB(`f7df699ade769315f36a3b75c7b3fb730b03af39329d7245f9b32520510a0059650531043a668df3abb0db0ecb854fc6cd70192b46649a0b9bcfc989`),
			Context:     internal.HB(`94373539bc462e4abc510baa81025eddf518764187c0baa1b255f9f567c90abaf5393eb252feba83f9d78dd6f3615594890bbd17c0859a730c450c73`),
			L:           0x200,
			K0:          internal.HB(`c2dfe8718c23b84d9b30feb1a76d219f7165e5d5996aa71ef402af6aff71034cc49e65fc3b5c78a3cce39eb5f294c75fb88611f835436c855bc3e0a08cf13e95`),
		},
		// □ 알고리즘: LEA-256
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`7fa2c46838db8a006eb35136639dbba3fc849ebaeb8312efc0c952aad7d75a35`),
			IV:          internal.HB(`104d35a9d9125d287b6342bab808d04c`),
			Label:       internal.HB(`c525eaefd675d4d181e629d6f7805d25d1dbabf23e4d7a81c118b5598f139730992f7610c4cad2e1ec71bf35fc3dcef7b3f1f04934350b1dd7b580d9`),
			Context:     internal.HB(`6e64d1ad170e0a44c5bd833146dd8393ba49d7b24f85c89f228f5ce71a97ed7f5eea81050d6e06c36acd18a8287827723b85916ac6335103a801a660`),
			L:           0x200,
			K0:          internal.HB(`67afd0ff64b67e7e3a3f234d5e7937344a2c6538fe11ea87327624a396640e7301820b42ed2c0dee5f675b88aa0dd815d4d716bc3568aa76e573b164a10474b8`),
		},
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`48b03d3cb04a16e050df8d9e9446a94bac2b4abde688263295fe6def4c37c76b`),
			IV:          internal.HB(`2fccd12d90e08db136c203f1c58359d8`),
			Label:       internal.HB(`7645696bd92c9e2684ef18136594d2e995bd45866e5cbdafa268a19285d38fa7de5bf3feb82366f37c57e064199bfa8f6f09eb37150e138565cbac09`),
			Context:     internal.HB(`1eca4dda2fc6ca02e222141c7ff9f87b703027fccf851a02b051e93f1c80a08e5ee5eb5fc571a8bff24425383d4a89ea0a7fea14a62f752a3a2c646a`),
			L:           0x200,
			K0:          internal.HB(`32a4877c6b62987c429ab0b9b4cb91b2a42922e7aa9d99303c7c51a54bf26ba647df80de29fdde30c4ad26a184250c6faa15396517d8fc56e4270b02cf2e62fb`),
		},
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 0,
			KI:          internal.HB(`b761eca8162944bcc0ba40f102719fd128d101dd6dc321102aeabacb5e7a222a`),
			IV:          internal.HB(`f3c4656c25757caf4589b21d4517a722`),
			Label:       internal.HB(`90cc7673b52850765b43565cd28bb7393c0f3fa2f89ba1f715e713002ee5bf3df02409eaab64b4b76585fc8e3e1b178bb66fd1d91f1cfbd288456350`),
			Context:     internal.HB(`7c950017d68bf8c96e68fb087c6d445d54813433202e0df2ad9ce744ff7cfbfe921be113bcba9221abed2695e92778e5b641e77133a1516ce0b3484a`),
			L:           0x200,
			K0:          internal.HB(`0575d031ed46efc1992b7ae3532117b6903d159444ca5cce2045ddd8dea63db29f2489a7781260492ed869a7ce06bbf7bc57ef7de0dabf8cee0becf7e6cf4e36`),
		},
		// □ 알고리즘: SEED
		{
			NewCipher:   seed.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`c0ec404a828088e5f79e2840b4dfda41`),
			IV:          internal.HB(`9e3ef8ab426d06cffccf9e7b016d4981`),
			Label:       internal.HB(`67bd99e57f8b859fd13b89a0250b945a9ddb7408d85e8008daea9cc1a57305e639f1550f5e2a8cbd9dcc368c4391759e66bb332977670a237cb60753`),
			Context:     internal.HB(`ad7d03f1d618fe349992a3be3a3c8641799da4f937169e010fcf738614a6ca2db04312c1f592d9fbf409f789227ad5c5b4f6c587baca990fe28b8c9a`),
			L:           0x200,
			K0:          internal.HB(`3c6dba1fee26203c0d221826148013bafe5ab7886c59420e10870ff241e6a8e0525b31d62d6907a4233104c03beb06b9511249ba4b9e327d0f56e546180524fe`),
		},
		{
			NewCipher:   seed.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`946e187f62347ef6f71714a42174871c`),
			IV:          internal.HB(`9fa4e7e0bffd91188aa50645084eb48c`),
			Label:       internal.HB(`c804f05f5674aa67b09f9959c0754e87697744f9b554f3dd16080ee75235f890059f0cdfcdae46610eda6a8ee21aa8d24c8029c6172f2ba13ca9330b`),
			Context:     internal.HB(`e42ff85f49551f54deaf634ac9cd79e0d3d9da67bcaeca667d26254cb8e3b783e5fe825d8ee2f710cd83956dd68b076089e73aedf9b6c528a07532fd`),
			L:           0x200,
			K0:          internal.HB(`2e607bf7204a8674017b6875709e9397f104a80e7585eba571b0ee749c3928f937a95a5b0e596e13f2e11c7e02ba70915f69ebcf6bab812654d8dcc731fdd2f1`),
		},
		{
			NewCipher:   seed.NewCipher,
			CounterSize: 0,
			KI:          internal.HB(`66e187fc7047253206d122f4cee6d62a`),
			IV:          internal.HB(`c633570173079aa2edf10d9e9190aad1`),
			Label:       internal.HB(`a468d915f7a5788625764636fb61fb1d2e9cfe8f4ebe8a99a21b8c002f3dc2acdef7d6756418f666e00ca9046fb8af2eb1e41c0a9f9b218bb2a12f40`),
			Context:     internal.HB(`c47374b18963c75362e80b958dca631c1565b6f67c03df6b4d819d0cd948b5a6390b86468b7c1c00154dd9e95a3e8b67ba4c68a1d5635134449422e3`),
			L:           0x200,
			K0:          internal.HB(`e127cc1264a0367ed2f797d5980e91d2a3c79fcc765be9358a6071bb34453acda8a2245b8ebd436ee51d64dc91c6cf57f8a99699307b8a8853603b6aa705bdb6`),
		},
	}
	testcasesCMACDP = []testVectorCMAC{
		// □ 알고리즘: ARIA-128
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`e9d69cf11d994561cb44a4d287f47027`),
			Label:       internal.HB(`4742592ba10bd2fe1c861e140aa689c6d08302fea12a51aa0312e18219ec7a536059e64de512cd2d80e35434c9cbb6cc31b89af357a7cd2a61860aca`),
			Context:     internal.HB(`78e16d721972459c6538e0981bb2edd42143700bf0af2e11612077e42bf9e1cfefda6cac5f6f1fa86a39eb9dc2a2c85de90065b96f8adbe806787f19`),
			L:           0x200,
			K0:          internal.HB(`25ef9ced759d0dc8b2b1a1703aa670c2eef56807fde63349f37837ddb867fb4a82a144f24616cde7c85ab4bb7d46bf6db948b9923366bd7351e4993d4333f284`),
		},
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`5674f7edd70f68c88ebe43b8db09c49a`),
			Label:       internal.HB(`463f13b7e9de09f664e5cda1491b314b17c7de7536f797b76f3b2f56bd85fad8157cf40f830840216aa05831f4a9c29fc55ffc5219160a890921b150`),
			Context:     internal.HB(`1c7dac079781f8cfbc3b9f6635837779bce6155b67d3b3c00dc9ec6a408bbe194190d6af54a892f9e5a181906c7a9d5e0c7127c79c25fb2d1a84b4c7`),
			L:           0x200,
			K0:          internal.HB(`aaf4e7c3b237b162af54f351835f587f361b28aa365ce84009e0ac1acf50edc287a0bb0628c47e9261f030bc5bc1e6628d390231cabed57357593d0700a9a162`),
		},
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 0,
			KI:          internal.HB(`b531b06fbb2efd59d6d291ccc986b5be`),
			Label:       internal.HB(`1501a06c5146b0eae2f32124a6115b1733bc37c72064ce26a4f2d84429211bd78c2f1aea3bbdaafff312361eedf2124595729dedadac2417d20d6084`),
			Context:     internal.HB(`169d76adc88e14d3a97cd0a26062642544115b498556ac345911be9d28a8db5602e7e5974a40bd71f825984ab46973d230e50ba19a51d1a956387498`),
			L:           0x200,
			K0:          internal.HB(`14ee0e2d320fab318092a509c68b4f746f74633c575d9668c5610a143eb962854793bcc1c8d31d28964de3d012af55068b41dca937a7dc1be67973784ff87318`),
		},
		// □ 알고리즘: ARIA-192
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`224c70fd148861536ac4ac04b3636517c0182a1d5f71a4be`),
			Label:       internal.HB(`8e1c56967ee0e8f0c77d54eb17f3d202f74b5bc49fadaed9fa54176edfda23f714b173ad249fab163fe4efa3b845360f63e39a228013c09a37a566ac`),
			Context:     internal.HB(`53054d5c7d0ecb2bc2467f226fd2689a6255f0f08a56a09aef2cf231e4d6554531a93ff5f17d15f74aaa3ebb29521bbfba238eac30531cfb8a45d6f5`),
			L:           0x200,
			K0:          internal.HB(`77feb8d322abc599c966caf4e38855a0f10f67ea3559a6f7b643e5f08fbcb749027105280002754dae3c0de9b86bc9d79ba0d2199a11696b24060117b75322b2`),
		},
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`aaed815efb067a067fdcc3a0fa6c8e020566b5055b0920a3`),
			Label:       internal.HB(`4513abb164a1574c32ed29e6b07740e326739d54650c540c5e0afbd6d6bace91c8f42c425456324728756fba904949196670ceab8fbbe4ca671a2fc7`),
			Context:     internal.HB(`f751948640e38d0b3a4c97bdb76ecd0c60001f41234697a0a83783ec395e6a872cd8c69569a8d8672d79743f3e1402ddd7e8be32b51f7742e68f8ee3`),
			L:           0x200,
			K0:          internal.HB(`39c63a8c59075726b5cf46ae040a57e737138d10278be5c2124b07e51c5d42181e0d91b81a8dbcb7d107c1b62c42448bcfcfc1e5ba8fc128c841ded55ef53457`),
		},
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 0,
			KI:          internal.HB(`a83e142bdb7137cbfeb1b9c7a3dd27f4c6cc8d83cadbe69c`),
			Label:       internal.HB(`b4cf4a4c9e31e464f0ca6c71ec1622d4899ee20a991bf29e7889869d3f166774b60264337f97268590689a07590cc8ca14d9f8e592928c4739e12605`),
			Context:     internal.HB(`9202182754f02023513274c358c8af9fd34aa7f34deb75ed69bb2fa243ec652eb3bf5fa18bfa7e530aa45f3fac85b5de491560cf7e1530650316b6e5`),
			L:           0x200,
			K0:          internal.HB(`33703eabdac09243cc6260a3b17621952f5be2a76b88dfb8325ed1d013282cd143a48a14d81cbd2f51ad589280dfb2c7c1951261caef15f25334872419051703`),
		},
		// □ 알고리즘: ARIA-256
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`736eea46389a3c0bdbc95ca2bc2ddeaff53457322aced461f5994c7176db6c96`),
			Label:       internal.HB(`97c545be0b5faeeebada5ce51d046c39210962a9acac6599ade7128c359db3513bb130d8d738317a3dc03a160338da84a6579a7e470e16676a23fb04`),
			Context:     internal.HB(`7aafe2abcff9e4e4f349a6818270df2c7343660f4cdf7675f459aedc11acf0fcdec23bd482befdb2da7c4200b9f85554511e54b8975aaa6963ae6d85`),
			L:           0x200,
			K0:          internal.HB(`e8438d327a101f32f520482487753e7ad0ec1430686e772e39e145933332a0ca44a86a550b2d10ec99d531635a30f55a6b692b9e3518b105e03d4370d2946ebf`),
		},
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`6d7be9d5fbd58a3d23c3f12b6cb425a29ba5e5ab460860db6ee775f04a8cf628`),
			Label:       internal.HB(`d70632637b70d400d0fa1d0f0b3dd15348098a85394bbc1dd5480eaeaefda22346fa64409f37e14f7eab965a5683e5c9772c886629120c379fe73ed1`),
			Context:     internal.HB(`9f55992366608b2170636928aae5716e781f08e819f0b65927553942dd604ee5ad8e871067aca5ea6980d487c09d091f2a7a4c1e3f449bb7fec8ff8c`),
			L:           0x200,
			K0:          internal.HB(`97b01e6927b1bedcecc022f0e1528e80b6983fd74bbd4a9bb44a8d324577f93b449e16c14834f10921faecc0ec2a37efc3f2a1a1e67f8975271647eba924ba0c`),
		},
		{
			NewCipher:   aria.NewCipher,
			CounterSize: 0,
			KI:          internal.HB(`dda4127e9ff8b25ebb021c4fabc5436818058e599d96e0e91bbdc68c388b2955`),
			Label:       internal.HB(`2cbd26f5151aca0f9d724e0256372a45f653eea2bef986f92a49079d6eae33c59d32a2a655084cca5d6c92acc0c1b386609bbe9f383d5e5960f26c0b`),
			Context:     internal.HB(`c2bf0f7e536707c37e6cf809fe0672c079cdd957f89212e97da08c4d204083eb15429a7743e9a2d7894d23575208edd5123a1aad3066c2ec483ed87d`),
			L:           0x200,
			K0:          internal.HB(`bc38752832938636f19a5727010ea2d91f0eec3030cd39d0ab3a7e7dfdf498170f39ff8cc6e1cd680d539202104b57e281d5fb011a8b07962272696610ed5b77`),
		},
		// □ 알고리즘: HIGHT
		{
			NewCipher:   hight.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`be2b8832ea5e6ed664afdea6040b5fd1`),
			Label:       internal.HB(`ab71f520eabf37cb8448b49c4b383faaeb989e10de0d20ebe41b1267f9209387f212d50a3c10fa0736e9256666205a95556f17ccb8069d03cf403cf5`),
			Context:     internal.HB(`c7cb3b9c0d11c98f49a89e011daae04c9ec258d6dd59d3987592d178d9877761f7f7084800c184ac7fb509930d2c4f333da7a8bb9fadc87e62c63fdd`),
			L:           0x200,
			K0:          internal.HB(`790a7716d70fedff842cbc8ee6ce46e8b51b7d159a9af177559474f90242cb3a239cc4ed57b121578d6ccb696e2ac90fcf522dcab6ce07e5d0e5a7f9abe96d07`),
		},
		{
			NewCipher:   hight.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`4c1d15857641d72811b3783d47649f52`),
			Label:       internal.HB(`7c026b1aedfa51837c786a8b9c25db7dd06567ac2d7ecb12a36f54bf1a3919ed83e627177649805860ff3cc7451dafcab5b4eebff13d39459ad90780`),
			Context:     internal.HB(`2de78ba52ce0359cb07215c1905fbf947c6a304edb763ef991f1bd6c67b19abb3378275c47b678a4a803f69eadd4236ccf48006871cad95942de2e67`),
			L:           0x200,
			K0:          internal.HB(`c45ddd2beeea3401637804ae45b4c39bc20aca788e6176415061a7f9b4bb29efd9d3d75e679464fe7224102dfd8778136df4f92cba1bde57dbfd5f453ca4c7e5`),
		},
		{
			NewCipher:   hight.NewCipher,
			CounterSize: 0,
			KI:          internal.HB(`3c3b799a4ee15f23a06dc5e4d1cf5069`),
			Label:       internal.HB(`1fcb1500e3100c4ed4b93d1e896b10f68650cae5f11cc9ebf505514f21b240e65672d4e0863ecd020622c372fb428a8c8f555d4721f948ea6ddeec67`),
			Context:     internal.HB(`56d0e5fe6c0c37db5380856680a5b4e13eab4c0aeca321985ad0f761a74bd3f4f341725ff12ffbc271b2e391f4e7927dbff4ddf1a66ff571460558dc`),
			L:           0x200,
			K0:          internal.HB(`8dc929b4099e7694c29fdba26f5eb70a887f76e71a3d5773444d2ce0d9cbda6a29c315bbf7b40bfbe6433c8ada47d1a2b599017943fa3226693c5ee0226e4f17`),
		},
		// □ 알고리즘: LEA-128
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`85a88fc16c0eb2acbc22695f510e5804`),
			Label:       internal.HB(`7be6a42f3bc9bdcb6324f214287c4482d043d760bb5049f2c8668d7d0c0c5354ccb30022a9194c96b49d51df0bb01b6e0a17115121c7e8da56c6aeb4`),
			Context:     internal.HB(`701473bba542578b8fa9b5e6539194b638b707c4dfc39271174c32c9fec6ad15c7247dfaee569a5b7f38ee942f0186fd872674f1a27a0267d1efd50c`),
			L:           0x200,
			K0:          internal.HB(`61bc313ffb4bc23f13fe894e2acda32c216118017595fab104930b761c5f012988f5954569c533dfbb784267f92bf38fb269488e7cf659027311f8722ad3409d`),
		},
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`87ad87c1054fae0979feb1db6f9988ea`),
			Label:       internal.HB(`40f17f70f76444603bd21252c0f78f75f75bfa3b372f03de5cad8a94ae2b6344a9962565bfc87ade2ed16b51b68f76e19acb48e9e48a8ab710b798b7`),
			Context:     internal.HB(`972103e249eb14f5044e62c7ddccc0939a0b18817da0d86385a2ef3a81cd09d75825ed745c3697afc1d184e6aa3db29518a67e7b8aae7756e0d08c25`),
			L:           0x200,
			K0:          internal.HB(`98cf5fb4e762b22db82eb1c73f6228bba4374af3abb43c37a14f6a5b7c77f24c56695ca63482e398470003ad9dd166be8729c6e6417759bf42dcd4b1df246728`),
		},
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 0,
			KI:          internal.HB(`1e9eb5aa7c61129db199eaeebbaa9715`),
			Label:       internal.HB(`b3f4d5779c72b039bcbd50c6906f63fba7d470230011d5877314ecbdac32215d35d87852e552e6c6a396b2bfceb1c16e004af9437846f9dc3de646d4`),
			Context:     internal.HB(`824835e9a9e12b59a1fc5384556cb0eb015bb4a9db89fc2b629ceef6de6a07a7ad830291a9879e5867aa249e6826f94610b4dbbea5455d12dadf9749`),
			L:           0x200,
			K0:          internal.HB(`3402d47b96efef63d66135b50af0b0ab76209d7db6dcda909b1b2597812b844a36d0a64d15fbad44e747a2e665fd43b404461e7900b6b0e18e027fbd60533ab5`),
		},
		// □ 알고리즘: LEA-192
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`d3e9a79d75aa6bf1833930e9d2d1b19e681756600bb0c29f`),
			Label:       internal.HB(`ba8aa7e94c452d65ed5296cd442c26ff6c8a844898a9d157b04afd50d25461311e5515264adb1eda09a4f20233516425980b889ad74bb4e792f00cd7`),
			Context:     internal.HB(`66d7a8f82cc3ffd088fab181892e0d2f5b454697c1759816649710b5642255a2a055e2da78d0e29dc342c1bc6f3e658dbf6983ac45e4cd5c896fcdf4`),
			L:           0x200,
			K0:          internal.HB(`626f5da75078e53e48c08ea0f8b33f1dd5fe47e42895eaff5b0660ed56ea01a97fa74142d020f114ac3d6474a153ce39bd86644a2bcec4cd0cb3225bfa04701a`),
		},
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`86ceb2a013067ea63605376ea2febac1fe106222bb39c98b`),
			Label:       internal.HB(`e8d392c7f11d0619f91f7dbdbbbc8a9d617983b17580b2cf2d63edbb09a82f8b9186e00316048ae797b61bff3d74e32c4f72a561ffa6057f5901af54`),
			Context:     internal.HB(`c91dd80a2b5d894f79295fd2b12c2d3c47bb2375718c74dc222126fddc5a8250672784c917eaa5073bcf19f24df7b48dfcac3c318d7b59af5d6e7809`),
			L:           0x200,
			K0:          internal.HB(`71d91f8e1708311c247346815558a89d594c12bf762b721ea1fe2ea872b47b76515ce4e831727d756896ae046b20a98e03714edd4cda5e331d0dacb058ca9d5d`),
		},
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 0,
			KI:          internal.HB(`529e385ceac88c3c7705fb41025e5e5c41da88806fa2a917`),
			Label:       internal.HB(`2584586c82027a5d3aa10418572c5214961029403b65de0d6b8989885fe9b16e5965f06c311d74ea8ddefa4424fabf700d16e46aa3e751497d8549cd`),
			Context:     internal.HB(`aaf62b11c8dd36c970a343775c4f12e518463acc5032800ca01b8382b28ac009768831daaba4484656f1fadc7850772555dd62ce0bf0bae23eb421ca`),
			L:           0x200,
			K0:          internal.HB(`e74ea7ebd4bea544d61d4d306f73891c48bc0842a2bc0105edc3f6fb9a3d78f7678bfc4434b0d425d8567a5ddc32e5524c649eb2f9a3ef538b54700ddcd87b53`),
		},
		// □ 알고리즘: LEA-256
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`ca52e219607027656ad342710a9c50e1490a7f1080c24f20362410635516d57d`),
			Label:       internal.HB(`fb68a0e609249f5d8af0fe5c21fbfdfb7d985aad386994ee165a11b438b32cb00af4b5a6657cbc401b6b9ebb4342e5933aecd9bf15fad21a00ea6de4`),
			Context:     internal.HB(`e02048b580b9eb4c4811402f9a3286324258dd73c007a9e3e3740d59202c1ac10129fb71aaf535b30008db167bb5f999b01dd439b83a4368fd42b43a`),
			L:           0x200,
			K0:          internal.HB(`ab78a4d865a2f9ea55012ff04d66e326c738c7a05b489999221f246864433c878b32f0e6f5c14e9f5eff0780777ab1c1f07da5a5f804fb5ac6cca818d1750726`),
		},
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`f243d19a7c9d985a11f737356bb7ccb619a4d4fb8a9f6554300725e9b50ba20a`),
			Label:       internal.HB(`6c928e8b322ea3b47acb7e9e7b3e3bcbecd1c6019d429733aa6e9ee62a587a9a7bcea6a6cfbd90cacf72cca08625bd4dcde4e644a7822676f33863cb`),
			Context:     internal.HB(`4bea13705292dcfd75059fd43cf804f3d5f9012ed12a10ca1d9dcf6b0deb97b790ca8893d8f03060c36278d4f5a7ff024c78dbe4a58650e73df12fd7`),
			L:           0x200,
			K0:          internal.HB(`b2f63313a48a90ff1511cc97dd1a3632ba0e2952da75a76ca52dfd2ddf18bfd853b864ea1a0528cb965235c7900ce5589717e005591642630a58bc4d0792d08a`),
		},
		{
			NewCipher:   lea.NewCipher,
			CounterSize: 0,
			KI:          internal.HB(`eea3023b2422714391afca675c84bd194e5b46cc810c374cd9219931200c9002`),
			Label:       internal.HB(`b95fc6a3b9e475f4479f0ddff22bed622764f04f335dcfc84714aa61c2e9acc0284b90399e25284f897d1735ff6edc9232c07843c564b2dbbcf78be0`),
			Context:     internal.HB(`faf488a972a26257c6a241d43f150eb91af154f7a39f0066c38dcf89c2a4dd9b097e49d55eda03846985a78f5b32c670a88c93278c083fa97da86e80`),
			L:           0x200,
			K0:          internal.HB(`4b44dd37c23a5a265b9a86a62c7461cbd83e20962461f3cfb87f57c6d21bc743d0c34738ec8d0a84197e27c1e1bc7d7ec4e8269b4dcc9a6e326c492a6b35f60e`),
		},
		// □ 알고리즘: SEED
		{
			NewCipher:   seed.NewCipher,
			CounterSize: 8,
			KI:          internal.HB(`5b898cdfdd9b46a8e3f8e5edb79247c4`),
			Label:       internal.HB(`50db604f57ac1864fb19657a60cef868c8e5c04025b08d78222cdeed819093a6bd5117b30d98b3bc69989be2f84e0043761b2686c2060db81c4cc1a3`),
			Context:     internal.HB(`c499e583a24617d0b741677665b75ab65bf95ea9971802d7e73d476b8267bce40697b89efa7d9efbf1a3fcd404f2e58f539df3326ae69adedf17f8fb`),
			L:           0x200,
			K0:          internal.HB(`13604558349b96745683036e8232ec5c78b923854d3b3d0e144befe773a406d1183b62cdaa05f7e560769fd286137a80b23faad6309f90bfcbd83f05c4f788a3`),
		},
		{
			NewCipher:   seed.NewCipher,
			CounterSize: 32,
			KI:          internal.HB(`d4b33a390b92de3bca3c806c19c44b36`),
			Label:       internal.HB(`aec9e38d96c1e2586f5054860cc0f6f2e98b76ff9fe481f5595eb76cecd3c58a9560469a31b10dd97f57bf178350b7f8ba60d5e3483056c0cec5a8a3`),
			Context:     internal.HB(`3afb990af85f083ff6af9409066101a99bc78754245cf47918d4d2f2469a6cbf2a6700e4146a4e2af1bde78706a73287a547c4982d2d0ffc39b5078d`),
			L:           0x200,
			K0:          internal.HB(`d70b91dd72af55bd22e7fc5029b9d6e8b68649b47c429bd1a49d893d41dea67f36271eab939aa6b5b2cf265a01d50f26bfb5ee26d17ce4c6470a055bc2e0b85c`),
		},
		{
			NewCipher:   seed.NewCipher,
			CounterSize: 0,
			KI:          internal.HB(`f666f6c71e65669e464060abc58f86cc`),
			Label:       internal.HB(`c3c4b1a34a5c3263cea91e3ae6e609f9e66d04a21f6fdda23ce8de919ba82008939f12c5b7d50ae29d4e7158226d630a8d43914d693d7a4504483001`),
			Context:     internal.HB(`3755274de4602482fc9775b1cbd693005a5d268e55756828ff3836833de0a056d2d9873bb5f7785f1f6d80e9d4fd409d0052b37ec235fde262d1d4df`),
			L:           0x200,
			K0:          internal.HB(`821b8985dab6022e8cb4b881ec82311eb333d50fe97c8ffe07f2fb48ba50ad335c7ba7315d223572205f1ec21183576307d93d71593268658d3b816a0692dfb7`),
		},
	}
)
