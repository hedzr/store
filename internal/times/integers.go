package times

import (
	"strconv"
)

func SmartParseInt(s string) (ret int64, err error) {
	ret, err = strconv.ParseInt(s, 0, 64)
	return
}

func MustSmartParseInt(s string) (ret int64) {
	ret, _ = strconv.ParseInt(s, 0, 64)
	return
}

func SmartParseUint(s string) (ret uint64, err error) {
	ret, err = strconv.ParseUint(s, 0, 64)
	return
}

func MstSmartParseUint(s string) (ret uint64) {
	ret, _ = strconv.ParseUint(s, 0, 64)
	return
}

func ParseFloat(s string) (ret float64, err error) {
	ret, err = strconv.ParseFloat(s, 64)
	return
}

func MustParseFloat(s string) (ret float64) {
	ret, _ = strconv.ParseFloat(s, 64)
	return
}

func ParseComplex(s string) (ret complex128, err error) {
	ret, err = strconv.ParseComplex(s, 64)
	return
}

func MustParseComplex(s string) (ret complex128) {
	ret, _ = strconv.ParseComplex(s, 64)
	return
}
