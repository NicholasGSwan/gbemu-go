package context

import (
	"fmt"
	"os"
)

type Emu struct {
	IsRunning bool
	IsPaused  bool
	Ticks     uint64
}

type Cart struct {
	filename string
	//romsize    uint32
	//byte slice of all the data in the rom file
	romdata    []byte
	headerData header
}

type header struct {
	//0100-0103
	entrypoint [4]byte
	//0104-0133
	ninLogo [48]byte
	//0134-0143, contains title, and later the manufacturers code as well as the color mode flag for later titles
	titleManCodeCGBFlag [16]byte
	//0144-145, only meaningful if old licensee is $33
	newLicenseeCd string
	//146 determines whether or not cart has SuperGameboy funcitonality
	sgbFlag byte
	//147 figure out how to take this byte and pull the type.  ENUM?
	cartType cartType
	//148 maybe this should be an enum? separate type?
	romSize romSize
	//149 also maybe enum.  If cartridge type is not one of the ram types or mbc2, set to 0
	ramSize ramSize
	//14A
	destCode byte
	//14b
	oldLicenseeCd byte
}

type cartType byte

const (
	ROM_ONLY                       cartType = 0x00
	MBC1                           cartType = 0x01
	MBC1_RAM                       cartType = 0x02
	MBC1_RAM_BATTERY               cartType = 0x03
	MBC2                           cartType = 0x05
	MBC2_BATTERY                   cartType = 0x06
	ROM_RAM                        cartType = 0x08
	ROM_RAM_BATTERY                cartType = 0x09
	MMM01                          cartType = 0x0B
	MMM01_RAM                      cartType = 0x0C
	MMM01_RAM_BATTERY              cartType = 0xD
	MBC3_TIMER_BATTERY             cartType = 0x0F
	MBC3_TIMER_RAM_BATTERY         cartType = 0x10
	MBC3                           cartType = 0x11
	MBC3_RAM                       cartType = 0x12
	MBC3_RAM_BATTERY               cartType = 0x13
	MBC5                           cartType = 0x19
	MBC5_RAM                       cartType = 0x1A
	MBC5_RAM_BATTERY               cartType = 0x1B
	MBC5_RUMBLE                    cartType = 0x1C
	MBC5_RUMBLE_RAM                cartType = 0x1D
	MBC5_RUMBLE_RAM_BATTERY        cartType = 0x1E
	MBC6                           cartType = 0x20
	MBC7_SENSOR_RUMBLE_RAM_BATTERY cartType = 0x22
	POCKET_CAMERA                  cartType = 0xFC
	BANDAI_TAMA5                   cartType = 0xFD
	HuC3                           cartType = 0xFE
	HuC1_RAM_BATTERY               cartType = 0xFF
)

var cartTypeName = map[cartType]string{
	ROM_ONLY:                       "ROM_ONLY",
	MBC1:                           "MBC1",
	MBC1_RAM:                       "MBC1_RAM",
	MBC1_RAM_BATTERY:               "MBC1_RAM_BATTERY",
	MBC2:                           "MBC2",
	MBC2_BATTERY:                   "MBC2_BATTERY",
	ROM_RAM:                        "ROM_RAM",
	ROM_RAM_BATTERY:                "ROM_RAM_BATTERY",
	MMM01:                          "MMM01",
	MMM01_RAM:                      "MMM01_RAM",
	MMM01_RAM_BATTERY:              "MMM01_RAM_BATTERY",
	MBC3_TIMER_BATTERY:             "MBC3_TIMER_BATTERY",
	MBC3_TIMER_RAM_BATTERY:         "MBC3_TIMER_RAM_BATTERY",
	MBC3:                           "MBC3",
	MBC3_RAM:                       "MBC3_RAM",
	MBC3_RAM_BATTERY:               "MBC3_RAM_BATTERY",
	MBC5:                           "MBC5",
	MBC5_RAM:                       "MBC5_RAM",
	MBC5_RAM_BATTERY:               "MBC5_RAM_BATTERY",
	MBC5_RUMBLE:                    "MBC5_RUMBLE",
	MBC5_RUMBLE_RAM:                "MBC5_RUMBLE_RAM ",
	MBC5_RUMBLE_RAM_BATTERY:        "MBC5_RUMBLE_RAM_BATTERY",
	MBC6:                           "MBC6",
	MBC7_SENSOR_RUMBLE_RAM_BATTERY: "MBC7_SENSOR_RUMBLE_RAM_BATTERY",
	POCKET_CAMERA:                  "POCKET_CAMERA",
	BANDAI_TAMA5:                   "BANDAI_TAMA5",
	HuC3:                           "HuC3",
	HuC1_RAM_BATTERY:               "HuC1_RAM_BATTERY",
}

type romSize byte

const (
	mem32KiB romSize = iota
	mem64KiB
	mem128KiB
	mem256KiB
	mem512KiB
	mem1MiB
	mem2MiB
	mem4MiB
	mem8MiB
	mem1_1MiB romSize = 0x52
	mem1_2MiB romSize = 0x53
	mem1_5MiB romSize = 0x54
)

var romSizeText = map[romSize]string{
	mem32KiB:  "32 kilobytes",
	mem64KiB:  "64 kilobytes",
	mem128KiB: "128 kilobytes",
	mem256KiB: "256 kilobytes",
	mem512KiB: "512 kilobytes",
	mem1MiB:   "1 megabyte",
	mem2MiB:   "2 megabytes",
	mem4MiB:   "4 megabytes",
	mem8MiB:   "8 megabytes",
	mem1_1MiB: "1.1 megabytes",
	mem1_2MiB: "1.2 megabytes",
	mem1_5MiB: "1.5 megabytes",
}

type ramSize byte

const (
	noRAM ramSize = iota
	Unused
	ram8KiB
	ram32KiB
	ram128KiB
	ram64KiB
)

var ramSizeText = map[ramSize]string{
	noRAM:     "No RAM",
	Unused:    "Unused",
	ram8KiB:   "8 KiB",
	ram32KiB:  "32 KiB",
	ram128KiB: "128 KiB",
	ram64KiB:  "64 KiB",
}

var destCodeMap = map[byte]string{
	0: "Japan",
	1: "Overseas only",
}

func OpenRom(fileName string) Cart {
	cart := Cart{}
	data, err := os.ReadFile(fileName)

	if err != nil {
		fmt.Println("Could not open file!")
		return cart
	}
	cart.filename = fileName
	cart.romdata = data
	parseHeaderData(&cart)

	return cart
}

func parseHeaderData(cart *Cart) {
	cart.headerData.entrypoint = [4]byte(cart.romdata[0x100:0x104])
	cart.headerData.ninLogo = [48]byte(cart.romdata[0x104:0x134])
	cart.headerData.titleManCodeCGBFlag = [16]byte(cart.romdata[0x134:0x144])
	cart.headerData.newLicenseeCd = string(cart.romdata[0x144:0x146])
	cart.headerData.sgbFlag = cart.romdata[0x146]
	cart.headerData.cartType = cartType(cart.romdata[0x147])
	cart.headerData.romSize = romSize(cart.romdata[0x148])
	cart.headerData.ramSize = ramSize(cart.romdata[0x149])
	cart.headerData.destCode = cart.romdata[0x14a]
	cart.headerData.oldLicenseeCd = cart.romdata[0x14b]
}

// not really using this one anymore.
func ParseNewCart(data []byte) Cart {
	cart := Cart{}
	cart.headerData.ninLogo = [48]byte(data[0x104:0x134])
	cart.headerData.titleManCodeCGBFlag = [16]byte(data[0x134:0x144])

	cart.headerData.cartType = cartType(data[0x147])
	cart.headerData.romSize = romSize(data[0x148])
	cart.headerData.oldLicenseeCd = data[0x14b]
	return cart
}

func PrintCartHeader(cart *Cart) {
	fmt.Printf("The title of the rom is: %s \n", string(cart.headerData.titleManCodeCGBFlag[:]))
	fmt.Printf("The cartridge type is: %v \n", cartTypeName[cart.headerData.cartType])
	fmt.Printf("The rom size is: %v \n", romSizeText[cart.headerData.romSize])
	fmt.Printf("The ram size is: %v \n", ramSizeText[cart.headerData.ramSize])
	fmt.Printf("The destination code is: %s \n", destCodeMap[cart.headerData.destCode])
	fmt.Printf("The old licensee code is: %v \n", cart.headerData.oldLicenseeCd)
	fmt.Printf("The checksum is valid: %v \n", performChecksum(cart))

}

func performChecksum(cart *Cart) bool {
	var checksum byte = 0
	for address := 0x0134; address <= 0x014c; address++ {
		checksum = checksum - cart.romdata[address] - 1
	}
	return checksum == cart.romdata[0x14d]
}

func RunEmulator(fileName string) {

	cart := OpenRom(fileName)
	fmt.Printf("Opening rom: %s", cart.headerData.titleManCodeCGBFlag)

	var emu Emu
	emu.IsRunning = true
	emu.IsPaused = false
	emu.Ticks = 0

	fmt.Println("Emulator is running")
}
