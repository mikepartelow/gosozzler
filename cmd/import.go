package cmd

import (
	"bufio"
	"fmt"
	"mp/sozzler/pkg/sozzler"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Parse a recipe from stdin into Sozzler Recipe YAML format and print it to stdout",
	Long: `Parse a recipe from standard input (example follows) and print a YAML representation to stdout.

Example input (between the --- lines):
---
Banana Fabrication

2 oz Rhum agricole
3/4 oz Clement Créole Shrubb
1/2 oz banana juice
1/8 dash bitters
Lime wheel

Shake ingredients with ice. Stir ingredients with dry ice. Correct spelling errors, strain into chilled coconut shell.
---
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return fmt.Errorf("no stdin detected: pipe a recipe into this command, e.g. `cat recipe.txt | gosozzler import`")
		}

		reader := bufio.NewReader(os.Stdin)

		parser := sozzler.RecipeParser{}
		recipe, err := parser.Parse(reader)
		if err != nil {
			return fmt.Errorf("error parsing markdown: %w", err)
		}

		enc := yaml.NewEncoder(os.Stdout)
		// enc.SetIndent("", "  ")
		if err := enc.Encode(recipe); err != nil {
			return fmt.Errorf("encoding JSON: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
