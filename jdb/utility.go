package jdb

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
)

/**
* RowsToItems return a items from a sql rows
* @param rows *sql.Rows
* @return et.Items
**/
func RowsToItems(rows *sql.Rows) et.Items {
	var result = et.Items{Result: []et.Json{}}
	for rows.Next() {
		var item et.Json
		item.ScanRows(rows)

		result.Ok = true
		result.Count++
		result.Result = append(result.Result, item)
	}

	return result
}
