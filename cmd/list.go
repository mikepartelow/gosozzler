package cmd

import (
	"fmt"
	"mp/sozzler/pkg/display"
	"mp/sozzler/pkg/sozzler"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		catalog := cmd.Context().Value("catalog").(*sozzler.RecipeCatalog)
		display := cmd.Context().Value("display").(display.Display)

		verbose, _ := cmd.Flags().GetBool("verbose")

		predicates, err := makePredicates(cmd.Flags())
		if err != nil {
			return err
		}

		if len(predicates) == 0 {
			display.List(catalog.Recipes)
			return nil
		}

		results := catalog.Search(predicates)

		if len(results) == 0 {
			display.String("no results\n")
			return nil
		}

		for recipe, result := range results {
			var matches []string
			for _, m := range result {
				matches = append(matches, fmt.Sprintf("%s (%s)", m.Match, m.Predicate.Name()))
			}
			fancyMatches := strings.Join(matches, ", ")

			display.String(recipe.Name + ": " + fancyMatches + "\n")
			if verbose {
				display.Show(recipe)
				display.String("\n")
			}
		}

		// FIXME: ensure sort order
		// FIXME: flags to change sort order

		return nil
	},
}

func makePredicates(flags *pflag.FlagSet) ([]sozzler.Predicate, error) {
	var predicates []sozzler.Predicate

	ingredients, err := flags.GetStringSlice("ingredients")
	if err != nil {
		return nil, fmt.Errorf("error reading ingredients flag: %w", err)
	}

	for _, i := range ingredients {
		predicates = append(predicates, sozzler.NewIngredientPredicate(i))
	}

	names, err := flags.GetStringSlice("names")
	if err != nil {
		return nil, fmt.Errorf("error reading names flag: %w", err)
	}

	for _, n := range names {
		predicates = append(predicates, sozzler.NewNamePredicate(n))
	}

	rating, err := flags.GetInt("rating")
	if err != nil {
		return nil, fmt.Errorf("error reading rating flag: %w", err)
	}
	if rating > 0 {
		predicates = append(predicates, sozzler.NewRatingPredicate(rating))
	}

	return predicates, nil
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceP("ingredients", "i", []string{}, "return recipes with ingredients, comma separated")
	listCmd.Flags().StringSliceP("names", "n", []string{}, "return recipes by name, comma separated")
	listCmd.Flags().IntP("rating", "r", 0, "return recipes with this rating or higher")
}
