package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"MOI-Evm-Project/evm"
)

func main() {
	evm := evm.NewEVM()

	// Example bytecode: PUSH1 0x60 PUSH2 0x20 PUSH3 0x5260 PUSH32 0x0064... MSTORE MSTORE ADD MUL SDIV EXP
	bytecode := []byte{0x60, 0x01, 0x60, 0x20, 0x60, 0x52, 0x60, 0x00, 0x64, /* ... */ 0x05, 0x60, 0x40, 0x52}
	for _, op := range bytecode {
		switch op {
		case 0x60:
			evm.Push1(uint64(bytecode[1]))
			bytecode = bytecode[2:]
		case 0x61:
			evm.Push2(uint64(bytecode[1])<<8 | uint64(bytecode[2]))
			bytecode = bytecode[3:]
		case 0x62:
			evm.Push3(uint64(bytecode[1])<<16 | uint64(bytecode[2])<<8 | uint64(bytecode[3]))
			bytecode = bytecode[4:]
		case 0x7f:
			evm.Push32(bytecode[1:33])
			bytecode = bytecode[33:]
		case 0x52:
			evm.Mstore()
			bytecode = bytecode[1:]
		case 0x53:
			evm.Mstore8()
			bytecode = bytecode[1:]
		case 0x01:
			evm.Add()
			bytecode = bytecode[1:]
		case 0x02:
			evm.Mul()
			bytecode = bytecode[1:]
		case 0x04:
			evm.Sdiv()
			bytecode = bytecode[1:]
		case 0x05:
			evm.Exp()
			bytecode = bytecode[1:]
		default:
			fmt.Println("Unknown opcode")
			return
		}
	}

	fmt.Println("Final Stack:", evm.GetStack())
	fmt.Println("Gas Consumed:", evm.GetGas())

	// Get the absolute path of the current directory
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// Set the output file name (e.g., evm-program)
	outputName := "evm-program"

	// Build the program
	cmd := exec.Command("go", "build", "-o", outputName, dir)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Build successful. Executable file created:", outputName)
}
