package W2

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kokizzu/gotro/S"
)

// loads .env file even when the binary/test not in project's root directory
// returns project's root directory (where `.env` should be located)
func LoadTestEnv() string {
	for z := 0; z < 4; z++ {
		dir := strings.Repeat(`../`, z)
		err := godotenv.Load(dir + `.env`)
		if err == nil {
			cwd, _ := os.Getwd()
			for i := 0; i < z; i++ {
				cwd = S.LeftOfLast(cwd, "/")
			}
			return cwd
		}
	}
	return ``
}
