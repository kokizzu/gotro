package S

import (
	"testing"
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
		t.Error(`%s %s should not equal`, firstInput, secondInput)
	}
}
