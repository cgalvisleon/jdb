package cli

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
)

var (
	PackageName        = "jdb"
	pidFile     string = "/tmp/jdb.pid"
	socketPath  string = "/tmp/jdb.sock"
	logFile     string = "/tmp/jdb.log"
)

type CommandHandler func(args string) string

type Cli struct {
	pidFile      string                    `json:"-"`
	socketPath   string                    `json:"-"`
	logFile      string                    `json:"-"`
	dataDir      string                    `json:"-"`
	tcpAddr      string                    `json:"-"`
	unixListener net.Listener              `json:"-"`
	tcpListener  net.Listener              `json:"-"`
	commands     map[string]CommandHandler `json:"-"`
}

/**
* newCli
* @return *Cli
**/
func newCli() *Cli {
	return &Cli{
		pidFile:    pidFile,
		socketPath: socketPath,
		logFile:    logFile,
		dataDir:    envar.GetStr("JDB_DATA", "./data"),
		tcpAddr:    envar.GetStr("JDB_PORT", ":8010"),
		commands:   make(map[string]CommandHandler),
	}
}

/**
* RegisterCommand
* @param name string, handler CommandHandler
**/
func RegisterCommand(name string, handler CommandHandler) {
	if cli.commands == nil {
		cli.commands = make(map[string]CommandHandler)
	}

	cli.commands[name] = handler
}

/**
* serialize
* @return ([]byte, error)
**/
func (s *Cli) serialize() ([]byte, error) {
	bt, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return bt, nil
}

/**
* ToJson
* @return et.Json
**/
func (s *Cli) ToJson() et.Json {
	bt, err := s.serialize()
	if err != nil {
		return et.Json{}
	}

	var result et.Json
	err = json.Unmarshal(bt, &result)
	if err != nil {
		return et.Json{}
	}

	return result
}

/**
* runServer
**/
func (s *Cli) runServer() {
	// --- Unix socket ---
	if s.socketPath != "" {
		if _, err := os.Stat(s.socketPath); err == nil {
			os.Remove(s.socketPath)
		}
		l, err := net.Listen("unix", s.socketPath)
		if err != nil {
			logs.Log(PackageName, "Error starting unix server:", err)
			return
		}
		s.unixListener = l
		logs.Logf(PackageName, "Unix server started in: %s", s.socketPath)
		go s.acceptLoop(s.unixListener, "unix")
	}

	// --- TCP socket ---
	if s.tcpAddr != "" {
		l, err := net.Listen("tcp", s.tcpAddr)
		if err != nil {
			logs.Log(PackageName, "Error starting tcp server:", err)
			return
		}
		s.tcpListener = l
		logs.Logf(PackageName, "TCP server started in: %s", s.tcpAddr)
		go s.acceptLoop(s.tcpListener, "tcp")
	}

	// Manejo de señales para detener
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	s.stop()
}

/**
* acceptLoop
* @param l net.Listener
* @param proto string
**/
func (d *Cli) acceptLoop(l net.Listener, proto string) {
	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		go d.handleConnection(conn, proto)
	}
}

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

/**
* stop
**/
func (s *Cli) stop() {
	if s.unixListener != nil {
		s.unixListener.Close()
		os.Remove(s.socketPath)
	}

	if s.tcpListener != nil {
		s.tcpListener.Close()
	}

	os.Remove(s.pidFile)
	logs.Logf(PackageName, "Server stopped.")
}

/**
* savePID
**/
func (s *Cli) savePID() error {
	return os.WriteFile(s.pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
}

/**
* checkExistingDaemon
**/
func (s *Cli) checkExistingDaemon() (int, bool) {
	data, err := os.ReadFile(s.pidFile)
	if err != nil {
		return 0, false // no existe pid file
	}

	var pid int
	fmt.Sscanf(string(data), "%d", &pid)

	// Verifica si proceso sigue vivo
	if err := syscall.Kill(pid, 0); err == nil {
		return pid, true
	}

	// Proceso muerto → limpiar pid file
	os.Remove(s.pidFile)
	return 0, false
}
