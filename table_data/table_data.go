package table_data

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	u "github.com/soerlemans/table/util"
	// "github.com/xuri/excelize/v2"
)

// Defines different supported table source formats.
type TableDataSource int

const (
	CSV TableDataSource = iota
	JSON
	EXCEL
)

type TableDataRow = []string

// Internal representation of the table.
type TableData struct {
	// Map strings to indexes.
	// Used for getting column data by header name.
	// We zero index instead of glorious Awk.
	HeadersMap map[string]int
	Headers    TableDataRow

	// TODO: See if this plays out or if we want to use []string.
	// And do the field separation on the fly.
	RowsData []TableDataRow
}

func (this *TableData) HeaderLength() int {
	return len(this.Headers)
}

// Return the amount of rows.
func (this *TableData) RowLength() int {
	return len(this.RowsData)
}

// Get a specific cell by just indices.
func (this *TableData) CellByIndices(t_row int, t_col int) (string, error) {
	var cell string
	defer func() { u.Logf("CellByIndices: %s.", u.Quote(cell)) }()

	if t_row < this.RowLength() {
		rowSlice := this.RowsData[t_row]
		if t_col < len(rowSlice) {
			cell = this.RowsData[t_row][t_col]
		} else {
			err := u.Errorf("Column index is out of bounds (%d).", t_col)
			return cell, err
		}
	} else {
		err := u.Errorf("Row index is out of bounds (%d).", t_row)
		return cell, err
	}

	return cell, nil
}

// Get a specific cell by index and column name.
func (this *TableData) CellByColName(t_row int, t_name string) (string, error) {
	// Fetch the index of a header by the resolving the header map.
	t_col, ok := this.HeadersMap[t_name]
	if !ok {
		errStr := fmt.Sprintf("Non existent column name: %s.", t_name)
		err := errors.New(errStr)
		return "", err
	}

	// Now get them by indices.
	return this.CellByIndices(t_row, t_col)
}

func (this *TableData) GetRow(t_row int) (TableDataRow, error) {
	var row TableDataRow
	defer func() { u.Logf("GetRow: %s.", u.Quote(row)) }()

	if t_row < this.RowLength() {
		rowArray := this.RowsData[t_row]

		// Append to the line string, for each cell.
		for _, cell := range rowArray {
			row = append(row, cell)
		}

	} else {
		err := u.Errorf("Row index is out of bounds (%d).", t_row)
		return row, err
	}

	return row, nil
}

func (this *TableData) RowAsStr(t_row int) (string, error) {
	var row string
	defer func() { u.Logf("Row2Str: %s.", u.Quote(row)) }()

	if t_row < this.RowLength() {
		rowArray := this.RowsData[t_row]

		// Append to the line string, for each cell.
		var sep string
		for _, cell := range rowArray {
			row += fmt.Sprintf("%s%s", sep, cell)

			sep = ", "
		}

	} else {
		err := u.Errorf("Row index is out of bounds (%d).", t_row)
		return row, err
	}

	return row, nil
}

// Convert a matrix of strings into a TableData struct.
func matrix2TableData(t_matrix [][]string) TableData {
	var table TableData

	// FIXME: There is no fix / checking yet if all columns have the same length.
	// This is an issue for the JSON input.

	// Initialize the headers map.
	table.HeadersMap = make(map[string]int)

	// An empty csv file is also valid csv so dont error.
	if len(t_matrix) > 1 {
		// Initialize headers:
		for index, header := range t_matrix[0] {
			table.Headers = append(table.Headers, header)
			table.HeadersMap[header] = index

			u.Logf("Header: (%d:%s)", index, header)
		}

		// Initialize fields:
		for index, row := range t_matrix[1:] {
			u.Logf("Added row(%d): %v", index, row)
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
	defer func() { u.LogStructName("records", records, u.ETC80) }()

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
	u.Println("raw: ", raw)

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

	err := errors.New("TODO: initTableData does not support, EXCEL yet.")
	return table, err

	// return table, nil

}

func InitTableData(t_buffer bytes.Buffer, t_source TableDataSource) (TableData, error) {
	var table TableData
	var err error

	defer func() { u.LogStructName("InitTableData", table, u.ETC80) }()

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
