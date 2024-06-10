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

| Algorithm | Package        | Reference          | 128 | 192 | 256 | SIMD Supports                  |
|:---------:|----------------|:------------------:|:---:|:---:|:---:|:------------------------------:|
| SEED-128  | `krypto/seed`  | TTAS.KO-12.0004/R1 | O   |     |     |                                |
| HIGHT     | `krypto/hight` | TTAS.KO-12.0040/R1 | O   |     |     |                                |
| ARIA      | `krypto/aria`  | KS X 1213-1        | O   | O   | O   | arm64(NEON), amd64(SSSE3)      |
| LEA       | `krypto/lea`   | TTAK.KO-12.0223    | O   | O   | O   | arm64(NEON), amd64(SSE2, AVX2) |

#### Block Cipher Mode Supports

- `crypto/cipher` package is available too.

| Mode  | Name                        | Reference       | Comment                     |
|:-----:|:---------------------------:|:---------------:|-----------------------------|
| Block | ECB (Electronic Codebook)   | NIST SP 800-38A |                             |
| Block | CBC (Cipher-Block Chaining) | NIST SP 800-38A |                             |
| Block | CFB (Cipher Feedback)       | NIST SP 800-38A | Supports CFB-8, CFG-32, ... |
| Block | OFB (Output Feedback)       | NIST SP 800-38A |                             |
| Block | CTR (Counter)               | NIST SP 800-38A |                             |
| AEAD  | CCM (Counter with CBC-MAC)  | NIST SP 800-38C |                             |
| AEAD  | GCM (Galois/Counter Mode)   | NIST SP 800-38D |                             |

### Hash Function Supports

| Algorithm  | Package         | Reference          | 160 | 224 | 256 | 384 | 512 | SIMD Supports                         |
|:----------:|-----------------|:------------------:|:---:|:---:|:---:|:---:|:---:|:-------------------------------------:|
| HAS-160    | `krypto/has160` | TTAS.KO-12.0011/R2 | O   |     |     |     |     |                                       |
| LSH-256    | `krypto/lsh256` | KS X 3262          |     | O   | O   |     |     | arm64(NEON), amd64(SSE2, SSSE3, AVX2) |
| LSH-512    | `krypto/lsh512` | KS X 3262          |     | O   | O   | O   | O   | arm64(NEON), amd64(SSE2, SSSE3, AVX2) |

### Digital Signature Supports

| Algorithm | Package          | Reference          |
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

| Algorithm | Package       | Reference                            |
|:---------:|---------------|--------------------------------------|
| CMAC      | `krypto/cmac` | KS X ISO/IEC 9797-1, NIST SP 800-38B |
| GMAC      | `krypto/gmac` | KS X ISO/IEC 9797-3, NIST SP 800-38D |

- use `crypto/hmac` for HMAC.

### Random Number Generator Supports

| Algorithm | Package       | Reference                            |
|:---------:|---------------|--------------------------------------|
| Hash_DRBG | `krypto/drbg`  | TTAK.KO-12.0331, NIST SP 800-90A    |
| HMAC_DRBG | `krypto/drbg`  | TTAK.KO-12.0332, NIST SP 800-90A    |
| CTR_DRBG  | `krypto/drbg`  | TTAK.KO-12.0189/R1, NIST SP 800-90A |

### Key Derivation Function Supports

| Algorithm     | Package       - | Reference                                           |
|:-------------:|-----------------|-----------------------------------------------------|
| KBKDF (CMAC)  | `krypto/kbkdf`  | TTAK.KO-12.0272, NIST SP 800-108                    |
| KBKDF (HMAC)  | `krypto/kbkdf`  | TTAK.KO-12.0333, NIST SP 800-108                    |
| PBKDF2 (HMAC) | `krypto/pbkdf2` | TTAK.KO-12.0334, NIST SP 800-132, RFC 2898(PKCS #5) |

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
