package testingutil

import (
	"os"
	"strings"
)

func init() {
	// 명령 라인 인수를 읽음
	args := os.Args

	// "go test" 명령을 실행하지 않았을 때 경고 메시지 출력
	if len(args) < 2 || !strings.Contains(args[1], "test") {
		print("warning: testingutil package is only for go test.")
	}
}
