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

var mnemonics []string = []string{"MOV", "ADDI", "ADD", "AND", "OR", "NOT", "PUSH", "POP", "SWAP", "CMP", "JMP", "RET"}
var registers []int = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var stack []int = []int{}

//////////
// MAIN //
//////////

func main() {
	content, err := os.ReadFile("assembler/src/program_test/sommePaire.txt")
	if err != nil {
		log.Fatal("Couldn't read file")
	}

	var program string = string(content)
	var assemblerProgram [][]string = readProgram(program)

	assemblerProgram = programCleaner(assemblerProgram)
	fmt.Println(assemblerProgram)

	var startTime time.Time = time.Now()
	finalProgram := assembleProgram(assemblerProgram)
	var elapsed time.Duration = time.Since(startTime)
	fmt.Printf("Temps : %s\n", elapsed)

	fmt.Println(finalProgram)
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

func programCleaner(assemblerProgram [][]string) [][]string {
	assemblerProgram = cleanEmptyOpe(assemblerProgram)
	checkUnexpectedCharacter(assemblerProgram)
	checkMnemonics(assemblerProgram)
	checkArgs(assemblerProgram)
	checkRegisters(assemblerProgram)
	assemblerProgram = countLines(assemblerProgram)

	// Jumps
	var labels map[string]int = checkJumps(assemblerProgram)
	assemblerProgram = createJumpAddress(assemblerProgram, labels)

	return assemblerProgram
}

func checkRegisters(assemblerProgram [][]string) {
	for i := range assemblerProgram {
		for j, arg := range assemblerProgram[i] {
			if arg[0] == 'R' {
				var reg int = strToInt(assemblerProgram[i][j][1:])
				if reg < 0 || reg > 15 {
					err := "Register out of range \"" + assemblerProgram[i][j] + "\" at line " + intToStr(i+1)
					log.Fatal(err)
				}
			}
		}
	}
}

func countLines(assemblerProgram [][]string) [][]string {
	for i := range len(assemblerProgram) {
		assemblerProgram[i] = append(assemblerProgram[i], intToStr(i+1))
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
				err := "Undefined label \"" + operation[1] + "\" at line " + operation[2]
				log.Fatal(err)
			}
			operation[1] = intToStr(targetLine - 1)
		}
	}
	return assemblerProgram
}

func checkJumps(assemblerProgram [][]string) map[string]int {
	var labels = make(map[string]int)
	for _, operation := range assemblerProgram {
		if string(operation[0][len(operation[0])-1]) == ":" {
			labels[operation[0][:len(operation[0])-1]] = strToInt(operation[1])
		}
	}
	return labels
}

func checkMnemonics(assemblerProgram [][]string) {
	for i, operation := range assemblerProgram {
		if !inList(mnemonics, operation[0]) && string(operation[0][len(operation[0])-1]) != ":" && len(operation) != 1 {
			err := "Unrecognized mnemonics \"" + operation[0] + "\" at line " + intToStr(i+1)
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
			err := "Wrong number of arguments for " + operation[0] + " at line " + operation[len(operation)-1]
			log.Fatal(err)
		} else if inList([]string{"NOT", "POP", "PUSH"}, operation[0]) && len(operation) != 2 {
			err := "Wrong numbers of arguments for " + operation[0] + " at line " + operation[len(operation)-1]
			log.Fatal(err)
		} else if string(operation[0][len(operation[0])-1]) == ":" && len(operation) != 2 {

		}
	}
}

/////////////////////////
// Execute the program //
/////////////////////////

func assembleProgram(assemblerProgram [][]string) string {
	var finalProgram string = ""
	for i := 0; i < len(assemblerProgram); i++ {
		if assemblerProgram[i][0] == "MOV" {
			var arg1 string = decimalToBinary(strToInt(assemblerProgram[i][1][1:]), 8)
			var arg2 string = decimalToBinary(strToInt(assemblerProgram[i][2]), 32)
			finalProgram = finalProgram + decimalToBinary(7, 8) + arg1 + arg2

		} else if assemblerProgram[i][0] == "ADD" {
			var arg1 string = decimalToBinary(strToInt(assemblerProgram[i][1][1:]), 4)
			var arg2 string = decimalToBinary(strToInt(assemblerProgram[i][2][1:]), 4)
			finalProgram = finalProgram + decimalToBinary(5, 8) + arg1 + arg2

		} else if assemblerProgram[i][0] == "ADDI" {
			var arg1 string = decimalToBinary(strToInt(assemblerProgram[i][1][1:]), 8)
			var arg2 string = decimalToBinary(strToInt(assemblerProgram[i][2]), 8)
			finalProgram = finalProgram + decimalToBinary(6, 8) + arg1 + arg2

		} else if assemblerProgram[i][0] == "PUSH" {
			var arg1 string = decimalToBinary(strToInt(assemblerProgram[i][1][1:]), 8)
			finalProgram = finalProgram + decimalToBinary(8, 8) + arg1

		} else if assemblerProgram[i][0] == "POP" {
			var arg1 string = decimalToBinary(strToInt(assemblerProgram[i][1][1:]), 8)
			finalProgram = finalProgram + decimalToBinary(9, 8) + arg1

		} else if assemblerProgram[i][0] == "AND" {
			var arg1 string = decimalToBinary(strToInt(assemblerProgram[i][1][1:]), 4)
			var arg2 string = decimalToBinary(strToInt(assemblerProgram[i][2][1:]), 4)
			finalProgram = finalProgram + decimalToBinary(2, 8) + arg1 + arg2

		} else if assemblerProgram[i][0] == "OR" {
			var arg1 string = decimalToBinary(strToInt(assemblerProgram[i][1][1:]), 4)
			var arg2 string = decimalToBinary(strToInt(assemblerProgram[i][2][1:]), 4)
			finalProgram = finalProgram + decimalToBinary(3, 8) + arg1 + arg2

		} else if assemblerProgram[i][0] == "NOT" {
			var arg1 string = decimalToBinary(strToInt(assemblerProgram[i][1][1:]), 8)
			finalProgram = finalProgram + decimalToBinary(4, 8) + arg1

		} else if assemblerProgram[i][0] == "SWAP" {
			var arg1 string = decimalToBinary(strToInt(assemblerProgram[i][1][1:]), 4)
			var arg2 string = decimalToBinary(strToInt(assemblerProgram[i][2][1:]), 4)
			finalProgram = finalProgram + decimalToBinary(14, 8) + arg1 + arg2

		} else if assemblerProgram[i][0] == "CMP" {
			var arg1 string = decimalToBinary(strToInt(assemblerProgram[i][1][1:]), 4)
			var arg2 string = decimalToBinary(strToInt(assemblerProgram[i][2][1:]), 4)
			var arg3 string = assemblerProgram[i][3]
			if inList([]string{"L", "G", "E"}, arg3) {
				if arg3 == "L" {
					finalProgram = finalProgram + decimalToBinary(6, 8) + arg1 + arg2 + decimalToBinary(2, 8)
				} else if arg3 == "G" {
					finalProgram = finalProgram + decimalToBinary(6, 8) + arg1 + arg2 + decimalToBinary(3, 8)
				} else if arg3 == "E" {
					finalProgram = finalProgram + decimalToBinary(6, 8) + arg1 + arg2 + decimalToBinary(1, 8)
				}
			} else {
				err := "Unrecognized comparison character \"" + arg3 + "\" at line " + assemblerProgram[i][4]
				log.Fatal(err)
			}
		} else if assemblerProgram[i][0] == "JMP" {
			var arg1 string = decimalToBinary(strToInt(assemblerProgram[i][1]), 16)
			finalProgram = finalProgram + decimalToBinary(11, 8) + arg1
		} else if assemblerProgram[i][0] == "RET" {
			i = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		} else {
			finalProgram += assemblerProgram[i][0]
		}
		//fmt.Println(i+1, registers, stack)
	}
	fmt.Println(registers, stack)
	return finalProgram
}

///////////
// UTILS //
///////////

func decimalToBinary(decimal int, size int) string {
	binary := ""
	for decimal > 0 {
		remainder := decimal % 2
		binary = fmt.Sprintf("%d%s", remainder, binary)
		decimal /= 2
	}
	for len(binary) < size {
		binary = "0" + binary
	}

	return binary
}

func strToInt(x string) int {
	num, err := strconv.Atoi(x)
	if err != nil {
		log.Fatal("Error in strToInt")
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
