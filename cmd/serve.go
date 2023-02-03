package cmd

import (
	"github.com/fatih/color"
	"github.com/sammyhass/web-ide/server/db"
	"github.com/sammyhass/web-ide/server/env"
	"github.com/sammyhass/web-ide/server/router"
	"github.com/sammyhass/web-ide/server/s3"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("port", "p", "8080", "Port to run the server on")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Long:  "Start serving the API on the specified port",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		startServer(port)
	},
}

func startServer(port string) {
	color.Green("Starting server on port %s", port)

	env.InitEnv()

	db.Connect()
	defer db.Close()

	s3.InitSession()

	if port != "" {
		env.Set(env.PORT, port)
	}

	router.Run(env.Get(env.PORT))
}
