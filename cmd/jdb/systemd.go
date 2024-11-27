package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
	"github.com/go-chi/chi/v5"
)

const (
	ServiceName = "daemon"
)

type Systemd struct {
	serviceName string
	port        int
	server      *http.Server
	isRunning   bool
	wg          sync.WaitGroup
}

func New() RepositoryCMD {
	return &Systemd{
		serviceName: "systemd",
		isRunning:   false,
		wg:          sync.WaitGroup{},
	}
}

func (s *Systemd) Version() string {
	result := "Version: 1.0.0"
	println(result)

	return result
}

func (s *Systemd) Help(key string) {
	if key == "" {
		println("Uso: daemon [opciones]")
		println("Opciones:")
		println("  --h, --help     Mostrar esta ayuda")
		println("  --v, --version  Mostrar la versión")
		println("  --s, --status   Mostrar el estado del servicio")
		println("  --conf          Configurar el servicio")
		println("  --start   			 Iniciar el servicio")
		println("  --stop   			 Detener el servicio")
		println("  --r, --restart  Reiniciar el servicio")
	}
}

func (s *Systemd) SetConfig(cfg string) {
	if cfg == "" {
		logs.Errorm(MSG_CONFIG_REQUIRED)
		return
	}

	config, err := et.Object(cfg)
	if err != nil {
		logs.Alert(err)
		return
	}

	for k, v := range config {
		k = strs.Uppcase(k)
		switch val := v.(type) {
		case string:
			envar.UpSetStr(k, val)
		case int:
			envar.UpSetInt(k, val)
		case bool:
			envar.UpSetBool(k, val)
		case float64:
			envar.UpSetFloat(k, val)
		}
	}
}

func (s *Systemd) Status() et.Item {
	pid, err := readPID()
	if err != nil {
		return et.Item{
			Ok:     false,
			Result: et.Json{"message": err.Error()},
		}
	}

	return et.Item{
		Ok:     true,
		Result: et.Json{"message": fmt.Sprintf("El servicio está en ejecución con el PID %d", pid)},
	}
}

func (s *Systemd) Start() et.Item {
	if s.port == 0 {
		return et.Item{
			Ok:     false,
			Result: et.Json{"message": MSG_PORT_REQUIRED},
		}
	}

	if s.isRunning {
		return et.Item{
			Ok:     true,
			Result: et.Json{"message": MSG_SERVICE_RUNNING},
		}
	}

	err := writePID()
	if err != nil {
		return et.Item{
			Ok:     false,
			Result: et.Json{"message": err.Error()},
		}
	}

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	s.server = &http.Server{
		Addr:    strs.Format(`:%d`, s.port),
		Handler: r,
	}

	s.wg.Add(1)
	s.isRunning = true
	go func() {
		defer s.wg.Done()
		logs.Logf(s.serviceName, `Iniciando el servidor en http://localhost:%d`, s.port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.isRunning = false
			logs.Alert(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		s.Stop()
	}()

	return et.Item{
		Ok:     true,
		Result: et.Json{"message": MSG_SERVICE_RUNNING},
	}
}

func (s *Systemd) Stop() et.Item {
	if !s.isRunning {
		return et.Item{
			Ok:     true,
			Result: et.Json{"message": MSG_SERVICE_CLOSED},
		}
	}

	fmt.Println("Deteniendo el servidor...")
	if err := s.server.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
	s.wg.Wait()
	s.isRunning = false
	logs.Log(ServiceName, "Servidor detenido.")

	return et.Item{
		Ok:     true,
		Result: et.Json{"message": MSG_SERVICE_CLOSED},
	}
}

func (s *Systemd) Restart() et.Item {
	return et.Item{}
}

func init() {
	Registry("systemd", New())
}
