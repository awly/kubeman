package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alytvynov/kubeman/client"
	"github.com/alytvynov/kubeman/ui"
)

func main() {
	c, err := client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	logf, err := os.OpenFile(filepath.Join(os.Getenv("HOME"), ".kubeman.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logf.Close()
	log.SetOutput(logf)
	log.SetFlags(log.Ltime | log.Lshortfile)

	u, err := ui.New(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer u.Close()

	<-u.ExitCh()
}
