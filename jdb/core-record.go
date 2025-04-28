package jdb

import (
	"database/sql"
	"net/http"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/response"
	"github.com/cgalvisleon/et/timezone"
	"github.com/cgalvisleon/et/utility"
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
	coreRecords.DefineColumn(CREATED_AT, CreatedAtField.TypeData())
	coreRecords.DefineColumn(UPDATED_AT, UpdatedAtField.TypeData())
	coreRecords.DefineColumn("schema_name", TypeDataText)
	coreRecords.DefineColumn("table_name", TypeDataText)
	coreRecords.DefineColumn("option", TypeDataShortText)
	coreRecords.DefineColumn("sync", TypeDataBool)
	coreRecords.DefineColumn(SYSID, SystemKeyField.TypeData())
	coreRecords.DefineIndexField()
	coreRecords.DefinePrimaryKey("schema_name", "table_name", SYSID)
	coreRecords.DefineIndex(true,
		"option",
		"sync",
		SYSID,
		INDEX,
	)
	if err := coreRecords.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

func (s *DB) upsertRecord(tx *sql.Tx, schema, name, sysid, option string) error {
	if sysid == "" {
		return mistake.New(MSG_SYSID_REQUIRED)
	}

	now := timezone.Now()
	item, err := coreRecords.
		Update(et.Json{
			UPDATED_AT: now,
			"option":   option,
			"sync":     false,
		}).
		Where("schema_name").Eq(schema).
		And("table_name").Eq(name).
		And(SYSID).Eq(sysid).
		ExecTx(tx)
	if err != nil {
		return err
	}

	if item.Ok {
		return nil
	}

	_, err = coreRecords.
		Insert(et.Json{
			CREATED_AT:    now,
			UPDATED_AT:    now,
			"schema_name": schema,
			"table_name":  name,
			"option":      option,
			"sync":        false,
			SYSID:         sysid,
			INDEX:         utility.GenIndex(),
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
	result, err := coreRecords.
		Query(query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

/**
* HandlerQueryRecords
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (s *DB) HandlerQueryRecords(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	result, err := s.QueryRecords(body)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.RESULT(w, r, http.StatusOK, result)
}
