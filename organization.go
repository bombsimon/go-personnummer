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

	return org.Valid()
}

// Valid returns if the parsed organization string is valid.
func (o *Organization) Valid() bool {
	if o.Month < 20 {
		return false
	}

	return o.Parsed.Valid()
}
