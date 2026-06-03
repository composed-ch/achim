package achim

import (
	"context"
	"fmt"
)

type InstanceGroup struct {
	Names     []string
	Key       string
	Autostart bool
	Image     string
	Size      string
	Labels    []Label
}

func (p *NewGroupParams) Compile(ctx context.Context) (*InstanceGroup, error) {
	/*
		exo := ctx.Value("exo").(*v3.Client)
		labels, err := ParseLabels(p.Labels)
		if err != nil {
			return nil, fmt.Errorf("parse labels: %w", err)
		}
		existingInstances, err := ListInstances(ctx, "")
		if err != nil {
			return nil, fmt.Errorf("list instances: %w", err)
		}
		existingNames := make(map[string]struct{}, 0)
		for _, i := range existingInstances {
			existingNames[i.Name] = struct{}{}
		}
		requiredNames := make(map[string]struct{}, 0)

		// TODO: skip instances that are not needed
		template, err := GetTemplateByName(ctx, p.Image)
		if err != nil {
			return nil, fmt.Errorf("find image: %w", err)
		}
		diskSizeGb := template.Size / (1024 * 1024 * 1024)
		if diskSizeGb == 0 {
			diskSizeGb = 10
		}
		instanceType, err := GetInstanceTypeBySize(ctx, p.Size)
		if err != nil {
			return nil, fmt.Errorf("find instance type: %w", err)
		}
		sshKey, err := GetSSHKeyByName(ctx, p.Key)
		if err != nil {
			return nil, fmt.Errorf("find SSH key: %w", err)
		}
	*/
	return nil, fmt.Errorf("hello")
}
