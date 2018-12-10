package dbr

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"github.com/joaosoft/errors"
	"io/ioutil"
	"os"
	"reflect"
)

func GetEnv() string {
	env := os.Getenv("env")
	if env == "" {
		env = "local"
	}

	return env
}

func Exists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ReadFile(file string, obj interface{}) ([]byte, error) {
	var err error

	if !Exists(file) {
		return nil, errors.New("0", "file don't exist")
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if obj != nil {
		if err := json.Unmarshal(data, obj); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func ReadFileLines(file string) ([]string, error) {
	lines := make([]string, 0)

	if !Exists(file) {
		return nil, errors.New("0", "file don't exist")
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func WriteFile(file string, obj interface{}) error {
	if !Exists(file) {
		return errors.New("0", "file don't exist")
	}

	jsonBytes, _ := json.MarshalIndent(obj, "", "    ")
	if err := ioutil.WriteFile(file, jsonBytes, 0644); err != nil {
		return err
	}

	return nil
}

func (stmt *StmtSelect) readRows(rows *sql.Rows, value reflect.Value) (int, error) {
	// read columns
	cols, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	// add columns to a map
	columns := make(map[string]bool)
	for _, col := range cols {
		columns[col] = true
	}

	value = value.Elem()
	isSlice := value.Kind() == reflect.Slice
	count := 0

	// load each row
	for rows.Next() {
		var elem reflect.Value
		if isSlice {
			elem = reflect.New(value.Type().Elem()).Elem()
		} else {
			elem = value
		}

		// load field values
		fields, err := getFields(columns, elem)
		if err != nil {
			return 0, err
		}

		// scan values from row
		err = rows.Scan(fields...)
		if err != nil {
			return 0, err
		}

		count++
		if isSlice {
			value.Set(reflect.Append(value, elem))
		} else {
			break
		}
	}

	return count, nil
}

func getFields(columns map[string]bool, object reflect.Value) ([]interface{}, error) {
	var fields []interface{}

	mappedValues := make(map[string]interface{})
	loadColumnStructValues(columns, object, mappedValues)

	for key, _ := range columns {
		fields = append(fields, mappedValues[key])
	}

	return fields, nil
}

func loadColumnStructValues(columns map[string]bool, object reflect.Value, mappedValues map[string]interface{}) {
	switch object.Kind() {
	case reflect.Ptr:
		if !object.IsNil() {
			loadColumnStructValues(columns, object.Elem(), mappedValues)
		}
	case reflect.Struct:
		t := object.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" && !field.Anonymous {
				// unexported
				continue
			}
			tag := field.Tag.Get("db")
			if tag == "-" || tag == "" {
				// ignore
				continue
			}

			if _, ok := columns[tag]; ok {
				mappedValues[tag] = object.Field(i).Addr().Interface()
			}
		}
	}
}

func loadStructValues(object reflect.Value, mappedValues map[string]reflect.Value) {
	switch object.Kind() {
	case reflect.Ptr:
		if !object.IsNil() {
			loadStructValues(object.Elem(), mappedValues)
		}
	case reflect.Struct:
		t := object.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" && !field.Anonymous {
				// unexported
				continue
			}
			tag := field.Tag.Get("db")
			if tag == "-" || tag == "" {
				// ignore
				continue
			}

			mappedValues[tag] = object.Field(i)
		}
	}
}