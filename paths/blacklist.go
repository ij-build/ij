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

	err := runOnPatterns(allPatterns, project, func(path string) error {
		blacklist[path] = struct{}{}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return blacklist, nil
}
