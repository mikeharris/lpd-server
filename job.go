package lpd

import (
	"bufio"
	"bytes"
	"io"
	"log"
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

func HandlePrintJob(c net.Conn, filePrefix string, out chan<- string) {
	log.Println("NEW LPD CONNECTION FROM ", c.RemoteAddr().String())
	state := STATE_IDLE
	defer c.Close()
	var printFilename string
	var job *PrintJob = new(PrintJob)

	for {
		var com = new(Command)
		switch state {
		case STATE_RECEIVE_DATA:
			err := job.Payload.unmarshal(c)
			if err != nil {
				log.Println("FAILED TO FETCH PRINT PAYLOAD: ", err)
			}
			// err = os.WriteFile(printFilename, payload, 0644)
			if err != nil {
				log.Println("FAILED TO PERSIST PRINT FILE: ", err)
			}
			state = STATE_IDLE
		case STATE_RECEIVE_CONTROL:
			data, err := bufio.NewReader(c).ReadBytes(ACK)
			if err != nil {
				log.Println("FAILED TO READ CONTROL FILE: ", err)
			}
			job.ControlFile.unmarshal(bytes.NewReader(data))
			log.Printf("CONTROL FILE CONTENT: %+v \n", job.ControlFile)
			state = STATE_RECEIVE_JOB
		case STATE_RECEIVE_JOB:
			line, err := bufio.NewReader(c).ReadBytes(LF)
			if err != nil {
				log.Println("ERROR READING LPD COMMAND LINE: ", err)
			}
			com.unmarshal(line)
			log.Println(com)
			if err != nil {
				log.Println("FAILED TO UNMARSHAL COMMAND: ", err)
				continue
			}

			switch com.Code {
			case JOB_RECEIVE_CONTROL:
				log.Println("RECEIVING CONTROL FILE")
				state = STATE_RECEIVE_CONTROL
			case JOB_RECEIVE_DATA:
				log.Println("RECEIVING DATA FILE")
				log.Println("FILENAME: ", printFilename)
				log.Println("FILESIZE (bytes): ", com.Operands[0])
				job.Payload.Filename = filePrefix + "-" + com.Operands[1] + ".ps"
				state = STATE_RECEIVE_DATA
			case JOB_ABORT:
				log.Println("ABORTING JOB")
				return
			}
		case STATE_IDLE:
			line, err := bufio.NewReader(c).ReadBytes(LF)
			if err != nil {
				if err == io.EOF {
					log.Println("END OF PRINT JOB DETECTED")
					out <- printFilename
					return
				}
				log.Println("ERROR READING LPD COMMAND LINE: ", err)
			}
			com.unmarshal(line)
			log.Println(com)
			if err != nil {
				log.Println("FAILED TO UNMARSHAL COMMAND: ", err)
				continue
			}

			switch com.Code {
			case CMD_RECEIVE_JOB:
				state = STATE_RECEIVE_JOB
			case CMD_PRINT_PENDING_JOBS:
				log.Println("IDLE STATE - PRINT PENDING")
			case CMD_SEND_Q_STATE_SHORT:
				log.Println("IDLE STATE - SEND Q STATE SHORT")
			case CMD_SEND_Q_STATE_LONG:
				log.Println("IDLE STATE - SEND Q STATE LONG")
			default:
				log.Println("IDLE STATE - UNHANDLED COMMAND: ", com.Code)
			}
		}
		_, err := c.Write([]byte{ACK})
		if err != nil {
			log.Println("UNABLE TO SEND COMMAND ACK: ", err)
		}

	}
}
