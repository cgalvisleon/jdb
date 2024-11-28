package jdb

type TypeRelation int

const (
	RelationOneToOne TypeRelation = iota
	RelationOneToMany
	RelationManyToMany
)

type Relation struct {
	Model           *Model
	To              *Model
	TypeRelation    TypeRelation
	ForeignKey      map[string]*Column
	ForeignKeyTo    map[string]*Column
	OnDeleteCascade bool
	OnUpdateCascade bool
}
