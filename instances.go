package achim

import (
	"context"
	"fmt"
	"maps"

	v3 "github.com/exoscale/egoscale/v3"
)

type NewInstanceParams struct {
	Name      string
	Key       string
	Autostart bool
	Image     string
	Size      string
	Labels    string
}

func CreateInstance(ctx context.Context, params NewInstanceParams) error {
	exo := ctx.Value("exo").(*v3.Client)
	labels, err := ParseLabels(params.Labels)
	if err != nil {
		return fmt.Errorf("parse labels: %w", err)
	}
	existingInstances, err := ListInstances(ctx, "")
	if err != nil {
		return fmt.Errorf("list instances: %w", err)
	}
	for _, instance := range existingInstances {
		if instance.Name == params.Name {
			return fmt.Errorf(`instance with name "%s" already exists`, params.Name)
		}
	}
	template, err := GetTemplateByName(ctx, params.Image)
	if err != nil {
		return fmt.Errorf("find image: %w", err)
	}
	diskSizeGb := template.Size / (1024 * 1024 * 1024)
	if diskSizeGb == 0 {
		diskSizeGb = 10
	}
	instanceType, err := GetInstanceTypeBySize(ctx, params.Size)
	if err != nil {
		return fmt.Errorf("find instance type: %w", err)
	}
	sshKey, err := GetSSHKeyByName(ctx, params.Key)
	if err != nil {
		return fmt.Errorf("find SSH key: %w", err)
	}
	_, err = exo.CreateInstance(ctx, v3.CreateInstanceRequest{
		AutoStart:    &params.Autostart,
		DiskSize:     diskSizeGb,
		InstanceType: instanceType,
		Labels:       AsMap(labels),
		Name:         params.Name,
		SSHKey:       sshKey,
		Template:     template,
	})
	if err != nil {
		return fmt.Errorf("create instance: %w", err)
	}
	return nil
}

func DestroyInstances(ctx context.Context, by string) error {
	instances, err := ListInstances(ctx, by)
	if err != nil {
		return fmt.Errorf("list instances: %w", err)
	}
	exo := ctx.Value("exo").(*v3.Client)
	for _, instance := range instances {
		_, err := exo.DeleteInstance(ctx, instance.ID)
		if err != nil {
			return fmt.Errorf(`destroy instance "%s": %w`, instance.ID, err)
		}
	}
	return nil
}

func EmbiggenDisk(ctx context.Context, by string, gb int64) error {
	instances, err := ListInstances(ctx, by)
	if err != nil {
		return fmt.Errorf("list instances: %w", err)
	}
	exo := ctx.Value("exo").(*v3.Client)
	for _, instance := range instances {
		if instance.DiskSize > gb {
			return fmt.Errorf(`instance %s "%s" disk size is %d GB; cannot shrink disk`,
				instance.ID, instance.Name, instance.DiskSize)
		}
		if instance.State != v3.InstanceStateStopped {
			return fmt.Errorf(`instance %s "%s" must be stopped but is %s`,
				instance.ID, instance.Name, instance.State)
		}
		_, err := exo.ResizeInstanceDisk(ctx, instance.ID, v3.ResizeInstanceDiskRequest{
			DiskSize: gb,
		})
		if err != nil {
			return fmt.Errorf(`resize disk of instance %s "%s": %w`,
				instance.ID, instance.Name, err)
		}
	}
	return nil
}

func ListInstances(ctx context.Context, by string) ([]v3.Instance, error) {
	exo := ctx.Value("exo").(*v3.Client)
	result, err := exo.ListInstances(ctx)
	if err != nil {
		return nil, fmt.Errorf("list instances: %w", err)
	}
	instances := make([]v3.Instance, 0)
	for _, item := range result.Instances {
		instance, err := getInstanceByID(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf(`get instance by ID "%s": %w`, item.ID, err)
		}
		instances = append(instances, *instance)
	}
	instances, err = filterInstances(instances, by)
	if err != nil {
		return nil, fmt.Errorf("filter instances: %w", err)
	}
	return instances, nil
}

func ListInstanceTypes(ctx context.Context, family string) ([]v3.InstanceType, error) {
	exo := ctx.Value("exo").(*v3.Client)
	res, err := exo.ListInstanceTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("list instance types: %w", err)
	}
	var result []v3.InstanceType
	for _, it := range res.InstanceTypes {
		if !*it.Authorized || string(it.Family) != family {
			continue
		}
		result = append(result, it)
	}
	return result, nil
}

func GetInstanceTypeBySize(ctx context.Context, size string) (*v3.InstanceType, error) {
	exo := ctx.Value("exo").(*v3.Client)
	res, err := exo.ListInstanceTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("list instance types: %w", err)
	}
	for _, it := range res.InstanceTypes {
		if string(it.Size) == size {
			return &it, nil
		}
	}
	return nil, fmt.Errorf(`no instance type for size "%s" found`, size)
}

func LabelInstances(ctx context.Context, label, value, by string) error {
	exo := ctx.Value("exo").(*v3.Client)
	instances, err := ListInstances(ctx, by)
	if err != nil {
		return fmt.Errorf(`list instances to label: %w`, err)
	}
	for _, instance := range instances {
		labels := make(map[string]string, 0)
		maps.Copy(labels, instance.Labels)
		labels[label] = value
		update := v3.UpdateInstanceRequest{Labels: labels}
		if _, err := exo.UpdateInstance(ctx, instance.ID, update); err != nil {
			return fmt.Errorf(`label instance %s: %w`, instance.ID, err)
		}
	}
	return nil
}

func ProtectInstances(ctx context.Context, by string) error {
	return changeInstanceProtection(ctx, by, true)
}

func DeprotectInstances(ctx context.Context, by string) error {
	return changeInstanceProtection(ctx, by, false)
}

func StartInstances(ctx context.Context, by string) error {
	return changeInstanceState(ctx, by, true)
}

func StopInstances(ctx context.Context, by string) error {
	return changeInstanceState(ctx, by, false)
}

func FormatInstance(instance v3.Instance) string {
	return fmt.Sprintf("%-40s %-32s %15s %-10s", instance.ID, instance.Name, instance.PublicIP, instance.State)
}

func changeInstanceProtection(ctx context.Context, by string, protect bool) error {
	exo := ctx.Value("exo").(*v3.Client)
	instances, err := ListInstances(ctx, by)
	if err != nil {
		return fmt.Errorf(`list instances to change protection to %v: %w`, protect, err)
	}
	for _, instance := range instances {
		if protect {
			if _, err := exo.AddInstanceProtection(ctx, instance.ID); err != nil {
				return fmt.Errorf(`protect instance %s: %w`, instance.ID, err)
			}
		} else {
			if _, err := exo.RemoveInstanceProtection(ctx, instance.ID); err != nil {
				return fmt.Errorf(`deprotect instance %s: %w`, instance.ID, err)
			}
		}
	}
	return nil
}

func changeInstanceState(ctx context.Context, by string, up bool) error {
	exo := ctx.Value("exo").(*v3.Client)
	instances, err := ListInstances(ctx, by)
	if err != nil {
		return fmt.Errorf(`list instances to change up state to %v: %w`, up, err)
	}
	for _, instance := range instances {
		if up {
			if _, err := exo.StartInstance(ctx, instance.ID, v3.StartInstanceRequest{}); err != nil {
				return err
			}
		} else {
			if _, err := exo.StopInstance(ctx, instance.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func getInstanceByID(ctx context.Context, id v3.UUID) (*v3.Instance, error) {
	exo := ctx.Value("exo").(*v3.Client)
	instance, err := exo.GetInstance(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(`get instance by ID "%v": %w`, id, err)
	}
	return instance, nil
}

func filterInstances(instances []v3.Instance, by string) ([]v3.Instance, error) {
	if by == "" {
		return instances, nil
	}
	selectors, err := ParseLabels(by)
	if err != nil {
		return nil, fmt.Errorf(`parse --by "%s": %w`, by, err)
	}
	filtered := make([]v3.Instance, 0)
	for _, instance := range instances {
		retain := true
		for _, selector := range selectors {
			if selector.Key == "name" {
				if instance.Name != selector.Value {
					retain = false
					break
				}
			} else {
				if value, ok := instance.Labels[selector.Key]; !ok {
					retain = false
					break
				} else if value != selector.Value {
					retain = false
					break
				}
			}
		}
		if retain {
			filtered = append(filtered, instance)
		}
	}
	return filtered, nil
}
