package environment

func Merge(environments ...Environment) Environment {
	target := Environment{}
	for _, env := range environments {
		for k, v := range env {
			target[k] = v
		}
	}

	return target
}
