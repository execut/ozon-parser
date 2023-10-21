package file

import (
	"os"
)

func File(filePath string) (string, error) {
	b, err := os.ReadFile(filePath)

	return string(b), err
}
