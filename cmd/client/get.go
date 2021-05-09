package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	featurepb "github.com/ajm188/goff/proto/feature"
)

var (
	getFeatureCmd = &cobra.Command{
		Use:          "get feature",
		Aliases:      []string{"get-feature"},
		Args:         cobra.ExactArgs(1),
		RunE:         getFeature,
		SilenceUsage: true,
	}
	getFeaturesCmd = &cobra.Command{
		Use:          "list [--name-only]",
		Aliases:      []string{"get-features", "list-all"},
		Args:         cobra.NoArgs,
		RunE:         getFeatures,
		SilenceUsage: true,
	}
)

func getFeature(cmd *cobra.Command, args []string) error {
	resp, err := client.GetFeature(ctx, &featurepb.GetFeatureRequest{
		Name: cmd.Flags().Arg(0),
	})
	if err != nil {
		return err
	}

	fmt.Printf("%s:%v\n", resp.Feature.Name, resp.Feature.Enabled)
	return nil
}

var getFeaturesOptions = struct {
	NamesOnly bool
}{}

func getFeatures(cmd *cobra.Command, args []string) error {
	resp, err := client.GetFeatures(ctx, &featurepb.GetFeaturesRequest{
		NamesOnly: getFeaturesOptions.NamesOnly,
	})
	if err != nil {
		return err
	}

	if getFeaturesOptions.NamesOnly {
		fmt.Printf("%s\n", strings.Join(resp.Names, "\n"))
		return nil
	}

	buf := &strings.Builder{}
	for _, feature := range resp.Features {
		fmt.Fprintf(buf, "%s:%v\n", feature.Name, feature.Enabled)
	}

	fmt.Print(buf.String())
	return nil
}

func init() {
	rootCmd.AddCommand(getFeatureCmd)

	getFeatureCmd.Flags().BoolVar(&getFeaturesOptions.NamesOnly, "name-only", false, "show feature names only")
	rootCmd.AddCommand(getFeaturesCmd)
}
