package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/composed-ch/achim"
	"github.com/urfave/cli/v3"
)

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
				Name: "image",
				Commands: []*cli.Command{
					{
						Name: "list",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "contains",
								Usage:   "Filter by Image Name",
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
				Name: "instance",
				Commands: []*cli.Command{
					{
						Name: "check",
					},
					{
						Name: "create",
					},
					{
						Name: "deprotect",
					},
					{
						Name: "destroy",
					},
					{
						Name: "list",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "by",
								Usage:   "Filter by Label/Value Selector",
								Aliases: []string{"b"},
							},
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
						Name: "label",
					},
					{
						Name: "probe",
					},
					{
						Name: "protect",
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
						Name: "type",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "family",
								Usage:   "Instance Type Family",
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
						Name: "add",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Usage:    "Tenant Name",
								Aliases:  []string{"n"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "key",
								Usage:    "Exoscale API Key",
								Aliases:  []string{"k"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "secret",
								Usage:    "Exoscale API Secret",
								Aliases:  []string{"s"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "zone",
								Usage:    "Exoscale Zone",
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
						Name: "default",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "name",
								Usage:   "Tenant Name",
								Aliases: []string{"n"},
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							return achim.SetDefaultTenant(cmd.String("name"), tenantsFilePath())
						},
					},
					{
						Name: "remove",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "name",
								Usage:   "Tenant Name",
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
