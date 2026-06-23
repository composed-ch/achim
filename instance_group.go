package achim

import (
	"context"
	"fmt"

	v3 "github.com/exoscale/egoscale/v3"
)

// FIXME: use this one as a receiver; parsing the groups file is a client issue
type NewInstancesParam struct {
	Names     []string
	Key       string
	Autostart bool
	Image     string
	Size      string
	Labels    string
}

func (p *NewInstancesParam) Compile(ctx context.Context) (*InstanceGroup, error) {
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
	missingNames := make([]string, 0)
	for _, required := range p.Names {
		if _, ok := existingNames[required]; !ok {
			missingNames = append(missingNames, required)
		}
	}
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
	return &InstanceGroup{
		Names:        missingNames,
		Key:          sshKey,
		Autostart:    p.Autostart,
		Template:     template,
		InstanceType: instanceType,
		DiskSizeGB:   diskSizeGb,
		Labels:       labels, // FIXME: merge generic labels, e.g. "owner"
	}, nil
}

type InstanceGroup struct {
	Names        []string
	Key          *v3.SSHKey
	Autostart    bool
	Template     *v3.Template
	InstanceType *v3.InstanceType
	DiskSizeGB   int64
	Labels       []Label
}

func (p *InstanceGroup) Create(ctx context.Context) error {
	exo := ctx.Value("exo").(*v3.Client)
	for _, name := range p.Names {
		_, err := exo.CreateInstance(ctx, v3.CreateInstanceRequest{
			AutoStart:    &p.Autostart,
			DiskSize:     p.DiskSizeGB,
			InstanceType: p.InstanceType,
			Labels:       AsMap(p.Labels),
			Name:         name,
			SSHKey:       p.Key,
			Template:     p.Template,
		})
		if err != nil {
			return fmt.Errorf("create instance: %w", err)
		}
	}
	return nil
}
