[![PkgGoDev](https://pkg.go.dev/badge/github.com/RyuaNerin/go-krypto)](https://pkg.go.dev/github.com/RyuaNerin/go-krypto)

# krypto

Golang implementation of cryptographic algorithms designed by Republic of Korea

## Supported

- Block cipher

    | Algorithm | 128 | 192 | 256 |
    |:---------:|:---:|:---:|:---:|
    | ARIA      | O   | O   | O   |
    | HIGHT     | O   |     |     |
    | LEA       | O   | O   | O   |
    | SEED      | O   |     | O   |

- Hash

    | Algorithm | 224 | 256 | 384 | 512 |
    |:---------:|:---:|:---:|:---:|:---:|
    | LSH-256   | O   | O   |     |     |
    | LSH-512   | O   | O   | O   | O   |

- SIMD support
    | Algorithm | Support |
    |:---------:|:-:|
    | ARIA      | |
    | HIGHT     | |
    | LEA       | `SSE2` `AVX2` |
    | SEED      | |
    | LSH-256   | |
    | LSH-512   | |

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

## Performance

```txt
goos: windows
goarch: amd64
pkg: krypto
Benchmark_CBC_Encrypt_1K_AES-16          1210395           987 ns/op    1037.11 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Decrypt_1K_AES-16          1305638           909 ns/op    1126.19 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Encrypt_1K_SEED128-16        87049         13605 ns/op      75.27 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Decrypt_1K_SEED128-16        88154         13364 ns/op      76.63 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Encrypt_1K_SEED256-16        60858         19736 ns/op      51.89 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Decrypt_1K_SEED256-16        61482         19568 ns/op      52.33 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Encrypt_1K_HIGHT-16          71362         16754 ns/op      61.12 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Decrypt_1K_HIGHT-16          71362         17366 ns/op      58.97 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Encrypt_1K_ARIA_16-16        40778         28998 ns/op      35.31 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Decrypt_1K_ARIA_16-16        40816         29434 ns/op      34.79 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Encrypt_1K_ARIA_24-16        35157         33737 ns/op      30.35 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Decrypt_1K_ARIA_24-16        35680         33354 ns/op      30.70 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Encrypt_1K_ARIA_32-16        31408         38092 ns/op      26.88 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Decrypt_1K_ARIA_32-16        31466         37646 ns/op      27.20 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Encrypt_1K_LEA_16-16        428188          2831 ns/op     361.74 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Decrypt_1K_LEA_16-16        272478          4309 ns/op     237.62 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Encrypt_1K_LEA_24-16        374660          3120 ns/op     328.20 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Decrypt_1K_LEA_24-16        239780          4880 ns/op     209.85 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Encrypt_1K_LEA_32-16        342548          3342 ns/op     306.44 MB/s           0 B/op          0 allocs/op
Benchmark_CBC_Decrypt_1K_LEA_32-16        217982          5407 ns/op     189.39 MB/s           0 B/op          0 allocs/op
Benchmark_HASH_SHA256_1K-16               544958          2218 ns/op         32 B/op           1 allocs/op
Benchmark_HASH_SHA512_1K-16               799306          1562 ns/op         64 B/op           1 allocs/op
Benchmark_HASH_LSH256_1K-16               444034          2668 ns/op         32 B/op           1 allocs/op
Benchmark_HASH_LSH512_1K-16               499533          2329 ns/op         64 B/op           1 allocs/op
PASS
ok      krypto    29.203s
```
