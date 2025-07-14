package jdb

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

var coreAudit *Model

func (s *DB) defineAudit() error {
	if s.driver.Name() == SqliteDriver {
		return nil
	}

	if err := s.defineSchema(); err != nil {
		return err
	}

	if coreAudit != nil {
		return nil
	}

	coreAudit = NewModel(coreSchema, "audit", 1)
	coreAudit.DefineColumn(cf.CreatedAt, TypeDataDateTime)
	coreAudit.DefineColumn("command", TypeDataText)
	coreAudit.DefineColumn("query", TypeDataMemo)
	coreAudit.definePrimaryKeyField()
	coreAudit.DefineIndexField()
	coreAudit.DefineIndex(true,
		cf.CreatedAt,
		"command",
	)
	coreAudit.isAudit = true
	if err := coreAudit.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
* audit
* @param command, query string
**/
func audit(command string, query string) {
	if coreAudit == nil || !coreAudit.isInit {
		return
	}

	result := utility.ToBase64(query)
	_, err := coreAudit.Insert(et.Json{
		cf.CreatedAt: utility.Now(),
		cf.Key:       coreAudit.GenId(),
		"command":    command,
		"query":      result,
	}).
		AfterInsert(func(tx *Tx, data et.Json) error {
			count, err := coreAudit.
				Counted()
			if err != nil {
				return err
			}

			limit := envar.GetInt("AUDIT_LIMIT", 10000)
			if count > limit {
				item, err := coreAudit.
					Where("command").Neg("exec").
					OrderBy(cf.Index).
					First(1)
				if err != nil {
					return err
				}

				id := item.Str(0, cf.Key)
				_, err = coreAudit.
					Delete(cf.Key).Eq(id).
					ExecTx(tx)
				if err != nil {
					return err
				}
			}

			return nil
		}).
		Exec()
	if err != nil {
		console.Alert(err)
	}

	debug := envar.Bool("DEBUG")

	if debug {
		console.Debug("Audit:", query)
	}
}
