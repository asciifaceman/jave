package main

import (
	"flag"
	"fmt"
)

func main() {
	showVersion := flag.Bool("version", false, "print javevm version")
	flag.Parse()

	if *showVersion {
		fmt.Println("javevm v0.1.0-bootstrap")
		return
	}

	fmt.Println("javevm: runtime bootstrap stub")
	fmt.Println("usage: javevm <program.jbin>")
}
