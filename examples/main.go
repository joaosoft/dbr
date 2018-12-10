package main

import (
	"dbr"
	"fmt"
)

type Person struct {
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Email     string `json:"email" db:"email"`
	Active    bool   `json:"active" db:"active"`
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
		FirstName: "joao-luis-ramos-ribeiro",
		LastName:  "ribeiro",
		Email:     "a@a.pt",
		Active:    true,
	}

	builder, _ := db.Insert().Into("buyer.contact").Record(person).Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Insert().Into("buyer.contact").Record(person).Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nSAVED PERSON: %+v", person)
}

func Select() {
	fmt.Println("\n\n:: SELECT")

	var person Person

	builder, _ := db.Select("first_name", "last_name", "email").From("buyer.contact").Where("first_name = ?", "joao-luis-ramos-ribeiro").Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Select("first_name", "last_name", "email").From("buyer.contact").Where("first_name = ?", "joao-luis-ramos-ribeiro").Load(&person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nLOADED PERSON: %+v", person)
}

func Update() {
	fmt.Println("\n\n:: UPDATE")

	builder, _ := db.Update("buyer.contact").Set("last_name", "arnaldo").Where("first_name = ?", "joao-luis-ramos-ribeiro").Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Update("buyer.contact").Set("last_name", "arnaldo").Where("first_name = ?", "joao-luis-ramos-ribeiro").Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nUPDATED PERSON")
}

func Delete() {
	fmt.Println("\n\n:: DELETE")

	builder, _ := db.Delete().From("buyer.contact").Where("first_name = ?", "joao-luis-ramos-ribeiro").Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Delete().From("buyer.contact").Where("first_name = ?", "joao-luis-ramos-ribeiro").Exec()
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
		FirstName: "joao-luis-ramos-ribeiro-2",
		LastName:  "ribeiro",
		Email:     "b@b.pt",
		Active:    true,
	}

	builder, _ := tx.Insert().Into("buyer.contact").Record(person).Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := tx.Insert().Into("buyer.contact").Record(person).Exec()
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

	builder, _ := db.Delete().From("buyer.contact").Where("first_name = ?", "joao-luis-ramos-ribeiro-2").Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Delete().From("buyer.contact").Where("first_name = ?", "joao-luis-ramos-ribeiro-2").Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nDELETED PERSON")
}
