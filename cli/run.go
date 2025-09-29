package cli

import (
	"github.com/spf13/cobra"
)

func Load(rootCmd *cobra.Command) {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(psCmd)
}

var runCmd = &cobra.Command{
	Use:   CMD_RUN,
	Short: CMD_RUN_SHORT,
	Run: func(cmd *cobra.Command, args []string) {
		daemonFlag, _ := cmd.Flags().GetBool("daemon")
		childFlag, _ := cmd.Flags().GetBool("child")

		if daemonFlag && !childFlag {
			cli.startDaemon()
			return
		}

		cli.runServer()
	},
}

func init() {
	runCmd.Flags().Bool("daemon", false, "Run in background (daemon)")
	runCmd.Flags().Bool("child", false, "Internal mode (do not use directly)")
}
