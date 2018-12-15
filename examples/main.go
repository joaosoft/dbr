package main

import (
	"dbr"
	"fmt"
)

type Person struct {
	IdPerson  int    `json:"id_person" db.read:"id_person"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Age       int    `json:"age" db:"age"`
	IdAddress *int   `json:"fk_address" db:"fk_address"`
}

type Address struct {
	IdAddress int    `json:"id_address" db:"id_address"`
	Street    string `json:"street" db:"street"`
	Number    int    `json:"number" db:"number"`
	Country   string `json:"country" db:"country"`
}

var db, _ = dbr.NewDbr()

func main() {
	Insert()
	Select()

	InsertValues()
	InsertRecords()
	SelectAll()
	SelectWith()
	SelectGroupBy()
	Join()

	Update()
	Select()

	UpdateReturning()
	Select()
	Delete()

	Transaction()
	DeleteTransactionData()

	DeleteAll()
}

func Insert() {
	fmt.Println("\n\n:: INSERT")

	person := Person{
		FirstName: "joao",
		LastName:  "ribeiro",
		Age:       30,
	}

	builder, _ := db.Insert().
		Into(dbr.Field("public.person").As("new_name")).
		Record(person).
		Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Insert().
		Into(dbr.Field("public.person").As("new_name")).
		Record(person).
		Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nSAVED PERSON: %+v", person)
}

func InsertValues() {
	fmt.Println("\n\n:: INSERT")

	builder, _ := db.Insert().
		Into(dbr.Field("public.person").As("new_name")).
		Columns("first_name", "last_name", "age").
		Values("a", "a", 1).
		Values("b", "b", 2).
		Values("c", "c", 3).
		Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Insert().
		Into(dbr.Field("public.person").As("new_name")).
		Columns("first_name", "last_name", "age").
		Values("a", "a", 1).
		Values("b", "b", 2).
		Values("c", "c", 3).
		Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nSAVED PERSONS!")
}

func InsertRecords() {
	fmt.Println("\n\n:: INSERT")

	person1 := Person{
		FirstName: "joao",
		LastName:  "ribeiro",
		Age:       30,
	}

	person2 := Person{
		FirstName: "luis",
		LastName:  "ribeiro",
		Age:       31,
	}

	builder, _ := db.Insert().
		Into(dbr.Field("public.person").As("new_name")).
		Record(person1).
		Record(person2).
		Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Insert().
		Into(dbr.Field("public.person").As("new_name")).
		Record(person1).
		Record(person2).
		Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nSAVED PERSON: %+v", person1)
}

func Select() {
	fmt.Println("\n\n:: SELECT")

	var person Person

	builder, _ := db.Select("id_person", "first_name", "last_name", "age").
		From("public.person").
		Where("first_name = ?", "joao").
		Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Select("id_person", "first_name", "last_name", "age").
		From("public.person").
		Where("first_name = ?", "joao").
		Load(&person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nLOADED PERSON: %+v", person)
}

func SelectAll() {
	fmt.Println("\n\n:: SELECT")

	var persons []Person

	builder, _ := db.Select("id_person", "first_name", "last_name", "age").
		From("public.person").
		OrderAsc("id_person").
		OrderDesc("first_name").
		Limit(5).
		Offset(1).
		Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Select("id_person", "first_name", "last_name", "age").
		From("public.person").
		OrderAsc("id_person").
		OrderDesc("first_name").
		Limit(5).
		Offset(1).
		Load(&persons)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nLOADED PERSONS: %+v", persons)
}

func SelectWith() {
	fmt.Println("\n\n:: SELECT WITH")

	var person Person

	builder := db.
		With("load_one",
			db.Select("first_name").
				From("public.person").
				Where("first_name = ?", "joao")).
		With("load_two",
			db.Select("id_person", "load_one.first_name", "last_name", "age").
				From("load_one").
				From(dbr.Field("public.person").As("person")).
				Where("person.first_name = ?", "joao")).
		Select("id_person", "first_name", "last_name", "age").
		From("load_two").
		Where("first_name = ?", "joao")

	stmt, _ := builder.Build()
	fmt.Printf("\nQUERY: %s", stmt)

	_, err := builder.Load(&person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nLOADED PERSON: %+v", person)
}

func SelectGroupBy() {
	fmt.Println("\n\n:: SELECT GROUP BY")

	var persons []Person

	builder, _ := db.Select("id_person", "first_name", "last_name", "age").
		From("public.person").
		OrderAsc("age").
		OrderDesc("first_name").
		GroupBy("id_person", "last_name", "first_name", "age").
		Having("age > 20").
		Limit(5).
		Offset(1).
		Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Select("id_person", "first_name", "last_name", "age").
		From("public.person").
		OrderAsc("age").
		OrderDesc("first_name").
		GroupBy("id_person", "last_name", "first_name", "age").
		Having("age > 20").
		Limit(5).
		Offset(1).
		Load(&persons)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nLOADED PERSONS: %+v", persons)
}

func Update() {
	fmt.Println("\n\n:: UPDATE")

	builder, _ := db.Update("public.person").
		Set("last_name", "males").
		Where("first_name = ?", "joao").
		Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Update("public.person").
		Set("last_name", "males").
		Where("first_name = ?", "joao").
		Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nUPDATED PERSON")
}

func UpdateReturning() {
	fmt.Println("\n\n:: UPDATE")

	builder, _ := db.Update("public.person").
		Set("last_name", "males").
		Where("first_name = ?", "joao").
		Build()
	fmt.Printf("\nQUERY: %s", builder)

	var age int
	err := db.Update("public.person").
		Set("last_name", "luis").
		Where("first_name = ?", "joao").
		Return("age").
		Load(&age)
	fmt.Printf("\n\nAGE: %d", age)

	if err != nil {
		panic(err)
	}

	fmt.Printf("\nUPDATED PERSON")
}

func Delete() {
	fmt.Println("\n\n:: DELETE")

	builder, _ := db.Delete().
		From("public.person").
		Where("first_name = ?", "joao").
		Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Delete().
		From("public.person").
		Where("first_name = ?", "joao").
		Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nDELETED PERSON")
}

func Join() {
	fmt.Println("\n\n:: JOIN")

	address := Address{
		IdAddress: 1,
		Street:    "street one",
		Number:    1,
		Country:   "portugal",
	}

	builder, _ := db.Insert().
		Into(dbr.Field("public.address").As("new_name")).
		Record(address).
		Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Insert().
		Into(dbr.Field("public.address").As("new_name")).
		Record(address).
		Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nSAVED ADDRESS: %+v", address)

	idAddress := 1
	person := Person{
		FirstName: "joao-join",
		LastName:  "ribeiro-join",
		Age:       30,
		IdAddress: &idAddress,
	}

	builder, _ = db.Insert().
		Into(dbr.Field("public.person").As("new_name")).
		Record(person).
		Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err = db.Insert().
		Into(dbr.Field("public.person").As("new_name")).
		Record(person).
		Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nSAVED PERSON: %+v", person)

	builder, _ = db.Select("address.street").
		From("public.person").
		Join("public.address", "fk_address = id_address").
		Where("first_name = ?", "joao-join").
		Build()
	fmt.Printf("\nQUERY: %s", builder)

	var street string
	_, err = db.Select("address.street").
		From("public.person").
		Join("public.address", "fk_address = id_address").
		Where("first_name = ?", "joao-join").
		Load(&street)
	fmt.Printf("\nSTREET: %s", street)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nSAVED ADDRESS: %+v", person)
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

	builder, _ := tx.Insert().
		Into("public.person").
		Record(person).
		Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := tx.Insert().
		Into("public.person").
		Record(person).
		Exec()
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

	builder, _ := db.Delete().
		From("public.person").
		Where("first_name = ?", "joao-2").
		Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Delete().
		From("public.person").
		Where("first_name = ?", "joao-2").
		Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nDELETED PERSON")
}

func DeleteAll() {
	fmt.Println("\n\n:: DELETE")

	builder, _ := db.Delete().
		From("public.person").
		Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Delete().
		From("public.person").
		Exec()
	if err != nil {
		panic(err)
	}

	builder, _ = db.Delete().
		From("public.address").
		Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err = db.Delete().
		From("public.address").
		Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nDELETED ALL")
}
