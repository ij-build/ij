package util

func ContainsAll(sub, sup []string) (bool, []string) {
	hash := map[string]struct{}{}
	for _, val := range sub {
		hash[val] = struct{}{}
	}

	missing := []string{}
	for _, val := range sup {
		if _, ok := hash[val]; !ok {
			missing = append(missing, val)
		}
	}

	return len(missing) == 0, missing
}
