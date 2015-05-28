package ui

import (
	"sort"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
)

type Event struct {
	Resource string
	Type     watch.EventType
	Data     interface{}
}

func handleUpdate(ui *UI, e Event) {
	ui.log.Printf("%+v", e)
	if e.Type == watch.Error {
		return
	}
	switch e.Resource {
	case "pod":
		p := *e.Data.(*api.Pod)
		ui.log.Printf("event: %v pod: %v", e.Type, p.Name)
		switch e.Type {
		case watch.Added:
			ui.pods.pods = append(ui.pods.pods, pod{p: p})
		case watch.Modified:
			found := false
			for i, up := range ui.pods.pods {
				if up.p.Name == p.Name {
					found = true
					ui.pods.pods[i].p = p
					break
				}
			}
			if !found {
				ui.pods.pods = append(ui.pods.pods, pod{p: p})
			}
		case watch.Deleted:
			for i, up := range ui.pods.pods {
				if up.p.Name == p.Name {
					ui.pods.Swap(i, ui.pods.Len()-1)
					ui.pods.pods = ui.pods.pods[:ui.pods.Len()-1]
					break
				}
			}
		}
		sort.Sort(ui.pods)
		ui.Redraw()
	}
	ui.Redraw()
}
