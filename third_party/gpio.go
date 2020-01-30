package main

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
)

type transmitOptions struct {
	PulseLength uint
	GpioPin     uint
	Protocol    int
}

func main() {

	options := &transmitOptions{
		PulseLength: gpio.DefaultPulseLength,
		GpioPin:     gpio.DefaultReceivePin,
		Protocol:    gpio.DefaultProtocol,
	}

	//const char* code[6] = { "00000", "10000", "01000", "00100", "00010", "00001" };
	codeOn := []byte{0000, 0000, 0001, 0101, 0001, 0101, 0101, 0100}

	off := uint64(00000000000101010001010001010100) // hex: 00 15 14 54, dec: 0 21 20 84
	on := uint64(00000000000101010001010101010100)  // hex: 00 15 15 54, dec: 0 21 21 84

	data := binary.BigEndian.Uint64(codeOn)
	fmt.Println(data)
	fmt.Println(strconv.FormatUint(data, 2))

	fmt.Println(off)
	fmt.Println(on)

	code := on

	t, err := gpio.NewTransmitter(options.GpioPin)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer t.Close()

	fmt.Printf("transmitting code=%d pulseLength=%d protocol=%d\n", code, options.PulseLength, options.Protocol)

	// returns no error if you have no gpio pin -.-
	err = t.Transmit(code, options.Protocol, options.PulseLength)
	if err != nil {
		fmt.Println(err)
		return
	}

	t.Wait()
}
