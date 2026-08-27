package cavp

import (
	"crypto/cipher"
	"fmt"

	"github.com/RyuaNerin/go-krypto/kipher"
)

// BlockMode is a confidentiality mode supported by the block-cipher CAVP
// vectors. CFB1 is intentionally omitted because kipher operates on bytes.
type BlockMode uint8

const (
	ECB BlockMode = iota
	CBC
	CFB8
	CFB32
	CFB64
	CFB128
	OFB
	CTR
)

// SegmentSize returns the CFB segment size in bytes, or zero for other modes.
func (m BlockMode) SegmentSize() int {
	switch m {
	case CFB8:
		return 1
	case CFB32:
		return 4
	case CFB64:
		return 8
	case CFB128:
		return 16
	default:
		return 0
	}
}

func (m BlockMode) String() string {
	switch m {
	case ECB:
		return "ECB"
	case CBC:
		return "CBC"
	case CFB8:
		return "CFB8"
	case CFB32:
		return "CFB32"
	case CFB64:
		return "CFB64"
	case CFB128:
		return "CFB128"
	case OFB:
		return "OFB"
	case CTR:
		return "CTR"
	default:
		return fmt.Sprintf("BlockMode(%d)", m)
	}
}

// CryptBlock processes one CAVP KAT/MMT message. For CTR, counter is used
// when present and IV is accepted as the counter otherwise.
func CryptBlock(block cipher.Block, mode BlockMode, iv, counter, dst, src []byte, encrypt bool) error {
	fn, err := NewBlockCrypt(block, mode, iv, counter, encrypt)
	if err != nil {
		return err
	}
	fn(dst, src)
	return nil
}

// NewBlockCrypt creates one stateful block-mode operation. It is a low-level
// kipher adapter; CAVP record processing belongs in each algorithm's ovs
// package.
func NewBlockCrypt(block cipher.Block, mode BlockMode, iv, counter []byte, encrypt bool) (func([]byte, []byte), error) {
	if block == nil {
		return nil, fmt.Errorf("cavp: nil block cipher")
	}
	if mode == CTR && len(counter) != 0 {
		iv = counter
	}
	if mode != ECB && len(iv) != block.BlockSize() {
		return nil, fmt.Errorf("cavp: %s requires a %d-byte IV", mode, block.BlockSize())
	}

	switch mode {
	case ECB:
		if encrypt {
			return kipher.NewECBEncryptor(block).CryptBlocks, nil
		} else {
			return kipher.NewECBDecryptor(block).CryptBlocks, nil
		}
	case CBC:
		if encrypt {
			return kipher.NewCBCEncrypter(block, iv).CryptBlocks, nil
		} else {
			return kipher.NewCBCDecrypter(block, iv).CryptBlocks, nil
		}
	case CFB8, CFB32, CFB64, CFB128:
		if encrypt {
			return kipher.NewCFBEncrypterWithBlockSize(block, iv, mode.SegmentSize()).XORKeyStream, nil
		} else {
			return kipher.NewCFBDecrypterWithBlockSize(block, iv, mode.SegmentSize()).XORKeyStream, nil
		}
	case OFB:
		return kipher.NewOFB(block, iv).XORKeyStream, nil
	case CTR:
		return kipher.NewCTR(block, iv).XORKeyStream, nil
	default:
		return nil, fmt.Errorf("cavp: unsupported block mode %d", mode)
	}
}
