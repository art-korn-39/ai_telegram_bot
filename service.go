package main

import (
	"math"
	"unicode"
)

func Unused(...any) {}

func Ternary(statement bool, a any, b any) any {
	if statement {
		return a
	}
	return b
}

func FR(v any, err error) any {
	if err != nil {
		panic("error encountered when none assumed:" + err.Error())
	}
	return v
}

func Round(x float64, decimals float64) float64 {

	multiplier := math.Pow(10, decimals)
	result := math.Round(x*multiplier) / multiplier
	return result

}

func IsEngByLoop(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func IsRusByUnicode(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Cyrillic, r) {
			return true
		}
	}
	return false
}
