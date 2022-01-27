package S

import (
	"testing"

	"github.com/kokizzu/gotro/L"
)

// Testing positive case function EqualsIgnoreCase //
func TestPositiveEqualsIgnoreCaseFunc(t *testing.T) {
	firstInput := "Bola"
	secondInput := "bola"

	result := EqualsIgnoreCase(firstInput, secondInput)

	if !result {
		t.Errorf(`%s %s should equal`, firstInput, secondInput)
	}
}

// Testing negative case function EqualsIgnoreCase //
func TestNegativeEqualsIgnoreCaseFunc(t *testing.T) {
	firstInput := "Bola"
	secondInput := "balon"

	result := EqualsIgnoreCase(firstInput, secondInput)

	if result {
		t.Errorf(`%s %s should not equal`, firstInput, secondInput)
	}
}

func TestMidFunction(t *testing.T) {
	text := "ABCDE"
	L.Print(`1` + Mid(text, -1, 1))
	L.Print(`2` + Mid(text, 0, -1))
	L.Print(`3` + Mid(text, 3, 1))
	L.Print(`4` + Mid(text, 4, 5))
	L.Print(`5` + Mid(text, 6, 1))
	L.Print(`6` + Mid(text, 2, 3))
}
