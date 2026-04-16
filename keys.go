package achim

import (
	"context"
	"fmt"

	v3 "github.com/exoscale/egoscale/v3"
)

func GetSSHKeyByName(ctx context.Context, name string) (*v3.SSHKey, error) {
	exo := ctx.Value("exo").(*v3.Client)
	res, err := exo.ListSSHKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("list SSH keys: %w", err)
	}
	for _, key := range res.SSHKeys {
		if key.Name == name {
			return &key, nil
		}
	}
	return nil, fmt.Errorf(`no SSH key for name "%s" found`)
}
