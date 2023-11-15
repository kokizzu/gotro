package conf

import (
	"fmt"
)

func WrapError(parent, child error) error {
	return fmt.Errorf("%w: %v", parent, child)
}
