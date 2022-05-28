package L

import (
	"bufio"
	"io"
	"os"
)

// FileExists check file exists
func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// FileEmpty check file missing or has zero size
func FileEmpty(name string) bool {
	stat, err := os.Stat(name)
	return os.IsNotExist(err) || stat.Size() <= 0
}

// CreateFile create file with specific content
func CreateFile(path string, content string) bool {
	var file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if IsError(err, `CreateFile.OpenFile: %s`, path) {
		return false
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if IsError(err, `CreateFile.WriteFile: %s`, path) {
		return false
	}

	err = file.Sync()
	if IsError(err, `CreateFile.SyncFile: %s`, path) {
		return false
	}
	return true
}

// CreateDir create directory recursively
func CreateDir(path string) bool {
	err := os.MkdirAll(path, 0777)
	if IsError(err, `CreateDir: `+path) {
		return false
	}
	return true
}

// ReadFile read file content as string
func ReadFile(path string) string {
	var buff, err = os.ReadFile(path)
	if IsError(err, `ReadFile: %s`, path) {
		return ``
	}
	return string(buff)
}

// ReadFileLines read file content line by line
func ReadFileLines(path string, eachLineFunc func(line string) (exitEarly bool)) (ok bool) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if IsError(err, `ReadFileLines.OpenFile: %s`, path) {
		return false
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			return true
		}
		if eachLineFunc(line) {
			return true
		}
	}
}
