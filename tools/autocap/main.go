package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	path := os.Args[1]

	output, err := exec.Command("setcap", "cap_net_raw+ep", path).CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		log.Fatal(err)
	}
}
