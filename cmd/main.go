package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/phantompunk/fig"
	"github.com/phantompunk/fig/internal/vcs"
	"github.com/phantompunk/fig/tui"
	"github.com/spf13/cobra"
)

var (
	fontName  string
	listFonts bool
	version   string
	commit    string
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
	rootCmd := &cobra.Command{
		Use:  "fig",
		RunE: figFunc,
	}

	rootCmd.Flags().BoolVarP(&listFonts, "list-fonts", "l", false, "List all available fonts")
	rootCmd.Flags().StringVarP(&fontName, "font", "f", "standard", "Specify a font, default is standard")

	version, commit = vcs.Version()
	setVersion(rootCmd)

	return rootCmd
}

func figFunc(cmd *cobra.Command, args []string) error {
	if listFonts {
		fonts := fig.ListFonts()
		fmt.Println("Supported fonts:", strings.Join(fonts, ", "))
		return nil
	}

	if len(args) == 0 {
		return tui.Start()
	}

	msg := strings.Join(args, "")
	font, err := fig.Font(fontName)
	if err != nil {
		return err
	}

	fmt.Println(font.Render(msg))
	return nil
}

func setVersion(cmd *cobra.Command) {
	vt := fmt.Sprintf("%s version %s (%s)\n", "fig", version, commit)
	cmd.SetVersionTemplate(vt)
	cmd.Version = version
}
