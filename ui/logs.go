package ui

import (
	"bufio"
	"log"
	"sync"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/gizak/termui"
)

type logTab struct {
	ui *UI

	closed chan struct{}
	mu     *sync.Mutex
	lines  []string
	height int
}

func showLogTab(ui *UI, p api.Pod) tab {
	lt := &logTab{
		ui:     ui,
		closed: make(chan struct{}),
		height: termui.TermHeight() - 2,
		mu:     &sync.Mutex{},
	}
	ui.body = lt
	go lt.stream(p.Name, p.Spec.Containers[0].Name)
	return lt
}

func (lt *logTab) stream(pod, cont string) {
	// initial feedback without waiting for stream to be opened
	lt.ui.redrawBody()

	in, err := lt.ui.api.Logs(pod, cont, true)
	if err != nil {
		log.Println(err)
		return
	}
	defer in.Close()
	log.Println(pod, cont, "log stream opened")
	go func() {
		t := time.NewTicker(100 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				lt.ui.redrawBody()
			case <-lt.closed:
				in.Close()
				return
			}
		}
	}()
	s := bufio.NewScanner(in)
	for s.Scan() {
		lt.mu.Lock()
		lt.lines = append(lt.lines, s.Text())
		if len(lt.lines) > lt.height {
			lt.lines = lt.lines[1:]
		}
		lt.mu.Unlock()
	}
	if s.Err() != nil {
		log.Println(s.Err())
	}
	log.Println(pod, cont, "log stream closed")
}

func (lt logTab) dataUpdate(e Event) {}
func (lt logTab) uiUpdate(e termui.Event) {
	switch e.Ch {
	case 'l':
		go lt.ui.SelectTab("pods")
	}
}

func (lt *logTab) toRows() []*termui.Row {
	lt.mu.Lock()
	defer lt.mu.Unlock()
	rows := make([]*termui.Row, 0, len(lt.lines))
	for _, l := range lt.lines {
		rows = append(rows, termui.NewRow(termui.NewCol(12, 0, label(l))))
	}
	return rows
}

func (lt *logTab) clean() {
	close(lt.closed)
}
