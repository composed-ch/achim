package achim

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

type User struct {
	Email  string `yaml:"email"`
	Name   string `yaml:"name"`
	SSHKey string `yaml:"ssh-key"`
}

type Group struct {
	Name  string `yaml:"name"`
	Users []User `yaml:"users"`
}

type NewGroupParams struct {
	GroupFile string
	Key       string
	Autostart bool
	Image     string
	Size      string
	Labels    string
}

func GroupFileFromText(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}
	name := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	var users []User
	for _, line := range strings.Split(string(data), "\n") {
		if user, ok := parseEmail(line); ok {
			users = append(users, user)
		}
	}
	slices.SortFunc(users, func(a, b User) int {
		return strings.Compare(a.Email, b.Email)
	})
	group := Group{Name: name, Users: users}
	out, err := yaml.Marshal(group)
	if err != nil {
		return fmt.Errorf("marshal group: %w", err)
	}
	if err := os.WriteFile(outputPath, out, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}
	return nil
}

func CreateGroup(ctx context.Context, params NewGroupParams) error {
	fmt.Println(params)
	return nil
}

func parseEmail(line string) (User, bool) {
	line = strings.ToLower(strings.TrimSpace(line))
	if line == "" {
		return User{}, false
	}
	username, _, found := strings.Cut(line, "@")
	if !found {
		return User{}, false
	}
	return User{Name: username, Email: line}, true
}
