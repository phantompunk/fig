package main

import (
	"fmt"
	"os"

	"github.com/phantompunk/stencil/pkg/font"
	"github.com/phantompunk/stencil/pkg/stencil"
	"github.com/spf13/cobra"
)

var stencilCmd = &cobra.Command{
	Use: "stencil",
	Short: "stencil renders figlet fonts",
}

var printCmd = &cobra.Command{
	Use: "print",
	Short: "print renders figlet fonts",
	Run: func (cmd *cobra.Command, args []string){
		font := args[1]
		phrase := args[0]
		st, err := stencil.NewStencil(phrase, font)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		st.DrawText()
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
}

