package main

import (
	"fmt"
	"log"
	"os"

	"github.com/caseylmanus/runtelchat"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	//read config
	config, err := runtelchat.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//setup logging
	log.SetOutput(&lumberjack.Logger{
		Filename:   config.LogFilePath,
		MaxAge:     30,
		MaxBackups: 3,
		MaxSize:    20,
		Compress:   false,
	})
	log.Println("Starting Server")
	//start chat server
	err = runtelchat.ServeTCP(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
