package achim

import (
	"context"
	"fmt"
	"slices"
	"strings"

	v3 "github.com/exoscale/egoscale/v3"
)

func ListImages(ctx context.Context, contains string) ([]string, error) {
	exo := ctx.Value("exo").(*v3.Client)
	res, err := exo.ListTemplates(ctx)
	if err != nil {
		return nil, fmt.Errorf("list templates: %w", err)
	}
	images := make([]string, 0)
	for _, template := range res.Templates {
		if contains == "" || strings.Contains(strings.ToLower(template.Name), contains) {
			images = append(images, template.Name)
		}
	}
	slices.Sort(images)
	return images, nil
}

func GetTemplateByName(ctx context.Context, name string) (*v3.Template, error) {
	exo := ctx.Value("exo").(*v3.Client)
	res, err := exo.ListTemplates(ctx)
	if err != nil {
		return nil, fmt.Errorf("list images: %w", err)
	}
	for _, template := range res.Templates {
		if template.Name == name {
			return &template, nil
		}
	}
	return nil, fmt.Errorf(`no template for image name "%s" found`, name)
}

func GetAllowedSizes(ctx context.Context) (map[string]*v3.InstanceType, error) {
	sizes := make(map[string]*v3.InstanceType)
	types, err := ListInstanceTypes(ctx, "standard")
	if err != nil {
		return nil, fmt.Errorf("list standard intance types: %w", err)
	}
	for _, t := range types {
		sizes[string(t.Size)] = &t
	}
	return sizes, nil
}
