package main

type Bus interface {}

type CPUBus struct {
  Memory *Memory
  Cartridge *Cartridge
}
