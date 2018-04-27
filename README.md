# Swedish SSN

[![Build Status](https://travis-ci.org/bombsimon/swedish-ssn.svg?branch=master)](https://travis-ci.org/bombsimon/swedish-ssn)
[![GoDoc](https://godoc.org/github.com/bombsimon/swedish-ssn?status.svg)](https://godoc.org/github.com/bombsimon/swedish-ssn)

This package aims to provide a toolset to handle Swedish social security numbers. A Swedish social security number uses the [Luhn algorithm](https://en.wikipedia.org/wiki/Luhn_algorithm) to calculate the checksum and create a valid number.

## Example

**Validate social security number**
```go
reader := bufio.NewReader(os.Stdin)
fmt.Println("Enter social security number: ")
ssnInput, _ := reader.ReadString('\n')

ssn := swessn.New(ssnInput)

if ssn.Valid() {
    fmt.Println("Valid SSN provided: %s", ssn)
}

minRequiredAge := 18
if ssn.IsOfAge(minRequiredAge) {
    fmt.Println("Person is %d which is above %d", ssn.Age(), minRequiredAge)
}
```

**Generate social security number**
```go
ssn1 := swessn.Generate(1990, 10, 10, "m")
ssn2 := swessn.Generate(1990, 10, 10, "f")
ssn3 := swessn.GenerateAny()

fmt.Printf("Male: %s, Female: %s, Random: %s\n", ssn1, ssn2, ssn3)
```
