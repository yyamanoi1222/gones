package main

type CPU struct {
  Memory Memory
  Register Register
}

type Register struct {
  A byte
  X byte
  Y byte
  S byte
  P byte
  PC uint16
}

func NewCPU() *CPU {
  return &CPU{}
}
