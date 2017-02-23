package B

// Boolean support package

// converts boolean type to string type, writing "true" or "false"
//   B.ToS(2 > 1)  // "true"
func ToS(b bool) string {
	if b {
		return `true`
	}
	return `false`
}
