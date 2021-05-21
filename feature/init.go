package feature

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

var ErrEmptyConfig = errors.New("empty config file")

func Init(m map[string]*Feature) {
	inst.m.Lock()
	defer inst.m.Unlock()

	inst.features = make(map[string]*Feature, len(m))

	for k, v := range m {
		inst.features[k] = v
	}
}

func InitFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return fmt.Errorf("%w %s", ErrEmptyConfig, path)
	}

	var m map[string]*Feature
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	for name, f := range m {
		f.Name = name

		ok, err := f.Validate()
		if err != nil {
			return fmt.Errorf("feature %s: %w", f.Name, err)
		}

		if !ok {
			return fmt.Errorf("%w for %s", ErrInvalidFeature, err)
		}
	}

	Init(m)
	return nil
}
