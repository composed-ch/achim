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
		return nil, fmt.Errorf("list templates: %v", err)
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
