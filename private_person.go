package swessn

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const (
	// https://www.skatteverket.se/privat/skatter/internationellt/bosattutomlands/samordningsnummer.4.53a97fe91163dfce2da80001279.html
	minCoordinationNumber = 60
)

// Gender represents a biological gender represented in a Swedish social
// security number.
type Gender int

const (
	Male Gender = iota
	Female
)

// Person represents what can be told about a person based on the social
// security number.
type Person struct {
	*Parsed
	Date           time.Time
	IsCoordination bool
	County         County
	Gender         Gender
}

// NewPerson parses and returns a pointer to a Person based on the input. If the
// input cannot be parsed an error will be returned. Upon creating a Person the
// century, date and county will be set.
func NewPerson(input string) (*Person, error) {
	parsed, err := Parse(input)
	if err != nil {
		return nil, err
	}

	return NewPersonFromParsed(parsed)
}

// NewPersonFromParsed returns a new person from a Parsed type. This may be used
// to skip parsing multiple times if a string should be tested as Parsed,
// Organization or Person.
func NewPersonFromParsed(parsed *Parsed) (*Person, error) {
	person := &Person{
		Parsed: parsed,
		Gender: GenderFromSerial(parsed.Serial),
	}

	if person.Day > minCoordinationNumber {
		person.Day %= 60
		person.IsCoordination = true
	}

	if err := person.SetCentury(); err != nil {
		return nil, err
	}

	if err := person.SetDate(); err != nil {
		return nil, err
	}

	if err := person.SetCounty(); err != nil {
		return nil, err
	}

	return person, nil
}

// IsValidPerson returns if the parsed person string is valid.
func IsValidPerson(input interface{}) bool {
	nr := StringFromInterface(input)

	person, err := NewPerson(nr)
	if err != nil {
		return false
	}

	return person.Valid()
}

// Valid returns if the parsed person string is valid.
func (p *Person) Valid() bool {
	if err := p.SetDate(); err != nil {
		return false
	}

	return p.Parsed.Valid()
}

// String returns the string representation of a person.
func (p *Person) String() string {
	return fmt.Sprintf(
		"%02d%02d%02d%s%03d%d",
		p.Year, p.Month, p.Day,
		p.Divider, p.Serial, p.ControlDigit,
	)
}

// SetCentury will update the century for the person based on the input data if
// no century was given.
//
// If a person is older than
// 100 years old the divider '+' is used. If no century value is given,
// calculate it with the following algorithm.
//  * If the year, month and date has passed this century
//    - If the divider is + -> last century
//    - If the divider is - -> this century
//  * If the year, month and date has NOT passed
//    - If the divider is + -> use the century before the last
//    - If the divider is - -> Use the last century.
func (p *Person) SetCentury() error {
	// Nothing to do if already set.
	if p.Century != 0 {
		return nil
	}

	personDateWithCurrentCentury, err := time.Parse(
		"2006-01-02",
		fmt.Sprintf(
			"%02d%02d-%02d-%02d",
			time.Now().Year()/100, p.Year, p.Month, p.Day,
		),
	)

	if err != nil {
		return errors.New("invalid format")
	}

	// If the date passed have not passed, assumed they meant last century.
	// 830101-1110 -> 19140101-1110
	// 140101-1110 -> 20140101-1110
	if personDateWithCurrentCentury.After(time.Now()) {
		personDateWithCurrentCentury = personDateWithCurrentCentury.AddDate(-100, 0, 0)
	}

	// If the divider '+' is used this means that the person is
	if p.Divider == DividerPlus {
		personDateWithCurrentCentury = personDateWithCurrentCentury.AddDate(-100, 0, 0)
	}

	p.Century = personDateWithCurrentCentury.Year() / 100 * 100

	return nil
}

// SetDate will set a time.Time type on the Person struct.
func (p *Person) SetDate() error {
	if !p.Date.IsZero() {
		return nil
	}

	if err := p.SetCentury(); err != nil {
		return err
	}

	t, err := time.Parse(
		"2006-01-02",
		fmt.Sprintf(
			"%d-%02d-%02d",
			p.Century+p.Year,
			p.Month,
			p.Day,
		),
	)

	if err != nil {
		return err
	}

	p.Date = t

	return nil
}

// SetCounty will set the count on the Person struct.
func (p *Person) SetCounty() error {
	if p.Century+p.Year > 1990 {
		return nil
	}

	c, err := CountyFromSerial(p.Serial)
	if err != nil {
		return err
	}

	p.County = c

	return nil
}

// Age returns the age of a person with a given personal number based on today's
// date (UTC+0).
func (p *Person) Age() int {
	if err := p.SetDate(); err != nil {
		panic(err)
	}

	duration := time.Since(p.Date)

	return int(math.Floor(duration.Hours() / 24 / 365))
}

// IsOfAge checks if the age of a person with a given social security number has
// been reached.
func (p *Person) IsOfAge(age int) bool {
	return p.Age() >= age
}

// Male returns true if the social security number serial number is uneven.
func (p *Person) Male() bool {
	return p.Gender == Male
}

// Female returns true if the social security number serial number is even.
func (p *Person) Female() bool {
	return p.Gender == Female
}

// GenderFromSerial will calculate gender from serial number. If the last digit
// in the serial number is even it's a female otherwise it's a male.
func GenderFromSerial(serial int) Gender {
	if (serial%10)%2 == 0 {
		return Female
	}

	return Male
}

// Generate will generate a valid Swedish social security number
// based on passed year, month, day and sex.
func Generate(date time.Time, sex Gender) (*Person, error) {
	if sex != Male && sex != Female {
		return nil, errors.New("invalid gender")
	}

	rand.Seed(time.Now().UnixNano())

	sexIndications := map[Gender][]int{
		Male:   {1, 3, 5, 7, 9},
		Female: {2, 4, 6, 8, 0},
	}

	randStart := rand.Intn(99)
	randSex := sexIndications[sex][rand.Intn(len(sexIndications[sex]))]
	randSerial, _ := strconv.Atoi(fmt.Sprintf("%02d%d", randStart, randSex))

	century := date.Year() / 100 * 100
	parsed := &Parsed{
		Century: century,
		Year:    date.Year() % century,
		Month:   int(date.Month()),
		Day:     date.Day(),
		Serial:  randSerial,
	}

	cs := parsed.LuhnChecksum()
	cd := parsed.LuhnControlDigit(cs)

	parsed.ControlDigit = &cd

	person := &Person{
		Parsed: parsed,
		Gender: sex,
	}

	if err := person.SetDate(); err != nil {
		return nil, err
	}

	return person, nil
}

// GenerateAny will generate a random valid Swedish social security number
// with a random date and random sex.
func GenerateAny() (*Person, error) {
	var (
		min   = time.Date(1974, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
		max   = time.Date(2014, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
		delta = max - min
		sec   = rand.Int63n(delta) + min
	)

	sexes := []Gender{Male, Female}

	return Generate(time.Unix(sec, 0), sexes[rand.Intn(len(sexes))])
}
