package cli

import (
	"net"

	"github.com/cgalvisleon/et/logs"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   CMD_STOP,
	Short: CMD_STOP_SHORT,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := net.Dial("unix", cli.socketPath)
		if err != nil {
			logs.Log(PackageName, "No se pudo conectar al servidor:", err)
			return
		}
		defer conn.Close()

		conn.Write([]byte("stop"))
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		logs.Log(PackageName, string(buf[:n]))
	},
}
