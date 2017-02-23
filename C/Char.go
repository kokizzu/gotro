package C

// Character/Rune support package

// check whether the character is a digit or not
//  C.IsDigit('9') // true
func IsDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// check whether the character is a valid identifier prefix (letter/underscore)
//  C.IsIdentStart('-') // false
//  C.IsIdentStart('_') // true
func IsIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

// check whether the character is a valid identifier suffix alphanumeric (letter/underscore/numeral)
//  C.IsIdent('9'))
func IsIdent(ch byte) bool {
	return IsDigit(ch) || IsIdentStart(ch)
}

// check whether the character is a safe file-name characters (alphanumeric/comma/full-stop/dash)
//   C.IsValidFilename(' ') // output bool(true)
func IsValidFilename(ch byte) bool {
	return ch == ' ' || IsIdent(ch) || ch == ',' || ch == '.' || ch == '-'
}
