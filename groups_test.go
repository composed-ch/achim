package achim

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestGroupFileFromText(t *testing.T) {
	input := filepath.Join(t.TempDir(), "students.txt")
	output := filepath.Join(t.TempDir(), "students.yaml")

	if err := os.WriteFile(input, []byte("terence_hill@composed.ch\nbud_spencer@composed.ch\n"), 0644); err != nil {
		t.Fatalf("write input file: %v", err)
	}
	if err := GroupFileFromText(input, output); err != nil {
		t.Fatalf("GroupFileFromText: %v", err)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}
	var group Group
	if err := yaml.Unmarshal(data, &group); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}
	if group.Name != "students" {
		t.Fatalf("expected group name 'students', got '%s'", group.Name)
	}
	if len(group.Users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(group.Users))
	}
	if group.Users[0].Name != "bud_spencer" {
		t.Fatalf("expected first user 'bud_spencer' (sorted), got '%s'", group.Users[0].Name)
	}
}

func TestGroupFileFromTextBlankLines(t *testing.T) {
	input := filepath.Join(t.TempDir(), "group.txt")
	output := filepath.Join(t.TempDir(), "group.yaml")

	if err := os.WriteFile(input, []byte("bud_spencer@composed.ch\n\nterence_hill@composed.ch\n"), 0644); err != nil {
		t.Fatalf("write input file: %v", err)
	}
	if err := GroupFileFromText(input, output); err != nil {
		t.Fatalf("GroupFileFromText: %v", err)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}
	var group Group
	if err := yaml.Unmarshal(data, &group); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}
	if len(group.Users) != 2 {
		t.Fatalf("expected 2 users (blank lines ignored), got %d", len(group.Users))
	}
}
