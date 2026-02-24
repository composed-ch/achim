package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name: "achim",
		Commands: []*cli.Command{
			{
				Name: "tenant",
				Commands: []*cli.Command{
					{
						Name: "add",
					},
					{
						Name: "default",
					},
					{
						Name: "remove",
					},
				},
			},
			{
				Name: "instance",
				Commands: []*cli.Command{
					{
						Name: "create",
					},
					{
						Name: "check",
					},
					{
						Name: "protect",
					},
					{
						Name: "deprotect",
					},
					{
						Name: "destroy",
					},
					{
						Name: "label",
					},
					{
						Name: "list",
					},
					{
						Name: "probe",
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
				Name: "group",
				Commands: []*cli.Command{
					{
						Name: "create",
					},
					{
						Name: "export-overview",
					},
					{
						Name: "export-inventory",
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
				Name: "network",
				Commands: []*cli.Command{
					{
						Name: "attach",
					},
					{
						Name: "create",
					},
					{
						Name: "cleanup",
					},
					{
						Name: "flush",
					},
					{
						Name: "list",
					},
					{
						Name: "destroy",
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
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
