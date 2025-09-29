package cli

import (
	"github.com/cgalvisleon/et/logs"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(psCmd)
}

var psCmd = &cobra.Command{
	Use:   CMD_PS,
	Short: CMD_PS_SHORT,
	Run: func(cmd *cobra.Command, args []string) {
		d := &Cli{pidFile: pidFile}
		if pid, alive := d.checkExistingDaemon(); alive {
			logs.Logf(PackageName, "✅ Daemon activo con PID %d\n", pid)
		} else {
			logs.Logf(PackageName, "❌ No hay daemon activo")
		}
	},
}
