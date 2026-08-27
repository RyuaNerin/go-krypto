package main

import "github.com/RyuaNerin/go-krypto/internal/cavp"

func processPBKDF(path, filename string) error {
	hashInfo := getHash(filename)
	if hashInfo == nil {
		return unknownAlgorithm(filename)
	}
	return cavp.ProcessPbkdfKAT(path, responsePath(path), hashInfo.New)
}
