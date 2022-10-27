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
  P byte
  PC uint16
}

const (
  REGISTER_STATUS_C = 1<<iota
  REGISTER_STATUS_Z = 1<<iota
  REGISTER_STATUS_I = 1<<iota
  REGISTER_STATUS_D = 1<<iota
  REGISTER_STATUS_B = 1<<iota
  REGISTER_STATUS_R = 1<<iota
  REGISTER_STATUS_V = 1<<iota
  REGISTER_STATUS_N = 1<<iota
)

var opecode [256]string = [256]string{
  "BRK", "ORA", "*",   "*", "*",   "ORA", "ASL", "*", "PHP", "ORA", "ASL", "*", "*",   "ORA", "ASL", "*",
  "BPL", "ORA", "*",   "*", "*",   "ORA", "ASL", "*", "CLC", "ORA", "*",   "*", "*",   "ORA", "ASL", "*",
  "JSL", "AND", "*",   "*", "BIT", "AND", "ROL", "*", "PLP", "ANS", "ROL", "*", "BIT", "AND", "ROL", "*",
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
}

func (c *CPU) Step() {
  v := c.Read(c.Register.PC)
  c.Register.PC++
  op := opecode[v]
  am := addressingModeMap[v]

  c.Register.PC++

  addr := c.getAddrFromMode(am)
  c.exec(op, addr)
}

func (c *CPU) getAddrFromMode(mode uint8) uint16 {
  switch mode {
  case 1:
    return 0
  case 2:
    return 0
  case 3:
    v := c.Read(c.Register.PC)
    c.Register.PC++
    return uint16(v)
  case 4:
    v := c.Read(c.Register.PC)
    c.Register.PC++
    return uint16(0x00 << 8 | uint16(v))
  case 5:
    v := c.Read(c.Register.PC)
    c.Register.PC++
    return uint16(0x00 << 8 | (v + c.Register.X))
  case 6:
    v := c.Read(c.Register.PC)
    c.Register.PC++
    return uint16(0x00 << 8 | (v + c.Register.Y))
  case 7:
    l := c.Read(c.Register.PC)
    c.Register.PC++
    u := c.Read(c.Register.PC)
    c.Register.PC++
    return uint16(u << 8 | l)
  case 8:
    l := c.Read(c.Register.PC)
    c.Register.PC++
    u := c.Read(c.Register.PC)
    c.Register.PC++
    return uint16(u << 8 | l) + uint16(c.Register.X)
  case 9:
    l := c.Read(c.Register.PC)
    c.Register.PC++
    u := c.Read(c.Register.PC)
    c.Register.PC++
    return uint16(u << 8 | l) + uint16(c.Register.Y)
  case 10:
    v := c.Read(c.Register.PC)
    c.Register.PC++
    return uint16(v) + c.Register.PC
  case 11:
    v := c.Read(c.Register.PC)
    c.Register.PC++
    b := uint16(0x00 << 8 | v) + uint16(c.Register.X)
    l := c.Read(b)
    h := c.Read(b+1)
    return uint16(h << 8 | l)
  case 12:
    v := c.Read(c.Register.PC)
    c.Register.PC++
    b := uint16(0x00 << 8 | v)
    l := c.Read(b)
    h := c.Read(b+1)
    return uint16(h << 8 | l) + uint16(c.Register.Y)
  case 13:
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

func (c *CPU) exec(operator string, addr uint16) {
  switch operator {
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

func (c *CPU) Write(addr uint16) {
}
