package main

import (
    "io"
    "os"
    "fmt"
    "flag"
    "./vm"
)

const BUFSIZE = 0x1000

func loadFile(cpu *vm.EmuCpu, data io.Reader) {
    var buf [BUFSIZE]byte
    var err error = nil
    var count int = 0
    var addr uint16 = 0
    for err == nil {
        count, err = data.Read(buf[:])
        for k := 0; k < (count - count % 2); k += 2 {
            cpu.Ram[addr] = uint16(buf[k]) << 8 + uint16(buf[k+1])
            addr++
        }
        if (err != nil && err != io.EOF) {
            panic(fmt.Sprintf("Unexpected error reading file: %a", err))
        }
    }
}

func display(cpu *vm.EmuCpu) {
    fmt.Printf("\033[2J\033[H")
    for k := 0; k < 80 * 25; k++ {
        c := cpu.Ram[0x8000 + k]
        if (c == 0) {
            fmt.Printf(" ")
        } else {
            fmt.Printf("%c", c)
        }
        if (k % 80 == 79) {
            fmt.Printf("\n")
        }
    }
}

func main() {
    defer (func() {
        if err := recover(); err != nil {
            fmt.Fprintf(os.Stderr, "%s\n", err)
        }
    })()
    flag.Parse()
    cpu := vm.NewEmuCpu()
    fname := flag.Arg(0)
    imagefile, err := os.Open(fname)
    if (err != nil) {
        panic(fmt.Sprintf("Cannot open file %s: %a", fname, err))
    }
    loadFile(cpu, imagefile)
    for {
        var dummy [1]byte
        display(cpu)
        asm, end := cpu.Disassemble(cpu.Pc)
        fmt.Printf("\nnext: %s (%d words, %x)\n(Enter to advance, ^C to quit)\n", asm, end - cpu.Pc, cpu.Ram[cpu.Pc])
        os.Stdin.Read(dummy[:])
        cpu.Step()
    }
}
