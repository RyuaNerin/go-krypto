[![PkgGoDev](https://pkg.go.dev/badge/github.com/RyuaNerin/go-krypto)](https://pkg.go.dev/github.com/RyuaNerin/go-krypto)

# krypto

- Golang implementation of cryptographic algorithms designed in Republic of Korea

- It is intended for compatibility with the `crypto` package.

## [LICNESE](/LICENSE)

## Supported

- `A` : Not tested.
- `*` : deprecated.

- Block cipher

    | Algorithm | Package        | Document           | 128 | 192 | 256 | SIMD Supports           |
    |:---------:|----------------|:------------------:|:---:|:---:|:---:|:-----------------------:|
    | SEED 128  | `krypto/seed`  | TTAS.KO-12.0004/R1 | O   |     |     |                         |
    | SEED 256  | `krypto/seed`  | Unknown            |     |     | A   |                         |
    | HIGHT     | `krypto/hight` | TTAS.KO-12.0040/R1 | O   |     |     |                         |
    | ARIA      | `krypto/aria`  | KS X 1213-1        | O   | O   | O   |                         |
    | LEA       | `krypto/lea`   | TTAK.KO-12.0223    | O   | O   | O   | `SSE2`, `SSSE3`, `AVX2` |

- Digital Signature

    | Algorithm | Package          | Document           |
    |:---------:|------------------|:------------------:|
    | KCDSA     | `krypto/kcdsa`   | TTAK.KO-12.0001/R4 |
    | EC-KCDSA  | `krypto/eckcdsa` | TTAK.KO-12.0015/R3 |

- Hash

    | Algorithm  | Package         | Document           | 160 | 224 | 256 | 384 | 512 | SIMD Supports         |
    |:----------:|-----------------|:------------------:|:---:|:---:|:---:|:---:|:---:|:---------------------:|
    | HAS-160`*` | `krypto/has160` | TTAS.KO-12.0011/R2 | O   |     |     |     |     |                       |
    | LSH-256    | `krypto/lsh256` | KS X 3262          |     | O   | O   |     |     | `SSE2` `SSSE3` `AVX2` |
    | LSH-512    | `krypto/lsh512` | KS X 3262          |     | O   | O   | O   | O   |                       |

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
