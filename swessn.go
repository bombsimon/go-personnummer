package swessn

/*
Swedish social security numbers are calculated with
the Luhn algorithm a.k.a. the mod 10 algorithm.

This package is used to validate these numbers to ensure
that a social security number, including the date within it,
is valid.

This package can also be used to generate valid SSNs to test
in your application.
*/

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

const (
	invalidSsn = "00000000-0001"
)

// SSN is a social security number (personal number)
type SSN struct {
	Original  string
	Formatted string
}

// New generates a new SSN type
func New(ssn string) *SSN {
	return &SSN{
		Original:  ssn,
		Formatted: format(ssn),
	}
}

func format(ssn string) string {
	var (
		extraReduce int        // Used to reduce 100 years if SSN contains '+' sign
		plusAsByte  byte = '+' // Byte representation of the '+' sig
		ssnBytes         = []byte(ssn)
	)

	// Make sure we at least got a valid format of date and control number.
	// Dates are allowed with or without century prefix.
	ssnRe := regexp.MustCompile(`^(\d{6}|\d{8})[+-]?(\d{4})$`)
	ssnGroups := ssnRe.FindSubmatch(ssnBytes)

	if len(ssnGroups) != 3 {
		return invalidSsn
	}

	ssnDate, ssnControlNo := ssnGroups[1], ssnGroups[2]

	// Check if a plus sign was present to indicate SSN from last century
	// and the person is over 100 years old (above 100 years).
	if ssnBytes[len(ssnBytes)-5] == plusAsByte {
		extraReduce = 100
	}

	// Add default century if not provided.
	// Will use current century if year has passed, otherwise last century.
	if len(ssnDate) == 6 {
		shortYear, _ := strconv.Atoi(string(ssnDate[:2]))
		currentYear, _, _ := time.Now().Date()
		currentCentury := (currentYear / 100) * 100

		if currentCentury+shortYear > currentYear {
			// The year has passed, should padd with current century.
			// Unless the SSN has a plus sign and the person is >= 100 years old.
			lastCenturyString := strconv.Itoa(((currentCentury - 100) - extraReduce) / 100)
			ssnDate = append([]byte(lastCenturyString), ssnDate...)
		} else {
			currentCenturyString := strconv.Itoa((currentCentury - extraReduce) / 100)
			ssnDate = append([]byte(currentCenturyString), ssnDate...)
		}
	}

	// remove co-ordination date for verification purposes.
	day, _ := strconv.Atoi(string(ssnDate[6:8]))
	noCoordinationSSNDate := []byte(fmt.Sprintf("%s%02d%s", ssnDate[:6], day%60, ssnDate[8:]))

	// Make sure the ssn date without co-ordination number can be parsed into a valid date.
	// Will support leap years and other handy things.
	_, err := time.Parse("20060102", string(noCoordinationSSNDate))
	if err != nil {
		return invalidSsn
	}

	return fmt.Sprintf("%s-%s", string(ssnDate), string(ssnControlNo))
}

// String ensures adaption of the stringer interface.
// The string returned will always be formated with full century
// and a dash separating the date and the control numbers.
func (s *SSN) String() string {
	return s.Formatted
}

// Age returns the age of a person with a given SSN based on todays date (UTC+0)
func (s *SSN) Age() int {
	ssnTime, err := time.Parse("20060102", s.Formatted[:8])
	if err != nil {
		return 0
	}

	duration := time.Now().Sub(ssnTime)

	return int(math.Floor(duration.Hours() / 24 / 365))
}

// IsOfAge checks if the age of a person with a given SSN is met.
func (s *SSN) IsOfAge(age int) bool {
	return s.Age() >= age
}

// Valid returns a boolean value telling if the social security number is vlaid.
// This is calculated with Luhn algorithm to ensure the checksum of the last digit
// in the social security number is valid.
func (s *SSN) Valid() bool {
	// Ommit all but 10 digits from the SSN
	ssn := fmt.Sprintf("%s%s", s.Formatted[2:8], s.Formatted[9:13])

	// Handle co-ordination numbers.
	day, _ := strconv.Atoi(ssn[4:6])
	ssn = fmt.Sprintf("%s%02d%s", ssn[:4], day%60, ssn[6:])

	// Add the Luhn iteration sum to the last digit and ensure
	// that we always get 0 rest from mod 10.
	sum, _ := strconv.Atoi(string(ssn[len(ssn)-1]))
	sum += calculateLuhnSum(ssn)

	return sum%10 == 0
}

func calculateLuhnSum(ssn string) int {
	sum := 0

	// Iterate through the SSN except the last digit
	for i := range ssn[:len(ssn)-1] {
		// Convert each index to integer
		digit, _ := strconv.Atoi(string(ssn[i]))

		// Multiply by two every second position in the SSN
		if i%2 == 0 {
			digit *= 2
		}

		// Values above 10 is not allowed and must be reduced by 9
		if digit > 9 {
			digit -= 9
		}

		sum += digit
	}

	return sum
}

func calculateLuhnChecksum(sum int) int {
	checksum := 10 - (sum % 10)
	if checksum == 10 {
		checksum = 0
	}

	return checksum
}

// Generate will generate a valid Swedish social security number
// based on passed year, month, day and sex where sex is passed
// as "m" for male and "f" for female.
func Generate(y, m, d int, sex string) *SSN {
	// Return invalid date if faulty sex is provided.
	if sex != "m" && sex != "f" {
		return New("")
	}

	rand.Seed(time.Now().UnixNano())

	date := fmt.Sprintf("%04d%02d%02d", y, m, d)
	sexIndications := map[string][]int{
		"m": {1, 3, 5, 7, 9},
		"f": {2, 4, 6, 8, 0},
	}

	randStart := rand.Intn(99)
	randSex := sexIndications[sex][rand.Intn(len(sexIndications[sex]))]

	controlNo := fmt.Sprintf("%02d%d0", randStart, randSex)
	ssn := fmt.Sprintf("%s%s", date[2:8], controlNo)

	sum := calculateLuhnSum(ssn)
	lastDigit := calculateLuhnChecksum(sum)

	ssn = fmt.Sprintf("%d%s%d", y, ssn[2:9], lastDigit)

	return New(ssn)
}

// GenerateAny will generate a random valid Swedish social security number
// with a random date and random sex.
func GenerateAny() *SSN {
	var (
		min   = time.Date(1974, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
		max   = time.Date(2014, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
		delta = max - min
		sec   = rand.Int63n(delta) + min
	)

	randomDate := time.Unix(sec, 0)
	year, month, day := randomDate.Date()

	sexes := []string{"m", "f"}

	return Generate(year, int(month), day, sexes[rand.Intn(len(sexes))])
}
