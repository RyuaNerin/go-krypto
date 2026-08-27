package hashdrbg

import (
	"encoding/binary"
	"hash"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

const (
	MaxLength                      uint64 = 1 << 35 >> 3 // 엔트로피 입력 최대 길이 (max_length)
	MaxPersonalizationStringLength uint64 = 1 << 35 >> 3 // 개별화 문자열 최대 허용 길이 (max_personalization_string_length) = 2 ** 35
	MaxAdditionalInputLength       uint64 = 1 << 35 >> 3 // 추가 입력 최대 허용 길이 (max_additional_input_length) = 2 ** 35
	MaxNoOfBytesPerRequest                = 1 << 19 >> 3 // 난수 최대 출력 길이 (max_no_of_bits_per_request) = 2 ** 19
	ReseedInterval                 uint64 = 1 << 48      // 시드별 출력값 생성 횟수 (reseed_interval) = 2 ** 48
)

const (
	maxSeedLen  = 111
	maxHashSize = 64 // SHA-512
)

type State struct {
	Hash                        hash.Hash // const
	ReseedInterval              uint64    // const
	PredictionResistanceRequest bool      // const

	SeedLen int // const, instantiated
	SumLen  int // const, instantiated

	V             [maxSeedLen]byte // seedlen 바이트 길이의 변수로, DRBG 인스턴스의 출력값 생성 시 갱신
	C             [maxSeedLen]byte // 시드로부터 계산되는 seedlen 바이트 길이의 상수
	ReseedCounter uint64           // 시드 생성 이후 DRBG 인스턴스의 출력값 생성 횟수

	tmpV      [maxSeedLen]byte
	tmpSum    [maxHashSize]byte
	tmpHeader [5]byte
	tmpDomain [1]byte
}

var _ krypto.Zeroizable = (*State)(nil)

func (state *State) Zeroize() {
	zeroize.Zeroize(state.Hash)
	state.Hash = nil
	state.ReseedInterval = 0
	state.PredictionResistanceRequest = false

	state.SeedLen = 0
	state.SumLen = 0

	clear(state.V[:])
	clear(state.C[:])
	state.ReseedCounter = 0

	clear(state.tmpV[:])
	clear(state.tmpSum[:])
	clear(state.tmpHeader[:])
	clear(state.tmpDomain[:])
}

func GetSecurityStrengthBytes(requestedStrength int) int {
	switch {
	case requestedStrength <= 14:
		return 14
	case requestedStrength <= 16:
		return 16
	case requestedStrength <= 24:
		return 24
	default:
		return 32
	}
}

// 6.2 유도 함수(derivation function)
//
// Hash_DRBG의 유도 함수 Hash_df는 내부 상태 구성 요소를 생성하거나 입력이 가진 엔트로피를 출력에 고르게 분포되도록 한다
func (state *State) df(dst []byte, data ...[]byte) {
	h := state.Hash

	// 1: temp ← Null
	// 2: len ← no_of_bits /outlen 

	noOfBytes := len(dst)

	sumBuf := &state.tmpSum
	defer clear(sumBuf[:])
	sum := sumBuf[:0]

	hashInput := &state.tmpHeader
	defer clear(hashInput[:])
	binary.BigEndian.PutUint32(hashInput[1:], uint32(noOfBytes*8))
	// 3: for i = 1 to len do
	for idx := 0; idx < len(dst); {
		hashInput[0]++

		// 4: temp ← temp ‖ Hash(i ‖ no_of_bits ‖ input_string )
		h.Reset()
		h.Write(hashInput[:])
		for _, v := range data {
			h.Write(v)
		}
		sum = h.Sum(sum[:0])
		idx += copy(dst[idx:], sum)
	}
	// 5: end do

	// 6: requested_bits ← leftmost(temp, no_of_bits )

	// 7: return requested_bits
}

// 6.3 인스턴스 생성 함수(instantiate function)
//
// Hash_DRBG의 인스턴스 생성 함수 Instantiate_Hash_DRBG는
// 엔트로피 입력, 논스, 개별화 문자열로부터 시드를 생성하고, 이 시드를 이용하여 내부 상태를 초기화한다.
func (state *State) Instantiate(
	h hash.Hash,
	entropyInput, nonce, personalizationString []byte,
	predictionResistanceRequest bool,
	reseedInterval uint64,
) {
	seedLen := 55
	if h.Size() > 32 {
		seedLen = 111
	}

	if reseedInterval <= 0 {
		reseedInterval = ReseedInterval
	}

	//  1: if (requested_strength > highest_supported_security_strength) then
	//  2:     return (ERROR_FLAG, Invalid) // ex. ("Invalid requested_security_strength", -1)
	//  3: end if

	//  4: if ((prediction_resistance_flag is set) AND
	//   (prediction resistance is not supported) then
	//  5:     return (ERROR_FLAG, Invalid)
	//  6: end if

	//  7: if (len(personalization_string ) > max_personalization_string_length) then
	//  8:     return (ERROR_FLAG, Invalid) // ex. ("Personalization_string too long", -1)
	//  9: end if

	// 10: if (requested_strength ≤ 112) then
	// 11:     security_strength ← 112
	// 12: else if (requested_strength ≤ 128) then
	// 13:     security_strength ← 128
	// 14: else if (requested_strength ≤ 192) then
	// 15:     security_strength ← 192
	// 16: else
	// 17:     security_strength ← 256
	// 18: end if

	// 19: (status, entropy_input) ← Get_entropy(min_entropy, min_len, max_len, request )
	//  // min_len = min_length, max_len = max_length, request = prediction_resistance_request
	// 20: if (status ≠ SUCCESS) then
	// 21:     return (status, Invalid)
	// 22: end if

	// 23: obtain a nonce and check its acceptability

	// 24: seed_material ← entropy_input ‖ nonce ‖ personalization_string
	// seedMaterial := drbg.WriteBytes(entropyInput, nonce, personalizationString)

	*state = State{
		Hash:                        h,
		ReseedInterval:              reseedInterval,
		PredictionResistanceRequest: predictionResistanceRequest,

		SeedLen: seedLen,
		SumLen:  h.Size(),

		ReseedCounter: 1,
	}

	// 25: seed ← Hash_df(seed_material, seedlen )
	state.df(state.V[:state.SeedLen], entropyInput, nonce, personalizationString)

	// 26: (status, state_handle) ← Find_state_space()
	// 27: if (status ≠ SUCCESS) then
	// 28:     return (status, Invalid)
	// 29: end if
	// 30: state(state_handle ).V ← seed
	// 31: state(state_handle ).C ← Hash_df((0x00 ‖ V ), seedlen )
	// 32: state(state_handle ).reseed_counter ← 1
	// 33: state(state_handle ).security_strength ← security_strength
	// 34: state(state_handle ).prediction_resistance_flag ← prediction_resistance_flag
	// 35: return (SUCCESS, state_handle)

	state.df(state.C[:state.SeedLen], state.tmpDomain[:], state.V[:state.SeedLen])
}

// 6.4 리시드 함수(reseeding function)
//
// Hash_DRBG의 리시드 함수 Reseed_Hash_DRBG는 인스턴스의 내부 상태와 엔트로피 입력,
// 그리고 추가 입력을 이용하여 새로운 시드를 생성하고,
// 이 시드를 이용하여 내부 상태를 초기화한다.
func (state *State) Reseed(entropyInput, additionalInput []byte) {
	//  1: if (a state(state_handle ) is not available) then
	//  2:     return ERROR_FLAG // ex. ‶State not available for the state_handle″
	//  3: end if

	//  4: V ← state(state_handle ).V
	//  5: security_strength ← state(state_handle ).security_strength
	//  6: (status, entropy_input ) ← Get_entropy(min_entropy, min_len, max_len, request )
	//  // min_len = min_length, max_len = max_length, request = prediction_resistance_request	}

	//  7: if (status ≠ SUCCESS) then
	//  8:     return (status )
	//  9: end if

	tmp := state.tmpV[:state.SeedLen]
	defer clear(state.tmpV[:])
	defer clear(state.tmpDomain[:])

	// 10: seed_material ← 0x01 ‖ V ‖ entropy_input ‖ additional_input
	copy(tmp, state.V[:state.SeedLen])
	// seed_material := drbg.WriteBytes([]byte{1}, state.tmpV, entropyInput, additionalInput)
	// 11: seed ← Hash_df(seed_material, seedlen )
	// 12: state(state_handle ).V ← seed
	state.tmpDomain[0] = 1
	state.df(state.V[:state.SeedLen], state.tmpDomain[:], tmp, entropyInput, additionalInput)
	// 13: state(state_handle ).C ← Hash_df((0x00 ‖ V ), seedlen )
	state.tmpDomain[0] = 0
	state.df(state.C[:state.SeedLen], state.tmpDomain[:], state.V[:state.SeedLen])
	// 14: state(state_handle ).reseed_counter ← 1
	state.ReseedCounter = 1
	// 15: return (SUCCESS)
}

// 6.5 생성 함수(generate function)
//
// Hash_DRBG의 생성 함수 Generate_Hash_DRBG는
// 인스턴스의 내부 상태 중 동작 상태의 V 를 이용하여 출력 비트열을 생성하고,
// 동작 상태의 C 와 reseed_counter 를 이용하여 V 를 갱신한다.
func (state *State) Generate(
	dst []byte,
	fnEntropyInput func() ([]byte, error),
	additionalInput []byte,
) error {
	//  1: if (a state(state_handle ) is not available) then
	//  2:     return (ERROR_FLAG, Null )
	//  3: end if

	//  4: V ← state(state_handle).V
	//  5: C ← state(state_handle).C
	//  6: reseed_counter ← state(state_handle).reseed_counter
	//  7: security_strength ← state(state_handle).security_strength
	//  8: prediction_resistance_flag ← state(state_handle).prediction_resistance_flag

	//  9: if (requested_no_of_bits > max_no_of_bits_per_request) then
	// 10:     return (ERROR_FLAG, Null ) // ex. (‶Too many bits requested″, Null)
	// 11: end if

	// 12: if (requested_strength > security_strength) then
	// 13:     return (ERROR_FLAG, Null ) // ex. (‶Invalid requested_security_strength″, Null)
	// 14: end if

	// 15: if (len(additional_input) > max_additional_input_length) then
	// 16:     return (ERROR_FLAG, Null ) // ex. (‶additional_input too long″, Null)
	// 17: end if

	// 18: if ((prediction_resistance_request is set) AND (prediction_resistance_flag is not set)) then
	// 19: return (ERROR_FLAG, Null )
	// 20: end if

	// 21: if ((reseed_counter > reseed_interval) OR (prediction_resistance_request is set)) then
	if state.ReseedCounter > state.ReseedInterval || state.PredictionResistanceRequest {
		// 22: if reseeding is not available then
		// 23:     return (ERROR_FLAG, Null)
		// 24: end if

		// 25: status ← Reseed_Hash_DRBG(handle, request, additional_input)
		//                 //handle = state_handle, request = prediction_resistance_request
		// 26: if (status ≠ SUCCESS) then
		// 27:     return (status, Null)
		// 28: end if
		entropyInput, err := fnEntropyInput()
		if err != nil {
			return err
		}
		state.Reseed(entropyInput, additionalInput)
		clear(entropyInput)

		// 29: V ← state(state_handle ).V
		// 30: C ← state(state_handle ).C
		// 31: reseed_counter ← state(state_handle ).reseed_counter
		// 32: additional_input ← Null
		additionalInput = nil
	}
	// 33: end if

	stateV := state.V[:state.SeedLen]

	sumBuf := &state.tmpSum
	defer clear(sumBuf[:])
	defer clear(state.tmpDomain[:])
	sum := sumBuf[:state.SumLen][:0]

	// 34: if (additional_input ≠ Null) then
	if len(additionalInput) > 0 {
		// 35: w ← Hash(0x02 ‖ V ‖ additional_input)
		state.Hash.Reset()
		state.tmpDomain[0] = 2
		state.Hash.Write(state.tmpDomain[:])
		state.Hash.Write(stateV)
		state.Hash.Write(additionalInput)
		w := state.Hash.Sum(sum)
		// 36: V ← (V + w ) mod (2 ** seedlen)
		internal.Add(stateV, stateV, w)
	}
	// 37: end if

	// 38: pseudorandom_bits ← Hashgen(requested_no_of_bits, V )
	state.hashgen(dst, stateV)

	// 39: H ← Hash(0x03 ‖ V )
	state.Hash.Reset()
	state.tmpDomain[0] = 0x03
	state.Hash.Write(state.tmpDomain[:])
	state.Hash.Write(stateV)
	H := state.Hash.Sum(sum)
	// 40: V ← (V + H + C + reseedCounter) mod (2 ** seedlen)
	// 41: state(state_handle).V ← V
	var reseedCounter [8]byte
	binary.BigEndian.PutUint64(reseedCounter[:], state.ReseedCounter)

	internal.Add(stateV, stateV, H, state.C[:state.SeedLen], reseedCounter[:])
	// 42: state(state_handle).reseed_counter ← reseed_counter + 1
	state.ReseedCounter++

	// 43: return (SUCCESS, pseudorandom_bits)
	return nil
}

// 알고리즘 4의 출력값 생성 과정(단계 39)에서 사용되는 함수 Hashgen의 구체적인 동작 방식은 알고리즘 5와 같다.
func (state *State) hashgen(dst, V []byte) {
	tmp := state.tmpV[:state.SeedLen]
	defer clear(state.tmpV[:])

	// Generate clears tmpSum after updating V.
	sum := state.tmpSum[:0]
	h := state.Hash
	sumLen := state.SumLen

	// W := dst
	copy(tmp, V)

	for len(dst) > 0 {
		h.Reset()
		h.Write(tmp)
		if len(dst) >= sumLen {
			// Append complete digests directly to the output.
			copy(dst, h.Sum(dst[:0:sumLen]))
			dst = dst[sumLen:]
		} else {
			copy(dst, h.Sum(sum))
			break
		}

		if len(dst) > 0 {
			// data ← (data + 1) mod 2 ** seedlen
			internal.IncCtr(tmp)
		}
	}
}
