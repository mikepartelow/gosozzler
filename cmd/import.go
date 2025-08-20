package cmd

import (
	"bufio"
	"fmt"
	"mp/sozzler/pkg/sozzler"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Read a Markdown recipe from stdin and print JSON",
	Long: `Reads Markdown from standard input, expects a format like EXHIBIT A,
and parses it into the JSON schema shown in EXHIBIT B (Recipe, Components, etc).
Outputs the resulting JSON to stdout.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return fmt.Errorf("no stdin detected: pipe Markdown into this command, e.g. `cat recipe.md | gosozzler import`")
		}

		reader := bufio.NewReader(os.Stdin)

		parser := sozzler.MarkdownParser{}
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
