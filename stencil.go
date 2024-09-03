package main

import (
	"fmt"
	"os"

	"github.com/phantompunk/stencil/internal/tui"
	"github.com/phantompunk/stencil/pkg/font"
	"github.com/phantompunk/stencil/pkg/stencil"
	"github.com/spf13/cobra"
)

var fontface string

var stencilCmd = &cobra.Command{
	Use: "stencil",
	Short: "stencil renders figlet fonts",
	Run: func (cmd *cobra.Command, args []string) {
		fmt.Println("Calling base")
		tui.Render()
		// tui.Run()
	},
}

var printCmd = &cobra.Command{
	Use: "print",
	Short: "print renders figlet fonts",
	Run: func (cmd *cobra.Command, args []string){
		if len(args) == 0 {
			fmt.Println("Needs a phrase")
			os.Exit(1)
		}

		phrase := args[0]
		st, err := stencil.NewStencil(phrase, fontface)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// st.DrawText()
		fmt.Println(st.Draw())
	},
}

var listCmd = &cobra.Command{
	Use: "list",
	Short: "list all available figlet fonts",
	Run: func (cmd *cobra.Command, args []string) {
		fmt.Println("List all fonts")
		for _, ff := range font.ListFonts() {
			st, err := stencil.NewStencil(ff, ff)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(ff+":")
			st.DrawText()
		}
	},
}

func main() {
	if err := stencilCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	stencilCmd.AddCommand(listCmd)
	stencilCmd.AddCommand(printCmd)
	stencilCmd.PersistentFlags().StringVarP(&fontface, "font", "f", "drpepper", "Figlet Font to use")
}

