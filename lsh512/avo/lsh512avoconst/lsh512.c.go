package lsh512avoconst

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

const (
	/* -------------------------------------------------------- *
	 * LSH: parameters
	 * -------------------------------------------------------- */
	MSG_BLK_WORD_LEN      = 32
	CV_WORD_LEN           = 16
	CONST_WORD_LEN        = 8
	HASH_VAL_MAX_WORD_LEN = 8

	WORD_BIT_LEN = 64

	/* -------------------------------------------------------- */

	NUM_STEPS = 28

	ROT_EVEN_ALPHA = 23
	ROT_EVEN_BETA  = 59
	ROT_ODD_ALPHA  = 7
	ROT_ODD_BETA   = 3
)

var (
	G_IV224         Mem
	G_IV256         Mem
	G_IV384         Mem
	G_IV512         Mem
	G_StepConstants Mem
)

func init() {
	var g_IV224 = []uint64{
		0x0C401E9FE8813A55, 0x4A5F446268FD3D35, 0xFF13E452334F612A, 0xF8227661037E354A,
		0xA5F223723C9CA29D, 0x95D965A11AED3979, 0x01E23835B9AB02CC, 0x52D49CBAD5B30616,
		0x9E5C2027773F4ED3, 0x66A5C8801925B701, 0x22BBC85B4C6779D9, 0xC13171A42C559C23,
		0x31E2B67D25BE3813, 0xD522C4DEED8E4D83, 0xA79F5509B43FBAFE, 0xE00D2CD88B4B6C6A,
	}

	var g_IV256 = []uint64{
		0x6DC57C33DF989423, 0xD8EA7F6E8342C199, 0x76DF8356F8603AC4, 0x40F1B44DE838223A,
		0x39FFE7CFC31484CD, 0x39C4326CC5281548, 0x8A2FF85A346045D8, 0xFF202AA46DBDD61E,
		0xCF785B3CD5FCDB8B, 0x1F0323B64A8150BF, 0xFF75D972F29EA355, 0x2E567F30BF1CA9E1,
		0xB596875BF8FF6DBA, 0xFCCA39B089EF4615, 0xECFF4017D020B4B6, 0x7E77384C772ED802,
	}

	var g_IV384 = []uint64{
		0x53156A66292808F6, 0xB2C4F362B204C2BC, 0xB84B7213BFA05C4E, 0x976CEB7C1B299F73,
		0xDF0CC63C0570AE97, 0xDA4441BAA486CE3F, 0x6559F5D9B5F2ACC2, 0x22DACF19B4B52A16,
		0xBBCDACEFDE80953A, 0xC9891A2879725B3E, 0x7C9FE6330237E440, 0xA30BA550553F7431,
		0xBB08043FB34E3E30, 0xA0DEC48D54618EAD, 0x150317267464BC57, 0x32D1501FDE63DC93,
	}

	var g_IV512 = []uint64{
		0xadd50f3c7f07094e, 0xe3f3cee8f9418a4f, 0xb527ecde5b3d0ae9, 0x2ef6dec68076f501,
		0x8cb994cae5aca216, 0xfbb9eae4bba48cc7, 0x650a526174725fea, 0x1f9a61a73f8d8085,
		0xb6607378173b539b, 0x1bc99853b0c0b9ed, 0xdf727fc19b182d47, 0xdbef360cf893a457,
		0x4981f5e570147e80, 0xd00c4490ca7d3e30, 0x5d73940c0e4ae1ec, 0x894085e2edb2d819,
	}

	var g_StepConstants = []uint64{
		0x97884283c938982a, 0xba1fca93533e2355, 0xc519a2e87aeb1c03, 0x9a0fc95462af17b1,
		0xfc3dda8ab019a82b, 0x02825d079a895407, 0x79f2d0a7ee06a6f7, 0xd76d15eed9fdf5fe,
		0x1fcac64d01d0c2c1, 0xd9ea5de69161790f, 0xdebc8b6366071fc8, 0xa9d91db711c6c94b,
		0x3a18653ac9c1d427, 0x84df64a223dd5b09, 0x6cc37895f4ad9e70, 0x448304c8d7f3f4d5,
		0xea91134ed29383e0, 0xc4484477f2da88e8, 0x9b47eec96d26e8a6, 0x82f6d4c8d89014f4,
		0x527da0048b95fb61, 0x644406c60138648d, 0x303c0e8aa24c0edc, 0xc787cda0cbe8ca19,
		0x7ba46221661764ca, 0x0c8cbc6acd6371ac, 0xe336b836940f8f41, 0x79cb9da168a50976,
		0xd01da49021915cb3, 0xa84accc7399cf1f1, 0x6c4a992cee5aeb0c, 0x4f556e6cb4b2e3e0,
		0x200683877d7c2f45, 0x9949273830d51db8, 0x19eeeecaa39ed124, 0x45693f0a0dae7fef,
		0xedc234b1b2ee1083, 0xf3179400d68ee399, 0xb6e3c61b4945f778, 0xa4c3db216796c42f,
		0x268a0b04f9ab7465, 0xe2705f6905f2d651, 0x08ddb96e426ff53d, 0xaea84917bc2e6f34,
		0xaff6e664a0fe9470, 0x0aab94d765727d8c, 0x9aa9e1648f3d702e, 0x689efc88fe5af3d3,
		0xb0950ffea51fd98b, 0x52cfc86ef8c92833, 0xe69727b0b2653245, 0x56f160d3ea9da3e2,
		0xa6dd4b059f93051f, 0xb6406c3cd7f00996, 0x448b45f3ccad9ec8, 0x079b8587594ec73b,
		0x45a50ea3c4f9653b, 0x22983767c1f15b85, 0x7dbed8631797782b, 0x485234be88418638,
		0x842850a5329824c5, 0xf6aca914c7f9a04c, 0xcfd139c07a4c670c, 0xa3210ce0a8160242,
		0xeab3b268be5ea080, 0xbacf9f29b34ce0a7, 0x3c973b7aaf0fa3a8, 0x9a86f346c9c7be80,
		0xac78f5d7cabcea49, 0xa355bddcc199ed42, 0xa10afa3ac6b373db, 0xc42ded88be1844e5,
		0x9e661b271cff216a, 0x8a6ec8dd002d8861, 0xd3d2b629beb34be4, 0x217a3a1091863f1a,
		0x256ecda287a733f5, 0xf9139a9e5b872fe5, 0xac0535017a274f7c, 0xf21b7646d65d2aa9,
		0x048142441c208c08, 0xf937a5dd2db5e9eb, 0xa688dfe871ff30b7, 0x9bb44aa217c5593b,
		0x943c702a2edb291a, 0x0cae38f9e2b715de, 0xb13a367ba176cc28, 0x0d91bd1d3387d49b,
		0x85c386603cac940c, 0x30dd830ae39fd5e4, 0x2f68c85a712fe85d, 0x4ffeecb9dd1e94d6,
		0xd0ac9a590a0443ae, 0xbae732dc99ccf3ea, 0xeb70b21d1842f4d9, 0x9f4eda50bb5c6fa8,
		0x4949e69ce940a091, 0x0e608dee8375ba14, 0x983122cba118458c, 0x4eeba696fbb36b25,
		0x7d46f3630e47f27e, 0xa21a0f7666c0dea4, 0x5c22cf355b37cec4, 0xee292b0c17cc1847,
		0x9330838629e131da, 0x6eee7c71f92fce22, 0xc953ee6cb95dd224, 0x3a923d92af1e9073,
		0xc43a5671563a70fb, 0xbc2985dd279f8346, 0x7ef2049093069320, 0x17543723e3e46035,
		0xc3b409b00b130c6d, 0x5d6aee6b28fdf090, 0x1d425b26172ff6ed, 0xcccfd041cdaf03ad,
		0xfe90c7c790ab6cbf, 0xe5af6304c722ca02, 0x70f695239999b39e, 0x6b8b5b07c844954c,
		0x77bdb9bb1e1f7a30, 0xc859599426ee80ed, 0x5f9d813d4726e40a, 0x9ca0120f7cb2b179,
		0x8f588f583c182cbd, 0x951267cbe9eccce7, 0x678bb8bd334d520e, 0xf6e662d00cd9e1b7,
		0x357774d93d99aaa7, 0x21b2edbb156f6eb5, 0xfd1ebe846e0aee69, 0x3cb2218c2f642b15,
		0xe7e7e7945444ea4c, 0xa77a33b5d6b9b47c, 0xf34475f0809f6075, 0xdd4932dce6bb99ad,
		0xacec4e16d74451dc, 0xd4a0a8d084de23d6, 0x1bdd42f278f95866, 0xeed3adbb938f4051,
		0xcfcf7be8992f3733, 0x21ade98c906e3123, 0x37ba66711fffd668, 0x267c0fc3a255478a,
		0x993a64ee1b962e88, 0x754979556301faaa, 0xf920356b7251be81, 0xc281694f22cf923f,
		0x9f4b6481c8666b02, 0xcf97761cfe9f5444, 0xf220d7911fd63e9f, 0xa28bd365f79cd1b0,
		0xd39f5309b1c4b721, 0xbec2ceb864fca51f, 0x1955a0ddc410407a, 0x43eab871f261d201,
		0xeaafe64a2ed16da1, 0x670d931b9df39913, 0x12f868b0f614de91, 0x2e5f395d946e8252,
		0x72f25cbb767bd8f4, 0x8191871d61a1c4dd, 0x6ef67ea1d450ba93, 0x2ea32a645433d344,
		0x9a963079003f0f8b, 0x74a0aeb9918cac7a, 0x0b6119a70af36fa3, 0x8d9896f202f0d480,
		0x654f1831f254cd66, 0x1318a47f0366a25e, 0x65752076250b4e01, 0xd1cd8eb888071772,
		0x30c6a9793f4e9b25, 0x154f684b1e3926ee, 0x6c7ac0b1fe6312ae, 0x262f88f4f3c5550d,
		0xb4674a24472233cb, 0x2bbd23826a090071, 0xda95969b30594f66, 0x9f5c47408f1e8a43,
		0xf77022b88de9c055, 0x64b7b36957601503, 0xe73b72b06175c11a, 0x55b87de8b91a6233,
		0x1bb16e6b6955ff7f, 0xe8e0a5ec7309719c, 0x702c31cb89a8b640, 0xfba387cfada8cde2,
		0x6792db4677aa164c, 0x1c6b1cc0b7751867, 0x22ae2311d736dc01, 0x0e3666a1d37c9588,
		0xcd1fd9d4bf557e9a, 0xc986925f7c7b0e84, 0x9c5dfd55325ef6b0, 0x9f2b577d5676b0dd,
		0xfa6e21be21c062b3, 0x8787dd782c8d7f83, 0xd0d134e90e12dd23, 0x449d087550121d96,
		0xecf9ae9414d41967, 0x5018f1dbf789934d, 0xfa5b52879155a74c, 0xca82d4d3cd278e7c,
		0x688fdfdfe22316ad, 0x0f6555a4ba0d030a, 0xa2061df720f000f3, 0xe1a57dc5622fb3da,
		0xe6a842a8e8ed8153, 0x690acdd3811ce09d, 0x55adda18e6fcf446, 0x4d57a8a0f4b60b46,
		0xf86fbfc20539c415, 0x74bafa5ec7100d19, 0xa824151810f0f495, 0x8723432791e38ebb,
		0x8eeaeb91d66ed539, 0x73d8a1549dfd7e06, 0x0387f2ffe3f13a9b, 0xa5004995aac15193,
		0x682f81c73efdda0d, 0x2fb55925d71d268d, 0xcc392d2901e58a3d, 0xaa666ab975724a42,
	}

	f := func(name string, arr []uint64) Mem {
		mem := GLOBL(name, NOPTR|RODATA)

		for i, v := range arr {
			DATA(i*8, U64(v))
		}

		return mem
	}

	G_IV224 = f("g_IV224", g_IV224)
	G_IV256 = f("g_IV256", g_IV256)
	G_IV384 = f("g_IV384", g_IV384)
	G_IV512 = f("g_IV512", g_IV512)
	G_StepConstants = f("g_StepConstants", g_StepConstants)
}
