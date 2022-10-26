package main

type Memory struct {
  data []byte
}

func NewMemory() *Memory {
  return &Memory{}
}

func (m *Memory) Read(addr uint16) byte {
  return m.data[addr]
}
