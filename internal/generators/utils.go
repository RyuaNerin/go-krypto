package generators

import (
	"bytes"
	"go/format"
	"io"
	"os"
)

func WriteBody(
	filename string,
	prefix, suffix string,
	fn func(io.Writer),
) {
	var buf bytes.Buffer

	buf.WriteString(prefix)
	fn(&buf)
	buf.WriteString(suffix)

	b, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filename, b, 0o644)
	if err != nil {
		panic(err)
	}
}
