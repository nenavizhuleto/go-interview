package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
)

func GetCheckSum(data string, len int) []byte {
	sum := 0
	for i := 0; i < len; i++ {
		b := rune(data[i])
		sum ^= int(b << 8)

		for j := 0; j < 8; j++ {
			sum <<= 1
			if sum&0x10000 != 0 {
				sum = (sum ^ 0x1021) & 0xFFFF
			}
		}

	}

	res := make([]byte, 2)
	binary.LittleEndian.PutUint16(res, uint16(sum))

	return res

}

type Options struct {
	InputFile string
	Length    int
	Verbose   bool
}

func PrintHelp() {
	help := `
Usage: %s -f [FILE] -l [NUMBER] [FLAGS]

Example: %s -f data.bin -l 0x400 -v

Arguments:
file	-f path		File to read
length	-l number	Length of bytes to read

Flags:
verify	-v		Verify checksum

Description:
Calculating checksum of the file from beginning to [length] bytes
If verify flag is set, program will verify calculated checksum
over next 2 bytes after [length]
	`
	fmt.Printf(help, os.Args[0], os.Args[0])
	fmt.Println()
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Printf("Not enough arguments passed.\nUse %s -h for help\n", os.Args[0])
		os.Exit(1)
	}
	opts := Options{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-h":
			PrintHelp()
			os.Exit(0)
		case "-f":
			if i+1 > len(args)-1 {
				fmt.Printf("No input file specified. Exiting...\n")
				PrintHelp()
				os.Exit(1)
			}
			opts.InputFile = args[i+1]
			i++
		case "-l":
			if i+1 > len(args)-1 {
				fmt.Printf("No length specified. Exiting...\n")
				PrintHelp()
				os.Exit(1)
			}
			length, err := strconv.ParseInt(args[i + 1], 0, 64)
			if err != nil {
				fmt.Printf("Error parsing length: %s\n", err.Error())
				os.Exit(1)
			}
			opts.Length = int(length)
			i++
		case "-v":
			opts.Verbose = true
		}
	}

	data, err := os.ReadFile(opts.InputFile)
	if err != nil {
		fmt.Printf("Error reading file %s. Exiting...\n", opts.InputFile)
		os.Exit(1)
	}

	if opts.Length == 0 {
		fmt.Println("Warning: length argument not specified.")
	}

	checksum := GetCheckSum(string(data), opts.Length)

	fmt.Printf("Result: ")
	for _, b := range checksum {
		fmt.Printf("0x%.02x ", b)
	}

	if opts.Verbose {
		validChecksum := data[opts.Length:opts.Length + 2]
		fmt.Printf("\nValid checksum: ")
		for _, b := range validChecksum {
			fmt.Printf("0x%.02x ", b)
		}

		for i := range checksum {
			if checksum[i] != validChecksum[i] {
				fmt.Println("\nFailed to verify")
				break
			}
		}
	}


}
