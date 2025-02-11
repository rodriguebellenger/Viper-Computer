package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

//////////
// DATA //
//////////

type Token struct {
	value     string
	tokenType int
	column    int
	row       int
}

const (
	integer int = iota // Ex : 14 ; 42 ; ...

	identifier      // Ex : x ; variable1 ; ...
	reservedKeyword // Ex : true ; int ; ...

	size        // Ex : @8 ; @64 ; ...
	assignement // Ex : =

	EOL // Ex : End Of Line
)

var dataTypes []string = []string{"int"}

var listReservedKeywords []string = []string{"int"}

var symbols []string = []string{"="}

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
	fmt.Println(program)

	var parsedProgram []Token = tokenize(program)
	fmt.Println(parsedProgram)

	var variables = make(map[string]string)
	//variables = lexer(parsedProgram)
	fmt.Println(variables)
}

///////////////
// TOKENIZER //
///////////////

func tokenize(program string) []Token {
	var parsedProgram []Token
	var column int
	var row int
	var word string
	var token Token
	for i := range len(program) {
		if strings.Contains(" \n", string(program[i])) {
			if len(word) != 0 {
				token = createToken(word, row, column-len(word))
				parsedProgram = append(parsedProgram, token)
				word = ""
			}

			if string(program[i]) == "\n" {
				token.column = column
				token.row = row
				token.value = word
				token.tokenType = EOL
				column = 0
				row += 1
			}
		} else {
			word += string(program[i])
		}
		column += 1

	}
	return parsedProgram
}

func createToken(word string, row int, column int) Token {
	var token Token
	token.column = column
	token.row = row
	token.value = word
	token.tokenType = findType(word)
	return token
}

func findType(word string) int {
	if isInt(word) {
		return integer
	} else if word[0] == '@' {
		if isInt(word[1:]) && isPowerOfTwo(strToInt(word[1:])) {
			return size
		}
	} else if inList(symbols, word) {
		return findSymbolType(word)
	} else if inList(listReservedKeywords, word) {
		return reservedKeyword
	} else {
		return identifier
	}
	return -1
}

func findSymbolType(word string) int {
	if word == "=" {
		return assignement
	}
	return -1
}

///////////
// LEXER //
///////////

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
