package cli

import (
	"os"

	"github.com/cgalvisleon/et/logs"
	"github.com/spf13/cobra"
)

var daemonMode bool

func init() {
	runCmd.Flags().BoolVar(&daemonMode, "daemon", false, "Ejecutar en background")
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   CMD_RUN,
	Short: CMD_RUN_SHORT,
	Run: func(cmd *cobra.Command, args []string) {
		cli = newCli()
		if pid, alive := cli.checkExistingDaemon(); alive {
			logs.Logf(PackageName, "⚠️ Ya existe un daemon con PID %d\n", pid)
			return
		}

		// Background
		if daemonMode {
			attr := &os.ProcAttr{
				Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
			}
			proc, err := os.StartProcess(os.Args[0], []string{os.Args[0], "run"}, attr)
			if err != nil {
				logs.Logf(PackageName, "Error: %s", err.Error())
				return
			}
			logs.Logf(PackageName, "Daemon iniciado con PID %d", proc.Pid)
			return
		}

		// Foreground
		if err := cli.savePID(); err != nil {
			logs.Logf(PackageName, "Error guardando pid: %s", err.Error())
			return
		}

		cli.runServer()
	},
}
