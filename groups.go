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

func ParseGroupFile(yamlPath string) (*Group, error) {
	raw, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf(`read YAML from file "%s": %w`, yamlPath, err)
	}
	var group Group
	if err = yaml.Unmarshal(raw, &group); err != nil {
		return nil, fmt.Errorf(`parse YAML content \n%s\n: %w`, string(raw), err)
	}
	return &group, nil
}

func CreateGroup(ctx context.Context, params NewGroupParams) error {
	group, err := ParseGroupFile(params.GroupFile)
	if err != nil {
		return fmt.Errorf("parse group file at %s: %w", params.GroupFile, err)
	}
	names := make([]string, len(group.Users))
	for i, u := range group.Users {
		names[i] = u.Name
	}
	newInstancesParam := NewInstancesParam{
		Names:     names,
		Key:       params.Key,
		Autostart: params.Autostart,
		Image:     params.Image,
		Size:      params.Size,
		Labels:    fmt.Sprintf("%s,group=%s", params.Labels, group.Name),
	}
	newInstanceGroup, err := newInstancesParam.Compile(ctx)
	if err != nil {
		return fmt.Errorf("comile instance group %v: %w", params, err)
	}
	return newInstanceGroup.Create(ctx)
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
