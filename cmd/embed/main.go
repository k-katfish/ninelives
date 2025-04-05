package main

import (
	"flag"
	"fmt"
	"os"

	"ninelives/internal/stego"
	"ninelives/internal/version"
)

func main() {
	inFile := flag.String("in", "", "Payload file to embed")
	coverFile := flag.String("cover", "", "Cover file to embed")
	outFile := flag.String("out", "", "Output file to write")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Commit: %s\n", version.Commit)
		fmt.Printf("BuildID: %s\n", version.BuildID)
		return
	}

	if *inFile == "" || *outFile == "" || *coverFile == "" {
		flag.Usage()
		return
	}

	payload, err := os.ReadFile(*inFile)
	if err != nil {
		panic(err)
	}

	err = stego.EmbedLSB(*coverFile, *outFile, payload)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Embedded %d bytes into %s\n", len(payload), *outFile)
}
