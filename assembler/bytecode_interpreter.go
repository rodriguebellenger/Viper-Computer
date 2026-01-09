package main

import (
	"fmt"
	"log"
)

const RAMSize uint32 = 1024

var RAM [RAMSize]uint8
var registers []uint64 = []uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, uint64(RAMSize - 1), uint64(RAMSize - 1)}

const (
	HLT int = iota
	AND
	ANDIB
	ANDIW
	OR
	ORIB
	ORIW
	NOT
	SHIL
	SHILI
	SHIR
	SHIRI
	ADD
	ADDIB
	ADDIW
	INCR
	DECR
	MUL
	MULIB
	MULIW
	DIV
	DIVIB
	DIVIW
	MOD
	MODIB
	MODIW
	CLEAR
	MOV1B
	MOV2B
	MOV3B
	MOV4B
	MOV1W
	MOV2W
	MOV3W
	MOV4W
	MOVR
	SWAP
	PUSH
	PUSHIB
	PUSHIW
	PUSHIT
	POP
	PEEK
	CMP
	JMP
	JMPB
	JMPW
	JMPT
	CALL
	CALLB
	CALLW
	CALLT
	RET
	WRT
	READ
)

/////////////////////////
// Execute the program //
/////////////////////////

func executeProgram() {
	var stackUpperBound uint32 = uint32(RAMSize - (RAMSize >> 2) - 1)
	var stackLowerBound uint32 = uint32(RAMSize - 1)
loop:
	for i := uint32(0); i < RAMSize; i++ {
		//var debugVariable uint32 = i
		switch RAM[i] {
		case uint8(HLT):
			break loop
		case uint8(AND):
			i += 1
			var arg1 uint8 = RAM[i]
			i += 1
			var arg2 uint8 = RAM[i]
			registers[arg1] = registers[arg1] & registers[arg2]
			i += 1
		case uint8(ANDIB):
			i += 1
			var arg1 uint8 = RAM[i]
			i += 1
			var arg2 uint8 = RAM[i]
			registers[arg1] = registers[arg1] & uint64(arg2) // TO TEST
			i += 1
		case uint8(ANDIW):
			i += 1
			var arg1 uint8 = RAM[i]
			i += 1
			var arg2 uint8 = RAM[i]
			i += 1
			var arg3 uint8 = RAM[i]
			registers[arg1] = registers[arg1] & (uint64(arg2) | (uint64(arg3) << 8))
		case uint8(OR):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint8 = RAM[i+2]
			registers[arg1] = registers[arg1] | registers[arg2]
			i += 3
		case uint8(ORIB):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint8 = RAM[i+2]
			registers[arg1] = registers[arg1] | uint64(arg2)
			i += 3
		case uint8(ORIW):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint8 = RAM[i+2]
			var arg3 uint8 = RAM[i+3]
			registers[arg1] = registers[arg1] | (uint64(arg2) | (uint64(arg3) << 8))
			i += 3
		case uint8(NOT):
			var arg uint8 = RAM[i+1]
			registers[arg] = ^registers[arg]
			i += 3
		case uint8(SHIL):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint8 = RAM[i+2]
			registers[arg1] = registers[arg1] << registers[arg2]
			i += 3
		case uint8(SHILI):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint8 = RAM[i+2]
			registers[arg1] = registers[arg1] >> arg2
			i += 3
		case uint8(SHIR):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint8 = RAM[i+2]
			registers[arg1] = registers[arg1] << registers[arg2]
			i += 3
		case uint8(SHIRI):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint8 = RAM[i+2]
			registers[arg1] = registers[arg1] >> arg2
			i += 3
		case uint8(ADD):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint8 = RAM[i+2]
			registers[arg1] = (registers[arg1] + registers[arg2])
			i += 3
		case uint8(ADDIB):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint64 = uint64(RAM[i+2])
			if arg2>>7 == 1 {
				registers[arg1] += arg2 | 0xFFFFFFFFFFFFFF00
			} else {
				registers[arg1] += arg2
			}
			i += 3
		case uint8(ADDIW):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint64 = uint64(RAM[i+2]) | uint64(RAM[i+3])<<8
			if arg2>>15 == 1 {
				registers[arg1] += arg2 | 0xFFFFFFFFFFFF0000
			} else {
				registers[arg1] += arg2
			}
			i += 3
		case uint8(INCR):
			registers[RAM[i+1]] += 1
			i += 3
		case uint8(DECR):
			registers[RAM[i+1]] -= 1
			i += 3
		// case MUL
		// case MULIB
		// case MULIW
		// case DIV
		// case DIVIB
		// case DIVIW
		// case MOD
		// case MODIB
		// case MODIW
		// case CLEAR
		case uint8(MOV1B):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint64 = uint64(RAM[i+2])
			registers[arg1] &= 0xFFFFFFFFFFFFFF00
			registers[arg1] |= arg2
			i += 3
		case uint8(MOV2B):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint64 = uint64(RAM[i+2])
			registers[arg1] &= 0xFFFFFFFFFF00FFFF
			registers[arg1] |= (arg2 << 16)
			i += 3
		case uint8(MOV3B):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint64 = uint64(RAM[i+2])
			registers[arg1] &= 0xFFFFFF00FFFFFFFF
			registers[arg1] |= (arg2 << 32)
			i += 3
		case uint8(MOV4B):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint64 = uint64(RAM[i+2])
			registers[arg1] &= 0xFF00FFFFFFFFFFFF
			registers[arg1] |= (arg2 << 48)
			i += 3
		case uint8(MOV1W):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint64 = uint64(RAM[i+2]) | (uint64(RAM[i+3]) << 8)
			registers[arg1] &= 0xFFFFFFFFFFFF0000
			registers[arg1] |= arg2
			i += 3
		case uint8(MOV2W):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint64 = uint64(RAM[i+2]) | (uint64(RAM[i+3]) << 8)
			registers[arg1] &= 0xFFFFFFFF0000FFFF
			registers[arg1] |= (arg2 << 16)
			i += 3
		case uint8(MOV3W):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint64 = uint64(RAM[i+2]) | (uint64(RAM[i+3]) << 8)
			registers[arg1] &= 0xFFFF0000FFFFFFFF
			registers[arg1] |= (arg2 << 32)
			i += 3
		case uint8(MOV4W):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint64 = uint64(RAM[i+2]) | (uint64(RAM[i+3]) << 8)
			registers[arg1] &= 0x0000FFFFFFFFFFFF
			registers[arg1] |= (arg2 << 48)
			i += 3
		// case MOVR
		// case SWAP
		// case PUSH
		// case PUSHIB
		// case PUSHIW
		// case PUSHIT
		// case POP
		// case PEEK
		case uint8(CMP):
			i += 1
			var arg1 uint8 = RAM[i]
			i += 1
			var arg2 uint8 = RAM[i]
			i += 1
			var arg3 uint8 = RAM[i]
			switch arg3 {
			case 1:
				if !(registers[arg1]^0x8000000000000000 < registers[arg2]^0x8000000000000000) {
					i += 3
				}
			case 2:
				if !(registers[arg1]^0x8000000000000000 > registers[arg2]^0x8000000000000000) {
					i += 3
				}
			case 3:
				if registers[arg1] != registers[arg2] {
					i += 3
				}
			case 4:
				if registers[arg1] == registers[arg2] {
					i += 3
				}
			}
		case uint8(JMPB):
			var offset uint32
			offset = uint32(RAM[i+1])
			if offset&0x80 != 0 {
				offset |= 0xFFFFFF00
			}
			i += offset
		case uint8(JMPW):
			var offset uint32
			offset = uint32(RAM[i+1]) | uint32(RAM[i+2])<<8
			if offset&0x8000 != 0 {
				offset |= 0xFFFF0000
			}
			i += offset
		case uint8(JMPT):
			var offset uint32
			offset = uint32(RAM[i+1]) | uint32(RAM[i+2])<<8 | uint32(RAM[i+3])<<16
			if offset&0x800000 != 0 {
				offset |= 0xFF000000
			}
			i += offset
		// case CALLB
		// case CALLW
		// case CALLT
		case uint8(RET):
			// TO DO
			if uint32(registers[15]) == uint32(RAMSize-1) {
				log.Fatal("Cannot return because stack is empty at memory address : " + intToStr(int(i)))
			}
			var newAddress uint32
			registers[15] += 8
			newAddress = uint32(RAM[uint32(registers[15])]) | uint32(RAM[uint32(registers[15])-1])<<8 | uint32(RAM[uint32(registers[15])-2])<<16 | uint32(RAM[uint32(registers[15])-3])<<24
			if newAddress > stackUpperBound {
				log.Fatal("Address out of bounds")
			}
			i = newAddress

		case uint8(PUSH):
			// TO Do
			var arg uint8 = RAM[i+1]
			if uint32(registers[15]) <= stackUpperBound {
				log.Fatal("Stack overflow (but not the website unfortunately)")
			}
			for j := range 8 {
				RAM[uint32(registers[15])-uint32(j)] = uint8(registers[arg] >> (8 * j))
			}
			registers[15] -= 8
			i += 1
		case uint8(POP):
			// TO DO
			var arg uint8 = RAM[i+1]
			if uint32(registers[15]) >= stackLowerBound {
				log.Fatal("Stack underflow")
			}
			registers[15] += 8
			var number uint64
			for j := range 8 {
				number += uint64(RAM[uint32(registers[15])-uint32(j)]) << (8 * j)
			}
			registers[arg] = number
			i += 1
		case uint8(WRT):
			// TO DO
			var arg1 uint8 = RAM[i+1]
			var arg2 uint8 = RAM[i+2]
			if registers[arg2] <= uint64(stackUpperBound) {
				log.Fatal("You cannot modify the program while running")
			} else if registers[arg2] > uint64(RAMSize) {
				log.Fatal("Address out of bounds")
			}
			var arg3 uint8 = RAM[i+3]
			var numberToStore uint64 = registers[arg3]
			var bytes uint8
			for j := range arg1 {
				bytes = uint8(numberToStore & 255)
				RAM[registers[arg2]+uint64(j)] = bytes
				numberToStore = numberToStore >> 8
			}
			i += 3
		case uint8(READ):
			// TO DO
			var arg1 uint8 = RAM[i+1]
			var arg2 uint8 = RAM[i+2]
			var arg3 uint8 = RAM[i+3]
			if registers[arg3] > uint64(RAMSize) {
				log.Fatal("Address out of bounds")
			}
			var storedNumber uint64 = 0
			for j := 0; uint8(j) < arg2; j++ {
				storedNumber += uint64(RAM[registers[arg3]+uint64(j)]) << (8 * j)
			}
			registers[arg1] = storedNumber
			i += 3
		}
		//fmt.Println(debugVariable, opcodeToMnemonics[int(RAM[debugVariable])], registers)
		//fmt.Println(RAM[3*(RAMSize>>2):])
		//fmt.Println(RAM[RAMSize>>2 : RAMSize-(RAMSize>>2)])
		//fmt.Println(RAM)
	}
	fmt.Println()
	fmt.Println(registers)
	fmt.Println(RAM)
}
