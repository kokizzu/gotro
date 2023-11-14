package wcAuth

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/kokizzu/gotro/L"
)

func (p *UsersMutator) SetEncryptPassword(password string) bool {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	p.SetPassword(string(pass))
	return !L.IsError(err, `bcrypt.GenerateFromPassword`)
}

// add more custom methods here
