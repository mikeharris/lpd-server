package lpd

import (
	"bufio"
	"io"
)

type Payload struct {
	PrintFile       []byte
	Filename        string
	FileSizeInBytes int
}

func (p *Payload) GetFileSizeInKB() float64 {
	return float64(p.FileSizeInBytes) / 1000.0
}

func (p *Payload) unmarshal(reader io.Reader) (err error) {
	p.PrintFile, err = bufio.NewReader(reader).ReadBytes(ACK)
	if err != nil && err != io.EOF {
		return err
	}
	p.FileSizeInBytes = len(p.PrintFile) - 1
	return nil
}
