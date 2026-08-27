package ctrdrbg

import (
	"crypto/subtle"
	"encoding/binary"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

const (
	MaxLength                      uint64 = 1 << 35 >> 3 // 엔트로피 입력 최대 길이 (max_length)
	MaxPersonalizationStringLength uint64 = 1 << 35 >> 3 // 개별화 문자열 최대 허용 길이 (max_personalization_string_length) = 2 ** 35
	MaxAdditionalInputLength       uint64 = 1 << 35 >> 3 // 추가 입력 최대 허용 길이 (max_additional_input_length) = 2 ** 35
	MaxDFInputLength               uint64 = 1<<32 - 1    // Maximum byte length representable by L.

	reseedInterval16Blocks uint64 = 1 << 48
	reseedIntervalEtc      uint64 = 1 << 32
)

var cbcKey = []byte{
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
	0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, // 128
	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, // 192
	0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, // 256
}

const (
	MacKeySize   = 256 / 8
	MaxBlockSize = 16
)

const maxSeedLen = MaxBlockSize + MacKeySize

type State struct {
	New                     func(key []byte) (krypto.Block8, error) // const
	KeySize                 int                                     // const
	ReseedInterval          uint64                                  // const
	UseDerivationFunction   bool                                    // const
	UsePredictionResistance bool                                    // const
	CounterSize             int                                     // const

	BlockSize              int           // const, instantiated
	SeedLen                int           // const, instantiated
	MaxNoOfBytesPerRequest int           // const, instantiated
	CBC                    krypto.Block8 // const, instantiated

	Key           [MacKeySize]byte
	V             [MaxBlockSize]byte
	ReseedCounter uint64

	tmpBlocks [8 * MaxBlockSize]byte
	tmpDFSeed [maxSeedLen]byte
}

var _ krypto.Zeroizable = (*State)(nil)

func (state *State) Zeroize() {
	state.New = nil
	state.KeySize = 0
	state.ReseedInterval = 0
	state.UseDerivationFunction = false
	state.UsePredictionResistance = false
	state.CounterSize = 0

	state.BlockSize = 0
	state.SeedLen = 0
	state.MaxNoOfBytesPerRequest = 0
	zeroize.Zeroize(state.CBC)
	state.CBC = nil

	clear(state.Key[:])
	clear(state.V[:])
	state.ReseedCounter = 0
	clear(state.tmpBlocks[:])
	clear(state.tmpDFSeed[:])
}

// 6.2 유도 함수(derivation function)
func (state *State) df(dst []byte, providedDataList ...[]byte) error {
	bs := state.BlockSize
	seedlen := state.SeedLen

	var providedDataLength uint64
	for _, v := range providedDataList {
		var err error
		providedDataLength, err = addDFInputLength(providedDataLength, uint64(len(v)))
		if err != nil {
			return err
		}
	}

	// a) L ← len(provided_data) / 8 (L은 32 비트 정수로 표현)
	// b) N ← len(seedlen ) / 8 (N은 32 비트 정수로 표현)
	var LN [4 + 4]byte
	defer clear(LN[:])
	binary.BigEndian.PutUint32(LN[0:4], uint32(providedDataLength))
	binary.BigEndian.PutUint32(LN[4:8], uint32(state.SeedLen))
	// c) S ← L ∥ N ∥ provided_data ∥ 0x80
	// d) S의 길이(len(S))가 blocklen 의 배수가 되도록 0으로 오른쪽 패딩

	// e) cnt ← 0 (cnt는 32 비트 정수로 표현. 즉 len(cnt) = 32)
	var cnt [4]byte
	defer clear(cnt[:])

	// f) temp ← Null
	temp := state.tmpDFSeed[:state.SeedLen]
	defer clear(state.tmpDFSeed[:])

	// 7) (7.1 ~ 7.6) 의 과정을 시드 길이(len_seed) 횟수만큼 반복한다
	block := state.tmpBlocks[:bs]
	defer clear(state.tmpBlocks[:])

	// g) while (len(temp) < seedlen ) do
	for off := 0; off < seedlen; off += bs {
		// 1) chaining_value ← 0 ** blocklen
		// 2) n ← len(cnt ∥ 0 ** (blocklen - len(cnt)) ∥ S) / blocklen
		// 3) (cnt ∥ 0 ** (blocklen - len(cnt)) ∥ S)를 n개의 블록(block1 ~ blockn)으로 나눔
		// 4) for i from 1 to n do
		//     (a) input_block ← chaining_value ⊕ block_i
		//     (b) chaining_value ← Block_Cipher_Enc(CBC_Key, input_block)

		//    Block 1    | Block 2   | B 3...| Block n
		// --------------|-----------|-------|---------------
		// C | 0x00 ...  | L | N | inputData    | 0x80 || pad

		// Block 1
		// C || pad
		copy(block, cnt[:])
		clear(block[4:])
		state.CBC.Encrypt(block, block)

		// Block 2
		// L || N, ...
		inputIdx := subtle.XORBytes(block, block, LN[:])
		if inputIdx == bs {
			state.CBC.Encrypt(block, block)
			inputIdx = 0
		}

		for _, inputData := range providedDataList {
			for len(inputData) > 0 {
				n := subtle.XORBytes(block[inputIdx:], block[inputIdx:], inputData)
				inputData = inputData[n:]
				inputIdx += n

				if inputIdx == bs {
					state.CBC.Encrypt(block, block)
					inputIdx = 0
				}
			}
		}

		// Block N
		// ... || 0x80 || pad
		block[inputIdx] ^= 0x80
		state.CBC.Encrypt(block, block)

		// 5) temp ← temp ∥ chaining_value
		copy(temp[off:], block)

		// 6) cnt ← cnt + 1
		internal.IncCtr(cnt[:])
	}

	// h) K ← leftmost(temp, keylen )
	// i) X ← select(temp, keylen + 1, seedlen )
	K := temp[:state.KeySize]
	X := temp[state.KeySize:]

	// j) temp ← Null

	// k) while (len(temp) < seedlen ) do
	c, _ := state.New(K)
	defer zeroize.Zeroize(c)
	for off := 0; off < seedlen; off += bs {
		// 1) X ← Block_Cipher_Enc(K, X)
		// 2) temp ← temp ∥ X
		c.Encrypt(X, X)
		copy(dst[off:], X)
	}

	// l) return leftmost(temp, seedlen )
	return nil
}

// addDFInputLength checks before adding, so even oversized aggregate lengths
// cannot overflow or silently truncate the 32-bit byte count encoded as L.
func addDFInputLength(total, length uint64) (uint64, error) {
	if total > MaxDFInputLength || length > MaxDFInputLength-total {
		return 0, DerivationInputIsTooLongError(MaxDFInputLength)
	}
	return total + length, nil
}

// MaxReseedInterval returns the request count limit for the cipher block size.
func MaxReseedInterval(blockSize int) uint64 {
	if blockSize == 16 {
		return reseedInterval16Blocks
	}
	return reseedIntervalEtc
}

// 6.3 갱신 함수(update function)
func (state *State) Update(providedData []byte) {
	b, _ := state.New(state.Key[:state.KeySize])
	defer zeroize.Zeroize(b)

	defer clear(state.tmpBlocks[:])
	state.update(b, providedData, state.tmpBlocks[:b.BlockSize()])
}

// update는 현재 Key로 초기화된 암호와 호출자의 작업 버퍼를 사용한다.
func (state *State) update(b krypto.Block8, providedData, block []byte) {
	stateV := state.V[:state.BlockSize]
	bs := b.BlockSize()
	_ = block[:bs]

	// a) temp ← Null
	var tempBuf [maxSeedLen]byte
	defer clear(tempBuf[:])
	temp := tempBuf[:state.SeedLen]

	// b) while (len(temp) < seedlen ) do
	for off := 0; off < state.SeedLen; off += bs {
		// 1) if (ctr_len < blocklen ) then
		//     (a) inc ← (rightmost(V, ctr_len ) + 1) mod 2 ** ctr_len
		//     (b) V ← leftmost(V, blocklen - ctr_len) ∥ inc
		// 2) else V ← (V + 1) mod 2 ** blocklen
		internal.IncCtr(stateV[bs-state.CounterSize:])

		// 3) output ← Block_Cipher_Enc(Key, V )
		b.Encrypt(block, stateV)
		// 4) temp ← temp ∥ output
		copy(temp[off:], block)
	}

	// c) temp ← leftmost(temp, seedlen )
	// d) temp ← temp ⊕ provided_data
	subtle.XORBytes(temp[:], temp[:], providedData)
	// e) Key ← leftmost(temp, keylen )
	copy(state.Key[:state.KeySize], temp[:state.KeySize])
	// f) V ← rightmost(temp, blocklen )
	copy(stateV, temp[state.KeySize:])
	// g) return (V, Key )
}

// 10. 블록 암호 알고리즘 기반 난수 발생기의 초기화 함수
func (state *State) Instantiate(
	newCipher func(key []byte) (krypto.Block8, error),
	keySize int,
	reseedInterval uint64,
	ctrLength int,
	entropyInput, nonce, personalizationString []byte,
	useDerivationFunction, usePredictionResistance bool,
) error {
	cbc, err := newCipher(cbcKey[:keySize])
	if err != nil {
		return err
	}

	if reseedInterval <= 0 {
		reseedInterval = MaxReseedInterval(cbc.BlockSize())
	}
	if ctrLength <= 0 {
		ctrLength = cbc.BlockSize()
	}

	blockSize := cbc.BlockSize()
	seedSize := keySize + blockSize

	maxNoOfBytesPerRequest := maxBytesPerRequest(blockSize, ctrLength)

	*state = State{
		New:                     newCipher,
		KeySize:                 keySize,
		ReseedInterval:          reseedInterval,
		UseDerivationFunction:   useDerivationFunction,
		UsePredictionResistance: usePredictionResistance,
		CounterSize:             ctrLength,

		BlockSize:              blockSize,
		SeedLen:                seedSize,
		MaxNoOfBytesPerRequest: maxNoOfBytesPerRequest,
		CBC:                    cbc,
	}

	var tmpReseedBuf [maxSeedLen]byte
	defer clear(tmpReseedBuf[:])
	tmpReseed := tmpReseedBuf[:state.SeedLen]

	// a) 유도 함수 사용: 단계 a.1)과 단계 a.2) 수행
	// b) 유도 함수 미사용: 단계 b.1)과 단계 b.2) 수행
	if state.UseDerivationFunction {
		// 1) seed_material ← entropy_input∥nonce∥personalization_string
		// 2) seed ← Block_Cipher_df(seed_material)
		if err := state.df(tmpReseed, entropyInput, nonce, personalizationString); err != nil {
			state.Zeroize()
			return err
		}
	} else {
		// 1) 개별화 문자열(personalization_string)의 길이가 seedlen 보다 작으면, seedlen 길이가 되도록 0으로 오른쪽 패딩
		// 2) seed ← entropy_input ⊕ personalization_string
		copy(tmpReseed, entropyInput)
		subtle.XORBytes(tmpReseed, tmpReseed, personalizationString)
	}

	// c) Key ← 0 ** keylen
	// d) V ← 0 ** blocklen

	// e) (V, Key ) ← CTR_DRBG_Update((V, Key ), seed)
	state.Update(tmpReseed)

	// f) reseed_counter ← 1
	state.ReseedCounter = 1

	// g) return (V, Key, reseed_counter )
	return nil
}

// 6.5 재초기화 함수(reseed function)
func (state *State) Reseed(entropyInput, additionalInput []byte) error {
	var tmpReseedBuf [maxSeedLen]byte
	defer clear(tmpReseedBuf[:])
	tmpReseed := tmpReseedBuf[:state.SeedLen]

	// a) 유도 함수 사용: 단계 a.1)과 단계 a.2)를 수행
	// b) 유도 함수 미사용: 단계 b.1)과 단계 b.2)를 수행
	if state.UseDerivationFunction {
		// 1) seed_material ← entropy_input ∥ additional_input
		// 2) seed ← Block_Cipher_df(seed_material)
		if err := state.df(tmpReseed, entropyInput, additionalInput); err != nil {
			return err
		}
	} else {
		// 1) 추가 입력(additional_input)의 길이가 seedlen 보다 작으면, seedlen 길이가 되도록 0으로 오른쪽 패딩
		// 2) seed ← entropy_input ⊕ additional_input
		copy(tmpReseed, entropyInput)
		subtle.XORBytes(tmpReseed, tmpReseed, additionalInput)
	}

	// c) (V, Key ) ← CTR_DRBG_Update((V, Key ), seed)
	state.Update(tmpReseed)

	// d) reseed_counter ← 1
	state.ReseedCounter = 1

	// e) return (V, Key, reseed_counter )
	return nil
}

// 6.6 출력 생성 함수(generate function)
func (state *State) Generate(
	dst []byte,
	fnEntropyInput func() ([]byte, error),
	additionalInput []byte,
) error {
	var tmpReseedBuf [maxSeedLen]byte
	defer clear(tmpReseedBuf[:])
	tmpReseed := tmpReseedBuf[:state.SeedLen]

	bs := state.BlockSize

	// a) 예측 내성이 활성화되어 있거나 출력 생성 횟수가 상태 갱신 주기보다 크면, 단계 a.1)과 단계 a.2)를 수행
	if state.UsePredictionResistance || state.ReseedCounter > state.ReseedInterval {
		entropyInput, err := fnEntropyInput()
		if err != nil {
			return err
		}

		// 1) (V, Key ) ← Reseed_CTR_DRBG((V, Key ), entropy_input, additional_input)
		err = state.Reseed(entropyInput, additionalInput)
		clear(entropyInput)
		if err != nil {
			return err
		}
		// 2) additional_input ← Null
		additionalInput = nil
	}

	// b) 추가 입력(additional_input)이 Null 이 아니면 단계 b.1)부터 단계 b.3)까지를 수행
	if len(additionalInput) > 0 {
		// 1) 유도 함수를 사용하는 경우 단계 (a)를 수행한다.
		//     (a) additional_input ← Block_Cipher_df(additional_input)
		// 2) 유도 함수를 사용하지 않고 추가 입력(additional_input)의 길이가 seedlen보다 작으면, seedlen 길이가 되도록 0으로 오른쪽 패딩
		if state.UseDerivationFunction {
			if err := state.df(tmpReseed, additionalInput); err != nil {
				return err
			}
			additionalInput = tmpReseed
		} /* else if len(additionalInput) < bs*state.seedLenByte {
			additionalInput = internal.FitSize(additionalInput, bs*state.seedLenByte)
		}*/
		// CTR_DRBG_Update에서 providedData == nil 인 경우에도 성립.
		//     subtle.XORBytes(temp, temp, nil) 하게되면 아무 작업 안되는데,
		//     xor 0 일 때도 마찬가지로 아무 작업 안됨.

		// 3) (V, Key ) ← CTR_DRBG_Update((V, Key ), additional_input)
		state.Update(additionalInput)
	} /** else {
		// c) 추가 입력이 Null 이면, additional_input ← 0 ** seedlen
		//additional_input = zeros[:state.SeedLenByte]
		// CTR_DRBG_Update에서 providedData == nil 인 경우에도 성립.
		//     subtle.XORBytes(temp, temp, nil) 하게되면 아무 작업 안되는데,
		//     xor 0 일 때도 마찬가지로 아무 작업 안됨.
	}*/

	// d) temp ← Null
	// e) num ← len_output
	// f) n ← ⎡num / blocklen⎤
	// g) for i from 1 to n do
	blocks := state.tmpBlocks[:8*bs]
	defer clear(state.tmpBlocks[:])
	b, _ := state.New(state.Key[:state.KeySize])
	defer zeroize.Zeroize(b)

	stateV := state.V[:state.BlockSize]
	for len(dst) > 0 {
		outputLen := min(len(dst), len(blocks))
		blocksToGenerate := internal.CeilDiv(outputLen, bs)

		// 1) if (ctr_len < blocklen ) then
		//     (a) inc ← (rightmost(V, ctr_len ) + 1) mod 2 ** ctr_len
		//     (b) V ← leftmost(V, blocklen - ctr_len) ∥ inc
		// 2) else V ← (V + 1) mod 2 ** blocklen
		if state.CounterSize == bs {
			internal.FillCtr(blocks[:blocksToGenerate*bs], stateV)
		} else {
			for off := 0; off < blocksToGenerate*bs; off += bs {
				internal.IncCtr(stateV[bs-state.CounterSize:])
				copy(blocks[off:], stateV)
			}
		}

		// 3) temp ← temp ∥ Block_Cipher_Enc(Key, V )
		// Keep partial batches on the available parallel encryption paths.
		switch {
		case blocksToGenerate == 1:
			b.Encrypt(blocks[:bs], blocks[:bs])
		case blocksToGenerate <= 4:
			b.Encrypt4(blocks, blocks)
		default:
			b.Encrypt8(blocks, blocks)
		}
		dst = dst[copy(dst, blocks[:outputLen]):]
	}

	// h) returned_bits ← leftmost(temp, num)
	// i) (V, Key ) ← CTR_DRBG_Update((V, Key ), additional_input)
	// 출력 생성 중 Key가 바뀌지 않으므로 암호와 작업 버퍼를 재사용한다.
	state.update(b, additionalInput, blocks[:bs])

	// j) reseed_counter ← reseed_counter + 1
	state.ReseedCounter++

	// k) return (returned_bits, V, Key, reseed_counter )
	return nil
}

func maxBytesPerRequest(blockSize, counterSize int) int {
	limit := uint64((1 << 19) / 8)
	if blockSize == 8 {
		limit = (1 << 13) / 8
	}
	if counterSize >= 8 {
		return int(limit)
	}
	blocks := (uint64(1) << (uint(counterSize) * 8)) - 4
	return int(min(limit, blocks*uint64(blockSize)))
}
