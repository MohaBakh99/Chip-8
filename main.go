package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

const MEMSIZE int = 4096

type machine struct {
	mem [MEMSIZE]uint8

	v  [16]uint8
	i  uint16
	pc uint16

	st    uint8
	dt    uint8
	sp    uint8
	stack [16]uint16

	display [32][64]uint8
	keys    [16]uint8
}

func main() {
	m := Init()
	err := m.Load("./roms/PONG")
	if err != nil {
		log.Fatal(err)
	}
	for z := 0; z < 80; z++ {
		fmt.Printf("sprite #%d = %x: \n", z, m.mem[z])
	}

	/*err = Draw()
	if err != nil {
		log.Fatal(err)
	}*/

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
	}

	return nil
}

func Init() machine {

	im := machine{
		pc: 0x200,
	}

	sprites := []uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0,
		0x20, 0x60, 0x20, 0x20, 0x70,
		0xF0, 0x10, 0xF0, 0x80, 0xF0,
		0xF0, 0x10, 0xF0, 0x10, 0xF0,
		0x90, 0x90, 0xF0, 0x10, 0x10,
		0xF0, 0x80, 0xF0, 0x10, 0xF0,
		0xF0, 0x80, 0xF0, 0x90, 0xF0,
		0xF0, 0x10, 0x20, 0x40, 0x40,
		0xF0, 0x90, 0xF0, 0x90, 0xF0,
		0xF0, 0x90, 0xF0, 0x10, 0xF0,
		0xF0, 0x90, 0xF0, 0x90, 0x90,
		0xE0, 0x90, 0xE0, 0x90, 0xE0,
		0xF0, 0x80, 0x80, 0x80, 0xF0,
		0xE0, 0x90, 0x90, 0x90, 0xE0,
		0xF0, 0x80, 0xF0, 0x80, 0xF0,
		0xF0, 0x80, 0xF0, 0x80, 0x80}

	for i := 0; i < 80; i++ {
		im.mem[i] = sprites[i]
	}

	for i := 0x200; i < MEMSIZE; i++ {
		im.mem[i] = 0x00
	}

	for i := 0; i < 16; i++ {
		im.stack[i] = 0
		im.v[i] = 0
	}

	keys := []uint8{
		1, 2, 3, 0x0c,
		4, 5, 6, 0x0d,
		7, 8, 9, 0x0e,
		0x0a, 0, 0x0b, 0x0f,
	}

	for i := 0; i < 16; i++ {
		im.keys[i] = keys[i]
	}

	return im
}

var mustQuit bool = false

func (m *machine) Run() error {
	for !mustQuit {
		opcode := (uint16(m.mem[m.pc]) << 8) | uint16(m.mem[m.pc+1])
		m.pc += 2
		if m.pc == uint16(MEMSIZE) {
			m.pc = 0x200
		} else {
			nnn := opcode & 0x0FFF
			kk := opcode & 0x0FF
			n := opcode & 0x0F
			x := uint16(opcode>>8) & 0x0F
			y := uint16(opcode>>4) & 0x0F
			switch op := opcode >> 12; op {
			case 0:
				if opcode == 0x00EE {
					fmt.Println("RET")
				} else if opcode == 0x00E0 {
					fmt.Println("CLS")
				} else if opcode == 0x0000 {
					m.pc = 0x200
				} else {
					fmt.Println("SYS")
					m.pc = nnn
				}
				break
			case 1:
				fmt.Printf("JP nnn:%x\n", nnn)
				m.pc = nnn
				break
			case 2:
				fmt.Printf("CALL nnn:%x\n", nnn)
				m.sp += 1
				m.stack[0] = m.pc
				m.pc = nnn
				break
			case 3:
				fmt.Printf("SE x:%x kk:%x\n", x, kk)
				if uint16(m.v[x]) == kk {
					m.pc += 2
				}
				break
			case 4:
				fmt.Printf("SNE x:%x kk:%x\n", x, kk)
				if !(uint16(m.v[x]) == kk) {
					m.pc += 2
				}
				break
			case 5:
				fmt.Printf("SE x:%x y:%x\n", x, y)
				if uint16(m.v[x]) == uint16(m.v[y]) {
					m.pc += 2
				}
				break
			case 6:
				fmt.Printf("LD x:%x kk:%x\n", x, kk)
				m.v[x] = uint8(kk)
				break
			case 7:
				fmt.Printf("ADD x:%x kk:%x\n", x, kk)
				m.v[x] += uint8(kk)
				break
			case 8:
				if n == 0 {
					fmt.Printf("LS x:%x y:%x\n", x, y)
					m.v[x] = m.v[y]
				} else if n == 1 {
					fmt.Printf("OR x:%x y:%x\n", x, y)
					m.v[x] |= m.v[y]
				} else if n == 2 {
					fmt.Printf("AND x:%x y:%x\n", x, y)
					m.v[x] &= m.v[y]
				} else if n == 3 {
					fmt.Printf("XOR x:%x y:%x\n", x, y)
					m.v[x] ^= m.v[y]
				} else if n == 4 {
					fmt.Printf("ADD x:%x y:%x\n", x, y)
					m.v[x] += m.v[y]
					if m.v[x] > 255 {
						m.v[15] = 1
					} else {
						m.v[15] = 0
					}
				} else if n == 5 {
					fmt.Printf("SUB x:%x y:%x\n", x, y)
					m.v[x] -= m.v[y]
					if m.v[x] > m.v[y] {
						m.v[15] = 1
					} else {
						m.v[15] = 0
					}
				} else if n == 6 {
					fmt.Printf("SHR x:%x y:%x\n", x, y)
					m.v[x] /= 2
					if (m.v[x] >> 6) == 1 {
						m.v[15] = 1
					} else if (m.v[x] >> 6) != 1 {
						m.v[15] = 0
					}
				} else if n == 7 {
					fmt.Printf("SUBN x:%x y:%x\n", x, y)
					m.v[x] -= m.v[y]
					if m.v[y] > m.v[x] {
						m.v[15] = 1
					} else {
						m.v[15] = 0
					}
				} else if n == 0xe {
					fmt.Printf("SHL x:%x y:%x\n", x, y)
					m.v[x] *= 2
					if (m.v[x] >> 6) == 1 {
						m.v[15] = 1
					} else if (m.v[x] >> 6) != 1 {
						m.v[15] = 0
					}
				}
				break
			case 9:
				if n == 0 {
					fmt.Printf("SNE x:%x y:%x\n", x, y)
					if m.v[x] != m.v[y] {
						m.pc += 2
					}
				}
				break
			case 0xa:
				fmt.Printf("LD nnn:%x\n", nnn)
				m.i = nnn
				break
			case 0xb:
				fmt.Printf("JP nnn:%x\n", nnn)
				m.pc = nnn + uint16(m.v[0])
				break
			case 0xc:
				fmt.Printf("RND x:%x kk:%x\n", x, kk)
				rand.Seed(time.Now().UnixNano())
				m.v[x] = uint8(rand.Intn(256)) & uint8(kk)
				break
			case 0xd:
				fmt.Printf("DRW x:%x y:%x n:%x\n", x, y, n) // DRAW
				//Draw(x, y)
				break
			case 0xe:
				if (op & 0x00ff) == 0x9e {
					fmt.Printf("SKP x:%x\n", x)
					//TECLADO
					/*if v[x] == m.keys {
						pos+1
					}*/
					m.pc += 2
				} else if (op & 0x00ff) == 0xa1 {
					fmt.Printf("SKNP x:%x\n", x)
					//TECLADO
					/*if v[x] == m.keys {
						pos-1
					}*/
					m.pc += 2
				}
				break
			case 0xf:
				if kk == 0x07 {
					fmt.Printf("LD x:%x\n", x)
					m.v[x] = m.dt
				} else if kk == 0x0a {
					fmt.Printf("LD x:%x\n", x)
					/*if m.keys == press {
						v[x] = m.keys
					}*/
				} else if kk == 0x15 {
					fmt.Printf("LD DT x:%x\n", x)
					m.dt = m.v[x]
				} else if kk == 0x18 {
					fmt.Printf("LD ST x:%x\n", x)
					m.v[x] = m.st
				} else if kk == 0x1e {
					fmt.Printf("ADD I x:%x\n", x)
					m.i += uint16(m.v[x])
				} else if kk == 0x29 {
					fmt.Printf("LD F x:%x\n", x)
					m.i = uint16(m.v[x]) * 0x5
					//DRAW SPRITES
				} else if kk == 0x33 {
					fmt.Printf("LD B x:%x\n", x)
					m.mem[m.i] = m.v[x] / 100
					m.mem[m.i+1] = (m.v[x] / 100) % 10
					m.mem[m.i+2] = (m.v[x] % 100) / 10
					// REVISAR
				} else if kk == 0x55 {
					fmt.Printf("LD [I] x:%x\n", x)
					m.v[0] = m.v[x]
					m.pc += 2 // REVISAR
				} else if kk == 0x65 {
					fmt.Printf("LD x:%x\n", x)
					m.v[0] = m.v[x]
					m.pc += 2 // REVISAR
				}
				break
			default:
				mustQuit = true
				break
			}
		}
	}

	if m.st > 0 {
		m.st -= 1
	}

	if m.dt > 0 {
		m.dt -= 1
	}

	return nil
}
