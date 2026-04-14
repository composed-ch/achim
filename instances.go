package achim

import (
	"context"
	"fmt"

	v3 "github.com/exoscale/egoscale/v3"
)

func ListInstances(ctx context.Context, by string) ([]v3.Instance, error) {
	return nil, nil
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
