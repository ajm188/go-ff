package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ajm188/go-ff/feature"
	featurepb "github.com/ajm188/go-ff/proto/feature"
)

var setFeatureCmd = &cobra.Command{
	Use:          "set feature-name [type]",
	Args:         cobra.RangeArgs(1, 2),
	RunE:         setFeature,
	SilenceUsage: false,
}

var setFeatureOptions featurepb.Feature

func setFeature(cmd *cobra.Command, args []string) error {
	var (
		t        *featurepb.Feature_Type
		typeName string
	)

	if cmd.Flags().NArg() > 1 {
		typeName = cmd.Flags().Arg(1)
		_t, err := feature.ParseType(typeName)
		if err != nil {
			return err
		}

		t = &_t
	}

	cmd.SilenceUsage = true
	name := cmd.Flags().Arg(0)

	resp, err := client.GetFeature(ctx, &featurepb.GetFeatureRequest{
		Name: name,
	})
	if err != nil {
		// Note: we cannot use errors.Is here because we're getting back a gRPC
		// error, which does not wrap the error returned by our server
		// implementation.
		if strings.Contains(err.Error(), feature.ErrNoFeature.Error()) {
			if t == nil {
				cmd.SilenceUsage = false
				return fmt.Errorf("no feature named %s, must specify a type", name)
			}

			log.Printf("no feature named %s, creating ...", name)
			setFeatureOptions.Name = name
			setFeatureOptions.Type = *t

			_, err := client.SetFeature(ctx, &featurepb.SetFeatureRequest{Feature: &setFeatureOptions})
			return err
		}

		return err
	}

	feat := resp.Feature

	if cmd.Flags().Changed("description") {
		feat.Description = setFeatureOptions.Description
	}

	if cmd.Flags().Changed("enabled") {
		feat.Enabled = setFeatureOptions.Enabled
	}

	if cmd.Flags().Changed("percentage") {
		feat.Percentage = setFeatureOptions.Percentage
	}

	if cmd.Flags().Changed("expression") {
		feat.Expression = setFeatureOptions.Expression
	}

	if t != nil {
		cmd.SilenceUsage = false

		switch *t {
		case featurepb.Feature_CONSTANT:
			if cmd.Flags().Changed("percentage") {
				return fmt.Errorf("--percentage is incompatible with feature type %s", typeName)
			}

			if cmd.Flags().Changed("expression") {
				return fmt.Errorf("--expression is incompatible with feature type %s", typeName)
			}
		case featurepb.Feature_PERCENTAGE_BASED:
			if cmd.Flags().Changed("enabled") {
				return fmt.Errorf("--enabled is incompatible with feature type %s", typeName)
			}

			if cmd.Flags().Changed("expression") {
				return fmt.Errorf("--expression is incompatible with feature type %s", typeName)
			}
		case featurepb.Feature_EXPRESSION:
			if cmd.Flags().Changed("enabled") {
				return fmt.Errorf("--enabled is incompatible with feature type %s", typeName)
			}

			if cmd.Flags().Changed("percentage") {
				return fmt.Errorf("--percentage is incompatible with feature type %s", typeName)
			}
		}

		cmd.SilenceUsage = true
	}

	_, err = client.SetFeature(ctx, &featurepb.SetFeatureRequest{
		Feature: feat,
	})
	return err
}

func init() {
	setFeatureCmd.Flags().StringVarP(&setFeatureOptions.Description, "description", "d", "", "description of the feature")
	setFeatureCmd.Flags().BoolVar(&setFeatureOptions.Enabled, "enabled", false, "enable this feature. only used for type=CONSTANT")
	setFeatureCmd.Flags().Uint32VarP(&setFeatureOptions.Percentage, "percentage", "p", 0, "percentage [0, 100] of requests for which the feature should be enabled. only used for type=PERCENTAGE_BASED")
	setFeatureCmd.Flags().StringVarP(&setFeatureOptions.Expression, "expression", "e", "", "govaluate expression string. only used for type=EXPRESSION")
	rootCmd.AddCommand(setFeatureCmd)
}
