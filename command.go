package lpd

import (
	"strings"
)

const (
	JOB_ABORT              = 1
	JOB_RECEIVE_CONTROL    = 2
	JOB_RECEIVE_DATA       = 3
	CMD_PRINT_PENDING_JOBS = 1
	CMD_RECEIVE_JOB        = 2
	CMD_SEND_Q_STATE_SHORT = 3
	CMD_SEND_Q_STATE_LONG  = 4
)

type Command struct {
	Code     byte
	Operands []string
}

func (c *Command) unmarshal(command []byte) {
	c.Code = command[0]
	s := string(command[1:])
	s = strings.TrimRight(s, "\n")
	ops := strings.Split(s, " ")
	c.Operands = ops
}
