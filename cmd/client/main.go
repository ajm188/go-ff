package main

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	featurepb "github.com/ajm188/goff/proto/feature"
)

var (
	ctx = context.Background()

	addr   string
	cc     *grpc.ClientConn
	client featurepb.FeaturesClient

	rootCmd = &cobra.Command{
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cc, err = grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				return err
			}

			client = featurepb.NewFeaturesClient(cc)
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return cc.Close()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&addr, "server", "s", ":15000", "server address to make requests against")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
