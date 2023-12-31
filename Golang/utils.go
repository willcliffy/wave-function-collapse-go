package main

func StringSliceContains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func sortFloat64Slice(slice []float64) {
	for i := 0; i < len(slice)-1; i++ {
		for j := i + 1; j < len(slice); j++ {
			if slice[j] < slice[i] {
				slice[i], slice[j] = slice[j], slice[i]
			}
		}
	}
}

func duplicateMap(originalMap map[string]WFCPrototype) map[string]WFCPrototype {
	duplicate := make(map[string]WFCPrototype, len(originalMap))
	for key, value := range originalMap {
		duplicate[key] = value
	}
	return duplicate
}
