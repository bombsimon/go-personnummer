package swessn

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	var (
		three = 3
		four  = 4
	)

	cases := []struct {
		input   string
		output  *Parsed
		wantErr bool
	}{
		{
			input:   "😸",
			wantErr: true,
		},
		{
			input: "8001013294",
			output: &Parsed{
				Century:      0,
				Year:         80,
				Month:        1,
				Day:          1,
				Serial:       329,
				ControlDigit: &four,
				Divider:      DividerMinus,
			},
		},
		{
			input: "18800101+3294",
			output: &Parsed{
				Century:      1800,
				Year:         80,
				Month:        1,
				Day:          1,
				Serial:       329,
				ControlDigit: &four,
				Divider:      DividerPlus,
			},
		},
		{
			input: "800101329",
			output: &Parsed{
				Century:      0,
				Year:         80,
				Month:        1,
				Day:          1,
				Serial:       329,
				ControlDigit: &four,
				Divider:      DividerMinus,
			},
		},
		{
			input: "090314-6603",
			output: &Parsed{
				Century:      0,
				Year:         9,
				Month:        3,
				Day:          14,
				Serial:       660,
				ControlDigit: &three,
				Divider:      DividerMinus,
			},
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("input %s", tc.input), func(t *testing.T) {
			result, err := Parse(tc.input)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)

				return
			}

			assert.Equal(t, tc.output, result)
		})
	}
}

func TestParsed_Valid(t *testing.T) {
	cases := []struct {
		description   string
		input         string
		invalidFormat bool
		validPnr      bool
		validOrg      bool
	}{
		{
			description:   "invalid input",
			input:         "😸",
			invalidFormat: true,
			validPnr:      false,
			validOrg:      false,
		},
		{
			description: "invalid input",
			input:       "000000-0000",
			validPnr:    false,
			validOrg:    false,
		},
		{
			description: "invalid input",
			input:       "000000-0001",
			validPnr:    false,
			validOrg:    false,
		},
		{
			description: "regular personal number last century",
			input:       "8001013294",
			validPnr:    true,
			validOrg:    false,
		},
		{
			description: "regular personal this century",
			input:       "090314-6603",
			validPnr:    true,
			validOrg:    false,
		},
		{
			description: "regular personal in the future",
			input:       "21800101-3294",
			validPnr:    true,
			validOrg:    false,
		},
		{
			description: "regular personal with + sign",
			input:       "800101+3294",
			validPnr:    true,
			validOrg:    false,
		},
		{
			description: "valid organization number",
			input:       "556703-7485",
			validPnr:    false,
			validOrg:    true,
		},
		{
			description: "valid organization number other form",
			input:       "252002-6135",
			validPnr:    false,
			validOrg:    true,
		},
		{
			description: "valid luhn but not valid personal or organization number",
			input:       "056703-7486",
			validPnr:    false,
			validOrg:    false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			parsed, err := Parse(tc.input)

			if tc.invalidFormat {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.validPnr, parsed.ValidPerson())
			assert.Equal(t, tc.validOrg, parsed.ValidOrganization())
		})
	}
}

func TestStringFromInterface(t *testing.T) {
	cases := []struct {
		input  interface{}
		output string
	}{
		{
			input:  "12",
			output: "12",
		},
		{
			input:  []byte("12"),
			output: "12",
		},
		{
			input:  12,
			output: "12",
		},
		{
			input:  int32(12),
			output: "12",
		},
		{
			input:  int64(12),
			output: "12",
		},
		{
			input:  float32(12),
			output: "12",
		},
		{
			input:  float64(12),
			output: "12",
		},
		{
			input:  Divider("12"),
			output: "",
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("input %s", tc.input), func(t *testing.T) {
			str := StringFromInterface(tc.input)

			assert.Equal(t, tc.output, str)
		})
	}
}
