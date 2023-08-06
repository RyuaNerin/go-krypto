[![PkgGoDev](https://pkg.go.dev/badge/github.com/RyuaNerin/go-krypto)](https://pkg.go.dev/github.com/RyuaNerin/go-krypto)

# krypto

Golang implementation of cryptographic algorithms designed by Republic of Korea

## [LICNESE](/LICENSE)

## Supported

- `A` : not tested

- Block cipher

    | Algorithm | Document           | 128 | 192 | 256 | SIMD Supports           |
    |:---------:|:------------------:|:---:|:---:|:---:|:-----------------------:|
    | SEED      | TTAS.KO-12.0004/R1 | O   |     | A   |                         |
    | HIGHT     | TTAS.KO-12.0040/R1 | O   |     |     |                         |
    | ARIA      | KS X 1213-1        | O   | O   | O   |                         |
    | LEA       | TTAK.KO-12.0223    | O   | O   | O   | `SSE2`, `SSSE3`, `AVX2` |

- Digital Signature

    | Algorithm | Document           | Note                                           s|
    |:---------:|:------------------:|:----------------------------------------------:|
    | KCDSA     | TTAK.KO-12.0001/R4 |                                                |
    | EC-KCDSA  | TTAK.KO-12.0015/R3 | Not tested: `B-233`, `B-283`, `K-233`, `K-283` |

- Hash

    | Algorithm | Document  | 224 | 256 | 384 | 512 | SIMD Supports         |
    |:---------:|:---------:|:---:|:---:|:---:|:---:|:---------------------:|
    | LSH-256   | KS X 3262 | O   | O   |     |     | `SSE2` `SSSE3` `AVX2` |
    | LSH-512   | KS X 3262 | O   | O   | O   | O   |                       |

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
