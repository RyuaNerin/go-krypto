# go-krypto Performance

## Package Performance

```txt
goos: windows
goarch: amd64
pkg: github.com/RyuaNerin/go-krypto
cpu: AMD Ryzen 7 3700X 8-Core Processor             
Benchmark_CBC_Encrypt_1K_AES-16        	 1272006	       943.5 ns/op	1085.38 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Decrypt_1K_AES-16        	 1581552	       760.7 ns/op	1346.08 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Encrypt_1K_SEED128-16    	  104672	     11524 ns/op	  88.86 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Decrypt_1K_SEED128-16    	  108507	     11103 ns/op	  92.23 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Encrypt_1K_SEED256-16    	   70472	     16891 ns/op	  60.63 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Decrypt_1K_SEED256-16    	   73366	     16405 ns/op	  62.42 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Encrypt_1K_HIGHT-16      	   54616	     21867 ns/op	  46.83 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Decrypt_1K_HIGHT-16      	   52513	     21857 ns/op	  46.85 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Encrypt_1K_ARIA_16-16    	   45429	     26165 ns/op	  39.14 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Decrypt_1K_ARIA_16-16    	   42463	     26516 ns/op	  38.62 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Encrypt_1K_ARIA_24-16    	   40006	     30579 ns/op	  33.49 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Decrypt_1K_ARIA_24-16    	   39322	     30135 ns/op	  33.98 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Encrypt_1K_ARIA_32-16    	   35148	     34179 ns/op	  29.96 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Decrypt_1K_ARIA_32-16    	   34988	     34057 ns/op	  30.07 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Encrypt_1K_LEA_16-16     	  406381	      2819 ns/op	 363.21 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Decrypt_1K_LEA_16-16     	  922296	      1342 ns/op	 762.87 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Encrypt_1K_LEA_24-16     	  379173	      3184 ns/op	 321.63 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Decrypt_1K_LEA_24-16     	  823954	      1433 ns/op	 714.73 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Encrypt_1K_LEA_32-16     	  336982	      3524 ns/op	 290.55 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Decrypt_1K_LEA_32-16     	  799540	      1481 ns/op	 691.39 MB/s	       0 B/op	       0 allocs/op
Benchmark_HASH_SHA256_1K-16            	  479318	      2321 ns/op	 441.25 MB/s	      32 B/op	       1 allocs/op
Benchmark_HASH_SHA512_1K-16            	  684993	      1632 ns/op	 627.60 MB/s	      64 B/op	       1 allocs/op
Benchmark_HASH_LSH256_1K-16            	  631651	      1872 ns/op	 546.94 MB/s	      64 B/op	       2 allocs/op
Benchmark_HASH_LSH512_1K-16            	  163092	      7295 ns/op	 140.38 MB/s	      64 B/op	       1 allocs/op
PASS
coverage: [no statements]
ok  	github.com/RyuaNerin/go-krypto	34.435s
```

## LEA Performance (purego vs SIMD)

```txt
goos: windows
goarch: amd64
pkg: github.com/RyuaNerin/go-krypto/lea
cpu: AMD Ryzen 7 3700X 8-Core Processor
Benchmark_New_LEA128-16                     	13652575	        86.86 ns/op	       0 B/op	       0 allocs/op
Benchmark_New_LEA192-16                     	16194288	        72.73 ns/op	       0 B/op	       0 allocs/op
Benchmark_New_LEA256-16                     	14659854	        81.54 ns/op	       0 B/op	       0 allocs/op

Benchmark_LEA128_Encrypt_1Block_Go-16       	31579611	        37.74 ns/op	 423.91 MB/s	       0 B/op	       0 allocs/op
Benchmark_LEA128_Encrypt_4Blocks_SSE2-16    	18358532	        57.87 ns/op	1105.96 MB/s	       0 B/op	       0 allocs/op
Benchmark_LEA128_Encrypt_8Blocks_AVX2-16    	 8417154	       145.8 ns/op	 878.12 MB/s	       0 B/op	       0 allocs/op

Benchmark_LEA128_Decrypt_1Block_Go-16       	15081079	        79.53 ns/op	 201.18 MB/s	       0 B/op	       0 allocs/op
Benchmark_LEA128_Decrypt_4Blocks_SSE2-16    	19994101	        58.99 ns/op	1084.91 MB/s	       0 B/op	       0 allocs/op
Benchmark_LEA128_Decrypt_8Blocks_AVX2-16    	 8363638	       146.8 ns/op	 872.04 MB/s	       0 B/op	       0 allocs/op

Benchmark_LEA192_Encrypt_1Block_Go-16       	27154419	        42.49 ns/op	 376.59 MB/s	       0 B/op	       0 allocs/op
Benchmark_LEA192_Encrypt_4Blocks_SSE2-16    	18286327	        65.75 ns/op	 973.38 MB/s	       0 B/op	       0 allocs/op
Benchmark_LEA192_Encrypt_8Blocks_AVX2-16    	 7780092	       158.6 ns/op	 807.21 MB/s	       0 B/op	       0 allocs/op

Benchmark_LEA192_Decrypt_1Block_Go-16       	13569962	        88.16 ns/op	 181.49 MB/s	       0 B/op	       0 allocs/op
Benchmark_LEA192_Decrypt_4Blocks_SSE2-16    	17940249	        68.24 ns/op	 937.93 MB/s	       0 B/op	       0 allocs/op
Benchmark_LEA192_Decrypt_8Blocks_AVX2-16    	 7658570	       154.8 ns/op	 826.88 MB/s	       0 B/op	       0 allocs/op

Benchmark_LEA256_Encrypt_1Block_Go-16       	24489844	        48.79 ns/op	 327.92 MB/s	       0 B/op	       0 allocs/op
Benchmark_LEA256_Encrypt_4Blocks_SSE2-16    	16670625	        81.64 ns/op	 783.92 MB/s	       0 B/op	       0 allocs/op
Benchmark_LEA256_Encrypt_8Blocks_AVX2-16    	 7241037	       164.9 ns/op	 776.01 MB/s	       0 B/op	       0 allocs/op

Benchmark_LEA256_Decrypt_1Block_Go-16       	12210720	        97.56 ns/op	 164.01 MB/s	       0 B/op	       0 allocs/op
Benchmark_LEA256_Decrypt_4Blocks_SSE2-16    	15064172	        77.76 ns/op	 823.07 MB/s	       0 B/op	       0 allocs/op
Benchmark_LEA256_Decrypt_8Blocks_AVX2-16    	 7447822	       163.1 ns/op	 784.97 MB/s	       0 B/op	       0 allocs/op

Benchmark_CBC_Decrypter_Go-16               	12569037	        93.93 ns/op	 170.35 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Decrypter_Asm_1Block-16       	12981884	        92.27 ns/op	 173.40 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Decrypter_Asm_4Blocks-16      	15765388	        75.79 ns/op	 844.46 MB/s	       0 B/op	       0 allocs/op
Benchmark_CBC_Decrypter_Asm_8Blocks-16      	 7122519	       167.1 ns/op	 765.90 MB/s	       0 B/op	       0 allocs/op

Benchmark_CTR_Go-16                         	22222550	        52.86 ns/op	 302.68 MB/s	       0 B/op	       0 allocs/op
Benchmark_CTR_Asm-16                        	36334017	        32.28 ns/op	 495.69 MB/s	       0 B/op	       0 allocs/op
PASS
coverage: 92.7% of statements
ok  	github.com/RyuaNerin/go-krypto/lea	40.245s
```
