package main

import (
	"flag"
	"fmt"
	"net/http"

	"ninelives/internal/stego"
	"ninelives/internal/version"
)

func main() {
	url := flag.String("url", "", "URL of the meme to download")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Commit: %s\n", version.Commit)
		fmt.Printf("BuildID: %s\n", version.BuildID)
		return
	}

	if *url == "" {
		flag.Usage()
		return
	}

	// Download the meme from the URL
	resp, err := http.Get(*url)
	if err != nil {
		fmt.Println("Error downloading meme:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response code:", resp.StatusCode)
		return
	}

	payload, err := stego.ExtractLSB(resp.Body)
	if err != nil {
		fmt.Printf("%v: %s\n", err, *url)
		return
	}
	fmt.Printf("Extracted %s from %s\n", payload, *url)
}
