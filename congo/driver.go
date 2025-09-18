package jdb

const (
	DriverPostgres = "postgres"
	DriverMysql    = "mysql"
	DriverSqlite   = "sqlite"
	DriverMssql    = "mssql"
	DriverOracle   = "oracle"
)

type Driver interface {
	Load(model *Model) error
}

var drivers map[string]Driver

func init() {
	drivers = make(map[string]Driver)
}

func Register(name string, driver Driver) {
	drivers[name] = driver
}
