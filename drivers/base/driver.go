package base

import (
	"github.com/cgalvisleon/et/et"
	jdb "github.com/cgalvisleon/jdb/jdb"
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

func (s *Base) SetKey(key string, value []byte) error {

	return nil
}

func (s *Base) GetKey(key string) (et.KeyValue, error) {
	return et.KeyValue{}, nil
}

func (s *Base) DeleteKey(key string) error {
	return nil
}

func (s *Base) FindKeys(search string, page, rows int) (et.List, error) {
	return et.List{}, nil
}

func init() {
	jdb.Register(DriverName, NewDriver)
}
