package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/NicholasGSwan/gbemu-go/context"
)

func main() {

	fmt.Println("This is the main function of my attempt at a gameboy emulator")

	//file, err := os.Open("Tetris (JUE) (V1.1) [!].gb")
	data, err := os.ReadFile("Tetris (JUE) (V1.1) [!].gb")
	check(err)

	//cart := context.ParseNewCart(data)
	cart := context.OpenRom("Tetris (JUE) (V1.1) [!].gb")
	//err = binary.Read(file, binary.LittleEndian, &val)
	ninLogo := make([]byte, 48)
	title := make([]byte, 16)
	for i := range len(ninLogo) {
		ninLogo[i] = data[i+0x104]
	}
	for i := range len(title) {
		title[i] = data[i+0x134]
	}
	fmt.Println(hex.EncodeToString(ninLogo))
	fmt.Println("printing from data directly")
	fmt.Printf("The title of the rom is: %s \n", string(title))
	fmt.Printf("The cartridge type is: %v \n", data[0x147])
	fmt.Printf("The rom size is: %v \n", data[0x148])
	fmt.Printf("The old licensee code is: %v \n", data[0x14b])

	fmt.Println("Printing from cart directly")
	context.PrintCartHeader(&cart)

	check(err)
	// for _, val := range b1 {
	// 	fmt.Printf("%d bytes were read. first byte: %v \n", n1, val)
	// }

	//fmt.Println("this is the data: " + string(data))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
