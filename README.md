# Viper computer

The viper computer is a virtual computer created for educationnal purposes.
This README is NOT up to date, but I'm too lazy to do it now, please check "instruction_set.txt" instead.

# Assembly

For the moment, you have the assembly interpreter.

## Operations

```
HLT
```
Stops the program no matter what.

---

```
CALL Label
```
Jump to the address of the label, and push the current address onto the stack.

---

```
RET
```
Return to the address at the top of the stack.

---

```
MOV R1 255
```
Move the number to the register specified (see registers section).
The number is 64 bits.

---

```
ADD R1 R2
```
Add the content of two registers and put the result back in the first one.

---

```
ADDI R1 175
```
Add an immediate to a register (immediate is 8 bits).

---

```
PUSH R1
```
Push the content of a register onto the stack.

---

```
POP R1
```
Pop the top of the stack and put it in a register.

---

```
AND R1 R2
```
Operate a bitwise and between two registers and put the result back in the first one.

---

```
OR R1 R2
```
Operate a bitwise or between two registers and put the result back in the first one.

---

```
NOT R1
```
Operate a bitwise not on a register.

---

```
SWAP R1 R2
```
Swap the content of two registers.

---

```
CMP R1 R2 (G/L/E)
```
Compare two registers with greater/less/equal operator and skip the next instruction if the comparison is false.

---

```
JMP Label
```
Jump inconditionally to a label.

---

```
WRT @32 *R1 R2
```
Write the value of a register at an address contained in a register, with the size indicated.
In this case, it puts the value of R2 at the addres contained in R1, plus the three next address because the value is 32bits.

---

```
READ R1 @32 *R2
```
Read the value at an address contained in a register, with the size indicated, and put it in a register.
In this case, it reads the value of the address contained in R2, plus the three next address because the value is 32bits, and puts it in R1.

---

```
Label:
```
Create a label that you can jump to.

## Registers

There are 16 registers, R0 through R15.

Convention : 
R0 is a constant and always is 0. Nothing stop you from changing it, but I don't think it should.
