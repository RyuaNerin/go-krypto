package elliptic2m

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
)

func Test_IsOnCurve_K283(t *testing.T) {
	testPoint(t, testcase_K283_PKV, K283())
}

// //////////////////////////////////////////////////////////////////////////////////////////////////
// 암호알고리즘 검증기준 V3.0
// 테스트 벡터
// ECDH_K-283_PKV.txt
var testcase_K283_PKV = []testCase{
	{
		Qx:   internal.HI(`59244BB3C44EBE2A6856844247BDBB807D9ECA7813AE660F38ABFA5B0E667544426C672`),
		Qy:   internal.HI(`2CB1699793A88ECEEEDA5CE792B79A994CD6B3EA634448EE7BF3AB3A4A9E42EE88B6B9E`),
		Fail: false,
	},
	{
		Qx:   internal.HI(`2F9EF2161F4B551668A54A8C8EABD9F03D04F8170CE85FB50D3B7C431254CB6477D13DD`),
		Qy:   internal.HI(`7C6E9902BA11F861F92D1608E48CEB8F699476345553AE5B806E97C0EFC5037CAD0942D`),
		Fail: false,
	},
	{
		Qx:   internal.HI(`400897863581AC554A9E43CD686A7AC8B8E2818410281753CE1CE3DEFCA58E1EBE4878B`),
		Qy:   internal.HI(`4F2ADC1FB3D2CA6A4D4F276E2C6CDDEDC8091F473F354C9C8D69B8E38026F95042BF7F5`),
		Fail: false,
	},
	{
		Qx:   internal.HI(`6BDD79203F7465127358CF463DBE8619D2D3301CDFC0194B7480D5A187DF4286FA664AB`),
		Qy:   internal.HI(`490886019687CE5C076BA82196397CC797EC71575DEFFAC01F1316052D13FDB69D0BDCC`),
		Fail: false,
	},
	{
		Qx:   internal.HI(`774A8F05D01408798FD00F09DB08EC3A9D4B6F2F574B988C658B3092035A3C828CDE943`),
		Qy:   internal.HI(`7CEA107492390EB625A1C012AA8185A2DDFEC2381AFD89047F74572CC9E3059898C945B`),
		Fail: false,
	},
	{
		Qx:   internal.HI(`386FA56D9B99B6DC9F90BA1BB7AE7167027AF7301FDB1869EB9AFA2EA456B834A5AAE55`),
		Qy:   internal.HI(`3D7515A284133985872FDD628FA9E64FA6EB2627264D601AF5F8F9D67FA54C1650D8E25`),
		Fail: false,
	},
	{
		Qx:   internal.HI(`63BC179295CFB0CE67CAE56CB72EF6707F73BE06EAE9C84E397410B3C20D110DCEC29C9`),
		Qy:   internal.HI(`44F1AE5C3FC5313233B21F404B4D07E1E1F169727E912FA986DA3F9CAF83B798485F0C0`),
		Fail: true,
	},
	{
		Qx:   internal.HI(`38B71F807BC97992EAEA84019CE93E33F079840D0E207E8CDD5EB92A9BF1B19E5F4BFC8`),
		Qy:   internal.HI(`159951059AC51F9FBA61CC547AA3B4E2D4EF07D654589C2ADEB64CDC4C730BC47000D7F`),
		Fail: true,
	},
	{
		Qx:   internal.HI(`334AA90E765CECE1F432158CB8058E802FB750ED8CD55C193B794BDCD215B00B7E5344`),
		Qy:   internal.HI(`1F0B5985B3DC98C57B0C425B6BF01AF9C894A981CED8B9673FC25B089A8300E39C2622B`),
		Fail: false,
	},
	{
		Qx:   internal.HI(`5BDFBB852990EC212179166A51432B5B53F5E36C6FD8AF2BDFC926BFFB2AD25646A66B9`),
		Qy:   internal.HI(`5F4B2C8A8EDD1C897CF20517BED4924B4560A730C695929734832ED1329EC4DF9AAB955`),
		Fail: false,
	},
}
