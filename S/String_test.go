package S

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testing positive case function EqualsIgnoreCase //
func TestPositiveEqualsIgnoreCaseFunc(t *testing.T) {
	firstInput := "Bola"
	secondInput := "bola"

	result := EqualsIgnoreCase(firstInput, secondInput)

	assert.Equal(t, true, result, "Result should return true, but actual get "+result)
}

// Testing negative case function EqualsIgnoreCase //
func TestNegativeEqualsIgnoreCaseFunc(t *testing.T) {
	firstInput := "Bola"
	secondInput := "balon"

	result := EqualsIgnoreCase(firstInput, secondInput)

	assert.Equal(t, false, result, "Result should return false, but actual get "+result)
}
