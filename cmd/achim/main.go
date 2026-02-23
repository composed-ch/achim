package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "achim",
		Usage: "advanced cloud hyperscaler infrastructure manager",
		Action: func(context.Context, *cli.Command) error {
			fmt.Println("achim says hi!")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
