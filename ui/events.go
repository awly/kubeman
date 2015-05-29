package ui

import "github.com/GoogleCloudPlatform/kubernetes/pkg/watch"

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
	t, ok := ui.tabs[e.Resource]
	if !ok {
		ui.log.Println("unsupported resource type", e.Resource)
		return
	}
	t.update(e)
	ui.RedrawBody()
}
