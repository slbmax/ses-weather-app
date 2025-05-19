package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	root := &cobra.Command{
		Use:   "weather-app",
		Short: "Weather App CLI",
	}

	root.AddCommand(migrateCmd, runCmd)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
