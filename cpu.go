package main

import "log"

type CPU struct {
  Register Register
  Bus *CPUBus
}

type Register struct {
  A byte
  X byte
  Y byte
  S byte
  PC uint16
  P *statusRegister
}

type statusRegister struct {
  C bool
  Z bool
  I bool
  D bool
  B bool
  R bool
  V bool
  N bool
}

var opecode [256]string = [256]string{
  "BRK", "ORA", "*",   "*", "*",   "ORA", "ASL", "*", "PHP", "ORA", "ASL", "*", "*",   "ORA", "ASL", "*",
  "BPL", "ORA", "*",   "*", "*",   "ORA", "ASL", "*", "CLC", "ORA", "*",   "*", "*",   "ORA", "ASL", "*",
  "JSR", "AND", "*",   "*", "BIT", "AND", "ROL", "*", "PLP", "ANS", "ROL", "*", "BIT", "AND", "ROL", "*",
  "BMI", "AND", "*",   "*", "*",   "AND", "ROL", "*", "SEC", "AND", "*",   "*", "*",   "AND", "ROL", "*",
  "RTI", "EOR", "*",   "*", "*",   "EOR", "LSR", "*", "PHA", "EOR", "LSR", "*", "JMP", "EOR", "LSR", "*",
  "BVC", "EOR", "*",   "*", "*",   "EOR", "LSR", "*", "CLI", "EOR", "*",   "*", "*",   "EOR", "LSR", "*",
  "RTS", "ADC", "*",   "*", "*",   "ADC", "ROR", "*", "PLA", "ADC", "ROR", "*", "JMP", "ADC", "ROR", "*",
  "BVS", "ADC", "*",   "*", "*",   "ADC", "ROR", "*", "SEI", "ADC", "*",   "*", "*",   "ADC", "ROR", "*",
  "*",   "STA", "*",   "*", "STY", "STA", "STX", "*", "DEY", "*",   "TXA", "*", "STY", "STA", "STX", "*",
  "BBC", "STA", "*",   "*", "STY", "STA", "STX", "*", "TYA", "STA", "TXS", "*", "*",   "STA", "*",   "*",
  "LDY", "LDA", "LDX", "*", "LDY", "LDA", "LDX", "*", "TAY", "LDA", "TAX", "*", "LDY", "LDA", "LDX", "*",
  "BCS", "LDA", "*",   "*", "LDY", "LDA", "LDX", "*", "CLV", "LDA", "TSX", "*", "LDY", "LDA", "LDX", "*",
  "CPY", "CMP", "*",   "*", "CPY", "CMP", "DEC", "*", "INY", "CMP", "DEX", "*", "CPY", "CMP", "DEC", "*",
  "BNE", "CMP", "*",   "*", "*",   "CMP", "DEC", "*", "CLD", "CMP", "*",   "*", "*",   "CMP", "DEC", "*",
  "CPX", "SBC", "*",   "*", "CPX", "SBC", "INC", "*", "INX", "SBC", "NOP", "*", "CPX", "SBC", "INC", "*",
  "BEQ", "SBC", "*",   "*", "*",   "SBC", "INC", "*", "SED", "SBC", "*",   "*", "*",   "SBC", "INC", "*",
}

const (
  _ = iota
  ADDRESSING_MODE_IMPL
  ADDRESSING_MODE_A
  ADDRESSING_MODE_IM
  ADDRESSING_MODE_ZPG
  ADDRESSING_MODE_ZPX
  ADDRESSING_MODE_ZPY
  ADDRESSING_MODE_ABS
  ADDRESSING_MODE_ABSX
  ADDRESSING_MODE_ABSY
  ADDRESSING_MODE_REL
  ADDRESSING_MODE_IIND
  ADDRESSING_MODE_INDI
  ADDRESSING_MODE_ABSI
)

var addressingModeMap [256]uint8 = [256]uint8{
  1,  11, 0, 0, 0, 4, 4, 0, 1, 3, 2, 0, 0,  7, 7, 0,
  10, 12, 0, 0, 0, 5, 5, 0, 1, 9, 0, 0, 0,  8, 8, 0,
  7,  11, 0, 0, 4, 4, 4, 0, 1, 3, 2, 0, 7,  7, 7, 0,
  10, 12, 0, 0, 0, 5, 5, 0, 1, 9, 0, 0, 0,  8, 8, 0,
  1,  11, 0, 0, 0, 4, 4, 0, 1, 3, 2, 0, 7,  7, 7, 0,
  10, 12, 0, 0, 0, 5, 5, 0, 1, 9, 0, 0, 0,  8, 8, 0,
  1,  11, 0, 0, 0, 4, 4, 0, 1, 3, 2, 0, 13, 7, 7, 0,
  10, 12, 0, 0, 0, 5, 5, 0, 1, 9, 0, 0, 0,  8, 8, 0,
  0,  11, 0, 0, 4, 4, 4, 0, 1, 0, 1, 0, 7,  7, 7, 0,
  10, 12, 0, 0, 5, 5, 6, 0, 1, 9, 1, 0, 0,  8, 0, 0,
  3,  11, 3, 0, 4, 4, 4, 0, 1, 3, 1, 0, 7,  7, 7, 0,
  10, 12, 0, 0, 5, 5, 6, 0, 1, 9, 1, 0, 8,  8, 9, 0,
  3,  11, 0, 0, 4, 4, 4, 0, 1, 3, 1, 0, 7,  7, 7, 0,
  10, 12, 0, 0, 0, 5, 5, 0, 1, 9, 0, 0, 0,  8, 8, 0,
  3,  11, 0, 0, 4, 4, 4, 0, 1, 3, 1, 0, 7,  7, 7, 0,
  10, 12, 0, 0, 0, 5, 5, 0, 1, 9, 0, 0, 0,  8, 8, 0,
}

func NewCPU(bus *CPUBus) *CPU {
  return &CPU{Bus: bus}
}

func (c *CPU) Reset() {
  c.Register.PC = c.ReadDouble(0xFFFC)
  c.Register.P = &statusRegister{}
}

func (c *CPU) Step() {
  log.Printf("---Start Step--- \n")
  log.Printf("pc %v \n", c.Register.PC)

  v := c.Read(c.Register.PC)
  op := opecode[v]
  am := addressingModeMap[v]

  log.Printf("op %v \n", op)
  log.Printf("am %v \n", am)

  c.Register.PC++

  addr := c.getAddrFromMode(am)
  c.exec(op, addr, am)

  log.Printf("---End Step--- \n\n")
}

func (c *CPU) getAddrFromMode(mode uint8) uint16 {
  switch mode {
  case ADDRESSING_MODE_IMPL:
    return 0
  case ADDRESSING_MODE_A:
    return 0
  case ADDRESSING_MODE_IM:
    v := c.Read(c.Register.PC)
    c.Register.PC++
    return uint16(v)
  case ADDRESSING_MODE_ZPG:
    v := c.Read(c.Register.PC)
    c.Register.PC++
    return uint16(0x00 << 8 | uint16(v))
  case ADDRESSING_MODE_ZPX:
    v := c.Read(c.Register.PC)
    c.Register.PC++
    return uint16(0x00 << 8 | (v + c.Register.X))
  case ADDRESSING_MODE_ZPY:
    v := c.Read(c.Register.PC)
    c.Register.PC++
    return uint16(0x00 << 8 | (v + c.Register.Y))
  case ADDRESSING_MODE_ABS:
    r := c.ReadDouble(c.Register.PC)
    c.Register.PC+=2
    return r
  case ADDRESSING_MODE_ABSX:
    l := c.Read(c.Register.PC)
    c.Register.PC++
    u := c.Read(c.Register.PC)
    c.Register.PC++
    return uint16(u << 8 | l) + uint16(c.Register.X)
  case ADDRESSING_MODE_ABSY:
    l := c.Read(c.Register.PC)
    c.Register.PC++
    u := c.Read(c.Register.PC)
    c.Register.PC++
    return uint16(u << 8 | l) + uint16(c.Register.Y)
  case ADDRESSING_MODE_REL:
    v := c.Read(c.Register.PC)
    c.Register.PC++

    if v < 0x80 {
      return uint16(v) + c.Register.PC
    } else {
      return uint16(v) + c.Register.PC - 256
    }
  case ADDRESSING_MODE_IIND:
    v := c.Read(c.Register.PC)
    c.Register.PC++
    b := uint16(0x00 << 8 | v) + uint16(c.Register.X)
    l := c.Read(b)
    h := c.Read(b+1)
    return uint16(h << 8 | l)
  case ADDRESSING_MODE_INDI:
    v := c.Read(c.Register.PC)
    c.Register.PC++
    b := uint16(0x00 << 8 | v)
    l := c.Read(b)
    h := c.Read(b+1)
    return uint16(h << 8 | l) + uint16(c.Register.Y)
  case ADDRESSING_MODE_ABSI:
    l := c.Read(c.Register.PC)
    c.Register.PC++
    u := c.Read(c.Register.PC)
    c.Register.PC++
    b := uint16(c.Read(uint16(u << 8 | l)))
    ll := c.Read(b)
    hh := c.Read(b+1)
    return uint16(hh << 8 | ll)
  default:
    log.Fatal("cannnot hadle mode ", mode)
  }
  return 0
}

func (c *CPU) exec(operator string, addr uint16, mode uint8) {
  switch operator {
  case "LDX":
    v := c.Read(addr)
    c.Register.X = v
    c.Register.P.Z = v == 0
    c.Register.P.N = v & 0x80 > 0
    return
  case "LDA":
    v := c.Read(addr)
    c.Register.A = v
    c.Register.P.Z = v == 0
    c.Register.P.N = v & 0x80 > 0
    return
  case "LDY":
    v := c.Read(addr)
    c.Register.Y = v
    c.Register.P.Z = v == 0
    c.Register.P.N = v & 0x80 > 0
  case "INX":
    v := c.Register.X + 1
    c.Register.X = v
    c.Register.P.Z = v == 0
    c.Register.P.N = v & 0x80 > 0
  case "DEY":
    v := c.Register.Y - 1
    c.Register.Y = v
    c.Register.P.Z = v == 0
    c.Register.P.N = v & 0x80 > 0
  case "TXS":
    c.Register.S = c.Register.X
  case "STA":
    c.Write(addr, c.Register.A)
    return
  case "ASL":
    c.Register.P.C = c.Register.A >> 7 & 1 == 1
    r := c.Register.A << 1
    c.Register.A = r
    c.Register.P.Z = r == 0
    c.Register.P.N = r & 0x80 > 0
  case "BPL":
    if !c.Register.P.N {
      c.Register.PC = addr
    }
  case "BNE":
    if !c.Register.P.Z {
      c.Register.PC = addr
    }
  case "SEI":
    c.Register.P.I = false
    return
  case "JSR":
    c.Register.PC = addr
    // Push to stack
  case "JMP":
    c.Register.PC = addr
    return
  case "BRK":
    if !c.Register.P.I {
      return
    }
    c.Register.P.B = true
    c.Register.PC++
    // TODO
    // Push to stack
    c.Register.P.I = true
    return
  default:
    log.Fatal("cannnot hadle operator ", operator)
  }
}

func (c *CPU) ReadDouble(addr uint16) uint16 {
  l := uint16(c.Read(addr))
  u := uint16(c.Read(addr + 1))
  return u<<8 | l
}

func (c *CPU) Read(addr uint16) byte {
  if addr < 0x0800 {
    // Read From RAM
    return c.Bus.Memory.Read(addr)
  } else if addr < 0x2000 {
    return c.Bus.Memory.Read(addr - 0x800)
  } else if addr < 0x2008 {
    // PPU Register
  } else if addr < 0x4000 {
    // PPU Register Mirror
  } else if addr < 0x4020 {
    // APU I/O PAD
  } else if addr < 0x8000 {
    // ext ROM
  } else if addr < 0xFFFF {
    // Read From PROGRAM ROM
    return c.Bus.Cartridge.Program.Read(addr - 0x8000)
  }

  var i byte
  return i
}

func (c *CPU) Write(addr uint16, data byte) {
}

func (c *CPU) push(data byte) {
}

func (c *CPU) pop() byte {
  var i byte
  return i
}
