//go:build amd64 && !purego
// +build amd64,!purego

package aria

func init() {
	processFin = processFinSSE2
}
