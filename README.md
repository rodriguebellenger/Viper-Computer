# Viper Computer

Viper Computer is a virtual computer made for educationnal purposes.

## Features

- Custom assembly language
- Bytecode generation
- Custom bytecode interpreter (C and Go implementation)
- Debug and performance measurement tools

## How to use

It should work as expected if you have an up-to-date version of Go, since their are no dependency.

### Commands

If you want to assemble a .vasm file and execute it with the Go implementation of the virtual machine, use `--run`.  
You can add `-time <n>` to measure the average execution time.  
You can also add `-debug` to output the bytecodes and the assembling duration.  
```
go run path/to/assembler --run <file.vasm> [-time <n>] [-debug]
```

If you want to check whether a .vasm file can be assembled, use `--check`.  
You can add `-debug` to output the bytecodes and the assembling duration.
```
go run path/to/assembler.go --check <file.vasm> [-debug]  
``` 

If you want to assemble a .vasm file and save the bytecodes into a new file, use `--emit`.
```
go run path/to/assembler.go --emit <file.vasm> <output.vbc> 
```

If you want to load and execute a .vbc file (assembled bytecode file), use `--load`.  
You must use either use `-c-vm` or `-go-vm` to specify which version of the vm you want to use.
```
go run path/to/assembler.go --load <file.vbc> [-c-vm/-go-vm]
```  

## Syntax

The instructions are almost all of the form : `INST [arg1] [arg2]` (with whatever number of args needed) and are mostly self-explanatory.   
Special cases are :  
- `CMP [register] [register] [COMP_OP]` with COMP_OP being either E (equal), G (greater), L (less) or NE (not equal)
The next instruction is executed only if the comparison is true.  
- `READ [@Size] [*register] [register]` with @Size being either @8, @16, @24, @32, @40, @48, @56 or @64.  
The size indicates the number of bytes which will be read. *register will take the value of the register as an address and the value read will be stored in register.  
- `WRT [register] [@Size] [*register]` with @Size being either @8, @16, @24, @32, @40, @48, @56 or @64.  
Same as READ except the order of the arguments is changed to indicate that the value in the register will be stored in the RAM at the address within *register with size of @Size.  
- `JMP Label` continues the program directly after where the label was defined.  
- `CALL Label` same as JMP, except it pushes the current execution address onto the stack.  
- `RET` jumps to the address at the top of the stack.  

To create a label, enter `TheNameOfTheLabel:`. You can then refer to it via a JMP or a CALL simply by using its name without the ":".  
`JMP Label` or `CALL Label`

Please note that the compiler will automatically ignore every useless characters, (not alphanumerical characters and not in ":@*")  
Hence, this is fine `,A$D,D,     %R|1/     -R_2,` and will be changed to simply `ADD R1 R2`.  

## Architecture

There are 16 registers of 64bits from R0 to R15.  
The RAM has a size of a kilobyte (but can easily be changed with RAMSize variable).

## Operations

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
