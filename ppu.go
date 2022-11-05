package main

import (
  "log"
)

type PPURegister [8]byte
/*
  PPUCTRL
  PPUMASK
  PPUSTATUS
  OAMADDR
  OAMDATA
  PPUSCROLL
  PPUADDR
  PPUDATA
*/

const clock = 341
const screenHeight = 240
const screenWidth = 256
const vBlank = 20
const lineMax = screenHeight + vBlank + 2
const tilePixel = 8

type PPU struct {
  Register *PPURegister
  Bus *PPUBus
  vramAddr uint16
  ppuaddrCount int
  cycle uint16
  line uint16
  bg [960]byte
}

type tile struct {
  sprite [16]byte
}

func NewPPU(bus *PPUBus) *PPU {
  var r PPURegister = [8]byte{1,2,3,4,5,6,7,8}
  return &PPU{
    Register: &r,
    Bus: bus,
    ppuaddrCount: 0,
    line: 0,
    bg: [960]byte{},
  }
}

func (p *PPU) Run(cycle uint16) {
  p.cycle += cycle

  if p.line == 0 {
    p.bg = [960]byte{}
  }

  if p.cycle >= clock {
    p.line++
    p.cycle -= clock

    if p.line <= screenHeight && p.line % tilePixel == 0 {
      p.addBgLine()
    }

    if p.line == lineMax {
      p.line = 0
      // complete rendering bg
    }
  }
}

func (p *PPU) addBgLine() {
  for i := 0; i < screenWidth / tilePixel; i++ {
    p.bg[((p.line / tilePixel) - 1) * (screenWidth / tilePixel) + uint16(i)] = 1
  }
}

func (p *PPU) ReadRegister(addr uint16) byte {
  var i byte
  return i
}

func (p *PPU) WriteRegister(addr uint16, data byte) {
  if addr == 0 {
    p.Register[addr] = data
  } else if addr == 1 {
    p.Register[addr] = data
  } else if addr == 5 {
    p.Register[addr] = data
  } else if addr == 6 {
    if p.ppuaddrCount == 0 {
      p.vramAddr = uint16(data << 8)
    } else {
      p.vramAddr += uint16(data)
    }
    p.ppuaddrCount++
  } else if addr == 7 {
    p.Bus.Memory.Write(p.vramAddr, data)
    p.vramAddr++
  } else {
    log.Fatal("unhandle ", addr)
  }
  return
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
