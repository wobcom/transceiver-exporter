package transceivercollector

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
