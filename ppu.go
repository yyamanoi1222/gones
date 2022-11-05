package main

type PPURegister [8]byte

type PPU struct {
  Register *PPURegister
}

/*
type PPURegister struct {
  PPUCTRL byte
  PPUMASK byte
  PPUSTATUS byte
  OAMADDR byte
  OAMDATA byte
  PPUSCROLL byte
  PPUADDR byte
  PPUDATA byte
}
*/

func NewPPU() *PPU {
  var r PPURegister = [8]byte{1,2,3,4,5,6,7,8}
  return &PPU{
    Register: &r,
  }
}

func (p *PPU) Read(addr uint16) byte {
  if addr < 0x1000 {
    // pattern table 0
  } else if addr < 0x2000 {
    // pattern table 1
  } else if addr < 0x23C0 {
    // name table 0
  } else if addr < 0x2400 {
    // attr table 0
  } else if addr < 0x27C0 {
    // name table 1
  } else if addr < 0x2800 {
    // attr table 1
  } else if addr < 0x2BC0 {
    // name table 2
  } else if addr < 0x2C00 {
    // attr table 2
  } else if addr < 0x2FC0 {
    // name table 3
  } else if addr < 0x3000 {
    // attr table 3
  } else if addr < 0x3F00 {
    // mirror
  } else if addr < 0x3F10 {
    // background palette
  } else if addr < 0x3F20 {
    // splite palette
  } else if addr <= 0x3FFF {
    // mirror
  }

  var i byte
  return i
}

func (p *PPU) Write(addr uint16, data byte) {
}
