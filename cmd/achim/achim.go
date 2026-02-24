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
				Name: "instance",
				Commands: []*cli.Command{
					{
						Name: "create",
					},
					{
						Name: "destroy",
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
						Name: "destroy",
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
