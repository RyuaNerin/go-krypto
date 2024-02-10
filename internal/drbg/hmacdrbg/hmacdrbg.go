package hmacdrbg

import (
	"crypto/hmac"
	"hash"
	"io"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/drbg"
	"github.com/RyuaNerin/go-krypto/internal/kryptoutil"
)

const (
	MaxLength                      = 0x8_0000_0000      // 엔트로피 입력 최대 길이 (max_length)
	MaxPersonalizationStringLength = 0x8_0000_0000      // 개별화 문자열 최대 허용 길이 (max_personalization_string_length) = 2 ** 35
	MaxAdditionalInputLength       = 0x8_0000_0000      // 추가 입력 최대 허용 길이 (max_additional_input_length) = 2 ** 35
	MaxNoOfBitsPerRequest          = 0x8_0000           // 난수 최대 출력 길이 (max_no_of_bits_per_request) = 2 ** 19
	ReseedInterval                 = 0x8_0000_0000_0000 // 시드별 출력값 생성 횟수 (reseed_interval) = 2 ** 48
)

type State struct {
	h func() hash.Hash

	outlen int
	sum    []byte

	v   []byte
	key []byte

	reseed_counter             uint64 // 시드 생성 이후 DRBG 인스턴스의 출력값 생성 횟수
	security_strength          int    // DRBG 인스턴스의 안전성 수준
	prediction_resistance_flag bool
}

func GetSecurityStrength(requested_strength int) int {
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
func (state *State) HMAC_DRBG_Update(provided_data func(w io.Writer)) {
	// 1: Key ← HMAC(Key, V ‖ 0x00 ‖ provided_data )
	h := hmac.New(state.h, state.key)
	h.Write(state.v)
	h.Write([]byte{0x00})
	if provided_data != nil {
		provided_data(h)
	}
	copy(state.key, h.Sum(state.sum[:0]))

	// 2: V ← HMAC(Key, V )
	h = hmac.New(state.h, state.key)
	h.Write(state.v)
	copy(state.v, h.Sum(state.sum[:0]))

	// 3: if (provided_data = Null) then
	// 4:     return (Key, V )
	// 5: end if
	if provided_data == nil {
		return
	}

	// 6: Key ← HMAC(Key, V ‖ 0x01 ‖ provided_data )
	h = hmac.New(state.h, state.key)
	h.Write(state.v)
	h.Write([]byte{0x01})
	provided_data(h)
	copy(state.key, h.Sum(state.sum[:0]))

	// 7: V ← HMAC(Key, V )
	h = hmac.New(state.h, state.key)
	h.Write(state.v)
	copy(state.v, h.Sum(state.sum[:0]))

	// 8: return (Key, V )
}

// 6.3 인스턴스 생성 함수(instantiate function)
//
// HMAC_DRBG의 인스턴스 생성 함수 Instantiate_HMAC_DRBG는
// 엔트로피 입력, 논스, 개별화 문자열로부터 시드 seed = (Key, V )를 생성하고,
// 이 시드를 이용하여 내부 상태를 초기화한다.
func Instantiate_HMAC_DRBG(
	h func() hash.Hash,
	requested_strength int,
	entropy_input []byte,
	nonce, personalization_string []byte,
	prediction_resistance_flag bool,
) *State {
	outlen := h().Size()

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
	security_strength := GetSecurityStrength(requested_strength)

	// 19: (status, entropy_input) ← Get_entropy(min_entropy, min_len, max_len, request )
	//  // min_len = min_length, max_len = max_length, request = prediction_resistance_request
	// 20: if (status ≠ SUCCESS) then
	// 21:     return (status, Invalid)
	// 22: end if

	// 23: obtain a nonce and check its acceptability

	// 24: seed_material ← entropy_input ‖ nonce ‖ personalization_string
	seed_material := drbg.WriteBytes(entropy_input, nonce, personalization_string)

	KV := make([]byte, outlen*2)
	// 25: Key ← 0x00 00...00
	// 26: V ← 0x01 01...01
	Key := KV[:outlen]
	V := KV[outlen:]
	kryptoutil.MemsetByte(V, 1)

	state := &State{
		h:      h,
		outlen: outlen,
		sum:    make([]byte, outlen),

		key: Key,
		v:   V,

		reseed_counter:             1,
		security_strength:          security_strength,
		prediction_resistance_flag: prediction_resistance_flag,
	}
	// 27: (Key, V ) ← HMAC_DRBG_Update(seed_material, Key, V )
	state.HMAC_DRBG_Update(seed_material)

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
func (state *State) Reseed_HMAC_DRBG(
	entropy_input []byte,
	additional_input []byte,
) {
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

	// 11: seed_material ← entropy_input ‖ additional_input
	seed_material := drbg.WriteBytes(entropy_input, additional_input)
	// 12: (Key, V ) ← HMAC_DRBG_Update(seed_material, Key, V )
	state.HMAC_DRBG_Update(seed_material)

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
func (state *State) Generate_HMAC_DRBG(
	dst []byte,
	entropy_input func() ([]byte, error),
	additional_input []byte,
) error {
	requested_no_of_bits := len(dst) * 8

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
	if state.reseed_counter > ReseedInterval || state.prediction_resistance_flag {
		// 22: if reseeding is not available then
		// 23:     return (ERROR_FLAG, Null)
		// 24: end if

		// 25: status ← Reseed_HMAC_DRBG(handle, request, additional_input)
		//         // handle = state_handle, request = prediction_resistance_request
		// 26: if (status ≠ SUCCESS) then
		// 27:     return (status, Null)
		// 28: end if
		entropy_input, err := entropy_input()
		if err != nil {
			return err
		}
		state.Reseed_HMAC_DRBG(entropy_input, additional_input)

		// 29: V ← state(state_handle ).V
		// 30: Key ← state(state_handle ).Key
		// 31: reseed_counter ← state(state_handle ).reseed_counter
		// 32: additional_input ← Null
		additional_input = nil
	}
	// 33: end if

	// 34: if (additional_input ≠ Null) then
	if len(additional_input) > 0 {
		// 35:     (Key, V ) ← HMAC_DRBG_Update(additional_input, Key, V )
		state.HMAC_DRBG_Update(drbg.WriteBytes(additional_input))
	}
	// 36: end if

	// 37: temp ← Null
	countMax := internal.CeilDiv(requested_no_of_bits, state.outlen*8)
	// 38: While (len(temp) < requested_no_of_bits ) do
	for count := 0; count < countMax; count++ {
		// 39:     V ← HMAC(Key, V )
		h := hmac.New(state.h, state.key)
		h.Write(state.v)
		copy(state.v, h.Sum(state.sum[:0]))

		// 40:     temp ← temp ‖ V
		copy(dst[state.outlen*count:], state.v)
	}
	// 41: end while

	// 42: pseudorandom_bits ← leftmost(temp, requested_no_of_bits )

	// 43: (Key, V ) ← HMAC_DRBG_Update(additional_input, Key, V )
	state.HMAC_DRBG_Update(drbg.WriteBytes(additional_input))

	// 44: state(state_handle).V ← V
	// 45: state(state_handle).Key ← Key
	// 46: state(state_handle).reseed_counter ← reseed_counter + 1
	state.reseed_counter++

	// 47: return (SUCCESS, pseudorandom_bits)
	return nil
}

// 6.6 인스턴스 소멸 함수(uninstantiate function)
//
// HMAC_DRBG의 인스턴스 소멸 함수 Uninstantiate_HMAC_DRBG는
// 인스턴스의 내부 상태를 해제하여 인스턴스가 더 이상 동작할 수 없도록 한다.
func (s *State) Uninstantiate_HMAC_DRBG() {
	kryptoutil.MemsetByte(s.v, 0)
	kryptoutil.MemsetByte(s.key, 0)
	s.v = nil
	s.key = nil
	s.reseed_counter = 0
	s.security_strength = 0
	s.prediction_resistance_flag = false
}
