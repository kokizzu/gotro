package C

// Character/Rune support package

type charOrRune interface {
	~byte | ~rune
}

// IsDigit check whether the character is a digit or not
//
//	C.IsDigit('9') // true
func IsDigit[T charOrRune](ch T) bool {
	return ch >= '0' && ch <= '9'
}

// IsAlpha check whether the character is a letter or not
//
//	C.IsDigit('a') // true
//	C.IsDigit('Z') // true
func IsAlpha[T charOrRune](ch T) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// IsIdentStart check whether the character is a valid identifier prefix (letter/underscore)
//
//	C.IsIdentStart('-') // false
//	C.IsIdentStart('_') // true
func IsIdentStart[T charOrRune](ch T) bool {
	return IsAlpha(ch) || ch == '_'
}

// IsIdent check whether the character is a valid identifier suffix alphanumeric (letter/underscore/numeral)
//
//	C.IsIdent('9'))
func IsIdent[T charOrRune](ch T) bool {
	return IsDigit(ch) || IsIdentStart(ch)
}

// IsValidFilename check whether the character is a safe file-name characters (alphanumeric/comma/full-stop/dash)
//
//	C.IsValidFilename(' ') // output bool(true)
func IsValidFilename[T charOrRune](ch T) bool {
	return ch == ' ' || IsIdent(ch) || ch == ',' || ch == '.' || ch == '-'
}
