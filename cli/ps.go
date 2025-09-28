package cli

import (
	"github.com/cgalvisleon/et/logs"
	"github.com/spf13/cobra"
)

var psCmd = &cobra.Command{
	Use:   CMD_PS,
	Short: CMD_PS_SHORT,
	Run: func(cmd *cobra.Command, args []string) {
		if pid, alive := cli.checkExistingDaemon(); alive {
			logs.Logf(PackageName, "✅ Daemon activo con PID %d\n", pid)
		} else {
			logs.Logf(PackageName, "❌ No hay daemon activo")
		}
	},
}
