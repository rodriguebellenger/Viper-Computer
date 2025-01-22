# Viper computer

The viper computer is a virtual computer created to learn what programmation is at the lowest level possible.

# Assembly

For the moment, you have the assembly interpreter, which is almost finished.

## Operations

```
HLT
```
Stops the program no matter what.

---

```
RET
```
Return to the "address" at the top of the stack (it isn't a real address for the moment).

---

```
MOV R1 255
```
Move the number to the register specified (see registers section).
The number is 32 bits.
No access yet to any kind of memory (except registers).

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
Label:
```
Create a label that you can jump to.

## Registers

Registers 0 through 15.
The R0 is constant and will always be 0.
There are still no convention about the use of any other registers.