package cmd

import (
	"context"
	"fmt"
	"mp/sozzler/pkg/display"
	"mp/sozzler/pkg/sozzler"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sozzler",
	Short: "The world's premier cocktail recipe app",
}

func init() {
	var (
		color bool
		plain bool
		tui   bool
	)
	rootCmd.PersistentFlags().BoolVarP(&color, "color", "c", false, "Force color mode when piping")
	rootCmd.PersistentFlags().BoolVarP(&plain, "plain", "p", false, "Plain Text")
	rootCmd.PersistentFlags().BoolVarP(&tui, "tui", "t", false, "Terminal User Interface")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		var catalog sozzler.RecipeCatalog
		if err := catalog.Load("./recipes"); err != nil {
			return err
		}

		var d display.Display = &display.StdoutDisplay{
			Plain: plain,
		}
		if tui {
			d = &display.TuiDisplay{}
		}

		if color {
			lipgloss.SetColorProfile(termenv.TrueColor) // or termenv.ANSI256

		}

		ctx := context.WithValue(cmd.Context(), catalogKey{}, &catalog)
		ctx = context.WithValue(ctx, displayKey{}, d)
		cmd.SetContext(ctx)

		return nil
	}
}

func Execute() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type catalogKey struct{}
type displayKey struct{}
