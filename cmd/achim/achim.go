package main

import (
	"context"
	"fmt"
	"os"

	"github.com/composed-ch/achim"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name: "achim",
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
				Name: "images",
				Commands: []*cli.Command{
					{
						Name: "list",
					},
					{
						Name: "list-types",
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
								Name:    "name",
								Usage:   "Tenant Name",
								Aliases: []string{"n"},
							},
							&cli.StringFlag{
								Name:    "key",
								Usage:   "Exoscale API Key",
								Aliases: []string{"k"},
							},
							&cli.StringFlag{
								Name:    "secret",
								Usage:   "Exoscale API Secret",
								Aliases: []string{"s"},
							},
							&cli.StringFlag{
								Name:    "zone",
								Usage:   "Exoscale Zone",
								Aliases: []string{"z"},
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							return achim.AddTenant(achim.Tenant{
								Name:   cmd.String("name"),
								Key:    cmd.String("key"),
								Secret: cmd.String("secret"),
								Zone:   cmd.String("zone"),
							}, "tenants.json")
						},
					},
					{
						Name: "default",
					},
					{
						Name: "remove",
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
