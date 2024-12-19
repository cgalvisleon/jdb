package jdb

type TypeReferences int

const (
	RelationOneToOne TypeReferences = iota
	RelationOneToMany
	RelationManyToMany
)

type References struct {
	Model           *Model
	To              string
	TypeRelation    TypeReferences
	Key             *Column
	ForeignKey      string
	OnDeleteCascade bool
	OnUpdateCascade bool
}
