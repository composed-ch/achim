package achim

import "fmt"

type Tenant struct {
	Name   string
	Key    string
	Secret string
	Zone   string
}

func AddTenant(tenant Tenant, path string) error {
	fmt.Println("add", tenant, "to", path)
	return nil
}
