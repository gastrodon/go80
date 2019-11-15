package disassembler

import (
	"fmt"
	"os"
	"strings"
)

var instructions map[byte]string = map[byte]string{
	// nothing
    0x00: "NOP",		0x10: "NOP",		0x08: "NOP",		0x18: "NOP",
    0x20: "NOP",		0x30: "NOP",		0x28: "NOP",		0x38: "NOP",
	// decimal adjust	halt
	0x27: "DAA",		0x76: "HLT",
	// databus out		in
	0xd3: "OUT %", 		0xdb: "IN %",
	// interrupts
	// disable			enable
	0xf3: "DI",			0xfb: "EI",
	// shift register A
	// left				right
	// use bit wrapping
	0x07: "RLC",		0x0f: "RRC",
	// use carry
	0x17: "RAL",		0x1f: "RAR",
	// A = !A
	0x2f: "CMA",
	// carry = !carry	carry = 1
	0x3f: "CMC",		0x37: "STC",
	// load value into register pair
	0x01: "LXI B, %%",	0x11: "LXI D, %%",
	0x21: "LXI H, %%",	0x31: "LXI SP, %%",
	// load register pair
	0x0a: "LDAX B",		0x1a: "LDAX D",
	// load address		store address
	0x3a: "LDA $%%",	0x32: "STA $%%",
	0x2a: "LHLD $%%",	0x22: "SHLD $%%",
	// Increment		Decrement			Store value
	// 8 bit
	0x04: "INR B",   	0x05: "DCR B",		0x06: "MVI B, %",
	0x0c: "INR C",		0x0d: "DCR C",		0x0e: "MVI C, %",
	0x14: "INR D",   	0x15: "DCR D",		0x16: "MVI D, %",
	0x1c: "INR E",		0x1d: "DCR E",		0x1e: "MVI E, %",
	0x24: "INR H",   	0x25: "DCR H",		0x26: "MVI H, %",
	0x2c: "INR L",		0x2d: "DCR L",		0x2e: "MVI L, %",
	0x34: "INR M",   	0x35: "DCR M",		0x36: "MVI M, %",
	0x3c: "INR A",		0x3d: "DCR A",		0x3e: "MVI A, %",
	// 16 bit register pairs
	0x03: "INX B",		0x0b: "DCX B",		0x02: "STAX B",
	0x13: "INX D",		0x1b: "DCX D",		0x12: "STAX D",
	0x23: "INX H",		0x2b: "DCX H",
	0x33: "INX SP",		0x3b: "DCX SP",
	// add 16 bit register pair to HL
	0x09: "DAD B",		0x19: "DAD D",
	0x29: "DAD H", 		0x39: "DAD SP",
	// Register B
	0x40: "MOV B, B", 	0x41: "MOV B, C",	0x42: "MOV B, D",
	0x43: "MOV B, E", 	0x44: "MOV B, H", 	0x45: "MOV B, L",
	0x46: "MOV B, M",	0x47: "MOV B, A",
	// Register C
	0x48: "MOV C, B", 	0x49: "MOV C, C",	0x4a: "MOV C, D",
	0x4b: "MOV C, E", 	0x4c: "MOV C, H", 	0x4d: "MOV C, L",
	0x4e: "MOV C, M",	0x4f: "MOV C, A",
	// Register D
	0x50: "MOV D, B", 	0x51: "MOV D, C",	0x52: "MOV D, D",
	0x53: "MOV D, E", 	0x54: "MOV D, H", 	0x55: "MOV D, L",
	0x56: "MOV D, M",	0x57: "MOV D, A",
	// Register E
	0x58: "MOV E, B", 	0x59: "MOV E, C",	0x5a: "MOV E, D",
	0x5b: "MOV E, E", 	0x5c: "MOV E, H", 	0x5d: "MOV E, L",
	0x5e: "MOV E, M",	0x5f: "MOV E, A",
	// Register H
	0x60: "MOV H, B", 	0x61: "MOV H, C",	0x62: "MOV H, D",
	0x63: "MOV H, E", 	0x64: "MOV H, H", 	0x65: "MOV H, L",
	0x66: "MOV H, M",	0x67: "MOV H, A",
	// Register L
	0x68: "MOV L, B", 	0x69: "MOV L, C",	0x6a: "MOV L, D",
	0x6b: "MOV L, E", 	0x6c: "MOV L, H", 	0x6d: "MOV L, L",
	0x6e: "MOV L, M",	0x6f: "MOV L, A",
	// Register M
	0x70: "MOV M, B", 	0x71: "MOV M, C",	0x72: "MOV M, D",
	0x73: "MOV M, E", 	0x74: "MOV M, H", 	0x75: "MOV M, L",
						0x77: "MOV M, A",
	// Register A
	0x78: "MOV A, B", 	0x79: "MOV A, C",	0x7a: "MOV A, D",
	0x7b: "MOV A, E", 	0x7c: "MOV A, H", 	0x7d: "MOV A, L",
	0x7e: "MOV A, M",	0x7f: "MOV A, A",
	// register A add
	0x80: "ADD B", 		0x81: "ADD C",		0x82: "ADD D",		0x83: "ADD E",
	0x84: "ADD H",		0x85: "ADD L",		0x86: "ADD M", 		0x87: "ADD A",
	0xc6: "ADI %",		// immediate
	// register A add with carry
	0x88: "ADC B", 		0x89: "ADC C",		0x8a: "ADC D",		0x8b: "ADC E",
	0x8c: "ADC H",		0x8d: "ADC L",		0x8e: "ADC M", 		0x8f: "ADC A",
	0xce: "ACI %",		// immediate
	// register A subtract
	0x90: "SUB B", 		0x91: "SUB C",		0x92: "SUB D",		0x93: "SUB E",
	0x94: "SUB H",		0x95: "SUB L",		0x96: "SUB M", 		0x97: "SUB A",
	0xd6: "SUI",		// immediate
	// register A subtract with carry
	0x98: "SBB B", 		0x99: "SBB C",		0x9a: "SBB D",		0x9b: "SBB E",
	0x9c: "SBB H",		0x9d: "SBB L",		0x9e: "SBB M", 		0x9f: "SBB A",
	0xde: "SBI %",		// immediate
	// register A AND
	0xa0: "ANA B", 		0xa1: "ANA C",		0xa2: "ANA D",		0xa3: "ANA E",
	0xa4: "ANA H",		0xa5: "ANA L",		0xa6: "ANA M", 		0xa7: "ANA A",
	0xe6: "ANI %",		// immediate
	// register A XOR
	0xa8: "XRA B", 		0xa9: "XRA C",		0xaa: "XRA D",		0xab: "XRA E",
	0xac: "XRA H",		0xad: "XRA L",		0xae: "XRA M", 		0xaf: "XRA A",
	0xee: "XRI %", 		// immediate
	// register A OR
	0xb0: "ORA B", 		0xb1: "ORA C",		0xb2: "ORA D",		0xb3: "ORA E",
	0xb4: "ORA H",		0xb5: "ORA L",		0xb6: "ORA M", 		0xb7: "ORA A",
	0xf6: "ORI %",		// immediate
	// register A comparison
	0xb8: "CMP B", 		0xb9: "CMP C",		0xba: "CMP D",		0xbb: "CMP E",
	0xbc: "CMP H",		0xbd: "CMP L",		0xbe: "CMP M", 		0xbf: "CMP A",
	0xfe: "CPI %",		// immediate
	// POP register pair
	0xc1: "POP B", 		0xd1: "POP D",		0xe1: "POP H",		0xf1: "POP PSW",
	// PUSH register pair
	0xc5: "PUSH B",		0xd5: "PUSH D",		0xe5: "PUSH H",		0xf5: "PUSH PSW",
	// HL <-> STACK[SP]	SP <- HL			H <-> D, L <-> E
	0xe3: "XTHL",		0xf9: "SPHL",		0xeb: "XCHG",
	// Restart after interrupt
	0xc7: "RST 0",		0xcf: "RST 1",		0xd7: "RST 2",		0xdf: "RST 3",
}

func bytes_of(path string) ([]byte, int64, error) {
	var stat os.FileInfo
	var err error
	stat, err = os.Stat(path)

	if err != nil {
		return make([]byte, 0), 0, err
	}

	var size int64
	size = stat.Size()

	var bytes []byte = make([]byte, size)

	var file *os.File
	file, err = os.Open(path)

	if err != nil {
		return bytes, size, err
	}

	defer file.Close()
	file.Read(bytes)

	return bytes, size, nil
}

func disassemble_bytes(bytes []byte, size int64) ([]string, error) {
	var index int64 = 0
	var argc int64
	var instruction string
	var args []byte

	for index < size {
		instruction = instructions[bytes[index]]
		argc = int64(strings.Count(instruction, "%"))
		args = bytes[index+1 : index+1+argc]
		instruction = strings.ReplaceAll(instruction, "%", "")
		index += argc + 1

		for argc > 0 {
			argc--
			instruction += fmt.Sprintf("%02X", args[argc])
		}

		fmt.Println(instruction)
	}

	return []string{}, nil
}

func push(array *[]string, item string) {
	var new []string = append(*array, item)
	array = &new
}

func T() {
	var bytes []byte
	var size int64
	bytes, size, _ = bytes_of("./source/invaders.h")
	disassemble_bytes(bytes, size)
}
