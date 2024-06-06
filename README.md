# [krypto](https://pkg.go.dev/github.com/RyuaNerin/go-krypto)

[![PkgGoDev](https://pkg.go.dev/badge/github.com/RyuaNerin/go-krypto)](https://pkg.go.dev/github.com/RyuaNerin/go-krypto)

- Golang implementation of cryptographic algorithms designed in Republic of Korea

- It is intended for compatibility with go's `crypto` package.

- `krypto` provides SIMD for some algorithms.

## [LICNESE](/LICENSE)

## Installation

```shell
go get -v "github.com/RyuaNerin/go-krypto"
```

```go
package main

import (
    ...
    krypto "github.com/RyuaNerin/go-krypto"
    ...
)
```

## Supports

### Block Cipher Supports

| Algorithm | Package        | Document           | 128 | 192 | 256 | SIMD Supports                  |
|:---------:|----------------|:------------------:|:---:|:---:|:---:|:------------------------------:|
| SEED-128  | `krypto/seed`  | TTAS.KO-12.0004/R1 | O   |     |     |                                |
| HIGHT     | `krypto/hight` | TTAS.KO-12.0040/R1 | O   |     |     |                                |
| ARIA      | `krypto/aria`  | KS X 1213-1        | O   | O   | O   | arm64(NEON), amd64(SSSE3)      |
| LEA       | `krypto/lea`   | TTAK.KO-12.0223    | O   | O   | O   | arm64(NEON), amd64(SSE2, AVX2) |

- Use `krypto/kipher` for block mode supports.
    
    -  `crypto/cipher` package is available, but recommend `krypto/kipher` package for performance.
    
    - Supports

        - Block Mode

            - ECB (Electronic Codebook)
            - CBC (Cipher-Block Chaining)
            - CFB (Cipher Feedback) : CFB-8, CFG-32, ...
            - OFB (Output Feedback)
            - CTR (Counter)

        - AEAD (Authenticated Encryption with Associated Data)

            - GCM (Galois/Counter Mode)
            - CCM (Counter with CBC-MAC)

### Hash Function Supports

| Algorithm  | Package         | Document           | 160 | 224 | 256 | 384 | 512 | SIMD Supports                         |
|:----------:|-----------------|:------------------:|:---:|:---:|:---:|:---:|:---:|:-------------------------------------:|
| HAS-160    | `krypto/has160` | TTAS.KO-12.0011/R2 | O   |     |     |     |     |                                       |
| LSH-256    | `krypto/lsh256` | KS X 3262          |     | O   | O   |     |     | arm64(NEON), amd64(SSE2, SSSE3, AVX2) |
| LSH-512    | `krypto/lsh512` | KS X 3262          |     | O   | O   | O   | O   | arm64(NEON), amd64(SSE2, SSSE3, AVX2) |

### Digital Signature Supports

| Algorithm | Package          | Document           |
|:---------:|------------------|:------------------:|
| KCDSA     | `krypto/kcdsa`   | TTAK.KO-12.0001/R4 |
| EC-KCDSA  | `krypto/eckcdsa` | TTAK.KO-12.0015/R3 |

- use `krypto/kx509` for marshaling and unmarshaling of the private/public key.

    | Algorithm | Format                | Reference    | Comment                                                         |
    |:---------:|:---------------------:|:------------:|-----------------------------------------------------------------|
    | KCDSA     | PKIX, PKCS#8          | NO NORMATIVE | Compatibility tested with [jCastle](http://www.jcastle.net/)    |
    | EC-KCDSA  | PKIX, PKCS#8          | NO NORMATIVE | Compatibility tested with [botan](https://botan.randombit.net/) |
    | EC-KCDSA  | SEC 1, ASN.1 DER form | NO NORMATIVE |                                                                 |


### Message Authentication Code Supports

| Algorithm | Package       | Document            |
|:---------:|---------------|---------------------|
| CMAC      | `krypto/cmac` | KS X ISO/IEC 9797-1 |
| GMAC      | `krypto/gmac` | KS X ISO/IEC 9797-1 |

- use `crypto/hmac` for HMAC.

### Random Number Generator Supports

| Algorithm | Package       | Document            |
|:---------:|---------------|---------------------|
| Hash_DRBG | `krypto/drbg`  | TTAK.KO-12.0331    |
| HMAC_DRBG | `krypto/drbg`  | TTAK.KO-12.0332    |
| CTR_DRBG  | `krypto/drbg`  | TTAK.KO-12.0189/R1 |

### Key Derivation Function Supports

| Algorithm    | Package        | Document        |
|:------------:|----------------|-----------------|
| KBKDF (HMAC) | `krypto/kbkdf` | TTAK.KO-12.0272 |
| KBKDF (CMAC) | `krypto/kbkdf` | TTAK.KO-12.0272 |
| PBKDF (HMAC) | `krypto/pbkdf` | TTAK.KO-12.0334 |

### SIMD Support

| Algorithm | SIMD Supports                         | Reference                                                   |
|:---------:|---------------------------------------|:-----------------------------------------------------------:|
| ARIA      | arm64(NEON), amd64(SSSE3)             | [CRYPTOPP 8.8.0 - aria_simd.cpp](https://github.com/weidai11/cryptopp/blob/CRYPTOPP_8_8_0/aria_simd.cpp) |
| LEA       | arm64(NEON), amd64(SSE2, AVX2)        | [KISA](https://seed.kisa.or.kr/kisa/Board/20/detailView.do) |
| LSH-256   | arm64(NEON), amd64(SSE2, SSSE3, AVX2) | [KISA](https://seed.kisa.or.kr/kisa/Board/22/detailView.do) |
| LSH-512   | arm64(NEON), amd64(SSE2, SSSE3, AVX2) | [KISA](https://seed.kisa.or.kr/kisa/Board/22/detailView.do) |

- The draft of the assembly code was created by clang and modifying verseion of the program below on MacMini M1.

    - [gorse-io/goat](https://github.com/gorse-io/goat)
    - [c2goasm](https://github.com/minio/c2goasm)

## [Performance](/PERFORMANCE.md)

## Usage

Todo

## TODO

- Supoorts Post-Quantum Cryptography
