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

func Test[T any](tb testing.TB, filename string, fn func(T)) {
	tb.Helper()

	filename, _ = strings.CutSuffix(filename, ".gz")
	filename, _ = strings.CutSuffix(filename, ".txt")

	fs, err := testdata.Open("testdata/" + filename + ".txt.gz")
	var input io.Reader
	if err == nil {
		defer fs.Close()
		gzipReader, err := gzip.NewReader(fs)
		require.NoError(tb, err)
		defer gzipReader.Close()
		input = gzipReader
	} else {
		fs, err = testdata.Open("testdata/" + filename + ".txt")
		require.NoError(tb, err)
		defer fs.Close()
		input = fs
	}

	tr := NewReader[T](input)

	for {
		f, err := tr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(tb, err)
		fn(f)
	}
}
