package tableport

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"unicode"

	"github.com/javad-majidi/farsi-reshaper"
	"github.com/signintech/gopdf"
)

//go:embed assets/fonts/Vazir.ttf
var myFont []byte

var (
	widthFactor        = 7
	widthFactorPersian = 5
	rowSize            = 20
)

var embeddedFont []byte

var tempFontPath string // Store the temp font path

// init() runs when the package is imported
func init() {
	tempDir := os.TempDir() // Get OS temp directory
	tempFontPath = filepath.Join(tempDir, "embedded_font.ttf")

	// Write font to the temp file
	err := os.WriteFile(tempFontPath, embeddedFont, 0644)
	if err != nil {
		panic("Failed to write the embedded font to a temp file")
	}

	// Optional: Delete the temp file when the program exits
	go func() {
		<-make(chan struct{}) // Prevent premature cleanup (adjust if needed)
		os.Remove(tempFontPath)
	}()
}

// isPersian checks if a rune belongs to the Persian character set
func isPersian(r rune) bool {
	// Persian characters in Unicode ranges
	return (r >= 0x0600 && r <= 0x06FF) || // Arabic (including Persian) block
		(r >= 0x0750 && r <= 0x077F) || // Arabic Supplement
		(r >= 0x08A0 && r <= 0x08FF) || // Arabic Extended-A
		(r >= 0xFB50 && r <= 0xFDFF) || // Arabic Presentation Forms-A
		(r >= 0xFE70 && r <= 0xFEFF) || // Arabic Presentation Forms-B
		(r >= 0x10E60 && r <= 0x10E7F) // Rumi Numeral Symbols (used in Persian)
}

// isAllPersian checks if all characters in a string are Persian
func isAllPersian(s string) bool {
	for _, r := range s {
		if !isPersian(r) && !unicode.IsSpace(r) { // Allow spaces
			return false
		}
	}
	return true
}

func getColWidth(s string) int {
	if isAllPersian(s) {
		return len(s)*widthFactorPersian + 10
	}
	return len(s)*widthFactor + 10
}

// ToPDF converts a slice or an array of structs to a PDF table and gives the bytes of the PDF
func ToPDF(in interface{}) (res []byte, err error) {
	// iterate over the fields of the struct and get the column names using reflect package
	t := reflect.TypeOf(in)

	// if the type is then deference it
	if t.Kind() == reflect.Ptr {
		in = reflect.ValueOf(in).Elem().Interface()
		t = reflect.TypeOf(in)
	}

	// the type should a slice or an array of structs
	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		return []byte{}, errors.New("input should be a slice or an array")
	}

	// get the type of the struct
	t = t.Elem()

	// get the length of the slice or array
	l := reflect.ValueOf(in).Len()

	// the type should be a struct
	if t.Kind() != reflect.Struct {
		return []byte{}, errors.New("input should be a slice or an array of structs")
	}

	columnsCount := t.NumField()

	columnNames := make([]string, columnsCount)

	// column width is columns length * width factor
	columnsWidths := make([]int, columnsCount)

	// iterate over the fields of the struct
	for i := 0; i < columnsCount; i++ {
		// the column name is the tag pdf of the field or the field name
		var colName string
		colName = t.Field(i).Tag.Get("pdf")
		if colName == "" {
			colName = t.Field(i).Name
		}

		columnsWidths[i] = getColWidth(colName)

		colName = FarsiReshaper.Reshape(colName)
		columnNames[i] = colName
	}

	var width, height int

	// iterate over the elements of the slice or array and fill the columnsLengths
	for i := 0; i < l; i++ {
		// iterate over the fields of the struct
		for j := 0; j < columnsCount; j++ {
			// get the value of the field
			v := reflect.ValueOf(in).Index(i).Field(j)
			// get the string representation of the value
			s := fmt.Sprintf("%v", v.Interface())
			// if the length of the string is greater than the width of the column then update the width of the column
			columnsWidths[j] = max(columnsWidths[j], getColWidth(s))
		}
	}

	for i := 0; i < columnsCount; i++ {
		columnsWidths[i] = max(40, columnsWidths[i])
		width += columnsWidths[i]
	}
	height = rowSize * (l + 1)

	// create a new pdf
	pdf := &gopdf.GoPdf{}
	// set the page size
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: float64(width + 20), H: float64(height + 20)}})
	// add a new page
	pdf.AddPage()

	err = pdf.AddTTFFont("Vazir", tempFontPath)
	if err != nil {
		return []byte{}, err
	}
	err = pdf.SetFont("Vazir", "", 14)
	if err != nil {
		return []byte{}, err
	}

	// set the starting position for the table
	tableStartX := 10
	tableStartY := 10

	table := pdf.NewTableLayout(float64(tableStartX), float64(tableStartY), float64(rowSize), l)

	// add columns to the table
	for i := 0; i < columnsCount; i++ {
		table.AddColumn(columnNames[i], float64(columnsWidths[i]), "center")
	}

	// add rows to the table
	for i := 0; i < l; i++ {
		var row []string
		for j := 0; j < columnsCount; j++ {
			v := reflect.ValueOf(in).Index(i).Field(j)
			s := fmt.Sprintf("%v", v.Interface())
			s = FarsiReshaper.Reshape(s)
			row = append(row, s)
		}
		table.AddRow(row)
	}

	// set the style for the table
	table.SetTableStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top: true, Left: true, Right: true, Bottom: true,
			Width:    0.4,
			RGBColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		},
		Font:      "Vazir",
		FontSize:  14,
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		FillColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
	})

	// set the stype for table header
	table.SetHeaderStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top: true, Left: true, Right: true, Bottom: true,
			Width:    0.5,
			RGBColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		},
		Font:      "Vazir",
		FontSize:  14,
		FillColor: gopdf.RGBColor{R: 240, G: 240, B: 240},
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
	})

	// set the style for table cells
	table.SetCellStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top: true, Left: true, Right: true, Bottom: true,
			Width:    0.4,
			RGBColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		},
		Font:      "Vazir",
		FontSize:  14,
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
	})

	// draw the table
	err = table.DrawTable()
	if err != nil {
		return []byte{}, err
	}

	// get the bytes of the pdf
	b := pdf.GetBytesPdf()

	return b, nil
}
