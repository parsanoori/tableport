package tableport

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

// ToExcel converts a slice of structs to an Excel file and returns it as bytes.
func ToExcel(in interface{}) (res []byte, err error) {
	keys, values, err := flatten(in, "excel")
	if err != nil {
		return
	}

	if len(keys) == 0 {
		err = errors.New("input should be a slice or an array of structs having fields")
		return
	}

	// Create a new Excel file
	f := excelize.NewFile()
	sheet := "Sheet1"

	// Write headers to the first row
	for col, header := range keys {
		cell := fmt.Sprintf("%c1", 'A'+col) // A1, B1, C1, etc.
		err = f.SetCellValue(sheet, cell, header)
		if err != nil {
			return
		}
	}

	// Write values to the rest of the rows
	for row, rowValues := range values {
		for col, value := range rowValues {
			cell := fmt.Sprintf("%c%d", 'A'+col, row+2) // A2, B2, C2, etc.
			err = f.SetCellValue(sheet, cell, value)
			if err != nil {
				return
			}
		}
	}

	// Save to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	res = buf.Bytes()
	return
}
