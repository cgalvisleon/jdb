package cli

import (
	"os"
	"syscall"

	"github.com/cgalvisleon/et/logs"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   CMD_STOP,
	Short: CMD_STOP_SHORT,
	Run: func(cmd *cobra.Command, args []string) {
		d := &Cli{pidFile: pidFile}
		pid, alive := d.checkExistingDaemon()
		if !alive {
			logs.Logf(PackageName, "❌ No hay daemon activo")
			return
		}

		syscall.Kill(pid, syscall.SIGTERM)
		os.Remove(d.pidFile)
		logs.Logf(PackageName, "Daemon detenido con PID %d", pid)
	},
}
