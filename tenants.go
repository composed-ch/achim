package achim

import (
	"encoding/json"
	"fmt"
	"os"
)

const ConfigFilePerm = 0600

type TenantFile struct {
	Tenants map[string]Tenant `json:"tenants"`
}

type Tenant struct {
	Name   string `json:"name"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
	Zone   string `json:"zone"`
}

// TODO: set tenant as default if it's the first one
// TODO: new method GetDefaultTenant
// TODO: implement operation to set default tenant

func AddTenant(tenant Tenant, path string) error {
	existing, err := ReadTenants(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "no tenants under %s, creating new tenants file\n", path)
		existing = &TenantFile{Tenants: make(map[string]Tenant)}
	}
	if _, ok := existing.Tenants[tenant.Name]; ok {
		return fmt.Errorf("tenant %s already exists", tenant.Name)
	}
	fmt.Println(existing)
	existing.Tenants[tenant.Name] = tenant
	data, err := json.MarshalIndent(existing, "", "    ")
	if err != nil {
		return fmt.Errorf("convert %v to JSON: %v", tenant, err)
	}
	os.WriteFile(path, data, ConfigFilePerm)
	return nil
}

func GetTenant(name, path string) (*Tenant, error) {
	tenantsFile, err := ReadTenants(path)
	if err != nil {
		return nil, fmt.Errorf("read tenants from %s: %v", path, err)
	}
	if tenant, ok := tenantsFile.Tenants[name]; !ok {
		return nil, fmt.Errorf("no such tenant %s in tenant file %s", name, path)
	} else {
		return &tenant, nil
	}
}

func ReadTenants(path string) (*TenantFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read tenants from %s: %v", path, err)
	}
	var tenantsFile TenantFile
	if err := json.Unmarshal(data, &tenantsFile); err != nil {
		return nil, fmt.Errorf("unmarshall tenant file %s: %v", path, err)
	}
	if tenantsFile.Tenants == nil {
		tenantsFile.Tenants = make(map[string]Tenant)
	}
	return &tenantsFile, nil
}
