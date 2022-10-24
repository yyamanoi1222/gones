package main

import (
  "os"
  "log"
  "io"
)

const NesHeaderSize = 16

type rom struct {
  data []byte
}

func (r *rom) Read() {
}

type Cartridge struct {
  Program *rom
  Char *rom
}


func NewCartridge(program *rom, char *rom) *Cartridge {
  return &Cartridge{
    Program: program,
    Char: char,
  }
}

func loadCartridge(path string) *Cartridge {
  file, err := os.Open(path)
  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }
  defer file.Close()

  hr := io.NewSectionReader(file, 0, NesHeaderSize)
  title := make([]byte, 3)

  _, err = hr.Read(title)
  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }

  if string(title) != "NES" {
    log.Fatal("Invalid NES Header")
  }

  return &Cartridge{}
}
