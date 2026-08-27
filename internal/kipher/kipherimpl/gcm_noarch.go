//go:build !(gc && (arm64 || amd64) && !noasm && !noasm_krypto && !noasm_krypto_kipher)

package kipherimpl

func init() {
	gcmInit = GCMInitGo
	gcmDeriveCounter = GCMDeriveCounterGo
	gcmUpdate = GCMUpdateGo
	gcmAuth = GCMAuthGo
	gcmFinish = GCMFinishGo
}
