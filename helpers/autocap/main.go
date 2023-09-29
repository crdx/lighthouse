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
		log.Fatal("Usage: autocap <path>")
	}

	output, err := exec.Command("setcap", "cap_net_raw+eip", os.Args[1]).CombinedOutput()

	if err != nil {
		fmt.Println(string(output))
		log.Fatal(err)
	}
}
