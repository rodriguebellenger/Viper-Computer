# Viper computer

The viper computer is a virtual computer created to learn what programmation is at the lowest level possible.

# Assembly

For the moment, you have the assembly interpreter.

## Operations

| Name       | Arg1  | Arg2  | Arg3  | Effect                                                                                                   |
| ---------- | ----- | ----- | ----- | -------------------------------------------------------------------------------------------------------- |
| **HLT**    | None  | None  | None  | Stops the program no matter what.                                                                        |
| **CALL**   | Label | None  | None  | Jump to the address of the label, and push the current address onto the stack.                           |
| **RET**    | None  | None  | None  | Return to the address at the top of the stack.                                                           |
| **MOV**    | R1    | Immed | None  | Move the immediate to the register specified. The immediate is 64 bits.                                  |
| **ADD**    | R1    | R2    | None  | Add the content of two registers and put the result back in the first one.                               |
| **ADDI**   | R1    | R2    | None  | Add an immediate to a register (immediate is 8 bits).                                                    |
| **PUSH**   | R1    | None  | None  | Push the content of a register onto the stack.                                                           |
| **POP**    | R1    | None  | None  | Pop the top of the stack and put it in a register.                                                       |
| **AND**    | R1    | R2    | None  | Operate a bitwise and between two registers and put the result back in the first one.                    |
| **OR**     | R1    | R2    | None  | Operate a bitwise or between two registers and put the result back in the first one.                     |
| **NOT**    | R1    | None  | None  | Operate a bitwise not on a register.                                                                     |
| **SWAP**   | R1    | R2    | None  | Swap the content of two registers.                                                                       |
| **CMP**    | R1    | R2    | G/L/E | Compare two registers with comparison operator and skip the next instruction if the comparison is false. |
| **JMP**    | Label | None  | None  | Jump inconditionally to a label.                                                                         |
| **WRT**    | @32   | *R1   | R2    | Write the value of a register at an address contained in an other register, with the size indicated.     |
| **READ**   | R1    | @32   | *R2   | Read the value at an address contained in a register, with the size indicated, and put it in a register. |
| **Label:** | None  | None  | None  | Create a label that you can jump to.                                                                     |

## Registers

There are 16 registers, R0 through R15.

Use cases : 
R0 is a constant and always is 0 (nothing in the assembler stops you from changing its value).
R15 is the stack pointer.
R14 is the base pointer.