package util

// ContainsAll determines if sup is a superset of sub. The
// first return value is a flag, and the second is the set
// of values which are missing from sup.
func ContainsAll(sup, sub []string) (bool, []string) {
	hash := map[string]struct{}{}
	for _, val := range sup {
		hash[val] = struct{}{}
	}

	missing := []string{}
	for _, val := range sub {
		if _, ok := hash[val]; !ok {
			missing = append(missing, val)
		}
	}

	return len(missing) == 0, missing
}
