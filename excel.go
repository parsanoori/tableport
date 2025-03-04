package tableport

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/xuri/excelize/v2"
)

// ToExcel converts a slice of structs to an Excel file and returns it as bytes.
func ToExcel(in interface{}) ([]byte, error) {
	// Dereference pointer if input is a pointer
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Validate input is a slice
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input must be a slice of structs or slice of struct pointers")
	}

	// Ensure slice is not empty
	if v.Len() == 0 {
		return nil, fmt.Errorf("input slice is empty")
	}

	// Get the struct type (handle pointer elements)
	elemType := v.Index(0).Type()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}
	if elemType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a slice of structs or slice of struct pointers")
	}

	// Create a new Excel file
	f := excelize.NewFile()
	sheet := "Sheet1"

	// Extract field names as headers
	var headers []string
	var fieldIndexes []int
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		header := field.Tag.Get("excel")
		if header == "" {
			header = field.Name
		}
		headers = append(headers, header)
		fieldIndexes = append(fieldIndexes, i)
	}

	// Write headers to the first row
	for col, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+col) // A1, B1, C1, etc.
		f.SetCellValue(sheet, cell, header)
	}

	// Write struct values to subsequent rows
	for row := 0; row < v.Len(); row++ {
		item := v.Index(row)
		if item.Kind() == reflect.Ptr {
			item = item.Elem() // Dereference pointer if it's a pointer
		}

		for col, fieldIndex := range fieldIndexes {
			cell := fmt.Sprintf("%c%d", 'A'+col, row+2) // A2, B2, C2, etc.
			f.SetCellValue(sheet, cell, item.Field(fieldIndex).Interface())
		}
	}

	// Save to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
