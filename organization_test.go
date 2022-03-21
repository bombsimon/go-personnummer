package personnummer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsValidOrganisation(t *testing.T) {
	cases := []struct {
		organizationNumber string
		valid              bool
	}{
		{organizationNumber: "16556703-7485", valid: true},
		{organizationNumber: "556703-7485", valid: true},
		{organizationNumber: "556074-7569", valid: true},
		{organizationNumber: "252002-6135", valid: true},
		{organizationNumber: "056703-7486", valid: false},
		{organizationNumber: "8001013294", valid: false},
		{organizationNumber: "556703+7485", valid: false},
		{organizationNumber: "19556703-7485", valid: false},
		{organizationNumber: "😸", valid: false},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s is %v", tc.organizationNumber, tc.valid), func(t *testing.T) {
			valid := IsValidOrganization(tc.organizationNumber)

			assert.Equal(t, tc.valid, valid)
		})
	}
}

func TestOrganizationCorporateForm(t *testing.T) {
	cases := []struct {
		organizationNumber string
		corporateForm      CorporateForm
	}{
		{organizationNumber: "556703-7485", corporateForm: CorporateFormLimitedCompany},
		{organizationNumber: "252002-6135", corporateForm: CorporateFormStateCCMunicipalities},
		{organizationNumber: "802405-0190", corporateForm: CorporateFormIdealFoundation},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s is %s", tc.organizationNumber, tc.corporateForm), func(t *testing.T) {
			org, err := NewOrganization(tc.organizationNumber)

			require.NoError(t, err)

			assert.Equal(t, tc.corporateForm, org.CorporateForm)
		})
	}
}
