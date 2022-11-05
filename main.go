package main

import "time"

func main() {
  cartridge := loadCartridge("sample.nes")
  ram := NewMemory(2048)
  ppu := NewPPU()
  bus := &CPUBus{
    PPU: ppu,
    Cartridge: cartridge,
    Memory: ram,
  }
  cpu := NewCPU(bus)
  cpu.Reset()

  for {
    time.Sleep(time.Millisecond * 1)
    cpu.Step()
  }
}
