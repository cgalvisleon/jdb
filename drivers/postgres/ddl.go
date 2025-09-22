package postgres

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

/**
* buildModel
* @param definition et.Json
* @return (string, error)
**/
func (s *Postgres) buildModel(definition et.Json) (string, error) {
	getType := func(tp string) string {
		switch tp {
		case "int":
			return "BIGINT"
		case "float":
			return "DOUBLE PRECISION"
		case "key":
			return "VARCHAR(80)"
		case "text":
			return "TEXT"
		case "datetime":
			return "TIMESTAMP"
		case "boolean":
			return "BOOLEAN"
		case "json":
			return "JSONB"
		case "index":
			return "BIGINT"
		case "bytes":
			return "BYTEA"
		case "geometry":
			return "JSONB"
		default:
			return tp
		}
	}

	columns := definition.ArrayJson("columns")
	columnsDef := ""
	for _, v := range columns {
		tp := v.String("type")
		if !jdb.TypeColumn[tp] {
			continue
		}
		tp = getType(tp)
		def := fmt.Sprintf("\n\t%s %s", v.String("name"), tp)
		columnsDef = strs.Append(columnsDef, def, ",")
	}

	schema := definition.String("schema")
	table := definition.String("table")
	result := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;\n\tCREATE TABLE IF NOT EXISTS %s (%s);", schema, table, columnsDef)

	return result, nil
}
