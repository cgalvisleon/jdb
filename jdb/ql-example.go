package jdb

func example() {

	From("table1.A").
		Join("table2.B").On("A.id").Eq("B.id").
		And("table3.id").Eq("1").
		Or("table3.id").Eq("2").
		LeftJoin("table4.B").On("A.id").Eq("C.id").
		RightJoin("table5.B").On("A.id").Eq("C.id").
		Select("A.id", "B.id").
		Where("table1.id").Eq(1).
		And("table2.id").Eq(2).
		All()

}
