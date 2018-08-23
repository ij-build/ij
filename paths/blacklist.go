package paths

var DefaultBlacklist = []string{
	".ij",
	".git",
}

func constructBlacklist(project string, patterns []string) (map[string]struct{}, error) {
	var (
		blacklist   = map[string]struct{}{}
		allPatterns = append(DefaultBlacklist, patterns...)
	)

	err := runOnPatterns(allPatterns, project, false, func(pair filePair) error {
		blacklist[pair.src] = struct{}{}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return blacklist, nil
}
