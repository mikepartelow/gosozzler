package cmd

import (
	"fmt"
	"mp/sozzler/pkg/sozzler"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <recipe name>",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		catalog := cmd.Context().Value("catalog").(*sozzler.RecipeCatalog)
		display := cmd.Context().Value("display").(sozzler.Display)

		name := args[0]
		recipe, ok := catalog.Find(name)
		if !ok {
			display.Error(fmt.Sprintf("couldn't find recipe %q", name))
			return
		}

		scale, _ := cmd.Flags().GetInt("scale")
		for idx, _ := range recipe.Components {
			recipe.Components[idx].Quantity *= sozzler.FractionFloat64(scale)
		}

		display.Show(recipe)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().IntP("scale", "s", 1, "scale recipe")
}
