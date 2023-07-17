package lsh256

import "hash"

var (
	newContext func(algType algType) hash.Hash = newContextGo
)
