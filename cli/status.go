package cli

import (
	"net"

	"github.com/cgalvisleon/et/logs"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   CMD_STATUS,
	Short: CMD_STATUS_SHORT,
	Run: func(cmd *cobra.Command, args []string) {
		d := &Cli{socketPath: socketPath}
		conn, err := net.Dial("unix", d.socketPath)
		if err != nil {
			logs.Logf(PackageName, "❌ No hay daemon activo")
			return
		}
		defer conn.Close()

		conn.Write([]byte("status"))
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		logs.Log(PackageName, string(buf[:n]))
	},
}
