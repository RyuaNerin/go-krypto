// drbg implemented deterministic random bit generator
// as defined in TTAK.KO-12.0331, TTAK.KO-12.0332
package drbg

import (
	"errors"
	"io"
)

type DRBG interface {
	io.Reader
	io.Closer

	Generate(dst []byte, length int, additionalInput []byte) ([]byte, error)
	Reseed(additionalInput []byte) error
}

var ErrUninstantiated = errors.New("krypto/drbg: state is uninstantiated")
