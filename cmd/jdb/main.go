package main

import (
	"os"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
)

const PackageName = "jdb"

type TypeCommand string

const (
	CMD_Help    TypeCommand = "help"
	CMD_Version TypeCommand = "version"
	CMD_Status  TypeCommand = "status"
	CMD_Start   TypeCommand = "start"
	CMD_Stop    TypeCommand = "stop"
	CMD_Restart TypeCommand = "restart"
	CMD_Conf    TypeCommand = "conf"
)

func ToTypeCommand(val string) TypeCommand {
	switch val {
	case "version", "--v":
		return CMD_Version
	case "status", "--s":
		return CMD_Status
	case "start", "--start":
		return CMD_Start
	case "stop", "--stop":
		return CMD_Stop
	case "restart", "--restart":
		return CMD_Restart
	case "conf", "--conf":
		return CMD_Conf
	}

	return CMD_Help
}

type RepositoryCMD interface {
	Help(key string)
	Version() string
	SetConfig(cfg string)
	Status() et.Item
	Start() et.Item
	Stop() et.Item
	Restart() et.Item
}

var apps map[string]RepositoryCMD
var app string = "systemd"

func main() {
	_app := apps[app]

	if len(os.Args) < 2 {
		_app.Help("")
		return
	}

	command := os.Args[1]
	switch ToTypeCommand(command) {
	case CMD_Version:
		_app.Version()
	case CMD_Help:
		if len(os.Args) > 2 {
			_app.Help(os.Args[2])
		} else {
			_app.Help("")
		}
	case CMD_Conf:
		if len(os.Args) > 2 {
			_app.SetConfig(os.Args[2])
		}
		println("Configuración: Sin parametros")
	case CMD_Status:
		result := _app.Status()
		if result.Ok {
			logs.Log(PackageName, result.Str("message"))
		} else {
			logs.Alertm(result.Str("message"))
		}
	case CMD_Start:
		result := _app.Start()
		logs.Debug(result.ToString())
		message := result.Str("message")
		if result.Ok {
			logs.Log(PackageName, result.Str("message"))
		} else {
			logs.Alertm(message)
		}
	case CMD_Stop:
		result := _app.Stop()
		if result.Ok {
			logs.Log(PackageName, result.Str("message"))
		} else {
			logs.Alertm(result.Str("message"))
		}
	case CMD_Restart:
		result := _app.Restart()
		if result.Ok {
			logs.Log(PackageName, result.Str("message"))
		} else {
			logs.Alertm(result.Str("message"))
		}
	}
}

func Registry(name string, cmd RepositoryCMD) {
	if apps == nil {
		apps = make(map[string]RepositoryCMD)
	}

	apps[name] = cmd
	app = name
}
