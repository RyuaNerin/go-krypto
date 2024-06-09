# subtle

```txt
_________  go1.20                                                       xor_go120.go
     \___ !go1.20 ______  purego _____ !go1.17                          xor.go + xor_generic.go
                    \              \__  go1.17                          xor.go + xor_generic_go117.go
                     \__ !purego _____ !go1.18                          xor.go + xor_generic_go117.go
                                   \__  go1.18 _____ !amd64 && !arm64   xor.go + xor_generic_go117.go
                                                 \__  amd64 ||  arm64   xor.go + *.s
```

```
>=go1.20      xor_go120
 <go1.20      xor.go
  go1.18      
  
```
