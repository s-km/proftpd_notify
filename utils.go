package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
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
	HandleErr("Unable to open log file: ", err)
	defer file.Close()

	stat, _ := file.Stat()
	size := stat.Size()

	var buf []byte
	cursor := int64(0)
	for cursor > -size {
		cursor--

		file.Seek(cursor, io.SeekEnd)
		b := make([]byte, 1)
		file.Read(b)

		if cursor != -1 && b[0] == 10 || b[0] == 13 {
			break
		}

		buf = append(b, buf...)
	}

	return string(buf)
}

func Expand(path string) string {
	if path[0] == '~' || strings.HasPrefix(path, "$HOME") {
		user, err := user.Current()
		HandleErr("Failed to get current user: ", err)

		rest := strings.TrimLeftFunc(path, func(r rune) bool {
			return r != os.PathSeparator
		})

		return filepath.Join(user.HomeDir, rest)
	}

	return path
}

func Hash(path string) string {
	file, err := os.Open(path)
	HandleErr("Failed to open file: ", err)
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		HandleErr("Failed to create hash: ", err)
	}

	return hex.EncodeToString(hash.Sum(nil))
}
