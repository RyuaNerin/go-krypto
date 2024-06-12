package hmacdrbg

import (
	"crypto/hmac"
	"hash"

	"github.com/RyuaNerin/go-krypto/internal/memory"
)

const (
	MaxLength                      uint64 = 0x8_0000_0000      // 엔트로피 입력 최대 길이 (max_length)
	MaxPersonalizationStringLength uint64 = 0x8_0000_0000      // 개별화 문자열 최대 허용 길이 (max_personalization_string_length) = 2 ** 35
	MaxAdditionalInputLength       uint64 = 0x8_0000_0000      // 추가 입력 최대 허용 길이 (max_additional_input_length) = 2 ** 35
	MaxNoOfBytesPerRequest                = 0x8_0000 / 8       // 난수 최대 출력 길이 (max_no_of_bits_per_request) = 2 ** 19
	ReseedInterval                 uint64 = 0x8_0000_0000_0000 // 시드별 출력값 생성 횟수 (reseed_interval) = 2 ** 48
)

type State struct {
	New                         func() hash.Hash // const
	ReseedInterval              uint64           // const
	PredictionResistanceRequest bool             // const

	V             []byte
	Key           []byte
	ReseedCounter uint64 // 시드 생성 이후 DRBG 인스턴스의 출력값 생성 횟수

	sum []byte
}

func GetSecurityStrengthBits(requested_strength int) int {
	switch {
	case requested_strength <= 112:
		return 112
	case requested_strength <= 128:
		return 128
	case requested_strength <= 192:
		return 192
	default:
		return 256
	}
}

// 6.2 갱신 함수(update function)
//
// 갱신 함수는 입력받은 데이터를 이용하여 HMAC_DRBG의 내부 상태를 갱신한다.
// 참고로 갱신 함수는 인스턴스 생성 함수와 리시드 함수에서
// 시드 생성을 위한 유도 함수(derivation function)의 역할도 수행한다.
func (state *State) Update(providedData ...[]byte) {
	// 1: Key ← HMAC(Key, V ‖ 0x00 ‖ provided_data )
	h := hmac.New(state.New, state.Key)
	h.Write(state.V)
	h.Write([]byte{0x00})
	for _, v := range providedData {
		h.Write(v)
	}
	copy(state.Key, h.Sum(state.sum[:0]))

	// 2: V ← HMAC(Key, V )
	h = hmac.New(state.New, state.Key)
	h.Write(state.V)
	copy(state.V, h.Sum(state.sum[:0]))

	// 3: if (provided_data = Null) then
	// 4:     return (Key, V )
	// 5: end if
	if len(providedData) == 0 {
		return
	}

	// 6: Key ← HMAC(Key, V ‖ 0x01 ‖ provided_data )
	h = hmac.New(state.New, state.Key)
	h.Write(state.V)
	h.Write([]byte{0x01})
	for _, v := range providedData {
		h.Write(v)
	}
	copy(state.Key, h.Sum(state.sum[:0]))

	// 7: V ← HMAC(Key, V )
	h = hmac.New(state.New, state.Key)
	h.Write(state.V)
	copy(state.V, h.Sum(state.sum[:0]))

	// 8: return (Key, V )
}

// 6.3 인스턴스 생성 함수(instantiate function)
//
// HMAC_DRBG의 인스턴스 생성 함수 Instantiate_HMAC_DRBG는
// 엔트로피 입력, 논스, 개별화 문자열로부터 시드 seed = (Key, V )를 생성하고,
// 이 시드를 이용하여 내부 상태를 초기화한다.
func Instantiate(
	h func() hash.Hash,
	entropyInput, nonce, personalizationString []byte,
	predictionResistanceRequest bool,
	reseedInterval uint64,
) *State {
	outlen := h().Size()

	if reseedInterval <= 0 {
		reseedInterval = ReseedInterval
	}

	//  1: if (requested_strength > highest_supported_security_strength) then
	//  2:     return (ERROR_FLAG, Invalid) // ex. ("Invalid requested_security_strength", -1)
	//  3: end if

	//  4: if ((prediction_resistance_flag is set) AND
	//         (prediction resistance is not supported) then
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

	// 24: seedMaterial ← entropy_input ‖ nonce ‖ personalization_string
	// seedMaterial := drbg.WriteBytes(entropyInput, nonce, personalizationString)

	buf := make([]byte, outlen*3)
	var (
		key = buf[:outlen]
		v   = buf[outlen : outlen*2]
		sum = buf[outlen*2:]
	)
	// 25: Key ← 0x00 00...00
	// 26: V ← 0x01 01...01
	memory.Memset(v, 1)

	state := &State{
		New:                         h,
		ReseedInterval:              reseedInterval,
		PredictionResistanceRequest: predictionResistanceRequest,

		Key:           key,
		V:             v,
		ReseedCounter: 1,

		sum: sum,
	}
	// 27: (Key, V ) ← HMAC_DRBG_Update(seed_material, Key, V )
	state.Update(entropyInput, nonce, personalizationString)

	// 28: (status, state_handle) ← Find_state_space()
	// 29: if (status ≠ SUCCESS) then
	// 30:     return (status, Invalid)
	// 31: end if
	// 32: state(state_handle ).V ← V
	// 33: state(state_handle ).Key ← Key
	// 34: state(state_handle ).reseed_counter ← 1
	// 35: state(state_handle ).security_strength ← security_strength
	// 36: state(state_handle ).prediction_resistance_flag ← prediction_resistance_flag
	// 37: return (SUCCESS, state_handle)

	return state
}

// 6.4 리시드 함수(reseeding function)
//
// HMAC_DRBG의 리시드 함수 Reseed_HMAC_DRBG는 인스턴스의 내부 상태와 엔트로피 입력,
// 그리고 추가 입력을 이용하여 새로운 시드 seed = (Key, V )를 생성하고,
// 이 값을 이용하여 내부 상태를 갱신한다.
func (state *State) Reseed(entropyInput, additionalInput []byte) {
	//  1: if (a state(state_handle ) is not available) then
	//  2:     return ERROR_FLAG // ex. ‶State not available for the state_handle″
	//  3: end if

	//  4: V ← state(state_handle ).V
	//  5: Key ← state(state_handle ).Key

	//  6: security_strength ← state(state_handle ).security_strength

	//  7: (status, entropy_input ) ← Get_entropy(min_entropy, min_len, max_len, request )
	//         // min_len = min_length, max_len = max_length, request = prediction_resistance_request
	//  8: if (status ≠ SUCCESS) then
	//  9:     return (status )
	// 10: end if

	// 11: seedMaterial ← entropy_input ‖ additional_input
	// seedMaterial := drbg.WriteBytes(entropyInput, additionalInput)
	// 12: (Key, V ) ← HMAC_DRBG_Update(seed_material, Key, V )
	state.Update(entropyInput, additionalInput)

	// 13: state(state_handle ).V ← V
	// 14: state(state_handle ).Key ← Key
	// 15: state(state_handle ).reseed_counter ← 1
	// 16: return (SUCCESS)
}

// 6.5 생성 함수(generate function)
//
// HMAC_DRBG의 생성 함수 Generate_HMAC_DRBG는 인스턴스의 내부 상태 중
// 동작 상태의 Key 를 고정하고 V 를 갱신하면서 출력값을 생성한 후,
// Key 와 V 를 갱신한다.
func (state *State) Generate(
	dst []byte,
	fnEntropyInput func() ([]byte, error),
	additionalInput []byte,
) error {
	//  1: if (a state(state_handle ) is not available) then
	//  2:     return (ERROR_FLAG, Null )
	//  3: end if

	//  4: V ← state(state_handle).V
	//  5: Key ← state(state_handle).Key
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

	// 18: if ((prediction_resistance_request is set) and
	//         (prediction_resistance_flag is not set)) then
	// 19:     return (ERROR_FLAG, Null )
	// 20: end if

	// 21: if ((reseed_counter > reseed_interval) OR
	//         (prediction_resistance_request is set)) then
	if state.ReseedCounter > state.ReseedInterval || state.PredictionResistanceRequest {
		// 22: if reseeding is not available then
		// 23:     return (ERROR_FLAG, Null)
		// 24: end if

		// 25: status ← Reseed_HMAC_DRBG(handle, request, additional_input)
		//         // handle = state_handle, request = prediction_resistance_request
		// 26: if (status ≠ SUCCESS) then
		// 27:     return (status, Null)
		// 28: end if
		entropyInput, err := fnEntropyInput()
		if err != nil {
			return err
		}
		state.Reseed(entropyInput, additionalInput)

		// 29: V ← state(state_handle ).V
		// 30: Key ← state(state_handle ).Key
		// 31: reseed_counter ← state(state_handle ).reseed_counter
		// 32: additional_input ← Null
		additionalInput = nil
	}
	// 33: end if

	// 34: if (additional_input ≠ Null) then
	if additionalInput != nil {
		// 35:     (Key, V ) ← HMAC_DRBG_Update(additional_input, Key, V )
		state.Update(additionalInput)
	}
	// 36: end if

	// 37: temp ← Null
	// 38: While (len(temp) < requested_no_of_bits ) do
	for len(dst) > 0 {
		// 39:     V ← HMAC(Key, V )
		h := hmac.New(state.New, state.Key)
		h.Write(state.V)
		copy(state.V, h.Sum(state.sum[:0]))

		// 40:     temp ← temp ‖ V
		dst = dst[copy(dst, state.V):]
	}
	// 41: end while

	// 42: pseudorandom_bits ← leftmost(temp, requested_no_of_bits )

	// 43: (Key, V ) ← HMAC_DRBG_Update(additional_input, Key, V )
	if additionalInput != nil {
		state.Update(additionalInput)
	} else {
		state.Update()
	}

	// 44: state(state_handle).V ← V
	// 45: state(state_handle).Key ← Key
	// 46: state(state_handle).reseed_counter ← reseed_counter + 1
	state.ReseedCounter++

	// 47: return (SUCCESS, pseudorandom_bits)
	return nil
}
