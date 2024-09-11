# subtle

```txt
!gc                                                                         xor_go115.go            !gc
gc _____  go1.20                                                            xor_go120.go            gc && go1.20
     \__ !go1.20 ______  purego _____ !go1.17                               xor_go115.go            gc && !go1.17 && purego
                   \              \__  go1.17                               xor_go117.go            gc && !go1.20 && purego && go1.17
                    \__ !purego _____ !go1.18                               xor_go117.go            gc && !go1.18 && !purego
                                  \__  go1.18 _____ amd64                   xor_go118_amd64.go      gc && !go1.20 && !purego && go1.18 && amd64
                                                |__ arm64                   xor_go118_arm64.go      gc && !go1.20 && !purego && go1.18 && arm64
                                                |__ ppc64 || ppc64le        xor_go118_ppc64x.go     gc && !go1.20 && !purego && go1.18 && (ppc64 || ppc64le)
                                                |__ all not                 xor_go117.go            gc && !go1.20 && !purego && go1.18 && !arm64 && !amd64 && !ppc64 && !ppc64le
```

- support `.s` file in `>=go1.18`
- not supported `.s` file in `!gc`

```
The GCC 10 releases include a complete implementation of the Go 1.14.6 release.
The GCC 11 releases include a complete implementation of the Go 1.16.3 release.
The GCC 12 and 13 releases include a complete implementation of the Go 1.18 standard library. However, GCC does not yet include support for generics.
```
