package jdb

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/timezone"
)

var coreSeries *Model

func (s *DB) defineSeries() error {
	if err := s.defineSchema(); err != nil {
		return err
	}

	if coreSeries != nil {
		return nil
	}

	coreSeries = NewModel(coreSchema, "series", 1)
	coreSeries.DefineAtribute(CREATED_AT, CreatedAtField.TypeData())
	coreSeries.DefineAtribute(UPDATED_AT, UpdatedAtField.TypeData())
	coreSeries.DefineAtribute("tag", TypeDataText)
	coreSeries.DefineAtribute("value", TypeDataInt)
	coreSeries.DefineAtribute(INDEX, IndexField.TypeData())
	coreSeries.DefinePrimaryKey("tag")
	coreSeries.DefineIndex(true,
		"tag",
		INDEX,
	)
	if err := coreSeries.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
* CurrentSerie
* @param tag string
* @return int64
**/
func (s *DB) CurrentSerie(tag string) int64 {
	if !s.UseCore {
		return 0
	}

	item, err := coreSeries.
		Where("tag").Eq(tag).
		One()
	if err != nil {
		return 0
	}

	if !item.Ok {
		return 0
	}

	return item.Int64("value")
}

/**
* GetSerie
* @param tag string
* @return int64
**/
func (s *DB) GetSerie(tag string) int64 {
	item, err := coreSeries.
		Update(et.Json{
			UPDATED_AT: timezone.Now(),
			"value":    "VALUE + 1",
		}).
		Where("tag").Eq(tag).
		Return("value").
		One()
	if err != nil {
		return 0
	}

	if item.Ok {
		return item.Int64("value")
	}

	now := timezone.Now()
	_, err = coreSeries.
		Insert(et.Json{
			CREATED_AT: now,
			UPDATED_AT: now,
			"tag":      tag,
			"value":    1,
		}).
		Return("value").
		One()
	if err != nil {
		return 0
	}

	return 1
}

/**
* NextCode
* @param tag, prefix string
* @return string
**/
func (s *DB) NextCode(tag, prefix string) string {
	num := s.GetSerie(tag)

	if len(prefix) == 0 {
		return strs.Format("%08v", num)
	} else {
		return strs.Format("%s%08v", prefix, num)
	}
}

/**
* SetSerie
* @param tag string, val int64
* @return int64
**/
func (s *DB) SetSerie(tag string, val int64) int64 {
	now := timezone.Now()
	item, err := coreSeries.
		Update(et.Json{
			UPDATED_AT: now,
			"value":    val,
		}).
		Where("tag").Eq(tag).
		Return("value").
		One()
	if err != nil {
		return 0
	}

	if item.Ok {
		return item.Int64("value")
	}

	_, err = coreSeries.
		Insert(et.Json{
			CREATED_AT: now,
			UPDATED_AT: now,
			"tag":      tag,
			"value":    val,
		}).
		Return("value").
		One()
	if err != nil {
		return 0
	}

	return val
}
