package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

//////////
// DATA //
//////////

var mnemonics []string = []string{"MOV", "ADDI", "ADD", "AND", "OR", "NOT", "PUSH", "POP", "SWAP", "CMP", "JMP", "RET", "HLT", "WRT", "READ", "CALL"}
var registersName []string = []string{"R0", "R1", "R2", "R3", "R4", "R5", "R6", "R7", "R8", "R9", "R10", "R11", "R12", "R13", "R14", "R15"}
var registers []uint64 = []uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

const RAMSize int = 512

var RAM [RAMSize]uint8
var stackPointer uint32 = uint32(RAMSize) - 1

const (
	HLT int = iota
	RET
	AND
	OR
	NOT
	ADD
	ADDI
	MOV
	PUSH
	POP
	CMP
	JMP
	WRT
	READ
	SWAP
	CALL
)

var opcodeToMnemonics = map[int]string{
	HLT:  "HLT",
	RET:  "RET",
	AND:  "AND",
	OR:   "OR",
	NOT:  "NOT",
	ADD:  "ADD",
	ADDI: "ADDI",
	MOV:  "MOV",
	PUSH: "PUSH",
	POP:  "POP",
	CMP:  "CMP",
	JMP:  "JMP",
	WRT:  "WRT",
	READ: "READ",
	SWAP: "SWAP",
	CALL: "CALL",
}

var syntaxRules = map[string][]string{
	"HLT":  {},
	"RET":  {},
	"AND":  {"Register", "Register"},
	"OR":   {"Register", "Register"},
	"NOT":  {"Register"},
	"ADD":  {"Register", "Register"},
	"ADDI": {"Register", "Number1"},
	"MOV":  {"Register", "Number8"},
	"PUSH": {"Register"},
	"POP":  {"Register"},
	"CMP":  {"Register", "Register", "Comparison"},
	"JMP":  {"Offset"},
	"WRT":  {"Size", "Address", "Register"},
	"READ": {"Register", "Size", "Address"},
	"SWAP": {"Register", "Register"},
	"CALL": {"Offset"},
}

var memorySize = map[string]int{
	"HLT":  1,
	"RET":  1,
	"AND":  3,
	"OR":   3,
	"NOT":  2,
	"ADD":  3,
	"ADDI": 3,
	"MOV":  10,
	"PUSH": 2,
	"POP":  2,
	"CMP":  4,
	"JMP":  5,
	"WRT":  7,
	"READ": 7,
	"SWAP": 3,
	"CALL": 5,
}

var forbiddenLabels []string = []string{"R0", "R1", "R2", "R3", "R4", "R5", "R6", "R7", "R8", "R9", "R10", "R11", "R12", "R13", "R14", "R15",
	"r0", "r1", "r2", "r3", "r4", "r5", "r6", "r7", "r8", "r9", "r10", "r11", "r12", "r13", "r14", "r15",
	"HLT", "RET", "AND", "OR", "NOT", "ADD", "ADDI", "MOV", "PUSH", "POP", "CMP", "JMP", "WRT", "READ", "SWAP",
	"E", "G", "L"}

//////////
// MAIN //
//////////

func main() {
	args := os.Args[1:] // Skip the program name
	content, err := os.ReadFile(args[0])
	if err != nil {
		log.Fatal("\rCouldn't read file")
	}

	var program string = string(content)
	var assemblerProgram [][]string = readProgram(program)

	var startTime time.Time = time.Now()
	var byteProgram []uint8 = programCleaner(assemblerProgram)
	var elapsed time.Duration = time.Since(startTime)
	fmt.Println(byteProgram)
	fmt.Println(len(byteProgram))
	fmt.Printf("Temps : %s\n", elapsed)

	writeToRAM(byteProgram)
	fmt.Println(RAM)

	startTime = time.Now()
	executeProgram()
	elapsed = time.Since(startTime)
	fmt.Printf("Temps : %s\n", elapsed)
}

/////////////////
// INTERPRETER //
/////////////////

func readProgram(program string) [][]string {
	var operations []string = strings.Split(program, "\n")
	var assemblerProgram [][]string

	for _, line := range operations {
		assemblerProgram = append(assemblerProgram, strings.Fields(line))
	}

	return assemblerProgram
}

///////////////////////
// Clean the program //
///////////////////////

func programCleaner(assemblerProgram [][]string) []uint8 {
	var assemblerProgramWithBlankLine [][]string = assemblerProgram
	var numberOfBlankLines int
	assemblerProgram = cleanEmpty(assemblerProgram)

	var labels = make(map[string]int)
	var tokenizedProgram [][][]string

	var memoryAddress int
	for i, line := range assemblerProgram {
		for isEmpty(assemblerProgramWithBlankLine[i+numberOfBlankLines]) {
			numberOfBlankLines += 1
		}
		line = checkUnexpectedCharacter(line)
		checkNumberOfArgs(line, i+numberOfBlankLines)
		tokenizedProgram = append(tokenizedProgram, checkWords(line, i+numberOfBlankLines))
		labels = checkJumps(tokenizedProgram[i], labels, memoryAddress)
		checkSyntax(tokenizedProgram[i], syntaxRules[tokenizedProgram[i][0][0]])
		memoryAddress += memorySize[tokenizedProgram[i][0][0]]
	}
	tokenizedProgram = delLabels(tokenizedProgram)

	memoryAddress = 0
	var opcodeProgram [][]uint64
	for i, line := range tokenizedProgram {
		if line[0][0] == "JMP" || line[0][0] == "CALL" {
			tokenizedProgram[i] = createJumpAddress(labels, line, memoryAddress)
		}
		opcodeProgram = append(opcodeProgram, mnemonicsToOpcode(line))
		memoryAddress += memorySize[tokenizedProgram[i][0][0]]
	}

	var byteProgram []uint8 = bytificationOfTheProgram(opcodeProgram)
	return byteProgram
}

func isEmpty(line []string) bool {
	for _, word := range line {
		if len(word) != 0 {
			return false
		}
	}
	return true
}

func cleanEmpty(assemblerProgram [][]string) [][]string {
	var cleanedProgram [][]string
	for _, line := range assemblerProgram {
		var cleanedLine []string
		for _, word := range line {
			if len(word) != 0 {
				cleanedLine = append(cleanedLine, word)
			}
		}
		if len(cleanedLine) != 0 {
			cleanedProgram = append(cleanedProgram, cleanedLine)
		}
	}
	return cleanedProgram
}

func checkUnexpectedCharacter(line []string) []string {
	validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890:-*@"
	for i := range len(line) {
		var cleanedString string = ""
		for _, character := range line[i] {
			if strings.Contains(validChars, string(character)) {
				cleanedString += string(character)
			}
		}
		line[i] = cleanedString
	}
	return line
}

func checkNumberOfArgs(line []string, i int) {
	if len(syntaxRules[line[0]]) != len(line)-1 {
		err := "Wrong number of args for \"" + line[0] + "\" at line " + intToStr(i+1)
		log.Fatal(err)
	}
}

func checkWords(line []string, i int) [][]string {
	var newLine [][]string
	for j, word := range line {
		if inList(mnemonics, word) {
			newLine = append(newLine, []string{word, "Operation"})
		} else if inList([]string{"G", "L", "E"}, word) {
			newLine = append(newLine, []string{word, "Comparison"})
		} else if word[len(word)-1] == ':' || (j > 0 && (line[j-1] == "JMP" || line[j-1] == "CALL")) {
			newLine = append(newLine, []string{word, "Offset"})
		} else if inList(registersName, word) {
			newLine = append(newLine, []string{word[1:], "Register"})
		} else if word[0] == '@' && isInt(word[1:]) && isPowerOfTwo(strToInt(word[1:])) && strToInt(word[1:]) >= 8 {
			newLine = append(newLine, []string{intToStr(exponentOfPowerOfTwo(strToInt(word[1:])) - 2), "Size"})
		} else if word[0] == '*' && inList(registersName, word[1:]) {
			newLine = append(newLine, []string{word[2:], "Address"})
		} else if isInt(word) {
			var number int = strToInt(word)
			if line[0] == "ADDI" && number < 129 && number > -128 {
				newLine = append(newLine, []string{word, "Number1"})
			} else if line[0] == "MOV" && number > int(-math.Pow(2, 63)) && number < int(math.Pow(2, 63))-1 {
				newLine = append(newLine, []string{word, "Number8"})
			} else {
				err := "Syntax error \"" + word + "\" at line " + intToStr(i+1)
				log.Fatal(err)
			}
		} else {
			err := "Unrecognized token \"" + word + "\" at line " + intToStr(i+1)
			log.Fatal(err)
		}
	}
	newLine = append(newLine, []string{intToStr(i), "Line"})
	return newLine
}

func checkJumps(line [][]string, labels map[string]int, memoryAddress int) map[string]int {
	if line[0][0][len(line[0][0])-1] == ':' {
		if !(inList(forbiddenLabels, line[0][0][:len(line[0][0])-1])) {
			labels[string(line[0][0][:len(line[0][0])-1])] = memoryAddress - 1
		} else {
			err := "Forbiddent label name \"" + string(line[0][0][:len(line[0][0])-1]) + "\" at line " + intToStr(strToInt(string(line[1][0]))+1)
			log.Fatal(err)
		}
	}
	return labels
}

func checkSyntax(line [][]string, rules []string) {
	var errorSyntax bool = false
	var numberLine int
	for j, rule := range rules {
		if rule == "Register" && line[j+1][1] != "Register" {
			errorSyntax = true
		} else if rule == "Comparison" && line[j+1][1] != "Comparison" {
			errorSyntax = true
		} else if rule == "Offset" && line[j+1][1] != "Offset" {
			errorSyntax = true
		} else if rule == "Address" && line[j+1][1] != "Address" {
			errorSyntax = true
		} else if rule == "Size" && line[j+1][1] != "Size" {
			errorSyntax = true
		} else if rule == "Number1" && line[j+1][1] != "Number1" {
			errorSyntax = true
		} else if rule == "Number8" && line[j+1][1] != "Number8" {
			errorSyntax = true
		}
		numberLine = j
	}
	if errorSyntax == true {
		fmt.Println(line)
		err := "Syntax error at line " + intToStr(strToInt(string(line[numberLine+2][0][0]))+1)
		log.Fatal(err)
	}
}

func delLabels(tokenizedProgram [][][]string) [][][]string {
	var cleanedProgram [][][]string
	for _, line := range tokenizedProgram {
		if line[0][1] != "Offset" {
			cleanedProgram = append(cleanedProgram, line)
		}
	}
	return cleanedProgram
}

func createJumpAddress(labels map[string]int, line [][]string, memoryAdress int) [][]string {
	var targetLine int = labels[line[1][0]]
	if targetLine == 0 {
		err := "Undefined label \"" + line[1][0] + "\""
		log.Fatal(err)
	}
	line[1][0] = intToStr(targetLine - memoryAdress)
	return line
}

func mnemonicsToOpcode(line [][]string) []uint64 {
	var newLine []uint64
	if string(line[0][0]) == "MOV" {
		var arg1 uint64 = uint64(strToInt(line[1][0]))
		var arg2 string = line[2][0]
		if arg2[0] == '-' {
			arg2 = arg2[1:]
			var arg2 uint64 = uint64(^strToInt(arg2) + 1)
			newLine = []uint64{uint64(MOV), uint64(arg1), arg2}
		} else {
			newLine = []uint64{uint64(MOV), uint64(arg1), uint64(strToInt(arg2))}
		}
	} else if string(line[0][0]) == "ADD" {
		var arg1 uint64 = uint64(strToInt(line[1][0]))
		var arg2 uint64 = uint64(strToInt(line[2][0]))
		newLine = []uint64{uint64(ADD), arg1, arg2}

	} else if string(line[0][0]) == "ADDI" {
		var arg1 uint64 = uint64(strToInt(line[1][0]))
		var arg2 string = line[2][0]
		if arg2[0] == '-' {
			arg2 = arg2[1:]
			var arg2 uint64 = uint64(^strToInt(arg2) + 1)
			newLine = []uint64{uint64(ADDI), arg1, arg2}
		} else {
			newLine = []uint64{uint64(ADDI), arg1, uint64(strToInt(arg2))}
		}
	} else if string(line[0][0]) == "PUSH" {
		var arg1 uint64 = uint64(strToInt(line[1][0]))
		newLine = []uint64{uint64(PUSH), arg1}

	} else if string(line[0][0]) == "POP" {
		var arg1 uint64 = uint64(strToInt(line[1][0]))
		newLine = []uint64{uint64(POP), arg1}

	} else if string(line[0][0]) == "AND" {
		var arg1 uint64 = uint64(strToInt(line[1][0]))
		var arg2 uint64 = uint64(strToInt(line[2][0]))
		newLine = []uint64{uint64(AND), arg1, arg2}

	} else if string(line[0][0]) == "OR" {
		var arg1 uint64 = uint64(strToInt(line[1][0]))
		var arg2 uint64 = uint64(strToInt(line[2][0]))
		newLine = []uint64{uint64(OR), arg1, arg2}

	} else if string(line[0][0]) == "NOT" {
		var arg1 uint64 = uint64(strToInt(line[1][0]))
		newLine = []uint64{uint64(NOT), arg1}

	} else if string(line[0][0]) == "SWAP" {
		var arg1 uint64 = uint64(strToInt(line[1][0]))
		var arg2 uint64 = uint64(strToInt(line[2][0]))
		newLine = []uint64{uint64(SWAP), arg1, arg2}

	} else if string(line[0][0]) == "CMP" {
		var arg1 uint64 = uint64(strToInt(line[1][0]))
		var arg2 uint64 = uint64(strToInt(line[2][0]))
		var arg3 string = line[3][0]
		if inList([]string{"L", "G", "E"}, arg3) {
			if arg3 == "L" {
				newLine = []uint64{uint64(CMP), arg1, arg2, uint64(1)}
			} else if arg3 == "G" {
				newLine = []uint64{uint64(CMP), arg1, arg2, uint64(2)}
			} else if arg3 == "E" {
				newLine = []uint64{uint64(CMP), arg1, arg2, uint64(3)}
			}
		}
	} else if string(line[0][0]) == "JMP" {
		var arg1 string = line[1][0]
		if arg1[0] == '-' {
			arg1 = arg1[1:]
			var arg1 uint64 = uint64(^strToInt(arg1) + 1)
			newLine = []uint64{uint64(JMP), arg1}
		} else {
			newLine = []uint64{uint64(JMP), uint64(strToInt(arg1))}
		}
	} else if string(line[0][0]) == "RET" {
		newLine = []uint64{uint64(RET)}
	} else if string(line[0][0]) == "HLT" {
		newLine = []uint64{uint64(HLT)}
	} else if string(line[0][0]) == "WRT" {
		var arg1 uint64 = uint64(strToInt(line[1][0]))
		var arg2 uint64 = uint64(strToInt(line[2][0]))
		var arg3 uint64 = uint64(strToInt(line[3][0]))
		newLine = []uint64{uint64(WRT), arg1, arg2, arg3}
	} else if string(line[0][0]) == "READ" {
		var arg1 uint64 = uint64(strToInt(line[1][0]))
		var arg2 uint64 = uint64(strToInt(line[2][0]))
		var arg3 uint64 = uint64(strToInt(line[3][0]))
		newLine = []uint64{uint64(READ), arg1, arg2, arg3}
	} else if string(line[0][0]) == "CALL" {
		var arg1 string = line[1][0]
		if arg1[0] == '-' {
			arg1 = arg1[1:]
			var arg1 uint64 = uint64(^strToInt(arg1) + 1)
			newLine = []uint64{uint64(CALL), arg1}
		} else {
			newLine = []uint64{uint64(CALL), uint64(strToInt(arg1))}
		}
	} else {
		log.Fatal("Err in mnemonicsToOpcode : " + string(line[0][0]))
	}
	return newLine
}

func bytificationOfTheProgram(opcodeProgram [][]uint64) []uint8 {
	var byteProgram []uint8
	for _, line := range opcodeProgram {
		switch line[0] {
		case uint64(HLT):
			byteProgram = append(byteProgram, uint8(HLT))
		case uint64(RET):
			byteProgram = append(byteProgram, uint8(RET))
		case uint64(MOV):
			byteProgram = append(byteProgram, uint8(MOV))
			byteProgram = append(byteProgram, uint8(line[1]))
			var argument uint8
			for range 8 {
				argument = uint8(line[2] & 255)
				byteProgram = append(byteProgram, argument)
				line[2] >>= 8
			}
		case uint64(ADD):
			byteProgram = append(byteProgram, uint8(ADD))
			byteProgram = append(byteProgram, uint8(line[1]))
			byteProgram = append(byteProgram, uint8(line[2]))
		case uint64(ADDI):
			byteProgram = append(byteProgram, uint8(ADDI))
			byteProgram = append(byteProgram, uint8(line[1]))
			byteProgram = append(byteProgram, uint8(line[2]&255))
		case uint64(PUSH):
			byteProgram = append(byteProgram, uint8(PUSH))
			byteProgram = append(byteProgram, uint8(line[1]))
		case uint64(POP):
			byteProgram = append(byteProgram, uint8(POP))
			byteProgram = append(byteProgram, uint8(line[1]))
		case uint64(AND):
			byteProgram = append(byteProgram, uint8(AND))
			byteProgram = append(byteProgram, uint8(line[1]))
			byteProgram = append(byteProgram, uint8(line[2]))
		case uint64(OR):
			byteProgram = append(byteProgram, uint8(OR))
			byteProgram = append(byteProgram, uint8(line[1]))
			byteProgram = append(byteProgram, uint8(line[2]))
		case uint64(NOT):
			byteProgram = append(byteProgram, uint8(NOT))
			byteProgram = append(byteProgram, uint8(line[1]))
		case uint64(CMP):
			byteProgram = append(byteProgram, uint8(CMP))
			byteProgram = append(byteProgram, uint8(line[1]))
			byteProgram = append(byteProgram, uint8(line[2]))
			byteProgram = append(byteProgram, uint8(line[3]))
		case uint64(JMP):
			byteProgram = append(byteProgram, uint8(JMP))
			var argument uint8
			for range 4 {
				argument = uint8(line[1] & 255)
				byteProgram = append(byteProgram, argument)
				line[1] >>= 8
			}
		case uint64(WRT):
			byteProgram = append(byteProgram, uint8(WRT))
			byteProgram = append(byteProgram, uint8(line[1]))
			byteProgram = append(byteProgram, uint8(line[2]))
			byteProgram = append(byteProgram, uint8(line[3]))
		case uint64(READ):
			byteProgram = append(byteProgram, uint8(READ))
			byteProgram = append(byteProgram, uint8(line[1]))
			byteProgram = append(byteProgram, uint8(line[2]))
			byteProgram = append(byteProgram, uint8(line[3]))
		case uint64(SWAP):
			byteProgram = append(byteProgram, uint8(SWAP))
			byteProgram = append(byteProgram, uint8(line[1]))
			byteProgram = append(byteProgram, uint8(line[2]))
		case uint64(CALL):
			byteProgram = append(byteProgram, uint8(CALL))
			var argument uint8
			for range 4 {
				argument = uint8(line[1] & 255)
				byteProgram = append(byteProgram, argument)
				line[1] >>= 8
			}
		}
	}
	return byteProgram
}

//////////////////
// Load program //
//////////////////

func writeToRAM(byteProgram []uint8) {
	for i, byte := range byteProgram {
		RAM[i] = byte
	}
}

/////////////////////////
// Execute the program //
/////////////////////////

func executeProgram() {
loop:
	for i := uint32(0); i < uint32(RAMSize/4); i++ {
		//var debugVariable uint32 = i
		switch RAM[i] {
		case uint8(HLT):
			break loop
		case uint8(RET):
			if stackPointer == uint32(RAMSize-1) {
				log.Fatal("Cannot return because stack is empty at memory address : " + intToStr(int(i)))
			}
			var newAddress uint32
			stackPointer -= 8
			for i := range 4 {
				newAddress += uint32(RAM[stackPointer+uint32(i)] << (8 * i))
			}
			if registers[newAddress] < 0 || registers[newAddress] >= uint64(RAMSize) {
				log.Fatal("Address out of bounds")
			}
			i = newAddress
		case uint8(MOV):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint64
			for j := range 8 {
				arg2 += uint64(RAM[i+2+uint32(j)]) << (8 * j)
			}
			registers[arg1] = arg2
			i += 9
		case uint8(ADD):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint8 = RAM[i+2]
			registers[arg1] = (registers[arg1] + registers[arg2])
			i += 2
		case uint8(ADDI):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint64 = uint64(RAM[i+2])
			registers[arg1] += arg2 | 0xFFFFFFFFFFFFFF00
			i += 2
		case uint8(PUSH):
			var arg uint8 = RAM[i+1]
			if stackPointer >= uint32(RAMSize-RAMSize/4-1) {
				log.Fatal("Stack overflow (but not the website unfortunately)")
			}
			for j := range 8 {
				RAM[stackPointer-uint32(j)] = uint8(registers[arg] >> (8 * j))
			}
			stackPointer -= 8
			i += 1
		case uint8(POP):
			var arg uint8 = RAM[i+1]
			if stackPointer <= uint32(RAMSize-1) {
				log.Fatal("Stack underflow")
			}
			stackPointer += 8
			for j := range 8 {
				registers[arg] += uint64(RAM[stackPointer-uint32(j)] << (8 * j))
			}
			i += 1
		case uint8(AND):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint8 = RAM[i+2]
			registers[arg1] = registers[arg1] & registers[arg2]
			i += 2
		case uint8(OR):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint8 = RAM[i+2]
			registers[arg1] = registers[arg1] | registers[arg2]
			i += 2
		case uint8(NOT):
			var arg uint8 = RAM[i+1]
			registers[arg] = ^registers[arg]
			i += 1
		case uint8(CMP):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint8 = RAM[i+2]
			var arg3 uint8 = RAM[i+3]
			switch arg3 {
			case 1:
				if !(registers[arg1]^0x8000000000000000 < registers[arg2]^0x8000000000000000) {
					i += uint32(memorySize[opcodeToMnemonics[int(RAM[int(i+4)])]])
				}
			case 2:
				if !(registers[arg1]^0x8000000000000000 > registers[arg2]^0x8000000000000000) {
					i += uint32(memorySize[opcodeToMnemonics[int(RAM[int(i+4)])]])
				}
			case 3:
				if registers[arg1] != registers[arg2] {
					i += uint32(memorySize[opcodeToMnemonics[int(RAM[int(i+4)])]])
				}
			}
			i += 3
		case uint8(JMP):
			var offset uint32
			for j := range 4 {
				offset += uint32(RAM[i+uint32(1+j)]) << (8 * j)
			}
			i += offset
		case uint8(WRT):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint8 = RAM[i+2]
			if registers[arg2] <= uint64(RAMSize/4-1) {
				log.Fatal("You cannot modify the program while running")
			} else if registers[arg2] < 0 || registers[arg2] >= uint64(RAMSize) {
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
			var arg1 uint8 = RAM[i+1]
			var arg2 uint8 = RAM[i+2]
			var arg3 uint8 = RAM[i+3]
			if registers[arg3] < 0 || registers[arg3] >= uint64(RAMSize) {
				log.Fatal("Address out of bounds")
			}
			var storedNumber uint64 = 0
			for j := 0; uint8(j) < arg2; j++ {
				storedNumber += uint64(RAM[registers[arg3]+uint64(j)] << (8 * j))
			}
			registers[arg1] = storedNumber
			i += 3
		case uint8(SWAP):
			var arg1 uint8 = RAM[i+1]
			var arg2 uint8 = RAM[i+2]
			intermediateVariable := registers[arg1]
			registers[arg1] = registers[arg2]
			registers[arg2] = intermediateVariable
			i += 2
		case uint8(CALL):

		}
		//fmt.Println(debugVariable, opcodeToMnemonics[int(RAM[debugVariable])], registers)
	}
	fmt.Println(registers, RAM)
}

///////////
// UTILS //
///////////

func strToInt(x string) int {
	num, err := strconv.Atoi(x)
	if err != nil {
		fmt.Println("Error in strToInt : " + x)
		return -1
	}
	return num
}

func isInt(x string) bool {
	for _, char := range x[1:] {
		if !(strings.Contains("0123456789", string(char))) {
			return false
		}
	}
	if !(strings.Contains("-0123456789", string(x[0]))) {
		return false
	}
	return true
}

func intToStr(x int) string {
	num := strconv.Itoa(x)
	return num
}

func inList(liste []string, item string) bool {
	for _, element := range liste {
		if element == item {
			return true
		}
	}
	return false
}

func inListBool(liste []bool, item bool) bool {
	for _, element := range liste {
		if element == item {
			return true
		}
	}
	return false
}

func isPowerOfTwo(x int) bool {
	return x >= 2 && (x&(x-1)) == 0
}

func exponentOfPowerOfTwo(x int) int {
	var exponent int = 0
	for x > 1 {
		x = x >> 1
		exponent += 1
	}
	return exponent
}
