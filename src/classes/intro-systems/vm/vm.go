package vm

import (
	"encoding/hex"
	"fmt"
	"os"
)

const (
	Load  = 0x01
	Store = 0x02
	Add   = 0x03
	Sub   = 0x04
	Halt  = 0xff
)

// Stretch goals
const (
	Addi = 0x05
	Subi = 0x06
	Jump = 0x07
	Beqz = 0x08
)

// Given a 256 byte array of "memory", run the stored program
// to completion, modifying the data in place to reflect the result
//
// The memory format is:
//
// 00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f ... ff
// __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ ... __
// ^==DATA===============^ ^==INSTRUCTIONS==============^
//
func compute(memory []byte) {
	registers := [3]byte{8, 0, 0} // PC, R1 and R2

	// Keep looping, like a physical computer's clock
	for {
		pc := registers[0]

		op := memory[pc]

		// // decode and execute
		switch op {
		case Load:
			registers[0] = pc + 3

			from := memory[pc+1]
			to := memory[pc+2]

			registers[to] = memory[from]

		case Store:
			registers[0] = pc + 3

			from := memory[pc+1]
			to := memory[pc+2]

			memory[to] = registers[from]

		case Add:
			registers[0] = pc + 3

			r1 := memory[pc+1]
			r2 := memory[pc+2]

			registers[r1] += registers[r2]

		case Addi:
			registers[0] = pc + 3

			r1 := memory[pc+1]
			n := memory[pc+2]

			registers[r1] += n

		case Sub:
			registers[0] = pc + 3

			r1 := memory[pc+1]
			r2 := memory[pc+2]

			registers[r1] -= registers[r2]

		case Subi:
			registers[0] = pc + 3

			r1 := memory[pc+1]
			n := memory[pc+2]

			registers[r1] -= n

		case Jump:
			next := memory[pc+1]

			registers[0] = next

		case Beqz:
			registers[0] = pc + 3

			r := memory[pc+1]
			value := registers[r]

			if value == 0 {
				offset := memory[pc+2]
				registers[0] += offset
			}

		case Halt:
			return

		default:
			fmt.Fprintf(os.Stderr, "unknown instruction 0x%s at memory location %d - halting\n", hex.EncodeToString([]byte{op}), pc)
			return
		}
	}
}
