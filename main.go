package main

import (
  	"fmt"
  	"log"
	"io/ioutil"
	//"github.com/veandco/go-sdl2/sdl"
) 

const MEMSIZE int = 4096

type machine struct {
  	mem [MEMSIZE]uint8

	v [16]uint8
	i uint16
	pc uint16

	st uint8
	dt uint8
	sp uint8
	stack [16]uint16

	display [32][64]uint8
  	keys [16]uint8
}

/*var sprites []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, //[0]
	0x20, 0x60, 0x20, 0x20, 0x70, //[1]
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //[2]
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //[3]
	0x90, 0x90, 0xF0, 0x10, 0x10, //[4]
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //[5]
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //[6]
	0xF0, 0x10, 0x20, 0x40, 0x40, //[7]
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //[8]
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //[9]
	0xF0, 0x90, 0xF0, 0x90, 0x90, //[A]
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //[B]
	0xF0, 0x80, 0x80, 0x80, 0xF0, //[C]
	0xE0, 0x90, 0x90, 0x90, 0xE0, //[D]
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //[E]
	0xF0, 0x80, 0xF0, 0x80, 0x80, //[F]
}*/

func main() {
	m := Init()
	err := m.Load("pong.rom")
	if err != nil {
		log.Fatal(err)
	}
	err = m.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func (m *machine) Load(f string) error {

	p, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}
  
	for i := 0; i < len(p); i++ {
		m.mem[i+512] = p[i]
		fmt.Printf("%x \n", p[i])
	}
  
	return nil
}

func Init() machine{
	
	im := machine{
		pc: 0x200,
	}

	/*for i := 0; i < len(sprites); i++ {
		im.mem[i] = sprites[i]
	}*/

	for i := 512; i < 4096; i++ {
		im.mem[i] = 0x00
	}

	for i := 0; i < 16; i++ {
		im.stack[i] = 0
		im.v[i] = 0
	}
	
	return im
}

func (m *machine) Run() error {
	opcode := (uint16(m.mem[m.pc]) << 8) | uint16(m.mem[m.pc+1])
	var nnn uint16 = (opcode & 0x0FFF)
	var kk uint8 = (opcode & 0x0FF)
	var n uint8 = opcode & 0x0F 
	var x uint8 = uint16(opcode >> 8) & 0x0F
	var y uint8 = uint16(opcode >> 4) & 0x0F
	//var mustQuit bool = false
	fmt.Printf("%x %x %x %x %x %x\n", opcode, nnn, kk, n, x, y)
	/*while(!mustQuit) {

		if (m.pc+2) == MEMSIZE {
			m.pc+=2
		}

		switch opcode {
			case 0x6a02: 
				fmt.Println("6a02")
			default:
				fmt.Println("NO")
		}
	}*/

  	if m.st > 0 {
    	m.st -=1
  	}
  
  	if m.dt > 0 {
    	m.dt -=1
  	}

  	return nil
}