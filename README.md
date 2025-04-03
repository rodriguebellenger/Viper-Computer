# Viper computer

The viper computer is a virtual computer created to learn what programmation is at the lowest level possible.

# Assembly

For the moment, you have the assembly interpreter.

## Operations

| Name       | Effect                                                                                                           |
| ---------- | ---------------------------------------------------------------------------------------------------------------- |
| **HLT**    | Stops the program no matter what.                                                                                |
| **CALL**   | Jump to the address of the label, and push the current address onto the stack.                                   |
| **RET**    | Return to the address at the top of the stack.                                                                   |
| **MOV**    | Move the number to the register specified (see registers section).
The number is 64 bits.|
| **ADD**    | Add the content of two registers and put the result back in the first one.                                       |
| **ADDI**   | Add an immediate to a register (immediate is 8 bits).                                                            |
| **PUSH**   | Push the content of a register onto the stack.                                                                   |
| **POP**    | Pop the top of the stack and put it in a register.                                                               |
| **AND**    | Operate a bitwise and between two registers and put the result back in the first one.                            |
| **OR**     | Operate a bitwise or between two registers and put the result back in the first one.                             |
| **NOT**    | Operate a bitwise not on a register.                                                                             |
| **SWAP**   | Swap the content of two registers.                                                                               |
| **CMP**    | Compare two registers with greater/less/equal operator and skip the next instruction if the comparison is false. |
| **JMP**    | Jump inconditionally to a label.                                                                                 |
| **WRT**    | Write the value of a register at an address contained in an other register, with the size indicated.             |
| **READ**   | Read the value at an address contained in a register, with the size indicated, and put it in a register.         |
| **Label:** | Create a label that you can jump to.                                                                             |

## Registers

There are 16 registers, R0 through R15.

Convention : 
R0 is a constant and always is 0.
