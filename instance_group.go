package achim

import (
	"context"
	"fmt"

	v3 "github.com/exoscale/egoscale/v3"
)

type InstanceGroup struct {
	Names        []string
	Key          *v3.SSHKey
	Autostart    bool
	Template     *v3.Template
	InstanceType *v3.InstanceType
	DiskSizeGB   int
	Labels       []Label
}

func (p *NewGroupParams) Compile(ctx context.Context) (*InstanceGroup, error) {
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
	group, err := ParseGroupFile(p.GroupFile)
	if err != nil {
		return nil, fmt.Errorf("parse group file at %s: %w", p.GroupFile, err)
	}
	requiredNames := make(map[string]struct{}, 0)
	for _, user := range group.Users {
		requiredNames[user.Name] = struct{}{}
	}
	missingNames := make([]string, 0)
	for required := range requiredNames {
		if _, ok := existingNames[required]; !ok {
			missingNames = append(missingNames, required)
		}
	}
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

	return &InstanceGroup{
		Names:        missingNames,
		Key:          sshKey,
		Autostart:    p.Autostart,
		Template:     template,
		InstanceType: instanceType,
		DiskSizeGB:   int(diskSizeGb),
		Labels:       labels,
	}, nil
}
