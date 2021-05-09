package main

import (
	"github.com/spf13/cobra"

	"github.com/ajm188/goff/feature"
	featurepb "github.com/ajm188/goff/proto/feature"
)

var setFeatureCmd = &cobra.Command{
	Use:          "set feature-name type",
	Args:         cobra.ExactArgs(2),
	RunE:         setFeature,
	SilenceUsage: false,
}

var setFeatureOptions featurepb.Feature

func setFeature(cmd *cobra.Command, args []string) error {
	t, err := feature.ParseType(cmd.Flags().Arg(1))
	if err != nil {
		return err
	}

	cmd.SilenceUsage = true

	setFeatureOptions.Type = t
	setFeatureOptions.Name = cmd.Flags().Arg(0)

	_, err = client.SetFeature(ctx, &featurepb.SetFeatureRequest{
		Feature: &setFeatureOptions,
	})
	return err
}

func init() {
	setFeatureCmd.Flags().BoolVar(&setFeatureOptions.Enabled, "enabled", false, "enable this feature. only used for type=CONSTANT")
	setFeatureCmd.Flags().Uint32VarP(&setFeatureOptions.Percentage, "percentage", "p", 0, "percentage [0, 100] of requests for which the feature should be enabled. only used for type=PERCENTAGE_BASED")
	setFeatureCmd.Flags().StringVarP(&setFeatureOptions.Expression, "expression", "e", "", "govaluate expression string. only used for type=EXPRESSION")
	rootCmd.AddCommand(setFeatureCmd)
}
