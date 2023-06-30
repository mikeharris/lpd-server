package lpd

import (
	"bufio"
	"io"
)

type ControlFile struct {
	BannerName     string
	HostName       string
	JobName        string
	SourceFilename string
	User           string
	Title          string
}

func (c *ControlFile) unmarshal(controlFile io.Reader) {
	scanner := bufio.NewScanner(controlFile)
	for scanner.Scan() {
		line := scanner.Text()
		code := line[0:1]
		value := line[1:]
		switch code {
		case "C":
			c.BannerName = value
		case "H":
			c.HostName = value
		case "J":
			c.JobName = value
		case "N":
			c.SourceFilename = value
		case "P":
			c.User = value
		case "T":
			c.Title = value

		}
	}
}
