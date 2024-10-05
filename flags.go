package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

var phase = cli.StringFlag{
	Name:  "phase",
	Usage: "pods in lifecycle `phase`",
	Action: func(ctx *cli.Context, s string) error {
		possible := []string{"pending", "running", "succeeded", "failed", "unknown"}
		for _, p := range possible {
			if ctx.String("phase") == p {
				return nil
			}
		}
		return fmt.Errorf("possible phases: %v", strings.Join(possible, ", "))
	},
}
