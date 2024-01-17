# SIMD

## go builds

| file                                 | etc  | `amd64` (SSE2) | `amd64` (SSSE3) | arm64 | purego | `//go:build`                     |
|--------------------------------------|:----:|:--------------:|:---------------:|:-----:|:------:|----------------------------------|
| `func newCipher`                     | `Go` | `Go`           | `Asm`           | `Asm` | `Go`   |                                  |
| [aria_amd64.go](./aria_amd64.go)     |      | O              | O               |       |        | `amd64 && !purego`               |
| [aria_arm64.go](./aria_arm64.go)     |      |                |                 | O     |        | `arm64 && !purego`               |
| [aria_asm.go](./aria_asm.go)         |      | O              | O               | O     |        | `(amd64 \|\| arm64) && !purego`  |
| [aria_const.go](./aria_const.go)     | O    | O              | O               | O     | O      |                                  |
| [aria_generic.go](./aria_generic.go) | O    | O              | O               |       | O      | `!arm64 \|\| purego`             |
| [aria_noasm.go](./aria_noasm.go)     | O    |                |                 |       | O      | `(!amd64 && !arm64) \|\| purego` |
