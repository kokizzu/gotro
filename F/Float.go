package F

// Floating-point/Real number support package

import (
	"strconv"
	"time"
)

// If simplified ternary operator (bool ? val : 0), returns second argument, if the condition (first arg) is true, returns 0 if not
//
//	F.If(true,3.12) // 3.12
//	F.If(false,3)   // 0
func If(b bool, yes float64) float64 {
	if b {
		return yes
	}
	return 0
}

// IfElse ternary operator (bool ? val1 : val2), returns second argument if the condition (first arg) is true, third argument if not
//
//	F.IfElse(true,3.12,3.45))  // 3.12
func IfElse(b bool, yes, no float64) float64 {
	if b {
		return yes
	}
	return no
}

// ToS convert float64 to string
//
//	F.ToS(3.1284)) // `3.1284`
func ToS(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}

// ToStr convert float64 to string with 2 digits behind the decimal point
//
//	F.ToStr(3.1284)) // `3.13`
func ToStr(num float64) string {
	return strconv.FormatFloat(num, 'f', 2, 64)
}

// ToIsoDateStr convert to ISO-8601 string
//
//	F.ToIsoDateStr(0) // `1970-01-01T00:00:00`
func ToIsoDateStr(num float64) string {
	n := int64(num)
	t := time.Unix(n, 0)
	return t.UTC().Format(`2006-01-02T15:04:05Z`)
}

// ToDateStr convert float64 unix to `YYYY-MM-DD`
func ToDateStr(num float64) string {
	n := int64(num)
	t := time.Unix(n, 0)
	return t.UTC().Format(`2006-01-02`)
}
