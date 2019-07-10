package swessn

// CorporateForm indicates what form the company is of. This could be told by
// reading the first digit in the organization number. This is not 100%
// guaranteed to be correct according to Bolagsverket and Bolagsverket is also a
// source that doesn't list all of these types.
// See the following sources for information:
// https://bolagsverket.se/ff/foretagsformer/organisationsnummer-1.7902
// https://www.solidinfo.se/hjalp/organisationsnummer
// https://www.ageras.se/ordlista/organisationsnummer
// https://sv.wikipedia.org/wiki/Organisationsnummer#Organisationsnummer
type CorporateForm int

const (
	CorporateFormEstate CorporateForm = iota + 1
	CorporateFormStateCCMunicipalities
	CorporateFormForeign
	CorporateFormUnknown
	CorporateFormLimitedCompany
	CorporateFormSimpleCompany
	CorporateFormEconomicTenantAssociation
	CorporateFormIdealFoundation
	CorporateFormTradingPartnershipSimple
)

func (cf CorporateForm) String() string {
	switch cf {
	case CorporateFormEstate:
		return "Dödsbo"
	case CorporateFormStateCCMunicipalities:
		return "Stat, landsting, kommun, församling"
	case CorporateFormForeign:
		return "Utländska företag som bedriver näringsverksamhet eller äger fastigheter i Sverige"
	case CorporateFormUnknown:
		return "Okänt"
	case CorporateFormLimitedCompany:
		return "Aktiebolag"
	case CorporateFormSimpleCompany:
		return "Enkelt bolag"
	case CorporateFormEconomicTenantAssociation:
		return "Ekonomisk förening, bostadsrättsförening"
	case CorporateFormIdealFoundation:
		return "Ideell förening och stiftelse"
	case CorporateFormTradingPartnershipSimple:
		return "Handelsbolag, kommanditbolag och enkelt bolag"
	}

	return "Okänt"
}

// Organization represents a parsed string to be used in the context of an
// organization.
type Organization struct {
	*Parsed
	CorporateForm CorporateForm
}

// NewOrganization parses and returns a pointer to an Organization based on the
// input. If the input cannot be parsed an error will be returned.
func NewOrganization(input string) (*Organization, error) {
	parsed, err := Parse(input)
	if err != nil {
		return nil, err
	}

	return NewOrganizationFromParsed(parsed)
}

// NewOrganizationFromParsed returns a new organization from a Parsed type. This
// may be used to skip parsing multiple times if a string should be tested as
// Parsed, Organization or Person.
func NewOrganizationFromParsed(parsed *Parsed) (*Organization, error) {
	organisation := &Organization{
		Parsed:        parsed,
		CorporateForm: CorporateForm(parsed.Year / 10),
	}

	return organisation, nil
}

// IsValidOrganization returns if the parsed organization string is valid.
func IsValidOrganization(input interface{}) bool {
	nr := stringFromInterface(input)

	org, err := NewOrganization(nr)
	if err != nil {
		return false
	}

	return org.Valid()
}

// Valid returns if the parsed organization string is valid.
func (o *Organization) Valid() bool {
	// May only be prefixed with 16.
	if o.Century != 0 && o.Century != 1600 {
		return false
	}

	// Third digit ("month") must be >= 2
	if o.Month < 20 {
		return false
	}

	// Organization numbers may never be divided with `+`.
	if o.Divider == DividerPlus {
		return false
	}

	// May never start with leading 0.
	if o.Year < 10 {
		return false
	}

	return o.Parsed.Valid()
}
