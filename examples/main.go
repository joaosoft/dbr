package main

import (
	"database/sql"
	"dbr"
	"fmt"
	"strings"
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

var db, _ = dbr.New(
	dbr.WithSuccessEventHandler(
		func(operation dbr.SqlOperation, table []string, query string, rows *sql.Rows, sqlResult sql.Result) error {
			fmt.Printf("\nSuccess event [operation: %s, tables: %s, query: %s]", operation, strings.Join(table, "; "), query)
			return nil
		}),
	dbr.WithErrorEventHandler(func(operation dbr.SqlOperation, table []string, query string, err error) error {
		fmt.Printf("\nError event [operation: %s, tables: %s, query: %s, error: %s]", operation, strings.Join(table, "; "), query, err.Error())
		return nil
	}))

func main() {
	DeleteAll()

	Insert()
	InsertOnConflict()
	Select()
	SelectOr()

	SelectMax()
	SelectMin()
	SelectSum()
	SelectAvg()
	SelectCount()
	SelectCountDistinct()
	SelectFunction()

	InsertValues()
	InsertRecords()
	SelectAll()
	SelectWith()
	SelectWithRecursive()
	InsertWith()
	SelectGroupBy()
	Join()

	Update()
	Select()
	SelectWithMultipleFrom()
	SelectCoalesce()
	SelectCase()

	UpdateReturning()
	Select()
	Delete()

	Execute()

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

	stmt := db.Insert().
		Into(dbr.As("person", "new_name")).
		Record(person)

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nSAVED PERSON: %+v", person)
}

func InsertOnConflict() {
	fmt.Println("\n\n:: INSERT")

	stmt := db.Insert().
		Into(dbr.As("person", "new_name")).
		Columns("first_name", "last_name", "age").
		Values("duplicated", "duplicated", 10)

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Exec()

	// on conflict do update
	stmt = db.Insert().
		Into(dbr.As("person", "new_name")).
		Columns("first_name", "last_name", "age").
		Values("duplicated", "duplicated", 10).
		OnConflict("id_person").
		DoUpdate("id_person", 100)

	query, err = stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	// on conflict do nothing
	stmt = db.Insert().
		Into(dbr.As("person", "new_name")).
		Columns("first_name", "last_name", "age").
		Values("duplicated", "duplicated", 10).
		OnConflict("id_person").
		DoNothing()

	query, err = stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}
}

func InsertValues() {
	fmt.Println("\n\n:: INSERT")

	stmt := db.Insert().
		Into(dbr.As("person", "new_name")).
		Columns("first_name", "last_name", "age").
		Values("a", "a", 1).
		Values("b", "b", 2).
		Values("c", "c", 3)

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Exec()
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

	stmt := db.Insert().
		Into(dbr.As("person", "new_name")).
		Record(person1).
		Record(person2)

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nSAVED PERSON: %+v", person1)
}

func Select() {
	fmt.Println("\n\n:: SELECT")

	var person Person

	stmt := db.Select("id_person", "first_name", "last_name", "age").
		From("person").
		Where("first_name = ?", "joao")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Load(&person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nLOADED PERSON: %+v", person)
}

func SelectMax() {
	fmt.Println("\n\n:: SELECT MAX")

	var age int

	stmt := db.Select(dbr.Max("age")).
		From("person")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Load(&age)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nMAX PERSON AGE: %+v", age)
}

func SelectCount() {
	fmt.Println("\n\n:: SELECT COUNT")

	var age int

	stmt := db.Select(dbr.Count("age")).
		From("person")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Load(&age)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nCOUNT PERSON AGE: %+v", age)
}

func SelectCountDistinct() {
	fmt.Println("\n\n:: SELECT COUNT DISTINCT")

	var age int

	stmt := db.Select(dbr.Count("age", true)).
		From("person")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Load(&age)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nCOUNT DISTINCT PERSON AGE: %+v", age)
}

func SelectAvg() {
	fmt.Println("\n\n:: SELECT AVG")

	var age float64

	stmt := db.Select(dbr.Avg("age")).
		From("person")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Load(&age)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nAVG PERSON AGE: %+v", age)
}

func SelectSum() {
	fmt.Println("\n\n:: SELECT SUM")

	var age int

	stmt := db.Select(dbr.Sum("age")).
		From("person")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Load(&age)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nSUM PERSON AGE: %+v", age)
}

func SelectMin() {
	fmt.Println("\n\n:: SELECT MIN")

	var age int

	stmt := db.Select(dbr.Min("age")).
		From("person")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Load(&age)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nMIN PERSON AGE: %+v", age)
}

func SelectFunction() {
	fmt.Println("\n\n:: SELECT FUNCTION")

	var age int

	stmt := db.Select(dbr.Function("MAX", "age")).
		From("person")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Load(&age)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nMAX PERSON AGE: %+v", age)
}

func SelectWithMultipleFrom() {
	fmt.Println("\n\n:: SELECT WITH MULTIPLE FROM")

	var person Person

	stmt := db.Select("id_person", "first_name", "last_name", "age", "street").
		From("person").
		From("address").
		Where("first_name = ?", "joao")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Load(&person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nLOADED PERSON: %+v", person)
}

func SelectCoalesce() {
	fmt.Println("\n\n:: SELECT COALESCE")

	var person Person

	stmt := db.Select("id_person", "first_name", "last_name", dbr.OnNull("age", "0", "age")).
		From("person").
		Where("first_name = ?", "joao")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Load(&person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nLOADED PERSON: %+v", person)
}

func SelectCase() {
	fmt.Println("\n\n:: SELECT CASE")

	var person Person

	stmt := db.Select("id_person", "first_name", "last_name",
		dbr.Case("age").
			When("age = ?", 0).Then(10).
			When("age = ? OR first_name = ?", 30, "joao").Then(100).
			Else(20)).
		From("person").
		Where("first_name = ?", "joao")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Load(&person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nLOADED PERSON: %+v", person)
}

func SelectOr() {
	fmt.Println("\n\n:: SELECT OR")

	var person Person

	stmt := db.Select("id_person", "first_name", "last_name", "age").
		From("person").
		Where("first_name = ?", "joao").
		WhereOr("last_name = ?", "maria")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Load(&person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nLOADED PERSON: %+v", person)
}

func SelectAll() {
	fmt.Println("\n\n:: SELECT")

	var persons []Person

	stmt := db.Select("id_person", "first_name", "last_name", "age").
		From("person").
		OrderAsc("id_person").
		OrderDesc("first_name").
		Limit(5).
		Offset(1)

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Load(&persons)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nLOADED PERSONS: %+v", persons)
}

func SelectWith() {
	fmt.Println("\n\n:: SELECT WITH")

	var person Person

	stmt := db.
		With("load_one",
			db.Select("first_name").
				From("person").
				Where("first_name = ?", "joao")).
		With("load_two",
			db.Select("id_person", "load_one.first_name", "last_name", "age").
				From("load_one").
				From(dbr.As("person", "person")).
				Where("person.first_name = ?", "joao")).
		Select("id_person", "first_name", "last_name", "age").
		From("load_two").
		Where("first_name = ?", "joao")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Load(&person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nLOADED PERSON: %+v", person)
}

func SelectWithRecursive() {
	fmt.Println("\n\n:: SELECT WITH RECURSIVE")

	var person Person

	stmt := db.
		WithRecursive("load_one",
			db.Select("first_name").
				From("person").
				Where("first_name = ?", "joao")).
		With("load_two",
			db.Select("id_person", "load_one.first_name", "last_name", "age").
				From("load_one").
				From(dbr.As("person", "person")).
				Where("person.first_name = ?", "joao")).
		Select("id_person", "first_name", "last_name", "age").
		From("load_two").
		Where("first_name = ?", "joao")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Load(&person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nLOADED PERSON: %+v", person)
}

func InsertWith() {
	fmt.Println("\n\n:: INSERT WITH")

	var person Person

	stmt := db.
		With("load_one",
			db.Select("first_name").
				From("person").
				Where("first_name = ?", "joao").
				Limit(1)).
		With("load_two",
			db.Select("id_person", "load_one.first_name", "last_name", "age").
				From("load_one").
				From(dbr.As("person", "person")).
				Where("person.first_name = ?", "joao").Limit(1)).
		Insert().
		Into("person").
		Columns("id_person", "first_name", "last_name", "age").
		FromSelect(
			db.Select(999, "first_name", "last_name", "age").
				From("load_two"))

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nINSERT PERSON 999: %+v", person)

	fmt.Println("\n\n:: SELECT")

	stmtSelect := db.Select("id_person", "first_name", "last_name", "age").
		From("person").
		Where("id_person = ?", 999)

	query, err = stmtSelect.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmtSelect.Load(&person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nLOADED PERSON 999: %+v", person)
}

func SelectGroupBy() {
	fmt.Println("\n\n:: SELECT GROUP BY")

	var persons []Person

	stmt := db.Select("id_person", "first_name", "last_name", "age").
		From("person").
		OrderAsc("age").
		OrderDesc("first_name").
		GroupBy("id_person", "last_name", "first_name", "age").
		Having("age > 20").
		Limit(5).
		Offset(1)

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Load(&persons)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nLOADED PERSONS: %+v", persons)
}

func Update() {
	fmt.Println("\n\n:: UPDATE")

	stmt := db.Update("person").
		Set("last_name", "males").
		Where("first_name = ?", "joao")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nUPDATED PERSON")
}

func UpdateReturning() {
	fmt.Println("\n\n:: UPDATE")

	stmt := db.Update("person").
		Set("last_name", "males").
		Where("first_name = ?", "joao").
		Return("age")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	var age int
	err = stmt.Load(&age)
	fmt.Printf("\n\nAGE: %d", age)

	if err != nil {
		panic(err)
	}

	fmt.Printf("\nUPDATED PERSON")
}

func Delete() {
	fmt.Println("\n\n:: DELETE")

	stmt := db.Delete().
		From("person").
		Where("first_name = ?", "joao")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Exec()
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

	stmtInsert := db.Insert().
		Into(dbr.As("address", "new_name")).
		Record(address)

	query, err := stmtInsert.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmtInsert.Exec()
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

	stmtInsert = db.Insert().
		Into(dbr.As("person", "new_name")).
		Record(person)

	query, err = stmtInsert.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmtInsert.Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nSAVED PERSON: %+v", person)

	stmtSelect := db.Select("address.street").
		From("person").
		Join("address", "fk_address = id_address").
		Where("first_name = ?", "joao-join")

	query, err = stmtSelect.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	var street string
	_, err = stmtSelect.Load(&street)
	fmt.Printf("\nSTREET: %s", street)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nSAVED ADDRESS: %+v", person)
}

func Execute() {
	fmt.Println("\n\n:: EXECUTE")

	stmt := db.Execute("SELECT * FROM person WHERE first_name LIKE ?").
		Values("%joao%")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n EXECUTE DONE")
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

	stmt := tx.Insert().
		Into("person").
		Record(person)

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Exec()
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

	stmt := db.Delete().
		From("person").
		Where("first_name = ?", "joao-2")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nDELETED PERSON")
}

func DeleteAll() {
	fmt.Println("\n\n:: DELETE")

	stmt := db.Delete().
		From("person")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	stmt = db.Delete().
		From("address")

	query, err = stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nDELETED ALL")
}
