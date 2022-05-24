[![PkgGoDev](https://pkg.go.dev/badge/github.com/RyuaNerin/go-krypto)](https://pkg.go.dev/github.com/RyuaNerin/go-krypto)

# krypto

Golang implementation of cryptographic algorithms designed by Republic of Korea

## [LICNESE](/LICENSE)

## Supported

- Block cipher

    | Algorithm | 128 | 192 | 256 |SIMD Support   |
    |:---------:|:---:|:---:|:---:|:-------------:|
    | ARIA      | O   | O   | O   |               |
    | HIGHT     | O   |     |     |               |
    | LEA       | O   | O   | O   | `SSE2` `AVX2` |
    | SEED      | O   |     | O   |               |

- Hash

    | Algorithm | 224 | 256 | 384 | 512 |
    |:---------:|:---:|:---:|:---:|:---:|
    | LSH-256   | O   | O   |     |     |
    | LSH-512   | O   | O   | O   | O   |

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
