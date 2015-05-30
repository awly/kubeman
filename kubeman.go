package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
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

	watchUpdates(c, u)
}

type watcher struct {
	resource string
	c        <-chan watch.Event
	watch    func() (<-chan watch.Event, error)
}

func watchUpdates(c *client.Client, u *ui.UI) {
	watches := []watcher{
		{resource: "pods", watch: c.WatchPods},
		{resource: "services", watch: c.WatchServices},
		{resource: "rcs", watch: c.WatchRCs},
		{resource: "nodes", watch: c.WatchNodes},
	}
	var err error
	for i, w := range watches {
		watches[i].c, err = w.watch()
		if err != nil {
			log.Println(err)
			return
		}
	}

	cases := make([]reflect.SelectCase, 0, len(watches))
	for _, w := range watches {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(w.c),
		})
	}
	cases = append(cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(u.ExitCh()),
	})

	for {
		i, val, ok := reflect.Select(cases)
		// if ui.ExitCh
		if i == len(cases)-1 {
			log.Println("exit")
			return
		}
		w := watches[i]
		if !ok {
			log.Println(w.resource, "watch closed, reconnecting")
			w.c, err = w.watch()
			if err != nil {
				log.Println(err)
				return
			}
			continue
		}
		e := val.Interface().(watch.Event)
		u.Updates <- ui.Event{
			Resource: w.resource,
			Type:     e.Type,
			Data:     e.Object,
		}
	}

}
