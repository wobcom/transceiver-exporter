package transceivercollector

import (
	"math"
	"regexp"
	"strconv"
)

func contains(l []string, test string) bool {
	for _, item := range l {
		match, _ := regexp.MatchString(item, test)
		if match {
			return true
		}
	}
	return false
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func milliwattsToDbm(mw float64) float64 {
	return 10 * math.Log10(mw)
}

func quoteToascii(str string) string {
	textQuoted := strconv.QuoteToASCII(str)
	textUnquoted := textQuoted[1 : len(textQuoted)-1]
	return textUnquoted
}
