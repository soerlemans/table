package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"

	"github.com/soerlemans/table/util"
	// "github.com/xuri/excelize/v2"
)

// Defines different supported table source formats.
type TableDataSource int

const (
	CSV TableDataSource = iota
	JSON
	EXCEL
)

// Internal representation of the table.
type TableData struct {
	// Keep track of string to index mapping.
	ColumnMap map[string]int

	Columns []string
	Rows    []string
}

func parseCsv(t_reader io.Reader) (TableData, error) {
	var table TableData

	reader := csv.NewReader(t_reader)
	records, err := reader.ReadAll()
	util.Logf("records: %+v", records)
	if err != nil {
		return table, err
	}

	return table, nil
}

func parseJson(t_reader io.Reader) (TableData, error) {
	var table TableData

	var raw []map[string]interface{}
	decoder := json.NewDecoder(t_reader)
	err := decoder.Decode(&raw)
	if err != nil {
		return table, err
	}

	if len(raw) == 0 {
		return table, nil
	}
	util.Println("raw: ", raw)

	// headers := []string{}
	// for k := range raw[0] {
	// 	headers = append(headers, k)
	// }

	// rows := [][]string{headers}
	// for _, obj := range raw {
	// 	row := []string{}
	// 	for _, h := range headers {
	// 		val := fmt.Sprintf("%v", obj[h])
	// 		row = append(row, val)
	// 	}
	// 	rows = append(rows, row)
	// }
	// return rows, nil

	return table, nil

}

func parseExcel(t_reader io.Reader) (TableData, error) {
	var table TableData

	err := errors.New("TODO: initTableMap does not support, EXCEL yet.")
	return table, err

	// return table, nil

}

func initTableData(t_buffer bytes.Buffer, t_source TableDataSource) (TableData, error) {
	var table TableData
	var err error

	switch t_source {
	case CSV:
		table, err = parseCsv(&t_buffer)

	case JSON:
		table, err = parseJson(&t_buffer)

	case EXCEL:
		table, err = parseExcel(&t_buffer)
	}

	util.Logf("initTableData: %+v.", table)

	return table, err
}
