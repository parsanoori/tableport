package tableport

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func getData() interface{} {
	// Test code here
	type Person struct {
		FirstName string `pdf:"نام" excel:"نام"`
		LastName  string `pdf:"نام خانوادگی" excel:"نام خانوادگی"`
		Age       int    `pdf:"سن" excel:"سن"`
	}

	type Student struct {
		Person
		StudentID int `pdf:"شماره دانشجویی" excel:"شماره دانشجویی"`
	}

	type TA struct {
		Student
		Course string `pdf:"درس" excel:"درس"`
	}

	s := TA{
		Student: Student{
			Person: Person{
				FirstName: "Joe",
				LastName:  "Doe",
				Age:       25,
			},
			StudentID: 423424,
		},
		Course: "Math",
	}

	p := []TA{s}

	return p
}

func TestToPDF(t *testing.T) {
	p := getData()
	pdf, err := ToPDF(p)
	assert.Nil(t, err)

	// write the pdf to a file
	err = os.WriteFile("test.pdf", []byte(pdf), 0644)
	assert.Nil(t, err)
}

func TestToExcel(t *testing.T) {
	p := getData()
	excel, err := ToExcel(p)
	assert.Nil(t, err)

	// write the excel to a file
	err = os.WriteFile("test.xlsx", []byte(excel), 0644)
	assert.Nil(t, err)
}
