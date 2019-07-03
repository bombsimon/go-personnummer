package swessn

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Divider represents the divider between birth date and control digits.
type Divider string

const (
	DividerPlus  Divider = "+"
	DividerMinus Divider = "-"
	DividerNone  Divider = ""
)

// Parsed represents a parsed string. The fields are named as date parts but may
// be of other types in case of an organisation number or coordination number.
type Parsed struct {
	Century      int
	Year         int
	Month        int
	Day          int
	Serial       int
	ControlDigit *int
	Divider      Divider
}

func Parse(input string) (*Parsed, error) {
	validFormatRe := regexp.MustCompile(`^(\d{2})?(\d{2})(\d{2})(\d{2})([-+])?(\d{3})(\d)?$`)

	matches := validFormatRe.FindStringSubmatch(input)

	if len(matches) != 8 {
		return nil, errors.New("invalid format")
	}

	var (
		century, _ = strconv.Atoi(matches[1])
		year, _    = strconv.Atoi(matches[2])
		month, _   = strconv.Atoi(matches[3])
		day, _     = strconv.Atoi(matches[4])
		serial, _  = strconv.Atoi(matches[6])
		divider    = Divider(strings.ToUpper(matches[5]))
	)

	p := &Parsed{
		Year:    year,
		Month:   month,
		Day:     day,
		Serial:  serial,
		Divider: divider,
	}

	if century > 0 {
		p.Century = century * 100
	}

	if p.Divider == DividerNone {
		p.Divider = DividerMinus
	}

	if cd, err := strconv.Atoi(matches[7]); err == nil {
		p.ControlDigit = &cd
	}

	if p.ControlDigit == nil {
		cd := p.LuhnControlDigit(p.LuhnChecksum())
		p.ControlDigit = &cd
	}

	return p, nil
}

func (p *Parsed) LuhnChecksum() int {
	var (
		sum    = 0
		digits = fmt.Sprintf("%02d%02d%02d%03d", p.Year, p.Month, p.Day, p.Serial)
	)

	for i := range digits {
		digit, err := strconv.Atoi(string(digits[i]))
		if err != nil {
			panic("invalid luhn iteration value")
		}

		if i%2 == 0 {
			digit *= 2
		}

		if digit > 9 {
			digit -= 9
		}

		sum += digit
	}

	return sum
}

func (p *Parsed) LuhnControlDigit(cs int) int {
	checksum := 10 - (cs % 10)

	if checksum == 10 {
		return 0
	}

	return checksum
}

func (p *Parsed) Valid() bool {
	var (
		controlDigit = p.LuhnControlDigit(p.LuhnChecksum())
		cd           = controlDigit
	)

	if p.ControlDigit != nil {
		cd = *p.ControlDigit
	}

	return controlDigit == cd
}

func StringFromInterface(input interface{}) string {
	var nr string

	switch v := input.(type) {
	case string:
		nr = v
	case []byte:
		nr = string(v)
	case int:
		nr = strconv.Itoa(v)
	case int32:
		nr = strconv.Itoa(int(v))
	case int64:
		nr = strconv.Itoa(int(v))
	case float32:
		nr = strconv.Itoa(int(v))
	case float64:
		nr = strconv.Itoa(int(v))
	default:
		nr = ""
	}

	return nr
}
