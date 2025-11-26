package postgres

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	"github.com/cgalvisleon/jdb/jdb"
)

/**
* buildModel
* @param model *jdb.Model
* @return (string, error)
**/
func (s *Driver) buildModel(model *jdb.Model) (string, error) {
	definition := model.ToJson()
	sql, err := s.buildSchema(definition)
	if err != nil {
		return "", err
	}

	def, err := s.buildTable(definition)
	if err != nil {
		return "", err
	}

	if def != "" {
		sql = strs.Append(sql, def, "\n")
	}

	def, err = s.buildForeignKeys(model, definition)
	if err != nil {
		return "", err
	}

	if def != "" {
		sql = strs.Append(sql, def, "\n")
	}

	def, err = s.buildIndexes(definition)
	if err != nil {
		return "", err
	}

	if def != "" {
		sql = strs.Append(sql, def, "\n")
	}

	def, err = s.buildUniqueIndex(definition)
	if err != nil {
		return "", err
	}

	if def != "" {
		sql = strs.Append(sql, def, "\n")
	}

	def, err = s.buildTriggerBeforeInsert(definition)
	if err != nil {
		return "", err
	}

	if def != "" {
		sql = strs.Append(sql, def, "\n")
	}

	return sql, nil
}

/**
* buildSchema
* @param definition et.Json
* @return (string, error)
**/
func (s *Driver) buildSchema(definition et.Json) (string, error) {
	schema := definition.String("schema")
	if !utility.ValidStr(schema, 0, []string{}) {
		return "", fmt.Errorf(jdb.MSG_SCHEMA_REQUIRED)
	}

	exist, err := ExistSchema(s.database.Db, schema)
	if err != nil {
		return "", err
	}

	if exist {
		return "", nil
	}

	return fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", schema), nil
}

/**
* buildTable
* @param definition et.Json
* @return (string, error)
**/
func (s *Driver) buildTable(definition et.Json) (string, error) {
	getType := func(tp string) string {
		switch tp {
		case "int":
			return "BIGINT"
		case "float":
			return "DOUBLE PRECISION"
		case "key":
			return "VARCHAR(80)"
		case "text":
			return "VARCHAR(250)"
		case "memo":
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

	defaultValue := func(tp string) string {
		switch tp {
		case "int":
			return "0"
		case "float":
			return "0.0"
		case "key":
			return "''"
		case "text":
			return "''"
		case "memo":
			return "''"
		case "datetime":
			return "NOW()"
		case "boolean":
			return "FALSE"
		case "json":
			return "'{}'"
		case "index":
			return "0"
		case "bytes":
			return "''"
		case "geometry":
			return "'{}'"
		default:
			return ""
		}
	}

	columns := definition.ArrayJson("columns")
	columnsDef := ""
	for _, v := range columns {
		tpVal := v.String("type")
		tp := tpVal
		if !jdb.TypeColumn[tpVal] {
			continue
		}
		tp = getType(tp)
		df := defaultValue(tpVal)
		def := fmt.Sprintf("\n\t%s %s DEFAULT %s", v.String("name"), tp, df)
		columnsDef = strs.Append(columnsDef, def, ",")
	}

	def, err := s.buildPrimaryKeys(definition)
	if err != nil {
		return "", err
	}

	if def != "" {
		columnsDef = strs.Append(columnsDef, def, ",\n\t")
	}

	table := definition.String("table")
	result := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", table, columnsDef)

	return result, nil
}

/**
* buildPrimaryKeys
* @param definition et.Json
* @return (string, error)
**/
func (s *Driver) buildPrimaryKeys(definition et.Json) (string, error) {
	primaryKeys := definition.ArrayStr("primary_keys")
	if len(primaryKeys) == 0 {
		return "", nil
	}

	columns := ""
	for _, v := range primaryKeys {
		columns = strs.Append(columns, v, ", ")
	}

	name := definition.String("name")
	result := fmt.Sprintf("CONSTRAINT pk_%s PRIMARY KEY (%s)", name, columns)

	return result, nil
}

/**
* buildForeignKeys
* @param mode *jdb.Model, definition et.Json
* @return (string, error)
**/
func (s *Driver) buildForeignKeys(mode *jdb.Model, definition et.Json) (string, error) {
	result := ""
	foreignKeys := definition.ArrayJson("foreign_keys")
	if len(foreignKeys) == 0 {
		return "", nil
	}

	for _, foreignKey := range foreignKeys {
		schema := foreignKey.String("schema")
		name := foreignKey.String("name")
		to := fmt.Sprintf("%s.%s", schema, name)
		references := foreignKey.Json("references")
		columns := references.ArrayJson("columns")
		fks := ""
		ks := ""
		for _, fk := range columns {
			for k, v := range fk {
				fks = strs.Append(fks, fmt.Sprintf("%s", k), ", ")
				ks = strs.Append(ks, fmt.Sprintf("%s", v), ", ")
			}
		}
		def := fmt.Sprintf("ALTER TABLE IF EXISTS %s ADD CONSTRAINT fk_%s_%s FOREIGN KEY(%s) REFERENCES %s(%s)", mode.Table, mode.Name, name, fks, to, ks)
		onDelete := references.String("on_delete")
		onUpdate := references.String("on_update")
		if onDelete != "" {
			def = strs.Append(def, fmt.Sprintf("ON DELETE %s", onDelete), " ")
		}
		if onUpdate != "" {
			def = strs.Append(def, fmt.Sprintf("ON UPDATE %s", onUpdate), " ")
		}
		def += ";"
		result = strs.Append(result, def, "\n")
	}

	return result, nil
}

/**
* buildIndexes
* @param definition et.Json
* @return (string, error)
**/
func (s *Driver) buildIndexes(definition et.Json) (string, error) {
	indexes := definition.ArrayStr("indexes")
	if len(indexes) == 0 {
		return "", nil
	}

	table := definition.String("table")
	name := definition.String("name")
	result := ""
	for _, v := range indexes {
		def := fmt.Sprintf("idx_%s_%s", name, v)
		if v == jdb.SOURCE {
			def = fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s USING GIN (%s);", def, table, v)
		} else {
			def = fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s(%s);", def, table, v)
		}
		result = strs.Append(result, def, "\n")
	}

	return result, nil
}

/**
* buildUniqueIndex
* @param definition et.Json
* @return (string, error)
**/
func (s *Driver) buildUniqueIndex(definition et.Json) (string, error) {
	uniqueIndexes := definition.ArrayStr("unique_indexes")
	if len(uniqueIndexes) == 0 {
		return "", nil
	}

	table := definition.String("table")
	name := definition.String("name")
	result := ""
	for _, v := range uniqueIndexes {
		def := fmt.Sprintf("idx_%s_%s", name, v)
		def = fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS %s ON %s(%s);", def, table, v)
		result = strs.Append(result, def, "\n")
	}

	return result, nil
}

/**
* buildTriggerBeforeInsert
* @param definition et.Json
* @return (string, error)
**/
func (s *Driver) buildTriggerBeforeInsert(definition et.Json) (string, error) {
	recordField := definition.String("record_field")
	if recordField == "" {
		return "", nil
	}

	isCore := definition.Bool("is_core")
	if isCore {
		return "", nil
	}

	table := definition.String("table")
	result := fmt.Sprintf(`
	DROP TRIGGER IF EXISTS RECORDS_SET ON %s CASCADE;
	CREATE TRIGGER RECORDS_SET
	AFTER INSERT OR UPDATE OR DELETE ON %s
	FOR EACH ROW
	EXECUTE FUNCTION core.after_records();
	`, table, table)

	return result, nil
}

/**
* mutateModel
* @param model *jdb.Model
* @return (string, error)
**/
func (s *Driver) mutateModel(model *jdb.Model) (string, error) {
	return "", nil
}
