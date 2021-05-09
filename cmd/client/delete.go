package main

import (
	"fmt"

	"github.com/spf13/cobra"

	featurepb "github.com/ajm188/go-ff/proto/feature"
)

var deleteFeatureCmd = &cobra.Command{
	Use:          "delete feature",
	Aliases:      []string{"remove"},
	Args:         cobra.ExactArgs(1),
	RunE:         deleteFeature,
	SilenceUsage: true,
}

func deleteFeature(cmd *cobra.Command, args []string) error {
	resp, err := client.DeleteFeature(ctx, &featurepb.DeleteFeatureRequest{
		Name: cmd.Flags().Arg(0),
	})
	if err != nil {
		return err
	}

	switch resp.Feature {
	case nil:
		fmt.Printf("no such feature %s\n", cmd.Flags().Arg(0))
	default:
		fmt.Printf("deleted feature %s:%v\n", resp.Feature.Name, resp.Feature.Enabled)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(deleteFeatureCmd)
}
