package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"strings"

	td "github.com/soerlemans/table/table_data"
	u "github.com/soerlemans/table/util"
	// "github.com/xuri/excelize/v2"
)

const NEWLINE = "\n"

// Convert a byte blob into a text array.
func bytes2Text(t_input []byte) []string {
	text := strings.Split(string(t_input), NEWLINE)

	// We should only strip the last line if we have at least one line.
	if len(text) > 1 {
		// Strip any empty lines from the end.
		for index, line := range slices.Backward(text) {
			if len(line) != 0 {
				text = text[:index+1]
				break
			}
		}
	}

	return text
}

// Count how many values are truthy.
func CountTruthy[T any](t_args ...T) int {
	var count int

	for _, value := range t_args {
		v := reflect.ValueOf(value)

		// Check if the value is its zero type.
		// False
		isTrue := !v.IsZero()
		if isTrue {
			count++
		}

	}

	return count
}

// Return truthy elements.
func MultipleTruthy[T any](t_args ...T) []T {
	var truthy []T

	// Count how
	for _, value := range t_args {
		v := reflect.ValueOf(value)

		isTrue := !v.IsZero()
		if isTrue {
			truthy = append(truthy, value)
		}
	}

	return truthy
}

// Append name to names slice if boolean is true.
func AppendWhen(t_bool bool, t_name string, t_names *[]string) bool {
	if t_bool {
		*t_names = append(*t_names, t_name)
	}

	return t_bool
}

func Arguments2TableDataSource(t_args Arguments) (td.TableDataSource, error) {
	var source td.TableDataSource

	// Selected supported source formats.
	var selected []string

	// Supported source formats.
	isCsv := AppendWhen(t_args.Csv, "csv", &selected)
	isJson := AppendWhen(t_args.Json, "json", &selected)
	isExcel := AppendWhen(t_args.Excel, "excel", &selected)

	// Verify that only one source format is selected.
	count := CountTruthy(isCsv, isJson, isExcel)
	if count > 1 {
		// TODO: This adds a space at the end but atm who cares, strip later.
		truthyStr := strings.Join(selected, ", ")
		truthyQuoted := u.Quote(truthyStr)
		errStr := fmt.Sprintf("Multiple input formats selected: %s.", truthyQuoted)
		err := errors.New(errStr)

		return source, err
	}

	// Actually set the source to return.
	if isCsv {
		source = td.CSV
	} else if isJson {
		source = td.JSON
	} else if isExcel {
		source = td.EXCEL
	} else {
		// TODO: Maybe just return an error and then have callee handle it?
		u.Logf("No input format was selected assuming CSV.")
		source = td.CSV
	}

	// TODO: Log end result.

	return source, nil
}

func readInputBuffer(t_args Arguments, t_reader io.Reader) (td.TableData, error) {
	var table td.TableData
	var buffer bytes.Buffer

	// Copy from stdin.
	bytesWritten, err := io.Copy(&buffer, t_reader)
	u.Logf("Bytes read from reader: %d", bytesWritten)
	if err != nil {
		return table, err
	}

	// Determine the source input..
	source, err := Arguments2TableDataSource(t_args)
	if err != nil {
		return table, err
	}

	// Initialize the table data.
	table, err = td.InitTableData(buffer, source)
	if err != nil {
		return table, err
	}

	return table, nil
}

func readInputStdin(t_args Arguments) ([]td.TableData, error) {
	var tables []td.TableData

	table, err := readInputBuffer(t_args, os.Stdin)
	if err != nil {
		return tables, err
	}

	// Return as a slice.
	tables = append(tables, table)

	return tables, nil
}

func readInputFiles(t_args Arguments) ([]td.TableData, error) {
	var tables []td.TableData

	for _, filePath := range t_args.InputFiles {
		// Read file.	// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			return tables, err
		}
		defer file.Close()

		table, err := readInputBuffer(t_args, file)
		if err != nil {
			return tables, err
		}

		// Append to the rest of the tables.
		tables = append(tables, table)
	}

	return tables, nil
}

func readInput(t_args Arguments) ([]td.TableData, error) {
	var (
		tables []td.TableData
		err    error
	)

	if t_args.FromStdin {
		tables, err = readInputStdin(t_args)
	} else {
		tables, err = readInputFiles(t_args)
	}

	return tables, err
}
