package main

import (
	"github.com/lamoda/clilocker/internal/command"
	config2 "github.com/lamoda/clilocker/internal/config"
	services2 "github.com/lamoda/clilocker/internal/services"
	"log"
	"os"
	"os/exec"
)

func main() {

	config, err := config2.New(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	services, err := services2.New(config)
	if err != nil {
		log.Fatal(err)
	}

	cmd, err := command.New(config, services)
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}

		log.Fatal(err)
	}
}
