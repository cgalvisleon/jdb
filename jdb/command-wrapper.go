package jdb

import "github.com/dop251/goja"

type CommandWrapper struct {
	vm      *goja.Runtime
	Command *Command
}
