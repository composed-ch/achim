package achim

import (
	"fmt"
	"os"
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
	path := fmt.Sprintf("%s%c%s", os.TempDir(), os.PathSeparator, filename)
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
	if err := AddTenant(first, path); err != nil {
		t.Fatalf("add tenant %v: %v\n", first, err)
	}
	defer mustCleanup(t, path)
	if err := AddTenant(second, path); err != nil {
		t.Fatalf("add tenant %v: %v\n", first, err)
	}
	firstRetrieved, err := GetTenant("1st", path)
	if err != nil {
		t.Fatalf("get tenant %s: %v\n", "1st", err)
	}
	if *firstRetrieved != first {
		t.Fatalf("checking tenant: expected %v, got %v\n", first, firstRetrieved)
	}
	secondRetrieved, err := GetTenant("2nd", path)
	if err != nil {
		t.Fatalf("get tenant %s: %v\n", "2nd", err)
	}
	if *secondRetrieved != second {
		t.Fatalf("checking tenant: expected %v, got %v\n", second, secondRetrieved)
	}
	if defaultTenant, err := GetDefaultTenant(path); err != nil {
		t.Fatalf("retrieve default tenant: %v\n", err)
	} else if *defaultTenant != first {
		t.Fatalf("retrieve default tenant: expected %v, got %v\n", first, defaultTenant)
	}
	if err := SetDefaultTenant("2nd", path); err != nil {
		t.Fatalf("set default tenant %s: %v\n", "2nd", err)
	}
	if defaultTenant, err := GetDefaultTenant(path); err != nil {
		t.Fatalf("retrieve default tenant: %v\n", err)
	} else if *defaultTenant != second {
		t.Fatalf("retrieve default tenant: expected %v, got %v\n", first, defaultTenant)
	}
}

func mustCleanup(t *testing.T, path string) {
	if err := os.Remove(path); err != nil {
		t.Fatalf("cleanup %s: %v\n", path, err)
	}
}
