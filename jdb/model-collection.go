package jdb

/**
* MakeCollection
* @param schema *Schema
* @param name string
* @return *Model
**/
func (s *Model) MakeCollection() *Model {
	s.DefineColumn(CreatedAtField.Str(), CreatedAtField.TypeData())
	s.DefineColumn(UpdatedAtField.Str(), UpdatedAtField.TypeData())
	s.DefineColumn(ProjectField.Str(), ProjectField.TypeData())
	s.DefineColumn(StateField.Str(), StateField.TypeData())
	s.DefineColumn(KeyField.Str(), KeyField.TypeData())
	s.DefineColumn(ClassField.Str(), ClassField.TypeData())
	s.DefineColumn(SourceField.Str(), SourceField.TypeData())
	s.DefineColumn(SystemKeyField.Str(), SystemKeyField.TypeData())
	s.DefineColumn(IndexField.Str(), IndexField.TypeData())
	s.DefineKey(KeyField.Str())
	s.DefineIndex(true,
		CreatedAtField.Str(),
		UpdatedAtField.Str(),
		ProjectField.Str(),
		StateField.Str(),
		ClassField.Str(),
		SourceField.Str(),
		SystemKeyField.Str(),
		IndexField.Str(),
	)

	return s
}
