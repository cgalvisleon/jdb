package jdb

import (
	"strconv"
	"strings"

	"github.com/cgalvisleon/et/strs"
)

type Field struct {
	Column *Column
	Schema string
	Table  string
	As     string
	Name   string
	Atrib  string
	Index  int
	Value  interface{}
}

/**
* TableName
* @return string
**/
func (s *Field) TableName() string {
	return strs.Format("%s.%s", s.Table, s.Name)
}

/**
* Tag
* @return string
**/
func (s *Field) Tag() string {
	result := ""
	result = strs.Append(result, s.Schema, "")
	result = strs.Append(result, s.Table, ".")
	result = strs.Append(result, s.Name, ".")
	result = strs.Append(result, s.Atrib, ".")

	return result
}

/**
* Caption
* @return string
**/
func (s *Field) Caption() string {
	if len(s.Atrib) == 0 {
		return s.Name
	}

	return s.Atrib
}

func NewField(name string) *Field {
	name = strs.Lowcase(name)
	list := strings.Split(name, ".")

	if len(list) == 1 {
		return &Field{
			Schema: "",
			Table:  "",
			As:     "",
			Name:   list[0],
			Atrib:  "",
			Index:  0,
		}
	}

	if len(list) == 2 {
		return &Field{
			Schema: "",
			Table:  list[0],
			As:     "",
			Name:   list[1],
			Atrib:  "",
			Index:  0,
		}
	}

	if len(list) == 3 {
		result := &Field{
			Schema: list[0],
			Table:  list[1],
			As:     "",
			Name:   list[2],
			Atrib:  "",
			Index:  0,
		}
		subList := strings.Split(list[1], ":")
		if len(subList) == 2 {
			result.Table = subList[0]
			result.As = subList[1]
		}

		return result
	}

	if len(list) == 4 {
		result := &Field{
			Schema: list[0],
			Table:  list[1],
			As:     "",
			Name:   list[2],
			Atrib:  list[3],
			Index:  0,
		}

		subList := strings.Split(list[1], ":")
		if len(subList) == 2 {
			result.Table = subList[0]
			result.As = subList[1]
		}

		subList = strings.Split(list[3], ":")
		if len(subList) == 2 {
			result.Atrib = subList[0]
			idx, err := strconv.Atoi(subList[1])
			if err != nil {
				idx = 0
			}
			result.Index = idx
		}

		return result
	}

	return nil
}
