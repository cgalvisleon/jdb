package jdb

type Driver interface {
	Load(model *Model) (string, error)
}

var drivers map[string]Driver

func init() {
	drivers = make(map[string]Driver)
}

func Register(name string, driver Driver) {
	drivers[name] = driver
}
