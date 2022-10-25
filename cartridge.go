package main

import (
  "os"
  "log"
  "io"
)

const NesHeaderSize int64 = 16
const ProgramRomBaseSize int64 = 16384
const CharRomBaseSize int64 = 8192

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
  header := make([]byte, 16)
  _, err = hr.Read(header)
  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }

  title := header[0:3]

  if string(title) != "NES" {
    log.Fatal("Invalid NES Header")
    os.Exit(1)
  }

  prgSize := int64(header[4]) * ProgramRomBaseSize
  charSize := int64(header[5]) * CharRomBaseSize

  pr := io.NewSectionReader(file, NesHeaderSize, prgSize)
  prgRom, err := io.ReadAll(pr)
  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }

  cr := io.NewSectionReader(file, prgSize + 1, charSize)
  chrRom, err := io.ReadAll(cr)
  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }

  return NewCartridge(
    &rom{data: prgRom},
    &rom{data: chrRom},
  )
}
