package base

import (
	jdb "github.com/cgalvisl/jdb/pkg"
	"github.com/cgalvisleon/et/et"
)

const DriverName = "base"

var driver Base

type Base struct {
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

func NewDriver() jdb.Driver {
	return &Base{}
}

func init() {
	jdb.Register(DriverName, NewDriver)
}
