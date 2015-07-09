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
	u.statusUpdate("connecting")
	watches := []watcher{
		{resource: "pods", watch: u.api.WatchPods},
		{resource: "services", watch: u.api.WatchServices},
		{resource: "rcs", watch: u.api.WatchRCs},
		{resource: "nodes", watch: u.api.WatchNodes},
	}
	var err error
	for i, w := range watches {
		log.Println("creating watch for", w.resource)
		u.statusUpdate("connecting " + w.resource + " watch")
		watches[i].c, err = w.watch()
		if err != nil {
			log.Println(err)
			u.statusUpdate(err.Error())
			//close(u.exitch)
			//return
		}
	}
	u.statusUpdate("connected")

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
			u.statusUpdate("reconnecting " + w.resource + " watch")
			w.c, err = w.watch()
			if err != nil {
				log.Println(err)
				u.statusUpdate(err.Error())
			}
			continue
		}
		e := val.Interface().(watch.Event)
		u.handleUpdate(Event{
			Resource: w.resource,
			Type:     e.Type,
			Data:     e.Object,
		})
	}

}
