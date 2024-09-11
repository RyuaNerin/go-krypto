//go:build !(gc && (arm64 || amd64) && !purego && go1.18)
// +build !gc !arm64,!amd64 purego !go1.18

package kipher

func init() {
	gcmInit = gcmInitGo
	gcmDeriveCounter = gcmDeriveCounterGo
	gcmUpdate = gcmUpdateGo
	gcmAuth = gcmAuthGo
	gcmFinish = gcmFinishGo
}
