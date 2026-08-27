package cavp

import (
	"bufio"
	"bytes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"

	"github.com/stretchr/testify/require"
)

type (
	BlockKatMmt struct {
		Key        []byte `cavp:"KEY,hex"`
		IV         []byte `cavp:"IV,hex,optional"`
		Counter    []byte `cavp:"CTR,hex,optional"`
		PlainText  []byte `cavp:"PT,hex,optional"`
		CipherText []byte `cavp:"CT,hex,optional"`
	}

	BlockMct struct {
		Count      int    `cavp:"COUNT,int"`
		Key        []byte `cavp:"KEY,hex"`
		IV         []byte `cavp:"IV,hex,optional"`
		Counter    []byte `cavp:"CTR,hex,optional"`
		PlainText  []byte `cavp:"PT,hex"`
		CipherText []byte `cavp:"CT,hex,optional"`
	}
)

type blockRequest struct {
	lines      []string
	key        []byte
	iv         []byte
	counter    []byte
	plainText  []byte
	cipherText []byte
	hasKey     bool
	hasIV      bool
	hasCounter bool
	hasPT      bool
	hasCT      bool
}

type NewCipherFunc func([]byte) (cipher.Block, error)

func ProcessBlockKAT(input, output string, mode BlockMode, newCipher NewCipherFunc) error {
	return processBlockKatMmt(input, output, mode, newCipher)
}

func ProcessBlockMMT(input, output string, mode BlockMode, newCipher NewCipherFunc) error {
	return processBlockKatMmt(input, output, mode, newCipher)
}

func processBlockKatMmt(
	input, output string,
	mode BlockMode,
	newCipher NewCipherFunc,
) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := os.Create(output)
	if err != nil {
		return err
	}
	defer w.Close()

	scanner := bufio.NewScanner(r)
	scanner.Split(scanBlockLines)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	writer := bufio.NewWriter(w)
	eol := "\n"
	lineNumber := 0
	wroteRecord := false
	var req blockRequest

	writeRecord := func() error {
		if len(req.lines) == 0 {
			return nil
		}
		if !req.hasKey {
			return fmt.Errorf("cavp: block request is missing KEY")
		}
		if req.hasPT == req.hasCT {
			return fmt.Errorf("cavp: block request must contain exactly one of PT or CT")
		}

		b, err := newCipher(req.key)
		if err != nil {
			return fmt.Errorf("create cipher: %w", err)
		}

		var label string
		var result []byte
		if req.hasPT {
			label = "CT"
			result = make([]byte, len(req.plainText))
			err = CryptBlock(b, mode, req.iv, req.counter, result, req.plainText, true)
		} else {
			label = "PT"
			result = make([]byte, len(req.cipherText))
			err = CryptBlock(b, mode, req.iv, req.counter, result, req.cipherText, false)
		}
		if err != nil {
			return err
		}

		if wroteRecord {
			if _, err = writer.WriteString(eol + eol); err != nil {
				return err
			}
		}
		for _, line := range req.lines {
			if _, err = writer.WriteString(line + eol); err != nil {
				return err
			}
		}
		if _, err = fmt.Fprintf(writer, "%s = %s", label, strings.ToUpper(hex.EncodeToString(result))); err != nil {
			return err
		}
		wroteRecord = true
		req = blockRequest{}
		return nil
	}

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if strings.HasSuffix(line, "\r") {
			eol = "\r\n"
			line = strings.TrimSuffix(line, "\r")
		}
		if strings.TrimSpace(line) == "" {
			if err = writeRecord(); err != nil {
				return err
			}
			continue
		}
		req.lines = append(req.lines, line)
		if err = req.parseLine(line); err != nil {
			return fmt.Errorf("%s:%d: %w", input, lineNumber, err)
		}
	}
	if err = scanner.Err(); err != nil {
		return err
	}
	if err = writeRecord(); err != nil {
		return err
	}
	return writer.Flush()
}

func TestBlockKAT(
	t *testing.T,
	filename string,
	mode BlockMode,
	newCipher NewCipherFunc,
) {
	t.Run("Encrypt", func(t *testing.T) {
		var got []byte
		Test(t, filename, func(v BlockKatMmt) {
			b, err := newCipher(v.Key)
			require.NoError(t, err)

			if len(v.PlainText) != 0 {
				got = internal.Grow(got, len(v.PlainText))

				err := CryptBlock(b, mode, v.IV, v.Counter, got, v.PlainText, true)
				require.NoError(t, err)
				require.Equal(t, v.CipherText, got, "CipherText should match")
			}
		})
	})

	t.Run("Decrypt", func(t *testing.T) {
		var got []byte
		Test(t, filename, func(v BlockKatMmt) {
			b, err := newCipher(v.Key)
			require.NoError(t, err)

			if len(v.CipherText) != 0 {
				got = internal.Grow(got, len(v.CipherText))

				err := CryptBlock(b, mode, v.IV, v.Counter, got, v.CipherText, false)
				require.NoError(t, err)
				require.Equal(t, v.PlainText, got, "PlainText should match")
			}
		})
	})
}

func ProcessBlockMCT(
	input, output string,
	mode BlockMode,
	newCipher NewCipherFunc,
	iterations int,
) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(output)
	if err != nil {
		return err
	}
	defer w.Close()

	scanner := bufio.NewScanner(r)
	scanner.Split(scanBlockLines)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	writer := bufio.NewWriter(w)
	eol := "\n"
	lineNumber := 0
	wroteRecord := false
	var req blockRequest

	writeRequest := func() error {
		if len(req.lines) == 0 {
			return nil
		}
		if !req.hasKey || !req.hasPT {
			return fmt.Errorf("cavp: block MCT request requires KEY and PT")
		}
		state := BlockMct{
			Key:       bytes.Clone(req.key),
			IV:        bytes.Clone(req.iv),
			Counter:   bytes.Clone(req.counter),
			PlainText: bytes.Clone(req.plainText),
		}
		for iteration := range iterations {
			result, next, err := blockRound(state, mode, newCipher)
			if err != nil {
				return err
			}
			if wroteRecord {
				if _, err = writer.WriteString(eol + eol); err != nil {
					return err
				}
			}
			if _, err = fmt.Fprintf(writer, "COUNT = %d%s", result.Count, eol); err != nil {
				return err
			}
			if iteration == 0 {
				for _, line := range req.lines {
					if _, err = writer.WriteString(line + eol); err != nil {
						return err
					}
				}
			} else {
				if err = writeBlockMCTState(writer, eol, result); err != nil {
					return err
				}
			}
			if _, err = fmt.Fprintf(writer, "CT = %s", strings.ToUpper(hex.EncodeToString(result.CipherText))); err != nil {
				return err
			}
			wroteRecord = true
			state = next
		}
		req = blockRequest{}
		return nil
	}

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if strings.HasSuffix(line, "\r") {
			eol = "\r\n"
			line = strings.TrimSuffix(line, "\r")
		}
		if strings.TrimSpace(line) == "" {
			if err = writeRequest(); err != nil {
				return err
			}
			continue
		}
		req.lines = append(req.lines, line)
		if err = req.parseLine(line); err != nil {
			return fmt.Errorf("%s:%d: %w", input, lineNumber, err)
		}
	}
	if err = scanner.Err(); err != nil {
		return err
	}
	if err = writeRequest(); err != nil {
		return err
	}
	return writer.Flush()
}

func (r *blockRequest) parseLine(line string) error {
	key, value, ok := strings.Cut(line, "=")
	if !ok {
		return nil
	}
	key = strings.ToUpper(strings.TrimSpace(key))
	if key != "KEY" && key != "IV" && key != "CTR" && key != "PT" && key != "CT" {
		return nil
	}
	decoded, err := hex.DecodeString(strings.TrimSpace(value))
	if err != nil {
		return fmt.Errorf("%s: invalid hexadecimal value: %w", key, err)
	}
	switch key {
	case "KEY":
		r.key, r.hasKey = decoded, true
	case "IV":
		r.iv, r.hasIV = decoded, true
	case "CTR":
		r.counter, r.hasCounter = decoded, true
	case "PT":
		r.plainText, r.hasPT = decoded, true
	case "CT":
		r.cipherText, r.hasCT = decoded, true
	}
	return nil
}

func scanBlockLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, data[:i], nil
	}
	if atEOF && len(data) != 0 {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func writeBlockMCTState(w *bufio.Writer, eol string, state BlockMct) error {
	lines := []string{"KEY = " + strings.ToUpper(hex.EncodeToString(state.Key))}
	if len(state.Counter) != 0 {
		lines = append(lines, "CTR = "+strings.ToUpper(hex.EncodeToString(state.Counter)))
	} else if len(state.IV) != 0 {
		lines = append(lines, "IV = "+strings.ToUpper(hex.EncodeToString(state.IV)))
	}
	lines = append(lines, "PT = "+strings.ToUpper(hex.EncodeToString(state.PlainText)))
	for _, line := range lines {
		if _, err := w.WriteString(line + eol); err != nil {
			return err
		}
	}
	return nil
}

func blockRound(
	state BlockMct,
	mode BlockMode,
	newCipher NewCipherFunc,
) (BlockMct, BlockMct, error) {
	b, err := newCipher(state.Key)
	if err != nil {
		return BlockMct{}, BlockMct{}, err
	}
	sz := len(state.PlainText)
	if sz == 0 {
		return BlockMct{}, BlockMct{}, fmt.Errorf("MCT requires plaintext")
	}
	if segment := mode.SegmentSize(); segment != 0 && segment != sz {
		return BlockMct{}, BlockMct{}, fmt.Errorf("%s MCT requires %d bytes", mode, segment)
	}
	if segment := mode.SegmentSize(); segment != 0 && b.BlockSize()%segment != 0 {
		return BlockMct{}, BlockMct{}, fmt.Errorf("CFB segment does not divide block size")
	}
	crypt, err := NewBlockCrypt(b, mode, state.IV, state.Counter, true)
	if err != nil {
		return BlockMct{}, BlockMct{}, err
	}
	result, plain := state.clone(), bytes.Clone(state.PlainText)
	history := make([]byte, 1000*sz)
	feedback := b.BlockSize() / sz
	iv := state.IV
	if mode == CTR && len(state.Counter) != 0 {
		iv = state.Counter
	}
	for iter := range 1000 {
		ct := history[iter*sz : (iter+1)*sz]
		crypt(ct, plain)
		switch mode {
		case ECB, CTR:
			copy(plain, ct)
		case CBC, OFB:
			if iter == 0 {
				copy(plain, iv)
			} else {
				copy(plain, history[(iter-1)*sz:iter*sz])
			}
		case CFB8, CFB32, CFB64, CFB128:
			if iter < feedback {
				copy(plain, state.IV[iter*sz:(iter+1)*sz])
			} else {
				copy(plain, history[(iter-feedback)*sz:(iter-feedback+1)*sz])
			}
		}
	}
	result.CipherText = bytes.Clone(history[len(history)-sz:])
	next := result.clone()
	next.Count++
	next.PlainText = bytes.Clone(plain)
	tail := history[len(history)-len(next.Key):]
	for i := range next.Key {
		next.Key[i] ^= tail[i]
	}
	switch mode {
	case CBC, OFB:
		next.IV = bytes.Clone(result.CipherText)
	case CFB8, CFB32, CFB64, CFB128:
		next.IV = bytes.Clone(history[len(history)-b.BlockSize():])
	case CTR:
		counter := &next.IV
		if len(state.Counter) != 0 {
			next.Counter = bytes.Clone(state.Counter)
			counter = &next.Counter
		}
		for range (1000*sz + b.BlockSize() - 1) / b.BlockSize() {
			internal.IncCtr(*counter)
		}
	}
	return result, next, nil
}

func (v BlockMct) clone() BlockMct {
	return BlockMct{
		Count:      v.Count,
		Key:        bytes.Clone(v.Key),
		IV:         bytes.Clone(v.IV),
		Counter:    bytes.Clone(v.Counter),
		PlainText:  bytes.Clone(v.PlainText),
		CipherText: bytes.Clone(v.CipherText),
	}
}
