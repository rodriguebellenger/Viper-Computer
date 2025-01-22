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

//////////
// MAIN //
//////////

func main() {
	content, err := os.ReadFile("assembler/debugging/program_test/loops.txt")
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
	assemblerProgram = cleanEmptyOpe(assemblerProgram)
	checkUnexpectedCharacter(assemblerProgram)
	checkMnemonics(assemblerProgram)
	checkArgs(assemblerProgram)
	checkRegisters(assemblerProgram)
	assemblerProgram = countLines(assemblerProgram)

	// Jumps
	var labels map[string]int
	labels, assemblerProgram = checkJumps(assemblerProgram)
	assemblerProgram = createJumpAddress(assemblerProgram, labels)
	fmt.Println((assemblerProgram))
	var opcodeProgram [][]int = mnemonicsToOpcode(assemblerProgram)

	return opcodeProgram
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
		} else {
			log.Fatal("Err in mnemonicsToOpcode")
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

func checkRegisters(assemblerProgram [][]string) {
	for i := range assemblerProgram {
		for j, arg := range assemblerProgram[i] {
			if arg[0] == 'R' {
				var reg int = strToInt(assemblerProgram[i][j][1:])
				if reg < 0 || reg > 15 {
					err := "\rRegister out of range \"" + assemblerProgram[i][j] + "\" at line " + intToStr(i+1)
					log.Fatal(err)
				}
			}
		}
	}
}

func countLines(assemblerProgram [][]string) [][]string {
	for i := range len(assemblerProgram) {
		assemblerProgram[i] = append(assemblerProgram[i], intToStr(i))
	}
	return assemblerProgram
}

func cleanEmptyOpe(assemblerProgram [][]string) [][]string {
	var cleanedProgram [][]string
	newLine := 0
	for i := range len(assemblerProgram) {
		if len(assemblerProgram[i][0]) != 0 {
			cleanedProgram = append(cleanedProgram, []string{})
			for j := range len(assemblerProgram[i]) {
				if len(assemblerProgram[i][j]) != 0 {
					cleanedProgram[i-newLine] = append(cleanedProgram[i-newLine], assemblerProgram[i][j])
				}
			}
		} else {
			newLine += 1
		}
	}
	return cleanedProgram
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

func checkJumps(assemblerProgram [][]string) (map[string]int, [][]string) {
	var labels = make(map[string]int)
	for _, operation := range assemblerProgram {
		if string(operation[0][len(operation[0])-1]) == ":" {
			labels[string(operation[0][:len(operation[0])-1])] = strToInt(operation[1]) - 1
		}
	}
	return labels, assemblerProgram
}

func checkMnemonics(assemblerProgram [][]string) {
	for i, operation := range assemblerProgram {
		if !inList(mnemonics, operation[0]) && string(operation[0][len(operation[0])-1]) != ":" && len(operation) != 1 {
			err := "\rUnrecognized mnemonics \"" + operation[0] + "\" at line " + intToStr(i+1)
			log.Fatal(err)
		}
	}
}

func checkUnexpectedCharacter(assemblerProgram [][]string) {
	validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890<>=:-"
	for i := 0; i < len(assemblerProgram); i++ {
		for j := 0; j < len(assemblerProgram[i]); j++ {
			cleanedString := ""
			for k := 0; k < len(assemblerProgram[i][j]); k++ {
				char := string(assemblerProgram[i][j][k])
				if strings.Contains(validChars, char) {
					cleanedString += char
				}
			}
			assemblerProgram[i][j] = cleanedString
		}
	}
}

func checkArgs(assemblerProgram [][]string) {
	for _, operation := range assemblerProgram {
		if inList([]string{"OR", "AND", "MOV", "ADD", "ADDI", "SWAP"}, operation[0]) && len(operation) != 3 {
			err := "\rWrong number of arguments for " + operation[0] + " at line " + operation[len(operation)-1]
			log.Fatal(err)
		} else if inList([]string{"NOT", "POP", "PUSH"}, operation[0]) && len(operation) != 2 {
			err := "\rWrong numbers of arguments for " + operation[0] + " at line " + operation[len(operation)-1]
			log.Fatal(err)
		} else if string(operation[0][len(operation[0])-1]) == ":" && len(operation) != 2 {

		}
	}
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
				if registers[arg1]^registers[arg2] != 0 {
					i += 1
				}
			}
		case JMP:
			i = assemblerProgram[i][1]
		}
		//fmt.Println(i, registers, stack)
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
