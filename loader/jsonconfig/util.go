package jsonconfig

import (
	"encoding/json"
	"fmt"
	"time"
)

type Duration struct {
	time.Duration
}

var ZeroDuration = Duration{time.Duration(0)}

func (d Duration) String() string {
	if d.Duration == 0 {
		return ""
	}

	return d.Duration.String()
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	val := ""
	err := json.Unmarshal(data, &val)

	if err == nil {
		parsed, err := time.ParseDuration(val)
		if err != nil {
			return err
		}

		*d = Duration{parsed}
		return nil
	}

	return fmt.Errorf(
		"%s is not a valid duration",
		string(data),
	)
}

func unmarshalStringList(raw json.RawMessage) ([]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	single := ""
	if err := json.Unmarshal(raw, &single); err == nil {
		return []string{single}, nil
	}

	multiple := []string{}
	if err := json.Unmarshal(raw, &multiple); err != nil {
		return nil, err
	}

	return multiple, nil
}
