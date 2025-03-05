package tableport

import (
	_ "embed"
	"errors"
	"github.com/javad-majidi/farsi-reshaper"
	"github.com/signintech/gopdf"
)

var (
	widthFactor        = 7
	widthFactorPersian = 5
	rowSize            = 20
)

func getColWidth(s string) int {
	if isAllPersian(s) {
		return len(s)*widthFactorPersian + 10
	}
	return len(s)*widthFactor + 10
}

// ToPDF converts a slice or an array of structs to a PDF table and gives the bytes of the PDF
func ToPDF(in interface{}) (res []byte, err error) {
	keys, values, err := flatten(in, "pdf")
	if err != nil {
		return
	}

	if len(keys) == 0 {
		err = errors.New("input should be a slice or an array of structs having fields")
		return
	}

	columnsWidths := make([]int, len(keys))
	for i, key := range keys {
		columnsWidths[i] = getColWidth(key) + 10
	}
	for i := range values {
		for j, value := range values[i] {
			columnsWidths[j] = max(columnsWidths[j], getColWidth(value)+10)
		}
	}
	for i := range columnsWidths {
		columnsWidths[i] = max(columnsWidths[i], 40)
	}

	columnsCount := len(keys)

	l := len(values)
	height := rowSize * (l + 1)
	var width int
	for _, w := range columnsWidths {
		width += w
	}

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
		table.AddColumn(FarsiReshaper.Reshape(keys[i]), float64(columnsWidths[i]), "center")
	}

	// add rows to the table
	for i := 0; i < l; i++ {
		var row []string
		for j := 0; j < columnsCount; j++ {
			v := FarsiReshaper.Reshape(values[i][j])
			row = append(row, v)
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

	// set the style for table header
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
