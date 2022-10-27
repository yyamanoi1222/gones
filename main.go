package main

func main() {
  cartridge := loadCartridge("sample.nes")
  ram := NewMemory(2048)
  bus := &CPUBus{
    Cartridge: cartridge,
    Memory: ram,
  }
  cpu := NewCPU(bus)
  cpu.Reset()
  // cpu.Step()
}
