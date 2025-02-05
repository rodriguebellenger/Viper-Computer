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

var mnemonics []string = []string{"MOV", "ADDI", "ADD", "AND", "OR", "NOT", "PUSH", "POP", "SWAP", "CMP", "JMP", "RET", "HLT"}
var registers []int = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
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
)

var numberOfArgs = map[string]int{
	"HLT":  0,
	"RET":  0,
	"AND":  2,
	"OR":   2,
	"NOT":  1,
	"ADD":  2,
	"ADDI": 2,
	"MOV":  2,
	"PUSH": 1,
	"POP":  1,
	"CMP":  3,
	"JMP":  1,
	"WRT":  0,
	"READ": 0,
	"SWAP": 2,
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
	"WRT":  {},
	"READ": {},
	"SWAP": {"Register", "Register"},
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
	content, err := os.ReadFile("assembler/program_test/" + args[0])
	if err != nil {
		log.Fatal("\rCouldn't read file")
	}

	var program string = string(content)
	var assemblerProgram [][]string = readProgram(program)

	var startTime time.Time = time.Now()
	var opcodeProgram = programCleaner(assemblerProgram)
	var elapsed time.Duration = time.Since(startTime)
	fmt.Println(opcodeProgram)
	fmt.Printf("Temps : %s\n", elapsed)

	startTime = time.Now()
	executeProgram(opcodeProgram)
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

func programCleaner(assemblerProgram [][]string) [][]int {
	assemblerProgram = cleanEmpty(assemblerProgram)
	fmt.Println(assemblerProgram)
	var labels = make(map[string]int)
	var tokenizedProgram [][][]string
	var skippedLine = 0
	for i, line := range assemblerProgram {
		var skipOrNot bool = skipEmptyLine(line)
		if !(skipOrNot) {
			line = checkUnexpectedCharacter(line)
			checkNumberOfArgs(line, i)
			tokenizedProgram = append(tokenizedProgram, checkWords(line, i))
			checkSyntax(tokenizedProgram[i-skippedLine], syntaxRules[tokenizedProgram[i-skippedLine][0][0]], i)
			labels = checkJumps(tokenizedProgram[i-skippedLine], labels, i)
		} else {
			skippedLine += 1
		}
	}

	var opcodeProgram [][]int
	var finishedLine []int
	skippedLine = 0
	for i, line := range tokenizedProgram {
		if line[0][0] == "JMP" {
			tokenizedProgram[i] = createJumpAddress(labels, line, i)
		}
		finishedLine = mnemonicsToOpcode(line)
		if len(finishedLine) != 0 {
			opcodeProgram = append(opcodeProgram, finishedLine)
		} else {
			skippedLine += 1
		}
	}
	return opcodeProgram
}

func cleanEmpty(assemblerProgram [][]string) [][]string {
	var skippedLine int = 0
	for i, line := range assemblerProgram {
		for j, ope := range line {
			if ope == "" || ope == " " {
				assemblerProgram[i] = append(assemblerProgram[i][:j], assemblerProgram[i][j+1:]...)
			}
		}
		if len(line) == 0 {
			assemblerProgram = append(assemblerProgram[:i-skippedLine], assemblerProgram[i+1-skippedLine:]...)
			skippedLine += 1
		}
	}
	return assemblerProgram
}

func skipEmptyLine(line []string) bool {
	return len(line) == 0
}

func checkUnexpectedCharacter(line []string) []string {
	validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890:-"
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
	if numberOfArgs[line[0]] != len(line)-1 {
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
		} else if word[len(word)-1] == ':' || (j > 0 && line[j-1] == "JMP") {
			newLine = append(newLine, []string{word, "Label"})
		} else if len(word[1:]) > 0 && strToInt(word[1:]) < 16 && strToInt(word[1:]) >= 0 && inList([]string{"r", "R"}, string(word[0])) {
			newLine = append(newLine, []string{word[1:], "Register"})
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
		if rule == "Register" && (line[j+1][1] != "Register" || strToInt(line[j+1][0]) > 15 || strToInt(line[j+1][0]) < 0) {
			errorSyntax = true
		} else if rule == "Comparison" && (line[j+1][1] != "Comparison" || !(inList([]string{"G", "L", "E"}, line[j+1][0]))) {
			errorSyntax = true
		} else if rule == "Label" && (line[j+1][1] != "Label") {
			errorSyntax = true
		} else if rule == "Number" && (line[j+1][1] != "Number") {
			errorSyntax = true
		}
	}
	if errorSyntax {
		err := "Args don't respect syntax rule for \"" + line[0][0] + "\" at line " + intToStr(i+1)
		log.Fatal(err)
	}
}

func checkJumps(line [][]string, labels map[string]int, i int) map[string]int {
	if string(line[0][len(line[0])-1]) == ":" {
		if !inList(forbiddenLabels, string(line[0][0][:len(line[0][0])-1])) {
			labels[string(line[0][0][:len(line[0])-1])] = i
		} else {
			err := "Forbiddent label name \"" + string(line[0][0][:len(line[0])-1]) + "\" at line " + intToStr(i)
			log.Fatal(err)
		}
	}
	return labels
}

func createJumpAddress(labels map[string]int, line [][]string, i int) [][]string {
	var targetLine int = labels[line[1][0]]
	if targetLine == 0 {
		err := "Undefined label \"" + line[1][0] + "\""
		log.Fatal(err)
	}
	line[1][0] = intToStr(targetLine - i)
	return line
}

func mnemonicsToOpcode(line [][]string) []int {
	var newLine []int
	//fmt.Println(line[1][0])
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
		} else {
			err := "Unrecognized comparison character \"" + arg3 + "\""
			log.Fatal(err)
		}
	} else if string(line[0][0]) == "JMP" {
		var arg1 int = strToInt(line[1][0])
		newLine = []int{JMP, arg1}
	} else if string(line[0][0]) == "RET" {
		newLine = []int{RET}
	} else if string(line[0][len(line[0])-1]) == ":" {

	} else if string(line[0][0]) == "HLT" {
		newLine = []int{HLT}
	} else if string(line[0][0][len(line[0][0])-1]) == ":" {
		newLine = []int{}
	} else {
		log.Fatal("Err in mnemonicsToOpcode : " + string(line[0][0]))
	}
	return newLine
}

/////////////////////////
// Execute the program //
/////////////////////////

func executeProgram(assemblerProgram [][]int) {
	for i := 0; i < len(assemblerProgram); i++ {
		switch assemblerProgram[i][0] {
		case HLT:
			break
		case RET:
			i = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		case MOV:
			var arg int = assemblerProgram[i][1]
			registers[arg] = assemblerProgram[i][2]
		case ADD:
			var arg1 int = assemblerProgram[i][1]
			var arg2 int = assemblerProgram[i][2]
			registers[arg1] = registers[arg1] + registers[arg2]
		case ADDI:
			var arg1 int = assemblerProgram[i][1]
			var arg2 int = assemblerProgram[i][2]
			registers[arg1] = registers[arg1] + arg2
		case PUSH:
			var arg int = assemblerProgram[i][1]
			stack = append(stack, registers[arg])
		case POP:
			var arg int = assemblerProgram[i][1]
			registers[arg] = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		case AND:
			var arg1 int = assemblerProgram[i][1]
			var arg2 int = assemblerProgram[i][2]
			registers[arg1] = registers[arg1] & registers[arg2]
		case OR:
			var arg1 int = assemblerProgram[i][1]
			var arg2 int = assemblerProgram[i][2]
			registers[arg1] = registers[arg1] | registers[arg2]
		case NOT:
			var arg int = assemblerProgram[i][1]
			registers[arg] = ^registers[arg]
		case SWAP:
			var arg1 int = assemblerProgram[i][1]
			var arg2 int = assemblerProgram[i][2]
			intermediateVariable := registers[arg1]
			registers[arg1] = registers[arg2]
			registers[arg2] = intermediateVariable
		case CMP:
			var arg1 int = assemblerProgram[i][1]
			var arg2 int = assemblerProgram[i][2]
			var arg3 int = assemblerProgram[i][3]
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
		case JMP:
			if assemblerProgram[i][1] > 0 {
				i = i + assemblerProgram[i][1] - 1
			} else {
				i = i + assemblerProgram[i][1]
			}
		}
		fmt.Println(i, assemblerProgram[i], registers, stack)
	}
	fmt.Println(registers, stack)
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
