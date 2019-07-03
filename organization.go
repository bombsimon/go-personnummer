package swessn

type Organisation struct {
	*Parsed
}

func NewOrganisation(input string) (*Organisation, error) {
	parsed, err := Parse(input)
	if err != nil {
		return nil, err
	}

	organisation := &Organisation{
		Parsed: parsed,
	}

	return organisation, nil
}

func IsValidOrganisation(input interface{}) bool {
	nr := StringFromInterface(input)

	org, err := NewOrganisation(nr)
	if err != nil {
		return false
	}

	return org.Valid()
}

func (o *Organisation) Valid() bool {
	if o.Month < 20 {
		return false
	}

	return o.Parsed.Valid()
}
