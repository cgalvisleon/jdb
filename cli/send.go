package cli

import (
	"net"
	"os"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/logs"
	"github.com/spf13/cobra"
)

var useTCP bool

func init() {
	sendCmd.Flags().BoolVar(&useTCP, "tcp", false, "Enviar comando por TCP en vez de Unix socket")
	rootCmd.AddCommand(sendCmd)
}

var sendCmd = &cobra.Command{
	Use:   CMD_SEND,
	Short: CMD_SEND_SHORT,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		message := args[0]
		var conn net.Conn
		var err error

		if useTCP {
			port := envar.GetStr("JDB_PORT", ":8010")
			conn, err = net.Dial("tcp", "127.0.0.1"+port)
		} else {
			conn, err = net.Dial("unix", socketPath)
		}
		if err != nil {
			logs.Logf(PackageName, "Error conectando al daemon: %s", err.Error())
			os.Exit(1)
		}
		defer conn.Close()

		// Enviar mensaje
		_, err = conn.Write([]byte(message))
		if err != nil {
			logs.Logf(PackageName, "Error enviando mensaje: %s", err.Error())
			return
		}

		// Leer respuesta
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		logs.Logf(PackageName, "Respuesta: %s\n", string(buf[:n]))
	},
}
