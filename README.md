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
|---|--------|--------|--------|-------|-|-|
|000 | HLT    | EMPTY  | EMPTY  | EMPTY || Yes |
|001 | AND    | Register | Register | EMPTY || Yes |
|002 | ANDIB  | Register | IMM    | EMPTY || Yes |
|003 | ANDIW  | Register | IMM    | IMM   || Yes |
|004 | OR     | Register | Register | EMPTY || Yes |
|005 | ORIB   | Register | IMM    | EMPTY || Yes |
|006 | ORIW   | Register | IMM    | IMM   || Yes |
|007 | NOT    | Register | EMPTY  | EMPTY || Yes |
|008 | SHIL   | Register | Register | EMPTY || Yes |
|009 | SHILI  | Register | IMM    | EMPTY || Yes |
|010 | SHIR   | Register | Register | EMPTY || Yes |
|011 | SHIRI  | Register | IMM    | EMPTY || Yes |
|012 | ADD    | Register | Register | EMPTY || Yes |
|013 | ADDIB  | Register | IMM    | EMPTY || Yes |
|014 | ADDIW  | Register | IMM    | IMM   || Yes |
|015 | INCR   | Register | EMPTY  | EMPTY || Yes |
|016 | DECR   | Register | EMPTY  | EMPTY || Yes |
|017 | MUL    | Register | Register | EMPTY || No |
|018 | MULIB  | Register | IMM    | EMPTY || No |
|019 | MULIW  | Register | IMM    | IMM   || No |
|020 | DIV    | Register | Register | EMPTY || No |
|021 | DIVIB  | Register | IMM    | EMPTY || No |
|022 | DIVIW  | Register | IMM    | IMM   || No |
|023 | MOD    | Register | Register | EMPTY || No |
|024 | MODIB  | Register | IMM    | EMPTY || No |
|025 | MODIW  | Register | IMM    | IMM   || No |
|026 | CLEAR  | Register | EMPTY  | EMPTY || No |
|027 | MOV1B  | Register | IMM    | EMPTY | least significant byte | Yes |
|028 | MOV2B  | Register | IMM    | EMPTY || Yes |
|029 | MOV3B  | Register | IMM    | EMPTY || Yes |
|030 | MOV4B  | Register | IMM    | EMPTY | most significant byte | Yes |
|031 | MOV1W  | Register | IMM    | IMM   | least significant byte | Yes |
|032 | MOV2W  | Register | IMM    | IMM   || Yes |
|033 | MOV3W  | Register | IMM    | IMM   || Yes |
|034 | MOV4W  | Register | IMM    | IMM   | most significant byte | Yes |
|035 | MOVR   | Register | Register | EMPTY || No |
|036 | SWAP   | Register | Register | EMPTY || No |
|037 | PUSH   | Register | EMPTY  | EMPTY || No |
|038 | PUSHIB | IMM    | EMPTY  | EMPTY || No |
|039 | PUSHIW | IMM    | IMM    | EMPTY || No |
|040 | PUSHIT | IMM    | IMM    | IMM   || No |
|041 | POP    | Register | EMPTY  | EMPTY || No |
|042 | PEEK   | Register | EMPTY  | EMPTY || No |
|043 | CMP    | Register | Register | COMP_OP | The COMP_OP can be G, L, E, or NE (greater, less, equal or not equal) | Yes |
|044 | JMP    | OFFSET | OFFSET | OFFSET | Jump to a label and continue execution from there | Yes |
|045 | JMPB   | OFFSET | EMPTY  | EMPTY | Inserted automatically by the assembler | Yes |
|046 | JMPW   | OFFSET | OFFSET | EMPTY | Inserted automatically by the assembler | Yes |
|047 | JMPT   | OFFSET | OFFSET | OFFSET | Inserted automatically by the assembler | Yes |
|048 | CALL   | OFFSET | OFFSET | OFFSET | Same as JMP, but push the current address before jumping | No |
|049 | CALLB  | OFFSET | EMPTY  | EMPTY | Inserted automatically by the assembler | No |
|050 | CALLW  | OFFSET | OFFSET | EMPTY | Inserted automatically by the assembler | No |
|051 | CALLT  | OFFSET | OFFSET | OFFSET | Inserted automatically by the assembler | No |
|052 | RET    | EMPTY  | EMPTY  | EMPTY | The execution continues at the address at the top of the stack | No |
|053 | WRT    | SIZE   | *Register | Register || Yes |
|054 | READ   | Register | SIZE   | *Register || Yes |
