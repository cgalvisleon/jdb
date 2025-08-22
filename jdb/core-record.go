package jdb

import (
	"net/http"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/reg"
	"github.com/cgalvisleon/et/response"
	"github.com/cgalvisleon/et/timezone"
)

var coreRecords *Model

func (s *DB) defineRecords() error {
	if err := s.defineSchema(); err != nil {
		return err
	}

	if coreRecords != nil {
		return nil
	}

	coreRecords = NewModel(coreSchema, "records", 1)
	coreRecords.DefineColumn(cf.CreatedAt, TypeDataDateTime)
	coreRecords.DefineColumn(cf.UpdatedAt, TypeDataDateTime)
	coreRecords.DefineColumn("schema_name", TypeDataText)
	coreRecords.DefineColumn("table_name", TypeDataText)
	coreRecords.DefineColumn("option", TypeDataShortText)
	coreRecords.DefineColumn("sync", TypeDataCheckbox)
	coreRecords.DefineColumn(cf.SystemId, TypeDataKey)
	coreRecords.DefineIndexField()
	coreRecords.DefinePrimaryKey("schema_name", "table_name", cf.SystemId)
	coreRecords.DefineIndex(true,
		"option",
		"sync",
		cf.SystemId,
		cf.Index,
	)
	if err := coreRecords.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

func (s *DB) upsertRecord(tx *Tx, schema, name, sysid, option string) error {
	if coreRecords == nil || !coreRecords.isInit || sysid == "" {
		return nil
	}

	now := timezone.Now()
	data := et.Json{
		"schema_name": schema,
		"table_name":  name,
		"option":      option,
		"sync":        false,
		cf.SystemId:   sysid,
	}
	_, err := coreRecords.
		Upsert(data).
		BeforeInsert(func(tx *Tx, data et.Json) error {
			data.Set(cf.CreatedAt, now)
			data.Set(cf.UpdatedAt, now)
			data.Set(cf.Index, reg.GenIndex())
			return nil
		}).
		BeforeUpdate(func(tx *Tx, data et.Json) error {
			data.Set(cf.UpdatedAt, now)
			return nil
		}).
		ExecTx(tx)
	if err != nil {
		return err
	}

	return nil
}

/**
* QueryRecords
* @param query et.Json
* @return interface{}, error
**/
func (s *DB) QueryRecords(query et.Json) (interface{}, error) {
	if coreRecords == nil || !coreRecords.isInit {
		return nil, mistake.New(MSG_DATABASE_NOT_CONCURRENT)
	}

	result, err := coreRecords.
		Query(query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

/**
* HttpQueryRecords
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (s *DB) HttpQueryRecords(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	result, err := s.QueryRecords(body)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.RESULT(w, r, http.StatusOK, result)
}
