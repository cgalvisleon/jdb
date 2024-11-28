package base

import (
	"github.com/cgalvisleon/et/et"
	jdb "github.com/cgalvisleon/jdb/pkg"
)

const DriverName = "base"

var driver Base

type Base struct {
}

func NewDriver() jdb.Driver {
	return &Base{}
}

func (s *Base) Name() string {
	return DriverName
}

func (s *Base) Connect(params et.Json) error {
	return nil
}

func (s *Base) Disconnect() error {
	return nil
}

func (s *Base) SetMain(params et.Json) error {

	return nil
}

func init() {
	jdb.Register(DriverName, NewDriver)
}
