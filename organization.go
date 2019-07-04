package swessn

// Organization represents a parsed string to be used in the context of an
// organization.
type Organization struct {
	*Parsed
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
		Parsed: parsed,
	}

	return organisation, nil
}

// IsValidOrganization returns if the parsed organization string is valid.
func IsValidOrganization(input interface{}) bool {
	nr := StringFromInterface(input)

	org, err := NewOrganization(nr)
	if err != nil {
		return false
	}

	if org.Divider == DividerPlus {
		return false
	}

	return org.Valid()
}

// Valid returns if the parsed organization string is valid.
func (o *Organization) Valid() bool {
	if o.Month < 20 {
		return false
	}

	return o.Parsed.Valid()
}
