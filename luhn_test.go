package swessn

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
