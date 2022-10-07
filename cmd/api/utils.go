package main

const (
	// add more environments
	// handle differences between envs
	develop = "develop"
)

func sliceContains(a []string, b string) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}

// compare if each element of 1 slice is present in another slice
func compareSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// this works just fine for smaller slices
	// otherwise this would require different solution due to the overhead
	for _, j := range a {
		if !sliceContains(b, j) {
			return false
		}
	}

	return true
}
