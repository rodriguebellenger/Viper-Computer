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
	"HLT":  []string{},
	"RET":  []string{},
	"AND":  []string{"Register", "Register"},
	"OR":   []string{"Register", "Register"},
	"NOT":  []string{"Register"},
	"ADD":  []string{"Register", "Register"},
	"ADDI": []string{"Register", "Number"},
	"MOV":  []string{"Register", "Number"},
	"PUSH": []string{"Register"},
	"POP":  []string{"Register"},
	"CMP":  []string{"Register", "Register", "Comparison"},
	"JMP":  []string{"Label"},
	"WRT":  []string{},
	"READ": []string{},
	"SWAP": []string{"Register", "Register"},
}

//////////
// MAIN //
//////////

func main() {
	args := os.Args[1:] // Skip the program name
	fmt.Println(args)
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
	var assemblerProgram [][]string = [][]string{}
	for i := range len(operations) {
		assemblerProgram = append(assemblerProgram, strings.Split(operations[i], " "))
	}
	return assemblerProgram
}

///////////////////////
// Clean the program //
///////////////////////

func programCleaner(assemblerProgram [][]string) [][]int {
	var labels map[string]int
	var tokenizedProgram [][][]string
	for i, line := range assemblerProgram {
		var skipOrNot bool = skipEmptyLine(line)
		if !(skipOrNot) {
			line = checkUnexpectedCharacter(line)
			assemblerProgram[i] = delEmptyArgs(line)
			checkNumberOfArgs(line, i)
			tokenizedProgram = append(tokenizedProgram, checkWords(line, i))
			checkSyntax(tokenizedProgram[i], syntaxRules[tokenizedProgram[i][0][0]], i)
			tokenizedProgram, labels = checkJumps(line, labels)
		}
	}

	// Jumps
	assemblerProgram = createJumpAddress(assemblerProgram, labels)
	fmt.Println((assemblerProgram))
	//var opcodeProgram [][]int = mnemonicsToOpcode(assemblerProgram)

	return opcodeProgram
}

func skipEmptyLine(line []string) bool {
	return len(line) == 0
}

// check change of account
func delEmptyArgs(line []string) []string {
	for i := range line {
		if len(line[i]) == 0 {
			line = append(line[:i], line[i+1:]...)
		}
	}
	return line
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
		err := "\rWrong number of args for \"" + line[0] + "\" at line " + intToStr(i+1)
		log.Fatal(err)
	}
}

func checkWords(line []string, i int) [][]string {
	var newLine [][]string
	for _, word := range line {
		if inList(mnemonics, word) {
			newLine = append(newLine, []string{word, "Operation"})
		} else if strToInt(word[1:]) < 16 && strToInt(word[1:]) >= 0 && inList([]string{"r", "R"}, string(word[0])) {
			newLine = append(newLine, []string{word[1:], "Register"})
		} else if inList([]string{"G", "L", "E"}, word) {
			newLine = append(newLine, []string{word, "Comparison"})
		} else if word[len(word)-1] == ':' {
			newLine = append(newLine, []string{word, "Label"})
		} else {
			for _, character := range word {
				if !(strings.Contains("0123456789", string(character))) {
					err := "\rUnrecognized word \"" + word + "\" at line " + intToStr(i+1)
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
		err := "\rArgs don't respect syntax rule for \"" + line[0][0] + "\" at line " + intToStr(i+1)
		log.Fatal(err)
	}
}

func checkJumps(line [][]string, labels map[string]int) (map[string]int, [][]string) {
	if string(line[0][0][len(line[0])-1]) == ":" {
		labels[line[0][0][:len(line[0])-1]] = strToInt(operation[1]) - 1
	}
	return labels, assemblerProgram
}

func mnemonicsToOpcode(assemblerProgram [][]string) [][]int {
	var opcodeProgram [][]int
	var newLine []int
	var load int = 0
	for _, line := range assemblerProgram {
		if line[0] == "MOV" {
			var arg1 int = strToInt(line[1][1:])
			var arg2 int = strToInt(line[2])
			newLine = []int{MOV, arg1, arg2}

		} else if line[0] == "ADD" {
			var arg1 int = strToInt(line[1][1:])
			var arg2 int = strToInt(line[2][1:])
			newLine = []int{ADD, arg1, arg2}

		} else if line[0] == "ADDI" {
			var arg1 int = strToInt(line[1][1:])
			var arg2 int = strToInt(line[2])
			newLine = []int{ADDI, arg1, arg2}

		} else if line[0] == "PUSH" {
			var arg1 int = strToInt(line[1][1:])
			newLine = []int{PUSH, arg1}

		} else if line[0] == "POP" {
			var arg1 int = strToInt(line[1][1:])
			newLine = []int{POP, arg1}

		} else if line[0] == "AND" {
			var arg1 int = strToInt(line[1][1:])
			var arg2 int = strToInt(line[2][1:])
			newLine = []int{AND, arg1, arg2}

		} else if line[0] == "OR" {
			var arg1 int = strToInt(line[1][1:])
			var arg2 int = strToInt(line[2][1:])
			newLine = []int{OR, arg1, arg2}

		} else if line[0] == "NOT" {
			var arg1 int = strToInt(line[1][1:])
			newLine = []int{NOT, arg1}

		} else if line[0] == "SWAP" {
			var arg1 int = strToInt(line[1][1:])
			var arg2 int = strToInt(line[2][1:])
			newLine = []int{SWAP, arg1, arg2}

		} else if line[0] == "CMP" {
			var arg1 int = strToInt(line[1][1:])
			var arg2 int = strToInt(line[2][1:])
			var arg3 string = line[3]
			if inList([]string{"L", "G", "E"}, arg3) {
				if arg3 == "L" {
					newLine = []int{CMP, arg1, arg2, 1}
				} else if arg3 == "G" {
					newLine = []int{CMP, arg1, arg2, 2}
				} else if arg3 == "E" {
					newLine = []int{CMP, arg1, arg2, 3}
				}
			} else {
				err := "Unrecognized comparison character \"" + arg3 + "\" at line " + line[4]
				log.Fatal(err)
			}
		} else if line[0] == "JMP" {
			var arg1 int = strToInt(line[1])
			newLine = []int{JMP, arg1}
		} else if line[0] == "RET" {
			newLine = []int{RET}
		} else if string(line[0][len(line[0])-1]) == ":" {
			load = 1
		} else if line[0] == "HLT" {
			newLine = []int{HLT}
		} else {
			log.Fatal("Err in mnemonicsToOpcode : " + line[0])
		}
		switch load {
		case 0:
			opcodeProgram = append(opcodeProgram, newLine)
		case 1:
			load = 0
		}
	}
	return opcodeProgram
}

func createJumpAddress(assemblerProgram [][]string, labels map[string]int) [][]string {
	for _, operation := range assemblerProgram {
		if operation[0] == "JMP" {
			var targetLine int = labels[operation[1]]
			if targetLine == 0 {
				err := "\rUndefined label \"" + operation[1] + "\" at line " + operation[2]
				log.Fatal(err)
			}
			operation[1] = intToStr(targetLine)
		}
	}
	return assemblerProgram
}

/////////////////////////
// Execute the program //
/////////////////////////

func executeProgram(assemblerProgram [][]int) {
	var name string
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
				if registers[arg1]^registers[arg2] != 0 {
					i += 1
				}
			}
		case JMP:
			i = assemblerProgram[i][1]
		}
		fmt.Println(i, assemblerProgram[i], registers, stack)
		fmt.Scan(&name)
	}
	fmt.Println(registers, stack)
}

///////////
// UTILS //
///////////

func strToInt(x string) int {
	num, err := strconv.Atoi(x)
	if err != nil {
		log.Fatal("\rError in strToInt")
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
