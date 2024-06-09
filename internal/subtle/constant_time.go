package subtle

import (
	go_subtle "crypto/subtle"
)

func ConstantTimeCompare(x, y []byte) int { return go_subtle.ConstantTimeCompare(x, y) }
func ConstantTimeSelect(v, x, y int) int  { return go_subtle.ConstantTimeSelect(v, x, y) }
func ConstantTimeByteEq(x, y uint8) int   { return go_subtle.ConstantTimeByteEq(x, y) }
func ConstantTimeEq(x, y int32) int       { return go_subtle.ConstantTimeEq(x, y) }
func ConstantTimeCopy(v int, x, y []byte) { go_subtle.ConstantTimeCopy(v, x, y) }
func ConstantTimeLessOrEq(x, y int) int   { return go_subtle.ConstantTimeLessOrEq(x, y) }
