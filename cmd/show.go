package cmd

import (
	"fmt"
	"mp/sozzler/pkg/display"
	"mp/sozzler/pkg/sozzler"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <recipe name>",
	Short: "Show a recipe.",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		catalog := cmd.Context().Value(catalogKey{}).(*sozzler.RecipeCatalog)
		display := cmd.Context().Value(displayKey{}).(display.Display)

		name := args[0]
		recipe, ok := catalog.Find(name)
		if !ok {
			display.Error(fmt.Sprintf("couldn't find recipe %q", name))
			return
		}

		scale, _ := cmd.Flags().GetInt("scale")
		for idx := range recipe.Components {
			recipe.Components[idx].Quantity *= sozzler.FractionFloat64(scale)
		}

		display.Show(recipe)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().IntP("scale", "s", 1, "scale recipe")
}
