package main

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/c2h5oh/datasize"
)

type Transfer struct {
	Date       string
	Duration   string
	RemoteHost string
	FileSize   string
	Filename   string
	Direction  string
	User       string
	Status     string
}

func parseDate(s string) string {
	loc, err := time.LoadLocation("America/Toronto")
	HandleErr("Failed to load timezone:", err)

	t, err := time.Parse(dateFormat, s)
	HandleErr("Failed to parse date:", err)

	return t.In(loc).Format(dateOutputFormat)
}

func parseFilesize(s string) string {
	var bs datasize.ByteSize

	err := bs.UnmarshalText([]byte(s))
	HandleErr("Failed to parse filesize string:", err)

	return bs.HumanReadable()
}

func parseDuration(s string) string {
	u, err := strconv.ParseUint(s, 10, 32)
	HandleErr("Failed to convert duration string to uint:", err)

	dur := time.Duration(u) * time.Second
	return dur.String()
}

func parseFilename(s string) string {
	return filepath.Base(s)
}

func parseDirection(s string) string {
	switch {
	case s == "o":
		return "Downloaded"
	case s == "i":
		return "Uploaded"
	case s == "d":
		return "Deleted"
	default:
		return "Unknown"
	}
}

func parseStatus(s string) string {
	switch {
	case s == "c":
		return "Complete"
	case s == "i":
		return "Incomplete"
	default:
		return "Unknown"
	}
}

func getDate(frags []string) string {
	return strings.Join([]string{frags[0], frags[1], frags[2], frags[3], frags[4]}, " ")
}

// See http://www.proftpd.org/docs/howto/Logging.html for format
func ParseXferLogEntry(t string) Transfer {
	frags := strings.Split(t, " ")

	xfer := Transfer{
		RemoteHost: frags[6],
		User:       frags[13],
	}

	xfer.Date = parseDate(getDate(frags))
	xfer.Duration = parseDuration(frags[5])
	xfer.FileSize = parseFilesize(frags[7])
	xfer.Filename = parseFilename(frags[8])
	xfer.Direction = parseDirection(frags[11])
	xfer.Status = parseStatus(frags[17])

	return xfer
}
