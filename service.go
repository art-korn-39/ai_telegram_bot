package main

import (
	"math"
	"sort"
	"time"
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

// С 0 (вкл.) до last (искл.)
func SubString(s string, first int, last int) string {

	runes := []rune(s)
	length := len(runes)

	if length <= last {
		last = length
	}

	return string(runes[first:last])

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

func SortMap(m map[int]string) (result map[int]string) {

	result = map[int]string{}

	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for _, k := range keys {
		result[k] = m[k]
	}

	return result

}

func MskTimeNow() time.Time {

	return time.Now().UTC().Add(3 * time.Hour)

}
