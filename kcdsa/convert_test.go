package kcdsa

import (
	"crypto/dsa"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
)

var (
	testPrivateKeyP = internal.HI(`b71a75bae95ba4536fc794cb9e585b92d0df42640481dd641e4001d048c0053f99f56bfee756b6012f9b032fe5ffcab568a9f6472dc4b07a5d4f2f0f8190244f7408fd6774240efd9de0393e63168476f6e8ed1ee8d8ebb4974a397fac8e5fd7d1df33e21550bb172d75ced1cfb6fd6a4a5e5da932e32f7f1ef6a97b7f3ddcb0ff07dd850e8f6e569536e1d13654491d54169a2f5c40e6fb21949a0dfb3f8f8f78924f95533962f9fa2874aefc64614b3e5da905f96b0fca4709a6677c47252208d0aef29d4bd976f31b7fbac4254d14e833a822449efa588953f6442b4b8e2bd340a85a614c9a5e0548f0699cb5b0d95d049513f67e3059a1be4b5175a73b3b`)
	testPrivateKeyQ = internal.HI(`a446f4cce4a1d3e6231f9cbd0a5cd63410f7d9d0d433ca3251d8c6c3`)
	testPrivateKeyG = internal.HI(`b0964ade6c7b05e0d8efd4445be89600f01ca19799eb9beb5a1020ab1485c6551b695f5653d9c939459dc8dba8c83532d8d207f32dc4a5fecdc87ad83dd568088ecf6a82730d969344c49f9467bcb8c99c4ac12039d88f1b6b2e101859a8d8dbbade66fa846a57bd94e82e16977fecbaa44189e6734f53fc8041bb9ef64a0812edd4aec70cc4bf4117c75676ecfe913f7d0232a92bfa76e90975f96f7cdb416855f0d2dc03a9b6b1322fa2b3eced48f4aac4c5b991fb5b6918f36fdb7f950db345f4aace178c7ed5af8f19f1ecb61be8f86895c34239d7be2750e1705e1b551ef66ffdb58f6fde4a3dd064949461aa21a1e11fd050155fc0e8f60073b2ae2ff1`)
	testPrivateKeyX = internal.HI(`80b32101220b1acd6d961e3d8b3dbf458ed15ffbd61fc6cf2ba3b167`)

	testPrivateKeyY_DSA   = internal.HI(`64956569fc9c18200fa82569a92d415120910a62173bfe56453ed5cd6022c60c5556d986ccd3eff288f4a066c328707f936a7e7b650ffdc4a9e8c45a60c8c500b9ca4921ba18c0373b5863695a0b18dd5d8e84a402712e716547ecfa8e7827a950dde220fa4c2ee889c81175b2281d8102f4b51c3879e1e2c30ec4f6113fc8b78a80081178f879fa5d0cd72c0b2a1199e6df78366660fd03681a31e9ee879984ec25a40fd84576577c0748ee9747a667c13b0f446ca5745a37526e1265ee8b96e26ca7d7c14a227c0ce5b2227c0b69765b614bf5b24e6e9ab837d40828402b0273b36550ed4d7573772dce958259b25adafb20adf2b04ff12ee3612ccbd8d443`)
	testPrivateKeyY_KCDSA = internal.HI(`538518573d1847afffad49f08208c403405e2a93ed5e1bf6a86778d55194dba846a84ebc88f02b82f1a357ba75d2d98c9bd266d58e075f2b083d3e96e6beef8939f3df8e7d1f1ef60f3c65b060db2a907427502b85cc1bd9ec237a8820ad06c000626516381a968d6810acb57aeaa8f6b75b090ec8234f69b3f9e9b5d1aaf148eb17960d8f3317aae9b23792317c15d421c80f6f50990a18c4ac1659e738cde34ef0e938d240fb994cc5adfb2eb22d43c9b137e146c531a0efabbc8b5db0ae45200849350c0ecaad7eba0e8336bae6cb193f558fcdf6eb6392a159e8897fb18a27ccba2531d644b5df6e60192403f13f6c10d457c5372a1b9a23441c958317dd`)
)

func TestPrivateKeyConversion(t *testing.T) {
	privKCDSA := PrivateKey{
		PublicKey: PublicKey{
			Parameters: Parameters{
				P: testPrivateKeyP,
				Q: testPrivateKeyQ,
				G: testPrivateKeyG,
			},
			Y: testPrivateKeyY_KCDSA,
		},
		X: testPrivateKeyX,
	}
	privDSA := dsa.PrivateKey{
		PublicKey: dsa.PublicKey{
			Parameters: dsa.Parameters{
				P: testPrivateKeyP,
				Q: testPrivateKeyQ,
				G: testPrivateKeyG,
			},
			Y: testPrivateKeyY_DSA,
		},
		X: testPrivateKeyX,
	}

	toDSA := privKCDSA.ToDSA()
	if toDSA.X.Cmp(privDSA.X) != 0 ||
		toDSA.Y.Cmp(privDSA.Y) != 0 ||
		toDSA.P.Cmp(privDSA.P) != 0 ||
		toDSA.Q.Cmp(privDSA.Q) != 0 ||
		toDSA.G.Cmp(privDSA.G) != 0 {
		t.Fail()
	}

	toKCDSA := FromDSA(&privDSA)
	if toKCDSA.X.Cmp(privKCDSA.X) != 0 ||
		toKCDSA.Y.Cmp(privKCDSA.Y) != 0 ||
		toKCDSA.P.Cmp(privKCDSA.P) != 0 ||
		toKCDSA.Q.Cmp(privKCDSA.Q) != 0 ||
		toKCDSA.G.Cmp(privKCDSA.G) != 0 {
		t.Fail()
	}
}
