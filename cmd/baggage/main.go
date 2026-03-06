package main

import (
	"flag"
	"fmt"
)

func main() {
	showVersion := flag.Bool("version", false, "print baggage version")
	flag.Parse()

	if *showVersion {
		fmt.Println("baggage v0.1.0-bootstrap")
		return
	}

	fmt.Println("baggage: package/build bootstrap stub")
	fmt.Println("usage: baggage <new|build|run|check|test|add>")
}
