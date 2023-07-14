package lpd

import (
	"fmt"
	"testing"
)

// Test Byte to KB conversion
func TestCommand_Unmarshal(t *testing.T) {
	// Test correct conversion
	commands := []string{"1mycomputer another thing\n",
		"Htwo operands\n",
		"LOneOperand\n",
		"2MissingNewline"}
	c := Command{}
	c.unmarshal([]byte(commands[0]))
	if string(c.Code) != commands[0][0:1] {
		t.Fatalf("Expected command code %s but got %s.", commands[0][0:1], string(c.Code))
	}
	if len(c.Operands) != 3 {
		for i, v := range c.Operands {
			fmt.Printf("index: %d - value: %s\n", i, v)
		}
		t.Fatalf("Expected 3 operands but got %d.", len(c.Operands))
	}
}
