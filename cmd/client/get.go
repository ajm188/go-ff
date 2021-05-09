package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ajm188/go-ff/feature"
	featurepb "github.com/ajm188/go-ff/proto/feature"
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
		Use:          "list [--name-only] [-j|--json]",
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
	UseJSON   bool
}{}

func getFeatures(cmd *cobra.Command, args []string) error {
	resp, err := client.GetFeatures(ctx, &featurepb.GetFeaturesRequest{
		NamesOnly: getFeaturesOptions.NamesOnly,
	})
	if err != nil {
		return err
	}

	if getFeaturesOptions.NamesOnly {
		if getFeaturesOptions.UseJSON {
			data, err := json.Marshal(resp.Names)
			if err != nil {
				return err
			}

			fmt.Printf("%s\n", data)
			return nil
		}

		fmt.Printf("%s\n", strings.Join(resp.Names, "\n"))
		return nil
	}

	if getFeaturesOptions.UseJSON {
		data, err := json.Marshal(feature.MapFromProtos(resp.Features))
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", data)
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

	getFeaturesCmd.Flags().BoolVar(&getFeaturesOptions.NamesOnly, "name-only", false, "show feature names only")
	getFeaturesCmd.Flags().BoolVarP(&getFeaturesOptions.UseJSON, "json", "j", false, "output features as JSON")
	rootCmd.AddCommand(getFeaturesCmd)
}
