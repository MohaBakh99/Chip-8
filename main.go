package main

import "emulador"

func main() {
	err := emulador.load_rom("pong.rom")
	if err != nil {
		panic(err)
	}
}