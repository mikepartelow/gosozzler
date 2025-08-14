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
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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

		ctx := context.WithValue(cmd.Context(), "catalog", &catalog)
		ctx = context.WithValue(ctx, "display", d)
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
