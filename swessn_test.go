package swessn

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidSSN(t *testing.T) {
	y, _, _ := time.Now().Date()
	ssn := New("640101-4136")

	assert.Equal(t, true, ssn.Valid())
	assert.Equal(t, y-1964, ssn.Age())

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
		"800161-3294":   "19800161-3294",
	}

	for in, expected := range valid {
		assert.Equal(t, expected, New(in).String())
		assert.Equal(t, true, New(in).Valid())
	}
}

func TestInvalidSSN(t *testing.T) {
	badDate := New("880435-3300")

	assert.Equal(t, "00000000-0001", badDate.Formatted)
	assert.Equal(t, false, badDate.Valid())

	badFormat := New("zeebra")

	assert.Equal(t, "00000000-0001", badFormat.String())
	assert.Equal(t, 0, badFormat.Age())
	assert.Equal(t, true, badFormat.IsOfAge(0))
}

func TestAge(t *testing.T) {
	y, m, d := time.Now().AddDate(-20, -1, 0).Date()
	canBuyAtSystembolaget := New(fmt.Sprintf("%04d%02d%02d-0000", y, m, d))

	assert.Equal(t, 20, canBuyAtSystembolaget.Age())
	assert.Equal(t, true, canBuyAtSystembolaget.IsOfAge(20))
	assert.Equal(t, false, canBuyAtSystembolaget.IsOfAge(21))

}

func TestGenerate(t *testing.T) {
	valid := map[string]*SSN{
		"19901110-": Generate(1990, 11, 10, "m"),
		"19880501-": Generate(1988, 5, 1, "f"),
		"20880501-": Generate(2088, 5, 1, "m"),
		"19880401-": Generate(88, 4, 1, "m"),
	}

	for expected, generated := range valid {
		assert.Equal(t, true, strings.HasPrefix(generated.String(), expected))
		assert.Equal(t, true, generated.Valid())
	}

	invalid := map[string]*SSN{
		"00000000-":  Generate(1990, 11, 10, "q"),
		"00000000-0": Generate(1990, 30, 40, "m"),
	}

	for expected, generated := range invalid {
		assert.Equal(t, true, strings.HasPrefix(generated.String(), expected))
		assert.Equal(t, false, generated.Valid())
	}

	for range make([]int, 100) {
		ssn := GenerateAny()

		assert.Equal(t, true, ssn.Valid())
		assert.Equal(t, true, ssn.Age() > 2)
	}
}
