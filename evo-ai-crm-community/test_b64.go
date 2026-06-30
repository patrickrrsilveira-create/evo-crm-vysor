package main

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func main() {
	url := "data:audio/ogg;base64,T2dnUwAC"
	parts := strings.SplitN(url, ",", 2)
	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success! Decoded %d bytes\n", len(decoded))
	}
}
