package templates

import (
	"fmt"
	"os"
	"text/template"
)

var templatesByName = map[string]string{
	"overview": Overview,
	"scenario": Scenario,
}

func Output[T any](data T, name, path string) error {
	t, ok := templatesByName[name]
	if !ok {
		return fmt.Errorf("no such template %s", name)
	}
	tmpl, err := template.New(name).Parse(t)
	if err != nil {
		return fmt.Errorf("parse overview template: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file %s: %w", path, err)
	}
	defer f.Close()
	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("execute %s template: %w", name, err)
	}
	fmt.Printf("wrote %s to %s\n", name, path)
	return nil
}
