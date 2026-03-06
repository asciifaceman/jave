package main

import (
	"flag"
	"fmt"
)

func main() {
	showVersion := flag.Bool("version", false, "print javec version")
	flag.Parse()

	if *showVersion {
		fmt.Println("javec v0.1.0-bootstrap")
		return
	}

	fmt.Println("javec: compiler bootstrap stub")
	fmt.Println("usage: javec [--version] <input.jave>")
}
