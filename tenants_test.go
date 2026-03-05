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
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v\n", path, err)
	}
	if info.Mode().Perm() != perm {
		t.Fatalf("%s expected mode %o, got %o\n", path, perm, info.Mode().Perm())
	}
}
