/*
 * codes.go
 * Copyright 2012 Jonathan Page. All rights reserved.
 *
 * Contains constants for opcodes, registers, etc.
 */

package vm

// general purpose registers
const (
    REG_A uint16 = iota
    REG_B
    REG_C
    REG_X
    REG_Y
    REG_Z
    REG_I
    REG_J
)

// OR these with register codes above
const MEM_AT uint16 = 0x08
const NW_MEM_AT uint16 = 0x10

const (
    // stack magic
    POP uint16 = iota + 0x18
    PEEK
    PUSH
    // special registers 
    REG_SP
    REG_PC
    REG_O
    // opcode magic
    NW     // [next word]
    LIT_NW // next word (literal)
    LIT    // literal value 0x00-0x1f (OR with value)
)

// basic opcodes
const (
    NON_L0 uint16 = iota
	SET // a, b - sets a to b
	ADD // a, b - sets a to a+b, sets O to 0x0001 if there's an overflow, 0x0 otherwise
	SUB // a, b - sets a to a-b, sets O to 0xffff if there's an underflow, 0x0 otherwise
	MUL // a, b - sets a to a*b, sets O to ((a*b)>>16)&0xffff
    DIV // a, b - sets a to a/b, sets O to ((a<<16)/b)&0xffff. if b==0, sets a and O to 0 instead.
	MOD // a, b - sets a to a%b. if b==0, sets a to 0 instead.
	SHL // a, b - sets a to a<<b, sets O to ((a<<b)>>16)&0xffff
	SHR // a, b - sets a to a>>b, sets O to ((a<<16)>>b)&0xffff
	AND // a, b - sets a to a&b
	BOR // a, b - sets a to a|b
	XOR // a, b - sets a to a^b
	IFE // a, b - performs next instruction only if a==b
	IFN // a, b - performs next instruction only if a!=b
	IFG // a, b - performs next instruction only if a>b
	IFB // a, b - performs next instruction only if (a&b)!=0
)

// nonbasic opcodes
const (
    _ = iota // reserved for future expansion
    JSR uint16 = iota // JSR a - pushes the address of the next instruction to the stack, then sets PC to a
    // 0x02-0x3f: reserved
)


