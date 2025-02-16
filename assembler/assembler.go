package main

import (
	"fmt"
	"log"
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
var registers []int = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var RAM [16]int
var stack []int = []int{}

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
	"ADDI": {"Register", "Number"},
	"MOV":  {"Register", "Number"},
	"PUSH": {"Register"},
	"POP":  {"Register"},
	"CMP":  {"Register", "Register", "Comparison"},
	"JMP":  {"Label"},
	"WRT":  {"Size", "Address", "Register"},
	"READ": {"Register", "Size", "Address"},
	"SWAP": {"Register", "Register"},
	"CALL": {"Label"},
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
	var byteProgram = programCleaner(assemblerProgram)
	var elapsed time.Duration = time.Since(startTime)
	fmt.Println(byteProgram)
	fmt.Printf("Temps : %s\n", elapsed)

	startTime = time.Now()
	executeProgram(byteProgram)
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

func programCleaner(assemblerProgram [][]string) []int {
	assemblerProgram = cleanEmpty(assemblerProgram)

	var labels = make(map[string]int)
	var tokenizedProgram [][][]string

	var memoryAddress int
	for i, line := range assemblerProgram {
		line = checkUnexpectedCharacter(line)
		checkNumberOfArgs(line, i)
		tokenizedProgram = append(tokenizedProgram, checkWords(line, i))
		checkSyntax(tokenizedProgram[i], syntaxRules[tokenizedProgram[i][0][0]], i)
		labels = checkJumps(tokenizedProgram[i], labels, i, memoryAddress)
		memoryAddress += memorySize[tokenizedProgram[i][0][0]]
	}

	tokenizedProgram = delLabels(tokenizedProgram)

	memoryAddress = 0
	var opcodeProgram [][]int
	for i, line := range tokenizedProgram {
		if line[0][0] == "JMP" || line[0][0] == "CALL" {
			tokenizedProgram[i] = createJumpAddress(labels, line, memoryAddress)
		}
		opcodeProgram = append(opcodeProgram, mnemonicsToOpcode(line))
		memoryAddress += memorySize[tokenizedProgram[i][0][0]]
	}

	var bytePogram []int = bytificationOfTheProgram(opcodeProgram)

	return bytePogram
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
			newLine = append(newLine, []string{word, "Label"})
		} else if inList(registersName, word) {
			newLine = append(newLine, []string{word[1:], "Register"})
		} else if word[0] == '@' && isInt(word[1:]) && isPowerOfTwo(strToInt(word[1:])) && strToInt(word[1:]) >= 8 {
			newLine = append(newLine, []string{intToStr(exponentOfPowerOfTwo(strToInt(word[1:])) - 2), "Size"})
		} else if word[0] == '*' && inList(registersName, word[1:]) {
			newLine = append(newLine, []string{word[2:], "Address"})
		} else {
			for _, character := range word {
				if !(strings.Contains("-0123456789", string(character))) {
					err := "Unrecognized word \"" + word + "\" at line " + intToStr(i+1)
					log.Fatal(err)
				}
			}
			newLine = append(newLine, []string{word, "Number"})
		}
	}
	return newLine
}

func checkSyntax(line [][]string, rules []string, i int) {
	var errorSyntax bool = false
	for j, rule := range rules {
		if rule == "Register" && line[j+1][1] != "Register" {
			errorSyntax = true
		} else if rule == "Comparison" && line[j+1][1] != "Comparison" {
			errorSyntax = true
		} else if rule == "Label" && line[j+1][1] != "Label" {
			errorSyntax = true
		} else if rule == "Number" && line[j+1][1] != "Number" {
			errorSyntax = true
		} else if rule == "Address" && line[j+1][1] != "Address" {
			errorSyntax = true
		} else if rule == "Size" && line[j+1][1] != "Size" {
			errorSyntax = true
		}
	}
	if errorSyntax {
		err := "Args don't respect syntax rule for \"" + line[0][0] + "\" at line " + intToStr(i+1)
		log.Fatal(err)
	}
}

func checkJumps(line [][]string, labels map[string]int, i int, memoryAddress int) map[string]int {
	if string(line[0][0][len(line[0][0])-1]) == ":" {
		if !inList(forbiddenLabels, string(line[0][0][:len(line[0][0])-1])) {
			labels[string(line[0][0][:len(line[0][0])-1])] = memoryAddress
		} else {
			err := "Forbiddent label name \"" + string(line[0][0][:len(line[0])-1]) + "\" at line " + intToStr(i)
			log.Fatal(err)
		}
	}
	return labels
}

func delLabels(tokenizedProgram [][][]string) [][][]string {
	var cleanedProgram [][][]string
	for _, line := range tokenizedProgram {
		if line[0][1] != "Label" {
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
	line[1][0] = intToStr(targetLine - memoryAdress - 5)
	return line
}

func mnemonicsToOpcode(line [][]string) []int {
	var newLine []int
	if string(line[0][0]) == "MOV" {
		var arg1 int = strToInt(line[1][0])
		var arg2 int = strToInt(line[2][0])
		newLine = []int{MOV, arg1, arg2}

	} else if string(line[0][0]) == "ADD" {
		var arg1 int = strToInt(line[1][0])
		var arg2 int = strToInt(line[2][0])
		newLine = []int{ADD, arg1, arg2}

	} else if string(line[0][0]) == "ADDI" {
		var arg1 int = strToInt(line[1][0])
		var arg2 int = strToInt(line[2][0])
		newLine = []int{ADDI, arg1, arg2}

	} else if string(line[0][0]) == "PUSH" {
		var arg1 int = strToInt(line[1][0])
		newLine = []int{PUSH, arg1}

	} else if string(line[0][0]) == "POP" {
		var arg1 int = strToInt(line[1][0])
		newLine = []int{POP, arg1}

	} else if string(line[0][0]) == "AND" {
		var arg1 int = strToInt(line[1][0])
		var arg2 int = strToInt(line[2][0])
		newLine = []int{AND, arg1, arg2}

	} else if string(line[0][0]) == "OR" {
		var arg1 int = strToInt(line[1][0])
		var arg2 int = strToInt(line[2][0])
		newLine = []int{OR, arg1, arg2}

	} else if string(line[0][0]) == "NOT" {
		var arg1 int = strToInt(line[1][0])
		newLine = []int{NOT, arg1}

	} else if string(line[0][0]) == "SWAP" {
		var arg1 int = strToInt(line[1][0])
		var arg2 int = strToInt(line[2][0])
		newLine = []int{SWAP, arg1, arg2}

	} else if string(line[0][0]) == "CMP" {
		var arg1 int = strToInt(line[1][0])
		var arg2 int = strToInt(line[2][0])
		var arg3 string = line[3][0]
		if inList([]string{"L", "G", "E"}, arg3) {
			if arg3 == "L" {
				newLine = []int{CMP, arg1, arg2, 1}
			} else if arg3 == "G" {
				newLine = []int{CMP, arg1, arg2, 2}
			} else if arg3 == "E" {
				newLine = []int{CMP, arg1, arg2, 3}
			}
		}
	} else if string(line[0][0]) == "JMP" {
		var arg1 int = strToInt(line[1][0])
		newLine = []int{JMP, arg1}
	} else if string(line[0][0]) == "RET" {
		newLine = []int{RET}
	} else if string(line[0][0]) == "HLT" {
		newLine = []int{HLT}
	} else if string(line[0][0]) == "WRT" {
		var arg1 int = strToInt(line[1][0])
		var arg2 int = strToInt(line[2][0])
		var arg3 int = strToInt(line[3][0])
		newLine = []int{WRT, arg1, arg2, arg3}
	} else if string(line[0][0]) == "READ" {
		var arg1 int = strToInt(line[1][0])
		var arg2 int = strToInt(line[2][0])
		var arg3 int = strToInt(line[3][0])
		newLine = []int{READ, arg1, arg2, arg3}
	} else if string(line[0][0]) == "CALL" {
		var arg1 int = strToInt(line[1][0])
		newLine = []int{CALL, arg1}
	} else {
		log.Fatal("Err in mnemonicsToOpcode : " + string(line[0][0]))
	}
	return newLine
}

func bytificationOfTheProgram(opcodeProgram [][]int) []int {
	var byteProgram []int
	for _, line := range opcodeProgram {
		switch line[0] {
		case HLT:
			byteProgram = append(byteProgram, HLT)
		case RET:
			byteProgram = append(byteProgram, RET)
		case MOV:
			byteProgram = append(byteProgram, MOV)
			byteProgram = append(byteProgram, line[1])
			var argument int
			for i := range 8 {
				line[2] >>= 8 * i
				argument = line[2] & 255
				byteProgram = append(byteProgram, argument)
			}
		case ADD:
			byteProgram = append(byteProgram, ADD)
			byteProgram = append(byteProgram, line[1])
			byteProgram = append(byteProgram, line[2])
		case ADDI:
			byteProgram = append(byteProgram, ADDI)
			byteProgram = append(byteProgram, line[1])
			byteProgram = append(byteProgram, line[2])
		case PUSH:
			byteProgram = append(byteProgram, PUSH)
			byteProgram = append(byteProgram, line[1])
		case POP:
			byteProgram = append(byteProgram, POP)
			byteProgram = append(byteProgram, line[1])
		case AND:
			byteProgram = append(byteProgram, AND)
			byteProgram = append(byteProgram, line[1])
			byteProgram = append(byteProgram, line[2])
		case OR:
			byteProgram = append(byteProgram, OR)
			byteProgram = append(byteProgram, line[1])
			byteProgram = append(byteProgram, line[2])
		case NOT:
			byteProgram = append(byteProgram, NOT)
			byteProgram = append(byteProgram, line[1])
		case CMP:
			byteProgram = append(byteProgram, CMP)
			byteProgram = append(byteProgram, line[1])
			byteProgram = append(byteProgram, line[2])
			byteProgram = append(byteProgram, line[3])
		case JMP:
			byteProgram = append(byteProgram, JMP)
			var argument int
			for i := range 4 {
				line[1] >>= 8 * i
				argument = line[1] & 255
				byteProgram = append(byteProgram, argument)
			}
		case WRT:
			byteProgram = append(byteProgram, WRT)
			byteProgram = append(byteProgram, line[1])
			byteProgram = append(byteProgram, line[2])
			byteProgram = append(byteProgram, line[3])
		case READ:
			byteProgram = append(byteProgram, READ)
			byteProgram = append(byteProgram, line[1])
			byteProgram = append(byteProgram, line[2])
			byteProgram = append(byteProgram, line[3])
		case SWAP:
			byteProgram = append(byteProgram, SWAP)
			byteProgram = append(byteProgram, line[1])
			byteProgram = append(byteProgram, line[2])
		case CALL:
			byteProgram = append(byteProgram, CALL)
			var argument int
			for i := range 4 {
				line[1] >>= 8 * i
				argument = line[1] & 255
				byteProgram = append(byteProgram, argument)
			}
		}
	}
	return byteProgram
}

/////////////////////////
// Execute the program //
/////////////////////////

func executeProgram(byteProgram []int) {
	for i := 0; i < len(byteProgram); i++ {
		switch byteProgram[i] {
		case HLT:
			break
		case RET:
			i = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		case MOV:
			var arg1 int = byteProgram[i+1]
			var arg2 int
			for j := range 8 {
				arg2 += byteProgram[i+2+j] << (8 * j)
			}
			registers[arg1] = arg2
			i += memorySize[opcodeToMnemonics[MOV]] - 1
		case ADD:
			var arg1 int = byteProgram[i+1]
			var arg2 int = byteProgram[i+2]
			registers[arg1] = registers[arg1] + registers[arg2]
			i += memorySize[opcodeToMnemonics[ADD]] - 1
		case ADDI:
			var arg1 int = byteProgram[i+1]
			var arg2 int = byteProgram[i+2]
			registers[arg1] += arg2
			i += memorySize[opcodeToMnemonics[ADDI]] - 1
		case PUSH:
			var arg int = byteProgram[i+1]
			stack = append(stack, registers[arg])
			i += memorySize[opcodeToMnemonics[PUSH]] - 1
		case POP:
			var arg int = byteProgram[i+1]
			registers[arg] = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			i += memorySize[opcodeToMnemonics[POP]] - 1
		case AND:
			var arg1 int = byteProgram[i+1]
			var arg2 int = byteProgram[i+2]
			registers[arg1] = registers[arg1] & registers[arg2]
			i += memorySize[opcodeToMnemonics[AND]] - 1
		case OR:
			var arg1 int = byteProgram[i+1]
			var arg2 int = byteProgram[i+2]
			registers[arg1] = registers[arg1] | registers[arg2]
			i += memorySize[opcodeToMnemonics[OR]] - 1
		case NOT:
			var arg int = byteProgram[i+1]
			registers[arg] = ^registers[arg]
			i += memorySize[opcodeToMnemonics[NOT]] - 1
		case CMP:
			//fmt.Println(i, assemblerProgram[i], registers, stack)
			var arg1 int = byteProgram[i+1]
			var arg2 int = byteProgram[i+2]
			var arg3 int = byteProgram[i+3]
			switch arg3 {
			case 1:
				if !(registers[arg1] < registers[arg2]) {
					i += 1
				}
			case 2:
				if !(registers[arg1] > registers[arg2]) {
					i += 1
				}
			case 3:
				if registers[arg1] != registers[arg2] {
					i += 1
				}
			}
			i += memorySize[opcodeToMnemonics[CMP]] - 1
		case JMP:
			fmt.Println(byteProgram[i], registers, stack, RAM)
			var offset int = 0
			for j := range 4 {
				offset += byteProgram[i+1+j] << (8 * j)
			}
			if offset > 0 {
				i = i + offset + 4
			} else {
				i = i + offset
			}
		case WRT:
			var arg1 int = byteProgram[i+1]
			var arg2 int = byteProgram[i+2]
			var arg3 int = byteProgram[i+3]
			var numberToStore int = registers[arg3]
			var bytes int
			for i := range arg1 {
				bytes = numberToStore & 255
				RAM[registers[arg2]+i] = bytes
				numberToStore = numberToStore >> 8
			}
			i += memorySize[opcodeToMnemonics[WRT]] - 1
		case READ:
			var arg1 int = byteProgram[i+1]
			var arg2 int = byteProgram[i+2]
			var arg3 int = byteProgram[i+3]
			var storedNumber int = 0
			for j := 0; j < arg2; j++ {
				storedNumber += RAM[registers[arg3]+j] << (8 * j)
			}
			registers[arg1] = storedNumber
			i += memorySize[opcodeToMnemonics[READ]] - 1
		case SWAP:
			var arg1 int = byteProgram[i+1]
			var arg2 int = byteProgram[i+2]
			intermediateVariable := registers[arg1]
			registers[arg1] = registers[arg2]
			registers[arg2] = intermediateVariable
			i += memorySize[opcodeToMnemonics[SWAP]] - 1
		case CALL:
			fmt.Println(byteProgram[i], registers, stack, RAM)
			stack = append(stack, i)
			var offset int
			for j := range 4 {
				offset += byteProgram[i+1+j] << (8 * j)
			}
			if offset > 0 {
				i = i + offset + 4
			} else {
				i = i + offset
			}
		}
		fmt.Println(byteProgram[i+1], registers, stack, RAM)
	}
	fmt.Println(registers, stack, RAM)
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
	for _, char := range x {
		if !(strings.Contains("0123456789", string(char))) {
			return false
		}
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
