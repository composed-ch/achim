package achim

import (
	"encoding/json"
	"fmt"
	"os"
)

const ConfigFilePerm = 0600

type Tenant struct {
	Name   string `json:"name"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
	Zone   string `json:"zone"`
}

func AddTenant(tenant Tenant, path string) error {
	data, err := json.MarshalIndent(tenant, "", "    ")
	if err != nil {
		return fmt.Errorf("convert %v to JSON: %v", tenant, err)
	}
	os.WriteFile(path, data, ConfigFilePerm)
	return nil
}
