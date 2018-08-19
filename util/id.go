package util

import "github.com/google/uuid"

func MakeID() (string, error) {
	raw, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return raw.String()[:6], nil
}
