package achim

import (
	"context"
	"fmt"
	"io"
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
	CloudInit string
}

func GroupFileFromText(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}
	name := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	var users []User
	for line := range strings.SplitSeq(string(data), "\n") {
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
	fmt.Printf("converted emails from %s to %s\n", inputPath, outputPath)
	return nil
}

func ParseGroupFile(yamlPath string) (*Group, error) {
	return ParseYAMLFile[Group](yamlPath)
}

func CreateGroup(ctx context.Context, params NewGroupParams) error {
	group, err := ParseGroupFile(params.GroupFile)
	if err != nil {
		return fmt.Errorf("parse group file at %s: %w", params.GroupFile, err)
	}
	names := make(map[string]User, len(group.Users))
	for _, u := range group.Users {
		names[u.Name] = u
	}
	newInstancesParam := NewInstanceGroupParam{
		Names:     names,
		Key:       params.Key,
		Autostart: params.Autostart,
		Image:     params.Image,
		Size:      params.Size,
		Labels:    fmt.Sprintf("%s,group=%s", params.Labels, group.Name),
		CloudInit: params.CloudInit,
	}
	newInstanceGroup, err := newInstancesParam.Compile(ctx)
	if err != nil {
		return fmt.Errorf("compile instance group %v: %w", params, err)
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

type Play struct {
	Become bool   `yaml:"become"`
	Hosts  string `yaml:"hosts"`
	Name   string `yaml:"name"`
	Tasks  []Task `yaml:"tasks"`
}

type Task struct {
	Name          string            `yaml:"name"`
	User          TaskUser          `yaml:"user,omitempty"`
	AuthorizedKey TaskAuthorizedKey `yaml:"authorized_key,omitempty"`
}

type TaskUser struct {
	Append     bool     `yaml:"append"`
	CreateHome bool     `yaml:"create_home"`
	Groups     []string `yaml:"groups"`
	Home       string   `yaml:"home"`
	Name       string   `yaml:"name"`
	Password   string   `yaml:"password"`
	Shell      string   `yaml:"shell"`
}
type TaskAuthorizedKey struct {
	Key  string `yaml:"key"`
	User string `yaml:"user"`
}

func ExportPlaybook(ctx context.Context, groupfile, playbookPath string) error {
	group, err := ParseGroupFile(groupfile)
	if err != nil {
		return fmt.Errorf("parse group file at %s: %w", groupfile, err)
	}
	f, err := os.Create(playbookPath)
	if err != nil {
		return fmt.Errorf(`create playbook file "%s": %w`, playbookPath, err)
	} else {
		defer f.Close()
	}
	playbook := make([]Play, len(group.Users))
	for i, user := range group.Users {
		play := Play{
			Become: true,
			Hosts:  user.Name,
			Name:   fmt.Sprintf("User Setup for %s", user.Name),
			Tasks: []Task{
				{
					Name: "User Created",
					User: TaskUser{
						Append:     true,
						CreateHome: true,
						Groups:     []string{"sudo"},
						Home:       "/home/user",
						Name:       "user",
						Password:   "*",
						Shell:      "/usr/bin/bash",
					},
				},
				{
					Name: "Key Authorized",
					AuthorizedKey: TaskAuthorizedKey{
						User: "user",
						Key:  user.SSHKey,
					},
				},
			},
		}
		playbook[i] = play
	}
	dump, err := yaml.Marshal(playbook)
	if err != nil {
		return fmt.Errorf(`marhsal playbook for group "%s": %w`, group.Name, err)
	}
	if _, err = io.Writer.Write(f, dump); err != nil {
		return fmt.Errorf(`write playbook to "%s": %w`, playbookPath, err)
	}
	fmt.Printf("created playbook based on group file %s under %s\n", groupfile, playbookPath)
	return nil
}
