// Package utils contains functions for the main maskcat program
package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// MakeMask performs substitution to make HC masks
func MakeMask(str string) string {
	var args []string
	for c := 'a'; c <= 'z'; c++ {
		args = append(args, string(c), "?l")
	}
	for c := 'A'; c <= 'Z'; c++ {
		args = append(args, string(c), "?u")
	}
	for c := '0'; c <= '9'; c++ {
		args = append(args, string(c), "?d")
	}
	specialChars := " !\"#$%&\\()*+,-./:;<=>?@[\\]^_`{|}~"
	for _, c := range specialChars {
		args = append(args, string(c), "?s")
	}
	replacer := strings.NewReplacer(args...)
	return replacer.Replace(str)
}

// MakeToken replaces all non-alpha characters to generate tokens
func MakeToken(str string) string {
	re := regexp.MustCompile(`[^a-zA-Z]+`)
	return re.ReplaceAllString(str, "")
}

// MakePartialMask creates a partial Hashcat mask
func MakePartialMask(str string, chars string) string {
	var lowerArgs, upperArgs, digitArgs []string
	for c := 'a'; c <= 'z'; c++ {
		lowerArgs = append(lowerArgs, string(c), "?l")
	}
	for c := 'A'; c <= 'Z'; c++ {
		upperArgs = append(upperArgs, string(c), "?u")
	}
	for c := '0'; c <= '9'; c++ {
		digitArgs = append(digitArgs, string(c), "?d")
	}
	specialChars := " !\"#$%&\\()*+,-./:;<=>?@[\\]^_`{|}~"
	specialArgs := make([]string, len(specialChars)*2)
	for i, c := range specialChars {
		specialArgs[i*2] = string(c)
		specialArgs[i*2+1] = "?s"
	}

	if strings.Contains(chars, "l") {
		str = strings.NewReplacer(lowerArgs...).Replace(str)
	}

	if strings.Contains(chars, "u") {
		str = strings.NewReplacer(upperArgs...).Replace(str)
	}

	if strings.Contains(chars, "d") {
		str = strings.NewReplacer(digitArgs...).Replace(str)
	}

	if strings.Contains(chars, "s") {
		if strings.Contains(chars, "u") || strings.Contains(chars, "l") || strings.Contains(chars, "d") {
			str = ""
		} else {
			str = strings.NewReplacer(specialArgs...).Replace(str)
		}
	}

	return str
}

// TestComplexity tests the complexity of an input string
func TestComplexity(str string) int {
	complexity := 0
	charTypes := []string{"?u", "?l", "?d", "?s"}
	for _, charType := range charTypes {
		if strings.Contains(str, charType) {
			complexity++
		}
	}
	return complexity
}

// TestEntropy calculates mask entropy
func TestEntropy(str string) int {
	entropy := 0
	charTypes := []struct {
		charType string
		count    int
	}{
		{"?u", 26},
		{"?l", 26},
		{"?d", 10},
		{"?s", 33},
	}
	for _, ct := range charTypes {
		entropy += strings.Count(str, ct.charType) * ct.count
	}
	return entropy
}

// ChunkString splits string into chunks
func ChunkString(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

// RemoveDuplicateStr removes duplicate strings from array
func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := make([]string, 0, len(strSlice))
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// ReplaceAtIndex replaces a rune at index in string
func ReplaceAtIndex(in string, r rune, i int) string {
	if i < 0 || i >= len(in) {
		CheckError(fmt.Errorf("index out of range"))
	}
	out := []rune(in)
	if i >= 0 && i < len(out) {
		out[i] = r
		// In instances where i is out of bounds go to the end
	} else if i >= 0 && i == len(out) {
		out[len(out)-1] = r
	}
	return string(out)
}

// ReplaceWord replaces a mask within an input string with a value
func ReplaceWord(stringword, mask string, value string) string {
	tokenmask := MakeMask(value)
	if strings.Contains(mask, tokenmask) {
		newword := strings.Replace(mask, tokenmask, value, -1)
		newword = strings.ReplaceAll(newword, "?u", "?")
		newword = strings.ReplaceAll(newword, "?l", "?")
		newword = strings.ReplaceAll(newword, "?d", "?")
		newword = strings.ReplaceAll(newword, "?s", "?")

		for i := range stringword {
			if i < len(newword) {

				if newword[i] == '?' {
					newword = ReplaceAtIndex(newword, rune(stringword[i]), i)
				}
			}
		}

		if strings.Contains(newword, value) && newword != value {
			return newword
		}
	}
	return ""
}

// CheckError is a general error handler
func CheckError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(0)
	}
}
