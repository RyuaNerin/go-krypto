package cavp

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/RyuaNerin/go-krypto/hmac"
	"github.com/stretchr/testify/require"
)

// HMACKAT is a hash-based message authentication code CAVP record.
// TagLength is expressed in bytes by the Korean CAVP format.
type HMACKAT struct {
	HashSize  int    `cavp:"L,int,header=plain"`
	Count     int    `cavp:"COUNT,int"`
	KeyLength int    `cavp:"Klen,int"`
	TagLength int    `cavp:"Tlen,int"`
	Key       []byte `cavp:"Key,hex"`
	Message   []byte `cavp:"Msg,hex"`
	MAC       []byte `cavp:"Mac,hex,optional"`
}

// ProcessHMAC reads an HMAC request file and writes its response file.
func ProcessHMAC(input, output string, newHash func() hash.Hash) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)

	response := make([]string, 0, 1024)
	record := make([]string, 0, 8)
	flush := func() error {
		if len(record) == 0 {
			return nil
		}
		var key, message []byte
		keyLength, tagLength := -1, -1
		haveKey, haveMessage, haveKeyLength, haveTagLength := false, false, false, false
		isRecord := false
		for _, raw := range record {
			entry := parseLine(raw)
			switch strings.ToLower(entry.key) {
			case "count":
				isRecord = true
			case "klen":
				isRecord = true
				keyLength, err = strconv.Atoi(entry.val)
				if err != nil {
					return fmt.Errorf("cavp: invalid Klen: %w", err)
				}
				haveKeyLength = true
			case "tlen":
				isRecord = true
				tagLength, err = strconv.Atoi(entry.val)
				if err != nil {
					return fmt.Errorf("cavp: invalid Tlen: %w", err)
				}
				haveTagLength = true
			case "key":
				isRecord = true
				key, err = hex.DecodeString(strings.Join(strings.Fields(entry.val), ""))
				if err != nil {
					return fmt.Errorf("cavp: invalid Key: %w", err)
				}
				haveKey = true
			case "msg":
				isRecord = true
				message, err = hex.DecodeString(strings.Join(strings.Fields(entry.val), ""))
				if err != nil {
					return fmt.Errorf("cavp: invalid Msg: %w", err)
				}
				haveMessage = true
			}
		}
		if !isRecord {
			response = append(response, record...)
			record = record[:0]
			return nil
		}
		if !haveKey || !haveMessage || !haveKeyLength || !haveTagLength {
			return fmt.Errorf("cavp: incomplete HMAC record (Klen=%t, Key=%t, Msg=%t, Tlen=%t)", haveKeyLength, haveKey, haveMessage, haveTagLength)
		}
		if keyLength < 0 || len(key) != keyLength {
			return fmt.Errorf("cavp: Key is %d bytes, want Klen %d", len(key), keyLength)
		}

		h := hmac.NewWithSize(newHash, key, tagLength)
		n, writeErr := h.Write(message)
		if writeErr != nil {
			return writeErr
		}
		if n != len(message) {
			return io.ErrShortWrite
		}
		sum := h.Sum(nil)
		response = append(response, record...)
		response = append(response, "Mac = "+upperHex(sum[:tagLength]))
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

// TestHMAC checks the MAC in every HMAC response record.
func TestHMAC(t testing.TB, input string, newHash func() hash.Hash) {
	t.Helper()
	Test(t, input, func(v HMACKAT) {
		h := hmac.NewWithSize(newHash, v.Key, v.TagLength)
		_, err := h.Write(v.Message)
		require.NoError(t, err)

		got := h.Sum(nil)
		require.Equalf(t, v.MAC, got, "HMAC mismatch for COUNT=%d", v.Count)
	})
}
