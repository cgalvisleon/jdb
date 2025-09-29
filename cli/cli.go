package cli

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
)

var (
	PackageName = "jdb"
	cli         *Cli
)

type Cli struct {
	pidFile      string       `json:"-"`
	socketPath   string       `json:"-"`
	logFile      string       `json:"-"`
	dataDir      string       `json:"-"`
	unixListener net.Listener `json:"-"`
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
* startDaemon
**/
func (s *Cli) startDaemon() {
	// Verifica si ya hay un pid file
	if pid, alive := s.checkExistingDaemon(); alive {
		logs.Logf(PackageName, "⚠️ Ya existe un daemon corriendo con PID %d\n", pid)
		return
	}

	// Abre log file
	f, err := os.OpenFile(s.logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logs.Log(PackageName, "Error creating log:", err)
		return
	}

	// Reejecuta binario en modo --child
	cmd := exec.Command(os.Args[0], "run", "--child")
	cmd.Stdout = f
	cmd.Stderr = f
	cmd.Stdin = nil

	if err := cmd.Start(); err != nil {
		logs.Log(PackageName, "Error starting daemon:", err)
		return
	}

	// Guarda el PID en archivo
	if err := os.WriteFile(s.pidFile, []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644); err != nil {
		fmt.Println("Error guardando pid file:", err)
		return
	}

	logs.Logf(PackageName, "Daemon started, PID %d (logs in %s)", cmd.Process.Pid, s.logFile)
}

/**
* runServer
**/
func (s *Cli) runServer() {
	// Limpia socket si ya existe
	if _, err := os.Stat(s.socketPath); err == nil {
		os.Remove(s.socketPath)
	}

	l, err := net.Listen("unix", s.socketPath)
	if err != nil {
		logs.Log(PackageName, "Error starting server:", err)
		return
	}
	s.unixListener = l
	defer l.Close()

	logs.Log(PackageName, "Server started in:", s.socketPath)

	// Captura señales
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			conn, err := s.unixListener.Accept()
			if err != nil {
				continue
			}
			go s.handleConnection(conn)
		}
	}()

	<-sigs
	s.stop()
}

/**
* stop
**/
func (s *Cli) stop() {
	if s.unixListener != nil {
		s.unixListener.Close() // esto desbloquea Accept()
	}
	os.Remove(s.socketPath)
	os.Remove(s.pidFile)
	logs.Logf(PackageName, "Server stopped.")
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

/**
* handleConnection
**/
func (s *Cli) handleConnection(c net.Conn) {
	defer c.Close()

	buf := make([]byte, 1024)
	n, _ := c.Read(buf)
	cmd := string(buf[:n])

	switch cmd {
	case "status":
		c.Write([]byte("✅ Server running\n"))
	case "stop":
		c.Write([]byte("⛔ Stopping server\n"))
		os.Exit(0)
	default:
		c.Write([]byte("❓ Command not recognized\n"))
	}
}

func init() {
	cli = &Cli{
		pidFile:    "/tmp/jdb.pid",
		socketPath: "/tmp/jdb.sock",
		logFile:    "/tmp/jdb.log",
		dataDir:    envar.GetStr("JDB_DATA", "./data"),
	}
}
