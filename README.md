# [krypto](https://pkg.go.dev/github.com/RyuaNerin/go-krypto)

[![PkgGoDev](https://pkg.go.dev/badge/github.com/RyuaNerin/go-krypto)](https://pkg.go.dev/github.com/RyuaNerin/go-krypto)

- Golang implementation of cryptographic algorithms designed in Republic of Korea

- It is intended for compatibility with the `crypto` package.

- `krypto` does not required any other C/C++ compiler to use SIMD :\)

    - But, may have not been optimized in compiler level.

## [LICNESE](/LICENSE)

## Supported

- `N` : Not tested.
- `^` : Deprecated algorithm.

- Block cipher

    | Algorithm | Package        | Document           | 128 | 192 | 256 | SIMD Supports                    |
    |:---------:|----------------|:------------------:|:---:|:---:|:---:|:--------------------------------:|
    | SEED 128  | `krypto/seed`  | TTAS.KO-12.0004/R1 | O   |     |     |                                  |
    | SEED 256  | `krypto/seed`  | ***Unknown***      |     |     | N   |                                  |
    | HIGHT     | `krypto/hight` | TTAS.KO-12.0040/R1 | O   |     |     |                                  |
    | ARIA      | `krypto/aria`  | KS X 1213-1        | O   | O   | O   | amd64(SSSE3)                     |
    | LEA       | `krypto/lea`   | TTAK.KO-12.0223    | O   | O   | O   | arm64(NEON), amd64(SSE2, AVX2),  |

- Digital Signature

    | Algorithm | Package          | Document           |
    |:---------:|------------------|:------------------:|
    | KCDSA     | `krypto/kcdsa`   | TTAK.KO-12.0001/R4 |
    | EC-KCDSA  | `krypto/eckcdsa` | TTAK.KO-12.0015/R3 |

- Hash

    | Algorithm  | Package         | Document           | 160 | 224 | 256 | 384 | 512 | SIMD Supports                         |
    |:----------:|-----------------|:------------------:|:---:|:---:|:---:|:---:|:---:|:-------------------------------------:|
    | HAS-160^   | `krypto/has160` | TTAS.KO-12.0011/R2 | O   |     |     |     |     |                                       |
    | LSH-256    | `krypto/lsh256` | KS X 3262          |     | O   | O   |     |     | arm64(NEON), amd64(SSE2, SSSE3, AVX2) |
    | LSH-512    | `krypto/lsh512` | KS X 3262          |     | O   | O   | O   | O   | arm64(NEON), amd64(SSE2, SSSE3, AVX2) |

### SIMD Support

- It was based on the below.

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

## Usage

Todo

## TODO

- Supoorts Post-Quantum Cryptography
