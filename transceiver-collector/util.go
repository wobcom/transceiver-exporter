package transceivercollector

import "math"

func contains(l []string, test string) bool {
	for _, item := range l {
		if item == test {
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
