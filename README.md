Viper Computer is a virtual computer made for educationnal purposes.
To use it, enter (from the project folder) : "go run assembler/assembler.go file_name.vasm" (the .vasm extension is not mandatory)
It should work as you have an up-to-date version of Go, since their are no dependency.
The different operations possible are listed below. For how to write an actual program, please refer to the examples in
assembler/assembly_test.

|   | 1byte  | 1byte  | 1byte  | 1byte |
|---|--------|--------|--------|-------|
000 | HLT    | EMPTY  | EMPTY  | EMPTY |
001 | AND    | R1     | R2     | EMPTY |
002 | ANDIB  | R1     | IMM    | EMPTY |
003 | ANDIW  | R1     | IMM    | IMM   |
004 | OR     | R1     | R2     | EMPTY |
005 | ORIB   | R1     | IMM    | EMPTY |
006 | ORIW   | R1     | IMM    | IMM   |
007 | NOT    | R1     | EMPTY  | EMPTY |
008 | SHIL   | R1     | R2     | EMPTY |
009 | SHILI  | R1     | IMM    | EMPTY |
010 | SHIR   | R1     | R2     | EMPTY |
011 | SHIRI  | R1     | IMM    | EMPTY |
012 | ADD    | R1     | R2     | EMPTY |
013 | ADDIB  | R1     | IMM    | EMPTY |
014 | ADDIW  | R1     | IMM    | IMM   |
015 | INCR   | R1     | EMPTY  | EMPTY |
016 | DECR   | R1     | EMPTY  | EMPTY |
017 | MUL    | R1     | R2     | EMPTY |
018 | MULIB  | R1     | IMM    | EMPTY |
019 | MULIW  | R1     | IMM    | IMM   |
020 | DIV    | R1     | R2     | EMPTY |
021 | DIVIB  | R1     | IMM    | EMPTY |
022 | DIVIW  | R1     | IMM    | IMM   |
023 | MOD    | R1     | R2     | EMPTY |
024 | MODIB  | R1     | IMM    | EMPTY |
025 | MODIW  | R1     | IMM    | IMM   |
026 | CLEAR  | R1     | EMPTY  | EMPTY |
027 | MOV1B  | R1     | IMM    | EMPTY | (least significant byte)
028 | MOV2B  | R1     | IMM    | EMPTY |
029 | MOV3B  | R1     | IMM    | EMPTY |
030 | MOV4B  | R1     | IMM    | EMPTY | (most significant byte)
031 | MOV1W  | R1     | IMM    | IMM   | (least significant byte)
032 | MOV2W  | R1     | IMM    | IMM   |
033 | MOV3W  | R1     | IMM    | IMM   |
034 | MOV4W  | R1     | IMM    | IMM   | (most significant byte)
035 | MOVR   | R1     | R2     | EMPTY |
036 | SWAP   | R1     | R2     | EMPTY |
037 | PUSH   | R1     | EMPTY  | EMPTY |
038 | PUSHIB | IMM    | EMPTY  | EMPTY |
039 | PUSHIW | IMM    | IMM    | EMPTY |
040 | PUSHIT | IMM    | IMM    | IMM   |
041 | POP    | R1     | EMPTY  | EMPTY |
042 | PEEK   | R1     | EMPTY  | EMPTY |
043 | CMP    | R1     | R2     | COMP_OP |
044 | JMP    | OFFSET | OFFSET | OFFSET |
045 | JMPB   | OFFSET | EMPTY  | EMPTY | (done automatically by the assembler)
046 | JMPW   | OFFSET | OFFSET | EMPTY | (done automatically by the assembler)
047 | JMPT   | OFFSET | OFFSET | OFFSET | (done automatically by the assembler)
048 | CALL   | OFFSET | OFFSET | OFFSET |
049 | CALLB  | OFFSET | EMPTY  | EMPTY | (done automatically by the assembler)
050 | CALLW  | OFFSET | OFFSET | EMPTY | (done automatically by the assembler)
051 | CALLT  | OFFSET | OFFSET | OFFSET | (done automatically by the assembler)
052 | RET    | EMPTY  | EMPTY  | EMPTY |
053 | WRT    | SIZE   | *R1    | R2 |
054 | READ   | R1     | SIZE   | *R2 |
