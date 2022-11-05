package main

import "time"

func main() {
  cartridge := loadCartridge("sample.nes")
  wram := NewMemory(2048)
  vram := NewMemory(2048)
  pb := &PPUBus{
    Memory: vram,
  }
  ppu := NewPPU(pb)
  cb := &CPUBus{
    PPU: ppu,
    Cartridge: cartridge,
    Memory: wram,
  }
  cpu := NewCPU(cb)
  cpu.Reset()

  for {
    time.Sleep(time.Millisecond * 0)
    cycle := cpu.Step()
    ppu.Run(uint16(cycle *3))
  }
}
