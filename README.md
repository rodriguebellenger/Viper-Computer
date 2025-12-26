Viper Computer is a virtual computer made for educationnal purposes.

To use it, enter (from the project folder) : 
"go run assembler/assembler.go file_name.vasm" (the .vasm extension is not mandatory)

It should work as expected if you have an up-to-date version of Go, since their are no dependency.

For how to write an actual program, please refer to the 
examples in assembler/assembly_test.

There are 16 registers of 64bits from R0 to R15.

The RAM has a size of a kilobyte (but can easily be changed with RAMSize variable).

The different operations possible are listed below. 

|   | 1byte  | 1byte  | 1byte  | 1byte |Additionnal info| Works |
|---|--------|--------|--------|-------|-| Yes |
000 | HLT    | EMPTY  | EMPTY  | EMPTY || Yes |
001 | AND    | R1     | R2     | EMPTY || Yes |
002 | ANDIB  | R1     | IMM    | EMPTY || Yes |
003 | ANDIW  | R1     | IMM    | IMM   || Yes |
004 | OR     | R1     | R2     | EMPTY || Yes |
005 | ORIB   | R1     | IMM    | EMPTY || Yes |
006 | ORIW   | R1     | IMM    | IMM   || Yes |
007 | NOT    | R1     | EMPTY  | EMPTY || Yes |
008 | SHIL   | R1     | R2     | EMPTY || Yes |
009 | SHILI  | R1     | IMM    | EMPTY || Yes |
010 | SHIR   | R1     | R2     | EMPTY || Yes |
011 | SHIRI  | R1     | IMM    | EMPTY || Yes |
012 | ADD    | R1     | R2     | EMPTY || Yes |
013 | ADDIB  | R1     | IMM    | EMPTY || Yes |
014 | ADDIW  | R1     | IMM    | IMM   || Yes |
015 | INCR   | R1     | EMPTY  | EMPTY || Yes |
016 | DECR   | R1     | EMPTY  | EMPTY || Yes |
017 | MUL    | R1     | R2     | EMPTY || No |
018 | MULIB  | R1     | IMM    | EMPTY || No |
019 | MULIW  | R1     | IMM    | IMM   || No |
020 | DIV    | R1     | R2     | EMPTY || No |
021 | DIVIB  | R1     | IMM    | EMPTY || No |
022 | DIVIW  | R1     | IMM    | IMM   || No |
023 | MOD    | R1     | R2     | EMPTY || No |
024 | MODIB  | R1     | IMM    | EMPTY || No |
025 | MODIW  | R1     | IMM    | IMM   || No |
026 | CLEAR  | R1     | EMPTY  | EMPTY || No |
027 | MOV1B  | R1     | IMM    | EMPTY | least significant byte | Yes |
028 | MOV2B  | R1     | IMM    | EMPTY || Yes |
029 | MOV3B  | R1     | IMM    | EMPTY || Yes |
030 | MOV4B  | R1     | IMM    | EMPTY | most significant byte | Yes |
031 | MOV1W  | R1     | IMM    | IMM   | least significant byte | Yes |
032 | MOV2W  | R1     | IMM    | IMM   || Yes |
033 | MOV3W  | R1     | IMM    | IMM   || Yes |
034 | MOV4W  | R1     | IMM    | IMM   | most significant byte | Yes |
035 | MOVR   | R1     | R2     | EMPTY || No |
036 | SWAP   | R1     | R2     | EMPTY || No |
037 | PUSH   | R1     | EMPTY  | EMPTY || No |
038 | PUSHIB | IMM    | EMPTY  | EMPTY || No |
039 | PUSHIW | IMM    | IMM    | EMPTY || No |
040 | PUSHIT | IMM    | IMM    | IMM   || No |
041 | POP    | R1     | EMPTY  | EMPTY || No |
042 | PEEK   | R1     | EMPTY  | EMPTY || No |
043 | CMP    | R1     | R2     | COMP_OP || Yes |
044 | JMP    | OFFSET | OFFSET | OFFSET | Jump to a label and continue execution from there | Yes |
045 | JMPB   | OFFSET | EMPTY  | EMPTY | Inserted automatically by the assembler | Yes |
046 | JMPW   | OFFSET | OFFSET | EMPTY | Inserted automatically by the assembler | Yes |
047 | JMPT   | OFFSET | OFFSET | OFFSET | Inserted automatically by the assembler | Yes |
048 | CALL   | OFFSET | OFFSET | OFFSET | Same as JMP, but push the current address before jumping | No |
049 | CALLB  | OFFSET | EMPTY  | EMPTY | Inserted automatically by the assembler | No |
050 | CALLW  | OFFSET | OFFSET | EMPTY | Inserted automatically by the assembler | No |
051 | CALLT  | OFFSET | OFFSET | OFFSET | Inserted automatically by the assembler | No |
052 | RET    | EMPTY  | EMPTY  | EMPTY | The execution continues at the address at the top of the stack | No |
053 | WRT    | SIZE   | *R1    | R2 || Yes |
054 | READ   | R1     | SIZE   | *R2 || Yes |
