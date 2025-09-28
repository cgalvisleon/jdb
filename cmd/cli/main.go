package main

import (
	"os"

	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/jdb/cli"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   cli.CMD_JDB,
		Short: cli.CMD_JDB_SHORT,
	}

	// Añadimos los comandos
	cli.Load(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		logs.Log(cli.PackageName, "Error:", err)
		os.Exit(1)
	}
}
