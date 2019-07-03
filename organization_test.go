package swessn

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidOrganisation(t *testing.T) {
	cases := []struct {
		organizationNumber string
		valid              bool
	}{
		{organizationNumber: "556703-7485", valid: true},
		{organizationNumber: "556074-7569", valid: true},
		{organizationNumber: "252002-6135", valid: true},
		{organizationNumber: "8001013294", valid: false},
		{organizationNumber: "😸", valid: false},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s is %v", tc.organizationNumber, tc.valid), func(t *testing.T) {
			valid := IsValidOrganisation(tc.organizationNumber)

			assert.Equal(t, tc.valid, valid)
		})
	}
}
