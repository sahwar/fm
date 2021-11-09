package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "fm",
	Short:   "FM is a simple, configurable, and fun to use file manager",
	Version: "0.8.0",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		startDir := cmd.Flags().Lookup("start-dir")
		selectionPath := cmd.Flags().Lookup("selection-path")

		config.LoadConfig(startDir, selectionPath)
		cfg := config.GetConfig()

		// If logging is enabled, logs will be output to debug.log.
		if cfg.Settings.EnableLogging {
			f, err := tea.LogToFile("debug.log", "debug")
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			defer func() {
				if err = f.Close(); err != nil {
					log.Fatal(err)
					os.Exit(1)
				}
			}()
		}

		m := ui.NewModel()
		var opts []tea.ProgramOption

		// Always append alt screen program option.
		opts = append(opts, tea.WithAltScreen())

		// If mousewheel is enabled, append it to the program options.
		if cfg.Settings.EnableMouseWheel {
			opts = append(opts, tea.WithMouseAllMotion())
		}

		// Initialize and start app.
		p := tea.NewProgram(m, opts...)
		if err := p.Start(); err != nil {
			log.Fatal("Failed to start fm", err)
			os.Exit(1)
		}
	},
}

// Execute runs the root command and starts the application.
func Execute() {
	rootCmd.PersistentFlags().String("selection-path", "", "Path to write to file on open.")
	rootCmd.PersistentFlags().String("start-dir", "", "Starting directory for FM")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
