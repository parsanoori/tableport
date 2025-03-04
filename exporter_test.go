package tableport

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestToPDF(t *testing.T) {
	// Test code here
	type Person struct {
		FirstName string `pdf:"نام"`
		LastName  string `pdf:"نام خانوادگی"`
		Age       int    `pdf:"سن"`
		A         string `pdf:"A"`
		B         string `pdf:"B"`
		C         string `pdf:"C"`
		D         string `pdf:"D"`
	}

	p := []Person{
		{"John", "Doe", 30, "A", "B", "C", "D"},
		{"Jane", "Doe", 29, "A", "B", "C", "D"},
		{"John", "Smith", 40, "A", "1213123213Bqwedrftghjwertgyhjukl", "C", "D"},
		{"Jane", "Smith", 39, "A", "B", "C", "D"},
		{"John", "Johnson", 50, "A", "B", "C", "D"},
		{"Jane", "Johnson", 49, "A", "B", "C", "D"},
		{"John", "Williams", 60, "A", "B", "C", "D"},
		{"Jane", "Williams", 59, "A", "B", "C", "D"},
		{"John", "Brown", 70, "A", "B", "C", "D"},
		{"Jane", "Brown", 69, "A", "B", "C", "D"},
	}

	pdf, err := ToPDF(p)
	assert.Nil(t, err)

	// write the pdf to a file
	err = os.WriteFile("test.pdf", []byte(pdf), 0644)
	assert.Nil(t, err)

}
