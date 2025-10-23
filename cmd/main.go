package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/phantompunk/fig"
	"github.com/spf13/cobra"
)

func main() {
	if err := execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func execute() error {
	rootCmd := initialize()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		fmt.Println("Received signal", sig)

		cancel()
	}()

	return rootCmd.ExecuteContext(ctx)
}

func initialize() *cobra.Command {
	root := &cobra.Command{
		Use:  "fig",
		RunE: figFunc,
	}

	return root
}

func figFunc(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		fmt.Println("TBD UI")
		return nil
	}

	msg := strings.Join(args, "")

	font, err := fig.Font("standard")
	if err != nil {
		return err
	}

	fmt.Println(font.Render(msg))
	return nil
}
