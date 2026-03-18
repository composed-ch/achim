package achim

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const ConfigFilePerm = 0600

type TenantFile struct {
	Default string            `json:"default"`
	Tenants map[string]Tenant `json:"tenants"`
}

type Tenant struct {
	Name   string `json:"name"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
	Zone   string `json:"zone"`
}

func (t Tenant) IsNonEmpty() bool {
	return strings.TrimSpace(t.Name) != "" &&
		strings.TrimSpace(t.Key) != "" &&
		strings.TrimSpace(t.Secret) != "" &&
		strings.TrimSpace(t.Zone) != ""
}

func AddTenant(tenant Tenant, path string) error {
	if !tenant.IsNonEmpty() {
		return fmt.Errorf("tenant %v lacks mandatory information", tenant)
	}
	existing, err := ReadTenants(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "no tenants under %s, creating new tenants file\n", path)
		existing = &TenantFile{
			Tenants: map[string]Tenant{tenant.Name: tenant},
			Default: tenant.Name,
		}
	} else if _, ok := existing.Tenants[tenant.Name]; ok {
		return fmt.Errorf("tenant %s already exists", tenant.Name)
	} else {
		existing.Tenants[tenant.Name] = tenant
	}
	return saveTenants(existing, path)
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

func GetDefaultTenant(path string) (*Tenant, error) {
	tenants, err := ReadTenants(path)
	if err != nil {
		return nil, fmt.Errorf("read tenants: %v", err)
	}
	if tenant, ok := tenants.Tenants[tenants.Default]; !ok {
		return nil, fmt.Errorf("retrieve default tenant %s from tenants: %v",
			tenants.Default, tenants)
	} else {
		return &tenant, nil
	}
}

func SetDefaultTenant(name, path string) error {
	tenants, err := ReadTenants(path)
	if err != nil {
		return fmt.Errorf("read tenants from %s: %v", path, err)
	}
	if _, ok := tenants.Tenants[name]; !ok {
		return fmt.Errorf("no tenant %s in %s", name, path)
	} else if tenants.Default != name {
		tenants.Default = name
		return saveTenants(tenants, path)
	}
	return nil
}

func RemoveTenant(name, path string) error {
	tenants, err := ReadTenants(path)
	if err != nil {
		return fmt.Errorf("read tenants from %s: %v", path, err)
	}
	if _, ok := tenants.Tenants[name]; !ok {
		return fmt.Errorf("no tenants %s in %s", name, path)
	}
	delete(tenants.Tenants, name)
	if tenants.Default == name {
		tenants.Default = ""
		for k := range tenants.Tenants {
			tenants.Default = k
			break
		}
	}
	saveTenants(tenants, path)
	return nil
}

func saveTenants(tenants *TenantFile, path string) error {
	data, err := json.MarshalIndent(tenants, "", "    ")
	if err != nil {
		return fmt.Errorf("convert %v to JSON: %v", tenants, err)
	}
	return os.WriteFile(path, data, ConfigFilePerm)
}
