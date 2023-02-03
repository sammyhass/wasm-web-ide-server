package cmd

import (
	"github.com/fatih/color"
	"github.com/sammyhass/web-ide/server/env"
	"github.com/sammyhass/web-ide/server/model"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run a database migration",
	Run: func(cmd *cobra.Command, args []string) {
		color.Green("Running migrations")

		env.InitEnv()

		model.Migrate()

		color.Green("Migrations complete")
	},
}
