package main

import (
	"os"
)

/*test prog for the external control binary*/
func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	switch os.Args[1] {
	case "0":
	case "1":
		return
	default:
		os.Exit(1)
	}
}
