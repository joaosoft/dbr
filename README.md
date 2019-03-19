Dbr
================

[![Build Status](https://travis-ci.org/joaosoft/dbr.svg?branch=master)](https://travis-ci.org/joaosoft/dbr) | [![codecov](https://codecov.io/gh/joaosoft/dbr/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/dbr) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/dbr)](https://goreportcard.com/report/github.com/joaosoft/dbr) | [![GoDoc](https://godoc.org/github.com/joaosoft/dbr?status.svg)](https://godoc.org/github.com/joaosoft/dbr)

A simple database client with support for master/slave databases.
The main goal of this project is to allow a application to write in a master database and read the data from a slave (replica).

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for databases
* Postgres 
* MySql
* SqlLite3

## With configuration options
* WithConfiguration
* WithLogger
* WithLogLevel
* WithManager
* WithDatabase
* WithSuccessEventHandler (call's the function when a query on db ends with success)
* WithErrorEventHandler (call's the function when a query on db ends with error)

## With support for methods
* Select, Where, Join, Distinct, Distinct on, Group by, Having, Order by, Union, Intersect, Except, Limit, Offset, Load
* Insert, Multi insert, Where, Record, Returning, Load
* Update, Where, Set, Record, Returning, Load
* Delete, Where, Returning, Load
* With, With Recursive
* OnConflict (DoNothing, DoUpdate)
* Execute
* UseOnlyRead, UseOnlyWrite (allows to use only read or write connection for the query)

## With support for type annotations
["-" when is to exclude a field]
* db -> used to read and write
* db.read -> used for select
* db.write -> used for insert and update

## Dependecy Management
>### Dependency

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Get dependency manager: `go get github.com/joaosoft/dependency`
* Install dependencies: `dependency get`

>### Go
```
go get github.com/joaosoft/dbr
```

>### Configuration
>>#### master / slave
```
{
  "dbr": {
    "read_db": {
      "driver": "postgres",
      "datasource": "postgres://user:password@localhost:7000/postgres?sslmode=disable&search_path=public"
    },
    "write_db": {
      "driver": "postgres",
      "datasource": "postgres://user:password@localhost:7100/postgres?sslmode=disable&search_path=public"
    },
    "log": {
      "level": "info"
    }
  }
}
```

>>#### one instance only
```
{
  "dbr": {
    "db": {
      "driver": "postgres",
      "datasource": "postgres://user:password@localhost:7000/postgres?sslmode=disable&search_path=public"
    },
    "log": {
      "level": "info"
    }
  }
}
```

## Usage 
This examples are available in the project at [dbr/examples](https://github.com/joaosoft/dbr/tree/master/examples)

```go
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
		Into(dbr.As("public.person", "new_name")).
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
		Into(dbr.As("public.person", "new_name")).
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
		Into(dbr.As("public.person", "new_name")).
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
		Into(dbr.As("public.person", "new_name")).
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
		Into(dbr.As("public.person", "new_name")).
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
		Into(dbr.As("public.person", "new_name")).
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
		From("public.person").
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
		From("public.person").
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
		From("public.person").
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
				From("public.person").
				Where("first_name = ?", "joao")).
		With("load_two",
			db.Select("id_person", "load_one.first_name", "last_name", "age").
				From("load_one").
				From(dbr.As("public.person", "person")).
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
				From("public.person").
				Where("first_name = ?", "joao")).
		With("load_two",
			db.Select("id_person", "load_one.first_name", "last_name", "age").
				From("load_one").
				From(dbr.As("public.person", "person")).
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
				From("public.person").
				Where("first_name = ?", "joao").
				Limit(1)).
		With("load_two",
			db.Select("id_person", "load_one.first_name", "last_name", "age").
				From("load_one").
				From(dbr.As("public.person", "person")).
				Where("person.first_name = ?", "joao").Limit(1)).
		Insert().
		Into("public.person").
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
		From("public.person").
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
		From("public.person").
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

	stmt := db.Update("public.person").
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

	stmt := db.Update("public.person").
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
		From("public.person").
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
		Into(dbr.As("public.address", "new_name")).
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
		Into(dbr.As("public.person", "new_name")).
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
		From("public.person").
		Join("public.address", "fk_address = id_address").
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

	stmt := db.Execute("SELECT * FROM public.person WHERE first_name LIKE ?").
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
		Into("public.person").
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
		From("public.person").
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
		From("public.person")

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
		From("public.address")

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
```

> ##### Result:
```
:: DELETE

QUERY: DELETE FROM public.person
Success event [operation: DELETE, tables: public.person, query: DELETE FROM public.person]
QUERY: DELETE FROM public.address
Success event [operation: DELETE, tables: public.address, query: DELETE FROM public.address]
DELETED ALL

:: INSERT

QUERY: INSERT INTO "public"."person" AS new_name ("first_name", "last_name", "age", "fk_address") VALUES ('joao', 'ribeiro', 30, NULL)
Success event [operation: INSERT, tables: "public"."person" AS new_name, query: INSERT INTO "public"."person" AS new_name ("first_name", "last_name", "age", "fk_address") VALUES ('joao', 'ribeiro', 30, NULL)]
SAVED PERSON: {IdPerson:0 FirstName:joao LastName:ribeiro Age:30 IdAddress:<nil>}

:: INSERT

QUERY: INSERT INTO "public"."person" AS new_name ("first_name", "last_name", "age") VALUES ('duplicated', 'duplicated', 10)
Success event [operation: INSERT, tables: "public"."person" AS new_name, query: INSERT INTO "public"."person" AS new_name ("first_name", "last_name", "age") VALUES ('duplicated', 'duplicated', 10)]
QUERY: INSERT INTO "public"."person" AS new_name ("first_name", "last_name", "age") VALUES ('duplicated', 'duplicated', 10) ON CONFLICT ("id_person") DO UPDATE SET id_person = 100
Success event [operation: INSERT, tables: "public"."person" AS new_name, query: INSERT INTO "public"."person" AS new_name ("first_name", "last_name", "age") VALUES ('duplicated', 'duplicated', 10) ON CONFLICT ("id_person") DO UPDATE SET id_person = 100]
QUERY: INSERT INTO "public"."person" AS new_name ("first_name", "last_name", "age") VALUES ('duplicated', 'duplicated', 10) ON CONFLICT ("id_person") DO NOTHING
Success event [operation: INSERT, tables: "public"."person" AS new_name, query: INSERT INTO "public"."person" AS new_name ("first_name", "last_name", "age") VALUES ('duplicated', 'duplicated', 10) ON CONFLICT ("id_person") DO NOTHING]

:: SELECT

QUERY: SELECT "id_person", "first_name", "last_name", "age" FROM public.person WHERE first_name = 'joao'
Success event [operation: SELECT, tables: public.person, query: SELECT "id_person", "first_name", "last_name", "age" FROM public.person WHERE first_name = 'joao']
LOADED PERSON: {IdPerson:238 FirstName:joao LastName:ribeiro Age:30 IdAddress:<nil>}

:: SELECT OR

QUERY: SELECT "id_person", "first_name", "last_name", "age" FROM public.person WHERE first_name = 'joao' OR last_name = 'maria'
Success event [operation: SELECT, tables: public.person, query: SELECT "id_person", "first_name", "last_name", "age" FROM public.person WHERE first_name = 'joao' OR last_name = 'maria']
LOADED PERSON: {IdPerson:238 FirstName:joao LastName:ribeiro Age:30 IdAddress:<nil>}

:: INSERT

QUERY: INSERT INTO "public"."person" AS new_name ("first_name", "last_name", "age") VALUES ('a', 'a', 1), ('b', 'b', 2), ('c', 'c', 3)
Success event [operation: INSERT, tables: "public"."person" AS new_name, query: INSERT INTO "public"."person" AS new_name ("first_name", "last_name", "age") VALUES ('a', 'a', 1), ('b', 'b', 2), ('c', 'c', 3)]
SAVED PERSONS!

:: INSERT

QUERY: INSERT INTO "public"."person" AS new_name ("first_name", "last_name", "age", "fk_address") VALUES ('joao', 'ribeiro', 30, NULL), ('luis', 'ribeiro', 31, NULL)
Success event [operation: INSERT, tables: "public"."person" AS new_name, query: INSERT INTO "public"."person" AS new_name ("first_name", "last_name", "age", "fk_address") VALUES ('joao', 'ribeiro', 30, NULL), ('luis', 'ribeiro', 31, NULL)]
SAVED PERSON: {IdPerson:0 FirstName:joao LastName:ribeiro Age:30 IdAddress:<nil>}

:: SELECT

QUERY: SELECT "id_person", "first_name", "last_name", "age" FROM public.person ORDER BY id_person asc, first_name desc LIMIT 5 OFFSET 1
Success event [operation: SELECT, tables: public.person, query: SELECT "id_person", "first_name", "last_name", "age" FROM public.person ORDER BY id_person asc, first_name desc LIMIT 5 OFFSET 1]
LOADED PERSONS: [{IdPerson:239 FirstName:duplicated LastName:duplicated Age:10 IdAddress:<nil>} {IdPerson:240 FirstName:duplicated LastName:duplicated Age:10 IdAddress:<nil>} {IdPerson:241 FirstName:duplicated LastName:duplicated Age:10 IdAddress:<nil>} {IdPerson:242 FirstName:a LastName:a Age:1 IdAddress:<nil>} {IdPerson:243 FirstName:b LastName:b Age:2 IdAddress:<nil>}]

:: SELECT WITH

QUERY: WITH load_one AS (SELECT "first_name" FROM public.person WHERE first_name = 'joao'), load_two AS (SELECT "id_person", "load_one"."first_name", "last_name", "age" FROM load_one, "public"."person" AS person WHERE person.first_name = 'joao')SELECT "id_person", "first_name", "last_name", "age" FROM load_two WHERE first_name = 'joao'
Success event [operation: SELECT, tables: load_two, query: WITH load_one AS (SELECT "first_name" FROM public.person WHERE first_name = 'joao'), load_two AS (SELECT "id_person", "load_one"."first_name", "last_name", "age" FROM load_one, "public"."person" AS person WHERE person.first_name = 'joao')SELECT "id_person", "first_name", "last_name", "age" FROM load_two WHERE first_name = 'joao']
LOADED PERSON: {IdPerson:238 FirstName:joao LastName:ribeiro Age:30 IdAddress:<nil>}

:: SELECT WITH RECURSIVE

QUERY: WITH RECURSIVE load_one AS (SELECT "first_name" FROM public.person WHERE first_name = 'joao'), load_two AS (SELECT "id_person", "load_one"."first_name", "last_name", "age" FROM load_one, "public"."person" AS person WHERE person.first_name = 'joao')SELECT "id_person", "first_name", "last_name", "age" FROM load_two WHERE first_name = 'joao'
Success event [operation: SELECT, tables: load_two, query: WITH RECURSIVE load_one AS (SELECT "first_name" FROM public.person WHERE first_name = 'joao'), load_two AS (SELECT "id_person", "load_one"."first_name", "last_name", "age" FROM load_one, "public"."person" AS person WHERE person.first_name = 'joao')SELECT "id_person", "first_name", "last_name", "age" FROM load_two WHERE first_name = 'joao']
LOADED PERSON: {IdPerson:238 FirstName:joao LastName:ribeiro Age:30 IdAddress:<nil>}

:: INSERT WITH

QUERY: WITH load_one AS (SELECT "first_name" FROM public.person WHERE first_name = 'joao' LIMIT 1), load_two AS (SELECT "id_person", "load_one"."first_name", "last_name", "age" FROM load_one, "public"."person" AS person WHERE person.first_name = 'joao' LIMIT 1)INSERT INTO public.person ("id_person", "first_name", "last_name", "age") SELECT 999, "first_name", "last_name", "age" FROM load_two
Success event [operation: INSERT, tables: public.person, query: WITH load_one AS (SELECT "first_name" FROM public.person WHERE first_name = 'joao' LIMIT 1), load_two AS (SELECT "id_person", "load_one"."first_name", "last_name", "age" FROM load_one, "public"."person" AS person WHERE person.first_name = 'joao' LIMIT 1)INSERT INTO public.person ("id_person", "first_name", "last_name", "age") SELECT 999, "first_name", "last_name", "age" FROM load_two]
INSERT PERSON 999: {IdPerson:0 FirstName: LastName: Age:0 IdAddress:<nil>}

:: SELECT

QUERY: SELECT "id_person", "first_name", "last_name", "age" FROM public.person WHERE id_person = '999'
Success event [operation: SELECT, tables: public.person, query: SELECT "id_person", "first_name", "last_name", "age" FROM public.person WHERE id_person = '999']
LOADED PERSON 999: {IdPerson:999 FirstName:joao LastName:ribeiro Age:30 IdAddress:<nil>}

:: SELECT GROUP BY

QUERY: SELECT "id_person", "first_name", "last_name", "age" FROM public.person GROUP BY id_person, last_name, first_name, age HAVING age > 20 ORDER BY age asc, first_name desc LIMIT 5 OFFSET 1
Success event [operation: SELECT, tables: public.person, query: SELECT "id_person", "first_name", "last_name", "age" FROM public.person GROUP BY id_person, last_name, first_name, age HAVING age > 20 ORDER BY age asc, first_name desc LIMIT 5 OFFSET 1]
LOADED PERSONS: [{IdPerson:238 FirstName:joao LastName:ribeiro Age:30 IdAddress:<nil>} {IdPerson:245 FirstName:joao LastName:ribeiro Age:30 IdAddress:<nil>} {IdPerson:246 FirstName:luis LastName:ribeiro Age:31 IdAddress:<nil>}]

:: JOIN

QUERY: INSERT INTO "public"."address" AS new_name ("id_address", "street", "number", "country") VALUES (1, 'street one', 1, 'portugal')
Success event [operation: INSERT, tables: "public"."address" AS new_name, query: INSERT INTO "public"."address" AS new_name ("id_address", "street", "number", "country") VALUES (1, 'street one', 1, 'portugal')]
SAVED ADDRESS: {IdAddress:1 Street:street one Number:1 Country:portugal}
QUERY: INSERT INTO "public"."person" AS new_name ("first_name", "last_name", "age", "fk_address") VALUES ('joao-join', 'ribeiro-join', 30, 1)
Success event [operation: INSERT, tables: "public"."person" AS new_name, query: INSERT INTO "public"."person" AS new_name ("first_name", "last_name", "age", "fk_address") VALUES ('joao-join', 'ribeiro-join', 30, 1)]
SAVED PERSON: {IdPerson:0 FirstName:joao-join LastName:ribeiro-join Age:30 IdAddress:0xc00015ff08}
QUERY: SELECT "address"."street" FROM public.person JOIN public.address ON (fk_address = id_address) WHERE first_name = 'joao-join'
Success event [operation: SELECT, tables: public.person, query: SELECT "address"."street" FROM public.person JOIN public.address ON (fk_address = id_address) WHERE first_name = 'joao-join']
STREET: street one
SAVED ADDRESS: {IdPerson:0 FirstName:joao-join LastName:ribeiro-join Age:30 IdAddress:0xc00015ff08}

:: UPDATE

QUERY: UPDATE public.person SET last_name = 'males' WHERE first_name = 'joao'
Success event [operation: UPDATE, tables: public.person, query: UPDATE public.person SET last_name = 'males' WHERE first_name = 'joao']
UPDATED PERSON

:: SELECT

QUERY: SELECT "id_person", "first_name", "last_name", "age" FROM public.person WHERE first_name = 'joao'
Success event [operation: SELECT, tables: public.person, query: SELECT "id_person", "first_name", "last_name", "age" FROM public.person WHERE first_name = 'joao']
LOADED PERSON: {IdPerson:238 FirstName:joao LastName:males Age:30 IdAddress:<nil>}

:: UPDATE

QUERY: UPDATE public.person SET last_name = 'males' WHERE first_name = 'joao' RETURNING "age"
Success event [operation: UPDATE, tables: public.person, query: UPDATE public.person SET last_name = 'males' WHERE first_name = 'joao' RETURNING "age"]

AGE: 30
UPDATED PERSON

:: SELECT

QUERY: SELECT "id_person", "first_name", "last_name", "age" FROM public.person WHERE first_name = 'joao'
Success event [operation: SELECT, tables: public.person, query: SELECT "id_person", "first_name", "last_name", "age" FROM public.person WHERE first_name = 'joao']
LOADED PERSON: {IdPerson:238 FirstName:joao LastName:males Age:30 IdAddress:<nil>}

:: DELETE

QUERY: DELETE FROM public.person WHERE first_name = 'joao'
Success event [operation: DELETE, tables: public.person, query: DELETE FROM public.person WHERE first_name = 'joao']
DELETED PERSON

:: EXECUTE

QUERY: SELECT * FROM public.person WHERE first_name LIKE '%joao%'
Success event [operation: EXECUTE, tables: , query: SELECT * FROM public.person WHERE first_name LIKE '%joao%']
 EXECUTE DONE

:: TRANSACTION

QUERY: INSERT INTO public.person ("first_name", "last_name", "age", "fk_address") VALUES ('joao-2', 'ribeiro', 30, NULL)
Success event [operation: INSERT, tables: public.person, query: INSERT INTO public.person ("first_name", "last_name", "age", "fk_address") VALUES ('joao-2', 'ribeiro', 30, NULL)]
SAVED PERSON: {IdPerson:0 FirstName:joao-2 LastName:ribeiro Age:30 IdAddress:<nil>}

:: DELETE

QUERY: DELETE FROM public.person WHERE first_name = 'joao-2'
Success event [operation: DELETE, tables: public.person, query: DELETE FROM public.person WHERE first_name = 'joao-2']
DELETED PERSON

:: DELETE

QUERY: DELETE FROM public.person
Success event [operation: DELETE, tables: public.person, query: DELETE FROM public.person]
QUERY: DELETE FROM public.address
Success event [operation: DELETE, tables: public.address, query: DELETE FROM public.address]
DELETED ALL
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
