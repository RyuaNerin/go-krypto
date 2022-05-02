//go:build !amd64

package lsh256

import "hash"

func newContext(algType algType) hash.Hash {
	return newContextGo(algType)
}
