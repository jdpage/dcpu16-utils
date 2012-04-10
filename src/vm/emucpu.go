package vm

import "fmt"

const VM_REG_INIT uint16 = 0x0
const VM_RAM_INIT uint16 = 0x0
const VM_RAMSIZE = 0x10000
const VM_STACK_START uint16 = 0xffff
const VM_PROG_START uint16 = 0x0000

const THREE_BITS   = 0x07
const FOUR_BITS    = 0x0f
const SIX_BITS     = 0x3f

type EmuCpu struct {
    Ram [VM_RAMSIZE]uint16
    Reg [8]uint16
    Pc, Sp, O uint16
    Skip bool
}

type Accessor struct {
    Get func() uint16
    Set func(uint16)
    Eval func()
}

func NewEmuCpu() (cpu *EmuCpu) {
    cpu = new(EmuCpu)
    cpu.Sp = VM_STACK_START
    cpu.Skip = false
    return
}

func (c *EmuCpu) Step() {
    word := c.Ram[c.Pc]
    next := c.Pc + 1
    /*
    * In bits, a basic instruction has the format: bbbbbbaaaaaaoooo
    * a nonbasic one has the format: aaaaaaoooooo0000
    */
    opcode := word & FOUR_BITS
    if (opcode != NON_L0) {
        // basic opcode
        var dest, src Accessor
        dest, next = c.MakeArg((word >>  4) & SIX_BITS, next)
        src,  next = c.MakeArg((word >> 10) & SIX_BITS, next)
        c.Pc = next
        if (!c.Skip) {
            dest.Eval()
            src.Eval()
            c.EvalOpL0(opcode, dest, src)
        } else {
            c.Skip = false
        }
    } else {
        // nonbasic opcode
        opcode = (word >> 4) & SIX_BITS
        var arg Accessor
        arg, next = c.MakeArg((word >> 10) & SIX_BITS, next)
        c.Pc = next
        if (!c.Skip) {
            arg.Eval()
            c.EvalOpL1(opcode, arg)
        } else {
            c.Skip = false
        }
    }
}

func (c *EmuCpu) MakeArg(a uint16, nextWord uint16) (Accessor, uint16) {
    switch {
    case a < MEM_AT:
        // register
        return register(c, a), nextWord
    case a < NW_MEM_AT:
        // memory at register
        return registerRef(c, a & THREE_BITS), nextWord
    case a < POP:
        // memory at (register + next word)
        return regPlusRef(c, a & THREE_BITS, c.Ram[nextWord]), (nextWord + 1)
    case a == POP:
        return literalRef(c, c.Sp, func(){ c.Sp += 1 }), nextWord
    case a == PEEK:
        return literalRef(c, c.Sp, func(){}), nextWord
    case a == PUSH:
        return literalRef(c, c.Sp - 1, func(){ c.Sp -= 1 }), nextWord
    case a == REG_SP:
        return regSp(c), nextWord
    case a == REG_PC:
        return regPc(c), nextWord
    case a == REG_O:
        return regO(c), nextWord
    case a == NW:
        return literalRef(c, c.Ram[nextWord], func(){}), (nextWord + 1)
    case a == LIT_NW:
        return literalVal(c.Ram[nextWord]), (nextWord + 1)
    }
    return literalVal(a ^ LIT), nextWord
}

func (c *EmuCpu) EvalOpL0(opcode uint16, a, b Accessor) {
    switch opcode {
    case SET:
        a.Set(b.Get())
    case ADD:
        var r uint32 = uint32(a.Get()) + uint32(b.Get())
        a.Set(uint16(r))
        if ((r & 0xffff) == r) {
            c.O = 0
        } else {
            c.O = 1
        }
    case SUB:
        if (a.Get() < b.Get()) {
            c.O = 0xffff
        } else {
            c.O = 0
        }
        a.Set(a.Get() - b.Get())
    case MUL:
        c.O = ((a.Get() * b.Get()) >> 16) & 0xffff
        a.Set(a.Get() * b.Get())
    case DIV:
        if (b.Get() == 0) {
            a.Set(0)
            c.O = 0
        } else {
            c.O = ((a.Get() << 16) / b.Get()) & 0xffff
            a.Set(a.Get() / b.Get())
        }
    case MOD:
        if (b.Get() == 0) {
            a.Set(0)
        } else {
            a.Set(a.Get() % b.Get())
        }
    case SHL:
        c.O = ((a.Get() << b.Get()) >> 16) & 0xffff
        a.Set(a.Get() << b.Get())
    case SHR:
        c.O = ((a.Get() << 16) >> b.Get()) & 0xffff
        a.Set(a.Get() >> b.Get())
    case AND:
        a.Set(a.Get() & b.Get())
    case BOR:
        a.Set(a.Get() | b.Get())
    case XOR:
        a.Set(a.Get() ^ b.Get())
    case IFE:
        c.Skip = !(a.Get() == b.Get())
    case IFN:
        c.Skip = !(a.Get() != b.Get())
    case IFG:
        c.Skip = !(a.Get() > b.Get())
    case IFB:
        c.Skip = (a.Get() & b.Get()) == 0
    default:
        panic(fmt.Sprintf("unknown L0 opcode %x", opcode))
    }
}

func (c *EmuCpu) EvalOpL1(opcode uint16, a Accessor) {
    switch opcode {
    case JSR:
        c.Sp += 1
        c.Ram[c.Sp] = c.Pc
        c.Pc = a.Get()
    default:
        panic(fmt.Sprintf("unknown L1 opcode %x", opcode))
    }
}

func register(cpu *EmuCpu, reg uint16) Accessor {
    return Accessor{
        func() uint16 { return cpu.Reg[reg] },
        func(val uint16) { cpu.Reg[reg] = val },
        func(){},
    }
}

func registerRef(cpu *EmuCpu, reg uint16) Accessor {
    return Accessor{
        func() uint16 { return cpu.Ram[cpu.Reg[reg]] },
        func(val uint16) { cpu.Ram[cpu.Reg[reg]] = val },
        func(){},
    }
}

func regPlusRef(cpu *EmuCpu, reg uint16, lit uint16) Accessor {
    return Accessor{
        func() uint16 { return cpu.Ram[cpu.Reg[reg] + lit] },
        func(val uint16) { cpu.Ram[cpu.Reg[reg] + lit] = val },
        func(){},
    }
}

func regSp(cpu *EmuCpu) Accessor {
    return Accessor{
        func() uint16 { return cpu.Sp },
        func(val uint16) { cpu.Sp = val },
        func(){},
    }
}

func regPc(cpu *EmuCpu) Accessor {
    return Accessor{
        func() uint16 { return cpu.Pc },
        func(val uint16) { cpu.Pc = val },
        func(){},
    }
}

func regO(cpu *EmuCpu) Accessor {
    return Accessor{
        func() uint16 { return cpu.O },
        func(val uint16) { cpu.O = val },
        func(){},
    }
}

func literalRef(cpu *EmuCpu, ref uint16, after func()) Accessor {
    return Accessor{
        func() uint16 { return cpu.Ram[ref] },
        func(val uint16) { cpu.Ram[ref] = val },
        after,
    }
}

func literalVal(val uint16) Accessor {
    return Accessor{
        func() uint16 { return val },
        func(val uint16) { /* fail silently */ },
        func(){},
    }
}

