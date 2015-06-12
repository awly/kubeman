package ui

import (
	"bufio"
	"log"
	"sync"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/alytvynov/termui"
)

type logTab struct {
	ui *UI

	closed chan struct{}
	mu     *sync.Mutex
	lines  []string
	offs   int
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
			lt.offs = len(lt.lines) - lt.height
		}
		lt.mu.Unlock()
	}
	if s.Err() != nil {
		log.Println(s.Err())
	}
	log.Println(pod, cont, "log stream closed")
}

func (lt logTab) dataUpdate(e Event) {}
func (lt *logTab) uiUpdate(e termui.Event) {
	switch e.Key {
	case termui.KeyArrowUp:
		lt.mu.Lock()
		if lt.offs > 0 {
			lt.offs--
		}
		lt.mu.Unlock()
		go lt.ui.redrawBody()
	case termui.KeyArrowDown:
		lt.mu.Lock()
		if lt.offs < len(lt.lines)-lt.height {
			lt.offs++
		}
		lt.mu.Unlock()
		go lt.ui.redrawBody()
	case termui.KeyCtrlD:
		lt.mu.Lock()
		lt.offs += 10
		if lt.offs > len(lt.lines)-lt.height {
			lt.offs = len(lt.lines) - lt.height
		}
		lt.mu.Unlock()
		go lt.ui.redrawBody()
	case termui.KeyCtrlU:
		lt.mu.Lock()
		lt.offs -= 10
		if lt.offs < 0 {
			lt.offs = 0
		}
		lt.mu.Unlock()
		go lt.ui.redrawBody()
	}
	switch e.Ch {
	case 'l':
		go lt.ui.SelectTab("pods")
	}
}

func (lt *logTab) toRows() []*termui.Row {
	lt.mu.Lock()
	defer lt.mu.Unlock()
	rows := make([]*termui.Row, 0, lt.height)
	lines := lt.lines[lt.offs:]
	if len(lines) > lt.height {
		lines = lines[:lt.height]
	}
	for _, l := range lines {
		rows = append(rows, termui.NewRow(termui.NewCol(12, 0, label(l))))
	}
	return rows
}

func (lt *logTab) clean() {
	close(lt.closed)
}
