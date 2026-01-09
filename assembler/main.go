package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var commands = map[string]func([]string){
	"--run":   runCommand,
	"--check": checkCommand,
	//"--emit":  testCommand,
	//"--load":  loadCommand,
	"--help": helpCommand,
}

//////////
// MAIN //
//////////

func main() {
	var args []string = os.Args[1:]
	command, ok := commands[args[0]]
	if !ok {
		log.Fatal("Unknown command, please enter --help to see the full command list.")
	}
	command(args[1:])
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
	return (strings.Contains("-0123456789", string(x[0])))
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
