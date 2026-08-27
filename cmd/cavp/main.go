package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <dir|zipfile>", filepath.Base(os.Args[0]))
	}

	if err := run(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}

func process(path, filename string) error {
	filename = strings.ToUpper(filename)

	switch {
	case strings.HasPrefix(filename, "ARIA"),
		strings.HasPrefix(filename, "HIGHT"),
		strings.HasPrefix(filename, "LEA"),
		strings.HasPrefix(filename, "SEED"):
		return processBlock(path, filename)
	case strings.HasPrefix(filename, "LSH"):
		return processHash(path, filename)
	case strings.HasPrefix(filename, "KCDSA"):
		return processKCDSA(path, filename)
	case strings.HasPrefix(filename, "EC-KCDSA"):
		return processECKCDSA(path, filename)
	case strings.HasPrefix(filename, "HMAC") && !strings.HasPrefix(filename, "HMAC_DRBG"):
		return processHMAC(path, filename)
	case strings.HasPrefix(filename, "CMAC") && !strings.HasPrefix(filename, "CMAC_DRBG"):
		return processCMAC(path, filename)
	case strings.HasPrefix(filename, "HASH_DRBG"):
		return processDRBGHash(path, filename)
	case strings.HasPrefix(filename, "HMAC_DRBG"):
		return processDRBGHMAC(path, filename)
	case strings.HasPrefix(filename, "CTR_DRBG"):
		return processDRBGCTR(path, filename)
	case strings.HasPrefix(filename, "PBKDF"):
		return processPBKDF(path, filename)
	case strings.HasPrefix(filename, "KDF"):
		return processKBKDF(path, filename)
	case strings.Contains(filename, "GCM"):
		return processGCM(path, filename)
	case strings.Contains(filename, "CCM"):
		return processCCM(path, filename)
	default:
		return unknownAlgorithm(filename)
	}
}

func responsePath(requestPath string) string {
	ext := filepath.Ext(requestPath)
	return strings.TrimSuffix(requestPath, ext) + ".rsp"
}

var errUnknownAlgorithm = errors.New("unknown algorithm")

func unknownAlgorithm(filename string) error {
	return fmt.Errorf("%w: %s", errUnknownAlgorithm, filename)
}
