package main

import (
	"log"
	"sczg/commands"
	"sczg/config"

	"github.com/spf13/cobra"
)

func main() {
	var settestDB bool
	var Env config.Env
	var rootCmd = &cobra.Command{Use: "app"}

	var setupDBCmd = &cobra.Command{
		Use:   "setupdb",
		Short: "create database",
		Run: func(cmd *cobra.Command, args []string) {
			commands.SetDefaultDB(settestDB)
		},
	}
	setupDBCmd.Flags().BoolVarP(&settestDB, "testdb", "t", false, "create testdb")

	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "start app server",
		PreRun: func(cmd *cobra.Command, args []string) {
			config.SetupEnv(&Env)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if Env.Err != nil {
				log.Fatal(Env.Err)
			}
			commands.StartServer(&Env)
		},
	}
	serveCmd.Flags().StringVar(&Env.CfgPath, "config", "", "server config path")
	serveCmd.Flags().StringVar(&Env.DbPath, "db", "", "database path")
	serveCmd.MarkFlagRequired("config")
	serveCmd.MarkFlagRequired("db")

	var fetchCmd = &cobra.Command{
		Use:   "fetch",
		Short: "fetch new ad data",
		PreRun: func(cmd *cobra.Command, args []string) {
			config.SetupEnv(&Env)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if Env.Err != nil {
				log.Fatal(Env.Err)
			}
			commands.StartFetcher(&Env)
		},
	}
	fetchCmd.Flags().StringVar(&Env.CfgPath, "config", "", "server config path")
	fetchCmd.Flags().StringVar(&Env.DbPath, "db", "", "database path")
	fetchCmd.MarkFlagRequired("config")
	fetchCmd.MarkFlagRequired("db")

	var archiveCmd = &cobra.Command{
		Use:   "archive",
		Short: "archive old ad data",
		PreRun: func(cmd *cobra.Command, args []string) {
			config.SetupEnv(&Env)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if Env.Err != nil {
				log.Fatal(Env.Err)
			}
			commands.StartArchiver(&Env)
		},
	}
	archiveCmd.Flags().StringVar(&Env.CfgPath, "config", "", "server config path")
	archiveCmd.Flags().StringVar(&Env.DbPath, "db", "", "database path")
	archiveCmd.MarkFlagRequired("config")
	archiveCmd.MarkFlagRequired("db")

	rootCmd.AddCommand(setupDBCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(fetchCmd)

	rootCmd.Execute()
}
