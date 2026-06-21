package gosysinfo

import (
	"os"
)

func getFileContent(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		// TODO: handle this error
	}

	return string(content)
}
