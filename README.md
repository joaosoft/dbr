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

## With support for methods
* Select, Join, Distinct, Distinct on, Order by, Union, Load
* Insert, Multi insert, Record, Returning, Load
* Update, Record, Returning, Load
* Delete, Returning, Load
* With

## With support for type annotations
["-" when is to exlude]
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

## Usage 
This examples are available in the project at [dbr/examples](https://github.com/joaosoft/dbr/tree/master/examples)

```go
type Person struct {
	IdPerson  int    `json:"id_person" db.read:"id_person"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Age       int    `json:"age" db:"age"`
}

var db, _ = dbr.NewDbr()

func main() {
	Insert()
	Select()

	InsertValues()
	InsertRecords()
	SelectAll()
	SelectWith()

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
		OrderAsc("id_person").
		OrderDesc("first_name").
		From("public.person").
		Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Select("id_person", "first_name", "last_name", "age").
		From("public.person").
		OrderAsc("id_person").
		OrderDesc("first_name").
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

	fmt.Printf("\nDELETED ALL")
}
```

> ##### Result:
```
:: INSERT

QUERY: INSERT INTO public.person AS new_name (first_name, last_name, age) VALUES ('joao', 'ribeiro', 30)
SAVED PERSON: {IdPerson:0 FirstName:joao LastName:ribeiro Age:30}

:: SELECT

QUERY: SELECT id_person, first_name, last_name, age FROM public.person WHERE first_name = 'joao'
LOADED PERSON: {IdPerson:266 FirstName:joao LastName:ribeiro Age:30}

:: INSERT

QUERY: INSERT INTO public.person AS new_name (first_name, last_name, age) VALUES ('a', 'a', 1), ('b', 'b', 2), ('c', 'c', 3)
SAVED PERSONS!

:: INSERT

QUERY: INSERT INTO public.person AS new_name (first_name, last_name, age) VALUES ('joao', 'ribeiro', 30), ('luis', 'ribeiro', 31)
SAVED PERSON: {IdPerson:0 FirstName:joao LastName:ribeiro Age:30}

:: SELECT

QUERY: SELECT id_person, first_name, last_name, age FROM public.person ORDER BY id_person asc, first_name desc
LOADED PERSONS: [{IdPerson:266 FirstName:joao LastName:ribeiro Age:30} {IdPerson:267 FirstName:a LastName:a Age:1} {IdPerson:268 FirstName:b LastName:b Age:2} {IdPerson:269 FirstName:c LastName:c Age:3} {IdPerson:270 FirstName:joao LastName:ribeiro Age:30} {IdPerson:271 FirstName:luis LastName:ribeiro Age:31}]

:: SELECT WITH

QUERY: WITH load_one AS (SELECT first_name FROM public.person WHERE first_name = 'joao'), load_two AS (SELECT id_person, load_one.first_name, last_name, age FROM load_one, public.person AS person WHERE person.first_name = 'joao') SELECT id_person, first_name, last_name, age FROM load_two WHERE first_name = 'joao'
LOADED PERSON: {IdPerson:266 FirstName:joao LastName:ribeiro Age:30}

:: UPDATE

QUERY: UPDATE public.person SET last_name = 'males' WHERE first_name = 'joao'
UPDATED PERSON

:: SELECT

QUERY: SELECT id_person, first_name, last_name, age FROM public.person WHERE first_name = 'joao'
LOADED PERSON: {IdPerson:266 FirstName:joao LastName:males Age:30}

:: UPDATE

QUERY: UPDATE public.person SET last_name = 'males' WHERE first_name = 'joao'

AGE: 30
UPDATED PERSON

:: SELECT

QUERY: SELECT id_person, first_name, last_name, age FROM public.person WHERE first_name = 'joao'
LOADED PERSON: {IdPerson:266 FirstName:joao LastName:luis Age:30}

:: DELETE

QUERY: DELETE FROM public.person WHERE first_name = 'joao'
DELETED PERSON

:: TRANSACTION

QUERY: INSERT INTO public.person (first_name, last_name, age) VALUES ('joao-2', 'ribeiro', 30)
SAVED PERSON: {IdPerson:0 FirstName:joao-2 LastName:ribeiro Age:30}

:: DELETE

QUERY: DELETE FROM public.person WHERE first_name = 'joao-2'
DELETED PERSON

:: DELETE

QUERY: DELETE FROM public.person
DELETED ALL
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
