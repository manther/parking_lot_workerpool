package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	parkinglot "github.com/manther/parking_lot_multithreaded/parking_lot"
)

func main() {
	// Dec / instantiation of parking lot and supervisor
	pl := parkinglot.NewParparkinglot(50, 50)
	sp := parkinglot.NewLotSupervisor(&pl, 8)

	// Welcome message
	fmt.Println("Hello, welcome to the parking lot.")
	fmt.Printf("The lot has %d spots.\n", pl.Total())
	fmt.Printf("There are %d compact spots and %d regular spots.\n", pl.CompactTot(), pl.RegularTot())

	// Compact cars prompts
	var CompactNum int
	fmt.Println("Would you like to park any compact cars? (Y/N)")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	input = strings.TrimSuffix(input, "\n")
	input = strings.ToLower(input)
	if input == "y" {
		fmt.Println("How many compact cars would you like to park? (int)")
		input, err = reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		input = strings.TrimSuffix(input, "\n")
		CompactNum = stringtoint(input)
		fmt.Println("We will park this many compact cars:", CompactNum)
	}

	// Regular cars prompts
	var regNum int
	fmt.Println("Would you like to park any regular cars? (Y/N)")
	input, err = reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	input = strings.TrimSuffix(input, "\n")
	input = strings.ToLower(input)
	if input == "y" {
		fmt.Println("How many regular cars would you like to park? (int)")
		input, err = reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		input = strings.TrimSuffix(input, "\n")
		regNum = stringtoint(input)
		fmt.Println("We will park this many regular cars:", regNum)
	}

	// App Driver
	// Check compact number
	if CompactNum > 1 {
		errs := sp.ManageGroupParkComp(CompactNum)
		for _, err := range errs {
			if err != nil {
				panic(err)
			}
		}
		fmt.Printf("%d compact cars parked.\n", CompactNum)
	}
	if CompactNum == 1 {
		err := sp.ManageSingleParkComp()
		if err != nil {
			panic(err)
		}
		fmt.Println("1 compact car parked.")
	}

	fmt.Println("regnum", regNum)
	if regNum > 1 {
		errs := sp.ManageGroupParkReg(regNum)
		for _, err := range errs {
			if err != nil {
				panic(err)
			}
		}
		fmt.Printf("%d regular cars parked.\n", regNum)
	}
	if regNum == 1 {
		err := sp.ManageSingleParkReg()
		if err != nil {
			panic(err)
		}
		fmt.Println("1 reg car parked.")
	}

	// Report results
	fmt.Printf("The parking lot currently has %d compact spots, and %d regular spots available\n", pl.CompactAvailable(), pl.RegularAvailable())
}

// No real error handeling in this simple main driver
// just a mechanism to show the parking lot in action
func stringtoint(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
