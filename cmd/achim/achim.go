package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/composed-ch/achim"
	"github.com/urfave/cli/v3"
)

var byFlag = &cli.StringFlag{
	Name:    "by",
	Usage:   "filter by label/value selector",
	Aliases: []string{"b"},
}

func main() {
	cmd := &cli.Command{
		Name:  "achim",
		Usage: "Advanced Cloud Hyperscaler Infrastructure Manager",
		Commands: []*cli.Command{
			{
				Name: "dns",
				Commands: []*cli.Command{
					{
						Name: "flush",
					},
					{
						Name: "sync",
					},
				},
			},
			{
				Name: "group",
				Commands: []*cli.Command{
					{
						Name: "create",
					},
					{
						Name: "export-inventory",
					},
					{
						Name: "export-overview",
					},
					{
						Name: "export-playbook",
					},
					{
						Name: "file-from-text",
					},
				},
			},
			{
				Name:  "image",
				Usage: "information on available images",
				Commands: []*cli.Command{
					{
						Name:  "list",
						Usage: "list available images",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "contains",
								Usage:   "filter by image name",
								Aliases: []string{"c"},
							},
						},
						Before: before,
						Action: func(ctx context.Context, c *cli.Command) error {
							contains := c.String("contains")
							images, err := achim.ListImages(ctx, contains)
							if err != nil {
								return fmt.Errorf(`list images containing "%s": %v`, contains, err)
							}
							for _, image := range images {
								fmt.Println(image)
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "instance",
				Usage: "manage compute instances",
				Commands: []*cli.Command{
					{
						Name: "check",
					},
					{
						Name: "create",
					},
					{
						Name:  "deprotect",
						Usage: "remove instance protection",
						Flags: []cli.Flag{
							byFlag,
						},
						Before: before,
						Action: func(ctx context.Context, c *cli.Command) error {
							by := c.String("by")
							return achim.DeprotectInstances(ctx, by)
						},
					},
					{
						Name: "destroy",
					},
					{
						Name:  "list",
						Usage: "list instances",
						Flags: []cli.Flag{
							byFlag,
						},
						Before: before,
						Action: func(ctx context.Context, c *cli.Command) error {
							by := c.String("by")
							instances, err := achim.ListInstances(ctx, by)
							if err != nil {
								return fmt.Errorf(`list instances by "%s": %w`, by, err)
							}
							for _, instance := range instances {
								fmt.Println(achim.FormatInstance(instance))
							}
							return nil
						},
					},
					{
						Name:  "label",
						Usage: "add a label to multiple instances",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "label",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "value",
								Required: true,
							},
							byFlag,
						},
						Before: before,
						Action: func(ctx context.Context, c *cli.Command) error {
							label := c.String("label")
							value := c.String("value")
							by := c.String("by")
							return achim.LabelInstances(ctx, label, value, by)
						},
					},
					{
						Name: "probe",
					},
					{
						Name:  "protect",
						Usage: "add instance protection",
						Flags: []cli.Flag{
							byFlag,
						},
						Before: before,
						Action: func(ctx context.Context, c *cli.Command) error {
							by := c.String("by")
							return achim.ProtectInstances(ctx, by)
						},
					},
					{
						Name: "resize",
					},
					{
						Name: "scale",
					},
					{
						Name: "start",
					},
					{
						Name: "stop",
					},
					{
						Name:  "type",
						Usage: "list instance types",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "family",
								Usage:   "instance type family",
								Aliases: []string{"f"},
								Value:   "standard",
							},
						},
						Before: before,
						Action: func(ctx context.Context, c *cli.Command) error {
							family := c.String("family")
							types, err := achim.ListInstanceTypes(ctx, family)
							if err != nil {
								return fmt.Errorf("list instance types: %v", err)
							}
							for _, t := range types {
								out, _ := json.Marshal(t)
								fmt.Println(string(out))
							}
							return nil
						},
					},
				},
			},
			{
				Name: "network",
				Commands: []*cli.Command{
					{
						Name: "attach",
					},
					{
						Name: "cleanup",
					},
					{
						Name: "create",
					},
					{
						Name: "destroy",
					},
					{
						Name: "flush",
					},
					{
						Name: "list",
					},
				},
			},
			{
				Name: "scenario",
				Commands: []*cli.Command{
					{
						Name: "create",
					},
					{
						Name: "destroy",
					},
					{
						Name: "export-overview",
					},
				},
			},
			{
				Name: "tenant",
				Commands: []*cli.Command{
					{
						Name:  "add",
						Usage: "add a tenant",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Usage:    "tenant name",
								Aliases:  []string{"n"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "key",
								Usage:    "Exoscale API key",
								Aliases:  []string{"k"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "secret",
								Usage:    "Exoscale API secret",
								Aliases:  []string{"s"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "zone",
								Usage:    "Exoscale zone",
								Aliases:  []string{"z"},
								Required: true,
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							return achim.AddTenant(achim.Tenant{
								Name:   cmd.String("name"),
								Key:    cmd.String("key"),
								Secret: cmd.String("secret"),
								Zone:   cmd.String("zone"),
							}, tenantsFilePath())
						},
					},
					{
						Name:  "default",
						Usage: "set the default tenant",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "name",
								Usage:   "tenant name",
								Aliases: []string{"n"},
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							return achim.SetDefaultTenant(cmd.String("name"), tenantsFilePath())
						},
					},
					{
						Name:  "remove",
						Usage: "remove a tenant",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "name",
								Usage:   "Tenant name",
								Aliases: []string{"n"},
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							return achim.RemoveTenant(cmd.String("name"), tenantsFilePath())
						},
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func before(ctx context.Context, c *cli.Command) (context.Context, error) {
	tenant, err := achim.GetDefaultTenant(tenantsFilePath())
	if err != nil {
		return nil, fmt.Errorf("get default tenant: %v", err)
	}
	client, err := tenant.Client()
	if err != nil {
		return nil, fmt.Errorf("get Exoscale client: %v", err)
	}
	return context.WithValue(ctx, "exo", client), nil
}

func tenantsFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s%c%s", home, os.PathSeparator, ".achim.json")
}
