package cavp

import (
	"bufio"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/kipher"

	"github.com/stretchr/testify/require"
)

type (
	GMACAe struct {
		KeyLen int `cavp:"KeyLen,header"`
		IVLen  int `cavp:"IVLen,header"`
		PTLen  int `cavp:"PTlen,header"`
		AADLen int `cavp:"AADlen,header"`
		TagLen int `cavp:"TagLen,header"`

		Count      int    `cavp:"COUNT"`
		Key        []byte `cavp:"Key,hex"`
		IV         []byte `cavp:"IV,hex"`
		PlainText  []byte `cavp:"PT,hex"`
		Adata      []byte `cavp:"Adata,hex"`
		Ciphertext []byte `cavp:"C,hex,optional"`
		Tag        []byte `cavp:"T,hex,optional"` // Output
	}

	GMACAd struct {
		KeyLen int `cavp:"KeyLen,header"`
		IVLen  int `cavp:"IVlen,header"`
		PTLen  int `cavp:"PTLen,header"`
		AADLen int `cavp:"AADLen,header"`
		TagLen int `cavp:"TagLen,header"`

		Count      int    `cavp:"COUNT"`
		Key        []byte `cavp:"Key,hex"`
		IV         []byte `cavp:"IV,hex"`
		Adata      []byte `cavp:"Adata,hex"`
		CipherText []byte `cavp:"C,hex"`
		Tag        []byte `cavp:"T,hex"`
		PlainText  []byte `cavp:"PT,hex,optional"`  // Output
		Invalid    bool   `cavp:"INVALID,optional"` // Output
	}
)

func ProcessGCMAe(
	input, output string,
	newCipher func([]byte) (cipher.Block, error),
) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	response := make([]string, 0, 4096)
	paragraph := make([]string, 0, 8)
	var config gcmConfig
	flush := func() error {
		if len(paragraph) == 0 {
			return nil
		}
		request, isRecord, parseErr := parseGCMParagraph(paragraph, &config, false)
		if parseErr != nil {
			return parseErr
		}
		response = append(response, paragraph...)
		paragraph = paragraph[:0]
		if !isRecord {
			return nil
		}
		key, iv, plaintext, adata, sizeErr := request.sizedInputs(config, false)
		if sizeErr != nil {
			return sizeErr
		}
		b, cipherErr := newCipher(key)
		if cipherErr != nil {
			return cipherErr
		}
		gcm, cipherErr := kipher.NewGCMWithSize(b, len(iv), config.tagBits/8)
		if cipherErr != nil {
			return cipherErr
		}
		sealed := gcm.Seal(nil, iv, plaintext, adata)
		response = append(
			response,
			"C = "+upperHex(sealed[:len(plaintext)]),
			"T = "+upperHex(sealed[len(plaintext):]),
		)
		return nil
	}
	for scanner.Scan() {
		raw := strings.TrimSuffix(scanner.Text(), "\r")
		if strings.TrimSpace(raw) == "" {
			if err := flush(); err != nil {
				return err
			}
			response = append(response, raw)
			continue
		}
		paragraph = append(paragraph, raw)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := flush(); err != nil {
		return err
	}
	return writeRecordResponse(output, response)
}

func TestGCMAe(
	t testing.TB,
	input string,
	newCipher func([]byte) (cipher.Block, error),
) {
	t.Helper()

	Test(t, input, func(v GMACAe) {
		keyLen := v.KeyLen / 8
		ivLen := v.IVLen / 8
		ptLen := v.PTLen / 8
		aadLen := v.AADLen / 8
		tagLen := v.TagLen / 8

		b, err := newCipher(v.Key[:keyLen])
		require.NoError(t, err)

		testAEAD(
			t,
			v.IV[:ivLen], v.PlainText[:ptLen], v.Adata[:aadLen],
			v.Ciphertext[:ptLen], v.Tag[:tagLen],
			func() (cipher.AEAD, error) {
				return kipher.NewGCMWithSize(b, ivLen, tagLen)
			},
		)
	})
}

func ProcessGCMAd(
	input, output string,
	newCipher func([]byte) (cipher.Block, error),
) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	response := make([]string, 0, 4096)
	paragraph := make([]string, 0, 8)
	var config gcmConfig
	flush := func() error {
		if len(paragraph) == 0 {
			return nil
		}
		request, isRecord, parseErr := parseGCMParagraph(paragraph, &config, true)
		if parseErr != nil {
			return parseErr
		}
		response = append(response, paragraph...)
		paragraph = paragraph[:0]
		if !isRecord {
			return nil
		}
		key, iv, ciphertext, adata, sizeErr := request.sizedInputs(config, true)
		if sizeErr != nil {
			return sizeErr
		}
		tagBytes := config.tagBits / 8
		if len(request.tag) < tagBytes {
			return fmt.Errorf("cavp: T is %d bytes, want at least %d", len(request.tag), tagBytes)
		}
		b, cipherErr := newCipher(key)
		if cipherErr != nil {
			return cipherErr
		}
		gcm, cipherErr := kipher.NewGCMWithSize(b, len(iv), tagBytes)
		if cipherErr != nil {
			return cipherErr
		}
		sealed := append(append(make([]byte, 0, len(ciphertext)+tagBytes), ciphertext...), request.tag[:tagBytes]...)
		plaintext, openErr := gcm.Open(nil, iv, sealed, adata)
		if openErr != nil {
			response = append(response, "Invalid")
		} else {
			response = append(response, "PT = "+upperHex(plaintext))
		}
		return nil
	}
	for scanner.Scan() {
		raw := strings.TrimSuffix(scanner.Text(), "\r")
		if strings.TrimSpace(raw) == "" {
			if err := flush(); err != nil {
				return err
			}
			response = append(response, raw)
			continue
		}
		paragraph = append(paragraph, raw)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := flush(); err != nil {
		return err
	}
	return writeRecordResponse(output, response)
}

type gcmConfig struct {
	keyBits, ivBits, ptBits, aadBits, tagBits int
	seen                                      map[string]bool
}

type gcmRequest struct {
	key, iv, plaintext, adata, ciphertext, tag []byte
}

func parseGCMParagraph(lines []string, config *gcmConfig, decrypt bool) (gcmRequest, bool, error) { //nolint:gocyclo
	// Each bracketed paragraph describes a complete independent configuration.
	// Reset it before parsing so an omitted header cannot inherit a value from
	// the preceding test group.
	for _, raw := range lines {
		if _, _, ok := parseGCMHeader(raw); ok {
			*config = gcmConfig{}
			break
		}
	}
	if config.seen == nil {
		config.seen = make(map[string]bool, 5)
	}
	var request gcmRequest
	seen := make(map[string]bool, 7)
	isRecord := false
	for _, raw := range lines {
		if name, value, ok := parseGCMHeader(raw); ok {
			key := strings.ToLower(name)
			var dst *int
			switch key {
			case "keylen":
				dst = &config.keyBits
			case "ivlen":
				dst = &config.ivBits
			case "ptlen":
				dst = &config.ptBits
			case "aadlen":
				dst = &config.aadBits
			case "taglen":
				dst = &config.tagBits
			}
			if dst != nil {
				n, err := strconv.Atoi(value)
				if err != nil {
					return request, false, fmt.Errorf("cavp: invalid %s: %w", name, err)
				}
				*dst = n
				config.seen[key] = true
			}
			continue
		}
		entry := parseLine(raw)
		key := strings.ToLower(entry.key)
		var dst *[]byte
		switch key {
		case "count":
			isRecord = true
		case "key":
			dst, isRecord = &request.key, true
		case "iv":
			dst, isRecord = &request.iv, true
		case "pt":
			dst, isRecord = &request.plaintext, true
		case "adata":
			dst, isRecord = &request.adata, true
		case "c":
			dst, isRecord = &request.ciphertext, true
		case "t":
			dst, isRecord = &request.tag, true
		}
		if dst != nil {
			var err error
			*dst, err = decodeGCMHex(entry.val, key)
			if err != nil {
				return request, true, err
			}
		}
		if key != "" {
			seen[key] = true
		}
	}
	if !isRecord {
		return request, false, nil
	}
	for _, key := range []string{"keylen", "ivlen", "ptlen", "aadlen", "taglen"} {
		if !config.seen[key] {
			return request, true, fmt.Errorf("cavp: GCM configuration is missing %s", key)
		}
	}
	required := []string{"key", "iv", "adata"}
	if decrypt {
		required = append(required, "c", "t")
	} else {
		required = append(required, "pt")
	}
	for _, key := range required {
		if !seen[key] {
			return request, true, fmt.Errorf("cavp: incomplete GCM record: missing %s", key)
		}
	}
	for name, bits := range map[string]int{"KeyLen": config.keyBits, "IVLen": config.ivBits, "PTLen": config.ptBits, "AADLen": config.aadBits, "TagLen": config.tagBits} {
		if bits < 0 || bits%8 != 0 {
			return request, true, fmt.Errorf("cavp: invalid %s %d", name, bits)
		}
	}
	return request, true, nil
}

func (r gcmRequest) sizedInputs(config gcmConfig, decrypt bool) (key, iv, text, adata []byte, err error) {
	textValue := r.plaintext
	textName := "PT"
	if decrypt {
		textValue = r.ciphertext
		textName = "C"
	}
	values := []struct {
		name string
		data []byte
		bits int
	}{
		{"Key", r.key, config.keyBits},
		{"IV", r.iv, config.ivBits},
		{textName, textValue, config.ptBits},
		{"Adata", r.adata, config.aadBits},
	}
	for _, value := range values {
		need := value.bits / 8
		if len(value.data) < need {
			return nil, nil, nil, nil, fmt.Errorf("cavp: %s is %d bytes, want at least %d", value.name, len(value.data), need)
		}
	}
	return r.key[:config.keyBits/8], r.iv[:config.ivBits/8], textValue[:config.ptBits/8], r.adata[:config.aadBits/8], nil
}

func parseGCMHeader(line string) (name, value string, ok bool) {
	line = strings.TrimSpace(line)
	if len(line) < 2 || line[0] != '[' || line[len(line)-1] != ']' {
		return "", "", false
	}
	entry := parseLine(line[1 : len(line)-1])
	return entry.key, entry.val, entry.has
}

func decodeGCMHex(value, name string) ([]byte, error) {
	decoded, err := hex.DecodeString(strings.Join(strings.Fields(value), ""))
	if err != nil {
		return nil, fmt.Errorf("cavp: invalid %s: %w", name, err)
	}
	return decoded, nil
}

func TestGCMAd(
	t testing.TB,
	input string,
	newCipher func([]byte) (cipher.Block, error),
) {
	t.Helper()

	Test(t, input, func(v GMACAd) {
		keyLen := v.KeyLen / 8
		ivLen := v.IVLen / 8
		ptLen := v.PTLen / 8
		aadLen := v.AADLen / 8
		tagLen := v.TagLen / 8

		b, err := newCipher(v.Key[:keyLen])
		require.NoError(t, err)

		gcm, err := kipher.NewGCMWithSize(b, ivLen, tagLen)
		require.NoError(t, err)

		ciphertext := internal.Grow(nil, ptLen+tagLen)[:0]
		ciphertext = append(append(ciphertext, v.CipherText[:ptLen]...), v.Tag[:tagLen]...)

		out, err := gcm.Open(nil, v.IV[:ivLen], ciphertext, v.Adata[:aadLen])
		if v.Invalid {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, v.PlainText[:ptLen], out)
		}
	})
}
