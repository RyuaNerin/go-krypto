package cavp

import (
	"compress/gzip"
	"embed"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed testdata
var testdata embed.FS

func Test[T any](TB testing.TB, filename string, fn func(T)) {
	TB.Helper()

	filename, _ = strings.CutSuffix(filename, ".gz")
	filename, _ = strings.CutSuffix(filename, ".txt")

	fs, err := testdata.Open("testdata/" + filename + ".txt.gz")
	var input io.Reader
	if err == nil {
		defer fs.Close()
		gzipReader, err := gzip.NewReader(fs)
		require.NoError(TB, err)
		defer gzipReader.Close()
		input = gzipReader
	} else {
		fs, err = testdata.Open("testdata/" + filename + ".txt")
		require.NoError(TB, err)
		defer fs.Close()
		input = fs
	}

	tr := NewReader[T](input)

	for {
		f, err := tr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(TB, err)
		fn(f)
	}
}
