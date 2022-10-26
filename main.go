package main

import "fmt"

func main() {
  cartridge := loadCartridge("sample.nes")
  bus := &CPUBus{
    Cartridge: cartridge,
  }
  cpu := NewCPU(bus)
  cpu.Reset()
  fmt.Printf("Hello world %v", cpu)
}
