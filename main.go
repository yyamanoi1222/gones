package main

import "fmt"

func main() {
  loadCartridge("sample.nes")
  cpu := NewCPU()
  fmt.Printf("Hello world %v", cpu)
}
