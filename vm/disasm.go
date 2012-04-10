package vm

import (
    "fmt"
)

func (c *EmuCpu) Disassemble(addr uint16) (string, uint16) {
    word := c.Ram[addr]
    next := addr + 1
    /*
    * In bits, a basic instruction has the format: bbbbbbaaaaaaoooo
    * a nonbasic one has the format: aaaaaaoooooo0000
    */
    opcode := word & FOUR_BITS
    if (opcode != NON_L0) {
        // basic opcode
        var op, dest, src string
        dest, next = c.PrintArg((word >>  4) & SIX_BITS, next)
        src,  next = c.PrintArg((word >> 10) & SIX_BITS, next)
        op = c.PrintOpL0(opcode)
        return fmt.Sprintf("%s %s, %s", op, dest, src), next
    }

    // nonbasic opcode
    opcode = (word >> 4) & SIX_BITS
    var op, arg string
    arg, next = c.PrintArg((word >> 10) & SIX_BITS, next)
    op = c.PrintOpL1(opcode)
    return fmt.Sprintf("%s %s", op, arg), next
}

func regToString(reg uint16) string {
    switch reg {
        case REG_A: return "A"
        case REG_B: return "B"
        case REG_C: return "C"
        case REG_X: return "X"
        case REG_Y: return "Y"
        case REG_Z: return "Z"
        case REG_I: return "I"
        case REG_J: return "J"
    }
    return "?"
}

func (c *EmuCpu) PrintArg(a uint16, nextWord uint16) (string, uint16) {
    switch {
    case a < MEM_AT:
        return regToString(a), nextWord
    case a < NW_MEM_AT:
        return fmt.Sprintf("[%s]", regToString(a & THREE_BITS)), nextWord
    case a < POP:
        return fmt.Sprintf("[%s + 0x%x]",regToString(a & THREE_BITS),c.Ram[nextWord]),nextWord+1
    case a == POP:
        return "POP", nextWord
    case a == PEEK:
        return "PEEK", nextWord
    case a == PUSH:
        return "PUSH", nextWord
    case a == REG_SP:
        return "SP", nextWord
    case a == REG_PC:
        return "PC", nextWord
    case a == REG_O:
        return "O", nextWord
    case a == NW:
        return fmt.Sprintf("[0x%x]", c.Ram[nextWord]), nextWord + 1
    case a == LIT_NW:
        return fmt.Sprintf("0x%x", c.Ram[nextWord]), nextWord + 1
    }
    return fmt.Sprintf("0x%x", a ^ LIT), nextWord
}

func (c *EmuCpu) PrintOpL0(opcode uint16) string {
    switch opcode {
        case SET: return "SET"
        case ADD: return "ADD"
        case SUB: return "SUB"
        case MUL: return "MUL"
        case DIV: return "DIV"
        case MOD: return "MOD"
        case SHL: return "SHL"
        case SHR: return "SHR"
        case AND: return "AND"
        case BOR: return "BOR"
        case XOR: return "XOR"
        case IFE: return "IFE"
        case IFN: return "IFN"
        case IFG: return "IFG"
        case IFB: return "IFB"
    }
    return "???"
}

func (c *EmuCpu) PrintOpL1(opcode uint16) string {
    switch opcode {
        case JSR: return "JSR"
    }
    return "???"
}

