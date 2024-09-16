# subtle

```txt
.
├── !gc    ---------------------------------------> xor_go115.go                !gc
└──  gc
    ├──  go1.20
    │   ├──  purego    ---------------------------> xor_go120.go                gc && go1.20
    │   └── !purego
    │       ├── !loong64    ----------------------> xor_go120.go                gc && go1.20 && !purego && !loong64
    │       └──  loong64    ----------------------> xor_go118_loong64.go        gc && go1.20 && !purego && loong64
    └── !go1.20
        ├── !go1.18
        │   ├── !go1.17    -----------------------> xor_go115.go                gc && !go1.17
        │   └──  go1.17    -----------------------> xor_go117.go                gc && !go1.18 && go1.17
        └── go1.18
            ├──  purego    -----------------------> xor_go117.go                gc && !go1.20 && go1.18 && purego
            └── !purego
                ├── noasm    ---------------------> xor_go117.go                gc && !go1.20 && go1.18 && !purego && !amd64 && !arm64 && !ppc64 && !ppc64le && !loong64
                └── asm    -----------------------> xor_go118.go                gc && !go1.20 && go1.18 && !purego && (amd64 || arm64 || ppc64 || ppc64le)
                    ├── amd64    -----------------> xor_go118_amd64.go          gc && !go1.20 && go1.18 && !purego && amd64
                    ├── arm64    -----------------> xor_go118_arm64.go          gc && !go1.20 && go1.18 && !purego && arm64
                    ├── loong64    ---------------> xor_go118_loong64.go        gc && !go1.20 && go1.18 && !purego && loong64
                    └── ppc64 || ppc64le    ------> xor_go118_ppc64x.go         gc && !go1.20 && go1.18 && !purego && (ppc64 || ppc64le)
```

- support `.s` file in `>=go1.18`
- not supported `.s` file in `!gc`

```
The GCC 10 releases include a complete implementation of the Go 1.14.6 release.
The GCC 11 releases include a complete implementation of the Go 1.16.3 release.
The GCC 12 and 13 releases include a complete implementation of the Go 1.18 standard library. However, GCC does not yet include support for generics.
```

- loong64 : [69827b5c8ddf93be65bfc8b17d331bb3ff7b704c](https://github.com/golang/go/commit/69827b5c8ddf93be65bfc8b17d331bb3ff7b704c)
