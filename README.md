Dbr
================

[![Build Status](https://travis-ci.org/joaosoft/dbr.svg?branch=master)](https://travis-ci.org/joaosoft/dbr) | [![codecov](https://codecov.io/gh/joaosoft/dbr/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/dbr) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/dbr)](https://goreportcard.com/report/github.com/joaosoft/dbr) | [![GoDoc](https://godoc.org/github.com/joaosoft/dbr?status.svg)](https://godoc.org/github.com/joaosoft/dbr)

A simple database client with support for master/slave databases.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* Postgres 
* MySql
* SqlLite3

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

	builder, _ := db.Insert().Into(dbr.Field("buyer.contact").As("new_name")).Record(person).Build()
	fmt.Printf("\nQUERY: %s", builder)

	_, err := db.Insert().Into(dbr.Field("buyer.contact").As("new_name")).Record(person).Exec()
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
```

> ##### Result:
```
:: INSERT

QUERY: INSERT INTO buyer.contact AS new_name (email, active, first_name, last_name) VALUES ('a@a.pt', TRUE, 'joao-luis-ramos-ribeiro', 'ribeiro')
SAVED PERSON: {FirstName:joao-luis-ramos-ribeiro LastName:ribeiro Email:a@a.pt Active:true}

:: SELECT

QUERY: SELECT first_name, last_name, email FROM buyer.contact WHERE first_name = 'joao-luis-ramos-ribeiro'
LOADED PERSON: {FirstName:a@a.pt LastName:joao-luis-ramos-ribeiro Email:ribeiro Active:false}

:: UPDATE

QUERY: UPDATE buyer.contact SET last_name = 'arnaldo' WHERE first_name = 'joao-luis-ramos-ribeiro'
UPDATED PERSON

:: SELECT

QUERY: SELECT first_name, last_name, email FROM buyer.contact WHERE first_name = 'joao-luis-ramos-ribeiro'
LOADED PERSON: {FirstName:arnaldo LastName:a@a.pt Email:joao-luis-ramos-ribeiro Active:false}

:: TRANSACTION

QUERY: INSERT INTO buyer.contact (first_name, last_name, email, active) VALUES ('joao-luis-ramos-ribeiro-2', 'ribeiro', 'b@b.pt', TRUE)
SAVED PERSON: {FirstName:joao-luis-ramos-ribeiro-2 LastName:ribeiro Email:b@b.pt Active:true}

:: DELETE

QUERY: DELETE FROM buyer.contact WHERE first_name = 'joao-luis-ramos-ribeiro-2'
DELETED PERSON

:: DELETE

QUERY: DELETE FROM buyer.contact WHERE first_name = 'joao-luis-ramos-ribeiro'
DELETED PERSON
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
