package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "jdb",
	Short: "jdb CLI",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
