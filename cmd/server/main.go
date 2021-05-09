package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/ajm188/goff/feature"
)

var (
	addr    string
	rootCmd = &cobra.Command{
		RunE:          serve,
		SilenceErrors: true,
	}
)

func serve(cmd *cobra.Command, args []string) error {
	s := grpc.NewServer()
	feature.RegisterServer(s)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	sigch := make(chan os.Signal, 8)
	signal.Notify(sigch, os.Interrupt, os.Kill)

	done := make(chan error, 2)

	// Listen for signals
	go func() {
		sig := <-sigch
		done <- fmt.Errorf("received signal %v", sig)
	}()

	go func() {
		done <- s.Serve(lis)
	}()

	log.Print(<-done)
	s.GracefulStop()
	return nil
}

func init() {
	rootCmd.Flags().StringVar(&addr, "addr", ":15000", "address to listen on")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
