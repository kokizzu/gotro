package L

import "os"

func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func CreateFile(path string, content string) bool {
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		Print(`CreateFile.OpenFile: ` + path + ` | ` + err.Error())
		return false
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		Print(`CreateFile.WriteFile: ` + path + ` | ` + err.Error())
		return false
	}

	err = file.Sync()
	if err != nil {
		Print(`CreateFile.SyncFile: ` + path + ` | ` + err.Error())
		return false
	}
	return true
}
