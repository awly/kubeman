package ui

import (
	"bufio"
	"io"
	"log"
	"time"

	"github.com/gizak/termui"
)

type logTab struct {
	in        io.ReadCloser
	lines     []string
	height    int
	redraw    func()
	uiUpdatef func(termui.Event)
	cleanf    func()
}

func (lt *logTab) stream() {
	defer lt.in.Close()
	go func() {
		t := time.NewTicker(100 * time.Millisecond)
		defer t.Stop()
		for range t.C {
			lt.redraw()
		}
	}()
	s := bufio.NewScanner(lt.in)
	for s.Scan() {
		lt.lines = append(lt.lines, s.Text())
		if len(lt.lines) > lt.height {
			lt.lines = lt.lines[1:]
		}
	}
	if s.Err() != nil {
		log.Println(s.Err())
	}
	log.Println("log source closed")
}

func (lt logTab) dataUpdate(e Event)      {}
func (lt logTab) uiUpdate(e termui.Event) { lt.uiUpdatef(e) }

func (lt logTab) toRows() []*termui.Row {
	rows := make([]*termui.Row, 0, len(lt.lines))
	for _, l := range lt.lines {
		rows = append(rows, termui.NewRow(termui.NewCol(12, 0, label(l))))
	}
	return rows
}

func (lt *logTab) clean() {
	lt.in.Close()
	lt.cleanf()
}
