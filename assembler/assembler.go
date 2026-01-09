package main

import (
	"fmt"
	"log"
	"strings"
)

var mnemonics []string = []string{"HLT", "AND", "ANDIB", "ANDIW", "OR", "ORIB", "ORIW", "NOT", "SHIL", "SHILI", "SHIR", "SHIRI", "ADD", "ADDIB", "ADDIW", "INCR", "DECR", "MUL", "MULIB", "MULIW", "DIV", "DIVIB", "DIVIW", "MOD", "MODIB", "MODIW", "CLEAR", "MOV1B", "MOV2B", "MOV3B", "MOV4B", "MOV1W", "MOV2W", "MOV3W", "MOV4W", "MOVR", "SWAP", "PUSH", "PUSHIB", "PUSHIW", "PUSHIT", "POP", "PEEK", "CMP", "JMP", "CALL", "RET", "WRT", "READ"}
var registersName []string = []string{"R0", "R1", "R2", "R3", "R4", "R5", "R6", "R7", "R8", "R9", "R10", "R11", "R12", "R13", "R14", "R15"}

var compileTimeBug []string

var opcodeToMnemonics = map[int]string{
	HLT: "HLT", AND: "AND", ANDIB: "ANDIB", ANDIW: "ANDIW", OR: "OR", ORIB: "ORIB", ORIW: "ORIW", NOT: "NOT", SHIL: "SHIL", SHILI: "SHILI", SHIR: "SHIR",
	SHIRI: "SHIRI", ADD: "ADD", ADDIB: "ADDIB", ADDIW: "ADDIW", INCR: "INCR", DECR: "DECR", MUL: "MUL", MULIB: "MULIB", MULIW: "MULIW", DIV: "DIV", DIVIB: "DIVIB", DIVIW: "DIVIW",
	MOD: "MOD", MODIB: "MODIB", MODIW: "MODIW", CLEAR: "CLEAR", MOV1B: "MOV1B", MOV2B: "MOV2B", MOV3B: "MOV3B", MOV4B: "MOV4B", MOV1W: "MOV1W", MOV2W: "MOV2W",
	MOV3W: "MOV3W", MOV4W: "MOV4W", MOVR: "MOVR", SWAP: "SWAP", PUSH: "PUSH", PUSHIB: "PUSHIB", PUSHIW: "PUSHIW", PUSHIT: "PUSHIT", POP: "POP", PEEK: "PEEK", CMP: "CMP",
	JMP: "JMP", JMPB: "JMPB", JMPW: "JMPW", JMPT: "JMPT", CALL: "CALL", CALLB: "CALLB", CALLW: "CALLW", CALLT: "CALLT", RET: "RET", WRT: "WRT", READ: "READ",
}

var mnemonicToOpcode = map[string]int{
	"HLT": HLT, "AND": AND, "ANDIB": ANDIB, "ANDIW": ANDIW, "OR": OR, "ORIB": ORIB, "ORIW": ORIW, "NOT": NOT, "SHIL": SHIL, "SHILI": SHILI, "SHIR": SHIR,
	"SHIRI": SHIRI, "ADD": ADD, "ADDIB": ADDIB, "ADDIW": ADDIW, "INCR": INCR, "DECR": DECR, "MUL": MUL, "MULIB": MULIB, "MULIW": MULIW, "DIV": DIV, "DIVIB": DIVIB, "DIVIW": DIVIW,
	"MOD": MOD, "MODIB": MODIB, "MODIW": MODIW, "CLEAR": CLEAR, "MOV1B": MOV1B, "MOV2B": MOV2B, "MOV3B": MOV3B, "MOV4B": MOV4B,
	"MOV1W": MOV1W, "MOV2W": MOV2W, "MOV3W": MOV3W, "MOV4W": MOV4W, "MOVR": MOVR, "SWAP": SWAP, "PUSH": PUSH, "PUSHIB": PUSHIB, "PUSHIW": PUSHIW, "PUSHIT": PUSHIT,
	"POP": POP, "PEEK": PEEK, "CMP": CMP, "JMP": JMP, "JMPB": JMPB, "JMPW": JMPW, "JMPT": JMPT, "CALL": CALL, "CALLB": CALLB, "CALLW": CALLW, "CALLT": CALLT, "RET": RET, "WRT": WRT, "READ": READ,
}

var comparOpToOpcode = map[string]string{
	"L": "1", "G": "2", "E": "3", "NE": "4",
}

var syntaxRules = map[string][]string{
	"HLT":    {},
	"AND":    {"Register", "Register"},
	"ANDIB":  {"Register", "Int8"},
	"ANDIW":  {"Register", "Int16"},
	"OR":     {"Register", "Register"},
	"ORIB":   {"Register", "Int8"},
	"ORIW":   {"Register", "Int16"},
	"NOT":    {"Register"},
	"SHIL":   {"Register", "Register"},
	"SHILI":  {"Register", "Int8"},
	"SHIR":   {"Register", "Register"},
	"SHIRI":  {"Register", "Int8"},
	"ADD":    {"Register", "Register"},
	"ADDIB":  {"Register", "Int8"},
	"ADDIW":  {"Register", "Int16"},
	"INCR":   {"Register"},
	"DECR":   {"Register"},
	"MUL":    {"Register", "Register"},
	"MULIB":  {"Register", "Int8"},
	"MULIW":  {"Register", "Int16"},
	"DIV":    {"Register", "Register"},
	"DIVIB":  {"Register", "Int8"},
	"DIVIW":  {"Register", "Int16"},
	"MOD":    {"Register", "Register"},
	"MODIB":  {"Register", "Int8"},
	"MODIW":  {"Register", "Int16"},
	"CLEAR":  {"Register"},
	"MOV1B":  {"Register", "Int8"},
	"MOV2B":  {"Register", "Int8"},
	"MOV3B":  {"Register", "Int8"},
	"MOV4B":  {"Register", "Int8"},
	"MOV1W":  {"Register", "Int16"},
	"MOV2W":  {"Register", "Int16"},
	"MOV3W":  {"Register", "Int16"},
	"MOV4W":  {"Register", "Int16"},
	"MOVR":   {"Register", "Register"},
	"SWAP":   {"Register", "Register"},
	"PUSH":   {"Register"},
	"PUSHIB": {"Int8"},
	"PUSHIW": {"Int16"},
	"PUSHIT": {"Int24"},
	"POP":    {"Register"},
	"PEEK":   {"Register"},
	"CMP":    {"Register", "Register", "Comparison"},
	"JMP":    {"Offset"},
	"JMPB":   {"Int8"},
	"JMPW":   {"Int16"},
	"JMPT":   {"Int24"},
	"CALL":   {"Offset"},
	"CALLB":  {"Int8"},
	"CALLW":  {"Int16"},
	"CALLT":  {"Int24"},
	"RET":    {},
	"WRT":    {"Size", "Address", "Register"},
	"READ":   {"Register", "Size", "Address"},
}

var forbiddenLabels []string = []string{"R0", "R1", "R2", "R3", "R4", "R5", "R6", "R7", "R8", "R9", "R10", "R11", "R12", "R13", "R14", "R15",
	"r0", "r1", "r2", "r3", "r4", "r5", "r6", "r7", "r8", "r9", "r10", "r11", "r12", "r13", "r14", "r15",
	"HLT", "AND", "ANDIB", "ANDIW", "OR", "ORIB", "ORIW", "NOT", "SHIL", "SHILI", "SHIR", "SHIRI", "ADD", "ADDIB", "ADDIW", "INCR", "DECR",
	"MUL", "MULIB", "MULIW", "DIV", "DIVIB", "DIVIW", "MOD", "MODIB", "MODIW", "CLEAR", "MOV1B", "MOV2B", "MOV3B", "MOV4B", "MOV1W", "MOV2W", "MOV3W", "MOV4W",
	"MOVR", "SWAP", "PUSH", "PUSHIB", "PUSHIW", "PUSHIT", "POP", "PEEK", "CMP", "JMP", "JMPB", "JMPW", "JMPT", "CALL", "CALLB", "CALLW", "CALLT", "RET", "WRT", "READ",
	"E", "G", "L", "NE"}

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
		memoryAddress, labels = checkJumps(tokenizedProgram[i], labels, memoryAddress)
		checkSyntax(tokenizedProgram[i], syntaxRules[tokenizedProgram[i][0][0]])
		memoryAddress += 4
	}
	tokenizedProgram = delLabels(tokenizedProgram)

	memoryAddress = 0
	for i, line := range tokenizedProgram {
		if line[0][0] == "JMP" || line[0][0] == "CALL" {
			tokenizedProgram[i] = createJumpAddress(labels, line, memoryAddress)
		}
		memoryAddress += 4
	}
	tokenizedProgram = optimizeJumps(tokenizedProgram)

	if len(compileTimeBug) != 0 {
		for _, err := range compileTimeBug {
			fmt.Println(err)
		}
		log.Fatal("Couldn't compile")
	} else {
		fmt.Println("No compile error")
	}

	var opcodeProgram [][]uint32
	for _, line := range tokenizedProgram {
		opcodeProgram = append(opcodeProgram, mnemonicsToOpcode(line))
	}

	var byteProgram []uint8 = bytificationOfTheProgram(opcodeProgram)

	return byteProgram
}

//////////////////
// CLEAN HELPER //
//////////////////

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
		compileTimeBug = append(compileTimeBug, "Wrong number of args for \""+line[0]+"\" at line "+intToStr(i+1))
	}
}

func checkWords(line []string, i int) [][]string {
	var newLine [][]string
	for j, word := range line {
		if inList(mnemonics, word) {
			newLine = append(newLine, []string{word, "Operation"})
		} else if inList([]string{"G", "L", "E", "NE"}, word) {
			newLine = append(newLine, []string{comparOpToOpcode[word], "Comparison"})
		} else if word[len(word)-1] == ':' || (j > 0 && (line[j-1] == "JMP" || line[j-1] == "CALL")) {
			newLine = append(newLine, []string{word, "Offset"})
		} else if inList(registersName, word) {
			newLine = append(newLine, []string{word[1:], "Register"})
		} else if word[0] == '@' && isInt(word[1:]) && isPowerOfTwo(strToInt(word[1:])) && strToInt(word[1:]) >= 8 && strToInt(word[1:]) <= 64 {
			newLine = append(newLine, []string{intToStr(strToInt(word[1:]) / 8), "Size"})
		} else if word[0] == '*' && inList(registersName, word[1:]) {
			newLine = append(newLine, []string{word[2:], "Address"})
		} else if isInt(word) {
			var number uint64 = uint64(strToInt(word))
			if inList([]string{"ANDIB", "ORIB", "SHILI", "SHIRI", "ADDIB", "MULIB", "DIVIB", "MODIB", "MOV1B", "MOV2B", "MOV3B", "MOV4B", "PUSHIB"}, line[0]) && number < 256 {
				newLine = append(newLine, []string{word, "Int8"})
			} else if inList([]string{"ANDIW", "ORIW", "ADDIW", "MULIW", "DIVIW", "MODIW", "MOV1W", "MOV2W", "MOV3W", "MOV4W", "PUSHIW"}, line[0]) && number < 65536 {
				newLine = append(newLine, []string{word, "Int16"})
			} else if line[0] == "PUSHIT" && number < 16777216 {
				newLine = append(newLine, []string{word, "Int24"})
			} else if inList([]string{"ANDIB", "ORIB", "SHILI", "SHIRI", "ADDIB", "MULIB", "DIVIB", "MODIB", "MOV1B", "MOV2B", "MOV3B", "MOV4B", "PUSHIB", "ANDIW", "ORIW", "ADDIW", "MULIW", "DIVIW", "MODIW", "MOV1W", "MOV2W", "MOV3W", "MOV4W", "PUSHIW", "PUSHIT"}, line[0]) {
				compileTimeBug = append(compileTimeBug, "Immediate \""+word+"\" is too big at line "+intToStr(i+1))
			}
		} else {
			compileTimeBug = append(compileTimeBug, "Unrecognized token \""+word+"\" at line "+intToStr(i+1))
		}
	}
	newLine = append(newLine, []string{intToStr(i), "Line"})
	return newLine
}

func checkJumps(line [][]string, labels map[string]int, memoryAddress int) (int, map[string]int) {
	if line[0][0][len(line[0][0])-1] == ':' {
		if !(inList(forbiddenLabels, line[0][0][:len(line[0][0])-1])) {
			labels[string(line[0][0][:len(line[0][0])-1])] = memoryAddress - 1
			memoryAddress -= 4
		} else {
			compileTimeBug = append(compileTimeBug, "Forbiddent label name \""+string(line[0][0][:len(line[0][0])-1])+"\" at line "+intToStr(strToInt(string(line[1][0]))+1))
		}
	}
	return memoryAddress, labels
}

func checkSyntax(line [][]string, rules []string) {
	var errorSyntax bool = false
	for j := 0; j < len(rules) && j < len(line); j++ {
		if rules[j] != line[j+1][1] {
			errorSyntax = true
		}
	}
	if errorSyntax {
		compileTimeBug = append(compileTimeBug, "Syntax error at line "+intToStr(strToInt(line[len(line)-1][0])+1))
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
		compileTimeBug = append(compileTimeBug, "Undefined label \""+line[1][0]+"\" at line "+intToStr(strToInt(line[len(line)-1][0])+1))
	}
	var offset int = targetLine - memoryAdress
	if offset > 0 {
		offset -= 1
	}
	line[1][0] = intToStr(offset)
	return line
}

func optimizeJumps(tokenizedProgram [][][]string) [][][]string {
	var optimizedProgram [][][]string
	for _, line := range tokenizedProgram {
		if line[0][0] == "JMP" || line[0][0] == "CALL" {
			if strToInt(line[1][0]) <= 127 && strToInt(line[1][0]) >= -128 {
				line[0][0] = line[0][0] + "B"
				optimizedProgram = append(optimizedProgram, line)
			} else if strToInt(line[1][0]) <= 32767 && strToInt(line[1][0]) >= -32768 {
				line[0][0] = line[0][0] + "W"
				optimizedProgram = append(optimizedProgram, line)
			} else {
				line[0][0] = line[0][0] + "T"
				optimizedProgram = append(optimizedProgram, line)
			}
		} else {
			optimizedProgram = append(optimizedProgram, line)
		}
	}
	return optimizedProgram
}

func mnemonicsToOpcode(line [][]string) []uint32 {
	var newLine []uint32
	if string(line[0][0]) == "HLT" || string(line[0][0]) == "RET" {
		newLine = []uint32{uint32(mnemonicToOpcode[line[0][0]])}
	} else if inList([]string{"NOT", "INCR", "DECR", "CLEAR", "PUSH", "PUSHIB", "PUSHIW", "PUSHIT", "POP", "PEEK", "JMPB", "JMPW", "JMPT", "CALLB", "CALLW", "CALLT"}, string(line[0][0])) {
		var arg1 uint32 = uint32(strToInt(line[1][0]))
		newLine = []uint32{uint32(mnemonicToOpcode[line[0][0]]), arg1}
	} else if inList([]string{"AND", "ANDIB", "ANDIW", "OR", "ORIB", "ORIW", "SHIL", "SHILI", "SHIR", "SHIRI", "ADD", "ADDIB", "ADDIW", "MUL", "MULIB", "MULIW", "DIV", "DIVIB", "DIVIW", "MOD", "MODIB", "MODIW", "MOV1B", "MOV2B", "MOV3B", "MOV4B", "MOV1W", "MOV2W", "MOV3W", "MOV4W", "MOVR", "SWAP"}, string(line[0][0])) {
		var arg1 uint32 = uint32(strToInt(line[1][0]))
		var arg2 uint32 = uint32(strToInt(line[2][0]))
		newLine = []uint32{uint32(mnemonicToOpcode[line[0][0]]), arg1, arg2}
	} else if inList([]string{"CMP", "WRT", "READ"}, string(line[0][0])) {
		var arg1 uint32 = uint32(strToInt(line[1][0]))
		var arg2 uint32 = uint32(strToInt(line[2][0]))
		var arg3 uint32 = uint32(strToInt(line[3][0]))
		newLine = []uint32{uint32(mnemonicToOpcode[line[0][0]]), arg1, arg2, arg3}
	}
	return newLine
}

func bytificationOfTheProgram(opcodeProgram [][]uint32) []uint8 {
	var byteProgram []uint8
	for _, line := range opcodeProgram {
		switch line[0] {
		case uint32(HLT), uint32(RET):
			byteProgram = append(byteProgram, uint8(line[0]))
			byteProgram = append(byteProgram, 0)
			byteProgram = append(byteProgram, 0)
			byteProgram = append(byteProgram, 0)
		case uint32(NOT), uint32(INCR), uint32(DECR), uint32(CLEAR), uint32(PUSH), uint32(PUSHIB), uint32(POP), uint32(PEEK), uint32(JMPB), uint32(CALLB):
			byteProgram = append(byteProgram, uint8(line[0]))
			byteProgram = append(byteProgram, uint8(line[1]))
			byteProgram = append(byteProgram, 0)
			byteProgram = append(byteProgram, 0)
		case uint32(AND), uint32(ANDIB), uint32(OR), uint32(ORIB), uint32(SHIL), uint32(SHILI), uint32(SHIR), uint32(SHIRI), uint32(ADD), uint32(ADDIB), uint32(MUL), uint32(MULIB), uint32(DIV), uint32(DIVIB), uint32(MOD), uint32(MODIB), uint32(MOV1B), uint32(MOV2B), uint32(MOV3B), uint32(MOV4B), uint32(MOVR), uint32(SWAP):
			byteProgram = append(byteProgram, uint8(line[0]))
			byteProgram = append(byteProgram, uint8(line[1]))
			byteProgram = append(byteProgram, uint8(line[2]))
			byteProgram = append(byteProgram, 0)
		case uint32(PUSHIW), uint32(JMPW), uint32(CALLW):
			byteProgram = append(byteProgram, uint8(line[0]))
			byteProgram = append(byteProgram, uint8(line[1]))
			byteProgram = append(byteProgram, uint8(line[1]>>8))
			byteProgram = append(byteProgram, 0)
		case uint32(ANDIW), uint32(ORIW), uint32(ADDIW), uint32(MULIW), uint32(DIVIW), uint32(MODIW), uint32(MOV1W), uint32(MOV2W), uint32(MOV3W), uint32(MOV4W):
			byteProgram = append(byteProgram, uint8(line[0]))
			byteProgram = append(byteProgram, uint8(line[1]))
			byteProgram = append(byteProgram, uint8(line[2]))
			byteProgram = append(byteProgram, uint8(line[2]>>8))
		case uint32(PUSHIT), uint32(JMPT), uint32(CALLT):
			byteProgram = append(byteProgram, uint8(line[0]))
			byteProgram = append(byteProgram, uint8(line[1]))
			byteProgram = append(byteProgram, uint8(line[1]>>8))
			byteProgram = append(byteProgram, uint8(line[1]>>16))
		case uint32(CMP), uint32(WRT), uint32(READ):
			byteProgram = append(byteProgram, uint8(line[0]))
			byteProgram = append(byteProgram, uint8(line[1]))
			byteProgram = append(byteProgram, uint8(line[2]))
			byteProgram = append(byteProgram, uint8(line[3]))
		}
	}
	return byteProgram
}
