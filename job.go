package lpd

import (
	"bufio"
	"bytes"
	"golang.org/x/exp/slog"
	"io"
	"net"
)

type PrintJob struct {
	Payload     *Payload
	ControlFile *ControlFile
}

const (
	ACK = 0x00
	LF  = 0x0A

	STATE_IDLE = iota
	STATE_RECEIVE_JOB
	STATE_RECEIVE_DATA
	STATE_RECEIVE_CONTROL
)

func HandleConnection(c net.Conn, filePrefix string, out chan<- PrintJob) {
	slog.Debug("NEW LPD CONNECTION", "remoteAddress", c.RemoteAddr().String())
	state := STATE_IDLE
	defer c.Close()
	var job = PrintJob{Payload: new(Payload), ControlFile: new(ControlFile)}

	for {
		var com = new(Command)
		switch state {
		case STATE_RECEIVE_DATA:
			err := job.Payload.unmarshal(c)
			if err != nil {
				slog.Error("FAILED TO FETCH PRINT PAYLOAD: ", err)
			}
			state = STATE_IDLE
		case STATE_RECEIVE_CONTROL:
			data, err := bufio.NewReader(c).ReadBytes(ACK)
			if err != nil {
				slog.Error("FAILED TO READ CONTROL FILE: ", err)
			}
			job.ControlFile.unmarshal(bytes.NewReader(data))
			slog.Debug("CONTROL FILE", "content", job.ControlFile)
			state = STATE_RECEIVE_JOB
		case STATE_RECEIVE_JOB:
			line, err := bufio.NewReader(c).ReadBytes(LF)
			if err != nil {
				slog.Error("ERROR READING LPD COMMAND LINE: ", err)
			}
			com.unmarshal(line)
			if err != nil {
				slog.Error("FAILED TO UNMARSHAL COMMAND: ", err)
				continue
			}

			switch com.Code {
			case JOB_RECEIVE_CONTROL:
				slog.Debug("RECEIVING CONTROL FILE")
				state = STATE_RECEIVE_CONTROL
			case JOB_RECEIVE_DATA:
				job.Payload.Filename = filePrefix + "-" + com.Operands[1] + ".ps"
				slog.Debug("RECEIVING DATA FILE", "filename", job.Payload.Filename, "filesize", com.Operands[0])
				state = STATE_RECEIVE_DATA
			case JOB_ABORT:
				slog.Debug("ABORTING JOB")
				return
			}
		case STATE_IDLE:
			line, err := bufio.NewReader(c).ReadBytes(LF)
			if err != nil {
				if err == io.EOF {
					slog.Debug("END OF PRINT JOB DETECTED")
					out <- job
					return
				}
				slog.Error("ERROR READING LPD COMMAND LINE: ", err)
			}
			com.unmarshal(line)
			if err != nil {
				slog.Error("FAILED TO UNMARSHAL COMMAND: ", err)
				continue
			}

			switch com.Code {
			case CMD_RECEIVE_JOB:
				state = STATE_RECEIVE_JOB
			case CMD_PRINT_PENDING_JOBS:
				slog.Debug("IDLE STATE - PRINT PENDING")
			case CMD_SEND_Q_STATE_SHORT:
				slog.Debug("IDLE STATE - SEND Q STATE SHORT")
			case CMD_SEND_Q_STATE_LONG:
				slog.Debug("IDLE STATE - SEND Q STATE LONG")
			default:
				slog.Warn("IDLE STATE - UNHANDLED COMMAND", "command", com.Code)
			}
		}
		_, err := c.Write([]byte{ACK})
		if err != nil {
			slog.Error("UNABLE TO SEND COMMAND ACK: ", err)
		}

	}
}
