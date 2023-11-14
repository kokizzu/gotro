package conf

import (
	"os"

	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
)

func EnvSuperAdmins() (res M.SB) {
	emailsStr := os.Getenv(`SUPERADMIN_EMAILS`)
	emails := S.Split(emailsStr, `,`)
	res = M.SB{}
	for _, email := range emails {
		res[email] = true
	}
	return res
}
