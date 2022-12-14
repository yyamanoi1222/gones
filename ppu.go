package main

import (
  "log"
  "os"
  "image"
  "image/jpeg"
  "image/color"
  "fmt"
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
  bg [960]tile
  renderable bool
  screen *image.RGBA
}

type sprite [8][8]byte
type tile struct {
  sprite sprite
  paletteId byte
}

func NewPPU(bus *PPUBus) *PPU {
  var r PPURegister = [8]byte{1,2,3,4,5,6,7,8}
  return &PPU{
    Register: &r,
    Bus: bus,
    ppuaddrCount: 0,
    line: 0,
    bg: [960]tile{},
    screen: image.NewRGBA(image.Rect(0, 0, 256, 240)),
  }
}

func (p *PPU) Run(cycle uint16) {
  p.cycle += cycle

  if p.line == 0 {
    p.bg = [960]tile{}
    p.renderable = false
  }

  if p.cycle >= clock {
    p.line++
    p.cycle -= clock

    if p.line <= screenHeight && p.line % tilePixel == 0 {
      p.addBgLine()
    }

    if p.line == lineMax {
      p.renderBackground()
      p.renderable = true
      p.line = 0
    }
  }
}

func (p *PPU) renderBackground() {
  for i, tile := range p.bg {
    x := (i % (screenWidth / tilePixel)) * 8
    y := i / (screenWidth / tilePixel) * 8

    for j := 0; j < 8; j++ {
      for k := 0; k < 8; k++ {
        c := tile.sprite[j][k]
        p.screen.Set(x+k, y+j, color.RGBA{c*100,0,0,0})
      }
    }
  }
  file, _ := os.Create("sample.jpg")
  jpeg.Encode(file, p.screen, &jpeg.Options{100})
  os.Exit(1)
}

func (p *PPU) addBgLine() {
  for i := 0; i < screenWidth / tilePixel; i++ {
    // tileX = 0 ~ 31
    // tileY = 0 ~ 31
    tileX := uint16(i)
    tileY := (p.line / tilePixel) - 1
    addr := tileY * (screenWidth / tilePixel) + tileX
    p.bg[addr] = p.fetchTile(addr, tileX, tileY)
  }
}

func (p *PPU) fetchTile(addr uint16, x uint16, y uint16) tile {
  spriteId := p.Bus.Memory.Read(addr)
  attr := p.fetchFromAttributeTable(x, y)
  fmt.Printf("attr %v \n", attr)

  return tile{
    sprite: p.buildSprite(spriteId),
  }
}

func (p *PPU) fetchFromAttributeTable(x uint16, y uint16) byte {
  addr := 0x03C0 + uint16((x / 4) + (y / 4))
  return p.Bus.Memory.Read(addr)
}

func (p *PPU) buildSprite(spriteId byte) sprite {
  var s sprite = [8][8]byte{}
  for i := 0; i < 16; i++ {
    for j := 0; j < 8; j++ {
      addr := uint16(spriteId) * 16 + uint16(i)
      v := p.Bus.Cartridge.Char.Read(addr)
      if v & (0x80 >> j) > 0 {
        s[i%8][j] += 0x01 << byte(i / 8)
      } else {
        s[i%8][j] += 0
      }
    }
  }
  return s
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
      p.vramAddr = uint16(data) << 8
      p.ppuaddrCount++
    } else {
      p.vramAddr += uint16(data)
      p.ppuaddrCount = 0
    }
  } else if addr == 7 {
    p.Write(p.vramAddr, data)
    p.vramAddr++
  } else {
    log.Fatal("unhandle ", addr)
  }
  return
}

func (p *PPU) Read(addr uint16) byte {
  return p.Bus.Memory.Read(addr)
  if addr > 0x2000 {
    return p.Bus.Memory.Read(addr - 0x2000)
  } else {
    // TODO read from characterRam
    log.Fatal("error ", addr)
  }

  /*
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
  */

  var i byte
  return i
}

func (p *PPU) Write(addr uint16, data byte) {
  if addr >= 0x2000 {
    if addr >= 0x3F00 {
    } else {
      if addr >= 0x3000 {
        p.Bus.Memory.Write(addr - 0x3000, data)
      } else {
        p.Bus.Memory.Write(addr - 0x2000, data)
      }
    }
  } else {
    // TODO write to character ram
    log.Fatal("error ", addr)
  }
}
