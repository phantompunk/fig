package main

import (
	"fmt"
	"os"

	"github.com/phantompunk/stencil/pkg/stencil"
	"github.com/spf13/cobra"
)

var stencilCmd = &cobra.Command{
	Use: "stencil",
	Short: "stencil renders figlet fonts",
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

func main() {
	if err := stencilCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}

