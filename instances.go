package achim

import (
	"context"
	"fmt"
	"maps"

	v3 "github.com/exoscale/egoscale/v3"
)

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
		return nil, fmt.Errorf("list instance types: %v", err)
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
		exo.UpdateInstance(ctx, instance.ID, v3.UpdateInstanceRequest{Labels: labels})
	}
	return nil
}

func FormatInstance(instance v3.Instance) string {
	return fmt.Sprintf("%-40s %-32s %15s %-10s", instance.ID, instance.Name, instance.PublicIP, instance.State)
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
	selectors, err := ParseSelector(by)
	if err != nil {
		return nil, fmt.Errorf(`parse --by "%s": %w`, by, err)
	}
	filtered := make([]v3.Instance, 0)
	for _, instance := range instances {
		retain := true
		for _, selector := range selectors {
			if selector.Label == "name" {
				if instance.Name != selector.Value {
					retain = false
					break
				}
			} else {
				if value, ok := instance.Labels[selector.Label]; !ok {
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
