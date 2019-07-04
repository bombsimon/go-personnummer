package swessn

import (
	"errors"
)

// County represents the counties within Sweden. This could be told from the
// serial number before 1990. See
// https://en.wikipedia.org/wiki/Personal_identity_number_(Sweden)#Format
// The naming and values of the counties is an ISO 3166-2 standard. See
// https://en.wikipedia.org/wiki/Counties_of_Sweden#Map
type County int

const (
	CountyA County = iota
	CountyAB
	CountyB
	CountyC
	CountyD
	CountyE
	CountyF
	CountyG
	CountyH
	CountyI
	CountyK
	CountyL
	CountyM
	CountyN
	CountyO
	CountyP
	CountyR
	CountyS
	CountyT
	CountyU
	CountyW
	CountyX
	CountyY
	CountyZ
	CountyAC
	CountyBD
	CountyQ
	CountyQQ
)

// String returns the name of the region where the person was born, if born
// before 1990 when this system was removed.
func (c County) String() string {
	switch c {
	case CountyA, CountyAB, CountyB:
		return "Stockholms Län"
	case CountyC:
		return "Uppsala län"
	case CountyD:
		return "Södermanlands län"
	case CountyE:
		return "Östergötlands län"
	case CountyF:
		return "Jönköpings län"
	case CountyG:
		return "Kronobergs län"
	case CountyH:
		return "Kalmar län"
	case CountyI:
		return "Gotlands län"
	case CountyK:
		return "Blekinge län"
	case CountyL:
		return "Kristianstads län"
	case CountyM:
		return "Malmöhus län"
	case CountyN:
		return "Hallands län"
	case CountyO:
		return "Göteborgs och Bohus län"
	case CountyP:
		return "Älvsborgs län"
	case CountyR:
		return "Skaraborgs län"
	case CountyS:
		return "Värmlands län"
	case CountyQ:
		return "Födda utomlands"
	case CountyT:
		return "Örebro län"
	case CountyU:
		return "Västmanlands län"
	case CountyW:
		return "Kopparbergs län"
	case CountyX:
		return "Gävleborgs län"
	case CountyY:
		return "Västernorrlands län"
	case CountyZ:
		return "Jämtlands län"
	case CountyAC:
		return "Västerbottens län"
	case CountyBD:
		return "Norrbottens län"
	case CountyQQ:
		return "Outside Sweden or non Swedish citizen"
	}

	return ""
}

// CountyFromSerial will calculate the appropriate county based on a serial
// number. The source for these values may be found here:
// https://sv.wikipedia.org/wiki/Personnummer_i_Sverige#F%C3%B6delsenumret
func CountyFromSerial(serial int) (County, error) {
	switch s := serial; {
	case s < 139:
		return CountyA, nil
	case s < 159:
		return CountyC, nil
	case s < 189:
		return CountyD, nil
	case s < 239:
		return CountyE, nil
	case s < 269:
		return CountyF, nil
	case s < 289:
		return CountyG, nil
	case s < 319:
		return CountyH, nil
	case s < 329:
		return CountyI, nil
	case s < 349:
		return CountyK, nil
	case s < 389:
		return CountyL, nil
	case s < 459:
		return CountyM, nil
	case s < 479:
		return CountyN, nil
	case s < 549:
		return CountyO, nil
	case s < 589:
		return CountyP, nil
	case s < 619:
		return CountyR, nil
	case s < 649:
		return CountyS, nil
	case s < 659:
		return CountyQ, nil
	case s < 689:
		return CountyT, nil
	case s < 709:
		return CountyU, nil
	case s < 739:
		return CountyW, nil
	case s < 779:
		return CountyX, nil
	case s < 819:
		return CountyY, nil
	case s < 849:
		return CountyZ, nil
	case s < 889:
		return CountyAC, nil
	case s < 929:
		return CountyBD, nil
	case s < 999:
		return CountyQQ, nil
	}

	return County(-1), errors.New("invalid serial")
}
