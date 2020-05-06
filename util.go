package main

func contains(slice []string, test string) bool {
	for _, item := range slice {
		if item == test {
			return true
		}
	}
	return false
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}
	return 0
}
