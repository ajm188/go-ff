package feature

import (
	"encoding/json"
	"io/ioutil"
)

func Init(m map[string]bool) {
	inst.m.Lock()
	defer inst.m.Unlock()

	inst.features = make(map[string]bool, len(m))

	for k, v := range m {
		inst.features[k] = v
	}
}

func InitFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	inst.m.Lock()
	defer inst.m.Unlock()

	return json.Unmarshal(data, &inst.features)
}
