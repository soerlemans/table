package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
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
	// Map strings to indexes.
	// Used for getting column data by header name.
	// We zero index instead of glorious Awk.
	HeadersMap map[string]int
	Headers    []string

	// TODO: See if this plays out or if we want to use []string.
	// And do the field separation on the fly.
	RowsData [][]string
}

func (this *TableData) CellByColName(t_row int, t_name string) (string, error) {
	index, ok := this.HeadersMap[t_name]
	if !ok {
		errStr := fmt.Sprintf("Non existent column name: %s.", t_name)
		err := errors.New(errStr)
		return "", err
	}

	// this.RowsData[t_row][index]
}

// Convert a matrix of strings into a TableData struct.
func matrix2TableData(t_matrix [][]string) TableData {
	var table TableData

	// Initialize the headers map.
	table.HeadersMap = make(map[string]int)

	// An empty csv file is also valid csv so dont error.
	if len(t_matrix) > 1 {
		// Initialize headers:
		for index, header := range t_matrix[0] {
			table.Headers = append(table.Headers, header)
			table.HeadersMap[header] = index

			util.Logf("Header: (%d:%s)", index, header)
		}

		// Initialize fields:
		for index, row := range t_matrix[1:] {
			util.Logf("Added row(%d): %v", index, row)
			table.RowsData = append(table.RowsData, row)
		}
	}

	return table
}

func parseCsv(t_reader io.Reader) (TableData, error) {
	var table TableData

	reader := csv.NewReader(t_reader)
	records, err := reader.ReadAll()
	if err != nil {
		return table, err
	}
	defer util.LogStructName("records", records, util.ETC80)

	table = matrix2TableData(records)

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

	defer util.LogStructName("initTableData", table, util.ETC80)

	switch t_source {
	case CSV:
		table, err = parseCsv(&t_buffer)

	case JSON:
		table, err = parseJson(&t_buffer)

	case EXCEL:
		table, err = parseExcel(&t_buffer)
	}

	return table, err
}
