package main

import (
	"fmt"
	"os"

	"github.com/caseylmanus/runtelchat"
)

func main() {
	config, err := runtelchat.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = runtelchat.ServeTCP(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
