package hight

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
)

var (
	testCases = []struct {
		Key    []byte
		Plain  []byte
		Secure []byte
	}{
		// TTAK.KO-12.0040_R1
		// p. 21
		// Ⅰ.1. 참조구현값 1
		{
			Key:    internal.Reverse(internal.HB(`00 11 22 33 44 55 66 77 88 99 aa bb cc dd ee ff`)),
			Plain:  internal.Reverse(internal.HB(`00 00 00 00 00 00 00 00`)),
			Secure: internal.Reverse(internal.HB(`00 f4 18 ae d9 4f 03 f2`)),
		},
		// p. 22
		// Ⅰ.2. 참조구현값 2
		{
			Key:    internal.Reverse(internal.HB(`ff ee dd cc bb aa 99 88 77 66 55 44 33 22 11 00`)),
			Plain:  internal.Reverse(internal.HB(`00 11 22 33 44 55 66 77`)),
			Secure: internal.Reverse(internal.HB(`23 ce 9f 72 e5 43 e6 d8`)),
		},
		// p. 23
		// Ⅰ.3. 참조구현값 3
		{
			Key:    internal.Reverse(internal.HB(`00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f`)),
			Plain:  internal.Reverse(internal.HB(`01 23 45 67 89 ab cd ef`)),
			Secure: internal.Reverse(internal.HB(`7a 6f b2 a2 8d 23 f4 66`)),
		},
		// p. 24
		// Ⅰ.4. 참조구현값 4
		{
			Key:    internal.Reverse(internal.HB(`28 db c3 bc 49 ff d8 7d cf a5 09 b1 1d 42 2b e7`)),
			Plain:  internal.Reverse(internal.HB(`b4 1e 6b e2 eb a8 4a 14`)),
			Secure: internal.Reverse(internal.HB(`cc 04 7a 75 20 9c 1f c6`)),
		},
		//////////////////////////////////////////////////
		// 암호알고리즘 검증기준 V3.0
		// 테스트 벡터
		// HIGHT(ECB)KAT.txt
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`8000000000000000`),
			Secure: internal.HB(`D2B366EE33648CCE`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`C000000000000000`),
			Secure: internal.HB(`C6FB1015230CC831`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`E000000000000000`),
			Secure: internal.HB(`4996F36EFA10A200`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`F000000000000000`),
			Secure: internal.HB(`5A32EF79EF37F039`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`F800000000000000`),
			Secure: internal.HB(`C27C45DB4826B081`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FC00000000000000`),
			Secure: internal.HB(`407DC4E6E2665373`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FE00000000000000`),
			Secure: internal.HB(`C943005CB54AC89A`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FF00000000000000`),
			Secure: internal.HB(`B240870B46532CCC`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FF80000000000000`),
			Secure: internal.HB(`274EB5A1E6C51720`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFC0000000000000`),
			Secure: internal.HB(`79C27255AB7525BA`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFE0000000000000`),
			Secure: internal.HB(`23626A73A12D9837`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFF0000000000000`),
			Secure: internal.HB(`BA8291096591F1B5`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFF8000000000000`),
			Secure: internal.HB(`F2375F5348192814`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFC000000000000`),
			Secure: internal.HB(`3B7430404A635B43`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFE000000000000`),
			Secure: internal.HB(`3AA1A277D46CE103`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFF000000000000`),
			Secure: internal.HB(`B59B79BEF5780981`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFF800000000000`),
			Secure: internal.HB(`E501F77D7F3A02E2`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFC00000000000`),
			Secure: internal.HB(`EDEE202D0DF9FC2D`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFE00000000000`),
			Secure: internal.HB(`98FA099FA226A026`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFF00000000000`),
			Secure: internal.HB(`DAD7E68A5D93780E`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFF80000000000`),
			Secure: internal.HB(`F4B615C18C65F1FB`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFC0000000000`),
			Secure: internal.HB(`77A915ADB6B9B87B`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFE0000000000`),
			Secure: internal.HB(`D0184F72DDBABE5E`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFF0000000000`),
			Secure: internal.HB(`43C504F4A4206A4D`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFF8000000000`),
			Secure: internal.HB(`66CCEB19928DF1DA`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFC000000000`),
			Secure: internal.HB(`106D1188E8019698`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFE000000000`),
			Secure: internal.HB(`A16B66FB2592AAB7`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFF000000000`),
			Secure: internal.HB(`95FCF3F11621DF7F`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFF800000000`),
			Secure: internal.HB(`AEE21298344AA294`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFC00000000`),
			Secure: internal.HB(`E5D6317B1F6F9CC2`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFE00000000`),
			Secure: internal.HB(`E37441CDC14ACC83`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFF00000000`),
			Secure: internal.HB(`7E56698097CC82CD`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFF80000000`),
			Secure: internal.HB(`EDC413DA94D2E9AF`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFC0000000`),
			Secure: internal.HB(`59FB312B17848FC9`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFE0000000`),
			Secure: internal.HB(`0EFE10D71BBB1028`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFF0000000`),
			Secure: internal.HB(`CD5F1F5EE6E77855`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFF8000000`),
			Secure: internal.HB(`4D53398F007A2189`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFC000000`),
			Secure: internal.HB(`7F842438056A6669`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFE000000`),
			Secure: internal.HB(`D12A323909B676DA`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFF000000`),
			Secure: internal.HB(`FBDC729AEDC6752F`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFF800000`),
			Secure: internal.HB(`9D15343AA1D0E4F6`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFC00000`),
			Secure: internal.HB(`5367E70AF3A05574`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFE00000`),
			Secure: internal.HB(`B6A198540C60481C`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFF00000`),
			Secure: internal.HB(`BE196FA98F01A1E2`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFF80000`),
			Secure: internal.HB(`2E5EAFADEA9A7BBE`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFC0000`),
			Secure: internal.HB(`5DBBF5B44271384C`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFE0000`),
			Secure: internal.HB(`BA6DCE6D958D5FA6`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFF0000`),
			Secure: internal.HB(`7660DC3193941B36`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFF8000`),
			Secure: internal.HB(`9184BC7B9A9F33B3`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFFC000`),
			Secure: internal.HB(`AB660E01B7789D6F`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFFE000`),
			Secure: internal.HB(`B00446094A5300F4`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFFF000`),
			Secure: internal.HB(`5AE010C0B545C7C9`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFFF800`),
			Secure: internal.HB(`DBD573FC610D614F`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFFFC00`),
			Secure: internal.HB(`41A75FAF0651C5E8`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFFFE00`),
			Secure: internal.HB(`1C3DB400B8C1D163`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFFFF00`),
			Secure: internal.HB(`7681D3FEDC95C0F8`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFFFF80`),
			Secure: internal.HB(`2E5F875423A6D155`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFFFFC0`),
			Secure: internal.HB(`E4BB3D7BEB1FC86B`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFFFFE0`),
			Secure: internal.HB(`D5C348E79F3508BC`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFFFFF0`),
			Secure: internal.HB(`A630D1CF46CB999B`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFFFFF8`),
			Secure: internal.HB(`87220E6D3BABECC3`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFFFFFC`),
			Secure: internal.HB(`A2FF263EF5DA956A`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFFFFFE`),
			Secure: internal.HB(`C682008189164C13`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FFFFFFFFFFFFFFFF`),
			Secure: internal.HB(`2CFFE4FF7938F7C0`),
		},
		{
			Key:    internal.HB(`80000000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`689DB471543EA251`),
		},
		{
			Key:    internal.HB(`C0000000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`FFC04D151996D7AE`),
		},
		{
			Key:    internal.HB(`E0000000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`FE5BA236C62BD3D3`),
		},
		{
			Key:    internal.HB(`F0000000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`6F5015A86E539A69`),
		},
		{
			Key:    internal.HB(`F8000000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`618D8E99D1D994FD`),
		},
		{
			Key:    internal.HB(`FC000000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`0014C6AF72089C84`),
		},
		{
			Key:    internal.HB(`FE000000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`31D234588FA1DD75`),
		},
		{
			Key:    internal.HB(`FF000000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`FB938357E2E3C5C8`),
		},
		{
			Key:    internal.HB(`FF800000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`FA0B03AB02062385`),
		},
		{
			Key:    internal.HB(`FFC00000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`AD582459CE145C06`),
		},
		{
			Key:    internal.HB(`FFE00000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`7180278EA10D1F39`),
		},
		{
			Key:    internal.HB(`FFF00000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`759294264B00B60B`),
		},
		{
			Key:    internal.HB(`FFF80000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`2B40B9330BDFD431`),
		},
		{
			Key:    internal.HB(`FFFC0000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`B523D37C7700F80A`),
		},
		{
			Key:    internal.HB(`FFFE0000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`E6634E05204AFC8B`),
		},
		{
			Key:    internal.HB(`FFFF0000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`00B852841E52A500`),
		},
		{
			Key:    internal.HB(`FFFF8000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`38F0CBE653586A1D`),
		},
		{
			Key:    internal.HB(`FFFFC000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`7957599D06476489`),
		},
		{
			Key:    internal.HB(`FFFFE000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`F366E7EE806DD8E2`),
		},
		{
			Key:    internal.HB(`FFFFF000000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`EF466A630C4EA3CA`),
		},
		{
			Key:    internal.HB(`FFFFF800000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`6D68D425FE466983`),
		},
		{
			Key:    internal.HB(`FFFFFC00000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`2BBA62781718967F`),
		},
		{
			Key:    internal.HB(`FFFFFE00000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`9A48E7348ACA78F2`),
		},
		{
			Key:    internal.HB(`FFFFFF00000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`CCA8B7CF972FEEA8`),
		},
		{
			Key:    internal.HB(`FFFFFF80000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`0BC7F645E439CBFE`),
		},
		{
			Key:    internal.HB(`FFFFFFC0000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`6B17CC2F07D50CE4`),
		},
		{
			Key:    internal.HB(`FFFFFFE0000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`F41C1F8D1990D156`),
		},
		{
			Key:    internal.HB(`FFFFFFF0000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`27B1D169B572D1FF`),
		},
		{
			Key:    internal.HB(`FFFFFFF8000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`5116DB6D3A5D386A`),
		},
		{
			Key:    internal.HB(`FFFFFFFC000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`52493EB7A4223D67`),
		},
		{
			Key:    internal.HB(`FFFFFFFE000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`157C173910462B8A`),
		},
		{
			Key:    internal.HB(`FFFFFFFF000000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`DF5FD74E4CD15050`),
		},
		{
			Key:    internal.HB(`FFFFFFFF800000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`7352CDAA640943CF`),
		},
		{
			Key:    internal.HB(`FFFFFFFFC00000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`C35DF749D46882A9`),
		},
		{
			Key:    internal.HB(`FFFFFFFFE00000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`3DD2DF5F888A0345`),
		},
		{
			Key:    internal.HB(`FFFFFFFFF00000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`DF25C0B6EFDACE1D`),
		},
		{
			Key:    internal.HB(`FFFFFFFFF80000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`C966B212515F14E5`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFC0000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`2F063EF30D59EB46`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFE0000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`994A1FB7A18A030D`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFF0000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`78895B803253450A`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFF8000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`1BA2A3E73797D6DE`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFC000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`CF00447566123EFC`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFE000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`4721F81D045A219E`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFF000000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`8E436F001945A175`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFF800000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`2B5EA8B643C03271`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFC00000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`859DD60FAFC06CBC`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFE00000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`7C815A6F59CA30D9`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFF00000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`FABD5F00B205917E`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFF80000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`0F096D670CA5844B`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFC0000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`BBAEAC969FE8469C`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFE0000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`0044A84CB41960BD`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFF0000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`3D65FBAAA4E0845A`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFF8000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`3000A49FDC7548F8`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFC000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`3B34B9FEE6DB12C4`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFE000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`3F2471D143F60AC1`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFF000000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`B97A37B5660C52FE`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFF800000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`50AD1E59F990FD57`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFC00000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`31694AB026FAA3B3`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFE00000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`447730C6421C39C5`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFF00000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`B1EDFFCEA8F5DBEF`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFF80000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`4DEEDAB4BA76E952`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFC0000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`499D672BEB6AF09C`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFE0000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`25644529D67AD993`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFF0000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`C5CFFB53D711799E`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFF8000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`2939DB3DAEE5E09A`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFC000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`45E1E8A6E038BBFC`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFE000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`0C7A1363F45F0503`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFF000000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`22BB26CC7D2F0198`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFF800000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`E2071CA91D33976B`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFC00000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`AAC1EC29B3CB9157`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFE00000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`2B8A686D43B73C7F`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFF00000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`5817E151C33574A0`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFF80000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`4EFBD547328D0884`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFC0000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`14AE645F897334D6`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFE0000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`10EB53CC021F6979`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFF0000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`825CAD8451D365EA`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFF8000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`3DA0138BD57125DF`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFC000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`8840F61B69F2AD37`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFE000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`20BF64B55B401BC9`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFF000000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`BC01D24BC626D6BA`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFF800000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`2AC8FA8B5EAFECF4`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFC00000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`A4259D5A2EDEE313`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFE00000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`5D4005ADAEE7DF04`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFF00000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`D28853BE7AF74B90`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFF80000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`3134F0F1A10448C7`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFC0000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`4E05A0308156F8F0`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFE0000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`325698888AC48513`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFF0000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`48B2B9422A13BE17`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFF8000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`E08ED667CDD9041C`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFC000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`A355C73358CA251C`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFE000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`98E2E9848D8C0144`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFF000000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`7F6B6F466D2DD369`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFF800000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`415E4101E78AE3B6`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFC00000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`4FC53722C90B9AD9`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFE00000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`94BE38F099EF7E49`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFF00000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`880BB83ED9B403BB`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFF80000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`97604D34A836176E`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFC0000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`53BB0964CBBE73D0`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFE0000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`E1E4042B6879E957`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFF0000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`F3AECC6016435BB8`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFF8000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`4AED3D640549D7FE`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFC000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`0483DB5A0DF14016`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFE000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`769B70A80F0F679E`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFF000000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`1BBCEF870305DF38`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFF800000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`6F0AD1CBDD56F2FB`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFC00000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`39E83DF7DA20C5DD`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFE00000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`4CB8E03FE4336E14`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFF00000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`36EC87377CDC5493`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFF80000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`EF032FB652A0D7D0`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFC0000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`08F9D570A0AA9EA7`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFE0000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`847D8692AB8AC657`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFF0000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`AED4D3318F2D8608`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFF8000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`1ED7AA3BC301F87A`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFFC000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`E773080EA5AB4919`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFFE000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`BDDCF059D3B84322`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFFF000`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`3C6ED713176FD619`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFFF800`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`CF153E8D2BC7E64A`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFFFC00`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`66138F4DB448F4C2`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFFFE00`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`7302DC25904F93D1`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`438E45925E710EEE`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFFFF80`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`85DDA75AA7C1BF7A`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFC0`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`0695F224F29DA9DE`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFE0`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`DB2A64467C5EBC2F`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF0`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`E8790E2442E6B224`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF8`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`2DF03D95E29441DD`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFC`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`F3A9353685BCEE6F`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFE`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`5311B20593300E76`),
		},
		{
			Key:    internal.HB(`FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF`),
			Plain:  internal.HB(`0000000000000000`),
			Secure: internal.HB(`802395603C089EBA`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`78471277A7E8E428`),
			Secure: internal.HB(`49B584704B58387D`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`6482E8E9F78BCB89`),
			Secure: internal.HB(`075C538007884EC8`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`82EBE2D23662070E`),
			Secure: internal.HB(`8362EDDCB792A66C`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`C624A4CA1CC70DD9`),
			Secure: internal.HB(`7B9C23C66A8FBC97`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`4C5F139CC6E8D681`),
			Secure: internal.HB(`72E64419D3F33CE1`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`CFE4DACA32445D90`),
			Secure: internal.HB(`392ADDB05E85EA6B`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`328CE807C4332503`),
			Secure: internal.HB(`4B4E3881B3DC4B9A`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`FB48F0BBC45CB0CD`),
			Secure: internal.HB(`1DE2F9A7F3D4BDDB`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`D39BE883DD3C0955`),
			Secure: internal.HB(`0D40068770A7132E`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`081C8BAFA0A53BF4`),
			Secure: internal.HB(`835BDF5486520C78`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`0BF5D7C2FF3CD878`),
			Secure: internal.HB(`F8BDEC208B42DED5`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`F2668FF5D2FB73A3`),
			Secure: internal.HB(`DF46E61C749EF745`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`F743BAB556AD26AC`),
			Secure: internal.HB(`6CAD68280131ABEB`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`F7712220A9750DBD`),
			Secure: internal.HB(`19178991061FAA0A`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`F46ED68C4F47C973`),
			Secure: internal.HB(`D7C56CCC3CB291FC`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`93C6A800B06DFE62`),
			Secure: internal.HB(`D456A893BCE8286D`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`9E9EB0B89703D58F`),
			Secure: internal.HB(`5357126F7EB9F286`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`842DC7A5B57B7AF3`),
			Secure: internal.HB(`6E950F179E145921`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`D43E0DEB1C1A9FFF`),
			Secure: internal.HB(`4B3FEFE2FE66F72D`),
		},
		{
			Key:    internal.HB(`00000000000000000000000000000000`),
			Plain:  internal.HB(`C6B06462C479F812`),
			Secure: internal.HB(`2207624A6B6A54CF`),
		},
	}
)

func TestEncryptDecrypt(t *testing.T) {
	plain := make([]byte, BlockSize)
	secure := make([]byte, BlockSize)

	for _, tc := range testCases {
		c, err := NewCipher(tc.Key)
		if err != nil {
			t.Error(err)
		}

		c.Encrypt(secure, tc.Plain)
		c.Decrypt(plain, secure)
		if !bytes.Equal(plain, tc.Plain) {
			t.Errorf("encrypt failed.\nresult: %s\nanswer: %s", hex.EncodeToString(plain), hex.EncodeToString(tc.Plain))
		}
	}
}

func TestEncrypt(t *testing.T) {
	dst := make([]byte, BlockSize)

	for _, tc := range testCases {
		c, err := NewCipher(tc.Key)
		if err != nil {
			t.Error(err)
		}

		c.Encrypt(dst, tc.Plain)
		if !bytes.Equal(dst, tc.Secure) {
			t.Errorf("encrypt failed.\nresult: %s\nanswer: %s", hex.EncodeToString(dst), hex.EncodeToString(tc.Secure))
		}
	}
}

func TestDecrypt(t *testing.T) {
	dst := make([]byte, BlockSize)

	for _, tc := range testCases {
		c, err := NewCipher(tc.Key)
		if err != nil {
			t.Error(err)
		}

		c.Decrypt(dst, tc.Secure)
		if !bytes.Equal(dst, tc.Plain) {
			t.Errorf("decrypt failed.\nresult: %s\nanswer: %s", hex.EncodeToString(dst), hex.EncodeToString(tc.Secure))
		}
	}
}
