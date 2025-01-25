package main

import (
	"fmt"
	"os"

	"github.com/realsunyz/geofeed-tools/plugin/validate"
)

const version = "1.0.0"

var help = fmt.Sprintf(
	"Geofeed Tools Version %s (compiled to binary)\n"+
		"Usage: ./geofeed-tools [flags] [filepath]\n\n"+
		"Flags:\n"+
		"  -h    Show this help message\n"+
		"  -v    Validate a geofeed file\n", version)

func main() {
	if len(os.Args) < 2 {
		fmt.Print(help)
		os.Exit(1)
	}

	flag := os.Args[1]

	if flag == "-h" {
		fmt.Print(help)
		os.Exit(1)
	} else if flag == "-v" {
		if len(os.Args) != 3 {
			fmt.Println(
				"Usage: ./geofeed-tools -v [filepath]")
			os.Exit(1)
		}
		filePath := os.Args[2]
		validate.Execute(filePath)
	} else {
		fmt.Printf(
			"Error: Invalid flag \"%s\"\n"+
				"Use \"./geofeed-tools -h\" to see available flags.\n", os.Args[1])
		os.Exit(1)
	}
}
