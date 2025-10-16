package context

import "fmt"

type Cart struct {
	filename   string
	romsize    uint32
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
	//0144-145
	newLicenseeCd [2]byte
	//146
	sgbFlag byte
	//147 figure out how to take this byte and pull the type.  ENUM?
	cartType cartType
	//148 maybe this should be an enum? separate type?
	romSize byte
	//149 also maybe enum.  If cartridge type is not one of the ram types or mbc2, set to 0
	ramSize byte
	//14A
	destCode byte
}

type cartType byte

//type romSize byte
//type ramSize byte

const (
	ROM_ONLY         cartType = 0x00
	MBC1             cartType = 0x01
	MBC1_RAM         cartType = 0x02
	MBC1_RAM_BATTERY cartType = 0x03
	MBC2             cartType = 0x05
	MBC2_BATTERY     cartType = 0x06

// $08	ROM_RAM
// $09	ROM_RAM_BATTERY
// $0B	MMM01
// $0C	MMM01_RAM
// $0D	MMM01_RAM_BATTERY
// $0F	MBC3_TIMER_BATTERY
// $10	MBC3_TIMER_RAM_BATTERY
// $11	MBC3
// $12	MBC3_RAM
// $13	MBC3_RAM_BATTERY
// $19	MBC5
// $1A	MBC5_RAM
// $1B	MBC5_RAM_BATTERY
// $1C	MBC5_RUMBLE
// $1D	MBC5_RUMBLE_RAM
// $1E	MBC5_RUMBLE_RAM_BATTERY
// $20	MBC6
// $22	MBC7_SENSOR_RUMBLE_RAM_BATTERY
// $FC	POCKET CAMERA
// $FD	BANDAI TAMA5
// $FE	HuC3
// $FF	HuC1_RAM_BATTERY
)

var cartTypeName = map[cartType]string{
	ROM_ONLY:         "ROM_ONLY",
	MBC1:             "MBC1",
	MBC1_RAM:         "MBC1_RAM",
	MBC1_RAM_BATTERY: "MBC1_RAM_BATTERY",
	MBC2:             "MBC2",
	MBC2_BATTERY:     "MBC2_BATTERY",
}

func ParseNewCart(data []byte) Cart {
	cart := Cart{}
	for i := range len(cart.headerData.ninLogo) {
		cart.headerData.ninLogo[i] = data[i+0x104]
	}
	for i := range len(cart.headerData.titleManCodeCGBFlag) {
		cart.headerData.titleManCodeCGBFlag[i] = data[i+0x134]
	}

	cart.headerData.cartType = cartType(data[0x147])
	return cart
}

func PrintCartHeader(cart *Cart) {
	fmt.Printf("The cartridge type is: %v \n", cartTypeName[cart.headerData.cartType])
	fmt.Printf("The rom size is: %v \n", data[0x148])
	fmt.Printf("The old licensee code is: %v \n", data[0x14b])
}
