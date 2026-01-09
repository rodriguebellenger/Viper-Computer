package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

//////////////
// COMMANDS //
//////////////

func runCommand(args []string) {
	var program string = readFile(args[0])
	args = args[1:]
	var debug bool = false
	var time_measurement uint64 = 1

	for i := 0; i < len(args); i++ {
		if args[i] == "debug" {
			debug = true
		} else if args[i] == "-time" {
			if !isInt(args[i+1]) {
				log.Fatal("-time needs a integer.")
			} else if args[i+1][0] == '-' {
				log.Fatal("-time needs a positive integer.")
			}
			time_measurement = uint64(strToInt(args[i+1]))
			i += 1
		} else {
			log.Fatal("Unrecognized argument for run command : " + args[i])
		}
	}

	var assemblerProgram [][]string = readProgram(program)

	var startTime time.Time = time.Now()
	var byteProgram []uint8 = programCleaner(assemblerProgram)
	var elapsed time.Duration = time.Since(startTime)
	if debug {
		fmt.Println(byteProgram)
		fmt.Printf("Time : %s\n\n", elapsed)
	}
	if time_measurement == 1 {
		startTime = time.Now()
		executeProgram()
		elapsed = time.Since(startTime)
		if debug {
			fmt.Printf("Time : %s\n", elapsed)
		}
	} else {
		var total_time time.Duration
		for i := 0; uint64(i) < time_measurement; i++ {
			startTime = time.Now()
			for i := 0; i < 1; i++ {
				executeProgram()
			}
			total_time += time.Since(startTime)
		}
		fmt.Printf("Time : %s\n", total_time/200)
		fmt.Printf("Total time : %s\n", total_time)
	}
}

func checkCommand(args []string) {
	if len(args) > 2 {
		log.Fatal("Too many argument for --check.")
	} else if args[0][len(args[0])-5:] != ".vasm" {
		log.Fatal("Unrecognized extension for \"" + args[0] + "\", need .vasm")
	}

	var program string = readFile(args[0])
	var assemblerProgram [][]string = readProgram(program)

	var startTime time.Time = time.Now()
	var byteProgram []uint8 = programCleaner(assemblerProgram)
	var elapsed time.Duration = time.Since(startTime)
	if len(args) > 1 && args[1] == "-debug" {
		fmt.Println(byteProgram)
		fmt.Printf("Time : %s\n\n", elapsed)
	}
}

func helpCommand(args []string) {
	if len(args) > 0 {
		log.Fatal("Too many arguments for --help.")
	}
	fmt.Println(`Commands:
  --run     Assemble a .vasm file and execute it with the Go implementation of the virtual machine
  --check   Check whether a .vasm file can be assembled
  --emit    Assemble a .vasm file and save bytecode into a new file
  --load    Load and execute an assembled bytecode file

Options:
  -debug        Enable debug output (only for --run and --check)
  -time <n>     Measure average execution time over <n> runs (--run only)
  -c-vm         Execute the file with the C implementation of the virtual machine (--load only)
  -go-vm        Execute the file with the Go implementation of the virtual machine (--load only)

Command usage:
  vasm --run   <file.vasm> [-time <n>] [-debug]
  vasm --check <file.vasm> [-debug]
  vasm --emit  <file.vasm> <output.vbc>
  vasm --load  <file.vbc> [-c-vm/-go-vm]`)
}

////////////////////
// COMMAND HELPER //
////////////////////

func writeToRAM(byteProgram []uint8) {
	for i, byte := range byteProgram {
		RAM[i] = byte
	}
}

func readProgram(program string) [][]string {
	var operations []string = strings.Split(program, "\n")
	var assemblerProgram [][]string

	for _, line := range operations {
		assemblerProgram = append(assemblerProgram, strings.Fields(line))
	}
	return assemblerProgram
}

func readFile(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("\rCouldn't read file : " + path)
	}
	var program string = string(content)
	return program
}
