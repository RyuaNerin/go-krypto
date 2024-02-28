package main

import (
	"crypto/rand"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/RyuaNerin/go-krypto/aria"
	_ "github.com/RyuaNerin/go-krypto/hight"
	_ "github.com/RyuaNerin/go-krypto/lea"
	_ "github.com/RyuaNerin/go-krypto/seed"
)

var rnd = rand.Reader

func main() {
	var reqList []string
	if len(os.Args) > 1 {
		reqList = os.Args[1:]
	} else {
		var err error
		reqList, err = filepath.Glob("*.req")
		if err != nil {
			panic(err)
		}
		sort.Strings(reqList)
	}

	for idx, path := range reqList {
		log.Printf("[%d / %d] %s\n", idx+1, len(reqList), path)
		filename := filepath.Base(path)
		if filepath.Ext(filename) != ".req" {
			return
		}

		process(path, filename)
	}
}

func process(path, filename string) {
	filename = strings.ToUpper(filename)

	switch {
	//////////////////////////////////////////////////
	// Block
	case strings.HasPrefix(filename, "ARIA"):
		processBlock(path, filename)
	case strings.HasPrefix(filename, "HIGHT"):
		processBlock(path, filename)
	case strings.HasPrefix(filename, "LEA"):
		processBlock(path, filename)
	case strings.HasPrefix(filename, "SEED"):
		processBlock(path, filename)

	//////////////////////////////////////////////////
	// Hash
	case strings.HasPrefix(filename, "LSH"):
		processHash(path, filename)

	//////////////////////////////////////////////////
	// Sign
	case strings.HasPrefix(filename, "KCDSA"):
		processKCDSA(path, filename)

	case strings.HasPrefix(filename, "EC-KCDSA"):
		processECKCDSA(path, filename)

	//////////////////////////////////////////////////
	// HMAC
	case strings.HasPrefix(filename, "HMAC") && !strings.HasPrefix(filename, "HMAC_DRBG"):
		processHMAC(path, filename)

	//////////////////////////////////////////////////
	// CMAC
	case strings.HasPrefix(filename, "CMAC") && !strings.HasPrefix(filename, "CMAC_DRBG"):
		processCMAC(path, filename)

	//////////////////////////////////////////////////
	// HASH_DRBG
	case strings.HasPrefix(filename, "HASH_DRBG"):
		processDRBGHash(path, filename)

	//////////////////////////////////////////////////
	// HMAC_DRBG
	case strings.HasPrefix(filename, "HMAC_DRBG"):
		processDRBGHMAC(path, filename)

	//////////////////////////////////////////////////
	// CTR_DRBG
	case strings.HasPrefix(filename, "CTR_DRBG"):
		processDRBGCTR(path, filename)

	//////////////////////////////////////////////////
	// PBKDF
	case strings.HasPrefix(filename, "PBKDF"):
		processPBKDF(path, filename)

	//////////////////////////////////////////////////
	// KBKDF
	case strings.HasPrefix(filename, "KDF"):
		processKBKDF(path, filename)

	//////////////////////////////////////////////////
	// GCM
	case strings.Contains(filename, "GCM"):
		processGCM(path, filename)

	//////////////////////////////////////////////////
	// CCM
	case strings.Contains(filename, "CCM"):
		processCCM(path, filename)

	//////////////////////////////////////////////////
	default:
		log.Println("Unknown algorithm: ", filename)
	}
}
