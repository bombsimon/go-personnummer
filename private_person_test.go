package swessn

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPerson(t *testing.T) {
	var (
		three = 3
		four  = 4
		d1, _ = time.Parse("2006-01-02", "1980-01-01")
	)

	cases := []struct {
		description string
		input       string
		output      *Person
		wantErr     bool
	}{
		{
			description: "invalid input",
			input:       "😸",
			wantErr:     true,
		},
		{
			description: "invalid date",
			input:       "19805101-3294",
			wantErr:     true,
		},
		{
			description: "valid without no century (year not occurred)",
			input:       "8001013294",
			output: &Person{
				Parsed: &Parsed{
					Century:      1900,
					Year:         80,
					Month:        1,
					Day:          1,
					Serial:       329,
					ControlDigit: &four,
					Divider:      DividerMinus,
				},
				Date:           d1,
				IsCoordination: false,
				County:         CountyK,
				Gender:         Male,
			},
		},
		{
			description: "valid without no century (year has occurred)",
			input:       "090314-6603",
			output: &Person{
				Parsed: &Parsed{
					Century:      2000,
					Year:         9,
					Month:        3,
					Day:          14,
					Serial:       660,
					ControlDigit: &three,
					Divider:      DividerMinus,
				},
				IsCoordination: false,
				County:         CountyA,
				Gender:         Female,
			},
		},
		{
			description: "valid without no century (year has occurred, plus divider)",
			input:       "090314+6603",
			output: &Person{
				Parsed: &Parsed{
					Century:      1900,
					Year:         9,
					Month:        3,
					Day:          14,
					Serial:       660,
					ControlDigit: &three,
					Divider:      DividerPlus,
				},
				IsCoordination: false,
				County:         CountyT,
				Gender:         Female,
			},
		},
		{
			description: "valid with century",
			input:       "198001013294",
			output: &Person{
				Parsed: &Parsed{
					Century:      1900,
					Year:         80,
					Month:        1,
					Day:          1,
					Serial:       329,
					ControlDigit: &four,
					Divider:      DividerMinus,
				},
				Date:           d1,
				IsCoordination: false,
				County:         CountyK,
				Gender:         Male,
			},
		},
		{
			description: "valid with no century last century",
			input:       "800101+3294",
			output: &Person{
				Parsed: &Parsed{
					Century:      1800,
					Year:         80,
					Month:        1,
					Day:          1,
					Serial:       329,
					ControlDigit: &four,
					Divider:      DividerPlus,
				},
				IsCoordination: false,
				County:         CountyK,
				Gender:         Male,
			},
		},
		{
			description: "valid with coordination number",
			input:       "800161-3294",
			output: &Person{
				Parsed: &Parsed{
					Century:      1900,
					Year:         80,
					Month:        1,
					Day:          1,
					Serial:       329,
					ControlDigit: &four,
					Divider:      DividerMinus,
				},
				IsCoordination: true,
				County:         CountyK,
				Gender:         Male,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			result, err := NewPerson(tc.input)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tc.output.Parsed, result.Parsed)
			assert.Equal(t, tc.output.IsCoordination, result.IsCoordination)
			assert.Equal(t, tc.output.County, result.County)
			assert.Equal(t, tc.output.Gender, result.Gender)

			if !tc.output.Date.IsZero() {
				assert.Equal(t, tc.output.Date, result.Date)
			}

			switch tc.output.Gender {
			case Male:
				assert.True(t, result.Male())
			case Female:
				assert.True(t, result.Female())
			}
		})
	}
}

func TestIsValidPerson(t *testing.T) {
	cases := []struct {
		pnr   string
		valid bool
	}{
		{pnr: "8001013294", valid: true},
		{pnr: "8001613294", valid: true},
		{pnr: "198001013294", valid: true},
		{pnr: "800101-3294", valid: true},
		{pnr: "19800101-3294", valid: true},
		{pnr: "090314-6603", valid: true},
		{pnr: "800101+3294", valid: true},
		{pnr: "18800101+3294", valid: true},
		{pnr: "15800101-3294", valid: true},
		{pnr: "158001013294", valid: true},
		{pnr: "21800101-3294", valid: true},
		{pnr: "218001013294", valid: true},
		{pnr: "800161-3294", valid: true},
		{pnr: "880435-3300", valid: false},
		{pnr: "00000000-0001", valid: false},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s is %v", tc.pnr, tc.valid), func(t *testing.T) {
			valid := IsValidPerson(tc.pnr)

			assert.Equal(t, tc.valid, valid)
		})
	}
}

func TestPerson_Age(t *testing.T) {
	d := time.Now().AddDate(-20, -1, 0)
	personWhoCanBuyAtSystembolaget, err := Generate(d, Male)

	require.NoError(t, err)

	assert.Equal(t, 20, personWhoCanBuyAtSystembolaget.Age())
	assert.Equal(t, true, personWhoCanBuyAtSystembolaget.IsOfAge(20))
	assert.Equal(t, false, personWhoCanBuyAtSystembolaget.IsOfAge(21))
}
