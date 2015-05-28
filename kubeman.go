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

	u, err := ui.New(log.New(logf, "ui: ", log.Ltime|log.Lshortfile))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer u.Close()

	pw, err := c.WatchPods()
	if err != nil {
		log.Println(err)
		return
	}

	for {
		select {
		case e, ok := <-pw:
			if !ok {
				log.Println("pod watch closed, reconnecting")
				pw, err = c.WatchPods()
				if err != nil {
					log.Println(err)
					return
				}
				continue
			}
			u.Updates <- ui.Event{
				Resource: "pods",
				Type:     e.Type,
				Data:     e.Object,
			}
		case <-u.ExitCh():
			log.Println("exit")
			return
		}
	}
}
