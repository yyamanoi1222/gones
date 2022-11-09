package main

type Bus interface {}

type CPUBus struct {
  Memory *Memory
  Cartridge *Cartridge
  PPU *PPU
}

type PPUBus struct {
  Memory *Memory
  Cartridge *Cartridge
}
