package josefina

import (
	"github.com/cgalvisl/jdb/jdb"
	"github.com/cgalvisleon/et/et"
)

const DriverName = "josefina"

var driver Josefina

type Josefina struct {
}

func (s *Josefina) Name() string {
	return DriverName
}

func (s *Josefina) Connect(params et.Json) error {
	return nil
}

func (s *Josefina) Disconnect() error {
	return nil
}

func NewDriver() jdb.Driver {
	return &Josefina{}
}

func init() {
	jdb.Register(DriverName, NewDriver)
}
