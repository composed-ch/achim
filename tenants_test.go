package achim

import (
	"fmt"
	"os"
	"path"
	"testing"
)

func TestAddFirstTenant(t *testing.T) {
	const filename = "achim.json"
	const perm = 0600
	path := fmt.Sprintf("%s%c%s", os.TempDir(), os.PathSeparator, filename)
	tenant := Tenant{
		Name:   "goachim",
		Key:    "EXO-123-456",
		Secret: "0123-4567-89ab-cdef",
		Zone:   "ch-xy1",
	}
	if err := AddTenant(tenant, path); err != nil {
		t.Fatalf("add tenant %v: %v\n", tenant, err)
	}
	defer mustCleanup(t, path)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v\n", path, err)
	}
	if info.Mode().Perm() != perm {
		t.Fatalf("%s expected mode %o, got %o\n", path, perm, info.Mode().Perm())
	}
	tenantRetrieved, err := GetTenant("goachim", path)
	if err != nil {
		t.Fatalf("retrieve tenant %s: %v", "goachim", err)
	}
	if *tenantRetrieved != tenant {
		t.Fatalf("checking tenant: expected %v, got %v\n", tenant, tenantRetrieved)
	}
	if defaultTenant, err := GetDefaultTenant(path); err != nil {
		t.Fatalf("retrieve default tenant: %v\n", err)
	} else if *defaultTenant != tenant {
		t.Fatalf("retrieve default tenant: expected %v, got %v\n", tenant, defaultTenant)
	}
}

func TestAddTwoTenants(t *testing.T) {
	const filename = "achim.json"
	filepath := path.Join(os.TempDir(), filename)
	first := Tenant{
		Name:   "1st",
		Key:    "EXO-123",
		Secret: "0123-4567",
		Zone:   "ch-xy1",
	}
	second := Tenant{
		Name:   "2nd",
		Key:    "EXO-456",
		Secret: "89ab-cdef",
		Zone:   "ch-xy2",
	}
	if err := AddTenant(first, filepath); err != nil {
		t.Fatalf("add tenant %v: %v\n", first, err)
	}
	defer mustCleanup(t, filepath)
	if err := AddTenant(second, filepath); err != nil {
		t.Fatalf("add tenant %v: %v\n", first, err)
	}
	firstRetrieved, err := GetTenant("1st", filepath)
	if err != nil {
		t.Fatalf("get tenant %s: %v\n", "1st", err)
	}
	if *firstRetrieved != first {
		t.Fatalf("checking tenant: expected %v, got %v\n", first, firstRetrieved)
	}
	secondRetrieved, err := GetTenant("2nd", filepath)
	if err != nil {
		t.Fatalf("get tenant %s: %v\n", "2nd", err)
	}
	if *secondRetrieved != second {
		t.Fatalf("checking tenant: expected %v, got %v\n", second, secondRetrieved)
	}
	if defaultTenant, err := GetDefaultTenant(filepath); err != nil {
		t.Fatalf("retrieve default tenant: %v\n", err)
	} else if *defaultTenant != first {
		t.Fatalf("retrieve default tenant: expected %v, got %v\n", first, defaultTenant)
	}
	if err := SetDefaultTenant("2nd", filepath); err != nil {
		t.Fatalf("set default tenant %s: %v\n", "2nd", err)
	}
	if defaultTenant, err := GetDefaultTenant(filepath); err != nil {
		t.Fatalf("retrieve default tenant: %v\n", err)
	} else if *defaultTenant != second {
		t.Fatalf("retrieve default tenant: expected %v, got %v\n", first, defaultTenant)
	}
}

func TestAddEmptyTenant(t *testing.T) {
	const filename = "achim.json"
	filepath := path.Join(os.TempDir(), filename)
	if err := AddTenant(Tenant{}, filepath); err == nil {
		defer mustCleanup(t, filepath)
		t.Fatalf("added empty tenant must not work\n")
	}
}

func TestAddAndRemoveTenants(t *testing.T) {
	const filename = "achim.json"
	filepath := path.Join(os.TempDir(), filename)
	one := Tenant{
		Name:   "one",
		Key:    "top",
		Secret: "secret",
		Zone:   "xy-z1",
	}
	two := Tenant{
		Name:   "two",
		Key:    "most",
		Secret: "secretly",
		Zone:   "ab-c2",
	}
	if err := AddTenant(one, filepath); err != nil {
		t.Fatalf("add tenant %v: %v\n", one, err)
	}
	defer mustCleanup(t, filepath)
	if err := AddTenant(two, filepath); err != nil {
		t.Fatalf("add tenant %v: %v\n", two, err)
	}
	if err := RemoveTenant(one.Name, filepath); err != nil {
		t.Fatalf("remove tenant %s: %v\n", one.Name, err)
	}
	tenant, err := GetDefaultTenant(filepath)
	if err != nil {
		t.Fatalf("get default tenant: %v\n", err)
	}
	if tenant.Name == one.Name {
		t.Fatalf("default tenant %s was removed", tenant.Name)
	}
	if err := RemoveTenant("phony", filepath); err == nil {
		t.Fatalf("removing phony tenant is supposed to fail\n")
	}
	tenants, err := ReadTenants(filepath)
	if err != nil {
		t.Fatalf("read tenants: %v", err)
	}
	fmt.Println(tenants)
	if len(tenants.Tenants) != 1 {
		t.Fatalf("expected one remaining tenant, has %d\n", len(tenants.Tenants))
	}
	if tenants.Default != two.Name {
		t.Fatalf("expected default tenant to be %s, was %s\n", two.Name, tenants.Default)
	}
	if err := RemoveTenant(two.Name, filepath); err != nil {
		t.Fatalf("remove tenant %s: %v\n", two.Name, err)
	}
}

func mustCleanup(t *testing.T, filepath string) {
	if err := os.Remove(filepath); err != nil {
		t.Fatalf("cleanup %s: %v\n", filepath, err)
	}
}
