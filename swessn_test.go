package swessn

import (
	"fmt"
	"strings"
	"testing"
	"time"

	. "gopkg.in/go-playground/assert.v1"
)

func TestValidSSN(t *testing.T) {
	y, _, _ := time.Now().Date()
	ssn := New("640101-4136")

	Equal(t, ssn.Valid(), true)
	Equal(t, ssn.Age(), y-1964)

	valid := map[string]string{
		"8001013294":    "19800101-3294",
		"198001013294":  "19800101-3294",
		"800101-3294":   "19800101-3294",
		"19800101-3294": "19800101-3294",
		"090314-6603":   "20090314-6603",
		"800101+3294":   "18800101-3294",
		"18800101+3294": "18800101-3294",
		"15800101-3294": "15800101-3294",
		"158001013294":  "15800101-3294",
		"21800101-3294": "21800101-3294",
		"218001013294":  "21800101-3294",
	}

	for in, expected := range valid {
		Equal(t, New(in).String(), expected)
		Equal(t, New(in).Valid(), true)
	}
}

func TestInvalidSSN(t *testing.T) {
	badDate := New("880435-3300")

	Equal(t, badDate.Formatted, "00000000-0001")
	Equal(t, badDate.Valid(), false)

	badFormat := New("zeebra")

	Equal(t, badFormat.String(), "00000000-0001")
	Equal(t, badFormat.Age(), 0)
	Equal(t, badFormat.IsOfAge(0), true)
}

func TestAge(t *testing.T) {
	y, m, d := time.Now().Date()
	canBuyAtSystembolaget := New(fmt.Sprintf("%04d%02d%02d-0000", y-20, m-1, d))

	Equal(t, canBuyAtSystembolaget.Age(), 20)
	Equal(t, canBuyAtSystembolaget.IsOfAge(20), true)
	Equal(t, canBuyAtSystembolaget.IsOfAge(21), false)

}

func TestGenerate(t *testing.T) {
	valid := map[string]*SSN{
		"19901110-": Generate(1990, 11, 10, "m"),
		"19880501-": Generate(1988, 5, 1, "f"),
		"20880501-": Generate(2088, 5, 1, "m"),
		"19880401-": Generate(88, 4, 1, "m"),
	}

	for expected, generated := range valid {
		Equal(t, strings.HasPrefix(generated.String(), expected), true)
		Equal(t, generated.Valid(), true)
	}

	invalid := map[string]*SSN{
		"00000000-":  Generate(1990, 11, 10, "q"),
		"00000000-0": Generate(1990, 30, 40, "m"),
	}

	for expected, generated := range invalid {
		Equal(t, strings.HasPrefix(generated.String(), expected), true)
		Equal(t, generated.Valid(), false)
	}

	for _ = range make([]int, 100) {
		ssn := GenerateAny()

		Equal(t, ssn.Valid(), true)
		Equal(t, ssn.Age() > 2, true)
	}
}
