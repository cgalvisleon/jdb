package cli

import (
	"net"
	"strings"

	"github.com/cgalvisleon/et/logs"
)

/**
* handleConnection
* @param c net.Conn, proto string
**/
func (s *Cli) handleConnection(c net.Conn, proto string) {
	defer c.Close()
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil {
		return
	}

	msg := strings.TrimSpace(string(buf[:n]))
	logs.Logf(PackageName, "[%s] Mensaje recibido: %s", proto, msg)

	parts := strings.SplitN(msg, " ", 2)
	cmd := parts[0]
	args := ""
	if len(parts) > 1 {
		args = parts[1]
	}

	if handler, ok := s.commands[cmd]; ok {
		resp := handler(args)
		c.Write([]byte(resp + "\n"))
	} else {
		c.Write([]byte("Comando no reconocido\n"))
	}
}
