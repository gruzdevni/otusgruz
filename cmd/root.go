package cmd

import (
	"context"

	"otusgruz/config"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func Run(ctx context.Context, conf config.Config) error {
	root := &cobra.Command{ //nolint:exhaustruct
		RunE: func(cmd *cobra.Command, _ []string) error {
			//nolint:wrapcheck
			return cmd.Usage()
		},
	}

	root.AddCommand(
		restCmd(ctx, conf),
	)

	return errors.Wrap(root.ExecuteContext(ctx), "run application")
}
