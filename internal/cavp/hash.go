package cavp

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type HashSmtLmt struct {
	HashSize      int    `cavp:"L,int,header=plain"`
	MessageLength int    `cavp:"Len,int"`
	Message       []byte `cavp:"Msg,hex"`
	Digest        []byte `cavp:"MD,hex,optional"`
}

type (
	HashMct struct {
		HashSize int    `cavp:"L,int,header=plain,optional"`
		Seed     []byte `cavp:"Seed,hex,optional"`
		Count    int    `cavp:"COUNT,int"`
		Digest   []byte `cavp:"MD,hex"`
	}
)

type (
	Digest interface {
		~[]byte |
			~[20]byte |
			~[28]byte | ~[32]byte |
			~[48]byte | ~[64]byte
	}
	NewHashFunc interface {
		func(size int) hash.Hash | func() hash.Hash
	}
	SumFunc[T Digest] interface {
		func(size int, data []byte) T | func(data []byte) T
	}
)

func normalizeHashFunction[F NewHashFunc](h F, size int) func() hash.Hash {
	switch f := any(h).(type) {
	case func(size int) hash.Hash:
		return func() hash.Hash { return f(size) }
	case func() hash.Hash:
		return f
	default:
		panic(fmt.Sprintf("invalid type %T", h))
	}
}

func normalizeSumFunction[T Digest, F SumFunc[T]](h F, size int) func([]byte) T {
	switch f := any(h).(type) {
	case func(size int, data []byte) T:
		return func(data []byte) T { return f(size, data) }
	case func(data []byte) T:
		return f
	default:
		panic(fmt.Sprintf("invalid type %T", h))
	}
}

func ProcessHashSMT(input, output string, newHash func() hash.Hash) error {
	return processHashSmtLmt(input, output, newHash)
}

func ProcessHashLMT(input, output string, newHash func() hash.Hash) error {
	return processHashSmtLmt(input, output, newHash)
}

func ProcessHashMCT(
	input, output string,
	newHash func() hash.Hash,
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
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)

	lines := make([]string, 0, iterations*3+4)
	hashSize := 0
	var requestSeed []byte
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := scanner.Text()
		lines = append(lines, line)

		name, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch strings.TrimSpace(name) {
		case "L":
			hashSize, err = strconv.Atoi(strings.TrimSpace(value))
			if err != nil {
				return fmt.Errorf("cavp: line %d: invalid L: %w", lineNumber, err)
			}
		case "Seed":
			requestSeed, err = hex.DecodeString(strings.TrimSpace(value))
			if err != nil {
				return fmt.Errorf("cavp: line %d: invalid Seed: %w", lineNumber, err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if requestSeed == nil {
		return fmt.Errorf("cavp: MCT requires a seed")
	}

	h := newHash()
	if h == nil {
		return fmt.Errorf("cavp: hash constructor returned nil")
	}
	if len(requestSeed) != h.Size() {
		return fmt.Errorf("cavp: MCT seed is %d bytes, want %d", len(requestSeed), h.Size())
	}
	if hashSize != 0 && hashSize != h.Size() {
		return fmt.Errorf("cavp: MCT L is %d bytes, want %d", hashSize, h.Size())
	}

	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	seed := slices.Clone(requestSeed)
	for count := range iterations {
		digest, err := hashMctRound(h, seed)
		if err != nil {
			return err
		}
		lines = append(lines, "", "COUNT = "+strconv.Itoa(count), "MD = "+upperHex(digest))
		seed = digest
	}

	response := strings.Join(lines, "\r\n")
	n, err := io.WriteString(w, response)
	if err == nil && n != len(response) {
		return io.ErrShortWrite
	}
	return err
}

func hashMctRound(h hash.Hash, seed []byte) ([]byte, error) {
	digestSize := h.Size()
	if digestSize <= 0 {
		return nil, fmt.Errorf("cavp: MCT requires a positive digest size")
	}
	if h.BlockSize() <= 0 {
		return nil, fmt.Errorf("cavp: MCT requires a positive block size")
	}
	if len(seed) != digestSize {
		return nil, fmt.Errorf("cavp: MCT seed is %d bytes, want %d", len(seed), digestSize)
	}

	n := h.BlockSize()/digestSize + 1
	history := make([]byte, (1000+n)*digestSize)
	for idx := range n {
		copy(history[idx*digestSize:], seed)
	}

	for idx := n; idx < 1000+n; idx++ {
		h.Reset()
		message := history[(idx-n)*digestSize : idx*digestSize]
		written, err := h.Write(message)
		if err != nil {
			return nil, err
		}
		if written != len(message) {
			return nil, io.ErrShortWrite
		}

		digest := h.Sum(nil)
		if len(digest) != digestSize {
			return nil, fmt.Errorf("cavp: hash returned %d bytes, want %d", len(digest), digestSize)
		}
		copy(history[idx*digestSize:], digest)
	}

	return slices.Clone(history[len(history)-digestSize:]), nil
}

func processHashSmtLmt(
	input, output string,
	newHash func() hash.Hash,
) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)

	response := make([]string, 0, 1024)
	record := make([]string, 0, 4)
	flush := func() error {
		if len(record) == 0 {
			return nil
		}

		messageLength := -1
		var message []byte
		haveLen, haveMsg := false, false
		for _, raw := range record {
			entry := parseLine(raw)
			switch strings.ToLower(entry.key) {
			case "len":
				messageLength, err = strconv.Atoi(entry.val)
				if err != nil {
					return fmt.Errorf("cavp: invalid Len: %w", err)
				}
				haveLen = true
			case "msg":
				message, err = hex.DecodeString(strings.Join(strings.Fields(entry.val), ""))
				if err != nil {
					return fmt.Errorf("cavp: invalid Msg: %w", err)
				}
				haveMsg = true
			}
		}

		// A header/instruction paragraph has neither record field and is copied
		// verbatim. A data record is evaluated only after its blank-line boundary,
		// so Len and Msg may appear in either order.
		if !haveLen && !haveMsg {
			response = append(response, record...)
			record = record[:0]
			return nil
		}
		if !haveLen || !haveMsg {
			return fmt.Errorf("cavp: incomplete hash record (Len=%t, Msg=%t)", haveLen, haveMsg)
		}
		if messageLength < 0 || messageLength%8 != 0 {
			return fmt.Errorf("cavp: Len must be a non-negative multiple of 8")
		}
		byteLength := messageLength / 8
		if len(message) < byteLength {
			return fmt.Errorf("cavp: Msg is %d bytes, want at least %d", len(message), byteLength)
		}

		h := newHash()
		if h == nil {
			return fmt.Errorf("cavp: hash constructor returned nil")
		}
		n, writeErr := h.Write(message[:byteLength])
		if writeErr != nil {
			return writeErr
		}
		if n != byteLength {
			return io.ErrShortWrite
		}
		response = append(response, record...)
		response = append(response, "MD = "+upperHex(h.Sum(nil)))
		record = record[:0]
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
		record = append(record, raw)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := flush(); err != nil {
		return err
	}
	for len(response) > 0 && strings.TrimSpace(response[len(response)-1]) == "" {
		response = response[:len(response)-1]
	}
	return os.WriteFile(output, []byte(strings.Join(response, "\r\n")), 0o644)
}

func TestHashSMT[F NewHashFunc](t testing.TB, filename string, fn F, size int) {
	t.Helper()

	newHash := normalizeHashFunction(fn, size)

	Test(t, filename, func(v HashSmtLmt) {
		if len(v.Digest) == 0 {
			return
		}

		testHashSMT(t, v.Message[:v.MessageLength/8], v.Digest, newHash)
	})
}

func testHashSMT(t testing.TB, input, want []byte, newHash func() hash.Hash) {
	t.Helper()

	h := newHash()
	n, err := h.Write(input)
	require.NoError(t, err)
	require.Equal(t, len(input), n, "short write")
	require.Equal(t, want, h.Sum(nil), "hash mismatch")
}

func TastHashSumSMTSize[T Digest](t testing.TB, filename string, fn func(size int, data []byte) T, size int) {
	t.Helper()

	TastHashSumSMT(t, filename, normalizeSumFunction[T](fn, size), size)
}

func TastHashSumSMT[T Digest](t testing.TB, filename string, sum func([]byte) T, size int) {
	t.Helper()

	Test(t, filename, func(v HashSmtLmt) {
		if len(v.Digest) == 0 {
			return
		}

		input := v.Message[:v.MessageLength/8]

		got := sum(input)
		gotB := make([]byte, len(got))
		for i := range gotB {
			gotB[i] = got[i]
		}
		require.GreaterOrEqual(t, len(gotB), size, "digest is shorter than expected")
		require.Equal(t, v.Digest, gotB[:size], "hash mismatch")
	})
}
