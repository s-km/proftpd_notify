package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	bytesPerLine     = 147
	dateFormat       = "Mon Jan 02 15:04:05 2006"
	dateOutputFormat = "Monday, January 2nd 2006 (03:04:05PM)"
)

func HandleErr(msg string, err error) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

func GetLatestXfer(path string) string {
	file, err := os.Open(path)
	HandleErr("Unable to open log file:", err)
	defer file.Close()

	buf := make([]byte, bytesPerLine)
	file.Seek(-bytesPerLine, 2)
	_, err = file.Read(buf)
	HandleErr("Unable to read file:", err)

	return string(buf)
}

func Expand(path string) string {
	if path[0] == '~' || strings.HasPrefix(path, "$HOME") {
		user, err := user.Current()
		HandleErr("Failed to get current user:", err)

		rest := strings.TrimLeftFunc(path, func(r rune) bool {
			return r != os.PathSeparator
		})

		return filepath.Join(user.HomeDir, rest)
	}

	return path
}
