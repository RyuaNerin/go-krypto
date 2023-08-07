package elliptic2m

import (
	"math/big"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
)

func Test_B233_Add(t *testing.T) {
	B233()
	c := &b233

	x, y := new(big.Int), new(big.Int)
	for _, tc := range testCase_B233_Add {
		add(x, y, tc.x1, tc.y1, tc.x2, tc.y2, c)
		if x.Cmp(tc.x) != 0 || y.Cmp(tc.y) != 0 {
			t.Fail()
			return
		}
	}
}

var testCase_B233_Add = []internalTestcase{
	{
		x1: internal.HI(`0106d6eb5aa30b22130fc176da5cbf31e21da5a5a903ff4bc340fbbadeb8`),
		y1: internal.HI(`b9b7253270795df2a2f415e7f4e6511f3e8668da8f95b8361f6f605f65`),
		x2: internal.HI(`0188095ee968b90a30d4cc8a89d24830da2cab24eb2c8ed40025f084c1d2`),
		y2: internal.HI(`9045c9f33b46f621d74babf883a84dc0e119a5e3b9dcfaae63c23740b8`),
		x:  internal.HI(`01fe29c5a0bb64f841a465fc83f4c36c92716eb7983984154160f39c54e1`),
		y:  internal.HI(`60c8cb247fcfb427c5c8fb24e1f3a7c764abbc2305e875ec2b03a7c8b4`),
	},
	{
		x1: internal.HI(`01b693d2db14af541516b87fe537ce4f39324fdd0113debc7b66fba77d06`),
		y1: internal.HI(`017bef5306de3ffa55e16d2311b4f5946252297ea89b22de53002096f190`),
		x2: internal.HI(`01a5ea9b569361834711defe0691b1efd7dcc3e88f1d48f1cf527176bb04`),
		y2: internal.HI(`015ecbf01e16824072eff63e3328423e1cf2ebc8f3d48c247271dfc82b09`),
		x:  internal.HI(`01b0d0181c1fc966f014ed4c8e8e9da1bfaf0ead25c55b77e3121489ff9e`),
		y:  internal.HI(`3f1f68a010cf7e56f9d5b7dbb4d9e0c854217f1e9c0737cf8e1489f808`),
	},
	{
		x1: internal.HI(`2586e9e466686bfe330c7494d6dd72ee7196c3b5dd380d8884a83a2fd3`),
		y1: internal.HI(`014cfce6b7713172b411636d6f0e32c60a0b33ca1ef32848306292cf7580`),
		x2: internal.HI(`5062d7012db19aba9b60bcbb2da5fac434dac4d33dcecd722c55eb86e1`),
		y2: internal.HI(`783f51b69c6df1c6a8c3e7d94faed46cdd703074e6330028571007c279`),
		x:  internal.HI(`015a4039b603e134d5f0f8c984372a988eb6b413f083b6bfa8f067d4fc5a`),
		y:  internal.HI(`7541cd3571e4831fde23559ce1c9492fcede29cf3c19d6b29f99343aac`),
	},
	{
		x1: internal.HI(`12cf0ea2719a91590304d8eed1b8493cea6229fabb590e8221769e9827`),
		y1: internal.HI(`016869b1fb11668f21ca53c86cd113f2cff8f9f2d0968bfc7e984c825229`),
		x2: internal.HI(`68219cce451271d53ec9af6f6cf418d0986fc814a5e73fdfec8913c930`),
		y2: internal.HI(`018334579ee32e77aa80011d3fcd2c60b770ead0398b94800fc9c218b199`),
		x:  internal.HI(`01d2ce9347f94bf3baf7d679e6fac619da8a4e2fcbdd8b33db973ac39278`),
		y:  internal.HI(`01d55bf305d73c4f97aabb6b52545800799db7c93cc1b5badee15080b8e1`),
	},
	{
		x1: internal.HI(`01ecf281d63d429ba6ea80f0f64a434940c5062f696d7eb27da7eee5697e`),
		y1: internal.HI(`0b158d7ba3ce2e53c1983ce70df99301693bce48f1444a5ba8d542c1a8`),
		x2: internal.HI(`a30d16e347ccd534134aa97df6b024844eb7b95349da7a2eecfd857405`),
		y2: internal.HI(`013fda62f385b79e508670a0914da0829f279944f4c1afeb312f3e477d3e`),
		x:  internal.HI(`01846eb86f2e099bbb8e2c98fb07ed02fb54e92630f59940d8fa91421537`),
		y:  internal.HI(`248e1506429ca69d3e8785218c16caf477ca23da62653a1a8770429e6c`),
	},
	{
		x1: internal.HI(`012ed6bf31c13bb522d4e3e15616b474d4ab214ddc7e35dc25c867737fde`),
		y1: internal.HI(`c61a06cd977c5bd04893538c0315ed5b305d21e64805dcc2ff7d027abe`),
		x2: internal.HI(`01624b0c8cb8e219bd3214d0ada08a1e29323cf401141e760d282207cef7`),
		y2: internal.HI(`c11ff1378ffc40eb434aaa5dda1dc2e62a2889f752b1a05a078328f210`),
		x:  internal.HI(`04672a9b9725e03a0e947f570aff7938a5c16430b24fb1119be38bd4fe`),
		y:  internal.HI(`01a9c127aa23e7cd388b1646fda0a392d4e34b1829662f8190d817844bb0`),
	},
	{
		x1: internal.HI(`01fb4ec11f8d0fe522bea7af35819ce1181a99aa8c4793303574e5e0169e`),
		y1: internal.HI(`ae5c79c2c3f6c4180daa370341e73d476428aa14ed52fbe9fc92af5bb3`),
		x2: internal.HI(`125dfdb7c5ee088ac5d9108ece15557c7bea6903d734f2b2e0fa542b4e`),
		y2: internal.HI(`ff22492f2ef876feeb67f506937d3a17bea61806ffb8e1fde774cdb697`),
		x:  internal.HI(`5d73703956b0fb1240346834dd49777eb55b60d2219a5ca52faf9d3f41`),
		y:  internal.HI(`0108ced26166b177f5c2c8e696db5e5af8429bc0152f183cdeda0f24d420`),
	},
	{
		x1: internal.HI(`0bb1105b58acc389814633cda990b54e05a9bb333ed5321fc05e97a5dd`),
		y1: internal.HI(`010af1c667e1326a7a1790e04fa6f88ee5eb34081fabd5452c5b9b9d49a9`),
		x2: internal.HI(`7371c8a9eb6ca9731700fd887111013cb3192eddc5cb18efe0cc427ebf`),
		y2: internal.HI(`d715d386985cd001fcd293d375c0e61169f2719de6e71d54069f385501`),
		x:  internal.HI(`25a09748e172d3faf75d2746ba11336037fe92992b9c78f0b2f165c6a9`),
		y:  internal.HI(`630760a00e0b1ec563187197974ae1df21b17970f891dbeed1c1552d0a`),
	},
	{
		x1: internal.HI(`34d1a270f5df09b90aeb9362c0f15470bd1b7f788e983d25bbc787c2fa`),
		y1: internal.HI(`01a21158e5df9e5af130b7376a894970628417b860f2db3514d8ab8295ef`),
		x2: internal.HI(`01974f7b608fe62dc7d279e061a0d5ac9118962042cbffc9112904b30bd5`),
		y2: internal.HI(`466154c6c5d1ac90f82f796c0fd11f4bc94475db5b3229ac96b5d010b5`),
		x:  internal.HI(`018483a2cce0b334cb93133a0f7bd0f470614c7daee2debe03c55c320cc3`),
		y:  internal.HI(`0123361427f9d8894f9e0f283368f33fa028dc39c10a8e43b77f0b37242d`),
	},
	{
		x1: internal.HI(`01a3aef4a879db97971fcc24f0eb19aef382856ecc4098ed78c9131c3114`),
		y1: internal.HI(`0126b1eb3831b9be35387728a9472ddb8003bcb2dc931b652bd7ce76184e`),
		x2: internal.HI(`3a88da03bc80c64ec520bd81b45b3cc52529d7e9778291868fc1e77471`),
		y2: internal.HI(`a2ac89e54dc3561b4fbe6fd270a302e6517bb658e5bac18025dcb4f267`),
		x:  internal.HI(`e0306149ff569c51c7fca1e2e002943450cf38b34153cd331613524a5e`),
		y:  internal.HI(`014d7906a58975db44ae45ebe5dddcc2d2343c4513f12d717ec93ed4e6cd`),
	},
	{
		x1: internal.HI(`635c023e8ebf4aaa7ab614d8a7b17400310731a853f521ce22b444c86e`),
		y1: internal.HI(`01aa426e9dc389b691d00f6120bce40ba7bf984f8369084e7d086e9396df`),
		x2: internal.HI(`d0c9347de6c1fb29aa2bf9b6af7372c76b07103ad928e311b011c69265`),
		y2: internal.HI(`6980211826cd8b39ebccd94e0132e535b85ee3657576914da0bcb0271e`),
		x:  internal.HI(`24ddee8c1257b2c15979eac40a0d518ade08ef75fa6a90b18ddf51277f`),
		y:  internal.HI(`01079df35f51659aebe187ffe9e3c9d78621725574f672eccd856c9f2134`),
	},
	{
		x1: internal.HI(`017024878566bd46b998df594add6f326bc8e9ed4d36ac4747f44fb8139a`),
		y1: internal.HI(`019a62cc11f3a91c0edeee34acd2f201202c8a875fbc8287cb2865bf8413`),
		x2: internal.HI(`16af0b53fd412709be9c903fd2f088b711c65c4783409f4f64c0a4c838`),
		y2: internal.HI(`01b1a5480cf99f2a7fcd90bef76e816cd733ab9658e3b259d964ff9dc08e`),
		x:  internal.HI(`20daa32babf6574c98ba758baba61a025f20182c0e18dd90792f01dc5a`),
		y:  internal.HI(`018b55a9786079c95aa4dd04c2f3b4ad4448d92cdab23f2da5d7c9770f78`),
	},
	{
		x1: internal.HI(`013dd1ed2936441434a50bba84c7507d327b9f14577e707ce7898d2afff5`),
		y1: internal.HI(`01bf572552dca3e592ea043cd1350c74b147cf05f1e9b3ceaa5d22932ddd`),
		x2: internal.HI(`0646cfb97315c48be392b63cb13958e4ed89e195252164a9a68f907384`),
		y2: internal.HI(`01b34eea654dbc04d2b6b69b2d83f77eee7c53e202e530fbb20942c65880`),
		x:  internal.HI(`018a98d43cca20547c0bcc9915069b4f5880c8d09f4c61930ca3253d4f46`),
		y:  internal.HI(`01e980c14f395ce608fa9dae8361e3ddbdc53fb89887a1070fe2a0f3fab9`),
	},
	{
		x1: internal.HI(`1201d4ca67a06418168f6ae800d50c3f49f22b5ea084e9aea814b29064`),
		y1: internal.HI(`01a0968010a22b623cbf7b461b77b7c00991da4efb36b138a0d7b26ba9a7`),
		x2: internal.HI(`019f787ff6b013490877fa501dc8cedd0d6be106d0a793a05379ca6dce4a`),
		y2: internal.HI(`c9f8ac11da58daae65da983f2baed05569e027276effa7d761077fbfb0`),
		x:  internal.HI(`bace3dc159fd3f5e2f51f53f00d4967e9786bec97908d2c60d074b1b7b`),
		y:  internal.HI(`0185144316f9007fd06fb8fd1a46f12ba49bd7e6905763081fa7ea207bc0`),
	},
	{
		x1: internal.HI(`01781b25ec1ca8df92bbab95f7433822b7a0e56ee53ef7602a145664a44e`),
		y1: internal.HI(`e0204e1cbdad24e430ec2b865df001add24d9dcc061b0ba229b793f92b`),
		x2: internal.HI(`0116ae3306cc298250287191d90e314c8c9a68462d7f3c9cb55bbf177321`),
		y2: internal.HI(`d894616acb321bf058db05a581b92bc8c26e191b3aa8653870201a348d`),
		x:  internal.HI(`facb133d85a0c0c416296995ba328ae07c019ca5f118b07d4b75c17a9d`),
		y:  internal.HI(`9391f840a41f1b8c379cd6f4bb765b771d14ed3483d12a6cdc2c4f4c7a`),
	},
	{
		x1: internal.HI(`e1368c09844fed8d6edd68e1d2be2af4a0d29aad42b7ce9db343805670`),
		y1: internal.HI(`017e1d9cfe0787d6a28fa7877b34a09e74665a61216255250cb81cf3ad70`),
		x2: internal.HI(`011176cd9614fbaba25ec98a924b5b101f416334c275c268187651d61cfc`),
		y2: internal.HI(`01fc0bc49cd798bf1221549c8b0b4eb94bb94013e47c773fd3ac7893d0a1`),
		x:  internal.HI(`eabccea6a24d7a45dc6f813341a899e52d7431ff972a783f101637f0f5`),
		y:  internal.HI(`010c8f2ffdf854ddeaa8b15bb25556f0bcf699f825082fd7a4f5e5b85647`),
	},
	{
		x1: internal.HI(`0156a31228c1c6430c89098dbef18f8284c42de47a4ec959ed50036a946f`),
		y1: internal.HI(`01bb651d894b7fad8496c29edfeaebecc23debde0fd4d9ee406b475ba737`),
		x2: internal.HI(`013f6a4c7e6d5447a829c47ae07d7afb5b2a73f257c86b6ff6d443332f20`),
		y2: internal.HI(`01eb5cdc7c93bc071d1d93cc18767987554a6946fef4e2aacad366ac6e84`),
		x:  internal.HI(`ee2094eefad3d533109791d9bac2f23ee143dcea2c5b96a2ff5b8b33c8`),
		y:  internal.HI(`01d5e0db503fdd74372460388dae8111f960e9fd393e7026fd58d540009a`),
	},
	{
		x1: internal.HI(`ed6f9fb6c2070007f3e52efb1df2422fbdee10b88b08bcbf4914d9fdf0`),
		y1: internal.HI(`0168bfca8404f3e31c48a7c35751e3b5eaff032c22cf249f57ac5d33f9a6`),
		x2: internal.HI(`01b2adb4bf363041959284074474480669bfbb80d7b7e1aa5ff930fff2bc`),
		y2: internal.HI(`219c6e2c8555308a443be97e2aaf5184c326ffcbd5f21ee8e88709026a`),
		x:  internal.HI(`01716394a9047b7b96c1d4e9c6217418607e2c57bc890b064849371ba2e2`),
		y:  internal.HI(`636adb8f47a98bf64f9d13d000fe66fb8086920b7766f2d3412a556f26`),
	},
	{
		x1: internal.HI(`6f2895fe58961fd9a69d12d395e57989245865f7923c80d54099754452`),
		y1: internal.HI(`011ab850cbfbf9ac1a871ee043023c95eb845513301676ab9c9604d1da9a`),
		x2: internal.HI(`18ab58d1a82835c451f6459e1c2aeb09300831ad2ae17c9e8a1d9f8682`),
		y2: internal.HI(`f497154b3f31c619915d92ee3f0e2556414f79fb4263781a7789cdd61d`),
		x:  internal.HI(`08fc93d9cc77a17bb4ee2ea6491eda72b971050b36039bc2ad0dc6fc60`),
		y:  internal.HI(`5875efc4d4dd7c2476614841509abf17adc7fdf76a140f2c339f3ed787`),
	},
	{
		x1: internal.HI(`2c83ef2b17bad80fe3a63d3a017b422a53af72f2fa257660ef9d83c06c`),
		y1: internal.HI(`59558fa1434199f7e259bf477cf36d5c4d1f443a2c1aed84a7a22ac093`),
		x2: internal.HI(`0176331e6e511615d5b894370d95bb276bc2f7febaaa44e45b35158ef624`),
		y2: internal.HI(`01338eb9d10eab78136895923f101bce1fe261beb7dbbcc1d748062e121d`),
		x:  internal.HI(`01d0cbd8c6c6e10d443c137d9f9f19697889abc01b2a560b16caf82f0193`),
		y:  internal.HI(`0114eb9df62e5bdee12097993835bc6f94dc384f61dcdd0114069008e20f`),
	},
	{
		x1: internal.HI(`db81c2c1a9c0f010e97e144749db8aa4d9e731a5aac761f929c2910819`),
		y1: internal.HI(`018691e30217a6688bad9b8acd6a9275bcea400675b56e803b878c420b04`),
		x2: internal.HI(`072f6078484c926a00fd0488a0ca32b6d23ca166a75f5405a3b9e23222`),
		y2: internal.HI(`013cbe4615f67b3066a6671cdbe9f57698c6187762aef7158b56329ee8fa`),
		x:  internal.HI(`e66fcf3c4d0ba4d1461c86eac67d516f0da0a9bc4ba2050478c2e2196b`),
		y:  internal.HI(`01b06dbe0582e15ec7c0e7ac1ed691952832b802fb9130e9116e2cf6b391`),
	},
	{
		x1: internal.HI(`93029e98f3233aa98f76c35f77c8b0f17edbc1ff9a5c49dd8250ec9799`),
		y1: internal.HI(`01a7adfdc465a4f398fde1b7bf2f85046718fdea50934f60591328587efc`),
		x2: internal.HI(`025d664b2a780a7fc237ce542fd6e5faeb48150e72e5534c8dee0f790f`),
		y2: internal.HI(`9b2d605625b73ff01043744fc03605c3e3231fe28f8fce2721d7365585`),
		x:  internal.HI(`924036a5c6eb43d303ac196665e898af9447289f654f4086ec6663e5f0`),
		y:  internal.HI(`3a539007fdb9e3480f71877702a9971a02a3fc1bb314afabd6e85903dd`),
	},
	{
		x1: internal.HI(`019458420ee1ec28cfa13e92f05734e6f252fcfea23556f9897c4775e409`),
		y1: internal.HI(`01036b4a97212fee15dfd845d4b7d5925628307ad57d66705ce1b8c536e8`),
		x2: internal.HI(`e607dfbb51617f59107affdd45407b570e53445019b3c99b0d2e2ea5cf`),
		y2: internal.HI(`88690f1287062b9b18045a47f43259c6c2ae8b0955bd416c81092987c9`),
		x:  internal.HI(`272a2d171c792c86a26027e2139246c86ceeb5d1bf036e6c1b97fe6aa1`),
		y:  internal.HI(`01ead168abddfa337410638cb76cea3fbae892716563cba4dbbb881e6807`),
	},
	{
		x1: internal.HI(`01e8ba15c0ab98c46343d83ce05291db7685c0dd8c14def628132bb59f14`),
		y1: internal.HI(`9b0eb4a59256ec757160a00d774dcd75694bdbfbf7292ca1a02b322b79`),
		x2: internal.HI(`01df5cfe5e2aeece49286a77fa4ab7875b87b63470150f1150892d12000f`),
		y2: internal.HI(`013475801a52553f7048d701e990da7a805ccefeb225c74a4671c2f042e2`),
		x:  internal.HI(`a7e9222bf9874fe6f4fb88565ad2291c36d9f37ba76be08832a80b95ae`),
		y:  internal.HI(`512aef11bd6ec573f2cf09dae71e2f3363d2a2501967979a76a63579ad`),
	},
	{
		x1: internal.HI(`ef5263c55ac708e63c7c45132e6068db114370d4928395725849bb1659`),
		y1: internal.HI(`1b9035c8c408513eeeeeadb0c526da204a5cf1f7150061af2a13b61d44`),
		x2: internal.HI(`01bf885c4f2c8488a8f89c5fb8d373829283107b5b44bdc47bbf3817cfa8`),
		y2: internal.HI(`4251f49a15503b8355767180f39af8d498398ce53593e19cabd1aee5b7`),
		x:  internal.HI(`01541084766c0244128048843be57da3232bf36d490bd56f2b4c002b7161`),
		y:  internal.HI(`ada6e0fe54fb7d3f8b4b51e42b68202a9899c0c02bb6fbff6b2fbe48d2`),
	},
	{
		x1: internal.HI(`94ff4ffe023f57e1877aaf81fa4df41c1365031aefc9bc70acf28484f3`),
		y1: internal.HI(`01fd3953206661bea6207decf2e26f36d5ad68d8e070782d4d3c8ae74b6e`),
		x2: internal.HI(`98cf29d7b1b62bce249a79a514e397fe64261612aa9d56d183ebb73f92`),
		y2: internal.HI(`f215da9ed3b969631d1bf570e86006a11d42582054bfc5ab5374f77e4e`),
		x:  internal.HI(`0e0476bc6b3298fe0b4180759ad69a7731ddb7e27b03bf227f50877f0f`),
		y:  internal.HI(`01b477cf699fe99d61ec1d597f55b0d618dc4f8dc7ea9761d01a47379d49`),
	},
	{
		x1: internal.HI(`01d7470ae240cfe5c547542c6a76e8bafe6a55cdeb56900d4e7f0c417c53`),
		y1: internal.HI(`ad0527563928d0345df69e1cd1faef5c7b7a0d007964722741d8c4e4ad`),
		x2: internal.HI(`158ea316877964726747473b3b8e3d8dc5c7a67ac75a3d5d23bb747947`),
		y2: internal.HI(`9e3776b7c7ed12465aad48e464b1e2c4be38e218ae3e329f1626c77958`),
		x:  internal.HI(`013c25935fb04b90704f14408438d7ec14cfb188e10b69b097717e850187`),
		y:  internal.HI(`01d3e404170d776a60997d20ac9ca33963001b6a3f6a1eebc11a2bd6f106`),
	},
	{
		x1: internal.HI(`01bcc39cb52eb633e7072c5ce317387a90c307bcd6334ca09297adde74c1`),
		y1: internal.HI(`4a86598fa98df52b2ff5e8bb0dac083fe4901cb3727dfa2936b00c630f`),
		x2: internal.HI(`0141ac77aafaedefb6f7ffbbef897572811fda450539a8cf2024b2335ca4`),
		y2: internal.HI(`0191d0d188872a1922dfb684781619616575b11a79a776d2fbfd893c5d69`),
		x:  internal.HI(`908c5ca9417aceb3a3b92de26911b09fd88fda9770273df7706a4481e6`),
		y:  internal.HI(`a231059bc65556fcae1d8016f813b0875ee80f9210b0e3f307bdb30512`),
	},
	{
		x1: internal.HI(`dcb96005f3d0206eb6dbab134a3162f5530bff6cbfc35129d8f87bc566`),
		y1: internal.HI(`66a9b1a3d8cdb3422b45fd5921499bcf25d115a7c77ef871d3de3c256e`),
		x2: internal.HI(`012a012f282ee415fd6be6b8ed60b1baa11a0a820b6905fdfc7a9e09ec58`),
		y2: internal.HI(`018b3c9f580c5efcb12b79e63791539d47f35baeef6346df2433470d70bb`),
		x:  internal.HI(`014ca78e15d37cce20fa7944bb8012d560fdd8acacd760a74551ad13306f`),
		y:  internal.HI(`018d9bb33b23b1c244deee8920e268dd65be35bdb1ccd281971ae9af3fea`),
	},
	{
		x1: internal.HI(`fba0e152469c4ad43ca3b28ccad309d90a4ceb00e58cdead0a3c24b286`),
		y1: internal.HI(`41938d48e3dc413cc2f6ee30ca1e35d5f3a005a844cf8bc915a3cc4095`),
		x2: internal.HI(`01f171e6b4a904c46b9db91071efbebfeed79568b8b3da12b890f6d31cba`),
		y2: internal.HI(`0178926b075deb3d194e09d9d75f249331d3772f1143c1f4dc7bebc0ebdb`),
		x:  internal.HI(`01ab70c5cd7b3a0653e8903c822b51785f95f7ae4e712f448d8400bfb248`),
		y:  internal.HI(`8c5475db31318bcde3b131ad09baa2368c2f2a44050f7c52ef7d97bf9d`),
	},
	{
		x1: internal.HI(`d03b7b691b6927e4753801d76334a1784ad0192a4a29254769fdc4cde1`),
		y1: internal.HI(`ab12835de79553967219dfd7ab2591298bb6deb94d7e2258e28abe07ed`),
		x2: internal.HI(`019b930b18487018dfb56ce833f68421e4c325308380a89ba58e693d5fbd`),
		y2: internal.HI(`08479e937524b709c82d535e472fceaf518f59b3d1bc51337008a98702`),
		x:  internal.HI(`04c1e9094f7fe9a3b97ae5e107e73da28fe79df79064e71fb5518fbd65`),
		y:  internal.HI(`01b0590f50e823f3df3202c204bce27ab93878afd5372e7b3814771b64fc`),
	},
	{
		x1: internal.HI(`019292d1e0c2ddd754754d8fb9f5c1ba7963e6372a17b655e8f9ebbde91b`),
		y1: internal.HI(`fa3a912b5f52dbe08a2a885aaa8f25ed1807c3f765d12463054241e4ad`),
		x2: internal.HI(`b58932a2ef928d6a79b2c7daebb5c153b569ccd9cb3c3fe5a72265ca60`),
		y2: internal.HI(`e7968db76852dc81881e2966480699a29a823935994e4ad9c4959029ef`),
		x:  internal.HI(`459520c2fade99663a1956ce901a5c272db37116800417e04663a31d3f`),
		y:  internal.HI(`01fcf3fb596905e0ad2a6c5753819c4ab1fa556629ad8d4cad543b1ea40b`),
	},
}