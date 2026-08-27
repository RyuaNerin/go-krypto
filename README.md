# 🔒 [go-krypto](https://pkg.go.dev/github.com/RyuaNerin/go-krypto)

[![PkgGoDev](https://pkg.go.dev/badge/github.com/RyuaNerin/go-krypto)](https://pkg.go.dev/github.com/RyuaNerin/go-krypto)

- Golang implementation of cryptographic algorithms designed in Republic of Korea

- It is intended for compatibility with go's `crypto` package.

- `go-krypto` supports SIMD for some algorithms.

## [LICENSE](/LICENSE)

- Apache License 2.0

## Reporting a Vulnerability

- Please send a mail to [`admin@ryuar.in`](mailto:RyuaNerin<admin@ryuar.in>) when you found ***a potential security issue in go-krypto***.

## Requirements

- `go-krypto` requires **`Go 1.25 or later`**.

- `go-krypto` uses [`github.com/RyuaNerin/elliptic2`](https://github.com/RyuaNerin/elliptic2) for Elliptic curves over binary fields (e.g. sect233k1, sect233r1, K-233, ...) in `EC-KCDSA` algorithm.

## Installation

- use `go get` command to install `go-krypto` package.

    ```shell
    go get -v "github.com/RyuaNerin/go-krypto"
    ```

- Sample code available here: [`example/main.go`](example/main.go).

## Supports

### SIMD

- `go-krypto` supports SIMD for some algorithms for better performance.

- if you want to disable SIMD, you can use `noasm_krypto` or `noasm_krypto_<algorithm>` build tag.

    ```sh
    $ go build -tags noasm_krypto
    $ go build -tags noasm_krypto_lea
    ```

- available build tags for disabling SIMD:

    | Algorithm Type | Algorithm | arm64 | amd64                   | Disable tag           | Reference                                                                                                |
    |:--------------:|:---------:|-------|-------------------------|:---------------------:|----------------------------------------------------------------------------------------------------------|
    | Block          | LEA       | NEON  | SSE2 / AVX2 / AVX512    | `noasm_krypto_lea`    |          [KISA](https://seed.kisa.or.kr/kisa/Board/20/detailView.do)                                     |
    | Hash           | LSH-256   | NEON  | SSE2 / SSSE3 / ~~AVX2~~ | `noasm_krypto_lsh256` |          [KISA](https://seed.kisa.or.kr/kisa/Board/20/detailView.do)                                     |
    | Hash           | LSH-256   |       | AVX512                  | `noasm_krypto_lsh256` | based on [KISA](https://seed.kisa.or.kr/kisa/Board/20/detailView.do)                                     |
    | Hash           | LSH-512   | NEON  | SSE2 / SSSE3 / AVX2     | `noasm_krypto_lsh512` |          [KISA](https://seed.kisa.or.kr/kisa/Board/20/detailView.do)                                     |
    | Hash           | LSH-512   |       | AVX512                  | `noasm_krypto_lsh512` | based on [KISA](https://seed.kisa.or.kr/kisa/Board/20/detailView.do)                                     |
    | MAC            | GMAC      | PMULL | PCLMULQDQ               | `noasm_krypto_kipher` | [go 1.22.4 `crypto/aes`](https://github.com/golang/go/tree/go1.22.4/src/crypto/aes)                      |
    | Block Mode     | GCM       | PMULL | PCLMULQDQ               | `noasm_krypto_kipher` | [go 1.22.4 `crypto/aes`](https://github.com/golang/go/tree/go1.22.4/src/crypto/aes)                      |

- C source code to Go assembly conversion uses [gorse-io/goat](https://github.com/gorse-io/goat).

    - AVX2 mode in LSH-256 is disabled because it is slower than SSSE3 mode in some CPU.

        -  if you want to enable AVX2 mode in LSH-256, you can use `krypto_enable_lsh256_avx2` build tag.

            ```sh
            $ go build -tags krypto_enable_lsh256_avx2
            ```

- example:

    | description                         | go build command                                         |
    |-------------------------------------|----------------------------------------------------------|
    | disable SIMD for all of `go-krypto` | `go build -tags noasm_krypto`                            |
    | disable SIMD for LEA, GCM(GCM)      | `go build -tags noasm_krypto_lea,noasm_krypto_kipher`    |
    | disable SIMD for LSH-256, LSH-512   | `go build -tags noasm_krypto_lsh256,noasm_krypto_lsh512` |

### Block Cipher Supports

| Algorithm | Package        | Reference          | 128 | 192 | 256 |
|:---------:|----------------|:------------------:|:---:|:---:|:---:|
| SEED-128  | `krypto/seed`  | TTAS.KO-12.0004/R1 | O   |     |     |
| HIGHT     | `krypto/hight` | TTAS.KO-12.0040/R1 | O   |     |     |
| ARIA      | `krypto/aria`  | KS X 1213-1        | O   | O   | O   |
| LEA       | `krypto/lea`   | TTAK.KO-12.0223    | O   | O   | O   |

- **use `krypto/kipher` package for block cipher mode instead of `crypto/cipher` package.**
    - `krypto/kipher` is sometimes faster than `crypto/cipher` package because it supports SIMD or multi-block processing for some algorithms.
    - `krypto/kipher` package is compatible with `crypto/cipher` package.

#### Block Cipher Mode Supports

- package `krypto/kipher` supports block cipher modes.

    - `crypto/cipher` package is available too.

| Mode  | Name                        | Reference       | SIMD Supports                  | Comment                     |
|:-----:|:---------------------------:|:---------------:|--------------------------------|-----------------------------|
| Block | ECB (Electronic Codebook)   | NIST SP 800-38A | `krypto/lea`                   |                             |
| Block | CBC (Cipher-Block Chaining) | NIST SP 800-38A |                                |                             |
| Block | CFB (Cipher Feedback)       | NIST SP 800-38A | `krypto/lea` (decrypt)         | Supports CFB-8, CFB-32, ... |
| Block | OFB (Output Feedback)       | NIST SP 800-38A | `krypto/lea`                   |                             |
| Block | CTR (Counter)               | NIST SP 800-38A |                                |                             |
| AEAD  | CCM (Counter with CBC-MAC)  | NIST SP 800-38C |                                |                             |
| AEAD  | GCM (Galois/Counter Mode)   | NIST SP 800-38D | arm64(PMULL), amd64(PCLMULQDQ) |                             |

### Hash Function Supports

| Algorithm  | Package         | Reference          | 160 | 224 | 256 | 384 | 512 |
|:----------:|-----------------|:------------------:|:---:|:---:|:---:|:---:|:---:|
| HAS-160    | `krypto/has160` | TTAS.KO-12.0011/R2 | O   |     |     |     |     |
| LSH-256    | `krypto/lsh256` | KS X 3262          |     | O   | O   |     |     |
| LSH-512    | `krypto/lsh512` | KS X 3262          |     | O   | O   | O   | O   |

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
| HMAC      | `krypto/hmac` | KS X ISO/IEC 9797-2, NIST SP 800-224 |
| GMAC      | `krypto/gmac` | KS X ISO/IEC 9797-3, NIST SP 800-38D |

- `crypto/hmac` works with the hash functions of `go-krypto` too; `krypto/hmac` additionally implements `Zeroize()` and truncated tags (`hmac.NewWithSize`).

- GMAC requires an explicit, nonempty IV that is unique for every distinct message under a key. `NewGMAC(block, nil)` now returns `gmac.ErrInvalidIV`. Create a new GMAC with a new IV for each message; `Reset()` retains the IV and is only suitable for recomputing the same message. Do not publish tags for different prefixes using one GMAC instance. Use HMAC or CMAC when a reusable MAC without nonce management is needed.

- Serialized GMAC states contain secret key-derived material and must be protected accordingly. The state format is now versioned and identifies the active GHASH implementation. Legacy states and states from a different implementation (for example, assembly versus `noasm`, or a different CPU architecture) are rejected instead of silently producing incorrect tags. Recreate the GMAC from the original key, IV, and message when migrating such a state.

### Random Number Generator Supports

| Algorithm | Package       | Reference                            |
|:---------:|---------------|--------------------------------------|
| Hash_DRBG | `krypto/drbg`  | TTAK.KO-12.0331, NIST SP 800-90A    |
| HMAC_DRBG | `krypto/drbg`  | TTAK.KO-12.0332, NIST SP 800-90A    |
| CTR_DRBG  | `krypto/drbg`  | TTAK.KO-12.0189/R1, NIST SP 800-90A |

- Hash/HMAC DRBG and CTR DRBG with a derivation function automatically generate a fresh nonce when `WithNonce` is omitted. For reproducible initialization, supply both the entropy input and a nonce meeting the security-strength requirements. An explicit empty or short nonce is accepted only when additional entropy makes up the required combined input; CTR without a derivation function does not require a nonce.

- `drbg.NewEntropy` can be shared concurrently. Each successful `Get()` returns an independent caller-owned buffer; clear that buffer after use. DRBGs clear the buffers they acquire internally. Custom entropy sources retain ownership of their buffers, which DRBGs copy before use; shared custom sources must support concurrent calls and keep returned bytes stable while copied.

### Key Derivation Function Supports

| Algorithm     | Package         | Reference                                           |
|:-------------:|-----------------|-----------------------------------------------------|
| KBKDF (CMAC)  | `krypto/kbkdf`  | TTAK.KO-12.0272, NIST SP 800-108                    |
| KBKDF (HMAC)  | `krypto/kbkdf`  | TTAK.KO-12.0333, NIST SP 800-108                    |
| PBKDF2 (HMAC) | `krypto/pbkdf2` | TTAK.KO-12.0334, NIST SP 800-132, RFC 2898(PKCS #5) |

### Zeroize (Secret Memory Erasure)

- Most values that hold key material — block ciphers, cipher modes, hashes, MACs, DRBGs — implement a `Zeroize()` method.
- Call `Zeroize()` when you are finished using such a value to erase its key material (round-key schedule, IV, accumulator, key buffer, ...) from memory.
- After `Zeroize()`, the value is **marked zeroized**: further use returns `krypto.ErrStateCleared` (for `krypto.DRBG`) or panics (for block ciphers and cipher modes).

```go
c, _ := lea.NewCipher(key)
defer krypto.Zeroize(c) // calls c.Zeroize(): erases the round-key schedule
// ... use c ...
```

- The DRBGs of `krypto/drbg` also have a `Close() error` method (not part of the `krypto.DRBG` interface). It zeroizes the DRBG like `Zeroize()`, and returns `krypto.ErrStateCleared` if the DRBG was already zeroized.
- The DRBGs of `krypto/drbg` are safe for concurrent use. A DRBG zeroizes only the entropy source it created itself; zeroize a source passed with `drbg.WithCustomEntropy` yourself.
- `krypto/kipher` modes wrap a block cipher; `Zeroize()` erases the **wrapped** block cipher's key schedule, not a separately held secret key.
- The convenience `krypto.Zeroize(vals ...any)` calls `Zeroize()` on values that implement `krypto.Zeroizable`, and clears `*big.Int` values and numeric **slices**. An array passed by value is a copy and **cannot be erased** this way; use `krypto.ZeroizeT(&array)` instead.

### Key Parsing and Validation

- `kx509.ParsePKIXPublicKey` and `ParsePKCS8PrivateKey` handle EC-KCDSA and KCDSA algorithm OIDs and delegate other key types to `crypto/x509`. Standard `id-ecPublicKey` keys are parsed as ECDSA; interpreting that OID as EC-KCDSA is not supported.

- KCDSA key import and export validate prime domain parameters and subgroup membership. For a KCDSA public key constructed directly from external values, call `pub.Validate()` once before verification and again after changing any field. `Verify` assumes its domain has been validated and avoids repeating expensive primality tests on each signature.

### Timing and Memory Limitations

- ARIA and SEED implementations use secret-dependent table lookups, and signature operations use `math/big`. These implementations do not provide a constant-time guarantee. Applications with a cache or timing side-channel threat model require a separately reviewed implementation.

- `Zeroize` clears the live state owned by the value. It does not erase previously serialized states, caller-held copies, or all copies that the Go runtime or compiler may have made.

### TODO

- Support Post-Quantum Cryptography
