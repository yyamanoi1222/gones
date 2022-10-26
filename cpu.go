package main

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

func NewCPU(bus *CPUBus) *CPU {
  return &CPU{Bus: bus}
}

func (c *CPU) Reset() {
  c.Register.PC = c.ReadDouble(0xFFFC)
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
