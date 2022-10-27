package main

type Memory struct {
  data []byte
}

func NewMemory(size uint16) *Memory {
  return &Memory{
    data: make([]byte, size),
  }
}

func (m *Memory) Read(addr uint16) byte {
  return m.data[addr]
}

func (m *Memory) Write(addr uint16) {
}
