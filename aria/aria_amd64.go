//go:build amd64 && gc && !purego

package aria

func init() {
	processFin = processFinSSE2
}
