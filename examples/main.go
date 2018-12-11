package main

import (
	"dbr"
	"fmt"
)

type Person struct {
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Age       int    `json:"age" db:"age"`
}

var db, _ = dbr.NewDbr()

func main() {
	defer Delete()
	defer DeleteTransactionData()

	Insert()
	Select()
	Update()
	Select()
	Transaction()
}

func Insert() {
	fmt.Println("\n\n:: INSERT")

	person := Person{
		FirstName: "joao",
		LastName:  "ribeiro",
		Age:       30,
	}

	builder, _ := db.Insert().Into(dbr.Field("public.person").As("new_name")).Record(person).Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Insert().Into(dbr.Field("public.person").As("new_name")).Record(person).Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nSAVED PERSON: %+v", person)
}

func Select() {
	fmt.Println("\n\n:: SELECT")

	var person Person

	builder, _ := db.Select("first_name", "last_name", "age").From("public.person").Where("first_name = ?", "joao").Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Select("first_name", "last_name", "age").From("public.person").Where("first_name = ?", "joao").Load(&person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nLOADED PERSON: %+v", person)
}

func Update() {
	fmt.Println("\n\n:: UPDATE")

	builder, _ := db.Update("public.person").Set("last_name", "males").Where("first_name = ?", "joao").Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Update("public.person").Set("last_name", "males").Where("first_name = ?", "joao").Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nUPDATED PERSON")
}

func Delete() {
	fmt.Println("\n\n:: DELETE")

	builder, _ := db.Delete().From("public.person").Where("first_name = ?", "joao").Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Delete().From("public.person").Where("first_name = ?", "joao").Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nDELETED PERSON")
}

func Transaction() {
	fmt.Println("\n\n:: TRANSACTION")

	tx, _ := db.Begin()
	defer tx.RollbackUnlessCommit()

	person := Person{
		FirstName: "joao-2",
		LastName:  "ribeiro",
		Age:       30,
	}

	builder, _ := tx.Insert().Into("public.person").Record(person).Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := tx.Insert().Into("public.person").Record(person).Exec()
	if err != nil {
		panic(err)
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}
	fmt.Printf("\nSAVED PERSON: %+v", person)
}

func DeleteTransactionData() {
	fmt.Println("\n\n:: DELETE")

	builder, _ := db.Delete().From("public.person").Where("first_name = ?", "joao-2").Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Delete().From("public.person").Where("first_name = ?", "joao-2").Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nDELETED PERSON")
}
