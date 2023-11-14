package conf

import (
	"os"
	"path"

	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
)

func UploadDir() string {
	// upload directory logic
	uploadDir := os.Getenv(`UPLOAD_DIR`)
	if !S.StartsWith(uploadDir, `/`) {
		workdDir, err := os.Getwd()
		L.PanicIf(err, `failed get working directory`)
		uploadDir = path.Join(workdDir, uploadDir)
	}
	dirStat, err := os.Stat(uploadDir)
	if err != nil {
		err = os.MkdirAll(uploadDir, 0770)
		L.PanicIf(err, `failed create upload directory: `+uploadDir)
		dirStat, _ = os.Stat(uploadDir)
	}
	if !dirStat.IsDir() {
		panic(`upload dir is not a directory: ` + uploadDir)
	}
	if !S.EndsWith(uploadDir, `/`) {
		uploadDir += `/`
	}
	return uploadDir
}
