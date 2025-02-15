package jdb

func example() {
	ejemplo := NewModel(nil, "ejemplo", 1)
	From(ejemplo).
		Where("table1.id").Eq(1).
		Select("A.id", "B.id").
		And("table2.id").Eq(2).
		All()

}
