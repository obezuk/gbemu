package main
import (
  "github.com/faiface/pixel"
  "github.com/faiface/pixelgl"
  // "bufio"
  // "strconv"
  "fmt"
  "log"
  "io/ioutil"
  // "reflect"
  // "os"
)

var _RAM [0xFFFF] uint8
var _PC uint16 = 0x0000;
var _SP uint16 = 0x0000;

var _A uint8 = 0x00;
var _B uint8 = 0x00;
var _C uint8 = 0x00;
var _D uint8 = 0x00;
var _E uint8 = 0x00;
var _F uint8 = 0x00;
var _HL uint16 = 0x0000;

var FlagZero bool;
var FlagSubtract bool;
var FlagHalfCarry bool;
var FlagCarry bool;

func debug() {
  fmt.Printf("\n");
  fmt.Printf("PC: 0x%X SP: 0x%X task: 0x%X\n", _PC, _SP, loadMemory(_PC));
  fmt.Printf("A: 0x%X B: 0x%X C: 0x%X D: 0x%X E: 0x%X \n", _A, _B, _C, _D, _E);
  fmt.Printf("Z: %t N: %t H: %t C: %t \n", FlagZero, FlagSubtract, FlagHalfCarry, FlagCarry);
  fmt.Printf("\n");
  log.Fatal("die\n");
}

func writeMemory(position uint16, value uint8) {
  fmt.Printf("RAM Write: 0x%x:0x%x\n", position, value);
  _RAM[position] = value;
}

func loadMemory(position uint16) uint8 {
  value := _RAM[position];
  fmt.Printf("RAM Read: 0x%x:0x%x\n", position, value);
  return value;
}

func CPU(task uint8) {
  switch task {
  case 0x0C: // INC C - Z 0 H -
    _C = _C + 1;
    // TODO FlagZero
    // TODO HalfFlag
    FlagSubtract = false;
    _PC = _PC + 1;
  case 0x0E: // LD C,d8
    _C = loadMemory(_PC + 1);
    _PC = _PC + 2;
  case 0x20: // JR NZ,r8
    if (FlagZero) {
      _PC = uint16(uint8(_PC + 2) + loadMemory(_PC + 1)); // Jump to relative memory position.
    } else {
      _PC = _PC + 2;
    }
  case 0x31: // LD SP
    Lo := uint16(loadMemory(_PC + 1));
    Hi := uint16(loadMemory(_PC + 2));
    _SP = ((Hi << 8) | Lo);
    _PC = _PC + 3;
  case 0x32: // LD (HL-),A
    writeMemory(_HL, _A);
    _HL = _HL - 1;
    _PC = _PC + 1;
  case 0x3E: // LD A,d8
    _A = loadMemory(_PC + 1);
    _PC = _PC + 2;
  case 0xAF: // XOR A
    _A = 0;
    FlagZero = true;
    FlagSubtract = false;
    FlagHalfCarry = false;
    FlagCarry = false;
    _PC = _PC + 1;
  case 0x21: // LD HL,d16
    Lo := uint16(loadMemory(_PC + 1));
    Hi := uint16(loadMemory(_PC + 2));
    _HL = ((Hi << 8) | Lo);
    _PC = _PC + 3;
  case 0x77: // LD (HL),A
    writeMemory(_HL, _A);
    _PC = _PC + 1;
  case 0xCB: // Prefix CB
    CB_CPU(loadMemory(_PC + 1)) // Delegate to CB CPU Function...
    _PC = _PC + 1; // Bump past 0xCB
  case 0xE0: // LDH (a8),A
    // THIS IS PROBABLY BROKEN... 0x001F
    _A = uint8(uint16(0xFF00) + uint16(loadMemory(_PC + 1)));
    _PC = _PC + 2;
  case 0xE2: // LD (C),A
    _A = _C
    _PC = _PC + 1;
  default:
    debug();
  }
}

func CB_CPU(task uint8) {

  switch task {
  case 0x7C: // BIT 7,H
    _H := uint8(_HL & 0xFF);
    FlagZero = ((_H >> 7) == 0);
    FlagSubtract = false;
    FlagHalfCarry = true;
    _PC = _PC + 1;
  default:
    log.Fatal("Unknown CB OP CODE...");
  }
}

func main() {

  // Load Boot Rom into RAM

  _ = _A;
  _ = _B;
  _ = _C;
  _ = _D;
  _ = _E;
  _ = _F;

  rom, err := ioutil.ReadFile("DMG_ROM.bin")

  if err != nil {
    log.Fatal(err)
	}

  for index, element := range rom {
    writeMemory(uint16(index), element);
  }

  for {
    task := loadMemory(_PC);
    CPU(task);
  }

}