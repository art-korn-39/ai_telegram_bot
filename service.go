package main

import (
	"math"
	"slices"
	"sort"
	"strings"
	"unicode"

	"github.com/sashabaranov/go-openai"
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

// func Tokenizer(s string) int {

// 	runes := []rune(s)
// 	for _, r := range runes {
// //		utf8.
// 	}

// }

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

func gptGetVoice(voice string) (v openai.SpeechVoice, isError bool) {

	array := []openai.SpeechVoice{openai.VoiceAlloy, openai.VoiceEcho, openai.VoiceFable,
		openai.VoiceOnyx, openai.VoiceNova, openai.VoiceShimmer}

	SV := openai.SpeechVoice(strings.ToLower(voice))

	i := slices.Index(array, SV)
	if i == -1 {
		return "", true
	} else {
		return array[i], false
	}

}

// func SpecialCommand(cmd string) bool {

// 	if slices.Contains(SpecialCMD, strings.ToLower(cmd)) {
// 		return true
// 	}

// 	return false

// }
