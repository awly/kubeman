package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alytvynov/kubeman/ui"
)

func main() {
	logf, err := os.OpenFile(filepath.Join(os.Getenv("HOME"), ".kubeman.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logf.Close()
	log.SetOutput(logf)
	log.SetFlags(log.Ltime | log.Lshortfile)

	ui, err := ui.New(log.New(logf, "ui: ", log.Ltime|log.Lshortfile))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ui.Close()

	select {
	case <-ui.ExitCh():
		log.Println("exit")
		return
	}
}
