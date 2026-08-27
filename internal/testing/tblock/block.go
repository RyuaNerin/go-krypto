package tblock

type CipherInfo struct {
	BlockSize int
	Sizes     []CipherSize
}
type CipherSize struct {
	Name        string
	KeySize     int // bytes
	KATFilename string
}
