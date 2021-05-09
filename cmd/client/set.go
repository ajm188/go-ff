package main

import (
	"strconv"

	"github.com/spf13/cobra"

	featurepb "github.com/ajm188/goff/proto/feature"
)

var setFeatureCmd = &cobra.Command{
	Use:          "set feature-name state",
	Args:         cobra.ExactArgs(2),
	RunE:         setFeature,
	SilenceUsage: false,
}

func setFeature(cmd *cobra.Command, args []string) error {
	state, err := strconv.ParseBool(cmd.Flags().Arg(1))
	if err != nil {
		return err
	}

	cmd.SilenceUsage = true

	_, err = client.SetFeature(ctx, &featurepb.SetFeatureRequest{
		Feature: &featurepb.Feature{
			Name:    cmd.Flags().Arg(0),
			Enabled: state,
		},
	})
	return err
}

func init() {
	rootCmd.AddCommand(setFeatureCmd)
}
