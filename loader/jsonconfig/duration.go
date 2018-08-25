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
