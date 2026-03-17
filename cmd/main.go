package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/phantompunk/fig/internal/font"
	"github.com/phantompunk/fig/internal/input"
	"github.com/phantompunk/fig/internal/render"
	"github.com/phantompunk/fig/internal/tui"
	"github.com/phantompunk/fig/internal/vcs"
	"github.com/spf13/cobra"
)

var (
	fontName  string
	listFonts bool
	center    bool
	right     bool
	// rtl       bool
)

func main() {
	if err := execute(); err != nil {
		fmt.Fprintf(os.Stderr, "fig: %v\n", err)
		os.Exit(1)
	}
}

func execute() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	return buildCmd().ExecuteContext(ctx)
}

func buildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "fig",
		Short:        "Render text as ASCII art using FIGlet fonts",
		Args:         cobra.ArbitraryArgs,
		RunE:         run,
		SilenceUsage: true,
	}

	cmd.Flags().BoolVarP(&listFonts, "list-fonts", "l", false, "List all available fonts")
	cmd.Flags().StringVarP(&fontName, "font", "f", "standard", "Specify a font, default is standard")
	cmd.Flags().BoolVarP(&center, "center", "c", false, "Center text in terminal")
	cmd.Flags().BoolVarP(&right, "right", "r", false, "Right align text in terminal")
	// cmd.Flags().BoolVarP(&rtl, "right-to-left", "rtl", false, "Print text right to left")

	v, commit := vcs.Version()
	cmd.Version = v
	cmd.SetVersionTemplate(fmt.Sprintf("%s version %s (%s)\n", "fig", v, commit))

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if listFonts {
		fonts := font.ListFonts()
		fmt.Println("Supported fonts:", strings.Join(fonts, ", "))
		return nil
	}

	src := input.Resolve(args)
	msg, err := src.Read()
	if err != nil {
		return err
	}

	if len(msg) == 0 {
		return tui.Start()
	}

	engine := render.New(font.BundledLoader())
	out, err := engine.Render(msg, render.RenderOptions{FontName: fontName})
	if err != nil {
		return err
	}
	fmt.Print(out)
	// font, err := font.LoadFont(fontName)
	// if err != nil {
	// 	return err
	// }
	//
	// fmt.Println(font.Render(msg))
	return nil
}
