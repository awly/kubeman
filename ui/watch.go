package ui

import (
	"log"
	"reflect"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
)

type watcher struct {
	resource string
	c        <-chan watch.Event
	watch    func() (<-chan watch.Event, error)
}

func (u *UI) watchUpdates() {
	watches := []watcher{
		{resource: "pods", watch: u.api.WatchPods},
		{resource: "services", watch: u.api.WatchServices},
		{resource: "rcs", watch: u.api.WatchRCs},
		{resource: "nodes", watch: u.api.WatchNodes},
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

	for {
		i, val, ok := reflect.Select(cases)
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
		handleUpdate(u, Event{
			Resource: w.resource,
			Type:     e.Type,
			Data:     e.Object,
		})
	}

}
