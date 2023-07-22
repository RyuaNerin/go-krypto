package avoutil

import (
	"bufio"
	"regexp"
	"strings"

	. "github.com/mmcloughlin/avo/operand"
)

var reRemove = regexp.MustCompile(`  |\/\/|\t`)

func CheckType(rules string, args ...Op) {
	sb := strings.NewReader(rules)
	br := bufio.NewReader(sb)
	for {
		lineRaw, _, err := br.ReadLine()
		if err != nil {
			break
		}

		line := strings.TrimSpace(string(lineRaw))
		for reRemove.MatchString(line) {
			line = reRemove.ReplaceAllString(line, " ")
		}
		line = strings.TrimSpace(line)

		ss := strings.Split(line, " ")
		if len(ss)-1 != len(args) {
			continue
		}

		matched := true

		for idx, rule := range ss[1:] {
			if !checkType(rule, args[idx]) {
				matched = false
				break
			}
		}

		if matched {
			return
		}
	}

	panic("not matched\n" + rules)
}

func checkType(rule string, value Op) bool {
	switch rule {
	case "m256", "m128", "m64", "m32", "m16", "m8":
		return IsMem(value)

	case "ymm":
		return IsYMM(value)
	case "xmm":
		return IsXMM(value)

	case "r64":
		return IsR64(value)
	case "r32":
		return IsR32(value)
	case "r16":
		return IsR16(value)
	case "r8":
		return IsR8(value)

	case "k":

	case "imm64":
		_, ok1 := value.(U64)
		_, ok2 := value.(I64)
		return ok1 || ok2

	case "imm32":
		_, ok1 := value.(U32)
		_, ok2 := value.(I32)
		return ok1 || ok2

	case "imm16":
		_, ok1 := value.(U16)
		_, ok2 := value.(I16)
		return ok1 || ok2

	case "imm8":
		_, ok1 := value.(U8)
		_, ok2 := value.(I8)
		return ok1 || ok2
	}

	return false
}
