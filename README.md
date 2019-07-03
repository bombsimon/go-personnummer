# Swedish Identification Numbers

[![Build Status](https://travis-ci.org/bombsimon/swedish-ssn.svg?branch=master)](https://travis-ci.org/bombsimon/swedish-ssn)
[![GoDoc](https://godoc.org/github.com/bombsimon/swedish-ssn?status.svg)](https://godoc.org/github.com/bombsimon/swedish-ssn)

This package aims to provide a toolbox to handle Swedish identification numbers
of three different types; social security numbers (or personal numbers),
coordination numbers and organization numbers. They are all validated using the
[Luhn algorithm](https://en.wikipedia.org/wiki/Luhn_algorithm) with some
exceptions to the coordination number.

In addition to the correct checksum calculated with the Luhn algorithm, the following rules applies:

* Digit between digits and control numbers may only be divided with `-`, or `+`
* A social security number must be a valid date
* An organization numbers third digit must be >= 2
* A coordination number must have a date where day value is > 60

## Extra data (private person)

Some extra data may be extracted from a social security number regarding the
person. The `Person` type holds and implements a few of these things.

* `IsCoordination` tells if the person has a coordination number
* `Date` is a `time.Time` type with the birth date
* `County` holds the county code for people born before 1990
* `Gender` holds wether the person was born `Male` or `Female`
* `Age()` can tell the persons age (in UTC timezone)
* `IsOfAge(n int)` can tell if the person is `n` (or above)
* `Male()` is true if it's a `Male`
* `Female()` is true if it's a `Female`

## Validation

Just validate number with Luhn algorithm independent of type.

```go
parsed, err := Parse("552099-1122")
if err != nil {
	panic("could not parse")
}

if !parsed.Valid() {
	panic("this is not a valid luhn number")
}
```

Validate a social security number or coordination number. The interface supports
strings, integers and floats of the most common types.

```go
// I just care for validation
if IsValidPersonalNumber("800101-3294") {
	return RealFood()
}

if !IsValidPersonalNumber(800101329) {
	return MaybeMetal()
}

// But for these I actually care!
person, err := NewPerson("20090314-6603")
if err != nil {
	panic("what now?!")
}

if !person.Valid() {
	return NotEvenForBabies()
}

if person.Female() && person.IsOfAge(16) {
	return TimeToStartDriving()
}
```

The interface to validate organizations is the same.

```go
if !IsValidOrganisation() {
	return NoGo("556703-7485")
}
```

## Generation

In addition to validation this package also provide support to generate social
security numbers for private persons. This is great for testing purposes. There
are two interfaces, one where everything is random (age and sex) and one where
you provided it.

```go
// Just give me something!
validPerson, err := GenerateAny()
if err != nil {
	panic("oh lord")
}

// Or maybe with some preferences.
t, _ := time.Parse("2006-01-02", "1999-02-20")

girlFrom99, err := Generate(t, Female)
if err != nil {
	panic("no Spice Girl I guess?!")
}
```
