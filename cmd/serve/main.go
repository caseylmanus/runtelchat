package main

import (
	"fmt"
	"os"

	"github.com/caseylmanus/runtelchat"
)

func main() {
	config, err := runtelchat.ReadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	runtelchat.ServeTCP(config)
}
